/*
 * Copyright © 2019 Hedzr Yeh.
 */

// Package flag is used to wrap some APIs from go stdlib flag
package flag

import (
	"github.com/hedzr/cmdr"
	"reflect"
	"time"
)

type (
	// Option is used by cmdr fluent API and flag compatible API
	Option func(flag cmdr.OptFlag)
)

var (
	pfRootCmd *cmdr.RootCmdOpt
	parsed    bool

	treatAsLongOpt bool
)

//
// -----------------------
//

// WithTitles setup short title, long title, and aliases titles
func WithTitles(short, long string, aliases ...string) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.Titles(short, long, aliases...)
	}
}

// WithShort sets the short title
func WithShort(short string) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.Short(short)
	}
}

// WithLong sets the Long/Full title
func WithLong(long string) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.Long(long)
	}
}

// WithAliases sets the aliases string list
func WithAliases(aliases ...string) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.Aliases(aliases...)
	}
}

// WithDescription sets the description string
func WithDescription(oneLine, long string) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.Description(oneLine, long)
	}
}

// WithExamples sets the example string for an option
func WithExamples(examples string) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.Examples(examples)
	}
}

// WithGroup sets the group name
func WithGroup(group string) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.Group(group)
	}
}

// WithHidden sets an hidden option that does not be displayed in any list or help screen.
func WithHidden(hidden bool) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.Hidden(hidden)
	}
}

// WithDeprecated sets a version string for an deprecation option
func WithDeprecated(deprecation string) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.Deprecated(deprecation)
	}
}

// WithAction to specify the action as the option was matched
func WithAction(action func(cmd *cmdr.Command, args []string) (err error)) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.Action(action)
	}
}

// WithToggleGroup allows to specify an group name, and any options in this group will be treated as an toggleable group, just like raido button group.
func WithToggleGroup(group string) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.ToggleGroup(group)
	}
}

// WithDefaultValue set the value with explicit data type, and its placeholder name.
func WithDefaultValue(val interface{}, placeholder string) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.DefaultValue(val, placeholder)
	}
}

// WithExternalTool allows launch an external program via an environment key yours specified.
// for example, while you setup by `WithExternalTool("EDITOR")`, cmdr will lookup it from os environment and launch that program.
// for EDITOR=vim, `vim` will be launched.
func WithExternalTool(envKeyName string) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.ExternalTool(envKeyName)
	}
}

// WithValidArgs enables enumerable values for an option.
func WithValidArgs(list ...string) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.ValidArgs(list...)
	}
}

// WithHeadLike enables `head -n` mode.
// min, max will be ignored at this version, its might be impl in the future
func WithHeadLike(enable bool, min, max int64) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.HeadLike(enable, min, max)
	}
}

// WithEnvKeys binds the environ variable keynames to an option.
func WithEnvKeys(keys ...string) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.EnvKeys(keys...)
	}
}

// WithOnSet binds the OnSet handler to an option.
func WithOnSet(f func(keyPath string, value interface{})) (opt Option) {
	return func(flag cmdr.OptFlag) {
		flag.OnSet(f)
	}
}

//
// -----------------------
//

// WithCommand define an (Sub-)Command
func WithCommand(cmdDefines func(newSubCmd cmdr.OptCmd)) (opt Option) {
	return func(flag cmdr.OptFlag) {
		var oo = flag.OwnerCommand().NewSubCommand()
		cmdDefines(oo)
	}
}

//
// -----------------------
//

// BoolVar defines a bool flag with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the flag.
func BoolVar(p *bool, name string, value bool, usage string, options ...Option) {
	// CommandLine.Var(newBoolValue(value, p), name, usage)

	*p = value
	f := pfRootCmd.Bool()
	// f:= pfRootCmd.NewFlag(cmdr.OptFlagTypeString)
	f.Description(usage, usage).DefaultValue(value, "")
	if treatAsLongOpt {
		f.Long(name)
	} else {
		f.Short(name)
	}

	for _, opt := range options {
		opt(f)
	}

	f.OnSet(func(keyPath string, val interface{}) {
		if b, ok := val.(bool); ok {
			*p = b
		}
	})
}

