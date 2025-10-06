package client

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/timeutil"
)

func NewServiceClientInterceptor(serviceName string, signingKey string) connect.Interceptor {
	return &ServiceInterceptor{
		serviceName: serviceName,
		signingKey:  signingKey,
	}
}

const bearerPrefix = "Bearer "
const authorizationHeaderKey = "Authorization"

type ServiceInterceptor struct {
	serviceName string
	signingKey  string
}

func (i ServiceInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if req.Spec().IsClient {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"service_name": i.serviceName,
				"nbf":          timeutil.NowRoundedForGranularity().Unix(),
			})
			// Sign and get the complete encoded token as a string using the secret
			if len(i.signingKey) == 0 {
				return nil, fmt.Errorf("invalid signing token, misconfigured instance")
			}
			tokenString, err := token.SignedString([]byte(i.signingKey))
			if err != nil {
				return nil, err
			}

			req.Header().Add(authorizationHeaderKey, bearerPrefix+tokenString)
			return next(ctx, req)
		}
		return next(ctx, req)
	})
}

func (i ServiceInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, s connect.Spec) connect.StreamingClientConn {
		conn := next(ctx, s)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"service_name": i.serviceName,
			"nbf":          timeutil.NowRoundedForGranularity().Unix(),
		})
		// Sign and get the complete encoded token as a string using the secret
		if len(i.signingKey) == 0 {
			return conn
		}
		tokenString, err := token.SignedString([]byte(i.signingKey))
		if err != nil {
			return conn
		}

		conn.RequestHeader().Add(authorizationHeaderKey, bearerPrefix+tokenString)
		return conn
	}
}

func (i ServiceInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		return next(ctx, conn)
	}
}
