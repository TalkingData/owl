package backend

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"owl/common/types"
	"time"

	"github.com/wuyingsong/tcp"
)

type repeaterBackend struct {
	tcpAddr *net.TCPAddr
	session *session
}

func NewRepeaterBackend(addr string) (*repeaterBackend, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	backend := new(repeaterBackend)
	backend.tcpAddr = tcpAddr
	backend.session = new(session)
	go backend.serve()
	return backend, nil
}

func (this *repeaterBackend) serve() {
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

func (this *repeaterBackend) Write(data *types.TimeSeriesData) error {
	if this.session.IsClosed() {
		return errors.New("backend session is closed.")
	}
	pkt := tcp.NewDefaultPacket(types.MsgRepeaterPostTimeSeriesData, data.Encode())
	head := make([]byte, 4)
	var buf bytes.Buffer
	pktBytes := pkt.Bytes()
	binary.BigEndian.PutUint32(head, uint32(len(pktBytes)))
	binary.Write(&buf, binary.BigEndian, head)
	binary.Write(&buf, binary.BigEndian, pktBytes)
	if _, err := this.session.Write(buf.Bytes()); err != nil {
		this.session.Close()
		return err
	}
	return nil
}
