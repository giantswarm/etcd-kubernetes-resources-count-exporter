package project

var (
	description string = "The etcd-kubernetes-resources-count-exporter exposes metrics about the number of k8s objects stored in etcd."
	gitSHA             = "n/a"
	name        string = "etcd-kubernetes-resources-count-exporter"
	source      string = "https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter"
	version            = "1.10.14"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
