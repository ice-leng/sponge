
## Changelog

### New Features

1. Added `pkg/sasynq` - a distributed task queue library for asynchronous background task processing.

### Bug Fixes

1. Fixed DBResolver compatibility issues with the latest GORM version [#111](https://github.com/go-dev-frame/sponge/issues/111)
2. Resolved postgresql table parsing errors when table comments contain single quotes [#112](https://github.com/go-dev-frame/sponge/issues/112)

### Dependency Updates

1. Upgraded GORM from v1.25.1 to v1.30.0

### Other Changes

1. Removed deprecated `pkg/ggorm` library (replaced by `pkg/sgorm`)
