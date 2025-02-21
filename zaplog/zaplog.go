package zaplog

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// creates configured logger
func New() *zap.Logger {
	cfg := zap.Config{
		Level: zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding: "console", // or "json"
		EncoderConfig: zapcore.EncoderConfig{
			CallerKey:     	"caller",
			MessageKey:    	"msg",
			LineEnding:    	zapcore.DefaultLineEnding,
			EncodeLevel:   	zapcore.CapitalColorLevelEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths: []string{
			"stdout",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Panicf("failed to build logger: %v", err)
	}

	return logger
}
