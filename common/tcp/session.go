package tcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	READCHAN_SIZE  = 1024
	WRITECHAN_SIZE = 0
)

var (
	ErrSessionClosing = errors.New("use of closed network connection")
	ErrBufferFull     = errors.New("the async buffer is full")
)

type Session struct {
	srv       *Server
	handler   Handler
	conn      net.Conn
	readChan  chan []byte
	writeChan chan []byte
	exitChan  chan struct{}
	execOnce  sync.Once
	exitFlag  int32
}

func (sess *Session) Serve() {
	go sess.readLoop()
	go sess.writeLoop()
	go sess.handleLoop()
}

func (sess *Session) readLoop() {
	defer func() {
		recover()
		sess.Close()
	}()

	var (
		err    error
		length uint32
	)
	head := make([]byte, 4)
	for {
		select {
		case <-sess.srv.exitChan:
			return
		case <-sess.exitChan:
			return
		default:
			if _, err = io.ReadFull(sess.conn, head); err != nil {
				return
			}
			if length = binary.BigEndian.Uint32(head); length > sess.srv.maxPacketSize {
				return
			}
			buf := make([]byte, length)
			if _, err = io.ReadFull(sess.conn, buf); err != nil {
				return
			}
			sess.readChan <- buf
		}
	}
}

func (sess *Session) writeLoop() {
	defer func() {
		recover()
		sess.Close()
	}()

	var buf bytes.Buffer
	head := make([]byte, 4)

	for {
		select {
		case <-sess.srv.exitChan:
			return
		case <-sess.exitChan:
			return
		case data := <-sess.writeChan:
			binary.BigEndian.PutUint32(head, uint32(len(data)))
			binary.Write(&buf, binary.BigEndian, head)
			binary.Write(&buf, binary.BigEndian, data)
			if _, err := sess.conn.Write(buf.Bytes()); err != nil {
				fmt.Println("writeLoop: %s", err.Error())
				return
			}
			buf.Reset()
		}
	}
}

func (sess *Session) handleLoop() {
	defer func() {
		recover()
		sess.Close()
	}()
	for {
		select {
		case _ = <-sess.srv.exitChan:
			return
		case _ = <-sess.exitChan:
			return
		case data := <-sess.readChan:
			sess.srv.Handler.Handle(sess, data)
		}
	}
}

func (sess *Session) Send(data []byte) error {
	if sess.IsClosed() {
		return ErrSessionClosing
	}
	select {
	case sess.writeChan <- data:
		return nil
	case <-time.After(time.Second * 1):
		return ErrBufferFull
	}
}

func (sess *Session) Close() {
	sess.execOnce.Do(func() {
		sess.srv.Handler.LostSession(sess)
		close(sess.exitChan)
		atomic.StoreInt32(&sess.exitFlag, 0)
		sess.srv.Sessions.Delete(sess.RemoteIPAddr())
		sess.conn.Close()
	})
}

func (sess *Session) IsClosed() bool {
	return atomic.LoadInt32(&sess.exitFlag) == 0
}

func (sess *Session) LocalAddr() string {
	return sess.conn.LocalAddr().String()
}

func (sess *Session) RemoteAddr() string {
	return sess.conn.RemoteAddr().String()
}

func (sess *Session) RemoteIPAddr() string {
	return strings.Split(sess.RemoteAddr(), ":")[0]
}
