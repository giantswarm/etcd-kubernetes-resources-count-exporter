# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.5.2...v1.0.0
[0.5.2]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.5.1...v0.5.2
[0.5.1]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.3.1...v0.4.0
[0.3.1]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v0.0.0...v0.1.0
