// Copyright © 2020 Hedzr Yeh.

package exec

import (
	"errors"
	"github.com/hedzr/log"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

// GetExecutableDir returns the executable file directory
func GetExecutableDir() string {
	// _ = ioutil.WriteFile("/tmp/11", []byte(strings.Join(os.Args,",")), 0644)
	// fmt.Printf("os.Args[0] = %v\n", os.Args[0])

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// fmt.Println(dir)
	return dir
}

// GetExecutablePath returns the executable file path
func GetExecutablePath() string {
	p, _ := filepath.Abs(os.Args[0])
	return p
}

// GetCurrentDir returns the current workingFlag directory
// it should be equal with os.Getenv("PWD")
func GetCurrentDir() string {
	dir, _ := os.Getwd()
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// fmt.Println(dir)
	return dir
}

// IsDirectory tests whether `path` is a directory or not
func IsDirectory(filepath string) (bool, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

// IsRegularFile tests whether `path` is a normal regular file or not
func IsRegularFile(filepath string) (bool, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return false, err
	}
	return fileInfo.Mode().IsRegular(), err
}

// FileModeIs tests the mode of 'filepath' with 'tester'. Examples:
//
//     var yes = exec.FileModeIs("/etc/passwd", exec.IsModeExecAny)
//     var yes = exec.FileModeIs("/etc/passwd", exec.IsModeDirectory)
//
func FileModeIs(filepath string, tester func(mode os.FileMode) bool) (ret bool) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return
	}
	ret = tester(fileInfo.Mode())
	return
}

// IsModeRegular give the result of whether a file is a regular file
func IsModeRegular(mode os.FileMode) bool { return mode.IsRegular() }

// IsModeDirectory give the result of whether a file is a directory
func IsModeDirectory(mode os.FileMode) bool { return mode&os.ModeDir != 0 }

// IsModeSymbolicLink give the result of whether a file is a symbolic link
func IsModeSymbolicLink(mode os.FileMode) bool { return mode&os.ModeSymlink != 0 }

// IsModeDevice give the result of whether a file is a device
func IsModeDevice(mode os.FileMode) bool { return mode&os.ModeDevice != 0 }

// IsModeNamedPipe give the result of whether a file is a named pipe
func IsModeNamedPipe(mode os.FileMode) bool { return mode&os.ModeNamedPipe != 0 }

// IsModeSocket give the result of whether a file is a socket file
func IsModeSocket(mode os.FileMode) bool { return mode&os.ModeSocket != 0 }

// IsModeSetuid give the result of whether a file has the setuid bit
func IsModeSetuid(mode os.FileMode) bool { return mode&os.ModeSetuid != 0 }

// IsModeSetgid give the result of whether a file has the setgid bit
func IsModeSetgid(mode os.FileMode) bool { return mode&os.ModeSetgid != 0 }

// IsModeCharDevice give the result of whether a file is a character device
func IsModeCharDevice(mode os.FileMode) bool { return mode&os.ModeCharDevice != 0 }

// IsModeSticky give the result of whether a file is a sticky file
func IsModeSticky(mode os.FileMode) bool { return mode&os.ModeSticky != 0 }

// IsModeIrregular give the result of whether a file is a non-regular file; nothing else is known about this file
func IsModeIrregular(mode os.FileMode) bool { return mode&os.ModeIrregular != 0 }

//

// IsModeExecOwner give the result of whether a file can be invoked by its unix-owner
func IsModeExecOwner(mode os.FileMode) bool { return mode&0100 != 0 }

// IsModeExecGroup give the result of whether a file can be invoked by its unix-group
func IsModeExecGroup(mode os.FileMode) bool { return mode&0010 != 0 }

// IsModeExecOther give the result of whether a file can be invoked by its unix-all
func IsModeExecOther(mode os.FileMode) bool { return mode&0001 != 0 }

// IsModeExecAny give the result of whether a file can be invoked by anyone
func IsModeExecAny(mode os.FileMode) bool { return mode&0111 != 0 }

// IsModeExecAll give the result of whether a file can be invoked by all users
func IsModeExecAll(mode os.FileMode) bool { return mode&0111 == 0111 }

//

// IsModeWriteOwner give the result of whether a file can be written by its unix-owner
func IsModeWriteOwner(mode os.FileMode) bool { return mode&0200 != 0 }

// IsModeWriteGroup give the result of whether a file can be written by its unix-group
func IsModeWriteGroup(mode os.FileMode) bool { return mode&0020 != 0 }

// IsModeWriteOther give the result of whether a file can be written by its unix-all
func IsModeWriteOther(mode os.FileMode) bool { return mode&0002 != 0 }

// IsModeWriteAny give the result of whether a file can be written by anyone
func IsModeWriteAny(mode os.FileMode) bool { return mode&0222 != 0 }

// IsModeWriteAll give the result of whether a file can be written by all users
func IsModeWriteAll(mode os.FileMode) bool { return mode&0222 == 0222 }

//

// IsModeReadOwner give the result of whether a file can be read by its unix-owner
func IsModeReadOwner(mode os.FileMode) bool { return mode&0400 != 0 }

// IsModeReadGroup give the result of whether a file can be read by its unix-group
func IsModeReadGroup(mode os.FileMode) bool { return mode&0040 != 0 }

// IsModeReadOther give the result of whether a file can be read by its unix-all
func IsModeReadOther(mode os.FileMode) bool { return mode&0004 != 0 }

// IsModeReadAny give the result of whether a file can be read by anyone
func IsModeReadAny(mode os.FileMode) bool { return mode&0444 != 0 }

