// Copyright © 2020 Hedzr Yeh.

package log

import (
	"fmt"
	"io"
	"os"
)

// NewDummyLogger return a dummy logger
func NewDummyLogger() Logger {
	return &dummyLogger{}
}

// NewDummyLoggerWithConfig return a dummy logger
func NewDummyLoggerWithConfig(config *LoggerConfig) Logger {
	return &dummyLogger{}
}

type dummyLogger struct{}

func (d *dummyLogger) Trace(args ...interface{}) {}
func (d *dummyLogger) Debug(args ...interface{}) {}
func (d *dummyLogger) Info(args ...interface{})  {}
func (d *dummyLogger) Warn(args ...interface{})  {}
func (d *dummyLogger) Error(args ...interface{}) {}

func (d *dummyLogger) Fatal(args ...interface{}) {
	if InTesting() {
		panic(fmt.Sprint(args...))
	}
	os.Exit(1)
}

func (d *dummyLogger) Panic(args ...interface{})              { panic(fmt.Sprint(args...)) }
func (d *dummyLogger) Print(args ...interface{})              {}
func (d *dummyLogger) Println(args ...interface{})            {}
func (d *dummyLogger) Tracef(msg string, args ...interface{}) {}
func (d *dummyLogger) Debugf(msg string, args ...interface{}) {}
func (d *dummyLogger) Infof(msg string, args ...interface{})  {}
func (d *dummyLogger) Warnf(msg string, args ...interface{})  {}
func (d *dummyLogger) Errorf(msg string, args ...interface{}) {}

func (d *dummyLogger) Fatalf(msg string, args ...interface{}) {
	// panic("implement me")
	if InTesting() {
		panic(fmt.Sprintf(msg, args...))
	}
	os.Exit(1)
}

func (d *dummyLogger) Panicf(msg string, args ...interface{}) { panic(fmt.Sprintf(msg, args...)) }
func (d *dummyLogger) Printf(msg string, args ...interface{}) {}
func (d *dummyLogger) SetLevel(lvl Level)                     {}
func (d *dummyLogger) GetLevel() Level                        { return InfoLevel }
func (d *dummyLogger) SetOutput(out io.Writer)                {}
func (d *dummyLogger) GetOutput() (out io.Writer)             { return }
func (d *dummyLogger) Setup()                                 {}
func (d *dummyLogger) AddSkip(skip int) Logger                { return d }

//
//
//
