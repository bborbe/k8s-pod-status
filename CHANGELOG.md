# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## Unreleased

- fix: Emit a `version` label on `build_info` by migrating to shared `github.com/bborbe/metrics` and dropping the private `pkg/libmetrics` copy. The local copy registered a bare unlabelled `Gauge`, and quant's fleet `BuildStale` rule selects `build_info{version!~"^v[0-9]+[.][0-9]+[.][0-9]+$"}` — PromQL treats an absent label as `""`, which does not match, so every label-less series stayed permanently in scope and re-fired 14 days after each rebuild. `v0.1.1` papered over this by refreshing the build; a `version` label fixes it for good, and tagged releases now self-exempt while untagged builds (`v0.1.4-3-gabc1234`) stay correctly covered. Same fix `sentry-proxy` and `alert-controller` took on 2026-07-25.
- fix: Pass `BUILD_GIT_VERSION` through to the binary. The Dockerfile and Makefile already baked it in, but `main.go` never read it, so the value was discarded at startup. Also wires up the previously unused `BuildGitCommit` field, now emitted as the `commit` label.
- fix: Register the canonical `/gc` admin endpoint, which was missing from the admin HTTP server. Present in `go-skeleton` and 59 other fleet services; used on-call to force garbage collection.
- chore: Adopting the shared metrics lib pulled transitive bumps (prometheus/client_golang v1.23.2→v1.24.1, getsentry/sentry-go v0.47.0→v0.48.0, protobuf, and several bborbe libs).

## v0.1.4

- chore: Bump errcheck to v1.20.0 and golangci-lint to v2.13.1 for Go 1.27 support
## v0.1.3

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
