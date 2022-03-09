package tcp

import (
	"net"
	"sync"
	"time"
)

// 已经迁移到 github.com/wuyingsong/tcp
// 测试稳定将移除

type Server struct {
	Addr    string //TCP address to listen on
	Listen  *net.TCPListener
	Handler Handler
	//	Protocol Protocol
	exitChan      chan struct{}
	maxPacketSize uint32
	wait          sync.WaitGroup
	Sessions      SessionBucket
}

type Handler interface {
	Handle(sess *Session, data []byte)
	MakeSession(sess *Session)
	LostSession(sess *Session)
}

func NewServer(addr string, handler Handler) *Server {
	return &Server{
		Addr:          addr,
		Handler:       handler,
		exitChan:      make(chan struct{}),
		maxPacketSize: 4096,
		Sessions: SessionBucket{
			m:    make(map[string]*Session),
			lock: sync.RWMutex{},
		},
	}
}

func (srv *Server) ListenAndServe() error {
	if srv.Addr == "" {
		srv.Addr = ":8888"
	}
	ln, err := net.Listen("tcp4", srv.Addr)
	if err != nil {
		return err
	}
	go srv.Serve(ln)
	return nil
}

func (srv *Server) Serve(l net.Listener) {
	defer l.Close()
	var (
		tempDelay time.Duration
	)
	for {
		conn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			return
		}
		tempDelay = 0
		if srv.Sessions.Get(conn.RemoteAddr().String()) != nil {
			conn.Close()
			continue
		}
		sess := srv.newSession(conn, nil)
		srv.Handler.MakeSession(sess)
		srv.Sessions.Add(sess)
	}
}

func (srv *Server) newSession(conn net.Conn, handler Handler) *Session {
	sess := &Session{
		srv:       srv,
		conn:      conn,
		handler:   handler,
		readChan:  make(chan []byte, READCHAN_SIZE),
		writeChan: make(chan []byte, WRITECHAN_SIZE),
		exitChan:  make(chan struct{}),
		exitFlag:  1,
	}
	if sess.handler == nil {
		sess.handler = srv.Handler
	}
	sess.Serve()
	return sess
}

func (srv *Server) Connect(addr string, handler Handler) (*Session, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Duration(time.Second*5))
	if err != nil {
		return nil, err
	}
	return srv.newSession(conn, handler), nil
}

//在ListenAndServe之前调用
func (srv *Server) SetMaxPacketSize(size uint32) {
	srv.maxPacketSize = size
}
