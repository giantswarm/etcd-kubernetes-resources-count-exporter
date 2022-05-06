package collector

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"go.etcd.io/etcd/clientv3"
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
		},
		nil,
	)
)

type EventsCollector struct {
	logger           micrologger.Logger
	cache            []cachedEvent
	etcdClientConfig *clientv3.Config
	etcdPrefix       string
}

type cachedEvent struct {
	count     float64
	kind      string
	namespace string
}

func NewEventsCollector(config EtcdConfig) (*Etcd, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.EtcdClientConfig == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.EtcdClientConfig must not be empty", config)
	}
	if !strings.HasSuffix(config.EtcdPrefix, "/") || !strings.HasPrefix(config.EtcdPrefix, "/") {
		return nil, microerror.Maskf(invalidConfigError, "%T.EtcdPrefix has to start and end with a '/'", config)
	}

	d := &EventsCollector{
		logger:           config.Logger,
		cache:            make([]cachedEvent, 0),
		etcdClientConfig: config.EtcdClientConfig,
		etcdPrefix:       config.EtcdPrefix,
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
			prometheus.GaugeValue,
			event.count,
			event.namespace,
			event.kind,
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

	resp, err := cli.Get(context.Background(), "/giantswarm.io/events", clientv3.WithPrefix())
	if err != nil {
		return microerror.Mask(err)
	}

	decoder := scheme.Codecs.UniversalDeserializer()
	encoder := jsonserializer.NewSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, true)

	for _, kv := range resp.Kvs {
		obj, _, err := decoder.Decode(kv.Value, nil, nil)
		if err != nil {
			d.logger.Debugf(context.Background(), "WARN: unable to decode %s: %v\n", kv.Key, err)
			continue
		}

		err = encoder.Encode(obj, os.Stdout)
		if err != nil {
			d.logger.Debugf(context.Background(), "WARN: unable to encode %s: %v\n", kv.Key, err)
			continue
		}

		fmt.Println(obj)

	}

	d.cache = newCache

	return nil
}
