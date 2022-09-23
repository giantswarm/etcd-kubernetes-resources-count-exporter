package service

import (
	"github.com/giantswarm/etcd-kubernetes-resources-count-exporter/v2/flag/service/etcd"
	"github.com/giantswarm/etcd-kubernetes-resources-count-exporter/v2/flag/service/events"
)

type Service struct {
	Etcd   etcd.Etcd
	Events events.Events
}
