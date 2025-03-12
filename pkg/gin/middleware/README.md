## middleware

Common gin middleware libraries.

<br>

## Example of use

### Logging middleware

You can set the maximum length for printing, add a request id field, ignore print path, customize [zap](go.uber.org/zap) log.

```go
    import "github.com/go-dev-frame/sponge/pkg/gin/middleware"

    r := gin.Default()

    // Print input parameters and return results
    // case 1:
    r.Use(middleware.Logging()) // default
    // case 2:
    r.Use(middleware.Logging( // custom
        middleware.WithMaxLen(400),
        middleware.WithRequestIDFromHeader(),
        //middleware.WithRequestIDFromContext(),
        //middleware.WithLog(log), // custom zap log
        //middleware.WithIgnoreRoutes("/hello"),
    ))    

    // ----------------------------------------

    // Print only return results
    // case 1:
    r.Use(middleware.SimpleLog()) // default
    // case 2:
    r.Use(middleware.SimpleLog( // custom
        middleware.WithRequestIDFromHeader(),
        //middleware.WithRequestIDFromContext(),
        //middleware.WithLog(log), // custom zap log
        //middleware.WithIgnoreRoutes("/hello"),
    ))    
```

<br>

### Allow cross-domain requests middleware

```go
    import "github.com/go-dev-frame/sponge/pkg/gin/middleware"

    r := gin.Default()
    r.Use(middleware.Cors())
```

<br>

### Rate limiter middleware

Adaptive flow limitation based on hardware resources.

```go
    import "github.com/go-dev-frame/sponge/pkg/gin/middleware"

    r := gin.Default()

    // case 1: default
    r.Use(middleware.RateLimit())

    // case 2: custom
    r.Use(middleware.RateLimit(
        middleware.WithWindow(time.Second*10),
        middleware.WithBucket(1000),
        middleware.WithCPUThreshold(100),
        middleware.WithCPUQuota(0.5),
    ))
```

<br>

### Circuit Breaker middleware

```go
    import "github.com/go-dev-frame/sponge/pkg/gin/middleware"

    r := gin.Default()
    r.Use(middleware.CircuitBreaker(
        //middleware.WithValidCode(http.StatusRequestTimeout), // add error code 408 for circuit breaker
        //middleware.WithDegradeHandler(handler), // add custom degrade handler
    ))
```

<br>

### JWT authorization middleware

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/go-dev-frame/sponge/pkg/gin/middleware"
    "github.com/go-dev-frame/sponge/pkg/gin/response"
    "github.com/go-dev-frame/sponge/pkg/jwt"
    "time"
)

func web() {
    r := gin.Default()

    // Case 1: default jwt options, signKey, signMethod(HS256), expiry time(24 hour)
    {
        r.POST("/auth/login", LoginDefault)
        r.GET("/demo1/user/:id", middleware.Auth(), GetByID)
        r.GET("/demo2/user/:id", middleware.Auth(middleware.WithReturnErrReason()), GetByID)
        r.GET("/demo3/user/:id", middleware.Auth(middleware.WithExtraVerify(extraVerifyFn)), GetByID)
    }

    // Case 2: custom jwt options, signKey, signMethod(HS512), expiry time(12 hour), fields, claims
    {
        signKey := []byte("custom-sign-key")
        jwtAuth1 := middleware.Auth(middleware.WithSignKey(signKey))
        jwtAuth2 := middleware.Auth(middleware.WithSignKey(signKey), middleware.WithReturnErrReason())
        jwtAuth3 := middleware.Auth(middleware.WithSignKey(signKey), middleware.WithExtraVerify(extraVerifyFn))

        r.POST("/auth/login", LoginCustom)
        r.GET("/demo4/user/:id", jwtAuth1, GetByID)
        r.GET("/demo5/user/:id", jwtAuth2, GetByID)
        r.GET("/demo6/user/:id", jwtAuth3, GetByID)
    }

    r.Run(":8080")
}

func LoginDefault(c *gin.Context) {
    // ......

    _, token, err := jwt.GenerateToken("100")

    response.Success(c, token)
}

