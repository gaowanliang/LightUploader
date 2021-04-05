// +build nacl

// Copyright © 2020 Hedzr Yeh.

package tool

// ReadPassword reads the password from stdin with safe protection
func ReadPassword() (text string, err error) {
	return randomStringPure(9), nil
}

// GetTtySize returns the window size in columns and rows in the active console window.
// The return value of this function is in the order of cols, rows.
func GetTtySize() (cols, rows int) {
	var sz struct {
		rows, cols, xPixels, yPixels uint16
	}
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&sz)))
	cols, rows = int(sz.cols), int(sz.rows)
	return
}
