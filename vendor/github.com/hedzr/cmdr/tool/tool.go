// Copyright © 2020 Hedzr Yeh.

package tool

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"gopkg.in/hedzr/errors.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// ParseComplex converts a string to complex number.
//
// Examples:
//
//    c1 := cmdr.ParseComplex("3-4i")
//    c2 := cmdr.ParseComplex("3.13+4.79i")
func ParseComplex(s string) (v complex128) {
	return a2complexShort(s)
}

// ParseComplexX converts a string to complex number.
// If the string is not valid complex format, return err not nil.
//
// Examples:
//
//    c1 := cmdr.ParseComplex("3-4i")
//    c2 := cmdr.ParseComplex("3.13+4.79i")
func ParseComplexX(s string) (v complex128, err error) {
	return a2complex(s)
}

func a2complexShort(s string) (v complex128) {
	v, _ = a2complex(s)
	return
}

func a2complex(s string) (v complex128, err error) {
	s = strings.TrimSpace(strings.TrimRightFunc(strings.TrimLeftFunc(s, func(r rune) bool {
		return r == '('
	}), func(r rune) bool {
		return r == ')'
	}))

	if i := strings.IndexAny(s, "+-"); i >= 0 {
		rr, ii := s[0:i], s[i:]
		if j := strings.Index(ii, "i"); j >= 0 {
			var ff, fi float64
			ff, err = strconv.ParseFloat(strings.TrimSpace(rr), 64)
			if err != nil {
				return
			}
			fi, err = strconv.ParseFloat(strings.TrimSpace(ii[0:j]), 64)
			if err != nil {
				return
			}

			v = complex(ff, fi)
			return
		}
		err = errors.New("for a complex number, the imaginary part should end with 'i', such as '3+4i'")
		return

		// err = errors.New("not valid complex number.")
	}

	var ff float64
	ff, err = strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return
	}
	v = complex(ff, 0)
	return
}

//
// external
//

func shellEditorRandomFilename() (fn string) {
	buf := make([]byte, 16)
	fn = os.Getenv("HOME") + ".CMDR_EDIT_FILE"
	if _, err := rand.Read(buf); err == nil {
		fn = fmt.Sprintf("%v/.CMDR_%x", os.Getenv("HOME"), buf)
	}
	return
}

// Launch executes a command setting both standard input, output and error.
func Launch(cmd string, args ...string) (err error) {
	c := exec.Command(cmd, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	err = c.Run()

	if err != nil {
		if _, isExitError := err.(*exec.ExitError); isExitError {
			err = nil
		}
	}
	return
}

// // LaunchSudo executes a command under "sudo".
// func LaunchSudo(cmd string, args ...string) error {
// 	return Launch("sudo", append([]string{cmd}, args...)...)
// }

//
// editor
//

// func getEditor() (string, error) {
// 	if GetEditor != nil {
// 		return GetEditor()
// 	}
// 	return exec.LookPath(DefaultEditor)
// }

// LaunchEditor launches the specified editor
func LaunchEditor(editor string) (content []byte, err error) {
	return launchEditorWith(editor, shellEditorRandomFilename())
}

// LaunchEditorWith launches the specified editor with a filename
func LaunchEditorWith(editor string, filename string) (content []byte, err error) {
	return launchEditorWith(editor, filename)
}

func launchEditorWith(editor, filename string) (content []byte, err error) {
	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		if _, isExitError := err.(*exec.ExitError); !isExitError {
			return
		}
	}

	content, err = ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, nil
	}
	return
}

// Max return the larger one of int
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min return the less one of int
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// StripPrefix strips the prefix 'p' from a string 's'
func StripPrefix(s, p string) string {
	return stripPrefix(s, p)
}

func stripPrefix(s, p string) string {
	if strings.HasPrefix(s, p) {
		return s[len(p):]
	}
	return s
}

// StripOrderPrefix strips the prefix string fragment for sorting order.
// see also: Command.Group, Flag.Group, ...
// An order prefix is a dotted string with multiple alphabet and digit. Such as:
// "zzzz.", "0001.", "700.", "A1." ...
func StripOrderPrefix(s string) string {
	if xre.MatchString(s) {
		s = s[strings.Index(s, ".")+1:]
	}
	return s
}

// HasOrderPrefix tests whether an order prefix is present or not.
// An order prefix is a dotted string with multiple alphabet and digit. Such as:
// "zzzz.", "0001.", "700.", "A1." ...
func HasOrderPrefix(s string) bool {
	return xre.MatchString(s)
}

var (
	xre = regexp.MustCompile(`^[0-9A-Za-z]+\.(.+)$`)
)

// IsDigitHeavy tests if the whole string is digit
func IsDigitHeavy(s string) bool {
	m, _ := regexp.MatchString("^\\d+$", s)
	// if err != nil {
	// 	return false
	// }
	return m
}

// PressEnterToContinue lets program pause and wait for user's ENTER key press in console/terminal
func PressEnterToContinue(in io.Reader, msg ...string) (input string) {
	if len(msg) > 0 && len(msg[0]) > 0 {
		fmt.Print(msg[0])
	} else {
		fmt.Print("Press 'Enter' to continue...")
	}
	b, _ := bufio.NewReader(in).ReadBytes('\n')
	return strings.TrimRight(string(b), "\n")
}

// PressAnyKeyToContinue lets program pause and wait for user's ANY key press in console/terminal
func PressAnyKeyToContinue(in io.Reader, msg ...string) (input string) {
	if len(msg) > 0 && len(msg[0]) > 0 {
		fmt.Print(msg[0])
	} else {
		fmt.Print("Press any key to continue...")
	}
	_, _ = fmt.Fscanf(in, "%s", &input)
	return
}

// SavedOsArgs is a copy of os.Args, just for testing
var SavedOsArgs []string

func init() {
	if SavedOsArgs == nil {
		// bug: can't copt slice to slice: _ = StandardCopier.Copy(&SavedOsArgs, &os.Args)
		for _, s := range os.Args {
			SavedOsArgs = append(SavedOsArgs, s)
		}
	}
}
