// Copyright © 2020 Hedzr Yeh.

// +build go1.13

// go1.14

package cmdr

import "reflect"

// IsZero reports whether v is the zero value for its type.
// It panics if the argument is invalid.
func IsZero(v reflect.Value) bool {
	// switch v.Kind() {
	// case reflect.Bool:
	// 	break
	// }
	return v.IsZero()
}
