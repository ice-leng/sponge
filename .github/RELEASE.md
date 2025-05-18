
## Change log

1. Removed the custom binding implementation and standardized the use of Gin's default binding mechanism.
2. Added a lightweight Gin-JWT middleware implementation to simplify the authentication process.
3. Standardized variable naming conventions in the generated code to ensure consistency, with special handling for proper nouns.
4. Deprecated the custom `$neq` operator in MongoDB queries to maintain consistency with native query syntax.
5. Added a whitelist validation mechanism for the `name` field in custom query APIs to effectively prevent SQL injection risks.
