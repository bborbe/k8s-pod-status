# Changelog

## Unreleased

- Refresh build to clear the critical `BuildStale` alert on quant (v0.1.0 image was 23 days old).
  No source changes — `BUILD_DATE` / `BUILD_GIT_VERSION` are baked in at image build time, so a
  rebuilt and re-tagged image is what resets `build_info`.

## v0.1.0

- Extract k8s-pod-status from the trading monorepo into a dedicated public repo (publish-only image)
- Decouple from trading/lib (kafka-free metrics + loglevel handler vendored into pkg/)
