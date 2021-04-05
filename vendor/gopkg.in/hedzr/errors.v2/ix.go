package errors

import "fmt"

// InterfaceCause is an interface with Cause
type InterfaceCause interface {
	// Cause returns the underlying cause of the error, if possible.
	// An error value has a cause if it implements the following
	// interface:
	//
	//     type causer interface {
	//            Cause() error
	//     }
	//
	// If the error does not implement Cause, the original error will
	// be returned. If the error is nil, nil will be returned without further
	// investigation.
	Cause() error
	// SetCause sets the underlying error manually if necessary.
	SetCause(cause error) error
}

// InterfaceCauses is an interface with Causes
type InterfaceCauses interface {
	Causes() []error
}

// InterfaceFormat is an interface with Format
type InterfaceFormat interface {
	// Format formats the stack of Frames according to the fmt.Formatter interface.
	//
	//    %s	lists source files for each Frame in the stack
	//    %v	lists the source file and line number for each Frame in the stack
	//
	// Format accepts flags that alter the printing of some verbs, as follows:
	//
	//    %+v   Prints filename, function, and line number for each Frame in the stack.
	Format(s fmt.State, verb rune)
}

// InterfaceIsAsUnwrap is an interface with Is, As, and Unwrap
type InterfaceIsAsUnwrap interface {
	// Is reports whether any error in err's chain matches target.
	Is(target error) bool
	// As finds the first error in err's chain that matches target, and if so, sets
	// target to that error value and returns true.
	As(target interface{}) bool
	// Unwrap returns the result of calling the Unwrap method on err, if err's
	// type contains an Unwrap method returning error.
	// Otherwise, Unwrap returns nil.
	Unwrap() error
}

// InterfaceAttachSpecial is an interface with Attach
type InterfaceAttachSpecial interface {
	// Attach appends errs
	Attach(errs ...error) *WithStackInfo
}

// InterfaceAttach is an interface with Attach
type InterfaceAttach interface {
	// Attach appends errs
	Attach(errs ...error)
}

// InterfaceContainer is an interface with IsEmpty
type InterfaceContainer interface {
	// IsEmpty tests has attached errors
	IsEmpty() bool
}

// InterfaceWithStackInfo is an interface for WithStackInfo
type InterfaceWithStackInfo interface {
	//error
	InterfaceCause
	InterfaceContainer
	InterfaceFormat
	InterfaceIsAsUnwrap
	InterfaceAttachSpecial
}

// Holder is an interface for WithCauses and InterfaceContainer
type Holder interface {
	Error() error
	//InterfaceCause
	InterfaceCauses
	InterfaceContainer
	InterfaceFormat
	InterfaceIsAsUnwrap
	InterfaceAttach
}
