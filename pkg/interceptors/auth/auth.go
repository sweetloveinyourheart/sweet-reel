package auth

import (
	"context"

	"connectrpc.com/connect"
)

// AuthFunc is the pluggable function that performs authentication.
//
// The passed in `Context` will contain the gRPC metadata.MD object (for header-based authentication) and
// the peer.Peer information that can contain transport-based credentials (e.g. `credentials.AuthInfo`).
//
// The returned context will be propagated to handlers, allowing user changes to `Context`. However,
// please make sure that the `Context` returned is a child `Context` of the one passed in.
//
// If error is returned, its `grpc.Code()` will be returned to the user as well as the verbatim message.
// Please make sure you use `codes.Unauthenticated` (lacking auth) and `codes.PermissionDenied`
// (authed, but lacking perms) appropriately.
type AuthFunc func(ctx context.Context, token string) (context.Context, error)

// ServiceAuthFuncOverride allows a given gRPC service implementation to override the global `AuthFunc`.
//
// If a service implements the AuthFuncOverride method, it takes precedence over the `AuthFunc` method,
// and will be called instead of AuthFunc for all method invocations within that service.
type ServiceAuthFuncOverride interface {
	AuthFuncOverride(ctx context.Context, token string, fullMethodName string) (context.Context, error)
}

// UnaryServerInterceptor returns a new unary server interceptors that performs per-request auth.
func UnaryServerInterceptor(authFunc AuthFunc, o ...Option) connect.UnaryInterceptorFunc {
	opts := evaluateOptions(o)
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (resp connect.AnyResponse, err error) {
			token := req.Header().Get("Authorization")

			var newCtx context.Context
			if opts.override != nil {
				newCtx, err = opts.override.AuthFuncOverride(ctx, token, req.Spec().Procedure)
			} else {
				newCtx, err = authFunc(ctx, token)
			}
			if err != nil {
				return nil, err
			}
			return next(newCtx, req)
		})
	}
}

func NewAuthInterceptor(authFunc AuthFunc, o ...Option) connect.Interceptor {
	opts := evaluateOptions(o)
	return &Interceptor{
		authFunc: authFunc,
		opts:     opts,
	}
}

type Interceptor struct {
	authFunc AuthFunc
	opts     *options
}

func (i Interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		token := req.Header().Get("Authorization")
		var newCtx context.Context
		var err error
		if i.opts.override != nil {
			newCtx, err = i.opts.override.AuthFuncOverride(ctx, token, req.Spec().Procedure)
		} else {
			newCtx, err = i.authFunc(ctx, token)
		}
		if err != nil {
			return nil, err
		}
		return next(newCtx, req)
	})
}

func (i Interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		token := conn.RequestHeader().Get("Authorization")
		var newCtx context.Context
		var err error
		if i.opts.override != nil {
			newCtx, err = i.opts.override.AuthFuncOverride(ctx, token, conn.Spec().Procedure)
		} else {
			newCtx, err = i.authFunc(ctx, token)
		}
		if err != nil {
			return err
		}
		return next(newCtx, conn)
	}
}
