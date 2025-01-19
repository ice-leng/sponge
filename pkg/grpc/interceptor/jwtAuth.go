package interceptor

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/go-dev-frame/sponge/pkg/jwt"
)

// ---------------------------------- client ----------------------------------

// SetJwtTokenToCtx set the token (excluding prefix Bearer) to the context in grpc client side
// Example:
//
// authorization := "Bearer jwt-token"
//
//	ctx := SetJwtTokenToCtx(ctx, authorization)
//	cli.GetByID(ctx, req)
func SetJwtTokenToCtx(ctx context.Context, token string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md.Set(headerAuthorize, authScheme+" "+token)
	} else {
		md = metadata.Pairs(headerAuthorize, authScheme+" "+token)
	}
	return metadata.NewOutgoingContext(ctx, md)
}

// SetAuthToCtx set the authorization (including prefix Bearer) to the context in grpc client side
// Example:
//
//	ctx := SetAuthToCtx(ctx, authorization)
//	cli.GetByID(ctx, req)
func SetAuthToCtx(ctx context.Context, authorization string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md.Set(headerAuthorize, authorization)
	} else {
		md = metadata.Pairs(headerAuthorize, authorization)
	}
	return metadata.NewOutgoingContext(ctx, md)
}

// ---------------------------------- server interceptor ----------------------------------

var (
	headerAuthorize = "authorization"

	// auth Scheme
	authScheme = "Bearer"

	// authentication information in ctx key name
	authCtxClaimsName = "tokenInfo"

	// collection of skip authentication methods
	authIgnoreMethods = map[string]struct{}{}
)

// GetAuthorization combining tokens into authentication information
func GetAuthorization(token string) string {
	return authScheme + " " + token
}

// GetAuthCtxKey get the name of Claims
func GetAuthCtxKey() string {
	return authCtxClaimsName
}

// ExtraDefaultVerifyFn extra default verify function, tokenTail10 is the last 10 characters of the token.
type ExtraDefaultVerifyFn = func(claims *jwt.Claims, tokenTail10 string) error

// ExtraCustomVerifyFn extra custom verify function, tokenTail10 is the last 10 characters of the token.
type ExtraCustomVerifyFn = func(claims *jwt.CustomClaims, tokenTail10 string) error

const (
	defaultAuthType = 1 // use default auth
	customAuthType  = 2 // use custom auth
)

// AuthOption setting the Authentication Field
type AuthOption func(*authOptions)

// authOptions settings
type authOptions struct {
	authScheme    string
	ctxClaimsName string
	ignoreMethods map[string]struct{}

	authType        int
	defaultVerifyFn ExtraDefaultVerifyFn
	customVerifyFn  ExtraCustomVerifyFn
}

func defaultAuthOptions() *authOptions {
	return &authOptions{
		authScheme:    authScheme,
		ctxClaimsName: authCtxClaimsName,
		ignoreMethods: make(map[string]struct{}), // ways to ignore forensics

		authType: defaultAuthType,
	}
}

