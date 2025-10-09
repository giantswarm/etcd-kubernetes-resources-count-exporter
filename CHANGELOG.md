# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.10.9] - 2025-10-09

### Changed

- Go: Update dependencies.

## [1.10.8] - 2025-10-02

### Changed

- Update Kyverno API to v2 for policy exceptions
- Go: Update dependencies.

## [1.10.7] - 2025-08-26

### Changed

- Go: Update dependencies.

## [1.10.6] - 2025-07-31

### Changed

- Go: Update dependencies.

## [1.10.5] - 2025-06-07

### Changed

- Go: Update dependencies.

## [1.10.4] - 2025-06-03

### Fixed

- Fix linting issues.
- Go: Update dependencies.

## [1.10.3] - 2025-03-18

### Changed

- Go: Update dependencies.

## [1.10.2] - 2025-03-17

### Changed

- Go: Update dependencies.

## [1.10.1] - 2025-02-17

### Changed 

- Set `readOnlyRootFilesystem` to true in the container security context.
- Update Kyverno `PolicyExceptions` to `v2beta1`.
- Go: Update `go.mod` and `.nancy-ignore`. ([#242](https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/pull/242))

## [1.10.0] - 2024-04-01

### Changed 

- Set min VPA settings and adjust CPU and memory resources.
- Use PodMonitor instead of legacy labels for monitoring.

## [1.9.0] - 2024-01-18

### Changed

- Configure `gsoci.azurecr.io` as the default container image registry.

## [1.8.0] - 2023-11-02

### Changed

- Avoid logging all events to stdout.

## [1.7.0] - 2023-10-10

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

[Unreleased]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.10.9...HEAD
[1.10.9]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.10.8...v1.10.9
[1.10.8]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.10.7...v1.10.8
[1.10.7]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.10.6...v1.10.7
[1.10.6]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.10.5...v1.10.6
[1.10.5]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.10.4...v1.10.5
[1.10.4]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.10.3...v1.10.4
[1.10.3]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.10.2...v1.10.3
[1.10.2]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.10.1...v1.10.2
[1.10.1]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.10.0...v1.10.1
[1.10.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.9.0...v1.10.0
[1.9.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.8.0...v1.9.0
[1.8.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.7.0...v1.8.0
[1.7.0]: https://github.com/giantswarm/etcd-kubernetes-resources-count-exporter/compare/v1.6.0...v1.7.0
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
