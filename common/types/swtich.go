package types

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"owl/common/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

var oids []string = []string{"ifHCOutOctets", "ifHCInOctets", "inErrors", "outErrors", "inDiscards", "outDiscards", "OperStatus"}

type Switch struct {
	ID              string   `json:"id"`
	IP              string   `json:"ip"`
	Hostname        string   `json:"hostname"`
	AgentVersion    string   `json:"agent_version"`
	CollectInterval int      `json:"collect_interval"`
	LegalPrefix     []string `json:"legal_prefix"`

	//采集数据成功更新
	LastUpdate time.Time             `json:"last_update"`
	Snmp       SnmpConfig            `json:"snmp"`
	Interfaces map[string]*Interface `json:"interfaces"`
}

type SnmpConfig struct {
	Port      int    `json:"port"`
	Version   string `json:"version"`
	Community string `json:"community"`
}

type Interface struct {
	Index       string    `json:"index"`
	Name        string    `json:"name"`
	OperStatus  uint64    `json:"oper_starus"`
	InBytes     [2]uint64 `json:"in_bytes"`
	OutBytes    [2]uint64 `json:"out_bytes"`
	InDiscards  [2]uint64 `json:"in_discards"`
	OutDiscards [2]uint64 `json:"out_discards"`
	InErrors    [2]uint64 `json:"in_errors"`
	OutErrors   [2]uint64 `json:"out_errors"`
	Speed       uint64    `json:"speed"`
}

func (this *Switch) walk(oid string) ([]byte, error) {
	return utils.RunCmdWithTimeout("snmpwalk", []string{"-v", this.Snmp.Version, "-c", this.Snmp.Community, this.IP, oid}, 5)
}

func (this *Switch) Do(buf1 chan<- *TimeSeriesData, buf2 chan<- *MetricConfig) {
retry:
	if err := this.BuildInterfaceIndex(); err != nil {
		if err == utils.ErrRunTimeout {
			fmt.Fprintln(os.Stderr, this.IP, "timeouts")
		}
		time.Sleep(time.Minute * 5)
		goto retry
	}
	this.initAllInterfaceData()
	this.CollectInterfaceName()
	this.getHostname()
	this.CollectIfaceSpeed()
	go this.loop(buf1)
	go this.postMetric(buf2)
}
func (this *Switch) postMetric(buffer chan<- *MetricConfig) {
	for {
		for _, i := range this.Interfaces {
			if !this.IsLegalPrefix(i.Name) {
				continue
			}
			buffer <- &MetricConfig{
				this.ID,
				TimeSeriesData{
					Metric:   "sw.if.InBytes",
					DataType: "COUNTER",
					Cycle:    this.CollectInterval,
					Tags:     map[string]string{"ifName": i.Name},
				},
			}
			buffer <- &MetricConfig{
				this.ID,
				TimeSeriesData{
					Metric:   "sw.if.OutBytes",
					DataType: "COUNTER",
					Cycle:    this.CollectInterval,
					Tags:     map[string]string{"ifName": i.Name},
				},
			}
			buffer <- &MetricConfig{
				this.ID,
				TimeSeriesData{
					Metric:   "sw.if.InErrors",
					DataType: "COUNTER",
					Cycle:    this.CollectInterval,
					Tags:     map[string]string{"ifName": i.Name},
				},
			}
			buffer <- &MetricConfig{
				this.ID,
				TimeSeriesData{
					Metric:   "sw.if.OutErrors",
					DataType: "COUNTER",
					Cycle:    this.CollectInterval,
					Tags:     map[string]string{"ifName": i.Name},
				},
			}
			buffer <- &MetricConfig{
				this.ID,
				TimeSeriesData{
					Metric:   "sw.if.InDiscards",
					DataType: "COUNTER",
					Cycle:    this.CollectInterval,
					Tags:     map[string]string{"ifName": i.Name},
				},
			}
			buffer <- &MetricConfig{
				this.ID,
				TimeSeriesData{
					Metric:   "sw.if.OutDiscards",
					DataType: "COUNTER",
					Cycle:    this.CollectInterval,
					Tags:     map[string]string{"ifName": i.Name},
				},
			}
			buffer <- &MetricConfig{
				this.ID,
				TimeSeriesData{
					Metric:   "sw.if.OutUsed.Percent",
					DataType: "GAUGE",
					Cycle:    this.CollectInterval,
					Tags:     map[string]string{"ifName": i.Name},
				},
			}
			buffer <- &MetricConfig{
				this.ID,
				TimeSeriesData{
					Metric:   "sw.if.InUsed.Percent",
					DataType: "GAUGE",
					Cycle:    this.CollectInterval,
					Tags:     map[string]string{"ifName": i.Name},
				},
			}
			buffer <- &MetricConfig{
				this.ID,
				TimeSeriesData{
					Metric:   "sw.if.OperStatus",
					DataType: "GAUGE",
					Cycle:    this.CollectInterval,
					Tags:     map[string]string{"ifName": i.Name},
				},
			}
		}
		buffer <- &MetricConfig{
			this.ID,
			TimeSeriesData{
				Metric:   "agent.alive",
				DataType: "GAUGE",
				Cycle:    this.CollectInterval,
			},
		}
		time.Sleep(time.Minute * 5)
	}
}

