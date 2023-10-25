FROM alpine:3.18.4

RUN apk add --no-cache ca-certificates

ADD ./etcd-kubernetes-resources-count-exporter /etcd-kubernetes-resources-count-exporter

ENTRYPOINT ["/etcd-kubernetes-resources-count-exporter"]
