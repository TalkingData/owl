package main

import (
	"net"
	"os"
	"os/exec"
	"owl/agent/builtin"
	"owl/common/tcp"
	"owl/common/types"
	"owl/common/utils"
	"strings"
	"time"
)

const (
	SendHostConfigInterval    = 5  //minute
	GetHostPluginListInterval = 5  //minute
	RunBuiltinMetricCycle     = 30 //second
)

type Agent struct {
	Addr string
	*tcp.Server
	cfc        *tcp.Session
	repeater   *tcp.Session
	SendChan   chan types.TimeSeriesData
	tsdHistory map[string]float64
}

var (
	agent        Agent
	AgentVersion string = "0.1"
)

func InitAgent() error {
	s := tcp.NewServer(GlobalConfig.TCP_BIND, &cfcHandle{})
	agent = Agent{
		"",
		s,
		&tcp.Session{},
		&tcp.Session{},
		make(chan types.TimeSeriesData, GlobalConfig.BUFFER_SIZE),
		make(map[string]float64),
	}
	return agent.ListenAndServe()
}

//TODO(yingsong.wu):需要优化
func (this *Agent) Dial(Type string) {
	var (
		tempDelay time.Duration
		err       error
		session   *tcp.Session
	)
retry:
	switch Type {
	case "cfc":
		session, err = this.Connect(GlobalConfig.CFC_ADDR, nil)
	case "repeater":
		session, err = this.Connect(GlobalConfig.REPEATER_ADDR, nil)
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
		this.Addr = session.LocalAddr()
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

func (this *Agent) SN() string {
	res, err := exec.Command("sh", "-c", "/usr/sbin/dmidecode -t system |awk -F ':' '/Serial Number/ {print $NF}'").Output()
	if err != nil {
		return err.Error()
	}
	return strings.Trim(strings.Trim(string(res), " "), "\n")
}

func (this *Agent) ID() string {
	var str string
	a, _ := net.Interfaces()
	for _, i := range a {
		for _, prifex := range []string{"en", "eth", "em"} {
			if strings.HasPrefix(i.Name, prifex) && len(i.HardwareAddr.String()) > 0 {
				str += i.HardwareAddr.String()
				continue
			}
		}
	}
	return utils.Md5(str)[:16]
}

func (this *Agent) Hostname() string {
	n, err := os.Hostname()
	if err != nil {
		return err.Error()
	}
	return n
}

func (this *Agent) RIP() string {
	res, err := exec.Command("sh", "-c", "/usr/bin/ipmitool lan print | awk -F ':' '/IP Address.*[1-9]/ {print $2}'").Output()
	if err != nil {
		return err.Error()
	}
	return strings.Trim(strings.Trim(string(res), "\n"), " ")
}

func (this *Agent) IP() string {
	if len(this.Addr) > 0 {
		arr := strings.Split(this.Addr, ":")
		return arr[0]
	}
	return ""
}

func (this *Agent) SendConfig2CFC() {
	go func() {
		time.Sleep(time.Second)
		for {
			if this.cfc.IsClosed() {
				goto sleep
			}
			this.cfc.Send(
				types.Pack(types.MESS_POST_HOST_CONFIG,
					NewHostConfig(),
				))
		sleep:
			time.Sleep(time.Minute * SendHostConfigInterval)
		}
	}()
}

func (this *Agent) SendHostAlive2CFC() {
	go func() {
		for {
			if this.cfc.IsClosed() {
				time.Sleep(time.Second * 10)
				continue
			}
			this.cfc.Send(types.Pack(types.MESS_POST_HOST_ALIVE, NewHostConfig()))
			time.Sleep(time.Minute * 1)
		}
	}()
}

func (this *Agent) SendTSD2Repeater() {
	go func() {
		var (
			//tsd types.TimeSeriesData
			pk string
		)
		tags := map[string]string{"uuid": this.ID(), "hostname": this.Hostname()}
		for {
			if this.repeater.IsClosed() {
				time.Sleep(time.Second * 5)
				continue
			}
			select {
			case tsd, ok := <-this.SendChan:
				pk = tsd.PK()
				curr := tsd.Value
				history, ok := this.tsdHistory[pk]
				if ok {
					switch tsd.DataType {
					case "GAUGE", "gauge":
					case "COUNTER", "counter":
						tsd.Value = (tsd.Value - history) / float64(tsd.Cycle)
					case "DERIVE", "derive":
						tsd.Value = tsd.Value - history
						if tsd.Value < 0 {
							continue
						}
					default:
						lg.Error("%s DataType is illegal", tsd)
						continue
					}
					//如果运算后的值小于0，则取当前值
					if tsd.Value < 0 {
						tsd.Value = curr
					}
					if tsd.Tags == nil {
						tsd.Tags = tags
					} else {
						for k, v := range tags {
							tsd.Tags[k] = v
						}
					}
					this.repeater.Send(
						types.Pack(types.MESS_POST_TSD,
							&tsd,
						))
					lg.Info("sender to repeater %s", tsd)
				} else {
					this.cfc.Send(
						types.Pack(types.MESS_POST_METRIC,
							&types.MetricConfig{this.ID(), tsd},
						))
					lg.Info("sender to cfc %s", tsd)
				}
				this.tsdHistory[pk] = curr
			}
		}
	}()
}

func (this *Agent) GetPluginList() {
	go func() {
		time.Sleep(time.Second)
		for {
			if agent.cfc.IsClosed() {
				goto sleep
			}
			agent.cfc.Send(
				types.Pack(
					types.MESS_GET_HOST_PLUGIN_LIST,
					NewHostConfig()))
		sleep:
			time.Sleep(time.Minute * GetHostPluginListInterval)
		}
	}()
}

func (this *Agent) SendAgentAlive2Repeater() {
	go func() {
		for {
			this.SendChan <- types.TimeSeriesData{
				Metric:    "agent.alive",
				DataType:  "GAUGE",
				Value:     1,
				Cycle:     30,
				Timestamp: time.Now().Unix(),
			}
			time.Sleep(time.Second * 30)
		}
	}()
}

func (this *Agent) RunBuiltinMetric() {
	go func() {
		now := time.Now().Unix()
		diff := 60 - (now % 60)
		time.Sleep(time.Second * time.Duration(diff))
		for {
			time.Sleep(time.Second * RunBuiltinMetricCycle)
			go builtin.MemoryMetrics(RunBuiltinMetricCycle, this.SendChan)
			go builtin.SwapMetrics(RunBuiltinMetricCycle, this.SendChan)
			go builtin.LoadMetrics(RunBuiltinMetricCycle, this.SendChan)
			go builtin.NetMetrics(RunBuiltinMetricCycle, this.SendChan)
			go builtin.DiskMetrics(RunBuiltinMetricCycle, this.SendChan)
			go builtin.FdMetrics(RunBuiltinMetricCycle, this.SendChan)
			go builtin.CpuMetrics(RunBuiltinMetricCycle, this.SendChan)
			//metrics := []*types.TimeSeriesData{}
			//metrics = append(metrics, builtin.MemoryMetrics(RunBuiltinMetricCycle)...)
			//metrics = append(metrics, builtin.SwapMetrics(RunBuiltinMetricCycle)...)
			//metrics = append(metrics, builtin.LoadMetrics(RunBuiltinMetricCycle)...)
			//metrics = append(metrics, builtin.NetMetrics(RunBuiltinMetricCycle)...)
			//metrics = append(metrics, builtin.DiskMetrics(RunBuiltinMetricCycle)...)
			//metrics = append(metrics, builtin.FdMetrics(RunBuiltinMetricCycle)...)
			//metrics = append(metrics, builtin.CpuMetrics(RunBuiltinMetricCycle)...)
			//for _, v := range metrics {
			//	if v == nil {
			//		continue
			//	}
			//	this.SendChan <- *v
			//}
		}
	}()
}

func NewHostConfig() *types.Host {
	return &types.Host{
		ID:           agent.ID(),
		SN:           agent.SN(),
		IP:           agent.IP(),
		Hostname:     agent.Hostname(),
		AgentVersion: AgentVersion,
	}
}
