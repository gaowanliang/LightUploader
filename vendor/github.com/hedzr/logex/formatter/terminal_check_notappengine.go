// +build !appengine,!js,!windows,!aix,!nacl

/*
 * Copyright © 2019 Hedzr Yeh.
 */

package formatter

import (
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// checkIfTerminal return bool
func checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}
