package tcp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrConnClosing = errors.New("use of closed network connection")
	ErrBufferFull  = errors.New("the async send buffer is full")
)

type TCPConn struct {
	callback CallBack
	protocol Protocol

	conn      *net.TCPConn
	readChan  chan Packet
	writeChan chan Packet

	readDeadline  time.Duration
	writeDeadline time.Duration

	exitChan  chan struct{}
	closeOnce sync.Once
	exitFlag  int32
}

func NewTCPConn(conn *net.TCPConn, callback CallBack, protocol Protocol) *TCPConn {
	c := &TCPConn{
		conn:     conn,
		callback: callback,
		protocol: protocol,

		readChan:  make(chan Packet, readChanSize),
		writeChan: make(chan Packet, writeChanSize),

		exitChan: make(chan struct{}),
		exitFlag: 0,
	}
	return c
}

func (c *TCPConn) Serve() error {
	defer func() {
		if r := recover(); r != nil {
			logger.Println("tcp conn(%v) Serve error, %v ", c.GetRemoteIPAddress(), r)
		}
	}()
	if c.callback == nil || c.protocol == nil {
		err := fmt.Errorf("callback and protocol are not allowed to be nil")
		c.Close()
		return err
	}
	atomic.StoreInt32(&c.exitFlag, 1)
	c.callback.OnConnected(c)
	go c.readLoop()
	go c.writeLoop()
	go c.handleLoop()
	return nil
}

func (c *TCPConn) readLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.exitChan:
			return
		default:
			if c.readDeadline > 0 {
				c.conn.SetReadDeadline(time.Now().Add(c.readDeadline))
			}
			p, err := c.protocol.ReadPacket(c.conn)
			if err != nil {
				if err != io.EOF {
					c.callback.OnError(err)
				}
				return
			}
			c.readChan <- p
		}
	}
}

func (c *TCPConn) ReadPacket() (Packet, error) {
	if c.protocol == nil {
		return nil, errors.New("no protocol impl")
	}
	return c.protocol.ReadPacket(c.conn)
}

func (c *TCPConn) writeLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for pkt := range c.writeChan {
		if pkt == nil {
			continue
		}
		if c.writeDeadline > 0 {
			c.conn.SetWriteDeadline(time.Now().Add(c.writeDeadline))
		}
		if err := c.protocol.WritePacket(c.conn, pkt); err != nil {
			c.callback.OnError(err)
			return
		}
	}
}

func (c *TCPConn) handleLoop() {
	defer func() {
		recover()
		c.Close()
	}()
	for p := range c.readChan {
		if p == nil {
			continue
		}
		c.callback.OnMessage(c, p)
	}
}

func (c *TCPConn) AsyncWritePacket(p Packet) error {
	if c.IsClosed() {
		return ErrConnClosing
	}
	select {
	case c.writeChan <- p:
		return nil
	default:
		return ErrBufferFull
	}
}

func (c *TCPConn) Close() {
	c.closeOnce.Do(func() {
		atomic.StoreInt32(&c.exitFlag, 0)
		close(c.exitChan)
		close(c.writeChan)
		close(c.readChan)
		if c.callback != nil {
			c.callback.OnDisconnected(c)
		}
		c.conn.Close()
	})
}

func (c *TCPConn) GetRawConn() *net.TCPConn {
	return c.conn
}

func (c *TCPConn) IsClosed() bool {
	return atomic.LoadInt32(&c.exitFlag) == 0
}

func (c *TCPConn) GetLocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

//LocalIPAddress 返回socket连接本地的ip地址
func (c *TCPConn) GetLocalIPAddress() string {
	return strings.Split(c.GetLocalAddr().String(), ":")[0]
}

func (c *TCPConn) GetRemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *TCPConn) GetRemoteIPAddress() string {
	return strings.Split(c.GetRemoteAddr().String(), ":")[0]
}

func (c *TCPConn) setReadDeadline(t time.Duration) {
	c.readDeadline = t
}

func (c *TCPConn) setWriteDeadline(t time.Duration) {
	c.writeDeadline = t
}
