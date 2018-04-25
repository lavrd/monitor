package log

import (
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func init() {
	logger = logrus.New()
}

// SetVerbose set verbose output
func SetVerbose(verbose bool) {
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	}
}

// Debug print debug log
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Error print error log
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Fatal print fatal log
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}
