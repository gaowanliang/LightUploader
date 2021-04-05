// Copyright © 2020 Hedzr Yeh.

package log

import (
	"github.com/hedzr/log/isdelve"
	"github.com/hedzr/log/trace"
)

// CmdrMinimal provides the accessors to debug/trace flags
type CmdrMinimal interface {
	InDebugging() bool
	GetDebugMode() bool
	GetTraceMode() bool
	SetDebugMode(b bool)
	SetTraceMode(b bool)
}

// MinimalEnv return the Env/CmdrMinimal object
func MinimalEnv() CmdrMinimal { return env }

// InDebugging check if the delve debugger presents
func InDebugging() bool { return env.InDebugging() }

// GetDebugMode return the debug boolean flag generally
func GetDebugMode() bool { return env.GetDebugMode() }

// GetTraceMode return the trace boolean flag generally
func GetTraceMode() bool { return env.GetTraceMode() }

// SetDebugMode set the debug boolean flag generally
func SetDebugMode(b bool) { env.SetDebugMode(b) }

// SetTraceMode set the trace boolean flag generally
func SetTraceMode(b bool) { env.SetTraceMode(b) }

// Env structure holds the debug/trace flags and provides CmdrMinimal accessors
type Env struct {
	debugMode bool
	traceMode bool
}

// InDebugging check if the delve debugger presents
func (e *Env) InDebugging() bool { return isdelve.Enabled }

// GetDebugMode return the debug boolean flag generally
func (e *Env) GetDebugMode() bool { return e.debugMode || isdelve.Enabled }

// GetTraceMode return the trace boolean flag generally
func (e *Env) GetTraceMode() bool { return e.traceMode || trace.IsEnabled() }

// SetDebugMode set the debug boolean flag generally
func (e *Env) SetDebugMode(b bool) { e.debugMode = b }

// SetTraceMode set the trace boolean flag generally
func (e *Env) SetTraceMode(b bool) { e.traceMode = b }

var env = &Env{}
