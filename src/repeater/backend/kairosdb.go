package backend

import (
	"errors"
	"fmt"
	"net"
	"owl/dto"
	"strings"
	"time"
)

// kairosDbBackend struct
type kairosDbBackend struct {
	tcpAddr *net.TCPAddr
	session *session
}

func newKairosDbBackend(addr string) (Backend, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	bEnd := &kairosDbBackend{
		tcpAddr: tcpAddr,
		session: newSession(),
	}

	go bEnd.serve()

	return bEnd, nil
}

func (kDb *kairosDbBackend) Write(data *dto.TsData) error {
	if kDb.session.IsClosed() {
		return errors.New("backend session is closed.")
	}
	content := []byte(fmt.Sprintf("put %s %d %f %s\n",
		data.Metric,
		data.Timestamp,
		data.Value,
		strings.Replace(data.Tags2Str(), ",", " ", -1)))
	if _, err := kDb.session.Write(content); err != nil {
		kDb.session.Close()
		return err
	}
	return nil
}

// serve
func (kDb *kairosDbBackend) serve() {
	var (
		err       error
		conn      *net.TCPConn
		tempDelay time.Duration
	)
retry:
	conn, err = net.DialTCP("tcp", nil, kDb.tcpAddr)
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
	kDb.session = &session{
		conn:     conn,
		exitFlag: 1,
	}
	for {
		if kDb.session.IsClosed() {
			goto retry
		}
		time.Sleep(time.Second * 1)
	}
}
