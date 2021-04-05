/*
 * Copyright © 2020 Hedzr Yeh.
 */

package trace

import (
	"sync"
	"sync/atomic"
)

var tracing struct {
	sync.Mutex
	enabled int32
}

// Start enables the trace mode
func Start() (err error) {
	if atomic.CompareAndSwapInt32(&tracing.enabled, 0, 1) {
		tracing.Lock()
		defer tracing.Unlock()

		// trace.Start()
		// tracing.enabled

		for _, fn := range handlers {
			fn(false, true)
		}
	}

	return
}

// Stop disables the trace mode
func Stop() {
	if atomic.CompareAndSwapInt32(&tracing.enabled, 1, 0) {
		tracing.Lock()
		defer tracing.Unlock()

		// ...

		for _, fn := range handlers {
			fn(true, false)
		}
	}
}

// IsEnabled return the trace mode status
func IsEnabled() bool {
	enabled := atomic.LoadInt32(&tracing.enabled)
	return enabled == 1
}

// RegisterOnTraceModeChanges register a handler to capture Start/Stop activities
func RegisterOnTraceModeChanges(onTraceModeChanged Handler) {
	handlers = append(handlers, onTraceModeChanged)
}

// Handler can capture Start/Stop activities
type Handler func(old, new bool)

var handlers []Handler
