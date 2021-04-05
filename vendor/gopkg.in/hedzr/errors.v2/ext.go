// Copyright © 2020 Hedzr Yeh.

package errors

import (
	"runtime"
)

// DumpStacksAsString returns stack tracing information like debug.PrintStack()
func DumpStacksAsString(allRoutines bool) string {
	buf := make([]byte, 16384)
	buf = buf[:runtime.Stack(buf, allRoutines)]
	// fmt.Printf("=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===\n", buf)
	return string(buf)
}

// CanAttach tests if err is attach-able
func CanAttach(err interface{}) (ok bool) {
	_, ok = err.(interface{ Attach(errs ...error) })
	return
}

// CanCause tests if err is cause-able
func CanCause(err interface{}) (ok bool) {
	_, ok = err.(causer)
	return
}

// // CanWalk tests if err is walkable
// func CanWalk(err error) (ok bool) {
// 	_, ok = err.(Walkable)
// 	return
// }
//
// // CanRange tests if err is range-able
// func CanRange(err error) (ok bool) {
// 	_, ok = err.(Ranged)
// 	return
// }

// CanUnwrap tests if err is unwrap-able
func CanUnwrap(err interface{}) (ok bool) {
	_, ok = err.(interface{ Unwrap() error })
	return
}

// CanIs tests if err is is-able
func CanIs(err interface{}) (ok bool) {
	_, ok = err.(interface{ Is(error) bool })
	return
}

// CanAs tests if err is as-able
func CanAs(err interface{}) (ok bool) {
	_, ok = err.(interface{ As(interface{}) bool })
	return
}
