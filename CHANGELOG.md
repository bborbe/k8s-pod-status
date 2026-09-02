# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## Unreleased

- chore: update github.com/bborbe/errors to v1.6.0, github.com/bborbe/http to v1.26.25, github.com/bborbe/k8s to v1.14.16, github.com/bborbe/log to v1.6.25, github.com/bborbe/run to v1.10.1, github.com/bborbe/sentry to v1.10.0, github.com/bborbe/service to v1.10.10, github.com/bborbe/time to v1.27.11, github.com/onsi/gomega to v1.43.0

## v0.1.7

- fix: add the missing `github.com/bborbe/metrics` dependency to `go.mod` and resync `go.sum`. The v0.1.5 migration to the shared metrics library landed the import without its module requirement, so `make check-go-mod` failed with 14 missing `go.sum` entries — on the `v0.1.6` tag *and* on master. The repo could not build at all, and because nothing publishes an image on release, no CI or weekly rebuild ever surfaced it. Found during the 2026-08-30 nuke rebuild while trying to publish `v0.1.6`. `go mod tidy` also brought several stale indirect pins current (`bborbe/errors` v1.5.21, `bborbe/collection` v1.20.24, `bborbe/math` v1.4.7, `bborbe/validation` v1.4.23, `getsentry/sentry-go` v0.49.0) and normalized the single-entry `exclude` block.

## v0.1.6

- chore: update Go to 1.27.0 and update dependencies

## v0.1.5

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
