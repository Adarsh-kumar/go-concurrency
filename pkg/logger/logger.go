package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	instance *zap.SugaredLogger
	once     sync.Once
)

// GetLogger returns the singleton logger instance.
func GetLogger() *zap.SugaredLogger {
	once.Do(func() {
		// Initialize the logger
		logger, err := zap.NewProduction() 
		if err != nil {
			panic("Failed to initialize logger: " + err.Error())
		}
		instance = logger.Sugar()
	})
	return instance
}
