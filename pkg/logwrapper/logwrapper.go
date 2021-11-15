package logwrapper

import "go.uber.org/zap"

type Logger struct {
	Logger *zap.SugaredLogger
}

func NewLogger(source string) Logger {
	logger := zap.NewExample().Sugar().With(
		"source", source,
	)
	return Logger{
		Logger: logger,
	}
}
