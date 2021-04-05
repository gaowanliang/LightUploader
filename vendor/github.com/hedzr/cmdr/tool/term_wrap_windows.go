// +build windows
// +build !nacl

// Copyright © 2020 Hedzr Yeh.

package tool

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
)

// ReadPassword reads the password from stdin with safe protection
func ReadPassword() (text string, err error) {
	var bytePassword []byte
	if bytePassword, err = terminal.ReadPassword(0); err == nil {
		fmt.Println() // it's necessary to add a new line after user's input
		text = string(bytePassword)
	} else {
		fmt.Println() // it's necessary to add a new line after user's input
	}
	return
}

// GetTtySize returns the window size in columns and rows in the active console window.
// The return value of this function is in the order of cols, rows.
func GetTtySize() (cols, rows int) {
	// return 0, 0
	cols, rows, _ = terminal.GetSize(0) // https://stackoverflow.com/a/45422726/6375060
	return
}
