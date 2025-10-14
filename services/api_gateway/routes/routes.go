package routes

import (
	"net/http"

	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/config"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/handlers"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/helpers"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/middleware"
)

// Router manages all routes for the API Gateway
type Router struct {
	mux      *http.ServeMux
	handlers *handlers.Handlers
	config   config.Config
}

// NewRouter creates a new router instance
func NewRouter(mux *http.ServeMux, handlers *handlers.Handlers, config config.Config) *Router {
	return &Router{
		mux:      mux,
		handlers: handlers,
		config:   config,
	}
}

// SetupRoutes configures all routes
func (r *Router) SetupRoutes() {
	r.setupPublicRoutes()
	r.setupProtectedRoutes()
}

// setupPublicRoutes sets up public API routes
func (r *Router) setupPublicRoutes() {
	r.mux.Handle("/api/v1/oauth", helpers.POST(r.handlers.AuthHandler.GoogleOAuth))
	r.mux.Handle("/api/v1/auth/refresh-token", helpers.GET(r.handlers.AuthHandler.RefreshToken))
}

// setupProtectedRoutes sets up authenticated API routes
func (r *Router) setupProtectedRoutes() {
	// Create authentication middleware
	authMiddleware := middleware.NewAuthMiddleware(middleware.AuthConfig{
		SigningKey: r.config.Security.JWTSecret,
		SkipPaths: []string{
			"/api/v1/public",
		},
	})

	// Video management routes
	r.mux.Handle("/api/v1/videos/presigned-url", authMiddleware(helpers.POST(r.handlers.VideoManagement.GeneratePresignedURL)))
}
