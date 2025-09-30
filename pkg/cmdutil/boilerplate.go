package cmdutil

import (
	"context"
	"strings"
	"sync"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/config"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

type AppRun struct {
	serviceType string
	serviceKey  string
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	readyChan   chan bool
}

func st(s string) (serviceType string, serviceKey string) {
	serviceType = strings.ToLower(s)
	serviceKey = strings.ReplaceAll(serviceType, "_", "-")
	serviceKey = strings.ReplaceAll(serviceKey, "-", ".")
	return
}

func BoilerplateRun(serviceType string) (*AppRun, error) {
	serviceType, serviceKey := st(serviceType)

	ctx, cancel := context.WithCancel(context.Background())

	// Logger init.
	logLevel := config.Instance().GetString("log.level")
	loggerInstance := logger.New().SetStringLevel(logLevel)
	loggerInstance.SetAsGlobal()
	defer func() {
		_ = loggerInstance.Sync()
	}()

	serviceName := config.Instance().GetString("service")
	loggerInstance.Infow("starting service",
		"type", serviceType,
		"service", serviceName,
	)

	healthCheckPort := config.Instance().GetInt("healthcheck.port")
	healthCheckWebPort := config.Instance().GetInt("healthcheck.web.port")
	readyChan := StartHealthServices(ctx, serviceName, healthCheckPort, healthCheckWebPort)

	return &AppRun{
		serviceType: serviceType,
		serviceKey:  serviceKey,
		wg:          sync.WaitGroup{},
		ctx:         ctx,
		cancel:      cancel,
		readyChan:   readyChan,
	}, nil
}
