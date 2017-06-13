package utils

import (
	"bytes"
	"errors"
	"os/exec"
	"time"
)

var ErrRunTimeout = errors.New("command execution timeouts")

//如果timeout为0， 则没有超时
func RunCmdWithTimeout(cmd string, args []string, timeout int) ([]byte, error) {
	var (
		buf  bytes.Buffer
		err  error
		done chan error = make(chan error)
	)
	c := exec.Command(cmd, args...)
	c.Stdout = &buf
	c.Stderr = &buf
	if err = c.Start(); err != nil {
		return nil, err
	}
	go func() {
		done <- c.Wait()
	}()

	if timeout == 0 {
		err = <-done
		return buf.Bytes(), err
	} else {
		select {
		case <-time.After(time.Second * time.Duration(timeout)):
			if err = c.Process.Kill(); err != nil {
				return nil, err
			}
			return nil, ErrRunTimeout
		case err = <-done:
			return buf.Bytes(), err
		}
	}
}

func RunShell(cmd string) ([]byte, error) {
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	return out, err
}
