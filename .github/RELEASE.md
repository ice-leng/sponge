
## Change Log

### New Features

1. Added a custom Copier utility library that supports automatic bidirectional conversion between time types and strings.
2. Implemented automatic conversion from Swagger 2.0 to OpenAPI 3.0 specification.
3. Added implementation of SSE (Server-Sent Events) for both server and client sides.
4. MongoDB ORM supports complex conditional group queries, with automatic type recognition and conversion for values (integer/date-time).
5. SGORM ORM supports automatic type recognition and conversion for values (integer/date-time).

### Bug Fixes

1. Fixed an issue where the Swagger API documentation generated from Protobuf was inconsistent with the actual API response format.

### Dependency Upgrades

1. Upgraded Gin framework from v1.9.1 to v1.10.1.
2. Upgraded Copier library from v0.3.5 to v0.4.0.
