// Copyright © 2019 Hedzr Yeh.

package cmdr

import (
	"bytes"
	"fmt"
	"gopkg.in/hedzr/errors.v2"
	"reflect"
)

var (
	// ErrShouldBeStopException tips `Exec()` cancelled the following actions after `PreAction()`
	ErrShouldBeStopException = newErrorWithMsg("stop me right now")

	// ErrBadArg is a generic error for user
	ErrBadArg = newErrorWithMsg("bad argument")

	errWrongEnumValue = newErrTmpl("unexpected enumerable value '%s' for option '%s', under command '%s'")
)

// ErrorForCmdr structure
type ErrorForCmdr struct {
	// Inner     error
	// Msg       string
	Ignorable bool
	causer    error
	msg       string
	livedArgs []interface{}
}

// newError formats a ErrorForCmdr object
func newError(ignorable bool, srcTemplate error, livedArgs ...interface{}) error {
	// log.Printf("--- newError: sourceTemplate args: %v", livedArgs)
	var e error
	switch v := srcTemplate.(type) {
	case *ErrorForCmdr:
		e = v.FormatNew(ignorable, livedArgs...)
	case *errors.WithStackInfo:
		if ex, ok := v.Cause().(*ErrorForCmdr); ok {
			e = ex.FormatNew(ignorable, livedArgs...)
		}
	}
	return e
	// if len(args) > 0 {
	// 	return &ErrorForCmdr{Inner: nil, Ignorable: ignorable, Msg: fmt.Sprintf(inner.Error(), args...)}
	// }
	// return &ErrorForCmdr{Inner: inner, Ignorable: ignorable}
}

// newErrorWithMsg formats a ErrorForCmdr object
func newErrorWithMsg(msg string, inners ...error) error {
	return newErr(msg).Attach(inners...)
	// return &ErrorForCmdr{Inner: inner, Ignorable: false, Msg: msg}
}

// func (s *ErrorForCmdr) Error() string {
// 	if s.Inner != nil {
// 		return fmt.Sprintf("Error: %v. Inner: %v", s.Msg, s.Inner.Error())
// 	}
// 	return s.Msg
// }

// newErr creates a *errors.WithStackInfo object
func newErr(msg string, args ...interface{}) *errors.WithStackInfo {
	// return &ErrorForCmdr{ExtErr: *errors.New(msg, args...)}
	return withIgnorable(false, nil, msg, args...).(*errors.WithStackInfo)
}

// newErrTmpl creates a *errors.WithStackInfo object
func newErrTmpl(tmpl string) *errors.WithStackInfo {
	return withIgnorable(false, nil, tmpl).(*errors.WithStackInfo)
}

// withIgnorable formats a wrapped error object with error code.
func withIgnorable(ignorable bool, err error, message string, args ...interface{}) error {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	err = &ErrorForCmdr{
		Ignorable: ignorable,
		causer:    err,
		msg:       message,
	}
	//n := errors.WithStack(err).(*errors.WithStackInfo)
	//return n.SetCause(err)
	return errors.WithStack(err)
}

func (w *ErrorForCmdr) Error() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprint(w.Ignorable))
	if len(w.msg) > 0 {
		buf.WriteRune('|')
		buf.WriteString(w.msg)
	}
	if w.causer != nil {
		buf.WriteRune('|')
		buf.WriteString(w.causer.Error())
	}
	return buf.String()
}

// FormatNew creates a new error object based on this error template 'w'.
//
// Example:
//
// 	   errTmpl1001 := BUG1001.NewTemplate("something is wrong %v")
// 	   err4 := errTmpl1001.FormatNew("ok").Attach(errBug1)
// 	   fmt.Println(err4)
// 	   fmt.Printf("%+v\n", err4)
//
func (w *ErrorForCmdr) FormatNew(ignorable bool, livedArgs ...interface{}) *errors.WithStackInfo {
	x := withIgnorable(ignorable, w.causer, w.msg, livedArgs...).(*errors.WithStackInfo)
	x.Cause().(*ErrorForCmdr).livedArgs = livedArgs
	return x
}

// Attach appends errs.
// For ErrorForCmdr, only one last error will be kept.
func (w *ErrorForCmdr) Attach(errs ...error) {
	for _, err := range errs {
		w.causer = err
	}
}

// Cause returns the underlying cause of the error recursively,
// if possible.
func (w *ErrorForCmdr) Cause() error {
	return w.causer
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func (w *ErrorForCmdr) Unwrap() error {
	return w.causer
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
func (w *ErrorForCmdr) As(target interface{}) bool {
	if target == nil {
		panic("errors: target cannot be nil")
	}
	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}
	if e := typ.Elem(); e.Kind() != reflect.Interface && !e.Implements(errorType) {
		panic("errors: *target must be interface or implement error")
	}
	targetType := typ.Elem()
	err := w.causer
	for err != nil {
		if reflect.TypeOf(err).AssignableTo(targetType) {
			val.Elem().Set(reflect.ValueOf(err))
			return true
		}
		if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}

// Is reports whether any error in err's chain matches target.
func (w *ErrorForCmdr) Is(target error) bool {
	if target == nil {
		return w.causer == target
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable && w.causer == target {
			return true
		}
		if x, ok := w.causer.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}
		// TODO: consider supporing target.Is(err). This would allow
		// user-definable predicates, but also may allow for coping with sloppy
		// APIs, thereby making it easier to get away with them.
		if err := errors.Unwrap(w.causer); err == nil {
			return false
		}
	}
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()
