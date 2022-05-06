package collector

import (
	"github.com/giantswarm/exporterkit/collector"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"go.etcd.io/etcd/clientv3"
)

const (
	MetricsNamespace = "etcd_kubernetes"
)

type SetConfig struct {
	Logger           micrologger.Logger
	EtcdClientConfig *clientv3.Config
	EtcdPrefix       string
}

// Set is basically only a wrapper for the operator's collector implementations.
// It eases the iniitialization and prevents some weird import mess so we do not
// have to alias packages.
type Set struct {
	*collector.Set
}

func NewSet(config SetConfig) (*Set, error) {
	var err error
	var collectors []collector.Interface

	{
		c := EtcdConfig{ //nolint
			Logger:           config.Logger,
			EtcdClientConfig: config.EtcdClientConfig,
			EtcdPrefix:       config.EtcdPrefix,
		}

		etcdCollector, err := NewEtcd(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		eventCollector, err := NewEventsCollector(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		collectors = append(collectors, etcdCollector)
		collectors = append(collectors, eventCollector)
	}

	var collectorSet *collector.Set
	{
		c := collector.SetConfig{
			Collectors: collectors,
			Logger:     config.Logger,
		}

		collectorSet, err = collector.NewSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	s := &Set{
		Set: collectorSet,
	}

	return s, nil
}
