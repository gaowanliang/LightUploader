package logrus

import (
	"github.com/hedzr/log"
	"github.com/sirupsen/logrus"
	"io"
)

type entry struct {
	*logrus.Entry
}

func (s *entry) SetLevel(lvl log.Level)     { s.Logger.SetLevel(logrus.Level(lvl)) }
func (s *entry) GetLevel() log.Level        { return log.Level(s.Logger.Level) }
func (s *entry) SetOutput(out io.Writer)    { s.Logger.Out = out }
func (s *entry) GetOutput() (out io.Writer) { return s.Logger.Out }
func (s *entry) Setup()                     {}

func (s *entry) AddSkip(skip int) log.Logger {
	return &entry{
		s.Entry.WithField("SKIP", skip),
	}
}

//
//

func (s *entry) Tracef(msg string, args ...interface{}) {
	if log.GetTraceMode() {
		s.Entry.Tracef(msg, args...)
	}
}

func (s *entry) Debugf(msg string, args ...interface{}) {
	s.Entry.Debugf(msg, args...)
}

func (s *entry) Infof(msg string, args ...interface{}) {
	s.Entry.Infof(msg, args...)
}

func (s *entry) Warnf(msg string, args ...interface{}) {
	s.Entry.Warnf(msg, args...)
}

func (s *entry) Errorf(msg string, args ...interface{}) {
	s.Entry.Errorf(msg, args...)
}

func (s *entry) Fatalf(msg string, args ...interface{}) {
	s.Entry.Fatalf(msg, args...)
}

func (s *entry) Panicf(msg string, args ...interface{}) {
	s.Entry.Panicf(msg, args...)
}

func (s *entry) Printf(msg string, args ...interface{}) {
	s.Entry.Infof(msg, args...)
}

//
//

type dzl struct {
	*logrus.Logger
	Config *log.LoggerConfig
}

func (s *dzl) AddSkip(skip int) log.Logger {
	return &entry{
		s.Logger.WithField("SKIP", skip),
	}
}

func (s *dzl) Tracef(msg string, args ...interface{}) {
	if log.GetTraceMode() {
		s.Logger.Tracef(msg, args...)
	}
}

func (s *dzl) Debugf(msg string, args ...interface{}) {
	s.Logger.Debugf(msg, args...)
}

func (s *dzl) Infof(msg string, args ...interface{}) {
	s.Logger.Infof(msg, args...)
}

func (s *dzl) Warnf(msg string, args ...interface{}) {
	s.Logger.Warnf(msg, args...)
}

func (s *dzl) Errorf(msg string, args ...interface{}) {
	s.Logger.Errorf(msg, args...)
}

func (s *dzl) Fatalf(msg string, args ...interface{}) {
	s.Logger.Fatalf(msg, args...)
}

func (s *dzl) Panicf(msg string, args ...interface{}) {
	s.Logger.Panicf(msg, args...)
}

func (s *dzl) Printf(msg string, args ...interface{}) {
	s.Logger.Infof(msg, args...)
}

//
//

func (s *dzl) Trace(args ...interface{}) {
	if log.GetTraceMode() {
		s.Logger.Trace(args...)
	}
}

func (s *dzl) Debug(args ...interface{}) {
	s.Logger.Debug(args...)
}

func (s *dzl) Info(args ...interface{}) {
	s.Logger.Info(args...)
}

func (s *dzl) Warn(args ...interface{}) {
	s.Logger.Warn(args...)
}

func (s *dzl) Error(args ...interface{}) {
	s.Logger.Error(args...)
}

func (s *dzl) Fatal(args ...interface{}) {
	s.Logger.Fatal(args...)
}

func (s *dzl) Print(args ...interface{}) {
	s.Logger.Print(args...)
}

//
//

func (s *dzl) SetLevel(lvl log.Level)     { s.Logger.SetLevel(logrus.Level(lvl)) }
func (s *dzl) GetLevel() log.Level        { return log.Level(s.Logger.Level) }
func (s *dzl) SetOutput(out io.Writer)    { s.Logger.Out = out }
func (s *dzl) GetOutput() (out io.Writer) { return s.Logger.Out }
func (s *dzl) Setup()                     {}

// func (s *dzl) AsFieldLogger() logx.FieldLogger {
//	return s
// }