// Bool defines a bool flag with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the flag.
func Bool(name string, value bool, usage string, options ...Option) *bool {
	var p = new(bool)
	BoolVar(p, name, value, usage, options...)
	return p
}

// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func IntVar(p *int, name string, value int, usage string, options ...Option) {
	// CommandLine.Var(newIntValue(value, p), name, usage)

	*p = value
	f := pfRootCmd.Int()
	// f:= pfRootCmd.NewFlag(cmdr.OptFlagTypeString)
	f.Description(usage, usage).DefaultValue(value, "")
	if treatAsLongOpt {
		f.Long(name)
	} else {
		f.Short(name)
	}

	for _, opt := range options {
		opt(f)
	}

	f.OnSet(func(keyPath string, val interface{}) {
		switch reflect.ValueOf(val).Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			*p = int(val.(uint64))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			*p = int(val.(int64))
		}
	})
}

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

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func Int(name string, value int, usage string, options ...Option) *int {
	// return CommandLine.Int(name, value, usage)
	var p = new(int)
	IntVar(p, name, value, usage, options...)
	return p
}

// Int64Var defines an int64 flag with specified name, default value, and usage string.
// The argument p points to an int64 variable in which to store the value of the flag.
func Int64Var(p *int64, name string, value int64, usage string, options ...Option) {
	// CommandLine.Var(newInt64Value(value, p), name, usage)

	*p = value
	f := pfRootCmd.Int64()
	// f:= pfRootCmd.NewFlag(cmdr.OptFlagTypeString)
	f.Description(usage, usage).DefaultValue(value, "")
	if treatAsLongOpt {
		f.Long(name)
	} else {
		f.Short(name)
	}

	for _, opt := range options {
		opt(f)
	}

	f.OnSet(func(keyPath string, val interface{}) {
		switch reflect.ValueOf(val).Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			*p = int64(val.(uint64))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			*p = int64(val.(int64))
		}
	})
}

// Int64 defines an int64 flag with specified name, default value, and usage string.
// The return value is the address of an int64 variable that stores the value of the flag.
func Int64(name string, value int64, usage string, options ...Option) *int64 {
	// return CommandLine.Int64(name, value, usage)
	var p = new(int64)
	Int64Var(p, name, value, usage, options...)
	return p
}

// UintVar defines a uint flag with specified name, default value, and usage string.
// The argument p points to a uint variable in which to store the value of the flag.
func UintVar(p *uint, name string, value uint, usage string, options ...Option) {
	// CommandLine.Var(newUintValue(value, p), name, usage)

	*p = value
	f := pfRootCmd.Uint()
	// f:= pfRootCmd.NewFlag(cmdr.OptFlagTypeString)
	f.Description(usage, usage).DefaultValue(value, "")
	if treatAsLongOpt {
		f.Long(name)
	} else {
		f.Short(name)
	}

	for _, opt := range options {
		opt(f)
	}

	f.OnSet(func(keyPath string, val interface{}) {
		switch reflect.ValueOf(val).Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			*p = uint(val.(uint64))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			*p = uint(val.(int64))
		}
	})
}

// Uint defines a uint flag with specified name, default value, and usage string.
// The return value is the address of a uint variable that stores the value of the flag.
func Uint(name string, value uint, usage string, options ...Option) *uint {
	// return CommandLine.Uint(name, value, usage)
	var p = new(uint)
	UintVar(p, name, value, usage, options...)
	return p
}

// Uint64Var defines a uint64 flag with specified name, default value, and usage string.
// The argument p points to a uint64 variable in which to store the value of the flag.
func Uint64Var(p *uint64, name string, value uint64, usage string, options ...Option) {
	// CommandLine.Var(newUint64Value(value, p), name, usage)

	*p = value
	f := pfRootCmd.Uint64()
	// f:= pfRootCmd.NewFlag(cmdr.OptFlagTypeString)
	f.Description(usage, usage).DefaultValue(value, "")
	if treatAsLongOpt {
		f.Long(name)
	} else {
		f.Short(name)
	}

	for _, opt := range options {
		opt(f)
	}

	f.OnSet(func(keyPath string, val interface{}) {
		switch reflect.ValueOf(val).Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			*p = uint64(val.(uint64))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			*p = uint64(val.(int64))
		}
	})
}

