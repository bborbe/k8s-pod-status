# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v0.1.1

- chore: Refresh build to clear the critical BuildStale alert on quant — the v0.1.0 image was 23 days old; BUILD_DATE and BUILD_GIT_VERSION are baked in at image build time, so only a rebuilt image resets build_info
- chore: Reformat the go.mod exclude directive to block form, as applied by go-modtool fmt in make precommit

## v0.1.0

- Extract k8s-pod-status from the trading monorepo into a dedicated public repo (publish-only image)
- Decouple from trading/lib (kafka-free metrics + loglevel handler vendored into pkg/)