func (o *authOptions) apply(opts ...AuthOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithAuthScheme set the message prefix for authentication
func WithAuthScheme(scheme string) AuthOption {
	return func(o *authOptions) {
		o.authScheme = scheme
	}
}

// WithAuthClaimsName set the key name of the information in ctx for authentication
func WithAuthClaimsName(claimsName string) AuthOption {
	return func(o *authOptions) {
		o.ctxClaimsName = claimsName
	}
}

// WithAuthIgnoreMethods ways to ignore forensics
// fullMethodName format: /packageName.serviceName/methodName,
// example /api.userExample.v1.userExampleService/GetByID
func WithAuthIgnoreMethods(fullMethodNames ...string) AuthOption {
	return func(o *authOptions) {
		for _, method := range fullMethodNames {
			o.ignoreMethods[method] = struct{}{}
		}
	}
}

// WithDefaultVerify set default verify type
func WithDefaultVerify(fn ...ExtraDefaultVerifyFn) AuthOption {
	return func(o *authOptions) {
		o.authType = defaultAuthType
		if len(fn) > 0 {
			o.defaultVerifyFn = fn[0]
		}
	}
}

// WithCustomVerify set custom verify type with extra verify function
func WithCustomVerify(fn ...ExtraCustomVerifyFn) AuthOption {
	return func(o *authOptions) {
		o.authType = customAuthType
		if len(fn) > 0 {
			o.customVerifyFn = fn[0]
		}
	}
}

// -------------------------------------------------------------------------------------------

// verify authorization from context, support default and custom verify processing
func jwtVerify(ctx context.Context, opt *authOptions) (context.Context, error) {
	if opt == nil {
		opt = defaultAuthOptions()
	}

	token, err := grpc_auth.AuthFromMD(ctx, authScheme) // key is authScheme
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "authFromMD error: %v", err)
	}

	if len(token) <= 100 {
		return ctx, status.Errorf(codes.Unauthenticated, "authorization is illegal")
	}

	// custom claims verify
	if opt.authType == customAuthType {
		var claims *jwt.CustomClaims
		claims, err = jwt.ParseCustomToken(token)
		if err != nil {
			return ctx, status.Errorf(codes.Unauthenticated, "use custom verify type, ParseCustomToken error: %v", err)
		}
		// extra custom verify function
		if opt.customVerifyFn != nil {
			tokenTail10 := token[len(token)-10:]
			err = opt.customVerifyFn(claims, tokenTail10)
			if err != nil {
				return ctx, status.Errorf(codes.Unauthenticated, "customVerifyFn error: %v", err)
			}
		}
		newCtx := context.WithValue(ctx, authCtxClaimsName, claims) //nolint
		return newCtx, nil
	}

	// default claims verify
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "use default verify type, %v", err)
	}
	if opt.defaultVerifyFn != nil {
		tokenTail10 := token[len(token)-10:]
		// extra default verify function
		err = opt.defaultVerifyFn(claims, tokenTail10)
		if err != nil {
			return ctx, status.Errorf(codes.Unauthenticated, "verifyFn error: %v", err)
		}
	}
	newCtx := context.WithValue(ctx, authCtxClaimsName, claims) //nolint
	return newCtx, nil
}

// GetJwtClaims get the jwt default claims from context, contains fixed fields uid and name
func GetJwtClaims(ctx context.Context) (*jwt.Claims, bool) {
	v, ok := ctx.Value(authCtxClaimsName).(*jwt.Claims)
	return v, ok
}

// GetJwtCustomClaims get the jwt custom claims from context, contains custom fields
func GetJwtCustomClaims(ctx context.Context) (*jwt.CustomClaims, bool) {
	v, ok := ctx.Value(authCtxClaimsName).(*jwt.CustomClaims)
	return v, ok
}

// UnaryServerJwtAuth jwt unary interceptor
func UnaryServerJwtAuth(opts ...AuthOption) grpc.UnaryServerInterceptor {
	o := defaultAuthOptions()
	o.apply(opts...)
	authScheme = o.authScheme
	authCtxClaimsName = o.ctxClaimsName
	authIgnoreMethods = o.ignoreMethods

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var err error

		if _, ok := authIgnoreMethods[info.FullMethod]; ok {
			newCtx = ctx
		} else {
			newCtx, err = jwtVerify(ctx, o)
			if err != nil {
				return nil, err
			}
		}

		return handler(newCtx, req)
	}
}

// StreamServerJwtAuth jwt stream interceptor
func StreamServerJwtAuth(opts ...AuthOption) grpc.StreamServerInterceptor {
	o := defaultAuthOptions()
	o.apply(opts...)
	authScheme = o.authScheme
	authCtxClaimsName = o.ctxClaimsName
	authIgnoreMethods = o.ignoreMethods

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		var err error

		if _, ok := authIgnoreMethods[info.FullMethod]; ok {
			newCtx = stream.Context()
		} else {
			newCtx, err = jwtVerify(stream.Context(), o)
			if err != nil {
				return err
			}
		}

		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}

// ----------------------------------------------------

// StandardVerifyFn default verify function, tokenTail10 is the last 10 characters of the token.
// Deprecated: use ExtraDefaultVerifyFn instead.
type StandardVerifyFn = ExtraDefaultVerifyFn

// CustomVerifyFn custom verify function, tokenTail10 is the last 10 characters of the token.
// Deprecated: use ExtraCustomVerifyFn instead.
type CustomVerifyFn = ExtraCustomVerifyFn

// WithStandardVerify set default verify type with extra verify function
// Deprecated: use WithExtraDefaultVerify instead.
func WithStandardVerify(fn StandardVerifyFn) AuthOption {
	return WithDefaultVerify(fn)
}