func (this *Switch) loop(buffer chan<- *TimeSeriesData) {
	for {
		ts := time.Now().Unix()
		this.CollectTraffic()
		interval := uint64(this.CollectInterval)
		flag := 0
		for _, i := range this.Interfaces {
			if !this.IsLegalPrefix(i.Name) || i.InBytes[0] == 12345678987654321 {
				continue
			}
			flag++
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.InBytes",
				DataType:  "COUNTER",
				Value:     float64((i.InBytes[1] - i.InBytes[0]) / interval),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "hostname": this.Hostname, "ifName": i.Name},
			}
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.OutBytes",
				DataType:  "COUNTER",
				Value:     float64((i.OutBytes[1] - i.OutBytes[0]) / interval),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "hostname": this.Hostname, "ifName": i.Name},
			}
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.InErrors",
				DataType:  "COUNTER",
				Value:     float64((i.InErrors[1] - i.InErrors[0]) / interval),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "hostname": this.Hostname, "ifName": i.Name},
			}
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.OutErrors",
				DataType:  "COUNTER",
				Value:     float64((i.OutErrors[1] - i.OutErrors[0]) / interval),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "hostname": this.Hostname, "ifName": i.Name},
			}
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.InDiscards",
				DataType:  "COUNTER",
				Value:     float64((i.InDiscards[1] - i.InDiscards[0]) / interval),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "hostname": this.Hostname, "ifName": i.Name},
			}
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.OutDiscards",
				DataType:  "COUNTER",
				Value:     float64((i.OutDiscards[1] - i.OutDiscards[0]) / interval),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "hostname": this.Hostname, "ifName": i.Name},
			}
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.InUsed.Percent",
				DataType:  "COUNTER",
				Value:     float64(((i.InBytes[1] - i.InBytes[0]) / interval) / i.Speed * 100),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "hostname": this.Hostname, "ifName": i.Name},
			}
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.OutUsed.Percent",
				DataType:  "COUNTER",
				Value:     float64(((i.OutBytes[1] - i.OutBytes[0]) / interval) / i.Speed * 100),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "hostname": this.Hostname, "ifName": i.Name},
			}
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.OperStatus",
				DataType:  "GAUGE",
				Value:     float64(i.OperStatus),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "hostname": this.Hostname, "ifName": i.Name},
			}
		}
		if flag != 0 {
			buffer <- &TimeSeriesData{
				Metric:    "agent.alive",
				DataType:  "GAUGE",
				Value:     1,
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "hostname": this.Hostname},
			}
		}
		time.Sleep(time.Second * time.Duration(this.CollectInterval))
	}
}

func (this *Switch) BuildInterfaceIndex() error {
	output, err := this.walk("ifIndex")
	if err != nil {
		return err
	}
	this.Interfaces = make(map[string]*Interface)
	buf := bytes.NewBuffer(output)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		fields := parseLine(line, " ")
		index := fields[len(fields)-1]
		this.Interfaces[index] = &Interface{
			Index: index,
		}
	}
	return nil
}

