package service

import (
	"github.com/giantswarm/etcd-kubernetes-resources-count-exporter/flag/service/etcd"
)

type Service struct {
	Etcd etcd.Etcd
}
