FROM alpine:3.24.3

RUN apk add --no-cache ca-certificates

ARG TARGETARCH
COPY etcd-kubernetes-resources-count-exporter-linux-${TARGETARCH} /etcd-kubernetes-resources-count-exporter

ENTRYPOINT ["/etcd-kubernetes-resources-count-exporter"]
