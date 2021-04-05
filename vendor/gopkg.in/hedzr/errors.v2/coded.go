package errors

import (
	"bytes"
	"fmt"
	"io"
)

// A Code is an signed 32-bit error code copied from gRPC spec but negatived.
type Code int32

const (
	// OK is returned on success.
	OK Code = 0

	// Canceled indicates the operation was canceled (typically by the caller).
	Canceled Code = -1

	// Unknown error. An example of where this error may be returned is
	// if a Status value received from another address space belongs to
	// an error-space that is not known in this address space. Also
	// errors raised by APIs that do not return enough error information
	// may be converted to this error.
	Unknown Code = -2

	// InvalidArgument indicates client specified an invalid argument.
	// Note that this differs from FailedPrecondition. It indicates arguments
	// that are problematic regardless of the state of the system
	// (e.g., a malformed file name).
	InvalidArgument Code = -3

	// DeadlineExceeded means operation expired before completion.
	// For operations that change the state of the system, this error may be
	// returned even if the operation has completed successfully. For
	// example, a successful response from a server could have been delayed
	// long enough for the deadline to expire.
	//
	// = HTTP 408 Timeout
	DeadlineExceeded Code = -4

	// NotFound means some requested entity (e.g., file or directory) was
	// not found.
	//
	// = HTTP 404
	NotFound Code = -5

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	AlreadyExists Code = -6

	// PermissionDenied indicates the caller does not have permission to
	// execute the specified operation. It must not be used for rejections
	// caused by exhausting some resource (use ResourceExhausted
	// instead for those errors). It must not be
	// used if the caller cannot be identified (use Unauthenticated
	// instead for those errors).
	PermissionDenied Code = -7

	// ResourceExhausted indicates some resource has been exhausted, perhaps
	// a per-user quota, or perhaps the entire file system is out of space.
	ResourceExhausted Code = -8

	// FailedPrecondition indicates operation was rejected because the
	// system is not in a state required for the operation's execution.
	// For example, directory to be deleted may be non-empty, an rmdir
	// operation is applied to a non-directory, etc.
	//
	// A litmus test that may help a service implementor in deciding
	// between FailedPrecondition, Aborted, and Unavailable:
	//  (a) Use Unavailable if the client can retry just the failing call.
	//  (b) Use Aborted if the client should retry at a higher-level
	//      (e.g., restarting a read-modify-write sequence).
	//  (c) Use FailedPrecondition if the client should not retry until
	//      the system state has been explicitly fixed. E.g., if an "rmdir"
	//      fails because the directory is non-empty, FailedPrecondition
	//      should be returned since the client should not retry unless
	//      they have first fixed up the directory by deleting files from it.
	//  (d) Use FailedPrecondition if the client performs conditional
	//      REST Get/Update/Delete on a resource and the resource on the
	//      server does not match the condition. E.g., conflicting
	//      read-modify-write on the same resource.
	FailedPrecondition Code = -9

	// Aborted indicates the operation was aborted, typically due to a
	// concurrency issue like sequencer check failures, transaction aborts,
	// etc.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	Aborted Code = -10

	// OutOfRange means operation was attempted past the valid range.
	// E.g., seeking or reading past end of file.
	//
	// Unlike InvalidArgument, this error indicates a problem that may
	// be fixed if the system state changes. For example, a 32-bit file
	// system will generate InvalidArgument if asked to read at an
	// offset that is not in the range [0,2^32-1], but it will generate
	// OutOfRange if asked to read from an offset past the current
	// file size.
	//
	// There is a fair bit of overlap between FailedPrecondition and
	// OutOfRange. We recommend using OutOfRange (the more specific
	// error) when it applies so that callers who are iterating through
	// a space can easily look for an OutOfRange error to detect when
	// they are done.
	OutOfRange Code = -11

	// Unimplemented indicates operation is not implemented or not
	// supported/enabled in this service.
	Unimplemented Code = -12

	// Internal errors. Means some invariants expected by underlying
	// system has been broken. If you see one of these errors,
	// something is very broken.
	Internal Code = -13

	// Unavailable indicates the service is currently unavailable.
	// This is a most likely a transient condition and may be corrected
	// by retrying with a backoff. Note that it is not always safe to retry
	// non-idempotent operations.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	Unavailable Code = -14

	// DataLoss indicates unrecoverable data loss or corruption.
	DataLoss Code = -15

	// Unauthenticated indicates the request does not have valid
	// authentication credentials for the operation.
	//
	// = HTTP 401 Unauthorized
	Unauthenticated Code = -16

	// RateLimited indicates some flow control algorithm is running and applied.
	// = HTTP Code 429
	RateLimited = -17

	// BadRequest generates a 400 error.
	// = HTTP 400
	BadRequest = -18

	// Conflict generates a 409 error.
	// = hTTP 409
	Conflict = -19

	// Forbidden generates a 403 error.
	Forbidden = -20

	// InternalServerError generates a 500 error.
	InternalServerError = -21

	// MethodNotAllowed generates a 405 error.
	MethodNotAllowed = -22

	// MinErrorCode is the lower bound
	MinErrorCode = -1000
)

