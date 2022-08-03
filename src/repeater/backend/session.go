package backend

import (
	"net"
	"sync/atomic"
)

// session struct
type session struct {
	conn     *net.TCPConn
	exitFlag int32
}

func newSession() *session {
	return new(session)
}

func (s *session) Write(b []byte) (int, error) {
	return s.conn.Write(b)
}

func (s *session) Close() {
	atomic.StoreInt32(&s.exitFlag, 0)
}

func (s *session) IsClosed() bool {
	return atomic.LoadInt32(&s.exitFlag) == 0
}
