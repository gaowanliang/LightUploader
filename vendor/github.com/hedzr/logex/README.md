# logex

![Go](https://github.com/hedzr/logex/workflows/Go/badge.svg)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/logex.svg?label=release)](https://github.com/hedzr/logex/releases)
[![Sourcegraph](https://sourcegraph.com/github.com/hedzr/logex/-/badge.svg)](https://sourcegraph.com/github.com/hedzr/logex?badge)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/hedzr/logex)
[![go.dev](https://img.shields.io/badge/go.dev-reference-green)](https://pkg.go.dev/github.com/hedzr/logex)
[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/logex)](https://goreportcard.com/report/github.com/hedzr/logex)
[![Coverage Status](https://coveralls.io/repos/github/hedzr/logex/badge.svg?branch=master)](https://coveralls.io/github/hedzr/logex?branch=master) <!--
 [![codecov](https://codecov.io/gh/hedzr/logex/branch/master/graph/badge.svg)](https://codecov.io/gh/hedzr/logex) --> 
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Flogex.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Flogex?ref=badge_shield)


an enhancement for [logrus](https://github.com/sirupsen/logrus). `logex` attaches the context caller info to the logging output.

Since v1.2.0, `logex` allows switching the logging backends (such as logrus, zap, ...) transparently.



![image-20200927083958978](https://i.loli.net/2020/09/27/LYlAcGUOa3CIeR7.png)



## Features

- Pre-setup logging backends with clickable caller info: logrus or zap
- Generic logging interface to cover the various logging backends via: [`log.Logger`](https://github.com/hedzr/log/blob/master/logger.go#L10), [`build.New(config)`](https://github.com/hedzr/logex/blob/master/build/builder.go#L14)
  - allow registering custom logging backend by [build.RegisterBuilder(backendName, log.BuilderFunc) ` ](https://github.com/hedzr/logex/blob/master/build/builder.go#L43).



## Usages



### Build logger transparently

We provide the ability to switch logging backends transparently now.

A sample config file looks like:

```yaml
app:

  # runmode: devel  # devel, prod

  logger:
    # The field 'level' will be reset to "debug" while the app is started up within a debugger
    # available values are:
    #   "disable"/"off", "panic", "fatal", "error", "warn", "info", "debug", "trace"
    level:  info
    format: text                  # text, json, logfmt, ...
    backend: sugar                # zap, sugar(sugared-zap) or logrus
    target: file                  # console, file
    directory: /var/log/$APPNAME
```

Load it to Config structure:

```go
import "github.com/hedzr/log"
var config *log.LoggerConfig = log.NewLoggerConfig()
// ...
```

And build the backend:

```go
import "github.com/hedzr/logex/build"
logger := build.New(config)
logger.Debugf("int value = %v", intVal)
```

#### Or build a logger backend directly

```go
import "github.com/hedzr/logex/logx/logrus"
logrus.New(level string, traceMode, debugMode bool, opts ...Opt)

import "github.com/hedzr/logex/logx/zap"
zap.New(level string, traceMode, debugMode bool, opts ...Opt)

import "github.com/hedzr/logex/logx/zap/sugar"
sugar.New(level string, traceMode, debugMode bool, opts ...Opt)

```

#### Or, build the logger with pure go codes

```go
import "github.com/hedzr/logex/build"
// config:=build.NewLoggerConfig()
config := build.NewLoggerConfigWith(true, "logrus", "debug")
logger := build.New(config)
logger.Debugf("int value = %v", intVal)
```



That's all stocked.



### Integrating your backend:

You can wrap a logging backend with `log.Logger` and register it into logex/build. Why we should do it like this? A universal logger creating interface from logex/build will simplify the application initiliazing coding, esp. in a framework.

```go
import "github.com/hedzr/logex/build"

build.RegisterBuilder("someone", createSomeLogger)

func createSomeLogger(config *log.LoggerConfig) log.Logger {
	//... wrapping your logging backend to log.Logger
}

// and use it:
build.New(build.NewLoggerConfigWith(true, "someone", "debug"))
build.New(build.NewLoggerConfigWith(false, "someone", "info"))
```



### Legacy tools

#### Enable logrus

```go
import "github.com/hedzr/logex"

func init(){
    logex.Enable()
    // Or:
    logex.EnableWith(logrus.DebugLevel)
}
```


#### Ignore the extra caller frames

If you are writing logging func wrappers, you might ignore the extra caller frames for those wrappers:

```go
func wrong(err error, fmt string, args interface{}) {
    logrus.WithError(err).WithFields(logrus.Fields{
        logex.SKIP: 1,  // ignore wrong() frame
    }).Errorf(fmt, args)
}

func wrongInner(err error, fields logrus.Fields, fmt string, args interface{}) {
    logrus.WithError(err).WithFields(fields).Errorf(fmt, args)
}

func wrongwrong(err error, fmt string, args interface{}) {
    wrongInner(err, logrus.Fields{
        logex.SKIP: 2,  // ignore wrongwrong() and wrongInner() frame
    }, fmt, args...)
}
```





## For `go test`

## make `logrus` works in `go test`

The codes is copied from:

<https://github.com/sirupsen/logrus/issues/834>

And in a test function, you could code now:

```go
   func TestFoo(t *testing.T) {
     defer logex.CaptureLog(t).Release()
     // …
   }
```





## LICENSE

MIT


[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Flogex.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Flogex?ref=badge_large)