func LoginCustom(c *gin.Context) {
    // ......
 
    uid := "100"
    fields := map[string]interface{}{
        "name":   "bob",
        "age":    10,
        "is_vip": true,
    }

    _, token, err := jwt.GenerateToken(
        uid,
        jwt.WithGenerateTokenSignKey([]byte("custom-sign-key")),
        jwt.WithGenerateTokenSignMethod(jwt.HS512),
        jwt.WithGenerateTokenFields(fields),
        jwt.WithGenerateTokenClaims([]jwt.RegisteredClaimsOption{
            jwt.WithExpires(time.Hour * 12),
            //jwt.WithIssuedAt(now),
            // jwt.WithSubject("123"),
            // jwt.WithIssuer("https://auth.example.com"),
            // jwt.WithAudience("https://api.example.com"),
            // jwt.WithNotBefore(now),
            // jwt.WithJwtID("abc1234xxx"),
        }...),
    )

    response.Success(c, token)
}

func GetByID(c *gin.Context) {
    uid := c.MustGet("id").(string)

    response.Success(c, gin.H{"id": uid})
}

func extraVerifyFn(claims *jwt.Claims, c *gin.Context) error {
    // judge whether the user is disabled, query whether jwt id exists from the blacklist
    //if CheckBlackList(uid, claims.ID) {
    //    return errors.New("user is disabled")
    //}

    // get fields from claims
    //uid := claims.UID
    //name, _ := claims.GetString("name")
    //age, _ := claims.GetInt("age")
    //isVip, _ := claims.GetBool("is_vip")

    return nil
}
```

<br>

### Tracing middleware

```go
import "github.com/go-dev-frame/sponge/pkg/tracer"
import "github.com/go-dev-frame/sponge/pkg/gin/middleware"

func InitTrace(serviceName string) {
    exporter, err := tracer.NewJaegerAgentExporter("192.168.3.37", "6831")
    if err != nil {
        panic(err)
    }

    resource := tracer.NewResource(
        tracer.WithServiceName(serviceName),
        tracer.WithEnvironment("dev"),
        tracer.WithServiceVersion("demo"),
    )

    tracer.Init(exporter, resource) // collect all by default
}

func NewRouter(
    r := gin.Default()
    r.Use(middleware.Tracing("your-service-name"))

    // ......
)

// if necessary, you can create a span in the program
func SpanDemo(serviceName string, spanName string, ctx context.Context) {
    _, span := otel.Tracer(serviceName).Start(
        ctx, spanName,
        trace.WithAttributes(attribute.String(spanName, time.Now().String())),
    )
    defer span.End()

    // ......
}
```

<br>

### Metrics middleware

```go
    import "github.com/go-dev-frame/sponge/pkg/gin/middleware/metrics"

    r := gin.Default()

    r.Use(metrics.Metrics(r,
        //metrics.WithMetricsPath("/demo/metrics"), // default is /metrics
        metrics.WithIgnoreStatusCodes(http.StatusNotFound), // ignore status codes
        //metrics.WithIgnoreRequestMethods(http.MethodHead),  // ignore request methods
        //metrics.WithIgnoreRequestPaths("/ping", "/health"), // ignore request paths
    ))
```

<br>

### Request id

```go
    import "github.com/go-dev-frame/sponge/pkg/gin/middleware"

    // Default request id
    r := gin.Default()
    r.Use(middleware.RequestID())

    // --- or ---

    // Customized request id key
    //r.User(middleware.RequestID(
    //    middleware.WithContextRequestIDKey("your ctx request id key"), // default is request_id
    //    middleware.WithHeaderRequestIDKey("your header request id key"), // default is X-Request-Id
    //))
    // If you change the ContextRequestIDKey, you have to set the same key name if you want to print the request id in the mysql logs as well.
    // example: 
    // db, err := mysql.Init(dsn,
        // mysql.WithLogRequestIDKey("your ctx request id key"),  // print request_id
        // ...
    // )
```

<br>

### Timeout

```go
    import "github.com/go-dev-frame/sponge/pkg/gin/middleware"

    r := gin.Default()

    // way1: global set timeout
    r.Use(middleware.Timeout(time.Second*5))

    // --- or ---

    // way2: router set timeout
    r.GET("/userExample/:id", middleware.Timeout(time.Second*3), h.GetByID)

    // Note: If timeout is set both globally and in the router, the minimum timeout prevails
```
