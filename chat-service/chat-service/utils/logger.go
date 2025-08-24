package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger(msg string, err error) {
	// Configure zap to include caller information and stack traces
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = "stacktrace"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // Shows file:line

	// Create a new logger with the configured options
	logger, err := config.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		// Fallback to a basic logger if configuration fails
		zap.NewExample().Sugar().Panicf("Failed to initialize logger: %v", err)
		return
	}
	defer logger.Sync() // Flushes buffer, if any

	// Use SugaredLogger for simpler API
	sugar := logger.Sugar()

	// Log the error with file, line, message, and full error description
	sugar.Infow("Error occurred",
		"message", msg,
		"error", err,
		"caller", zapcore.EntryCaller{}, // Automatically includes file:line
	)

	// If no error is provided, log only the message with caller
	if err == nil {
		sugar.Infow("Info message",
			"message", msg,
			"caller", zapcore.EntryCaller{},
		)
	}
}
