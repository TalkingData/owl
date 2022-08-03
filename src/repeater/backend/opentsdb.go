package backend

import (
	"errors"
	"fmt"
	"net"
	"owl/dto"
	"strings"
	"time"
)

// opentsdbBackend struct
type opentsdbBackend struct {
	tcpAddr *net.TCPAddr
	session *session
}

// newOpentsdbBackend
func newOpentsdbBackend(addr string) (Backend, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	bEnd := &opentsdbBackend{
		tcpAddr: tcpAddr,
		session: newSession(),
	}

	go bEnd.serve()

	return bEnd, nil
}

// Write
func (tsdb *opentsdbBackend) Write(data *dto.TsData) error {
	if tsdb.session.IsClosed() {
		return errors.New("backend session is closed.")
	}
	content := []byte(fmt.Sprintf("put %s %d %f %s\n",
		data.Metric,
		data.Timestamp,
		data.Value,
		strings.Replace(data.Tags2Str(), ",", " ", -1)))
	if _, err := tsdb.session.Write(content); err != nil {
		tsdb.session.Close()
		return err
	}
	return nil
}

// serve
func (tsdb *opentsdbBackend) serve() {
	var (
		err       error
		conn      *net.TCPConn
		tempDelay time.Duration
	)
retry:
	conn, err = net.DialTCP("tcp", nil, tsdb.tcpAddr)
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
	tsdb.session = &session{
		conn:     conn,
		exitFlag: 1,
	}
	for {
		if tsdb.session.IsClosed() {
			goto retry
		}
		time.Sleep(time.Second * 1)
	}
}
