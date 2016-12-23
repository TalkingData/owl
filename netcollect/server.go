package main

import (
	"owl/common/tcp"
	"owl/common/types"
	"owl/common/utils"
	"strings"
	"time"
)

var (
	netCollect   *NetCollect
	AgentVersion string = "0.1"
)

type NetCollect struct {
	srv          *tcp.Server
	tsdBuffer    chan *types.TimeSeriesData
	histroyMap   map[string]float64
	repeater     *tcp.Session
	cfc          *tcp.Session
	metricBuffer chan *types.MetricConfig
	switchs      []*types.Switch
}

func InitNetCollect() error {
	s := tcp.NewServer("", &handle{})
	c := &NetCollect{}
	c.srv = s
	c.tsdBuffer = make(chan *types.TimeSeriesData, GlobalConfig.BUFFER_SIZE)
	c.metricBuffer = make(chan *types.MetricConfig, GlobalConfig.BUFFER_SIZE)
	c.switchs = []*types.Switch{}
	c.repeater = &tcp.Session{}
	c.cfc = &tcp.Session{}
	netCollect = c
	return nil
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
			},
		}
		netCollect.switchs = append(netCollect.switchs, s)
		lg.Info("do %s, %#v", s.IP, s.Snmp)
		go s.Do(netCollect.tsdBuffer, netCollect.metricBuffer)
	}

	go netCollect.SendConfig2CFC()
	go netCollect.SendMetri2CFC()

	return nil
}

//TODO:需要优化
func (this *NetCollect) Dial(Type string) {
	var (
		tempDelay time.Duration
		err       error
		session   *tcp.Session
	)
retry:
	switch Type {
	case "cfc":
		session, err = this.srv.Connect(GlobalConfig.CFC_ADDR, nil)
	case "repeater":
		session, err = this.srv.Connect(GlobalConfig.REPEATER_ADDR, nil)
	default:

	}
	if err != nil {
		lg.Error("%s: %s", Type, err.Error())
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
	switch Type {
	case "cfc":
		this.cfc = session
	case "repeater":
		this.repeater = session
	}
	for {
		switch Type {
		case "cfc":
			if this.cfc.IsClosed() {
				goto retry
			}
		case "repeater":
			if this.repeater.IsClosed() {
				goto retry
			}
		}
		time.Sleep(time.Second * 1)
	}
}

func (this *NetCollect) SendConfig2CFC() {
	for {
		if this.cfc.IsClosed() {
			goto sleep
		}
		for _, s := range this.switchs {
			if s.Hostname == "" || s.Hostname == "Unknown" {
				continue
			}
			h := &types.Host{
				ID:           s.ID,
				IP:           s.IP,
				Hostname:     s.Hostname,
				AgentVersion: AgentVersion,
			}
			//config
			this.cfc.Send(
				types.Pack(
					types.MESS_POST_HOST_CONFIG,
					h,
				))
			//heartbeat
			this.cfc.Send(
				types.Pack(
					types.MESS_POST_HOST_ALIVE,
					h,
				))
		}
	sleep:
		time.Sleep(time.Minute * 1)
	}
}

func (this *NetCollect) SendMetri2CFC() {
	time.Sleep(time.Second * time.Duration(GlobalConfig.COLLECT_INTERVAL))
	for {
		if this.cfc.IsClosed() {
			time.Sleep(time.Second * 30)
			continue
		}
		select {
		case tsd, ok := <-this.metricBuffer:
			if ok {
				this.cfc.Send(types.Pack(types.MESS_POST_METRIC, tsd))
				lg.Info("sender to cfc %s", tsd)
			}
		}
	}
}

func (this *NetCollect) SendTSD2Repeater() {
	go func() {
		for {
			if this.repeater.IsClosed() {
				time.Sleep(time.Second * 5)
				continue
			}
			select {
			case tsd, ok := <-this.tsdBuffer:
				if ok {
					this.repeater.Send(
						types.Pack(types.MESS_POST_TSD,
							tsd,
						))
					lg.Info("sender to repeater %s", tsd)
				}
			}
		}
	}()
}
