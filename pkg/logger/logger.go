package logger

import (
	"github.com/lavrs/dlm/pkg/context"
	"github.com/sirupsen/logrus"
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
