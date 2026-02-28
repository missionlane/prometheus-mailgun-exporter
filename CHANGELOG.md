# Changelog

## [1.0.1](https://github.com/missionlane/prometheus-mailgun-exporter/compare/v1.0.0...v1.0.1) (2026-02-28)


### Bug Fixes

* **deps:** update module github.com/mailgun/mailgun-go/v5 to v5.10.0 ([#68](https://github.com/missionlane/prometheus-mailgun-exporter/issues/68)) ([599cb28](https://github.com/missionlane/prometheus-mailgun-exporter/commit/599cb28b1dfef3b6a5de3019874a7eaf669b5b8d))
* **deps:** update module github.com/mailgun/mailgun-go/v5 to v5.10.1 ([#71](https://github.com/missionlane/prometheus-mailgun-exporter/issues/71)) ([47acbf7](https://github.com/missionlane/prometheus-mailgun-exporter/commit/47acbf7449711a0be11e8dfb9700fed00df094bb))
* **deps:** update module github.com/mailgun/mailgun-go/v5 to v5.11.0 ([#77](https://github.com/missionlane/prometheus-mailgun-exporter/issues/77)) ([0ae86a6](https://github.com/missionlane/prometheus-mailgun-exporter/commit/0ae86a6406d74f7c0b32c5ddbe246ede1882bcd9))
* **deps:** update module github.com/mailgun/mailgun-go/v5 to v5.12.0 ([#82](https://github.com/missionlane/prometheus-mailgun-exporter/issues/82)) ([5c2217e](https://github.com/missionlane/prometheus-mailgun-exporter/commit/5c2217ed99295b63caabc4fe719275bd478f0b05))
* **deps:** update module github.com/mailgun/mailgun-go/v5 to v5.13.0 ([#87](https://github.com/missionlane/prometheus-mailgun-exporter/issues/87)) ([91f91ee](https://github.com/missionlane/prometheus-mailgun-exporter/commit/91f91eef418aca691f1e55fc867de836c7a474a8))
* **deps:** update module github.com/mailgun/mailgun-go/v5 to v5.13.1 ([#90](https://github.com/missionlane/prometheus-mailgun-exporter/issues/90)) ([fae72ce](https://github.com/missionlane/prometheus-mailgun-exporter/commit/fae72ce3f67151101b292212afccbc0250ed156f))
* **deps:** update module github.com/mailgun/mailgun-go/v5 to v5.13.2 ([#91](https://github.com/missionlane/prometheus-mailgun-exporter/issues/91)) ([ce486a0](https://github.com/missionlane/prometheus-mailgun-exporter/commit/ce486a0c0b073ee0e73ec966c46bb8b3d7cf5fee))
* **deps:** update module github.com/mailgun/mailgun-go/v5 to v5.14.0 ([#93](https://github.com/missionlane/prometheus-mailgun-exporter/issues/93)) ([ce8b6d1](https://github.com/missionlane/prometheus-mailgun-exporter/commit/ce8b6d1b9069a508ebe402e7b22260f749ad2768))
* **deps:** update module github.com/mailgun/mailgun-go/v5 to v5.9.1 ([#66](https://github.com/missionlane/prometheus-mailgun-exporter/issues/66)) ([a635371](https://github.com/missionlane/prometheus-mailgun-exporter/commit/a635371855bba18a557ddd0c18e25f8c5dbe391a))
* **deps:** update module github.com/prometheus/common to v0.67.5 ([#67](https://github.com/missionlane/prometheus-mailgun-exporter/issues/67)) ([c66bd1a](https://github.com/missionlane/prometheus-mailgun-exporter/commit/c66bd1a128193769738373ac250c7c1586ba6574))

## [0.1.0](https://github.com/missionlane/prometheus-mailgun-exporter/compare/v0.0.7...v0.1.0) (2025-12-22)


### ⚠ BREAKING CHANGES

* MG_DOMAIN environment variable is no longer required. The exporter now uses the new ListMetrics API instead of the deprecated GetStats API. All Prometheus metrics remain unchanged.

### Features

* upgrade mailgun-go from v3 to v5 ([2e386e1](https://github.com/missionlane/prometheus-mailgun-exporter/commit/2e386e16c0bffe5b06fe57f6279243e7e55ebe51))


### Bug Fixes

* **deps:** update golang deps ([c5ff76b](https://github.com/missionlane/prometheus-mailgun-exporter/commit/c5ff76b301e5a8bfcd49e5d9ca92b2c5ea42d911))
* **deps:** update golang deps ([b7a3d01](https://github.com/missionlane/prometheus-mailgun-exporter/commit/b7a3d0146a7ce9dcf2f404e6cf706166e5a868e5))
* **deps:** update module github.com/prometheus/client_golang to v1.23.0 ([acbac11](https://github.com/missionlane/prometheus-mailgun-exporter/commit/acbac11bfddb44406f47457e5b31cca9363f0137))
* **deps:** update module github.com/prometheus/common to v0.64.0 ([3d0bc35](https://github.com/missionlane/prometheus-mailgun-exporter/commit/3d0bc35cc3e9e86e38fdbb00127e836cda833f90))
* **deps:** update module github.com/prometheus/common to v0.65.0 ([18154cd](https://github.com/missionlane/prometheus-mailgun-exporter/commit/18154cd62d45543b98f204d58a491a36ea89546b))
* **deps:** update module github.com/prometheus/common to v0.67.3 ([#46](https://github.com/missionlane/prometheus-mailgun-exporter/issues/46)) ([59c3fb0](https://github.com/missionlane/prometheus-mailgun-exporter/commit/59c3fb0836b2193371647cf4743728b1c98afca5))
* **deps:** update module github.com/prometheus/common to v0.67.4 ([#50](https://github.com/missionlane/prometheus-mailgun-exporter/issues/50)) ([ed35bdf](https://github.com/missionlane/prometheus-mailgun-exporter/commit/ed35bdf77e62cb981238e993707202c433417014))
* public only build ([f24acb4](https://github.com/missionlane/prometheus-mailgun-exporter/commit/f24acb425facf53e3593e9cd7e8ba82eb849d0e9))
* public only build ([fbe22a7](https://github.com/missionlane/prometheus-mailgun-exporter/commit/fbe22a7f4d6576332d5595be1eacc575ea7609d5))
