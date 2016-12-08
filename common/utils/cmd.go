package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

var ErrRunTimeout = errors.New("command execution timeouts")

func RunCmdWithTimeout(cmd string, args []string, timeout int) ([]byte, error) {
	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
		err    error
		done   chan error = make(chan error)
	)
	c := exec.Command(cmd, args...)
	c.Stdout = &stdout
	c.Stderr = &stderr
	if err = c.Start(); err != nil {
		return nil, err
	}

	go func() {
		done <- c.Wait()
	}()

	select {
	case <-time.After(time.Second * time.Duration(timeout)):
		if err = c.Process.Kill(); err != nil {
			return nil, err
		}
		return nil, ErrRunTimeout
	case err = <-done:
		if err != nil {
			return nil, fmt.Errorf("%s %s", stderr.Bytes(), err)
		}
		return stdout.Bytes(), nil
	}
}

func RunShell(cmd string) ([]byte, error) {
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	return out, err
}
