package main

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/microkit/command"
	microserver "github.com/giantswarm/microkit/server"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/versionbundle"
	"github.com/spf13/viper"

	"github.com/giantswarm/etcd-kubernetes-resources-count-exporter/flag"
	"github.com/giantswarm/etcd-kubernetes-resources-count-exporter/pkg/project"
	"github.com/giantswarm/etcd-kubernetes-resources-count-exporter/server"
	"github.com/giantswarm/etcd-kubernetes-resources-count-exporter/service"
)

var (
	f = flag.New()
)

func main() {
	err := mainError()
	if err != nil {
		panic(fmt.Sprintf("%#v\n", err))
	}
}
func mainError() error {
	var err error

	ctx := context.Background()
	logger, err := micrologger.New(micrologger.Config{})
	if err != nil {
		return microerror.Mask(err)
	}

	// We define a server factory to create the custom server once all command
	// line flags are parsed and all microservice configuration is sorted out.
	serverFactory := func(v *viper.Viper) microserver.Server {
		// Create a new custom service which implements business logic.
		var newService *service.Service
		{
			c := service.Config{
				Flag:   f,
				Logger: logger,
				Viper:  v,

				Description: project.Description(),
				GitCommit:   project.GitSHA(),
				ProjectName: project.Name(),
				Source:      project.Source(),
				Version:     project.Version(),
			}

			newService, err = service.New(c)
			if err != nil {
				panic(fmt.Sprintf("%#v", microerror.Mask(err)))
			}

			go newService.Boot(ctx)
		}

		// Create a new custom server which bundles our endpoints.
		var newServer microserver.Server
		{
			c := server.Config{
				Logger:  logger,
				Service: newService,
				Viper:   v,

				ProjectName: project.Name(),
			}

			newServer, err = server.New(c)
			if err != nil {
				panic(fmt.Sprintf("%#v", microerror.Mask(err)))
			}
		}

		return newServer
	}

	// Create a new microkit command which manages our custom microservice.
	var newCommand command.Command
	{
		c := command.Config{
			Logger:        logger,
			ServerFactory: serverFactory,

			Description:    project.Description(),
			GitCommit:      project.GitSHA(),
			Name:           project.Name(),
			Source:         project.Source(),
			Version:        project.Version(),
			VersionBundles: []versionbundle.Bundle{project.NewVersionBundle()},
		}

		newCommand, err = command.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	daemonCommand := newCommand.DaemonCommand().CobraCommand()

	daemonCommand.PersistentFlags().StringArray(f.Service.Etcd.Endpoints, []string{"https://127.0.0.1:2379"}, "endpoint to connect to etcd. Can be repeated multiple times for multiple endpoints. Default: https://127.0.0.1:2379")
	daemonCommand.PersistentFlags().String(f.Service.Etcd.CaCertPath, "", "path of the CA certificate file for tls validation.")
	daemonCommand.PersistentFlags().String(f.Service.Etcd.KeyPath, "", "path of the client key file for tls authentication.")
	daemonCommand.PersistentFlags().String(f.Service.Etcd.CertPath, "", "path of the client certificate file for tls authentication.")
	daemonCommand.PersistentFlags().Int(f.Service.Etcd.DialTimeout, 10, "dial timeout in seconds for connecting to etcd.")
	daemonCommand.PersistentFlags().String(f.Service.Etcd.Prefix, "", "prefix used to store k8s data in etcd as specified in the '--etcd-prefix' flag of API server.")
	daemonCommand.PersistentFlags().String(f.Service.Events.Prefix, "", "prefix used to store event data in etcd")

	err = newCommand.CobraCommand().Execute()
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
