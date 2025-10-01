package apputils

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

var (
	retryCount      int
	retryDelay      time.Duration
	method          string
	body            string
	timeout         time.Duration
	grpcServiceName string
)

var checkCmd = &cobra.Command{
	Use:   `check`,
	Short: "Health check commands",
	Long:  "Commands for running health checks",
}

var checkUrlCmd = &cobra.Command{
	Use:     `http [address]`,
	Short:   "HTTP(S) health check",
	Example: `  check http https://httpbin.org/anything --method=POST --retry=2 --delay=2s --timeout=3s --body='{"test"=1}'`,
	RunE:    runCheckUrl,
}

var checkGRPCCmd = &cobra.Command{
	Use:     `grpc [address]`,
	Short:   "gRPC health check",
	Example: `  check grpc localhost:8080 --service=podinfo --retry=1 --delay=2s --timeout=2s`,
	RunE:    runCheckGPRC,
}

func CheckCommand() *cobra.Command {
	checkUrlCmd.Flags().StringVar(&method, "method", "GET", "HTTP method")
	checkUrlCmd.Flags().StringVar(&body, "body", "", "HTTP POST/PUT content")
	checkUrlCmd.Flags().IntVar(&retryCount, "retry", 0, "times to retry the HTTP call")
	checkUrlCmd.Flags().DurationVar(&retryDelay, "delay", 1*time.Second, "wait duration between retries")
	checkUrlCmd.Flags().DurationVar(&timeout, "timeout", 5*time.Second, "timeout")
	checkCmd.AddCommand(checkUrlCmd)

	checkGRPCCmd.Flags().IntVar(&retryCount, "retry", 0, "times to retry the TCP check")
	checkGRPCCmd.Flags().DurationVar(&retryDelay, "delay", 1*time.Second, "wait duration between retries")
	checkGRPCCmd.Flags().DurationVar(&timeout, "timeout", 5*time.Second, "timeout")
	checkGRPCCmd.Flags().StringVar(&grpcServiceName, "service", "", "gRPC service name")
	checkCmd.AddCommand(checkGRPCCmd)

	return checkCmd
}

func runCheckUrl(cmd *cobra.Command, args []string) error {
	if retryCount < 0 {
		return fmt.Errorf("--retry is required")
	}
	if len(args) < 1 {
		return fmt.Errorf("address is required! example: check http https://httpbin.org")
	}

	address := args[0]
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		address = fmt.Sprintf("http://%s", address)
	}

	for n := 0; n <= retryCount; n++ {
		if n != 1 {
			time.Sleep(retryDelay)
		}

		req, err := http.NewRequestWithContext(context.Background(), method, address, bytes.NewBuffer([]byte(body)))
		if err != nil {
			logger.Global().Info("check failed",
				zap.String("address", address),
				zap.Error(err))
			os.Exit(1)
		}

		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		resp, err := http.DefaultClient.Do(req.WithContext(ctx))
		cancel()
		if err != nil {
			logger.Global().InfoContext(ctx, "check failed",
				zap.String("address", address),
				zap.Error(err))
			continue
		}

		if resp.Body != nil {
			resp.Body.Close()
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			logger.Global().InfoContext(ctx, "check succeed",
				zap.String("address", address),
				zap.Int("status code", resp.StatusCode),
				zap.String("response size", fmtContentLength(resp.ContentLength)))
			os.Exit(0)
		} else {
			logger.Global().InfoContext(ctx, "check failed",
				zap.String("address", address),
				zap.Int("status code", resp.StatusCode))
			continue
		}
	}

	os.Exit(1)
	return nil
}

func runCheckGPRC(cmd *cobra.Command, args []string) error {
	if retryCount < 0 {
		return fmt.Errorf("--retry is required")
	}
	if len(args) < 1 {
		return fmt.Errorf("address is required! example: check grpc localhost:8080")
	}
	address := args[0]

	for n := 0; n <= retryCount; n++ {
		if n != 1 {
			time.Sleep(retryDelay)
		}

		conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Global().Info("check failed",
				zap.String("address", address),
				zap.Error(err))
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		resp, err := grpc_health_v1.NewHealthClient(conn).Check(ctx, &grpc_health_v1.HealthCheckRequest{
			Service: grpcServiceName,
		})
		cancel()

		if err != nil {
			if stat, ok := status.FromError(err); ok && stat.Code() == codes.Unimplemented {
				logger.Global().InfoContext(ctx, "gPRC health protocol not implemented")
				os.Exit(1)
			} else {
				logger.Global().InfoContext(ctx, "check failed",
					zap.String("address", address),
					zap.Error(err))
			}
			continue
		}

		conn.Close()
		logger.Global().InfoContext(ctx, "check succeed",
			zap.String("status", resp.GetStatus().String()))
		os.Exit(0)

	}

	os.Exit(1)
	return nil
}

func fmtContentLength(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
