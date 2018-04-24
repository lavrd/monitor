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

// Debugf print formatted debug log
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Errorf print formatted error log
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Error print error log
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Fatal print fatal log
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}