// IsModeReadAll give the result of whether a file can be read by all users
func IsModeReadAll(mode os.FileMode) bool { return mode&0444 == 0444 }

//

// FileExists returns the existence of an directory or file
func FileExists(filepath string) bool {
	if _, err := os.Stat(os.ExpandEnv(filepath)); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// EnsureDir checks and creates the directory.
func EnsureDir(dir string) (err error) {
	if len(dir) == 0 {
		return errors.New("empty directory")
	}
	if !FileExists(dir) {
		err = os.MkdirAll(dir, 0755)
	}
	return
}

// EnsureDirEnh checks and creates the directory, via sudo if necessary.
func EnsureDirEnh(dir string) (err error) {
	if len(dir) == 0 {
		return errors.New("empty directory")
	}
	if !FileExists(dir) {
		err = os.MkdirAll(dir, 0755)
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.EACCES {
			var u *user.User
			u, err = user.Current()
			if _, _, err = Sudo("mkdir", "-p", dir); err == nil {
				_, _, err = Sudo("chown", u.Username+":", dir)
			}

			//if _, _, err = exec.Sudo("mkdir", "-p", dir); err != nil {
			//	logrus.Warnf("Failed to create directory %q, using default stderr. error is: %v", dir, err)
			//} else if _, _, err = exec.Sudo("chown", u.Username+":", dir); err != nil {
			//	logrus.Warnf("Failed to create directory %q, using default stderr. error is: %v", dir, err)
			//}
		}
	}
	return
}

// RemoveDirRecursive removes a directory and any children it contains.
func RemoveDirRecursive(dir string) (err error) {
	// RemoveContentsInDir(dir)
	err = os.RemoveAll(dir)
	return
}

// // RemoveContentsInDir removes all file and sub-directory in a directory
// func RemoveContentsInDir(dir string) error {
// 	d, err := os.Open(dir)
// 	if err != nil {
// 		return err
// 	}
// 	defer d.Close()
// 	names, err := d.Readdirnames(-1)
// 	if err != nil {
// 		return err
// 	}
// 	for _, name := range names {
// 		err = os.RemoveAll(filepath.Join(dir, name))
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// NormalizeDir make dir name normalized
func NormalizeDir(s string) string {
	return normalizeDir(s)
}

func normalizeDir(s string) string {
	p := normalizeDirBasic(s)
	p = filepath.Clean(p)
	return p
}

func normalizeDirBasic(s string) string {
	if len(s) == 0 {
		return s
	}

	s = os.Expand(s, os.Getenv)
	if s[0] == '/' {
		return s
	} else if strings.HasPrefix(s, "./") {
		return path.Join(GetCurrentDir(), s)
	} else if strings.HasPrefix(s, "../") {
		return path.Dir(path.Join(GetCurrentDir(), s))
	} else if strings.HasPrefix(s, "~/") {
		return path.Join(os.Getenv("HOME"), s[2:])
	} else {
		return s
	}
}

// AbsPath returns a clean, normalized and absolute path string for the given pathname.
func AbsPath(pathname string) string {
	return absPath(pathname)
}

func absPath(pathname string) (abs string) {
	abs = normalizePath(pathname)
	if s, err := filepath.Abs(abs); err == nil {
		abs = s
	}
	return
}

// NormalizePath cleans up the given pathname
func NormalizePath(pathname string) string {
	return normalizePath(pathname)
}

func normalizePath(pathname string) string {
	p := normalizePathBasic(pathname)
	p = filepath.Clean(p)
	return p
}

func normalizePathBasic(pathname string) string {
	if len(pathname) == 0 {
		return pathname
	}

	pathname = os.Expand(pathname, os.Getenv)
	if pathname[0] == '/' {
		return pathname
	} else if strings.HasPrefix(pathname, "~/") {
		return path.Join(os.Getenv("HOME"), pathname[2:])
	} else {
		return pathname
	}
}

// ForDir walks on `root` directory and its children
func ForDir(root string, cb func(depth int, cwd string, fi os.FileInfo) (stop bool, err error)) (err error) {
	err = ForDirMax(root, 0, -1, cb)
	return
}

// ForDirMax walks on `root` directory and its children with nested levels up to `maxLength`.
//
// Example - discover folder just one level
//
//      _ = ForDirMax(dir, 0, 1, func(depth int, cwd string, fi os.FileInfo) (stop bool, err error) {
//			if fi.IsDir() {
//				return
//			}
//          // ... doing something for a file,
//			return
//		})
//
// maxDepth = -1: no limit.
// initialDepth: 0 if no idea.
func ForDirMax(root string, initialDepth, maxDepth int, cb func(depth int, cwd string, fi os.FileInfo) (stop bool, err error)) (err error) {
	if maxDepth > 0 && initialDepth >= maxDepth {
		return
	}

	var files []os.FileInfo
	files, err = ioutil.ReadDir(os.ExpandEnv(root))
	if err != nil {
		// Logger.Fatalf("error in ForDirMax(): %v", err)
		return
	}

	var stop bool
	for _, f := range files {
		//Logger.Printf("  - %v", f.Name())
		if stop, err = cb(initialDepth, root, f); stop {
			return
		}
		if err != nil {
			log.NewStdLogger().Errorf("error in ForDirMax().cb: %v", err)
		} else if f.IsDir() && (maxDepth <= 0 || (maxDepth > 0 && initialDepth+1 < maxDepth)) {
			dir := path.Join(root, f.Name())
			if err = ForDirMax(dir, initialDepth+1, maxDepth, cb); err != nil {
				log.NewStdLogger().Errorf("error in ForDirMax(): %v", err)
			}
		}
	}

	return
}
