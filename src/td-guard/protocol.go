package main

import (
	"encoding/binary"
	"errors"
	"io"
)

var ErrPacketTooLarger = errors.New("the size of packet is larger than the limit")

type protocol struct {
}

func (this *protocol) ReadPacket(r io.Reader, packetLimitSize uint32) ([]byte, error) {
	var (
		head   = make([]byte, 4)
		length uint32
	)
	if _, err := io.ReadFull(r, head); err != nil {
		return nil, err
	}

	if length = binary.BigEndian.Uint32(head); length > packetLimitSize {
		return nil, ErrPacketTooLarger
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}
