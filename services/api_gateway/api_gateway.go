package apigateway

import (
	"context"

	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/config"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/server"
)

// NewServer creates a new API Gateway server using the new internal structure
func NewServer(ctx context.Context, port uint64, signingKey string) *server.Server {
	return server.NewServer(ctx, config.LoadServerConfig(port, signingKey))
}

// InitializeRepos initializes any repositories or dependencies specific to API Gateway
func InitializeRepos(ctx context.Context) error {
	return nil
}
