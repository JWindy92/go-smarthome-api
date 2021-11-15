package logwrapper

import "go.uber.org/zap"

type Logger struct {
	Logger *zap.SugaredLogger
}

func NewLogger() Logger {
	return Logger{
		Logger: zap.NewExample().Sugar(),
	}
}
