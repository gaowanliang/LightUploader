package logrus

import (
	"fmt"
	"github.com/hedzr/log"
	"github.com/hedzr/log/exec"
	"github.com/hedzr/logex"
	"github.com/hedzr/logex/formatter"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
)

// New create a sugared logrus logger
//
// level can be: "disable", "panic", "fatal", "error", "warn", "info", "debug", "trace"
//
func New(level string, traceMode, debugMode bool, opts ...Opt) log.Logger {
	log.SetTraceMode(traceMode)
	log.SetDebugMode(debugMode)
	//// ll := cmdr.GetStringR("logger.level", "info")
	//lvl, _ := log.ParseLevel(level)
	//if log.GetDebugMode() {
	//	if lvl < log.DebugLevel {
	//		lvl = log.DebugLevel
	//		level = "debug"
	//	}
	//}
	//if log.GetTraceMode() {
	//	if lvl < log.TraceLevel {
	//		lvl = log.TraceLevel
	//		level = "trace"
	//	}
	//}

	config := log.NewLoggerConfig()
	zl := initLogger(config)

	for _, opt := range opts {
		opt(zl)
	}

	logger := &dzl{Logger: zl, Config: config}
	logger.Setup()
	log.SetLogger(logger) // .AddSkip(extraSkip))
	return logger
}

const extraSkip = 1

func NewWithConfigSimple(config *log.LoggerConfig) log.Logger { return NewWithConfig(config) }

func NewWithConfig(config *log.LoggerConfig, opts ...Opt) log.Logger {
	log.SetTraceMode(config.TraceMode)
	log.SetDebugMode(config.DebugMode)
	// ll := cmdr.GetStringR("logger.level", "info")
	lvl, _ := log.ParseLevel(config.Level)
	if log.GetDebugMode() {
		if lvl < log.DebugLevel {
			lvl = log.DebugLevel
			config.Level = "debug"
		}
	}
	if log.GetTraceMode() {
		if lvl < log.TraceLevel {
			lvl = log.TraceLevel
			config.Level = "trace"
		}
	}

	zl := initLogger(config)

	for _, opt := range opts {
		opt(zl)
	}

	logger := &dzl{Logger: zl, Config: config}
	logger.Setup()
	log.SetLogger(logger) // .AddSkip(extraSkip))
	return logger
}

type Opt func(logger *logrus.Logger)

func WithLoggingFormat(format string) Opt {
	return func(logger *logrus.Logger) {
		logex.SetupLoggingFormat(format, extraSkip)
	}
}

func initLogger(config *log.LoggerConfig) *logrus.Logger {
	var ll log.Level
	ll, _ = log.ParseLevel(config.Level)
	if ll == log.OffLevel {
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.SetOutput(ioutil.Discard)
		return logrus.New()
	}

	var err error
	if config.Target == "file" {
		logrus.SetLevel(logrus.Level(ll))

		logrus.SetFormatter(&formatter.TextFormatter{ForceColors: true})
		logrus.SetReportCaller(true)

		var file *os.File
		fPath := path.Join(os.ExpandEnv(config.Directory), "output.log")
		fDir := path.Dir(fPath)
		err = exec.EnsureDir(fDir)
		if err != nil {
			fmt.Printf(`

You're been prompt with a "sudo" requesting because this folder was been creating but need more privileges:

- %v

We must have created the logging output file in it.

`, fDir)
			err = exec.EnsureDirEnh(fDir)
		}

		if err == nil {
			if file, err = os.OpenFile(fPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660); err == nil {
				logrus.SetOutput(file)
				return logrus.StandardLogger()
			}
		}

		logrus.Warnf("Failed to log to file %q, using default stderr", fPath)
	}

	// setupLoggingFormat(format, 0)
	logex.EnableWith(ll)

	format := "text" // cmdr.GetStringR("logger.format", "text")
	logex.SetupLoggingFormat(format, extraSkip)

	logger := logrus.StandardLogger()
	// logger.Infof("hello, logLevel = %q", logLevel)
	// logrus.Infof("hello, logLevel = %q", logLevel)
	return logger
}

// const extraSkip = 1

// func setupLoggingFormat(format string, logexSkipFrames int) {
//	switch format {
//	case "json":
//		logrus.SetFormatter(&logrus.JSONFormatter{
//			TimestampFormat:  "2006-01-02 15:04:05.000",
//			DisableTimestamp: false,
//			PrettyPrint:      false,
//		})
//	default:
//		e := false
//		if logexSkipFrames > 0 {
//			e = true
//		}
//		logrus.SetFormatter(&formatter.TextFormatter{
//			ForceColors:               true,
//			DisableColors:             false,
//			FullTimestamp:             true,
//			TimestampFormat:           "2006-01-02 15:04:05.000",
//			Skip:                      logexSkipFrames,
//			EnableSkip:                e,
//			EnvironmentOverrideColors: true,
//		})
//	}
// }
