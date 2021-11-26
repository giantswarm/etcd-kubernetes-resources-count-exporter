package collector

import (
	"github.com/giantswarm/exporterkit/collector"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

const (
	MetricsNamespace = "etcd_kubernetes_resources_count"
)

type SetConfig struct {
	Logger micrologger.Logger
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
		c := DeploymentConfig{
			Logger: config.Logger,
		}

		etcdCollector, err := NewEtcd(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		collectors = append(collectors, etcdCollector)
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
