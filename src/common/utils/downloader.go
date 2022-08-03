package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	DefaultDownloaderEndOfFileExitCode = 1024
	DefaultDownloaderNormallyExitCode  = 1025
)

var (
	ErrEndOfFileExit = errors.New(fmt.Sprintf("%d", DefaultDownloaderEndOfFileExitCode))
	ErrNormallyExit  = errors.New(fmt.Sprintf("%d", DefaultDownloaderNormallyExitCode))
)

const (
	defaultBufferSize = 4096
)

type Downloader struct {
	wg sync.WaitGroup
	mu sync.Mutex

	exitChans []chan int
}

func (d *Downloader) Download(dst string, fn func(buffer []byte) error) error {
	d.wg.Add(1)
	defer d.wg.Done()

	exitChan := d.NewExitChan()
	defer func() {
		d.FreeExitChan(exitChan)
	}()

	fp, err := os.Open(dst)
	if err != nil {
		return err
	}
	defer func() {
		_ = fp.Close()
	}()

	n := 0
	buffer := make([]byte, defaultBufferSize)

	for {
		select {
		case <-exitChan:
			return ErrNormallyExit
		default:
			n, err = fp.Read(buffer)
			if err != nil {
				if err != io.EOF {
					return err
				}
			}
			if n > 0 {
				err = fn(buffer[:n])
				if err != nil {
					return err
				}
			}
		}
		// 收到包的size小于defaultBufferSize的，说明接收已经完成
		if n < defaultBufferSize {
			return ErrEndOfFileExit
		}
	}
}

func (d *Downloader) WaitForDone() {
	d.wg.Wait()
}

func (d *Downloader) CancelAllAndWaitForDone() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.SendToExitChans()
	d.wg.Wait()
	d.freeAllExitChans()
}

func (d *Downloader) NewExitChan() chan int {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.newExitChan()
}

func (d *Downloader) FreeExitChan(c chan int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.freeExitChan(c)
}

func (d *Downloader) FreeAllExitChans() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.freeAllExitChans()
}

func (d *Downloader) SendToExitChans() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i := range d.exitChans {
		go func(i int) {
			d.exitChans[i] <- 1
		}(i)
	}
}

func (d *Downloader) newExitChan() chan int {
	if d.exitChans == nil {
		d.exitChans = make([]chan int, 0)
	}

	chn := make(chan int, 1)
	d.exitChans = append(d.exitChans, chn)

	return chn
}

func (d *Downloader) freeExitChan(c chan int) {
	for i, chn := range d.exitChans {
		if chn == c {
			d.exitChans[i] = nil
			d.exitChans = append(d.exitChans[:i], d.exitChans[i+1:]...)
			break
		}
	}
}

func (d *Downloader) freeAllExitChans() {
	for i := range d.exitChans {
		d.exitChans[i] = nil
	}
	d.exitChans = nil
}
