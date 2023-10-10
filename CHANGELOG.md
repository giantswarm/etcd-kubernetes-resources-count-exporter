# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Replace condition for PSP CR installation.
- Disable debugging.

## [1.6.0] - 2023-10-03

### Changed

- Set VPA max memory to 2Gi.

## [1.5.0] - 2023-09-19

### Changed

- Set `priorityClassName` to the deployment to mitigate scheduling issues.

## [1.4.0] - 2023-08-01

### Changed

- Add Max memory (default 500Mi) for VPA.

## [1.3.0] - 2023-06-23

### Changed

- Configurable path for etcd certificates.
- Enable installation of `etcd-kubernetes-resources-count-exporter` for CAPI clusters.

## [1.2.1] - 2023-06-08

### Removed

- Removed debug code that was dumping all events' contents to stdout.

## [1.2.0] - 2023-05-04

### Changed

- Disable PSPs for k8s 1.25 and newer.

## [1.1.2] - 2023-04-25

### Changed

- Update icon.

## [1.1.1] - 2023-04-20

### Fixed

- Use Port 10999 for listening as the previous value was overlapping with the TCP ephemeral port range.

## [1.1.0] - 2023-04-20

### Added

- Added the use of the runtime/default seccomp profile.

### Fixed

- Use ClusterIP service rather than NodePort one.

## [1.0.0] - 2022-09-23

### Changed

- Push to azure and AWS app collection.
- Run in kube-system by default.
- Change default registry to docker.io.

## [0.5.2] - 2022-09-06

### Changed

- Include namespaces in dotted scopes.

## [0.5.1] - 2022-08-03

### Changed

- Bump build image to alpine 3.16.1.

## [0.5.0] - 2022-07-05

### Changed

- Also push to default catalog.

## [0.4.0] - 2022-05-20

## [0.3.1] - 2022-03-24

## [0.3.0] - 2022-03-24

### Changed

- Collect data asyncronously to shorten scrape time and avoid timeouts.

## [0.2.0] - 2022-03-21

### Added

- Add VerticalPodAutoscaler CR.

## [0.1.0] - 2021-11-29

### Added

- First release. 

[Unreleased]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.6.0...HEAD
[1.6.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.5.0...v1.6.0
[1.5.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.4.0...v1.5.0
[1.4.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.2.1...v1.3.0
[1.2.1]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.2.0...v1.2.1
[1.2.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.1.2...v1.2.0
[1.1.2]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.1.1...v1.1.2
[1.1.1]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.5.2...v1.0.0
[0.5.2]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.5.1...v0.5.2
[0.5.1]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.3.1...v0.4.0
[0.3.1]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.0.0...v0.1.0
