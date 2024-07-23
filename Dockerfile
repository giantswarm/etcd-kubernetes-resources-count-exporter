FROM alpine:3.20.2

RUN apk add --no-cache ca-certificates

ADD ./etcd-kubernetes-resources-count-exporter /etcd-kubernetes-resources-count-exporter

ENTRYPOINT ["/etcd-kubernetes-resources-count-exporter"]
