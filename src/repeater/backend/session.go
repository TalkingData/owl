package backend

import (
	"net"
	"sync/atomic"
)

type session struct {
	conn     *net.TCPConn
	exitFlag int32
}

func (this *session) Write(b []byte) (int, error) {
	return this.conn.Write(b)
}

func (this *session) Close() {
	atomic.StoreInt32(&this.exitFlag, 0)
}

func (this *session) IsClosed() bool {
	return atomic.LoadInt32(&this.exitFlag) == 0
}
