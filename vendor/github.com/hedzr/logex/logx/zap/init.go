package zap

import (
	"fmt"
	"github.com/hedzr/log"
	"github.com/hedzr/log/exec"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	log2 "log"
	"os"
	"path"
)

// New create a sugared zap logger
//
// level can be: "disable", "panic", "fatal", "error", "warn", "info", "debug", "trace"
//
func New(level string, traceMode, debugMode bool, opts ...Opt) log.Logger {
	log.SetTraceMode(traceMode)
	log.SetDebugMode(debugMode)
	// ll := cmdr.GetStringR("logger.level", "info")
	lvl, _ := log.ParseLevel(level)
	if log.GetDebugMode() {
		if lvl < log.DebugLevel {
			lvl = log.DebugLevel
			level = "debug"
		}
	}
	if log.GetTraceMode() {
		if lvl < log.TraceLevel {
			lvl = log.TraceLevel
			level = "debug"
		}
	}

	zl := initLogger(log.NewLoggerConfig())

	for _, opt := range opts {
		opt(zl)
	}

	logger := &dzl{zl, zl.Sugar()}
	logger.Setup()
	log.SetLogger(logger)
	return logger
}

// NewWithConfigSimple create a sugared zap sugared logger
func NewWithConfigSimple(config *log.LoggerConfig) log.Logger { return NewWithConfig(config) }

// NewWithConfig create a sugared zap sugared logger
//
// level can be: "disable", "panic", "fatal", "error", "warn", "info", "debug", "trace"
//
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
			config.Level = "debug" // zap hasn't `trace` level
		}
	}

	zl := initLogger(config)

	for _, opt := range opts {
		opt(zl)
	}

	logger := &dzl{zl, zl.Sugar()}
	logger.Setup()
	log.SetLogger(logger)
	return logger
}

type Opt func(logger *zap.Logger)

func initLogger(config *log.LoggerConfig) *zap.Logger {
	var level zapcore.Level
	_ = level.Set(config.Level)

	if config.Target == "file" {
		var w zapcore.WriteSyncer

		fPath := path.Join(os.ExpandEnv(config.Directory), "output.log")
		fDir := path.Dir(fPath)
		err := exec.EnsureDir(fDir)
		if err != nil {
			fmt.Printf(`

You're been prompt with a "sudo" requesting because this folder was been creating but need more privileges:

- %v

We must have created the logging output file in it.

`, fDir)
			err = exec.EnsureDirEnh(fDir)
		}

		if err != nil {
			log2.Printf("cannot create logging dir %q, error: %v", fDir, err)
			return nil
		}

		hook := lumberjack.Logger{
			Filename:   fPath,             // the logging file path
			MaxSize:    config.MaxSize,    // megabytes
			MaxBackups: config.MaxBackups, // 3 backups kept at most
			MaxAge:     config.MaxAge,     // 7 days kept at most
			Compress:   config.Compress,   // disabled by default
		}
		w = zapcore.AddSync(&hook)

		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			w,
			level,
		)
		logger := zap.New(core)
		return logger.WithOptions(zap.AddCallerSkip(extraSkip)) // .Sugar()

	} else {
		logCfg := zap.NewDevelopmentConfig()
		logCfg.Level = zap.NewAtomicLevelAt(level)
		logCfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		logCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logCfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		logger, _ := logCfg.Build()
		return logger.WithOptions(zap.AddCallerSkip(extraSkip)) // .Sugar()
	}

}

const extraSkip = 2

func initLoggerConsole(logLevel zapcore.Level) *zap.Logger {
	// alevel := zap.NewAtomicLevel()
	// http.HandleFunc("/handle/level", alevel.ServeHTTP)
	// logCfg.Level = alevel

	logCfg := zap.NewDevelopmentConfig()
	logCfg.Level = zap.NewAtomicLevelAt(logLevel)
	logCfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	logger, _ := logCfg.Build()
	return logger
}

// func initLogger(logPath string, logLevel string) *zap.Logger {
//	hook := lumberjack.Logger{
//		Filename:   logPath, // the logging file path
//		MaxSize:    1024,    // megabytes
//		MaxBackups: 3,       // 3 backups kept at most
//		MaxAge:     7,       // 7 days kept at most
//		Compress:   true,    // disabled by default
//	}
//	w := zapcore.AddSync(&hook)
//
//	var level zapcore.Level
//	switch logLevel {
//	case "debug":
//		level = zap.DebugLevel
//	case "info":
//		level = zap.InfoLevel
//	case "error":
//		level = zap.ErrorLevel
//	default:
//		level = zap.InfoLevel
//	}
//
//	encoderConfig := zap.NewProductionEncoderConfig()
//	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
//	core := zapcore.NewCore(
//		zapcore.NewConsoleEncoder(encoderConfig),
//		w,
//		level,
//	)
//	logger := zap.New(core)
//	return logger
// }
//
// var sugarLogger *zap.SugaredLogger
//
// func mainEntry() {
//	initOneLogger()
//	defer sugarLogger.Sync()
//	simpleHttpGet("www.sogo.com")
//	simpleHttpGet("http://www.sogo.com")
// }
//
// func initOneLogger() {
//	encoder := getEncoder()
//	writeSyncer := getLogWriter()
//	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
//
//	logger := zap.New(core, zap.AddCaller())
//	sugarLogger = logger.Sugar()
// }
//
// func getEncoder() zapcore.Encoder {
//	encoderConfig := zap.NewProductionEncoderConfig()
//	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
//	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
//	return zapcore.NewConsoleEncoder(encoderConfig)
// }
//
// func getLogWriter() zapcore.WriteSyncer {
//	lumberJackLogger := &lumberjack.Logger{
//		Filename:   "./test.log",
//		MaxSize:    1,
//		MaxBackups: 5,
//		MaxAge:     30,
//		Compress:   false,
//	}
//	return zapcore.AddSync(lumberJackLogger)
// }
//
// func simpleHttpGet(url string) {
//	sugarLogger.Debugf("Trying to hit GET request for %s", url)
//	sugarLogger.Info("Success..",
//		zap.String("statusCode", "200"),
//		zap.String("url", "https://z.cn"))
//	resp, err := http.Get(url)
//	if err != nil {
//		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
//	} else {
//		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
//		resp.Body.Close()
//	}
// }
