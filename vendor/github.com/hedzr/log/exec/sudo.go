// Copyright © 2020 Hedzr Yeh.

package exec

import (
	"bytes"
	"fmt"
	"gopkg.in/hedzr/errors.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

// Run runs an OS command
func Run(command string, arguments ...string) error {
	_, _, err := RunCommand(command, false, arguments...)
	return err
}

// Sudo runs an OS command with sudo prefix
func Sudo(command string, arguments ...string) (int, string, error) {
	sudocmd, err := exec.LookPath("sudo")
	if err != nil {
		return -1, "'sudo' not found", Run(command, arguments...)
	}

	rc, output, err1 := RunCommand(sudocmd, true, append([]string{command}, arguments...)...)
	return rc, output, err1
}

// RunWithOutput runs an OS command and collect the result outputting
func RunWithOutput(command string, arguments ...string) (int, string, error) {
	return RunCommand(command, true, arguments...)
}

// RunCommand runs an OS command
func RunCommand(command string, readStdout bool, arguments ...string) (retCode int, stdoutText string, err error) {
	cmd := exec.Command(command, arguments...)

	var output string
	var stdout io.ReadCloser
	var stderr io.ReadCloser

	if readStdout {
		// Connect pipe to read Stdout
		stdout, err = cmd.StdoutPipe()
		if err != nil {
			// Failed to connect pipe
			return 0, "", fmt.Errorf("%q failed to connect stdout pipe: %v", command, err)
		}

		defer stdout.Close()
	} else {
		cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
	}

	// Connect pipe to read Stderr
	stderr, err = cmd.StderrPipe()
	if err != nil {
		// Failed to connect pipe
		return 0, "", fmt.Errorf("%q failed to connect stderr pipe: %v", command, err)
	}

	defer stderr.Close()

	// Do not use cmd.Run()
	if err = cmd.Start(); err != nil {
		// Problem while copying stdin, stdout, or stderr
		return 0, "", fmt.Errorf("%q failed: %v", command, err)
	}

	// Zero exit status
	// Darwin: launchctl can fail with a zero exit status,
	// so check for emtpy stderr
	if command == "launchctl" {
		slurp, _ := ioutil.ReadAll(stderr)
		if len(slurp) > 0 && !bytes.HasSuffix(slurp, []byte("Operation now in progress\n")) {
			return 0, "", fmt.Errorf("%q failed with stderr: %s", command, slurp)
		}
	}

	slurp, _ := ioutil.ReadAll(stderr)

	if err = cmd.Wait(); err != nil {
		exitStatus, ok := IsExitError(err)
		if ok {
			// Command didn't exit with a zero exit status.
			return exitStatus, output, errors.New("%q failed with stderr:\n%v\n  ", command, string(slurp)).Attach(err)
		}

		// An error occurred and there is no exit status.
		//return 0, output, fmt.Errorf("%q failed: %v |\n  stderr: %s", command, err.Error(), slurp)
		return 0, output, errors.New("%q failed with stderr:\n%v\n  ", command, string(slurp)).Attach(err)
	}

	if readStdout {
		var out []byte
		out, err = ioutil.ReadAll(stdout)
		if err != nil {
			return 0, "", fmt.Errorf("%q failed while attempting to read stdout: %v", command, err)
		} else if len(out) > 0 {
			output = string(out)
		}
	}

	return 0, output, nil
}

// IsExitError checks the error object
func IsExitError(err error) (int, bool) {
	if ee, ok := err.(*exec.ExitError); ok {
		if status, ok := ee.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus(), true
		}
	}

	return 0, false
}

// IsEAccess detects whether err is a EACCESS errno or not
func IsEAccess(err error) bool {
	if e, ok := err.(*os.PathError); ok && e.Err == syscall.EACCES {
		return true
	}
	return false
}
