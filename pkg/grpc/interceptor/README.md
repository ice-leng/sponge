## interceptor

Common interceptors for gRPC server and client side, including:

- Logging
- Recovery
- Retry
- Rate limiter
- Circuit breaker
- Timeout
- Tracing
- Request id
- Metrics
- JWT authentication

<br>

### Example of use

All import paths are "github.com/go-dev-frame/sponge/pkg/grpc/interceptor".

#### Logging interceptor

**gRPC server side**

```go
// set unary server logging
func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption
	
	option := grpc.ChainUnaryInterceptor(
		// if you don't want to log reply data, you can use interceptor.StreamServerSimpleLog instead of interceptor.UnaryServerLog,
		interceptor.UnaryServerLog(
			logger.Get(),
			interceptor.WithReplaceGRPCLogger(),
			//interceptor.WithMarshalFn(fn), // customised marshal function, default is jsonpb.Marshal
			//interceptor.WithLogIgnoreMethods(fullMethodNames), // ignore methods logging
			//interceptor.WithMaxLen(400), // logging max length, default 300
		),
	)
	options = append(options, option)

	return options
}


// you can also set stream server logging
```

**gRPC client side**

```go
// set unary client logging
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	option := grpc.WithChainUnaryInterceptor(
		interceptor.UnaryClientLog(
			logger.Get(),
			interceptor.WithReplaceGRPCLogger(),
		),
	)
	options = append(options, option)

	return options
}

// you can also set stream client logging
```

<br>

#### Recovery interceptor

**gRPC server side**

```go
func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	option := grpc.ChainUnaryInterceptor(
		interceptor.UnaryServerRecovery(),
	)
	options = append(options, option)

	return options
}
```

**gRPC client side**

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	option := grpc.WithChainUnaryInterceptor(
		interceptor.UnaryClientRecovery(),
	)
	options = append(options, option)

	return options
}
```

<br>

#### Retry interceptor

**gRPC client side**

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// use insecure transfer
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// retry
	option := grpc.WithChainUnaryInterceptor(
		interceptor.UnaryClientRetry(
			//middleware.WithRetryTimes(5), // modify the default number of retries to 3 by default
			//middleware.WithRetryInterval(100*time.Millisecond), // modify the default retry interval, default 50 milliseconds
			//middleware.WithRetryErrCodes(), // add trigger retry error code, default is codes.Internal, codes.DeadlineExceeded, codes.Unavailable
		),
	)
	options = append(options, option)

	return options
}
```

<br>

#### Adaptive rate limiter interceptor

**gRPC server side**

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// use insecure transfer
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// rate limiter
	option := grpc.ChainUnaryInterceptor(
		interceptor.UnaryServerRateLimit(
			//interceptor.WithWindow(time.Second*5),
			//interceptor.WithBucket(200),
			//interceptor.WithCPUThreshold(600),
			//interceptor.WithCPUQuota(0),
		),
	)
	options = append(options, option)

	return options
}
```

<br>

#### Adaptive circuit breaker interceptor

**gRPC server side**

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// use insecure transfer
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// circuit breaker
	option := grpc.ChainUnaryInterceptor(
		interceptor.UnaryServerCircuitBreaker(
			//interceptor.WithValidCode(codes.DeadlineExceeded), // add error code 4 for circuit breaker
			//interceptor.WithUnaryServerDegradeHandler(handler), // add custom degrade handler
		),
	)
	options = append(options, option)

	return options
}
```

<br>

#### Timeout interceptor

**gRPC client side**

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// use insecure transfer
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// timeout
	option := grpc.WithChainUnaryInterceptor(
		interceptor.UnaryClientTimeout(time.Second), // set timeout
	)
	options = append(options, option)

	return options
}
```

<br>

#### Tracing interceptor

**gRPC server side**

```go
// initialize tracing
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

// set up trace on the client side
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// use insecure transfer
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// use tracing
	option := grpc.WithUnaryInterceptor(
		interceptor.UnaryClientTracing(),
	)
	options = append(options, option)

	return options
}