func (this *Switch) CollectInterfaceName() error {
	output, err := this.walk("ifName")
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(output)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		fields := parseLine(line, " ")
		indexField := parseLine(fields[0], ".")
		index := indexField[len(indexField)-1]
		name := fields[len(fields)-1]
		if iface, ok := this.Interfaces[index]; ok {
			iface.Name = name
		}
	}
	return nil
}

func (this *Switch) CollectIfaceSpeed() error {
	output, err := this.walk("ifSpeed")
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(output)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		fields := parseLine(line, " ")
		indexField := parseLine(fields[0], ".")
		index := indexField[len(indexField)-1]
		val := fields[len(fields)-1]
		speed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			continue
		}
		if iface, ok := this.Interfaces[index]; ok {
			iface.Speed = speed
		}
	}
	return nil
}

func (this *Switch) CollectTraffic() {
	var wg sync.WaitGroup
	for _, oid := range oids {
		wg.Add(1)
		go func(oid string) {
			this.CollectPerformanceData(oid)
			wg.Done()
		}(oid)
	}
	wg.Wait()
}
func parseLine(s string, sep string) []string {
	return strings.Split(strings.TrimSpace(s), sep)
}

func (this *Switch) CollectPerformanceData(oid string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "recover :%s", r)
		}
	}()
	output, err := this.walk(oid)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		this.initAllInterfaceData()
		return
	}
	buf := bytes.NewBuffer(output)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		fields := parseLine(line, " ")
		indexField := parseLine(fields[0], ".")
		index := indexField[len(indexField)-1]
		valField := fields[len(fields)-1]

		if _, ok := this.Interfaces[index]; !ok {
			continue
		}

		var val uint64
		if oid == "OperStatus" {
			if strings.Contains(valField, "up") {
				val = 1
			} else {
				val = 0
			}
		} else {
			val, err = strconv.ParseUint(valField, 10, 64)
			if err != nil {
				continue
			}
		}

		switch oid {
		case "ifHCInOctets":
			this.Interfaces[index].InBytes[0] = this.Interfaces[index].InBytes[1]
			this.Interfaces[index].InBytes[1] = val
		case "ifHCOutOctets":
			this.Interfaces[index].OutBytes[0] = this.Interfaces[index].OutBytes[1]
			this.Interfaces[index].OutBytes[1] = val
		case "inDiscards":
			this.Interfaces[index].InDiscards[0] = this.Interfaces[index].InDiscards[1]
			this.Interfaces[index].InDiscards[1] = val
		case "outDiscards":
			this.Interfaces[index].OutDiscards[0] = this.Interfaces[index].OutDiscards[1]
			this.Interfaces[index].OutDiscards[1] = val
		case "inErrors":
			this.Interfaces[index].InErrors[0] = this.Interfaces[index].InErrors[1]
			this.Interfaces[index].InErrors[1] = val
		case "outErrors":
			this.Interfaces[index].OutErrors[0] = this.Interfaces[index].OutErrors[1]
			this.Interfaces[index].OutErrors[1] = val
		case "OperStatus":
			this.Interfaces[index].OperStatus = val
		default:
		}
	}
}

func (this *Switch) IsLegalPrefix(name string) bool {
	for _, v := range this.LegalPrefix {
		if strings.HasPrefix(name, v) {
			return true
		}
	}
	return false
}

func (this *Switch) getHostname() {
	output, err := this.walk("1.3.6.1.2.1.1.5.0")
	if err != nil {
		this.Hostname = "Unknown"
	}
	fields := parseLine(string(output), " ")
	this.Hostname = fields[len(fields)-1]
}

func (this *Switch) initAllInterfaceData() {
	var init uint64 = 12345678987654321
	for _, i := range this.Interfaces {
		i.InBytes = [2]uint64{init, init}
		i.OutBytes = [2]uint64{init, init}
		i.InErrors = [2]uint64{init, init}
		i.OutErrors = [2]uint64{init, init}
		i.InDiscards = [2]uint64{init, init}
		i.OutDiscards = [2]uint64{init, init}
	}
}
