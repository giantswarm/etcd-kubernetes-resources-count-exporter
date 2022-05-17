package collector

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"go.etcd.io/etcd/clientv3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/kubectl/pkg/scheme"
)

var (
	eventsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(MetricsNamespace, "events", "count"),
		"Count of events recorded in etcd.",
		[]string{
			"namespace",
			"kind",
			"objectNamespace",
			"reason",
			"source",
		},
		nil,
	)
)

type EventsCollectorConfig struct {
	Logger           micrologger.Logger
	EtcdClientConfig *clientv3.Config
	EventsPrefix     string
}

type state struct {
	uids map[string]string
	keys map[string]float64
}

type EventsCollector struct {
	logger           micrologger.Logger
	cache            []cachedEvent
	state            state
	etcdClientConfig *clientv3.Config
	eventsPrefix     string
}

type cachedEvent struct {
	count           float64
	kind            string
	namespace       string
	objectNamespace string
	reason          string
	source          string
}

func NewEventsCollector(config EventsCollectorConfig) (*EventsCollector, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.EtcdClientConfig == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.EtcdClientConfig must not be empty", config)
	}
	if !strings.HasSuffix(config.EventsPrefix, "/") || !strings.HasPrefix(config.EventsPrefix, "/") {
		return nil, microerror.Maskf(invalidConfigError, "%T.EventsPrefix has to start and end with a '/'", config)
	}

	emptyState := state{
		uids: map[string]string{},
		keys: map[string]float64{},
	}

	d := &EventsCollector{
		logger:           config.Logger,
		cache:            make([]cachedEvent, 0),
		state:            emptyState,
		etcdClientConfig: config.EtcdClientConfig,
		eventsPrefix:     config.EventsPrefix,
	}

	go func() {
		for {
			err := d.refreshCache()

			if err != nil {
				d.logger.Errorf(context.Background(), err, "Error refreshing cache")
			}

			time.Sleep(30 * time.Second)
		}
	}()

	return d, nil
}

func (d *EventsCollector) Collect(ch chan<- prometheus.Metric) error {
	for _, event := range d.cache {
		ch <- prometheus.MustNewConstMetric(
			eventsDesc,
			prometheus.CounterValue,
			event.count,
			event.namespace,
			event.kind,
			event.objectNamespace,
			event.reason,
			event.source,
		)
	}

	return nil
}

func (d *EventsCollector) Describe(ch chan<- *prometheus.Desc) error {
	ch <- eventsDesc
	return nil
}

func (d *EventsCollector) refreshCache() error {
	newCache := make([]cachedEvent, 0)

	cli, err := clientv3.New(*d.etcdClientConfig)
	if err != nil {
		return microerror.Mask(err)
	}

	defer cli.Close()

	resp, err := cli.Get(context.Background(), d.eventsPrefix, clientv3.WithPrefix())

	if err != nil {
		return microerror.Mask(err)
	}

	decoder := scheme.Codecs.UniversalDeserializer()
	encoder := jsonserializer.NewSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, true)

	for _, kv := range resp.Kvs {
		event := d.getEventFromResponse(kv, decoder, encoder)

		d.logger.Debugf(context.Background(), "Event", event)

		cachedEventObj := cachedEvent{
			count:           float64(event.Count),
			kind:            event.InvolvedObject.Kind,
			namespace:       event.Namespace,
			objectNamespace: event.InvolvedObject.Namespace,
			reason:          event.Reason,
			source:          event.Source.Component,
		}

		eventKey, err := getKey(cachedEventObj)

		if err != nil {
			d.logger.Debugf(context.Background(), "WARN: unable to get event metric key %s: %v\n", kv.Key, err)
		}

		//We keep track of events that have been counted.
		//We don't want to recounted existing events
		if _, exists := d.state.uids[string(kv.Key)]; exists {
			continue
		}

		d.state.uids[string(kv.Key)] = event.ObjectMeta.Name

		if count, exists := d.state.keys[eventKey]; exists {
			// We've already counted an event metric like this before. So we add to get the total
			count = count + float64(event.Count)

			cachedEventObj.count = count
			d.state.keys[eventKey] = count

			newCache = append(newCache, cachedEventObj)
			continue
		}

		d.state.keys[eventKey] = cachedEventObj.count
		newCache = append(newCache, cachedEventObj)

		d.logger.Debugf(context.Background(), "EventCounts", d.state)
	}

	d.cache = newCache

	return nil
}

func getKey(event cachedEvent) (string, error) {
	event.count = 0
	out, err := json.Marshal(event)

	return string(out), err
}

func (d *EventsCollector) getEventFromResponse(kv *mvccpb.KeyValue, decoder runtime.Decoder, encoder runtime.Encoder) corev1.Event {
	obj, _, err := decoder.Decode(kv.Value, nil, nil)
	if err != nil {
		d.logger.Debugf(context.Background(), "WARN: unable to decode %s: %v\n", kv.Key, err)
	}

	err = encoder.Encode(obj, os.Stdout)
	if err != nil {
		d.logger.Debugf(context.Background(), "WARN: unable to encode %s: %v\n", kv.Key, err)
	}

	marshalledObj, _ := json.Marshal(obj)
	event := &corev1.Event{}
	_ = json.Unmarshal([]byte(marshalledObj), event)

	return *event
}
