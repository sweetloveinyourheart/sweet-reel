package interceptors

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/grpc"
)

// ConnectAuthHandler is used by a middleware to authenticate requests
func ConnectAuthHandler(signingKey string) func(context.Context, string) (context.Context, error) {
	return func(ctx context.Context, tokenString string) (context.Context, error) {
		splits := strings.SplitN(tokenString, " ", 2)
		if len(splits) < 2 {
			return nil, status.Errorf(codes.Unauthenticated, "Bad authorization string")
		}
		if !strings.EqualFold(splits[0], "bearer") {
			return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated with bearer")
		}
		return authHandler(ctx, signingKey, splits[1])
	}
}

// ConnectServerAuthHandler is used by a middleware to authenticate server to server requests
func ConnectServerAuthHandler(signingKey string) func(context.Context, string) (context.Context, error) {
	return func(ctx context.Context, tokenString string) (context.Context, error) {
		splits := strings.SplitN(tokenString, " ", 2)
		if len(splits) < 2 {
			return nil, status.Errorf(codes.Unauthenticated, "Bad authorization string")
		}
		if !strings.EqualFold(splits[0], "bearer") {
			return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated with bearer")
		}
		return serverAuthHandler(ctx, signingKey, splits[1])
	}
}

// AuthHandler is used by a middleware to authenticate requests
func authHandler(ctx context.Context, signingKey string, tokenString string) (context.Context, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})

	if token == nil || !token.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		rawUserID, ok := claims["user_id"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid token, malformed id")
		}
		userID, err := uuid.FromString(rawUserID)
		if err != nil {
			return nil, fmt.Errorf("invalid token, malformed id")
		}

		grpc_ctxtags.Extract(ctx).Set("auth.sub", claims)
		newCtx := context.WithValue(ctx, grpc.AuthToken, userID)

		return newCtx, nil
	} else {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

}

// ServerAuthHandler is used by a middleware to authenticate server to server requests
func serverAuthHandler(ctx context.Context, signingKey string, tokenString string) (context.Context, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})

	if token == nil || !token.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		rawServiceName, ok := claims["service_name"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid token, malformed service name")
		}

		grpc_ctxtags.Extract(ctx).Set("auth.sub", claims)
		newCtx := context.WithValue(ctx, grpc.AuthServiceToken, rawServiceName)

		return newCtx, nil
	} else {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}
}
