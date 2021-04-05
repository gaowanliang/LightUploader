package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

// NewStdLogger return a stdlib `log` logger
func NewStdLogger() Logger {
	return &stdLogger{Level: InfoLevel}
}

// NewStdLoggerWith return a stdlib `log` logger
func NewStdLoggerWith(lvl Level) Logger {
	return &stdLogger{Level: lvl}
}

// NewStdLoggerWithConfig return a stdlib `log` logger
func NewStdLoggerWithConfig(config *LoggerConfig) Logger {
	l, _ := ParseLevel(config.Level)
	return &stdLogger{Level: l}
}

type stdLogger struct {
	Level
}

func (s *stdLogger) out(args ...interface{}) {
	str := fmt.Sprint(args...)
	_ = log.Output(2, str)
}

func (s *stdLogger) Trace(args ...interface{}) {
	if s.Level >= TraceLevel {
		s.out(args...)
	}
}

func (s *stdLogger) Debug(args ...interface{}) {
	if s.Level >= DebugLevel {
		s.out(args...)
	}
}

func (s *stdLogger) Info(args ...interface{}) {
	if s.Level >= InfoLevel {
		s.out(args...)
	}
}

func (s *stdLogger) Warn(args ...interface{}) {
	s.out(args...)
}

func (s *stdLogger) Error(args ...interface{}) {
	s.out(args...)
}

func (s *stdLogger) Fatal(args ...interface{}) {
	s.out(args...)
	if InTesting() {
		panic(fmt.Sprint(args...))
	}
	os.Exit(1)
}

func (s *stdLogger) Panic(args ...interface{}) {
	str := fmt.Sprint(args...)
	_ = log.Output(2, str)
	panic(str)
}

func (s *stdLogger) Print(args ...interface{}) {
	s.out(args...)
}

func (s *stdLogger) Println(args ...interface{}) {
	str := fmt.Sprintln(args...)
	_ = log.Output(2, str)
}

func (s *stdLogger) outf(msg string, args ...interface{}) {
	str := fmt.Sprintf(msg, args...)
	_ = log.Output(2, str)
}

func (s *stdLogger) Tracef(msg string, args ...interface{}) {
	if s.Level >= TraceLevel {
		s.outf(msg, args...)
	}
}

func (s *stdLogger) Debugf(msg string, args ...interface{}) {
	if s.Level >= DebugLevel {
		s.outf(msg, args...)
	}
}

func (s *stdLogger) Infof(msg string, args ...interface{}) {
	if s.Level >= InfoLevel {
		s.outf(msg, args...)
	}
}

func (s *stdLogger) Warnf(msg string, args ...interface{}) {
	s.outf(msg, args...)
}

func (s *stdLogger) Errorf(msg string, args ...interface{}) {
	s.outf(msg, args...)
}

func (s *stdLogger) Fatalf(msg string, args ...interface{}) {
	s.outf(msg, args...)
	if InTesting() {
		panic(fmt.Sprintf(msg, args...))
	}
	os.Exit(1)
}

func (s *stdLogger) Panicf(msg string, args ...interface{}) {
	str := fmt.Sprintf(msg, args...)
	_ = log.Output(2, str)
	panic(str)
}

func (s *stdLogger) Printf(msg string, args ...interface{}) {
	s.outf(msg, args...)
}

func (s *stdLogger) SetLevel(lvl Level)         { s.Level = lvl }
func (s *stdLogger) GetLevel() Level            { return s.Level }
func (d *stdLogger) SetOutput(out io.Writer)    {}
func (d *stdLogger) GetOutput() (out io.Writer) { return }
func (s *stdLogger) Setup()                     {}
func (d *stdLogger) AddSkip(skip int) Logger    { return d }
