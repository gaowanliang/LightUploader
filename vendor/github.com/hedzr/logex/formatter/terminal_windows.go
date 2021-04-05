// +build !appengine,!js,windows

/*
 * Copyright © 2019 Hedzr Yeh.
 */

package formatter

import (
	sequences "github.com/konsorten/go-windows-terminal-sequences"
	"io"
	"os"
	"syscall"
)

func initTerminal(w io.Writer) {
	switch v := w.(type) {
	case *os.File:
		sequences.EnableVirtualTerminalProcessing(syscall.Handle(v.Fd()), true)
	}
}
