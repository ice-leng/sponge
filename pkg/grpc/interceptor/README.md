## interceptor

Common interceptors for client-side and server-side gRPC.

<br>

### Example of use

```go
import "github.com/go-dev-frame/sponge/pkg/grpc/interceptor"
```

#### logging

**server-side gRPC**

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

**client-side gRPC**

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

#### recovery

**server-side gRPC**

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

**client-side gRPC**

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

#### retry

**client-side gRPC**

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

#### rate limiter

**server-side gRPC**

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

#### Circuit Breaker

**server-side gRPC**

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

#### timeout

**client-side gRPC**

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

#### tracing

**server-side gRPC**

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

#### metrics

example [metrics](../metrics/README.md).

<br>

#### Request id

**server-side gRPC**

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

**client-side gRPC**

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

#### jwt authentication

JWT supports two verification methods:

- The default verification method includes fixed fields `uid` and `name` in the claim, and supports additional custom verification functions.
- The custom verification method allows users to define the claim themselves and also supports additional custom verification functions.

**client-side gRPC**

```go
package main

import (
	"context"
	"github.com/go-dev-frame/sponge/pkg/jwt"
	"github.com/go-dev-frame/sponge/pkg/grpc/interceptor"
	"github.com/go-dev-frame/sponge/pkg/grpc/grpccli"
	userV1 "user/api/user/v1"
)

func main() {
	ctx := context.Background()
	conn, _ := grpccli.NewClient("127.0.0.1:8282")
	cli := userV1.NewUserClient(conn)

	token := "xxxxxx" // no Bearer prefix
	ctx = interceptor.SetJwtTokenToCtx(ctx, token)

	req := &userV1.GetUserByIDRequest{Id: 100}
	cli.GetByID(ctx, req)
}
```

**server-side gRPC**

```go
package main

import (
	"context"
	"net"
	"github.com/go-dev-frame/sponge/pkg/jwt"
	"github.com/go-dev-frame/sponge/pkg/grpc/interceptor"
	"google.golang.org/grpc"
	userV1 "user/api/user/v1"
)

func main()  {
	list, err := net.Listen("tcp", ":8282")
	server := grpc.NewServer(getUnaryServerOptions()...)
	userV1.RegisterUserServer(server, &user{})
	server.Serve(list)
	select{}
}

func getUnaryServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// other interceptors ...

	options = append(options, grpc.UnaryInterceptor(
	    interceptor.UnaryServerJwtAuth(
	        // Choose to use one of the following 4 authorization
			
	        // Case 1: default authorization
	        // interceptor.WithDefaultVerify(), // can be ignored
	        // default authorization with extra verification
	        // interceptor.WithDefaultVerify(extraDefaultVerifyFn),

	        // Case 2: custom authorization
	        // interceptor.WithCustomVerify(),
	        // custom authorization with extra verification
	        // interceptor.WithCustomVerify(extraCustomVerifyFn),

	        // specify the gRPC API to ignore token verification(full path)
	        interceptor.WithAuthIgnoreMethods(
	            "/api.user.v1.User/Register",
	            "/api.user.v1.User/Login",
	        ),
	    ),
	))

	return options
}


type user struct {
	userV1.UnimplementedUserServer
}

// Login ...
func (s *user) Login(ctx context.Context, req *userV1.LoginRequest) (*userV1.LoginReply, error) {
	// check user and password success

	// Case 1: default authorization
	token, err := jwt.GenerateToken("123", "admin")
	
	// Case 2: custom authorization
	fields := jwt.KV{"id": uint64(100), "name": "tom", "age": 10}
	token, err := jwt.GenerateCustomToken(fields)

	return &userV1.LoginReply{Token: token},nil
}

// GetByID ...
func (s *user) GetByID(ctx context.Context, req *userV1.GetUserByIDRequest) (*userV1.GetUserByIDReply, error) {
	// if token is valid, won't get here, because the interceptor has returned an error message 

	// if you want get jwt claims, you can use the following code
	// Case 1: default authorization
	claims, err := interceptor.GetJwtClaims(ctx)
	
	// Case 2: custom authorization
	customClaims, err := interceptor.GetJwtCustomClaims(ctx)
	
	// ......

	return &userV1.GetUserByIDReply{},nil
}

func extraDefaultVerifyFn(claims *jwt.Claims, tokenTail10 string) error {
	// In addition to jwt certification, additional checks can be customized here.

	// err := errors.New("verify failed")
	// if claims.Name != "admin" {
	//     return err
	// }
	// token := getToken(claims.UID) // from cache or database
	// if tokenTail10 != token[len(token)-10:] { return err }

	return nil
}

func extraCustomVerifyFn(claims *jwt.CustomClaims, tokenTail10 string) error {
	// In addition to jwt certification, additional checks can be customized here.

	// err := errors.New("verify failed")
	// token, fields := getToken(id) // from cache or database
	// if tokenTail10 != token[len(token)-10:] { return err }

	// id, exist := claims.GetUint64("id")
	// if !exist || id != fields["id"].(uint64) { return err }

	// name, exist := claims.GetString("name")
	// if !exist || name != fields["name"].(string) { return err }

	// age, exist := claims.GetInt("age")
	// if !exist || age != fields["age"].(int) { return err }

	return nil
}
```

<br>
