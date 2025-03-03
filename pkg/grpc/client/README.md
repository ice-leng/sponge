## gRPC client

`client` is a gRPC client library for Go. It provides a simple way to connect to a gRPC server and call its methods.

### Example of use

```go
	import "github.com/go-dev-frame/sponge/pkg/grpc/client"

	conn, err := client.Dial(context.Background(), "127.0.0.1:8282",
		//client.WithServiceDiscover(builder),
		//client.WithLoadBalance(),
		//client.WithSecure(credentials),
		//client.WithUnaryInterceptor(unaryInterceptors...),
		//client.WithStreamInterceptor(streamInterceptors...),
	)
```

Examples of practical use https://github.com/zhufuyi/grpc_examples/blob/main/usage/client/main.go
