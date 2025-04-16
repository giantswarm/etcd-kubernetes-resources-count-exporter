package collector

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	k8sResourcesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(MetricsNamespace, "resources", "count"),
		"Count of kubernetes resources stored in etcd.",
		[]string{
			"namespace",
			"kind",
		},
		nil,
	)
)

type EtcdConfig struct {
	Logger           micrologger.Logger
	EtcdClientConfig *clientv3.Config
	EtcdPrefix       string
}

type Etcd struct {
	logger           micrologger.Logger
	cache            []cacheEntry
	etcdClientConfig *clientv3.Config
	etcdPrefix       string
}

type cacheEntry struct {
	count     float64
	kind      string
	namespace string
}

// NewEtcd exposes metrics about the number of k8s resources stored in etcd.
func NewEtcd(config EtcdConfig) (*Etcd, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.EtcdClientConfig == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.EtcdClientConfig must not be empty", config)
	}
	if !strings.HasSuffix(config.EtcdPrefix, "/") || !strings.HasPrefix(config.EtcdPrefix, "/") {
		return nil, microerror.Maskf(invalidConfigError, "%T.EtcdPrefix has to start and end with a '/'", config)
	}

	d := &Etcd{
		logger:           config.Logger,
		cache:            make([]cacheEntry, 0),
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

func (d *Etcd) Collect(ch chan<- prometheus.Metric) error {
	for _, entry := range d.cache {
		ch <- prometheus.MustNewConstMetric(
			k8sResourcesDesc,
			prometheus.GaugeValue,
			entry.count,
			entry.namespace,
			entry.kind,
		)
	}

	return nil
}

func (d *Etcd) Describe(ch chan<- *prometheus.Desc) error {
	ch <- k8sResourcesDesc
	return nil
}

func (d *Etcd) refreshCache() error {
	grouped := map[string]map[string]int64{}

	newCache := make([]cacheEntry, 0)

	cli, err := clientv3.New(*d.etcdClientConfig)
	if err != nil {
		return microerror.Mask(err)
	}

	defer cli.Close() //nolint:errcheck

	resp, err := cli.Get(context.Background(), "/", clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return microerror.Mask(err)
	}

	for _, kv := range resp.Kvs {
		line := string(kv.Key)

		line = strings.TrimPrefix(line, d.etcdPrefix)

		namespace, kind, err := parseLine(line)
		if IsEmptyLine(err) {
			continue
		} else if err != nil {
			d.logger.Debugf(context.Background(), "Error parsing line %q: %s", line, err)
			continue
		}

		if _, found := grouped[kind]; !found {
			grouped[kind] = make(map[string]int64)
		}

		grouped[kind][namespace] = grouped[kind][namespace] + 1
	}

	for kind, namespacedCount := range grouped {
		for namespace, count := range namespacedCount {
			newCache = append(newCache, cacheEntry{
				count:     float64(count),
				kind:      kind,
				namespace: namespace,
			})
		}
	}

	d.cache = newCache

	return nil
}

func parseLine(line string) (namespace string, kind string, err error) {
	err = nil

	// Strip etcd prefix
	line = strings.TrimRight(line, "\n")

	if line == "" {
		err = microerror.Maskf(emptyLineError, "Line was empty")
		return
	}

	// cert-manager.io/clusterissuers/letsencrypt-giantswarm

	tokens := strings.Split(line, "/")

	namespace = "Not namespaced"
	kind = tokens[0]

	if strings.Contains(tokens[0], ".") {
		// the first token contains a dot, we consider this an api version of non namespaced objects.
		kind = fmt.Sprintf("%s.%s", tokens[1], tokens[0])

		if len(tokens) > 3 {
			namespace = tokens[2]
		}
		return
	}

	if len(tokens) > 2 {
		namespace = tokens[1]
	}

	return
}
