package main

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"
)

func GuardHandle() {

START:
	conn, err := net.DialTimeout("tcp", cfg.GUARDBIND, time.Second*5)
	if err != nil {
		log.Warn("connect to guard error %s", err.Error())
		StartGuard()
		time.Sleep(time.Second * 10)
		goto START
	}
	for {
		_, err := conn.Write(NewGuardHBPacket())
		if err != nil {
			log.Warn("send packet to guard error %s", err.Error())
			StartGuard()
			time.Sleep(time.Second * 10)
			goto START
		}
		time.Sleep(time.Second * 1)
	}

}

func NewGuardHBPacket() []byte {
	var buf bytes.Buffer
	head := make([]byte, 4)
	binary.BigEndian.PutUint32(head, uint32(1))
	binary.Write(&buf, binary.BigEndian, head)
	binary.Write(&buf, binary.BigEndian, byte(HOSTHB))
	return buf.Bytes()
}
