package tcpserver

import (
	"errors"
	"net"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	ErrConnClosing   = errors.New("use of closed network connection")
	ErrWriteBlocking = errors.New("write packet was blocking")
)

type Conn struct {
	server            *Server
	conn              *net.TCPConn
	closeChan         chan struct{}
	closeOnce         sync.Once
	closeFlag         int32
	packetSendChan    chan []byte
	packetReceiveChan chan []byte
}

func newConn(conn *net.TCPConn, server *Server) *Conn {

	return &Conn{
		server:            server,
		conn:              conn,
		closeChan:         make(chan struct{}),
		packetSendChan:    make(chan []byte, server.packetLimitSize),
		packetReceiveChan: make(chan []byte, server.packetReceiveChanLimit),
	}
}

func (this *Conn) Run() {
	this.server.waitGroup.Add(3)
	go this.witeLoop()   //从packetSendChan中读取数据并发送到对端
	go this.readLoop()   //从网络连接中读取数据包并放入packetReceiveChan
	go this.handleLoop() //从packetReceiveChan中读取数据并处理
	this.server.handler.Connect(this)
}

func (this *Conn) Close() {
	this.closeOnce.Do(func() {
		atomic.StoreInt32(&this.closeFlag, 1)
		close(this.closeChan)
		this.conn.Close()
		this.server.handler.Disconnect(this)
	})
}

func (this *Conn) IsClosed() bool {
	return atomic.LoadInt32(&this.closeFlag) == 1
}

func (this *Conn) handleLoop() {

	defer func() {
		recover()
		this.Close()
		this.server.waitGroup.Done()
	}()

	for {
		select {
		case <-this.closeChan:
			return
		case <-this.server.exitChan:
			return

		case p := <-this.packetReceiveChan:
			this.server.handler.HandlerMessage(this, p)
		}
	}
}

func (this *Conn) witeLoop() {
	defer func() {
		recover()
		this.Close()
		this.server.waitGroup.Done()
	}()

	for {
		select {
		case <-this.closeChan:
			return

		case <-this.server.exitChan:
			return

		case p := <-this.packetSendChan:
			if _, err := this.conn.Write(p); err != nil {
				return
			}
		}
	}
}

func (this *Conn) readLoop() {
	defer func() {
		recover()
		this.Close()
		this.server.waitGroup.Done()
	}()

	for {
		select {
		case <-this.closeChan:
			return
		case <-this.server.exitChan:
			return
		default:
		}
		p, err := this.server.protocol.ReadPacket(this.conn, this.server.packetLimitSize)
		if err != nil {
			return
		}
		this.packetReceiveChan <- p
	}
}

func (this *Conn) AsyncWriteData(data []byte) error {
	if this.IsClosed() {
		return ErrConnClosing
	}

	select {
	case this.packetSendChan <- data:
		return nil
	default:
		return ErrWriteBlocking
	}

}

func (this *Conn) GetLocalIp() string {
	return strings.Split(this.conn.LocalAddr().String(), ":")[0]
}

func (this *Conn) RemoteAddr() string {
	return this.conn.RemoteAddr().String()
}
