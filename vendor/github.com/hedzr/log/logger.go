// Copyright © 2020 Hedzr Yeh.

// Package log provide the standard interface of logging for what any go
// libraries want strip off the direct dependency from a known logging
// library.
package log

import (
	"io"
	"os"
	"strings"
)

type (
	L interface {
		// Trace prints all args to stdin if logging level is greater than TraceLevel
		Trace(args ...interface{})
		// Debug prints all args to stdin if logging level is greater than DebugLevel
		Debug(args ...interface{})
		// Info prints all args to stdin if logging level is greater than InfoLevel
		Info(args ...interface{})
		// Warn prints all args to stderr
		Warn(args ...interface{})
		// Error prints all args to stderr
		Error(args ...interface{})
		// Fatal is equivalent to Printf() followed by a call to os.Exit(1).
		Fatal(args ...interface{})
		// Panic is equivalent to Printf() followed by a call to panic().
		Panic(args ...interface{})
		// Print calls Output to print to the standard logger.
		// Arguments are handled in the manner of fmt.Print.
		Print(args ...interface{})
		// Println calls Output to print to the standard logger.
		// Arguments are handled in the manner of fmt.Println.
		Println(args ...interface{})
	}

	LF interface {
		// Tracef prints the text to stdin if logging level is greater than TraceLevel
		Tracef(msg string, args ...interface{})
		// Debugf prints the text to stdin if logging level is greater than DebugLevel
		Debugf(msg string, args ...interface{})
		// Infof prints the text to stdin if logging level is greater than InfoLevel
		Infof(msg string, args ...interface{})
		// Warnf prints the text to stderr
		Warnf(msg string, args ...interface{})
		// Errorf prints the text to stderr
		Errorf(msg string, args ...interface{})
		// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
		Fatalf(msg string, args ...interface{})
		// Panicf is equivalent to Printf() followed by a call to panic().
		Panicf(msg string, args ...interface{})
		// Printf calls Output to print to the standard logger.
		// Arguments are handled in the manner of fmt.Printf.
		Printf(msg string, args ...interface{})
	}

	// Logger is a minimal logger with no more dependencies
	Logger interface {
		LF

		// SetLevel sets the logging level
		SetLevel(lvl Level)
		// GetLevel returns the current logging level
		GetLevel() Level
		// SetOutput setup the logging output device
		SetOutput(out io.Writer)
		// SetOutput returns the current logging output device
		GetOutput() (out io.Writer)

		// Setup will be invoked once an instance created
		Setup()

		// AddSkip adds an extra count to skip stack frames
		AddSkip(skip int) Logger

		// AsFieldLogger() FieldLogger
	}

	// LoggerExt is a minimal logger with no more dependencies
	LoggerExt interface {
		L
		Logger
	}

	// LoggerConfig is used for creating a minimal logger with no more dependencies
	LoggerConfig struct {
		Enabled          bool
		Backend          string // zap, sugar, logrus
		Level            string // level
		Format           string // text, json, ...
		Target           string // console, file, console+file
		Directory        string // logdir, for file
		AllToErrorDevice bool   //
		DebugMode        bool   `json:"-" yaml:"-"`
		TraceMode        bool   `json:"-" yaml:"-"`

		// the following options are copied from zap rotator

		// MaxSize is the maximum size in megabytes of the log file before it gets
		// rotated. It defaults to 100 megabytes.
		MaxSize int `json:"maxsize" yaml:"maxsize"`

		// MaxAge is the maximum number of days to retain old log files based on the
		// timestamp encoded in their filename.  Note that a day is defined as 24
		// hours and may not exactly correspond to calendar days due to daylight
		// savings, leap seconds, etc. The default is not to remove old log files
		// based on age.
		MaxAge int `json:"maxage" yaml:"maxage"`

		// MaxBackups is the maximum number of old log files to retain.  The default
		// is to retain all old log files (though MaxAge may still cause them to get
		// deleted.)
		MaxBackups int `json:"maxbackups" yaml:"maxbackups"`

		// LocalTime determines if the time used for formatting the timestamps in
		// backup files is the computer's local time.  The default is to use UTC
		// time.
		LocalTime bool `json:"localtime" yaml:"localtime"`

		// Compress determines if the rotated log files should be compressed
		// using gzip. The default is not to perform compression.
		Compress bool `json:"compress" yaml:"compress"`
	}

	// BuilderFunc provides a function prototype for creating a hedzr/log & hedzr/logex -compliant creator.
	BuilderFunc func(config *LoggerConfig) (logger Logger)
)

// InTesting detects whether is running under go test mode
func InTesting() bool { return InTestingT(os.Args) }

// InTestingT detects whether is running under go test mode
func InTestingT(args []string) bool {
	if !strings.HasSuffix(args[0], ".test") &&
		!strings.Contains(args[0], "/T/___Test") {

		// [0] = /var/folders/td/2475l44j4n3dcjhqbmf3p5l40000gq/T/go-build328292371/b001/exe/main
		// !strings.Contains(SavedOsArgs[0], "/T/go-build")

		for _, s := range args {
			if s == "-test.v" || s == "-test.run" {
				return true
			}
		}
		return false

	}
	return true
}

