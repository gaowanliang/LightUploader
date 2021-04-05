// +build nacl plan9

// Copyright © 2019 Hedzr Yeh.

package formatter

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return false
}
