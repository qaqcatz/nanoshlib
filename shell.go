package nanoshlib

import (
	"bytes"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// Exec use /bin/bash -c to execute cmdStr, wait for the result, or timeout, return out stream, error stream,
// and an error, which can be nil, normal error or *TimeoutError.
//
// timeoutMS <= 0 means timeoutMS = inf
func Exec(cmdStr string, timeoutMS int) (string, string, error) {
	// child process
	cmd := exec.Command("/bin/bash", "-c", cmdStr)

	// Use a bytes.Buffer to get the output
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	if err := cmd.Start(); err != nil {
		return "", "", err
	}

	// Use a channel to signal completion so we can use a select statement
	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	if timeoutMS > 0 {
		// Start a timer
		timeout := time.After(time.Duration(timeoutMS) * time.Millisecond)

		// The select statement allows us to execute based on which channel we get a message from first.
		select {
		case <-timeout:
			// Timeout happened first, kill the process and print a message.
			// The reason why I don't use context.WithTimeout() is that sometimes it can not kill the child process
			_ = cmd.Process.Kill()
			return outBuf.String(), errBuf.String(), &TimeoutError{}
		case err := <-done:
			// Command completed before timeout. Print output and error if it exists.
			if err != nil {
				// This branch means that the return value of cmd != 0
				return outBuf.String(), errBuf.String(), err
			}
			return outBuf.String(), errBuf.String(), nil
		}
	} else {
		select {
		case err := <-done:
			// Command completed before timeout. Print output and error if it exists.
			if err != nil {
				// This branch means that the return value of cmd != 0
				return outBuf.String(), errBuf.String(), err
			}
			return outBuf.String(), errBuf.String(), nil
		}
	}
}

// Exec0 is an extension of Exec. When Exec0 is executed, it will return immediately.
// Exec0 does not care about the result, it is only responsible for
// making sure the process can be started/killed successfully and
// checking if the process is still running.
// It can be used to start/monitor/kill a service.
//
// Specifically, Exec0 will return doneChan, killChan, err:
//
// - doneChan, you can use select case: <-errChan default: ... to check if the process is still running
//
// - killChan, you can use killChan<-0(any number) to kill the process.
//
// - err, command start error.
func Exec0(myCmdStr string) (chan error, chan int, error) {
	// child process
	cmd := exec.Command("/bin/bash", "-c", myCmdStr)

	// Use a bytes.Buffer to get the output
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	// Use a channel to signal completion so we can use a select statement
	doneChan := make(chan error)
	go func() { doneChan <- cmd.Wait() }()

	// user can use killChan<-0(any number) to kill the process.
	killChan := make(chan int)
	go func() {
		select {
		case <- killChan:
			_ = cmd.Process.Kill()
		}
	} ()

	return doneChan, killChan, nil
}

// Exec0s is an extension of Exec0. We use cmd.SysProcAttr = &syscall.SysProcAttr{Setsid:true}
// to avoid the process being killed when the program ends.
func Exec0s(myCmdStr string) (chan error, chan int, error) {
	// child process
	cmd := exec.Command("/bin/bash", "-c", myCmdStr)

	// Use a bytes.Buffer to get the output
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid:true}

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	// Use a channel to signal completion so we can use a select statement
	doneChan := make(chan error)
	go func() { doneChan <- cmd.Wait() }()

	// user can use killChan<-0(any number) to kill the process.
	killChan := make(chan int)
	go func() {
		select {
		case <- killChan:
			_ = cmd.Process.Kill()
		}
	} ()

	return doneChan, killChan, nil
}

// ExecStd: likes Exec, but write the result to standard output stream and standard error stream.
func ExecStd(cmdStr string, timeoutMS int) error {
	// child process
	cmd := exec.Command("/bin/bash", "-c", cmdStr)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	// Use a channel to signal completion so we can use a select statement
	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	if timeoutMS > 0 {
		// Start a timer
		timeout := time.After(time.Duration(timeoutMS) * time.Millisecond)

		// The select statement allows us to execute based on which channel we get a message from first.
		select {
		case <-timeout:
			// Timeout happened first, kill the process and print a message.
			// The reason why I don't use context.WithTimeout() is that sometimes it can not kill the child process
			_ = cmd.Process.Kill()
			return &TimeoutError{}
		case err := <-done:
			// Command completed before timeout. Print output and error if it exists.
			if err != nil {
				// This branch means that the return value of cmd != 0
				return err
			}
			return nil
		}
	} else {
		select {
		case err := <-done:
			// Command completed before timeout. Print output and error if it exists.
			if err != nil {
				// This branch means that the return value of cmd != 0
				return err
			}
			return nil
		}
	}
}