var strToCode = map[string]Code{
	`OK`:                  OK,
	`CANCELLED`:           Canceled,
	`UNKNOWN`:             Unknown,
	`INVALID_ARGUMENT`:    InvalidArgument,
	`DEADLINE_EXCEEDED`:   DeadlineExceeded,
	`NOT_FOUND`:           NotFound,
	`ALREADY_EXISTS`:      AlreadyExists,
	`PERMISSION_DENIED`:   PermissionDenied,
	`RESOURCE_EXHAUSTED`:  ResourceExhausted,
	`FAILED_PRECONDITION`: FailedPrecondition,
	`ABORTED`:             Aborted,
	`OUT_OF_RANGE`:        OutOfRange,
	`UNIMPLEMENTED`:       Unimplemented,
	`INTERNAL`:            Internal,
	`UNAVAILABLE`:         Unavailable,
	`DATA_LOSS`:           DataLoss,
	`UNAUTHENTICATED`:     Unauthenticated,
	`RATE_LIMITED`:        RateLimited,
}

var codeToStr = map[Code]string{
	OK:                 `OK`,
	Canceled:           `CANCELLED`,
	Unknown:            `UNKNOWN`,
	InvalidArgument:    `INVALID_ARGUMENT`,
	DeadlineExceeded:   `DEADLINE_EXCEEDED`,
	NotFound:           `NOT_FOUND`,
	AlreadyExists:      `ALREADY_EXISTS`,
	PermissionDenied:   `PERMISSION_DENIED`,
	ResourceExhausted:  `RESOURCE_EXHAUSTED`,
	FailedPrecondition: `FAILED_PRECONDITION`,
	Aborted:            `ABORTED`,
	OutOfRange:         `OUT_OF_RANGE`,
	Unimplemented:      `UNIMPLEMENTED`,
	Internal:           `INTERNAL`,
	Unavailable:        `UNAVAILABLE`,
	DataLoss:           `DATA_LOSS`,
	Unauthenticated:    `UNAUTHENTICATED`,
	RateLimited:        `RATE_LIMITED`,
}

// WithCodeInfo is a type integrating both error code, cause, message, and template
type WithCodeInfo struct {
	code      Code
	causer    error
	msg       string
	livedArgs []interface{}
}

// Code returns the error code value
func (w *WithCodeInfo) Code() Code {
	return w.code
}

// Equal tests if equals with code 'c'
func (w *WithCodeInfo) Equal(c Code) bool {
	return w.code == c
}

