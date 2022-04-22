package logs

import (
	"os"

	"go.uber.org/zap"
)

var Log *zap.Logger

func init() {
	config := zap.NewDevelopmentConfig()

	// Set basic configuration for Logger
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.Encoding = os.Getenv("LOGGING_ENCODING")
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	Log, _ = config.Build()
}
