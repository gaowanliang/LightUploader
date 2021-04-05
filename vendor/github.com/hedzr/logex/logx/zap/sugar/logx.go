package sugar

import (
	"github.com/hedzr/log"
	"go.uber.org/zap"
	"io"
)

type dzl struct {
	Logger *zap.SugaredLogger
}

func (s *dzl) AddSkip(skip int) log.Logger {
	return s
}

func (s *dzl) Tracef(msg string, args ...interface{}) {
	if log.GetTraceMode() {
		s.Logger.Debugf(msg, args...)
	}
}

func (s *dzl) Debugf(msg string, args ...interface{}) {
	if log.GetDebugMode() {
		s.Logger.Debugf(msg, args...)
	}
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
		// s.Logger.Debugw(msg, fields...)
	}
}

func (s *dzl) Debug(args ...interface{}) {
	if log.GetDebugMode() {
		// s.Logger.Debugw(msg, fields...)
	}
}

func (s *dzl) Info(args ...interface{}) {
	// s.Logger.Infow(msg, fields...)
}

func (s *dzl) Warn(args ...interface{}) {
	// s.Logger.Warnw(msg, fields...)
}

func (s *dzl) Error(args ...interface{}) {
	// s.Logger.Errorw(msg, fields...)
}

func (s *dzl) Fatal(args ...interface{}) {
	// s.Logger.Fatalw(msg, fields...)
}

func (s *dzl) Print(args ...zap.Field) {
	// s.Logger.Infow(msg, fields...)
}

//
//

func (s *dzl) SetLevel(lvl log.Level) {
	// panic("implement me")
}

func (s *dzl) GetLevel() log.Level {
	// panic("implement me")
	return log.DebugLevel
}

func (s *dzl) SetOutput(out io.Writer) {
	// s.Logger.Out = out
}

func (s *dzl) GetOutput() (out io.Writer) {
	return
}

func (s *dzl) Setup() {
	// initLogger("", "")
}

// func (s *dzl) AsFieldLogger() logx.FieldLogger {
//	return s
// }

func AsFieldLogger(s log.Logger) FieldLogger {
	if l, ok := s.(FieldLogger); ok {
		return l
	}
	return nil
}

type FieldLogger interface {
	log.Logger
	Trace(args ...zap.Field)
	Debug(args ...zap.Field)
	Info(args ...zap.Field)
	Warn(args ...zap.Field)
	Error(args ...zap.Field)
	Fatal(args ...zap.Field)
	Print(args ...zap.Field)
}
