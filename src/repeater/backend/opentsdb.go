package backend

import (
	"errors"
	"fmt"
	"net"
	"owl/common/types"
	"strings"
	"time"
)

type OpentsdbBackend struct {
	tcpAddr *net.TCPAddr
	session *session
}

func NewOpentsdbBackend(addr string) (*OpentsdbBackend, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	backend := new(OpentsdbBackend)
	backend.tcpAddr = tcpAddr
	backend.session = new(session)
	go backend.serve()
	return backend, nil
}

func (this *OpentsdbBackend) serve() {
	var (
		err       error
		conn      *net.TCPConn
		tempDelay time.Duration
	)
retry:
	conn, err = net.DialTCP("tcp", nil, this.tcpAddr)
	if err != nil {
		if tempDelay == 0 {
			tempDelay = 5 * time.Millisecond
		} else {
			tempDelay *= 2
		}
		if max := 5 * time.Second; tempDelay > max {
			tempDelay = max
		}
		time.Sleep(tempDelay)
		goto retry
	}
	this.session = &session{
		conn:     conn,
		exitFlag: 1,
	}
	for {
		if this.session.IsClosed() {
			goto retry
		}
		time.Sleep(time.Second * 1)
	}
}

func (this *OpentsdbBackend) Write(data *types.TimeSeriesData) error {
	if this.session.IsClosed() {
		return errors.New("backend session is closed.")
	}
	content := []byte(fmt.Sprintf("put %s %d %f %s\n",
		data.Metric,
		data.Timestamp,
		data.Value,
		strings.Replace(data.Tags2String(), ",", " ", -1)))
	if _, err := this.session.Write(content); err != nil {
		this.session.Close()
		return err
	}
	return nil
}
