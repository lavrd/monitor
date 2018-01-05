package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/spacelavr/dlm/pkg/context"
)

// Info print info log
func Info(info ...interface{}) {
	if context.Get().Verbose {
		logrus.Info(info)
	}
}

// Panic print panic log
func Panic(err ...interface{}) {
	logrus.Panic(err)
}
