package service

import (
	"github.com/giantswarm/etcd-kubernetes-resources-count-exporter/flag/service/etcd"
	"github.com/giantswarm/etcd-kubernetes-resources-count-exporter/flag/service/events"
)

type Service struct {
	Etcd   etcd.Etcd
	Events events.Events
}
