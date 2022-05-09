package service

import (
	"context"
	"sync"
	"time"

	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/versionbundle"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"

	"github.com/giantswarm/etcd-kubernetes-resources-count-exporter/flag"
	"github.com/giantswarm/etcd-kubernetes-resources-count-exporter/pkg/project"
	"github.com/giantswarm/etcd-kubernetes-resources-count-exporter/service/collector"
)

// Config represents the configuration used to create a new service.
type Config struct {
	Logger micrologger.Logger

	Flag  *flag.Flag
	Viper *viper.Viper

	Description string
	GitCommit   string
	ProjectName string
	Source      string
	Version     string
}

type Service struct {
	Version *version.Service

	bootOnce          sync.Once
	operatorCollector *collector.Set
}

// New creates a new configured service object.
func New(config Config) (*Service, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flag must not be empty", config)
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Viper must not be empty", config)
	}
	if config.Description == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Description must not be empty", config)
	}
	if config.GitCommit == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.GitCommit must not be empty", config)
	}
	if config.ProjectName == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.ProjectName must not be empty", config)
	}
	if config.Source == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Source must not be empty", config)
	}

	// Validate flags.
	if len(config.Viper.GetStringSlice(config.Flag.Service.Etcd.Endpoints)) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "Config.Service.Etcd.Endpoints must not be empty")
	}
	if config.Viper.GetString(config.Flag.Service.Etcd.CaCertPath) == "" {
		return nil, microerror.Maskf(invalidConfigError, "Config.Service.Etcd.CaCertPath must not be empty")
	}
	if config.Viper.GetString(config.Flag.Service.Etcd.CertPath) == "" {
		return nil, microerror.Maskf(invalidConfigError, "Config.Service.Etcd.CertPath must not be empty")
	}
	if config.Viper.GetString(config.Flag.Service.Etcd.KeyPath) == "" {
		return nil, microerror.Maskf(invalidConfigError, "Config.Service.Etcd.KeyPath must not be empty")
	}

	var err error

	var operatorCollector *collector.Set
	{
		tlsInfo := transport.TLSInfo{
			TrustedCAFile: config.Viper.GetString(config.Flag.Service.Etcd.CaCertPath),
			CertFile:      config.Viper.GetString(config.Flag.Service.Etcd.CertPath),
			KeyFile:       config.Viper.GetString(config.Flag.Service.Etcd.KeyPath),
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return nil, microerror.Mask(err)
		}

		c := collector.SetConfig{
			Logger: config.Logger,
			EtcdClientConfig: &clientv3.Config{
				Endpoints:   config.Viper.GetStringSlice(config.Flag.Service.Etcd.Endpoints),
				DialTimeout: time.Second * time.Duration(config.Viper.GetInt(config.Flag.Service.Etcd.DialTimeout)),
				TLS:         tlsConfig,
			},
			EtcdPrefix:   config.Viper.GetString(config.Flag.Service.Etcd.Prefix),
			EventsPrefix: config.Viper.GetString(config.Flag.Service.Events.Prefix),
		}

		operatorCollector, err = collector.NewSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var versionService *version.Service
	{
		c := version.Config{
			Description:    config.Description,
			GitCommit:      config.GitCommit,
			Name:           config.ProjectName,
			Source:         config.Source,
			Version:        config.Version,
			VersionBundles: []versionbundle.Bundle{project.NewVersionBundle()},
		}

		versionService, err = version.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	s := &Service{
		Version: versionService,

		bootOnce:          sync.Once{},
		operatorCollector: operatorCollector,
	}

	return s, nil
}

func (s *Service) Boot(ctx context.Context) {
	s.bootOnce.Do(func() {
		go s.operatorCollector.Boot(ctx) // nolint: errcheck
	})
}
