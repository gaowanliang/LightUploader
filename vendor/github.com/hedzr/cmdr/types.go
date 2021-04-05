/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import "reflect"

// func isTypeInt(kind reflect.Kind) bool {
// 	switch kind {
// 	case reflect.Int:
// 	case reflect.Int8:
// 	case reflect.Int16:
// 	case reflect.Int32:
// 	case reflect.Int64:
// 	case reflect.Uint:
// 	case reflect.Uint8:
// 	case reflect.Uint16:
// 	case reflect.Uint32:
// 	case reflect.Uint64:
// 	default:
// 		return false
// 	}
// 	return true
// }

func isTypeUint(kind reflect.Kind) bool {
	switch kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func isTypeSInt(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

func isTypeFloat(kind reflect.Kind) bool {
	switch kind {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func isTypeComplex(kind reflect.Kind) bool {
	switch kind {
	case reflect.Complex64, reflect.Complex128:
		return true
	default:
		return false
	}
}

func isBool(v interface{}) bool {
	_, ok := v.(bool)
	return ok
}

func isNil1(v interface{}) bool {
	return v == nil
}