func (w *WithCodeInfo) Error() string {
	var buf bytes.Buffer
	buf.WriteString(w.code.String())
	if len(w.msg) > 0 {
		buf.WriteRune('|')
		if len(w.livedArgs) > 0 {
			buf.WriteString(fmt.Sprintf(w.msg, w.livedArgs))
		} else {
			buf.WriteString(w.msg)
		}
	}
	if w.causer != nil {
		buf.WriteRune('|')
		buf.WriteString(w.causer.Error())
	}
	return buf.String()
}

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//    %s	lists source files for each Frame in the stack
//    %v	lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+v   Prints filename, function, and line number for each Frame in the stack.
func (w *WithCodeInfo) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			msg := w.msg
			if len(w.livedArgs) > 0 {
				msg = fmt.Sprintf(w.msg, w.livedArgs...)
			}
			_, _ = fmt.Fprintf(s, "%d|%+v|%s", int(w.code), w.code.String(), msg)
			if w.causer != nil {
				_, _ = fmt.Fprintf(s, "|%+v", w.causer)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
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
func (w *WithCodeInfo) FormatNew(livedArgs ...interface{}) *WithStackInfo {
	x := WithCode(w.code, w.causer, w.msg)
	x.error.(*WithCodeInfo).livedArgs = livedArgs
	return x
}

// Attach appends errs
func (w *WithCodeInfo) Attach(errs ...error) {
	for _, err := range errs {
		if err != nil {
			w.causer = err
		}
	}
	if len(errs) > 1 {
		panic("*WithCodeInfo.Attach() can only wrap one child error object.")
	}
}

// Cause returns the underlying cause of the error recursively,
// if possible.
func (w *WithCodeInfo) Cause() error {
	return w.causer
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func (w *WithCodeInfo) Unwrap() error {
	return w.causer
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
func (w *WithCodeInfo) As(target interface{}) bool {
	return As(w.causer, target)
	//if target == nil {
	//	panic("errors: target cannot be nil")
	//}
	//val := reflect.ValueOf(target)
	//typ := val.Type()
	//if typ.Kind() != reflect.Ptr || val.IsNil() {
	//	panic("errors: target must be a non-nil pointer")
	//}
	//if e := typ.Elem(); e.Kind() != reflect.Interface && !e.Implements(errorType) {
	//	panic("errors: *target must be interface or implement error")
	//}
	//targetType := typ.Elem()
	//err := w.causer
	//for err != nil {
	//	if reflect.TypeOf(err).AssignableTo(targetType) {
	//		val.Elem().Set(reflect.ValueOf(err))
	//		return true
	//	}
	//	if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) {
	//		return true
	//	}
	//	err = Unwrap(err)
	//}
	//return false
}

// Is reports whether any error in err's chain matches target.
func (w *WithCodeInfo) Is(target error) bool {
	return w.causer == target || Is(w.causer, target)
	//if target == nil {
	//	return w.causer == target
	//}
	//
	//isComparable := reflect.TypeOf(target).Comparable()
	//for {
	//	if isComparable && w.causer == target {
	//		return true
	//	}
	//	if x, ok := w.causer.(interface{ Is(error) bool }); ok && x.Is(target) {
	//		return true
	//	}
	//	// TO/DO: consider supporting target.Is(err). This would allow
	//	// user-definable predicates, but also may allow for coping with sloppy
	//	// APIs, thereby making it easier to get away with them.
	//	//if err := Unwrap(w.causer); err == nil {
	//	//	return false
	//	//}
	//
	//	return w.causer == target
	//}
}

//
// ----------------------------
//

// WithCode formats a wrapped error object with error code.
func WithCode(code Code, err error, message string, args ...interface{}) *WithStackInfo {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	err = &WithCodeInfo{
		code:   code,
		causer: err,
		msg:    message,
	}
	return &WithStackInfo{
		error: err,
		Stack: callers(),
	}
}

// Equal compares error object and error Code
func Equal(err error, code Code) bool {
	x, ok := err.(interface{ Equal(Code) bool })
	if !ok {
		return false
	}
	return x.Equal(code)
}

// EqualR compares error object and error Code recursively
func EqualR(err error, code Code) bool {
retryEqual:
	x, ok := err.(interface{ Equal(Code) bool })
	if !ok {
		err = Cause1(err) // unwrap one level of inner wrapped error
		if err != nil {
			goto retryEqual
		}
		return false
	}
	return x.Equal(code)
}

// New create a new *CodedErr object based an error code
func (c Code) New(msg string, args ...interface{}) *WithStackInfo {
	return WithCode(c, nil, msg, args...)
}

// NewTemplate create an error template so that you may `FormatNew(liveArgs...)` late.
func (c Code) NewTemplate(tmpl string) *WithCodeInfo {
	err := &WithCodeInfo{
		code:   c,
		causer: nil,
		msg:    tmpl,
	}
	return err
}

// Error for error interface
func (c Code) Error() string { return c.String() }

// String for stringer interface
func (c Code) String() string {
	if x, ok := codeToStr[c]; ok {
		return x
	}
	return codeToStr[Unknown]
}

// Register register a code and its token string for using later
func (c Code) Register(codeName string) (errno Code) {
	errno = AlreadyExists
	if c <= MinErrorCode || c > 0 {
		if _, ok := strToCode[codeName]; !ok {
			if _, ok = codeToStr[c]; !ok {
				strToCode[codeName] = c
				codeToStr[c] = codeName
				errno = OK
			}
		}
	}
	return
}
