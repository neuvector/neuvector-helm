// Package logger contains different methods to log.
package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

// Logf logs the given format and arguments, formatted using fmt.Sprintf, to stdout, along with a timestamp and information
// about what test and file is doing the logging. This is an alternative to t.Logf that logs to stdout immediately,
// rather than buffering all log output and only displaying it at the very end of the test. This is useful because:
//
// 1. It allows you to iterate faster locally, as you get feedback on whether your code changes are working as expected
//    right away, rather than at the very end of the test run.
//
// 2. If you have a bug in your code that causes a test to never complete or if the test code crashes,, t.Logf would
//    show you no log output whatsoever, making debugging very hard, where as this method will show you all the log
//    output available.
//
// 3. If you have a test that takes a long time to complete, some CI systems will kill the test suite prematurely
//    because there is no log output with t.Logf (e.g., CircleCI kills tests after 10 minutes of no log output). With
//    this log method, you get log output continuously.
//
// Note that there is a proposal to improve t.Logf (https://github.com/golang/go/issues/24929), but until that's
// implemented, this method is our best bet.
func Logf(t *testing.T, format string, args ...interface{}) {
	DoLog(t, 2, os.Stdout, fmt.Sprintf(format, args...))
}

// Log logs the given arguments to stdout, along with a timestamp and information about what test and file is doing the
// logging. This is an alternative to t.Logf that logs to stdout immediately, rather than buffering all log output and
// only displaying it at the very end of the test. See the Logf method for more info.
func Log(t *testing.T, args ...interface{}) {
	DoLog(t, 2, os.Stdout, args...)
}

// DoLog logs the given arguments to the given writer, along with a timestamp and information about what test and file is
// doing the logging.
func DoLog(t *testing.T, callDepth int, writer io.Writer, args ...interface{}) {
	date := time.Now()
	prefix := fmt.Sprintf("%s %s %s:", t.Name(), date.Format(time.RFC3339), CallerPrefix(callDepth+1))
	allArgs := append([]interface{}{prefix}, args...)
	fmt.Fprintln(writer, allArgs...)
}

// CallerPrefix returns the file and line number information about the methods that called this method, based on the current
// goroutine's stack. The argument callDepth is the number of stack frames to ascend, with 0 identifying the method
// that called CallerPrefix, 1 identifying the method that called that method, and so on.
//
// This code is adapted from testing.go, where it is in a private method called decorate.
func CallerPrefix(callDepth int) string {
	_, file, line, ok := runtime.Caller(callDepth)
	if ok {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = 1
	}

	return fmt.Sprintf("%s:%d", file, line)
}
