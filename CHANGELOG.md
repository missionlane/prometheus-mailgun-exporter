# Changelog

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
