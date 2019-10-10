package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/logger"
)

// Command is a simpler struct for defining commands than Go's built-in Cmd.
type Command struct {
	Command           string            // The command to run
	Args              []string          // The args to pass to the command
	WorkingDir        string            // The working directory
	Env               map[string]string // Additional environment variables to set
	OutputMaxLineSize int               // The max line size of stdout and stderr (in bytes)
}

// RunCommand runs a shell command and redirects its stdout and stderr to the stdout of the atomic script itself.
func RunCommand(t *testing.T, command Command) {
	err := RunCommandE(t, command)
	if err != nil {
		t.Fatal(err)
	}
}

// RunCommandE runs a shell command and redirects its stdout and stderr to the stdout of the atomic script itself.
func RunCommandE(t *testing.T, command Command) error {
	_, err := RunCommandAndGetOutputE(t, command)
	return err
}

// RunCommandAndGetOutput runs a shell command and returns its stdout and stderr as a string. The stdout and stderr of that command will also
// be printed to the stdout and stderr of this Go program to make debugging easier.
func RunCommandAndGetOutput(t *testing.T, command Command) string {
	out, err := RunCommandAndGetOutputE(t, command)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// RunCommandAndGetOutputE runs a shell command and returns its stdout and stderr as a string. The stdout and stderr of that command will also
// be printed to the stdout and stderr of this Go program to make debugging easier.
func RunCommandAndGetOutputE(t *testing.T, command Command) (string, error) {
	allOutput := []string{}
	err := runCommandAndStoreOutputE(t, command, &allOutput, &allOutput)

	output := strings.Join(allOutput, "\n")
	return output, err
}

// RunCommandAndGetStdOut runs a shell command and returns solely its stdout (but not stderr) as a string. The stdout
// and stderr of that command will also be printed to the stdout and stderr of this Go program to make debugging easier.
// If there are any errors, fail the test.
func RunCommandAndGetStdOut(t *testing.T, command Command) string {
	output, err := RunCommandAndGetStdOutE(t, command)
	require.NoError(t, err)
	return output
}

// RunCommandAndGetStdOutE runs a shell command and returns solely its stdout (but not stderr) as a string. The stdout
// and stderr of that command will also be printed to the stdout and stderr of this Go program to make debugging easier.
func RunCommandAndGetStdOutE(t *testing.T, command Command) (string, error) {
	stdout := []string{}
	stderr := []string{}
	err := runCommandAndStoreOutputE(t, command, &stdout, &stderr)

	output := strings.Join(stdout, "\n")
	return output, err
}

// runCommandAndStoreOutputE runs a shell command and stores each line from stdout and stderr in the given
// storedStdout and storedStderr variables, respectively. The stdout and stderr of that command will also
// be printed to the stdout and stderr of this Go program to make debugging easier.
func runCommandAndStoreOutputE(t *testing.T, command Command, storedStdout *[]string, storedStderr *[]string) error {
	logger.Logf(t, "Running command %s with args %s", command.Command, command.Args)

	cmd := exec.Command(command.Command, command.Args...)
	cmd.Dir = command.WorkingDir
	cmd.Stdin = os.Stdin
	cmd.Env = formatEnvVars(command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	if err := readStdoutAndStderr(t, stdout, stderr, storedStdout, storedStderr, command.OutputMaxLineSize); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

// This function captures stdout and stderr into the given variables while still printing it to the stdout and stderr
// of this Go program
func readStdoutAndStderr(t *testing.T, stdout io.ReadCloser, stderr io.ReadCloser, storedStdout *[]string, storedStderr *[]string, maxLineSize int) error {
	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	if maxLineSize > 0 {
		stdoutScanner.Buffer(make([]byte, maxLineSize), maxLineSize)
		stderrScanner.Buffer(make([]byte, maxLineSize), maxLineSize)
	}

	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}
	wg.Add(2)
	go readData(t, stdoutScanner, wg, mutex, storedStdout)
	go readData(t, stderrScanner, wg, mutex, storedStderr)
	wg.Wait()

	if err := stdoutScanner.Err(); err != nil {
		return err
	}

	if err := stderrScanner.Err(); err != nil {
		return err
	}

	return nil
}

func readData(t *testing.T, scanner *bufio.Scanner, wg *sync.WaitGroup, mutex *sync.Mutex, allOutput *[]string) {
	defer wg.Done()
	for scanner.Scan() {
		logTextAndAppendToOutput(t, mutex, scanner.Text(), allOutput)
	}
}

func logTextAndAppendToOutput(t *testing.T, mutex *sync.Mutex, text string, allOutput *[]string) {
	defer mutex.Unlock()
	logger.Log(t, text)
	mutex.Lock()
	*allOutput = append(*allOutput, text)
}

// GetExitCodeForRunCommandError tries to read the exit code for the error object returned from running a shell command. This is a bit tricky to do
// in a way that works across platforms.
func GetExitCodeForRunCommandError(err error) (int, error) {
	// http://stackoverflow.com/a/10385867/483528
	if exitErr, ok := err.(*exec.ExitError); ok {
		// The program has exited with an exit code != 0

		// This works on both Unix and Windows. Although package
		// syscall is generally platform dependent, WaitStatus is
		// defined for both Unix and Windows and in both cases has
		// an ExitStatus() method with the same signature.
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus(), nil
		}
		return 1, errors.New("could not determine exit code")
	}

	return 0, nil
}

func formatEnvVars(command Command) []string {
	env := os.Environ()
	for key, value := range command.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	return env
}
