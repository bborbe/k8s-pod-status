# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## Unreleased

- chore: Ignore the generated /vendor/ dir — `make build` vendors it via check-go-mod (matching the trading repo's Makefile.docker), `make precommit` removes it again via ensure, and it is never tracked; leaving it unignored made `git add .` liable to commit the whole vendor tree
- chore: Update Go from 1.26.5 to 1.26.6 (go.mod, Dockerfile, CI) and golang.org/x/mod from v0.37.0 to v0.40.0 to clear the vulncheck failures blocking CI

## v0.1.2

- docs: add a License section to the README

## v0.1.1

- chore: Refresh build to clear the critical BuildStale alert on quant — the v0.1.0 image was 23 days old; BUILD_DATE and BUILD_GIT_VERSION are baked in at image build time, so only a rebuilt image resets build_info
- chore: Reformat the go.mod exclude directive to block form, as applied by go-modtool fmt in make precommit

## v0.1.0

- Extract k8s-pod-status from the trading monorepo into a dedicated public repo (publish-only image)
- Decouple from trading/lib (kafka-free metrics + loglevel handler vendored into pkg/)
