package interceptors

import (
	"connectrpc.com/connect"

	auth_interceptors "github.com/sweetloveinyourheart/sweet-reel/pkg/interceptors/auth"
	client_interceptors "github.com/sweetloveinyourheart/sweet-reel/pkg/interceptors/client"
)

func CommonConnectInterceptors(serviceName string, signingKey string, authFunc auth_interceptors.AuthFunc, authOpts ...auth_interceptors.Option) []connect.Interceptor {
	if authFunc == nil {
		authFunc = ConnectAuthHandler(signingKey)
	}

	return []connect.Interceptor{
		auth_interceptors.NewAuthInterceptor(authFunc, authOpts...),
	}
}

func CommonConnectClientInterceptors(serviceName string, signingKey string) []connect.Interceptor {
	return []connect.Interceptor{
		client_interceptors.NewServiceClientInterceptor(serviceName, signingKey),
	}
}
