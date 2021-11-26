package project

import (
	"github.com/giantswarm/versionbundle"
)

func NewVersionBundle() versionbundle.Bundle {
	return versionbundle.Bundle{
		Changelogs: []versionbundle.Changelog{
			{
				Component:   Name(),
				Description: "Initial commit",
				Kind:        versionbundle.KindAdded,
				URLs:        []string{"https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter"},
			},
		},
		Name:    Name(),
		Version: Version(),
	}
}