// set up trace on the server side
func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// use tracing
	option := grpc.UnaryInterceptor(
		interceptor.UnaryServerTracing(),
	)
	options = append(options, option)

	return options
}

// if necessary, you can create a span in the program
func SpanDemo(serviceName string, spanName string, ctx context.Context) {
	_, span := otel.Tracer(serviceName).Start(
		ctx, spanName,
		trace.WithAttributes(attribute.String(spanName, time.Now().String())), // customised attributes
	)
	defer span.End()

	// ......
}
```

<br>

#### Metrics interceptor

Click to view [metrics examples](../metrics/README.md).

<br>

#### Request id interceptor

**gRPC server side**

```go
func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	option := grpc.ChainUnaryInterceptor(
		interceptor.UnaryServerRequestID(),
	)
	options = append(options, option)

	return options
}
```

<br>

**gRPC client side**

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// use insecure transfer
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	option := grpc.WithChainUnaryInterceptor(
		interceptor.UnaryClientRequestID(),
	)
	options = append(options, option)

	return options
}
```

<br>

#### JWT authentication interceptor

**gRPC server side**

```go
package main

import (
	"context"
	"github.com/go-dev-frame/sponge/pkg/grpc/interceptor"
	"github.com/go-dev-frame/sponge/pkg/jwt"
	"google.golang.org/grpc"
	"net"
	"time"
	userV1 "user/api/user/v1"
)

func main() {
	list, err := net.Listen("tcp", ":8282")
	server := grpc.NewServer(getUnaryServerOptions()...)
	userV1.RegisterUserServer(server, &user{})
	server.Serve(list)
	select {}
}

func getUnaryServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// Case1: default options
	{
		options = append(options, grpc.UnaryInterceptor(
			interceptor.UnaryServerJwtAuth(),
		))
	}

	// Case 2: custom options, signKey, extra verify function, rpc method
	{
		options = append(options, grpc.UnaryInterceptor(
			interceptor.UnaryServerJwtAuth(
				interceptor.WithSignKey([]byte("your_secret_key")),
				interceptor.WithExtraVerify(extraVerifyFn),
				interceptor.WithAuthIgnoreMethods(// specify the gRPC API to ignore token verification(full path)
					"/api.user.v1.User/Register",
					"/api.user.v1.User/Login",
				),
			),
		))
	}

	return options
}

type user struct {
	userV1.UnimplementedUserServer
}

// Login ...
func (s *user) Login(ctx context.Context, req *userV1.LoginRequest) (*userV1.LoginReply, error) {
	// check user and password success

	uid := "100"
	fields := map[string]interface{}{"name":   "bob","age":    10,"is_vip": true}

	// Case 1: default jwt options, signKey, signMethod(HS256), expiry time(24 hour)
	{
		_, token, err := jwt.GenerateToken("100")
	}

	// Case 2: custom jwt options, signKey, signMethod(HS512), expiry time(12 hour), fields, claims
	{
		_, token, err := jwt.GenerateToken(
			uid,
			jwt.WithGenerateTokenSignKey([]byte("your_secret_key")),
			jwt.WithGenerateTokenSignMethod(jwt.HS384),
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
	}

	return &userV1.LoginReply{Token: token}, nil
}

func extraVerifyFn(ctx context.Context, claims *jwt.Claims) error {
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

**gRPC client side**

```go
package main

import (
	"context"
	"github.com/go-dev-frame/sponge/pkg/grpc/grpccli"
	"github.com/go-dev-frame/sponge/pkg/grpc/interceptor"
	userV1 "user/api/user/v1"
)

func main() {
	conn, _ := grpccli.NewClient("127.0.0.1:8282")
	cli := userV1.NewUserClient(conn)

	uid := "100"
	ctx := context.Background()

	// Case 1: get authorization from header key is "authorization", value is "Bearer xxx"
	{
		ctx = interceptor.SetAuthToCtx(ctx, authorization)
	}
	// Case 2: get token from grpc server response result
	{
		ctx = interceptor.SetJwtTokenToCtx(ctx, token)
	}

	cli.GetByID(ctx, &userV1.GetUserByIDRequest{Id: 100})
}

```

<br>
