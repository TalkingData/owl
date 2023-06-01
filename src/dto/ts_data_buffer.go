package dto

import (
	commonpb "owl/common/proto"
	"sync"
)

type TsDataBuffer struct {
	content []*commonpb.TsData

	mu sync.Mutex
}

func NewTsDataBuffer() *TsDataBuffer {
	return &TsDataBuffer{
		content: make([]*commonpb.TsData, 0),
	}
}

func (buf *TsDataBuffer) Put(data ...*commonpb.TsData) {
	buf.mu.Lock()
	defer buf.mu.Unlock()

	buf.content = append(buf.content, data...)
}

func (buf *TsDataBuffer) Get(size int) []*commonpb.TsData {
	buf.mu.Lock()
	defer buf.mu.Unlock()

	if size > len(buf.content) {
		size = len(buf.content)
	}

	batch := buf.content[:size]
	buf.content = buf.content[size:]
	return batch
}

func (buf *TsDataBuffer) Len() int {
	buf.mu.Lock()
	defer buf.mu.Unlock()

	return len(buf.content)
}
