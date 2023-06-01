package backend

import (
	"fmt"
	"github.com/silenceper/pool"
	"net"
	"owl/dto"
)

// kairos struct
type kairos struct {
	connPool pool.Pool
}

func newKairos(addr string, maxIdleConns, maxOpenConns int) (Backend, error) {
	pConf := &pool.Config{
		InitialCap: maxIdleConns,
		MaxIdle:    maxIdleConns,
		MaxCap:     maxOpenConns,
		Factory:    connKairosDb(addr),
		Close:      closeKairosDb,
	}

	p, err := pool.NewChannelPool(pConf)
	if err != nil {
		return nil, err
	}

	return &kairos{
		connPool: p,
	}, nil
}

func (kdb *kairos) Write(data *dto.TsData) error {
	v, err := kdb.connPool.Get()
	if err != nil {
		return err
	}
	defer func() {
		_ = kdb.connPool.Put(v)
	}()

	content := []byte(fmt.Sprintf("put %s %d %f %s\n",
		data.Metric,
		data.Timestamp,
		data.Value,
		data.Tags2Str(" "),
	))

	_, err = v.(net.Conn).Write(content)
	if err != nil {
		_ = kdb.connPool.Close(v)
	}

	return err
}

func (kdb *kairos) Close() {
	kdb.connPool.Release()
}

func connKairosDb(addr string) func() (interface{}, error) {
	return func() (interface{}, error) {
		kDbAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			return nil, err
		}

		return net.DialTCP("tcp", nil, kDbAddr)
	}
}

func closeKairosDb(v interface{}) error {
	return v.(net.Conn).Close()
}
