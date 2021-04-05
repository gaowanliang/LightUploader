package errors

import "fmt"

// NewContainer wraps a group of errors and msg as one and return it.
// The returned error object is a container to hold many sub-errors.
//
// For Example:
//
// 		c := errors.NewContainer("sample error")
// 		... for a long loop
// 		errors.AttachTo(c, io.EOF, io.ErrUnexpectedEOF, io.ErrShortBuffer, io.ErrShortWrite)
// 		...
// 		err = c.Error()
//
//
func NewContainer(message string, args ...interface{}) *WithCauses {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	err := &WithCauses{
		msg:   message,
		Stack: callers(),
	}
	return err
}

// ContainerIsEmpty tests if 'container' is empty (no more wrapped/attached sub-errors)
func ContainerIsEmpty(container error) bool {
	if x, ok := container.(interface{ IsEmpty() bool }); ok {
		return x.IsEmpty()
	}
	return false
}

// AttachTo appends more errors into 'container' error container.
func AttachTo(container *WithCauses, errs ...error) {
	if container == nil {
		panic("nil error container not allowed")
	}
	container.Attach(errs...)
}
