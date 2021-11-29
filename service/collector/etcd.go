package collector

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"go.etcd.io/etcd/clientv3"
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
	etcdClientConfig *clientv3.Config
	etcdPrefix       string
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
		etcdClientConfig: config.EtcdClientConfig,
		etcdPrefix:       config.EtcdPrefix,
	}

	return d, nil
}

func (d *Etcd) Collect(ch chan<- prometheus.Metric) error {
	grouped := map[string]map[string]int64{}

	cli, err := clientv3.New(*d.etcdClientConfig)
	if err != nil {
		return microerror.Mask(err)
	}

	defer cli.Close()

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
			ch <- prometheus.MustNewConstMetric(
				k8sResourcesDesc,
				prometheus.GaugeValue,
				float64(count),
				namespace,
				kind,
			)
		}
	}

	return nil
}

func (d *Etcd) Describe(ch chan<- *prometheus.Desc) error {
	ch <- k8sResourcesDesc
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

	if strings.Contains(tokens[0], ".") {
		// the first token contains a dot, we consider this an api version.
		kind = fmt.Sprintf("%s.%s", tokens[1], tokens[0])
		return
	}
	if len(tokens) > 2 {
		namespace = tokens[len(tokens)-2]
	}
	kind = tokens[0]
	if len(tokens) > 3 {
		kind = fmt.Sprintf("%s.%s", tokens[1], kind)
	}

	return
}
