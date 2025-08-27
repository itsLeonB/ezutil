package ezutil

import "github.com/itsLeonB/ezutil/v2/internal"

type Logger interface {
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Fatal(args ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}

func NewSimpleLogger(namespace string, useColor bool, minLevel int) Logger {
	return &internal.SimpleLogger{
		Namespace: namespace,
		UseColor:  useColor,
		MinLevel:  minLevel,
	}
}
