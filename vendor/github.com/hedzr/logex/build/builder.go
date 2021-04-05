// Copyright © 2020 Hedzr Yeh.

package build

import (
	"github.com/hedzr/log"
	"github.com/hedzr/logex/logx/logrus"
	"github.com/hedzr/logex/logx/zap"
	"github.com/hedzr/logex/logx/zap/sugar"
)

// New creates and returns a Logger instance from LoggerConfig
func New(config *log.LoggerConfig) log.Logger {
	l, _ := log.ParseLevel(config.Level)
	if l == log.OffLevel {
		return log.NewDummyLogger()
	}

	if bf, ok := builders[config.Backend]; ok {
		return bf(config)
	}

	var logger log.Logger
	switch config.Backend {
	case "dummy", "none", "off":
		logger = log.NewDummyLogger()
	case "std", "standard":
		logger = log.NewStdLogger()
	case "logrus":
		logger = logrus.NewWithConfig(config)
	case "sugar":
		logger = sugar.NewWithConfig(config)
	default:
		logger = zap.NewWithConfig(config)
		//logger = zap.New(config.Level, config.TraceMode, config.DebugMode)
	}
	return logger
}

var builders map[string]log.BuilderFunc

// RegisterBuilder register a builder for your logger.
func RegisterBuilder(backendName string, builderFunc log.BuilderFunc) {
	builders[backendName] = builderFunc
}

func init() {
	builders = make(map[string]log.BuilderFunc)

	builders["dummy"] = log.NewDummyLoggerWithConfig
	builders["none"] = log.NewDummyLoggerWithConfig
	builders["off"] = log.NewDummyLoggerWithConfig

	builders["std"] = log.NewStdLoggerWithConfig
	builders["standard"] = log.NewStdLoggerWithConfig
	builders["go"] = log.NewStdLoggerWithConfig

	builders["logrus"] = logrus.NewWithConfigSimple
	builders["sugar"] = sugar.NewWithConfigSimple
	builders["zap"] = zap.NewWithConfigSimple
}

// NewLoggerConfig returns a default LoggerConfig
func NewLoggerConfig() *log.LoggerConfig {
	c := log.NewLoggerConfig()
	//c.DebugMode = log.GetDebugMode()
	//c.TraceMode = log.GetTraceMode()
	return c
}

// NewLoggerConfigWith returns a default LoggerConfig
func NewLoggerConfigWith(enabled bool, backend, level string) *log.LoggerConfig {
	c := log.NewLoggerConfigWith(enabled, backend, level)
	return c
}
