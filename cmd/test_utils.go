package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"time"
)

// execCommand is a package-level variable to allow mocking exec.Command in tests.
var execCommand = exec.Command

// timeNow is a package-level variable to allow mocking time.Now in tests.
var timeNow = time.Now

// osExit is a package-level variable to allow mocking os.Exit in tests.
var osExit = os.Exit

// captureOutput is a generic helper to capture output from a given *os.File (e.g., os.Stdout, os.Stderr).
// It returns a bytes.Buffer to read the output from, and a cleanup function.
func captureOutput(target *os.File) (outputBuffer *bytes.Buffer, cleanup func()) {
	oldTarget := target
	r, w, _ := os.Pipe()
	*target = *r // Redirect target (os.Stdout or os.Stderr) to the read end of the pipe

	outputBuffer = new(bytes.Buffer)
	// Start a goroutine to copy data from the read end of the pipe to the buffer
	// and then restore the original target.
	go func() {
		io.Copy(outputBuffer, r)
		r.Close()
	}()

	cleanup = func() {
		w.Close() // Close the write end of the pipe
		*target = *oldTarget // Restore the original target
	}

	return outputBuffer, cleanup
}

// captureStdout captures the output that writes to os.Stdout.
// It returns the captured string.
func captureStdout(f func()) string {
	outputBuffer, cleanup := captureOutput(os.Stdout)
	defer cleanup()
	f()
	return outputBuffer.String()
}

// captureStderr captures the output that writes to os.Stderr.
// It returns the captured string.
func captureStderr(f func()) string {
	outputBuffer, cleanup := captureOutput(os.Stderr)
	defer cleanup()
	f()
	return outputBuffer.String()
}

// simulateInput is a helper function to simulate user input for os.Stdin.
func simulateInput(input string) func() {
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	go func() {
		defer w.Close()
		io.WriteString(w, input)
	}()
	return func() {
		os.Stdin = oldStdin
	}
}