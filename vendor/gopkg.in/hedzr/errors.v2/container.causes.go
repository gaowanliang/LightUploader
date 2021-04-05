// Copyright © 2020 Hedzr Yeh.

package errors

import (
	"bytes"
	"fmt"
	"io"
)

type causes struct {
	Causers []error
	*Stack
}

func (w *causes) Error() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, c := range w.Causers {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(c.Error())
	}
	buf.WriteString("]")
	// buf.WriteString(w.Stack)
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
func (w *causes) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", w.Error())
			w.Stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}

func (w *causes) Cause() error {
	if len(w.Causers) == 0 {
		return nil
	}
	return w.Causers[0]
}

func (w *causes) Causes() []error {
	if len(w.Causers) == 0 {
		return nil
	}
	return w.Causers
}

func (w *causes) Unwrap() error {
	// return w.Cause()

	for _, err := range w.Causers {
		//u, ok := err.(interface {
		//	Unwrap() error
		//})
		//if ok {
		//	return u.Unwrap()
		//}
		return err // just return the first cause
	}
	return nil
}

func (w *causes) Is(target error) bool {
	return IsSlice(w.Causers, target)
	//if target == nil {
	//	//for _, e := range w.Causers {
	//	//	if e == target {
	//	//		return true
	//	//	}
	//	//}
	//	return false
	//}
	//
	//isComparable := reflect.TypeOf(target).Comparable()
	//for {
	//	if isComparable {
	//		for _, e := range w.Causers {
	//			if e == target {
	//				return true
	//			}
	//		}
	//		// return false
	//	}
	//
	//	for _, e := range w.Causers {
	//		if x, ok := e.(interface{ Is(error) bool }); ok && x.Is(target) {
	//			return true
	//		}
	//		//if err := Unwrap(e); err == nil {
	//		//	return false
	//		//}
	//	}
	//	return false
	//}
}

// As finds the first error in `err`'s chain that matches target, and if so, sets
// target to that error value and returns true.
func (w *causes) As(target interface{}) bool {
	return AsSlice(w.Causers, target)
}
