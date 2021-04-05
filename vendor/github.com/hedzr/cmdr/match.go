// Copyright © 2020 Hedzr Yeh.

package cmdr

import (
	"os"
	"strings"
)

// Match try parsing the input command-line, the result is the last hit *Command.
func Match(inputCommandlineWithoutArg0 string, opts ...ExecOption) (last *Command, err error) {
	saved := internalGetWorker()
	savedUnknownOptionHandler := unknownOptionHandler
	defer func() {
		uniqueWorkerLock.Lock()
		uniqueWorker = saved
		unknownOptionHandler = savedUnknownOptionHandler
		uniqueWorkerLock.Unlock()
	}()

	rootCmd := internalGetWorker().rootCommand

	w := internalResetWorkerNoLock()

	for _, opt := range opts {
		opt(w)
	}

	w.noDefaultHelpScreen = true
	w.noUnknownCmdTip = true
	w.noCommandAction = true
	unknownOptionHandler = emptyUnknownOptionHandler

	line := os.Args[0] + " " + inputCommandlineWithoutArg0
	last, err = w.InternalExecFor(rootCmd, strings.Split(line, " "))
	return
}
