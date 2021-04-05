# errors.v2

[![Build Status](https://travis-ci.org/hedzr/errors.svg?branch=master)](https://travis-ci.org/hedzr/errors)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/errors.svg?label=release)](https://github.com/hedzr/errors/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/hedzr/errors) 
[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/errors)](https://goreportcard.com/report/github.com/hedzr/errors)
[![codecov](https://codecov.io/gh/hedzr/errors/branch/v2/graph/badge.svg)](https://codecov.io/gh/hedzr/errors)

Wrapped errors and more for golang developing (not just for go1.13+).

`hedzr/errors` provides the compatbilities to your old project up to go 1.13.

`hedzr/errors` provides some extra enhancements for better context environment saving on error occurred.



## Import

```go
// wrong: import "github.com/hedzr/errors/v2"
import "gopkg.in/hedzr/errors.v2"
```



## Features




#### stdlib `errors' compatibilities

- `func As(err error, target interface{}) bool`
- `func Is(err, target error) bool`
- `func New(text string) error`
- `func Unwrap(err error) error`

#### `pkg/errors` compatibilities

- `func Wrap(err error, message string) error`
- `func Cause(err error) error`: unwraps recursively, just like Unwrap()
- [x] `func Cause1(err error) error`: unwraps just one level
- `func WithCause(cause error, message string, args ...interface{}) error`, = `Wrap`
- supports Stacktrace
  - in an error by `Wrap()`, stacktrace wrapped;
  - for your error, attached by `WithStack(cause error)`;

#### Enhancements

- `New(msg, args...)` combines New and `Newf`(if there is a name), WithMessage, WithMessagef, ...
- `WithCause(cause error, message string, args...interface{})`
- `Wrap(err error, message string, args ...interface{}) error`, no Wrapf
- `DumpStacksAsString(allRoutines bool)`: returns stack tracing information like debug.PrintStack()
- `CanXXX`:
   - `CanAttach(err interface{}) bool`
   - `CanCause(err interface{}) bool`
   - `CanUnwrap(err interface{}) bool`
   - `CanIs(err interface{}) bool`
   - `CanAs(err interface{}) bool`

#### Extras

- Container/Holder for a group of sub-errors
- Coded error: the predefined errno



## error Container and sub-errors (wrapped, attached or nested)

- `NewContainer(message string, args ...interface{}) *withCauses`
- `ContainerIsEmpty(container error) bool`
- `AttachTo(container *withCauses, errs ...error)`
- `withCauses.Attach(errs ...error)`

For example:

```go
func a() (err error){
	container = errors.NewContainer("sample error")
    // ...
    for {
        // ...
        // in a long loop, we can add many sub-errors into container 'c'...
        errors.AttachTo(container, io.EOF, io.ErrUnexpectedEOF, io.ErrShortBuffer, io.ErrShortWrite)
        // Or:
        // container.Attach(someFuncReturnsErr(xxx))
        // ... break
    }
	// and we extract all of them as a single parent error object now.
	err = container.Error()
	return
}

func b(){
    err := a()
    // test the containered error 'err' if it hosted a sub-error `io.ErrShortWrite` or not.
    if errors.Is(err, io.ErrShortWrite) {
        panic(err)
    }
}
```



## Coded error

- `Code` is a generic type of error codes
- `WithCode(code, err, msg, args...)` can format an error object with error code, attached inner err, message or msg template, and stack info.
- `Code.New(msg, args...)` is like `WithCode`.
- `Code.Register(codeNameString)` declares the name string of an error code yourself.
- `Code.NewTemplate(tmpl)` create an coded error template object `*WithCodeInfo`.
- `WithCodeInfo.FormateNew(livedArgs...)` formats the err msg till used.
- `Equal(err, code)`: compares `err` with `code`

Try it at: <https://play.golang.org/p/Y2uThZHAvK1>

### Builtin Codes

The builtin Codes are errors, such as `OK`, `Canceled`, `Unknown`, `InvalidArgument`, `DeadlineExceeded`, `NotFound`, `AlreadyExists`,  etc..

```go
// Uses a Code as an error
var err error = errors.OK
var err2 error = errors.InvalidArgument
fmt.Println("error is: %v", err2)

// Uses a Code as enh-error (hedzr/errors)
err := InvalidArgument.New("wrong").Attach(io.ErrShortWrite)
```

### Customized Codes

```go
// customizing the error code
const MyCode001 errors.Code=1001

// and register the name of MyCode001
MyCode001.Register("MyCode001")

// and use it as a builtin Code
fmt.Println("error is: %v", MyCode001)
err := MyCode001.New("wrong 001: no config file")
```

### Error Template: formatting the coded-error late


```go
const BUG1001 errors.Code=1001
errTmpl1001 := BUG1001.NewTemplate("something is wrong, %v")
err4 := errTmpl1001.FormatNew("unsatisfied conditions").Attach(io.ShortBuffer)
fmt.Println(err4)
fmt.Printf("%+v\n", err4)
```




## ACK

- stack.go is an copy from pkg/errors
- withStack is an copy from pkg/errors
- Is, As, Unwrap are inspired from go1.13 errors
- Cause, Wrap are inspired from pkg/errors

## LICENSE

MIT
