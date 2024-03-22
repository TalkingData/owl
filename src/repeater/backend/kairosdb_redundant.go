package backend

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/silenceper/pool"
	"net"
	"owl/dto"
)

// kairosRedundant struct
type kairosRedundant struct {
	connPool1 pool.Pool
	connPool2 pool.Pool

	promMetric prometheus.Counter
}

func newKairosRedundant(addr1, addr2 string, maxIdleConns, maxOpenConns int) (Backend, error) {
	p1Conf := &pool.Config{
		InitialCap: maxIdleConns,
		MaxIdle:    maxIdleConns,
		MaxCap:     maxOpenConns,
		Factory:    connKairosDb(addr1),
		Close:      closeKairosDb,
	}

	p1, err := pool.NewChannelPool(p1Conf)
	if err != nil {
		return nil, err
	}

	p2Conf := &pool.Config{
		InitialCap: maxIdleConns,
		MaxIdle:    maxIdleConns,
		MaxCap:     maxOpenConns,
		Factory:    connKairosDbRe(addr2),
		Close:      closeKairosDbRe,
	}

	p2, err := pool.NewChannelPool(p2Conf)
	if err != nil {
		return nil, err
	}

	pm := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "write_kairosdb_redundant_count",
	})
	prometheus.MustRegister(pm)
	pm.Add(0)

	return &kairosRedundant{
		connPool1: p1,
		connPool2: p2,

		promMetric: pm,
	}, nil
}

func (kdbRe *kairosRedundant) Write(data *dto.TsData) error {
	kdbRe.promMetric.Inc()

	content := []byte(fmt.Sprintf("put %s %d %f %s\n",
		data.Metric,
		data.Timestamp,
		data.Value,
		data.Tags2Str(" "),
	))

	err1 := kdbRe.writeConnPool1(content)
	err2 := kdbRe.writeConnPool2(content)

	// 判断两个连接池的写入结果，如果有一个失败，则返回错误
	if err1 != nil || err2 != nil {
		return fmt.Errorf("write kairosdb redundant error, connPool1: %s, connPool2: %s", err1, err2)
	}

	return nil
}

func (kdbRe *kairosRedundant) writeConnPool1(content []byte) error {
	v, err := kdbRe.connPool1.Get()
	if err != nil {
		return err
	}
	defer func() {
		_ = kdbRe.connPool1.Put(v)
	}()

	_, err = v.(net.Conn).Write(content)
	if err != nil {
		_ = kdbRe.connPool1.Close(v)
	}

	return err
}

func (kdbRe *kairosRedundant) writeConnPool2(content []byte) error {
	v, err := kdbRe.connPool2.Get()
	if err != nil {
		return err
	}
	defer func() {
		_ = kdbRe.connPool2.Put(v)
	}()

	_, err = v.(net.Conn).Write(content)
	if err != nil {
		_ = kdbRe.connPool2.Close(v)
	}

	return err
}

func (kdbRe *kairosRedundant) Close() {
	kdbRe.connPool1.Release()
	kdbRe.connPool2.Release()
}

func connKairosDbRe(addr string) func() (interface{}, error) {
	return func() (interface{}, error) {
		kDbAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			return nil, err
		}

		return net.DialTCP("tcp", nil, kDbAddr)
	}
}

func closeKairosDbRe(v interface{}) error {
	return v.(net.Conn).Close()
}
