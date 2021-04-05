// Copyright © 2020 Hedzr Yeh.

package logex

import "github.com/hedzr/log"

// InDebugging check if the delve debugger presents
func InDebugging() bool { return log.InDebugging() }

// GetDebugMode return the debug boolean flag generally
func GetDebugMode() bool { return log.GetDebugMode() }

// GetTraceMode return the trace boolean flag generally
func GetTraceMode() bool { return log.GetTraceMode() }

// SetDebugMode set the debug boolean flag generally
func SetDebugMode(b bool) { log.SetDebugMode(b) }

// SetTraceMode set the trace boolean flag generally
func SetTraceMode(b bool) { log.SetTraceMode(b) }