// Uint64 defines a uint64 flag with specified name, default value, and usage string.
// The return value is the address of a uint64 variable that stores the value of the flag.
func Uint64(name string, value uint64, usage string, options ...Option) *uint64 {
	// return CommandLine.Uint64(name, value, usage)
	var p = new(uint64)
	Uint64Var(p, name, value, usage, options...)
	return p
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func StringVar(p *string, name string, value string, usage string, options ...Option) {
	// CommandLine.Var(newStringValue(value, p), name, usage)

	*p = value
	f := pfRootCmd.String()
	// f:= pfRootCmd.NewFlag(cmdr.OptFlagTypeString)
	f.Description(usage, usage).DefaultValue(value, "")
	if treatAsLongOpt {
		f.Long(name)
	} else {
		f.Short(name)
	}

	for _, opt := range options {
		opt(f)
	}

	f.OnSet(func(keyPath string, val interface{}) {
		if b, ok := val.(string); ok {
			*p = b
		}
	})
}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func String(name string, value string, usage string, options ...Option) *string {
	// return CommandLine.String(name, value, usage)
	var p = new(string)
	StringVar(p, name, value, usage, options...)
	return p
}

// Float64Var defines a float64 flag with specified name, default value, and usage string.
// The argument p points to a float64 variable in which to store the value of the flag.
func Float64Var(p *float64, name string, value float64, usage string, options ...Option) {
	// CommandLine.Var(newFloat64Value(value, p), name, usage)

	*p = value
	f := pfRootCmd.Float64()
	// f:= pfRootCmd.NewFlag(cmdr.OptFlagTypeString)
	f.Description(usage, usage).DefaultValue(value, "")
	if treatAsLongOpt {
		f.Long(name)
	} else {
		f.Short(name)
	}

	for _, opt := range options {
		opt(f)
	}

	f.OnSet(func(keyPath string, val interface{}) {
		if b, ok := val.(float64); ok {
			*p = b
		}
	})
}

// Float64 defines a float64 flag with specified name, default value, and usage string.
// The return value is the address of a float64 variable that stores the value of the flag.
func Float64(name string, value float64, usage string, options ...Option) *float64 {
	// return CommandLine.Float64(name, value, usage)
	var p = new(float64)
	Float64Var(p, name, value, usage, options...)
	return p
}

// Float32Var defines a float32 flag with specified name, default value, and usage string.
// The argument p points to a float64 variable in which to store the value of the flag.
func Float32Var(p *float32, name string, value float32, usage string, options ...Option) {
	// CommandLine.Var(newFloat64Value(value, p), name, usage)

	*p = value
	f := pfRootCmd.Float32()
	// f:= pfRootCmd.NewFlag(cmdr.OptFlagTypeString)
	f.Description(usage, usage).DefaultValue(value, "")
	if treatAsLongOpt {
		f.Long(name)
	} else {
		f.Short(name)
	}

	for _, opt := range options {
		opt(f)
	}

	f.OnSet(func(keyPath string, val interface{}) {
		if b, ok := val.(float32); ok {
			*p = b
		}
	})
}

// Float32 defines a float64 flag with specified name, default value, and usage string.
// The return value is the address of a float64 variable that stores the value of the flag.
func Float32(name string, value float32, usage string, options ...Option) *float32 {
	// return CommandLine.Float64(name, value, usage)
	var p = new(float32)
	Float32Var(p, name, value, usage, options...)
	return p
}

// TODO: complex64, complex128, []uint, []int64, []uint64, ...

// DurationVar defines a time.Duration flag with specified name, default value, and usage string.
// The argument p points to a time.Duration variable in which to store the value of the flag.
// The flag accepts a value acceptable to time.ParseDuration.
func DurationVar(p *time.Duration, name string, value time.Duration, usage string, options ...Option) {
	// CommandLine.Var(newDurationValue(value, p), name, usage)

	*p = value
	f := pfRootCmd.Duration()
	// f:= pfRootCmd.NewFlag(cmdr.OptFlagTypeString)
	f.Description(usage, usage).DefaultValue(value, "")
	if treatAsLongOpt {
		f.Long(name)
	} else {
		f.Short(name)
	}

	for _, opt := range options {
		opt(f)
	}

	f.OnSet(func(keyPath string, val interface{}) {
		if b, ok := val.(time.Duration); ok {
			*p = b
		}
	})
}

// Duration defines a time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a time.Duration variable that stores the value of the flag.
// The flag accepts a value acceptable to time.ParseDuration.
func Duration(name string, value time.Duration, usage string, options ...Option) *time.Duration {
	// return CommandLine.Duration(name, value, usage)
	var p = new(time.Duration)
	DurationVar(p, name, value, usage, options...)
	return p
}

// StringSliceVar defines a string slice flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func StringSliceVar(p *[]string, name string, value []string, usage string, options ...Option) {
	// CommandLine.Var(newStringValue(value, p), name, usage)

	*p = value
	f := pfRootCmd.StringSlice()
	// f:= pfRootCmd.NewFlag(cmdr.OptFlagTypeString)
	f.Description(usage, usage).DefaultValue(value, "")
	if treatAsLongOpt {
		f.Long(name)
	} else {
		f.Short(name)
	}

	for _, opt := range options {
		opt(f)
	}

	f.OnSet(func(keyPath string, val interface{}) {
		if b, ok := val.([]string); ok {
			*p = b
		}
	})
}

