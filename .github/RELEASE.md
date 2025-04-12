
## Change log

### New Features
1. **Enhanced Code Generation Tool**
    - Added command and UI for generating gRPC + HTTP service code based on SQL
    - Introduced `goast` library for Go source code parsing
    - Added Gemini AI Assistant SDK
    - Added command and UI for the AI Assistant to generate and merge code
    - `make run` command now supports specifying a configuration file

### Refactoring & Optimization
1. **Core Logic Refactoring**
    - Optimized the logic for code generation and merging using the `protoc` plugin
    - Refactored the authentication module:
        - Improved the `pkg/jwt` package
        - Enhanced JWT authentication middleware for the Gin framework (`pkg/gin/middleware/jwtAuth.go`)

### Bug Fixes
1. **Database Related**
    - Fixed an issue where the `sgorm.Bool` type could not properly read or assign PostgreSQL boolean fields
2. **Cross-Platform Compatibility**
    - Resolved an issue where code archives appeared empty when extracted using the built-in tool on Windows
3. **Dependency Management**
    - Fixed a version conflict issue with the `go.opentelemetry.io/otel` dependency [#97](https://github.com/go-dev-frame/sponge/issues/97)