// AsL converts a logger to L type (with Info(...), ... prototypes)
func AsL(logger LF) L {
	if l, ok := logger.(L); ok {
		//if l1, ok := l.(Logger); ok {
		//	return l1.AddSkip(1).(L)
		//}
		return l
	}
	return nil
}

// AsLogger converts a logger to LF or Logger type (with Infof(...), ... prototypes)
func AsLogger(logger L) Logger {
	if l, ok := logger.(Logger); ok {
		return l // .AddSkip(1)
	}
	return nil
}

// NewLoggerConfig returns a default LoggerConfig
func NewLoggerConfig() *LoggerConfig {
	return NewLoggerConfigWith(true, "sugar", "info")
}

// NewLoggerConfigWith returns a default LoggerConfig
func NewLoggerConfigWith(enabled bool, backend, level string) *LoggerConfig {
	var dm, tm = GetDebugMode(), GetTraceMode()
	if dm {
		level = "debug"
	}
	if tm {
		level = "trace"
	}
	var l Level
	l, _ = ParseLevel(level)
	SetDebugMode(l >= DebugLevel)
	SetTraceMode(l >= TraceLevel)
	dm, tm = GetDebugMode(), GetTraceMode()
	return &LoggerConfig{
		Enabled:   enabled,
		Backend:   backend,
		Level:     level,
		Format:    "text",
		Target:    "console",
		Directory: "/var/log",
		DebugMode: dm,
		TraceMode: tm,

		MaxSize:    1024, // megabytes
		MaxBackups: 3,    // 3 backups kept at most
		MaxAge:     7,    // 7 days kept at most
		Compress:   true, // disabled by default
	}
}

// SetLevel sets the logging level
func SetLevel(l Level) { logger.SetLevel(l) }

// GetLevel returns the current logging level
func GetLevel() Level { return logger.GetLevel() }

var logger = NewStdLogger()

// SetOutput setup the logging output device
func SetOutput(w io.Writer) { logger.SetOutput(w) }

// SetLogger transfer an instance into log package-level value
func SetLogger(l Logger) { l.SetLevel(logger.GetLevel()); logger = l }

// GetLogger returns the package-level logger globally
func GetLogger() Logger { return logger }

// Tracef prints the text to stdin if logging level is greater than TraceLevel
func Tracef(msg string, args ...interface{}) {
	logger.Tracef(msg, args...)
}

// Debugf prints the text to stdin if logging level is greater than DebugLevel
func Debugf(msg string, args ...interface{}) {
	logger.Debugf(msg, args...)
}

// Infof prints the text to stdin if logging level is greater than InfoLevel
func Infof(msg string, args ...interface{}) {
	logger.Infof(msg, args...)
}

// Warnf prints the text to stderr
func Warnf(msg string, args ...interface{}) {
	logger.Warnf(msg, args...)
}

// Errorf prints the text to stderr
func Errorf(msg string, args ...interface{}) {
	logger.Errorf(msg, args...)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(msg string, args ...interface{}) {
	if InTesting() {
		logger.Panicf(msg, args)
	}
	logger.Fatalf(msg, args...)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(msg string, args ...interface{}) {
	logger.Panicf(msg, args...)
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(msg string, args ...interface{}) {
	logger.Printf(msg, args...)
}

// Trace prints all args to stdin if logging level is greater than TraceLevel
func Trace(args ...interface{}) {
	if l := AsL(logger); l != nil {
		l.Trace(args...)
	}
}

// Debug prints all args to stdin if logging level is greater than DebugLevel
func Debug(args ...interface{}) {
	if l := AsL(logger); l != nil {
		l.Debug(args...)
	}
}

// Info prints all args to stdin if logging level is greater than InfoLevel
func Info(args ...interface{}) {
	if l := AsL(logger); l != nil {
		l.Info(args...)
	}
}

// Warn prints all args to stderr
func Warn(args ...interface{}) {
	if l := AsL(logger); l != nil {
		l.Warn(args...)
	}
}

// Error prints all args to stderr
func Error(args ...interface{}) {
	if l := AsL(logger); l != nil {
		l.Error(args...)
	}
}

// Fatal is equivalent to Printf() followed by a call to os.Exit(1).
func Fatal(args ...interface{}) {
	if l := AsL(logger); l != nil {
		if InTesting() {
			l.Panic(args)
		}
		l.Fatal(args...)
	}
}

// Panic is equivalent to Printf() followed by a call to panic().
func Panic(args ...interface{}) {
	if l := AsL(logger); l != nil {
		l.Panic(args...)
	}
}

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(args ...interface{}) {
	if l := AsL(logger); l != nil {
		l.Print(args...)
	}
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(args ...interface{}) {
	if l := AsL(logger); l != nil {
		l.Println(args...)
	}
}
