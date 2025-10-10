package server

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/config"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/handlers"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/middleware"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/routes"
)

// Server represents the API Gateway server
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	config     config.Config
	ctx        context.Context
	cancel     context.CancelFunc
	handlers   *handlers.Handlers
}

// NewServer creates a new API Gateway server
func NewServer(ctx context.Context, cfg config.Config) *Server {
	ctx, cancel := context.WithCancel(ctx)

	// Create HTTP mux
	mux := http.NewServeMux()

	// Initialize handlers
	handlerInstances := handlers.NewHandlers()

	server := &Server{
		mux:      mux,
		config:   cfg,
		ctx:      ctx,
		cancel:   cancel,
		handlers: handlerInstances,
	}

	// Setup routes
	server.setupRoutes()

	// Create HTTP server with middleware chain
	server.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      server.createMiddlewareChain(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.ShutdownTimeout,
	}

	return server
}

// Apply middleware in reverse order (last middleware wraps first)
// createMiddlewareChain creates the middleware chain for the HTTP server
func (s *Server) createMiddlewareChain() http.Handler {
	var handler http.Handler = s.mux

	// Error recovery middleware
	handler = middleware.RecoveryMiddleware(handler)

	// Request ID middleware
	handler = middleware.RequestIDMiddleware(handler)

	// CORS middleware
	handler = middleware.CORSMiddleware(handler, middleware.CORSConfig{
		AllowOrigins: s.config.Security.AllowOrigins,
		AllowMethods: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
	})

	// Logging middleware (if enabled)
	if s.config.Logging.RequestLog {
		handler = middleware.LoggingMiddleware(handler, middleware.LoggingConfig{
			SkipPaths:   []string{},
			LogHeaders:  false,
			LogBody:     false,
			MaxBodySize: 1024,
		})
	}

	return handler
}

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	router := routes.NewRouter(s.mux, s.handlers, s.config)
	router.SetupRoutes()
}

// Start starts the server
func (s *Server) Start(port uint64) {
	if port != 0 {
		s.httpServer.Addr = fmt.Sprintf("%s:%d", s.config.Server.Host, port)
	}

	logger.Global().Info("Starting API Gateway server",
		zap.String("address", s.httpServer.Addr),
		zap.String("version", "1.0.0"),
	)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Global().Error("failed to start server", zap.Error(err))
		}
	}()

	<-s.ctx.Done()
	if err := s.Shutdown(); err != nil {
		logger.Global().Error("failed to shutdown server", zap.Error(err))
	}
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	logger.Global().Info("Shutting down API Gateway server")

	ctx, cancel := context.WithTimeout(context.Background(), s.config.Server.ShutdownTimeout)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

// GetMux returns the HTTP mux instance for testing
func (s *Server) GetMux() *http.ServeMux {
	return s.mux
}

// Stop stops the server
func (s *Server) Stop() {
	s.cancel()
}
