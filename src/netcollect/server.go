package main

import (
	"fmt"
	"net"
	"owl/common/types"
	"strings"
	"time"

	"github.com/wuyingsong/tcp"
	"github.com/wuyingsong/utils"
)

var (
	netCollect *NetCollect
)

type NetCollect struct {
	tsdBuffer chan *types.TimeSeriesData
	//historyMap   map[string]float64
	repeater     *tcp.TCPConn
	cfc          *tcp.TCPConn
	metricBuffer chan *types.MetricConfig
	switches     []*types.Switch
}

func InitNetCollect() error {
	c := &NetCollect{}
	c.tsdBuffer = make(chan *types.TimeSeriesData, GlobalConfig.BUFFER_SIZE)
	c.metricBuffer = make(chan *types.MetricConfig, GlobalConfig.BUFFER_SIZE)
	c.switches = []*types.Switch{}
	c.repeater = &tcp.TCPConn{}
	c.cfc = &tcp.TCPConn{}
	netCollect = c
	return nil
}

// 连接cfc
func (nc *NetCollect) dialCFC() error {
	if !nc.cfc.IsClosed() {
		return fmt.Errorf("cfc is already connected")
	}
	conn, err := nc.newTCPConn(GlobalConfig.CFC_ADDR)
	if err != nil {
		return err
	}
	nc.cfc = conn
	return nil
}

// 连接repeater
func (nc *NetCollect) dialRepeater() error {
	if !nc.repeater.IsClosed() {
		return fmt.Errorf("repeater is already connected")
	}
	conn, err := nc.newTCPConn(GlobalConfig.REPEATER_ADDR)
	if err != nil {
		return err
	}
	nc.repeater = conn
	return nil

}

func (nc *NetCollect) newTCPConn(addr string) (*tcp.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	tcpConn := tcp.NewTCPConn(conn, &callback{}, &tcp.DefaultProtocol{})
	return tcpConn, tcpConn.Serve()
}

// 守护cfc和repeater连接，失败重连
func (nc *NetCollect) watchConnLoop() {
	for {
		if nc.cfc.IsClosed() {
			lg.Error("cfc reconnect %v", nc.dialCFC())
		}
		if nc.repeater.IsClosed() {
			lg.Error("repeater reconnect %v", nc.dialRepeater())
		}
		time.Sleep(time.Second * 5)
	}
}

func InitIpRange() error {
	ipList := []string{}
	for _, ipRange := range GlobalConfig.IP_RANGE {
		var (
			ipStart string
			ipEnd   string
		)
		if len(ipRange) == 0 {
			continue
		}
		if strings.Contains(ipRange, "-") {
			field := strings.Split(ipRange, "-")
			if len(field) != 2 {
				continue
			}
			ipStart, ipEnd = field[0], field[1]
		} else {
			ipStart = ipRange
			ipEnd = ipRange
		}
		ips, err := utils.GetIPRange(ipStart, ipEnd)
		if err != nil {
			return err
		}
		ipList = append(ipList, ips...)
	}
	for _, ip := range ipList {
		s := &types.Switch{
			ID:              utils.Md5(ip)[:16],
			IP:              ip,
			Hostname:        "",
			LegalPrefix:     GlobalConfig.LEGAL_PREFIX,
			CollectInterval: GlobalConfig.COLLECT_INTERVAL,
			Snmp: types.SnmpConfig{
				Port:      GlobalConfig.SNMP_PORT,
				Version:   GlobalConfig.SNMP_VERSION,
				Community: GlobalConfig.SNMP_COMMUNITY,
				Timeout:   GlobalConfig.SNMP_TIMEOUT,
			},
		}
		s.GetHostname()
		s.GetVendor()
		netCollect.switches = append(netCollect.switches, s)
		lg.Info("do %s, %#v", s.IP, s.Snmp)
		go s.Do(netCollect.tsdBuffer, netCollect.metricBuffer)
	}
	netCollect.registerSwitchs()
	go netCollect.SendHeartbeat2CFC()
	go netCollect.SendMetri2CFC()

	return nil
}