// StringSlice defines a string slice flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func StringSlice(name string, value []string, usage string, options ...Option) *[]string {
	// return CommandLine.String(name, value, usage)
	var p = new([]string)
	StringSliceVar(p, name, value, usage, options...)
	return p
}

// IntSliceVar defines a int slice flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func IntSliceVar(p *[]int, name string, value []int, usage string, options ...Option) {
	// CommandLine.Var(newStringValue(value, p), name, usage)

	*p = value
	f := pfRootCmd.IntSlice()
	// f:= pfRootCmd.NewFlag(cmdr.OptFlagTypeString)
	f.Description(usage, usage).DefaultValue(value, "")
	if treatAsLongOpt {
		f.Long(name)
	} else {
		f.Short(name)
	}

	for _, opt := range options {
		opt(f)
	}

	f.OnSet(func(keyPath string, val interface{}) {
		if b, ok := val.([]int); ok {
			*p = b
		}
	})
}

// IntSlice defines a int slice flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func IntSlice(name string, value []int, usage string, options ...Option) *[]int {
	// return CommandLine.String(name, value, usage)
	var p = new([]int)
	IntSliceVar(p, name, value, usage, options...)
	return p
}

//
//
// ---------------------------------------------------------------------------
//
//
//
// // -- bool Value
// type boolValue bool
//
// func newBoolValue(val bool, p *bool) *boolValue {
// 	*p = val
// 	return (*boolValue)(p)
// }
//
// func (b *boolValue) Set(s string) error {
// 	v, err := strconv.ParseBool(s)
// 	if err != nil {
// 		err = errParse
// 	}
// 	*b = boolValue(v)
// 	return err
// }
//
// func (b *boolValue) Get() interface{} { return bool(*b) }
//
// func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }
//
// func (b *boolValue) IsBoolFlag() bool { return true }
//
// // optional interface to indicate boolean flags that can be
// // supplied without "=value" text
// type boolFlag interface {
// 	conf.Value
// 	IsBoolFlag() bool
// }
//
// // -- int Value
// type intValue int
//
// func newIntValue(val int, p *int) *intValue {
// 	*p = val
// 	return (*intValue)(p)
// }
//
// func (i *intValue) Set(s string) error {
// 	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
// 	if err != nil {
// 		err = numError(err)
// 	}
// 	*i = intValue(v)
// 	return err
// }
//
// func (i *intValue) Get() interface{} { return int(*i) }
//
// func (i *intValue) String() string { return strconv.Itoa(int(*i)) }
//
// // -- int64 Value
// type int64Value int64
//
// func newInt64Value(val int64, p *int64) *int64Value {
// 	*p = val
// 	return (*int64Value)(p)
// }
//
// func (i *int64Value) Set(s string) error {
// 	v, err := strconv.ParseInt(s, 0, 64)
// 	if err != nil {
// 		err = numError(err)
// 	}
// 	*i = int64Value(v)
// 	return err
// }
//
// func (i *int64Value) Get() interface{} { return int64(*i) }
//
// func (i *int64Value) String() string { return strconv.FormatInt(int64(*i), 10) }
//
// // -- uint Value
// type uintValue uint
//
// func newUintValue(val uint, p *uint) *uintValue {
// 	*p = val
// 	return (*uintValue)(p)
// }
//
// func (i *uintValue) Set(s string) error {
// 	v, err := strconv.ParseUint(s, 0, strconv.IntSize)
// 	if err != nil {
// 		err = numError(err)
// 	}
// 	*i = uintValue(v)
// 	return err
// }
//
// func (i *uintValue) Get() interface{} { return uint(*i) }
//
// func (i *uintValue) String() string { return strconv.FormatUint(uint64(*i), 10) }
//
// // -- uint64 Value
// type uint64Value uint64
//
// func newUint64Value(val uint64, p *uint64) *uint64Value {
// 	*p = val
// 	return (*uint64Value)(p)
// }
//
// func (i *uint64Value) Set(s string) error {
// 	v, err := strconv.ParseUint(s, 0, 64)
// 	if err != nil {
// 		err = numError(err)
// 	}
// 	*i = uint64Value(v)
// 	return err
// }
//
// func (i *uint64Value) Get() interface{} { return uint64(*i) }
//
// func (i *uint64Value) String() string { return strconv.FormatUint(uint64(*i), 10) }
//
// // -- string Value
// type stringValue string
//
// func newStringValue(val string, p *string) *stringValue {
// 	*p = val
// 	return (*stringValue)(p)
// }
//
// func (s *stringValue) Set(val string) error {
// 	*s = stringValue(val)
// 	return nil
// }
//
// func (s *stringValue) Get() interface{} { return string(*s) }
//
// func (s *stringValue) String() string { return string(*s) }
//
// // -- float64 Value
// type float64Value float64
//
// func newFloat64Value(val float64, p *float64) *float64Value {
// 	*p = val
// 	return (*float64Value)(p)
// }
//
// func (f *float64Value) Set(s string) error {
// 	v, err := strconv.ParseFloat(s, 64)
// 	if err != nil {
// 		err = numError(err)
// 	}
// 	*f = float64Value(v)
// 	return err
// }
//
// func (f *float64Value) Get() interface{} { return float64(*f) }
//
// func (f *float64Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }
//
// // -- time.Duration Value
// type durationValue time.Duration
//
// func newDurationValue(val time.Duration, p *time.Duration) *durationValue {
// 	*p = val
// 	return (*durationValue)(p)
// }
//
// func (d *durationValue) Set(s string) error {
// 	v, err := time.ParseDuration(s)
// 	if err != nil {
// 		err = errParse
// 	}
// 	*d = durationValue(v)
// 	return err
// }
//
// func (d *durationValue) Get() interface{} { return time.Duration(*d) }
//
// func (d *durationValue) String() string { return (*time.Duration)(d).String() }
//
// //
// // ------------------------
// //
//
// // errParse is returned by Set if a flag's value fails to parse, such as with an invalid integer for Int.
// // It then gets wrapped through failf to provide more information.
// var errParse = errors.New("parse error")
//
// // errRange is returned by Set if a flag's value is out of range.
// // It then gets wrapped through failf to provide more information.
// var errRange = errors.New("value out of range")
//
// func numError(err error) error {
// 	ne, ok := err.(*strconv.NumError)
// 	if !ok {
// 		return err
// 	}
// 	if ne.Err == strconv.ErrSyntax {
// 		return errParse
// 	}
// 	if ne.Err == strconv.ErrRange {
// 		return errRange
// 	}
// 	return err
// }
