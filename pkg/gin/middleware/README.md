## middleware

Common gin middleware libraries.

<br>

## Example of use

### Logging middleware

You can set the maximum length for printing, add a request id field, ignore print path, customize [zap](go.uber.org/zap) log.

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/go-dev-frame/sponge/pkg/gin/middleware"
)

func NewRouter() *gin.Engine {
    r := gin.Default()
    // ......

    // Print input parameters and return results
    // Case 1: default
    {
        r.Use(middleware.Logging()
    }
    // Case 2: custom
    {
        r.Use(middleware.Logging(
            middleware.WithLog(logger.Get()))
            middleware.WithMaxLen(400),
            middleware.WithRequestIDFromHeader(),
            //middleware.WithRequestIDFromContext(),
            //middleware.WithIgnoreRoutes("/hello"),
        ))
    }

    /*******************************************
    TIP: You can use middleware.SimpleLog instead of
           middleware.Logging, it only prints return results
    *******************************************/

    // ......
    return r
}
```

<br>

### Allow cross-domain requests middleware

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/go-dev-frame/sponge/pkg/gin/middleware"
)

func NewRouter() *gin.Engine {
    r := gin.Default()
    // ......

    r.Use(middleware.Cors())

    // ......
    return r
}
```

<br>

### Rate limiter middleware

Adaptive flow limitation based on hardware resources.

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/go-dev-frame/sponge/pkg/gin/middleware"
)

func NewRouter() *gin.Engine {
    r := gin.Default()
    // ......

    // Case 1: default
    r.Use(middleware.RateLimit())

    // Case 2: custom
    r.Use(middleware.RateLimit(
        middleware.WithWindow(time.Second*10),
        middleware.WithBucket(1000),
        middleware.WithCPUThreshold(100),
        middleware.WithCPUQuota(0.5),
    ))

    // ......
    return r
}
```

<br>

### Circuit Breaker middleware

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/go-dev-frame/sponge/pkg/gin/middleware"
)

func NewRouter() *gin.Engine {
    r := gin.Default()
    // ......

    r.Use(middleware.CircuitBreaker(
        //middleware.WithValidCode(http.StatusRequestTimeout), // add error code 408 for circuit breaker
        //middleware.WithDegradeHandler(handler), // add custom degrade handler
    ))

    // ......
    return r
}
```

<br>

### JWT authorization middleware

```go
package main

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/go-dev-frame/sponge/pkg/gin/middleware"
    "github.com/go-dev-frame/sponge/pkg/gin/response"
    "github.com/go-dev-frame/sponge/pkg/jwt"
)

func main() {
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
import (
    "github.com/gin-gonic/gin"
    "github.com/go-dev-frame/sponge/pkg/gin/middleware"
    "github.com/go-dev-frame/sponge/pkg/tracer"
)

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

func NewRouter() *gin.Engine {
    r := gin.Default()
    // ......

    r.Use(middleware.Tracing("your-service-name"))

    // ......
    return r
}

// if necessary, you can create a span in the program
func CreateSpanDemo(serviceName string, spanName string, ctx context.Context) {
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
import (
    "github.com/gin-gonic/gin"
    "github.com/go-dev-frame/sponge/pkg/gin/middleware"
    "github.com/go-dev-frame/sponge/pkg/gin/middleware/metrics"
)

func NewRouter() *gin.Engine {
    r := gin.Default()
    // ......

    r.Use(metrics.Metrics(r,
        //metrics.WithMetricsPath("/demo/metrics"), // default is /metrics
        metrics.WithIgnoreStatusCodes(http.StatusNotFound), // ignore status codes
        //metrics.WithIgnoreRequestMethods(http.MethodHead),  // ignore request methods
        //metrics.WithIgnoreRequestPaths("/ping", "/health"), // ignore request paths
    ))

    // ......
    return r
```

<br>

### Request id

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/go-dev-frame/sponge/pkg/gin/middleware"
)

func NewRouter() *gin.Engine {
    r := gin.Default()
    // ......

    // Case 1: default request id
    {
        r.Use(middleware.RequestID())
    }
    // Case 2: custom request id key
    {
        //r.User(middleware.RequestID(
        //    middleware.WithContextRequestIDKey("your ctx request id key"), // default is request_id
        //    middleware.WithHeaderRequestIDKey("your header request id key"), // default is X-Request-Id
        //))
        // If you change the ContextRequestIDKey, you have to set the same key name if you want to print the request id in the mysql logs as well.
        // example:
        //     db, err := mysql.Init(dsn,mysql.WithLogRequestIDKey("your ctx request id key"))  // print request_id
    }

    // ......
    return r
}
```

<br>

### Timeout

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/go-dev-frame/sponge/pkg/gin/middleware"
)

func NewRouter() *gin.Engine {
    r := gin.Default()
    // ......

    // Case 1: global set timeout
    {
        r.Use(middleware.Timeout(time.Second*5))
    }
    // Case 2: set timeout for specifyed router
    {
        r.GET("/userExample/:id", middleware.Timeout(time.Second*3), GetByID)
    }
    // Note: If timeout is set both globally and in the router, the minimum timeout prevails

    // ......
    return r
}

func GetByID(c *gin.Context) {
    // do something
}
```
