package logwrapper

import "go.uber.org/zap"

type Logger struct {
	Logger *zap.SugaredLogger
}

func NewLogger() Logger {
	logger := zap.NewExample().Sugar()
	return Logger{
		Logger: logger,
	}
}
