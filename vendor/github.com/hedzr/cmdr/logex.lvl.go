// Copyright © 2020 Hedzr Yeh.

package cmdr

import (
	"fmt"
	"strings"
)

//
// The following codes copied and modified from hedzr/log
//

// Level type
type Level uint32

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level Level) String() string {
	if b, err := level.MarshalText(); err == nil {
		return string(b)
	}
	return "unknown"
}

// ParseLevel takes a string level and returns the hedzr/log logging level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug", "devel", "dev", "develop":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	case "off", "no", "disable":
		return OffLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid logger Level: %q", lvl)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (level *Level) UnmarshalText(text []byte) error {
	l, err := ParseLevel(string(text))
	if err != nil {
		return err
	}

	*level = Level(l)

	return nil
}

// Available level names are:
// "disable"
// "fatal"
// "error"
// "warn"
// "info"
// "debug"

// MarshalText convert Level to string and []byte
func (level Level) MarshalText() ([]byte, error) {
	switch level {
	case OffLevel:
		return []byte("off"), nil
	case TraceLevel:
		return []byte("trace"), nil
	case DebugLevel:
		return []byte("debug"), nil
	case InfoLevel:
		return []byte("info"), nil
	case WarnLevel:
		return []byte("warning"), nil
	case ErrorLevel:
		return []byte("error"), nil
	case FatalLevel:
		return []byte("fatal"), nil
	case PanicLevel:
		return []byte("panic"), nil
	}

	return nil, fmt.Errorf("not a valid logger level %d", level)
}

// AllLevels is a constant exposing all logging levels
var AllLevels = []Level{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
	OffLevel,
}

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `log.NewDummyLogger()`/`log.NewStdLogger()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
	// OffLevel level. The logger will be shutdown.
	OffLevel
)