//TODO:需要优化
// func (this *NetCollect) Dial(Type string) {
// 	var (
// 		tempDelay time.Duration
// 		err       error
// 		session   *tcp.TCPConn
// 	)
// retry:
// 	switch Type {
// 	case "cfc":
// 		tcp.AsyncTCPServer.Connect
// 		session, err = this.srv.Connect(GlobalConfig.CFC_ADDR, nil)
// 	case "repeater":
// 		session, err = this.srv.Connect(GlobalConfig.REPEATER_ADDR, nil)
// 	default:

// 	}
// 	if err != nil {
// 		lg.Error("%s: %s", Type, err.Error())
// 		if tempDelay == 0 {
// 			tempDelay = 5 * time.Millisecond
// 		} else {
// 			tempDelay *= 2
// 		}
// 		if max := 5 * time.Second; tempDelay > max {
// 			tempDelay = max
// 		}
// 		time.Sleep(tempDelay)
// 		goto retry
// 	}
// 	switch Type {
// 	case "cfc":
// 		this.cfc = session
// 	case "repeater":
// 		this.repeater = session
// 	}
// 	for {
// 		switch Type {
// 		case "cfc":
// 			if this.cfc.IsClosed() {
// 				goto retry
// 			}
// 		case "repeater":
// 			if this.repeater.IsClosed() {
// 				goto retry
// 			}
// 		}
// 		time.Sleep(time.Second * 1)
// 	}
// }

func (this *NetCollect) registerSwitchs() {
	for _, s := range this.switches {
		if !s.IsInit() {
			continue
		}
		h := &types.Host{
			ID:           s.ID,
			IP:           s.IP,
			Hostname:     s.Hostname,
			AgentVersion: Version,
			Metadata:     GlobalConfig.Metadata,
		}
		//config
		if err := this.cfc.AsyncWritePacket(
			tcp.NewDefaultPacket(
				types.MsgAgentRegister,
				h.Encode(),
			)); err != nil {
			lg.Error("register ")
		}
	}
}

func (this *NetCollect) SendHeartbeat2CFC() {
	for {
		if this.cfc.IsClosed() {
			goto sleep
		}
		for _, s := range this.switches {
			if !s.IsInit() {
				continue
			}
			h := &types.Host{
				ID:           s.ID,
				IP:           s.IP,
				Hostname:     s.Hostname,
				AgentVersion: Version,
				Metadata:     GlobalConfig.Metadata,
			}
			//config
			//this.cfc.AsyncWritePacket(
			//	tcp.NewDefaultPacket(
			//		types.MsgAgentRegister,
			//		h.Encode(),
			//	))

			//heartbeat
			this.cfc.AsyncWritePacket(
				tcp.NewDefaultPacket(
					types.MsgAgentSendHeartbeat,
					h.Encode(),
				),
			)
		}
	sleep:
		time.Sleep(time.Minute * 1)
	}
}

func (nc *NetCollect) SendMetri2CFC() {
	time.Sleep(time.Second * time.Duration(GlobalConfig.COLLECT_INTERVAL))
	for {
		if nc.cfc.IsClosed() {
			time.Sleep(time.Second * 30)
			continue
		}
		select {
		case tsd, ok := <-nc.metricBuffer:
			if ok {
				nc.cfc.AsyncWritePacket(
					tcp.NewDefaultPacket(
						types.MsgAgentSendMetricInfo,
						tsd.Encode(),
					),
				)
				lg.Info("sender to cfc %s", tsd)
			}
		}
	}
}

func (nc *NetCollect) SendTSD2Repeater() {
	for {
		if nc.repeater.IsClosed() {
			time.Sleep(time.Second * 5)
			continue
		}
		select {
		case tsd, ok := <-nc.tsdBuffer:
			if ok {
				nc.repeater.AsyncWritePacket(
					tcp.NewDefaultPacket(
						types.MsgAgentSendTimeSeriesData,
						tsd.Encode(),
					))
				lg.Info("sender to repeater %s", tsd)
			}
		}
	}
}
