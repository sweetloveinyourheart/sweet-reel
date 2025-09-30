package logger

import (
	"encoding/json"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	*SugaredLogger
	atomicLevel zap.AtomicLevel
}

// New creates a logger with default config and log level.
// Default Level is INFO.
func New() *logger {
	return new(zap.InfoLevel)
}

// NewFromJson creates a logger from the raw JSON config.
// See zap.Config for all the available params.
func NewFromJson(config []byte) (*logger, error) {
	var cfg zap.Config
	if err := json.Unmarshal(config, &cfg); err != nil {
		return nil, err
	}
	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	log := NewOtel(zapLogger, WithTraceIDField(true))

	return &logger{log.Sugar(), cfg.Level}, nil
}

// SetZapLevel allows to change the log Level at runtime.
// INFO is default.
func (l *logger) SetZapLevel(level zapcore.Level) *logger {
	l.atomicLevel.SetLevel(level)
	return l
}

// SetStringLevel parses the log level from string and sets it.
// INFO is default.
func (l *logger) SetStringLevel(level string) *logger {
	// ignore an error, default log level is INFO
	zapLevel, _ := zapcore.ParseLevel(level)
	return l.SetZapLevel(zapLevel)
}

// SetAsGlobal sets a receiver as the global logger.
// Use Global() to access it from anywhere.
func (l *logger) SetAsGlobal() {
	_ = ReplaceGlobals(l.Desugar())
}

func new(level zapcore.Level) *logger {
	al := zap.NewAtomicLevelAt(level)
	ws := zapcore.AddSync(os.Stdout)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000Z0700")
	encoderJSON := zapcore.NewJSONEncoder(encoderConfig)

	var zapLogger *zap.Logger
	coreStd := zapcore.NewCore(encoderJSON, ws, al)
	zapLogger = zap.New(coreStd, zap.AddCaller())

	log := NewOtel(zapLogger, WithTraceIDField(true), WithMinLevel(level))

	return &logger{log.Sugar(), al}
}

// Global returns the global logger.
// Use logger.SetAsGlobal() to replace it.
func Global() *Logger {
	return L()
}

// GlobalSugared returns the global sugared logger.
// See zap.SugaredLogger for details.
// Use logger.SetAsGlobal() to replace it.
func GlobalSugared() *SugaredLogger {
	return S()
}
