package collector

import (
	"context"
	"encoding/json"
	"fmt"
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

type EventsCollector struct {
	logger           micrologger.Logger
	cache            map[string]cachedEvent
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

	d := &EventsCollector{
		logger:           config.Logger,
		cache:            make(map[string]cachedEvent),
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
	newCache := make(map[string]cachedEvent)

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

		cachedEventObj := cachedEvent{
			count:           float64(event.Count),
			kind:            event.InvolvedObject.Kind,
			namespace:       event.Namespace,
			objectNamespace: event.InvolvedObject.Namespace,
			reason:          event.Reason,
			source:          event.Source.Component,
		}

		eventKey := getKey(cachedEventObj)

		if e, exists := newCache[eventKey]; exists {
			cachedEventObj.count = e.count + cachedEventObj.count
		}

		newCache[eventKey] = cachedEventObj

		d.logger.Debugf(context.Background(), "Cached Event", cachedEventObj)
	}

	d.cache = newCache

	return nil
}

func getKey(event cachedEvent) string {
	eventKey := fmt.Sprint(
		event.namespace,
		event.kind,
		event.objectNamespace,
		event.reason,
		event.source,
	)

	return eventKey
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
