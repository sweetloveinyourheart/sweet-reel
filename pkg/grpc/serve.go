package grpc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/rs/cors"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

func ServeBuf(ctx context.Context, path string, rpcHandler http.Handler, port uint64, serviceName string) {
	mux := http.NewServeMux()
	mux.Handle(path, rpcHandler)

	// Use h2c so we can serve HTTP/2 without TLS.
	handler := h2c.NewHandler(newCORS().Handler(mux), &http2.Server{})

	logger.GlobalSugared().Infof("Buf %s listening on port %d\n", serviceName, port)
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logger.GlobalSugared().Panicf("%s failed to serve: %v", serviceName, err)
			} else {
				logger.GlobalSugared().Infof("Buf %s server closed", serviceName)
			}
		}
	}()
	<-ctx.Done()
	logger.GlobalSugared().Infof("Buf %s shutting down", serviceName)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.GlobalSugared().Panicf("%s failed to shutdown: %v", serviceName, err)
	}
}

func newCORS() *cors.Cors {
	// To let web developers play with the demo service from browsers, we need a
	// very permissive CORS setup.
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowOriginFunc: func(origin string) bool {
			// Allow all origins, which effectively disables CORS.
			return true
		},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{
			// Content-Type is in the default safelist.
			"Accept",
			"Accept-Encoding",
			"Accept-Post",
			"Authorization",
			"Cache-Control",
			"Connect-Accept-Encoding",
			"Connect-Content-Encoding",
			"Connect-Protocol-Version",
			"Content-Encoding",
			"Content-Length",
			"Content-Type",
			"DNT",
			"If-Modified-Since",
			"Keep-Alive",
			"Grpc-Accept-Encoding",
			"Grpc-Encoding",
			"Grpc-Message",
			"Grpc-Status",
			"Grpc-Status-Details-Bin",
			"Grpc-Timeout",
			"Grpc-Web",
			"TraceParent",
			"TraceState",
			"Timeout",
			"User-Agent",
			"X-CSRF-Token",
			"X-Datadog-Origin",
			"X-Datadog-Parent-Id",
			"X-Datadog-Trace-Id",
			"X-Datadog-Sampling-Priority",
			"X-Grpc-Web",
			"X-User-Agent",
			"X-Requested-With",
			"X-Robots-Tag",
		},
		// Let browsers cache CORS information for longer, which reduces the number
		// of preflight requests. Any changes to ExposedHeaders won't take effect
		// until the cached data expires. FF caps this value at 24h, and modern
		// Chrome caps it at 2h.
		MaxAge: int(2 * time.Hour / time.Second),
	})
}
