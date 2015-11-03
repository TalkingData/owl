package tcpserver

import (
	"net"
	"sync"
	"time"
)

type Config struct {
}

type Server struct {
	listener *net.TCPListener //网络监听

	acceptTimeOut time.Duration
	readTimeOut   time.Duration //网络数据包读取超时时间
	writeTimeOut  time.Duration //数据包写入超时时间

	packetLimitSize        uint32 //数据包最大值
	packetSendChanLimit    uint32 //数据包发送队列大小
	packetReceiveChanLimit uint32 //数据包接受队列大小

	exitChan  chan struct{}   //退出信号
	waitGroup *sync.WaitGroup //同步

	protocol Protocol //协议
	handler  Handler  //数据包处理
}

func NewServer(listener *net.TCPListener, handler Handler) *Server {

	return &Server{
		listener:               listener,
		acceptTimeOut:          30 * time.Second,
		readTimeOut:            5 * time.Minute,
		writeTimeOut:           5 * time.Minute,
		packetLimitSize:        4096,
		packetSendChanLimit:    20,
		packetReceiveChanLimit: 20,
		exitChan:               make(chan struct{}),
		waitGroup:              &sync.WaitGroup{},
		protocol:               Protocol{},
		handler:                handler,
	}
}

func (this *Server) SetAcceptTimeOut(d time.Duration) {
	this.acceptTimeOut = d
}

func (this *Server) SetReadTimeOut(d time.Duration) {
	this.readTimeOut = d
}

func (this *Server) SetWriteTimeOut(d time.Duration) {
	this.writeTimeOut = d
}

func (this *Server) SetPacketLimitSize(s uint32) {
	this.packetLimitSize = s
}

func (this *Server) SetPacketSendChanLimit(s uint32) {
	this.packetSendChanLimit = s
}

func (this *Server) SetPacketReceiveChanLimit(s uint32) {
	this.packetReceiveChanLimit = s
}

func (this *Server) Start() {
	this.waitGroup.Add(1)
	defer func() {
		this.listener.Close()
		this.waitGroup.Done()
	}()

	for {
		select {
		case <-this.exitChan:
			return
		default:
		}

		this.listener.SetDeadline(time.Now().Add(this.acceptTimeOut))

		conn, err := this.listener.AcceptTCP()
		if err != nil {
			continue
		}

		go newConn(conn, this).Run()
	}
}

func (this *Server) Stop() {
	close(this.exitChan)
	this.waitGroup.Wait()
}

func (this *Server) Connect(address string) (*Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return newConn(conn, this), nil
}
