package collector

import (
	"bufio"
	_ "embed"
	"fmt"
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
)

const etcdPrefix = "/giantswarm.io/"

//go:embed sampledata
var sampledata string

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
	Logger micrologger.Logger
}

type Deployment struct {
	logger micrologger.Logger
}

// NewEtcd exposes metrics about the number of k8s resources stored in etcd.
func NewEtcd(config EtcdConfig) (*Deployment, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	d := &Deployment{
		logger: config.Logger,
	}

	return d, nil
}

func (d *Deployment) Collect(ch chan<- prometheus.Metric) error {
	grouped := map[string]map[string]int64{}

	scanner := bufio.NewScanner(strings.NewReader(sampledata))
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := scanner.Text()
		// Strip etcd prefix
		line = strings.TrimPrefix(line, etcdPrefix)
		line = strings.TrimRight(line, "\n")

		if line == "" {
			continue
		}

		tokens := strings.Split(line, "/")
		namespace := tokens[len(tokens)-2]
		kind := tokens[0]
		if len(tokens) > 3 {
			kind = fmt.Sprintf("%s.%s", tokens[1], kind)
		}

		if _, found := grouped[kind]; !found {
			grouped[kind] = make(map[string]int64)
		}

		grouped[kind][namespace] = grouped[kind][namespace] + 1
	}

	if err := scanner.Err(); err != nil {
		return err
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

func (d *Deployment) Describe(ch chan<- *prometheus.Desc) error {
	ch <- k8sResourcesDesc
	return nil
}
