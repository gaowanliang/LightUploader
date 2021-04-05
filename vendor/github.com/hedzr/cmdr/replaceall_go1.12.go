// Copyright © 2020 Hedzr Yeh.

// +build go1.12

package cmdr

import (
	"reflect"
	"strings"
)

func replaceAll(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}

func isNil(to reflect.Value) bool {
	switch to.Kind() {
	case reflect.Uintptr:
		return to.UnsafeAddr() == 0
	case reflect.Array:
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr:
		return to.IsNil()
	case reflect.Interface:
		return to.IsNil()
	case reflect.Slice:
	case reflect.String:
	case reflect.Struct:
		return false
	case reflect.UnsafePointer:
		return to.IsNil()
	}
	return false
}
