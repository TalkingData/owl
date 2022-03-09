package types

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wuyingsong/utils"
)

var oids []string = []string{
	"ifHCOutOctets",
	"ifHCInOctets",
	"inErrors",
	"outErrors",
	"inDiscards",
	"outDiscards",
	"OperStatus",
}

const (
	H3C_CPU_OID    = ".1.3.6.1.4.1.25506.2.6.1.1.1.1.6"
	H3C_MEM_OID    = ".1.3.6.1.4.1.25506.2.6.1.1.1.1.8"
	CISCO_CPU_OID  = ".1.3.6.1.4.1.9.9.109.1.1.1.1.7"
	CISCO_MEM_OID  = ".1.3.6.1.4.1.9.9.109.1.1.1.1.12"
	HUAWEI_MEM_OID = ".1.3.6.1.4.1.2011.5.25.31.1.1.1.1.7"
	HUAWEI_CPU_OID = ".1.3.6.1.4.1.2011.5.25.31.1.1.1.1.5"
)

const (
	defaultCounterV = 12345678987654321
)

type Counter struct {
	Last, Prev uint64
}

func newCounter() Counter {
	return Counter{
		defaultCounterV,
		defaultCounterV,
	}
}

func (c *Counter) getValue(i uint64) float64 {
	return float64((c.Last - c.Prev)) / float64(i)
}

type Switch struct {
	ID              string   `json:"id"`
	IP              string   `json:"ip"`
	Hostname        string   `json:"hostname"`
	Vendor          string   `json:"vendor"`
	AgentVersion    string   `json:"agent_version"`
	CollectInterval int      `json:"collect_interval"`
	LegalPrefix     []string `json:"legal_prefix"`

	//采集数据成功更新
	LastUpdate time.Time             `json:"last_update"`
	Snmp       SnmpConfig            `json:"snmp"`
	Interfaces map[string]*Interface `json:"interfaces"`
	Cpu        map[string]uint64
	Mem        map[string]uint64
}

type SnmpConfig struct {
	Port      int    `json:"port"`
	Version   string `json:"version"`
	Community string `json:"community"`
	Timeout   int    `json:"timeout"`
}

type Interface struct {
	Index       string  `json:"index"`
	Name        string  `json:"name"`
	OperStatus  uint64  `json:"oper_starus"`
	InBytes     Counter `json:"in_bytes"`
	OutBytes    Counter `json:"out_bytes"`
	InDiscards  Counter `json:"in_discards"`
	OutDiscards Counter `json:"out_discards"`
	InErrors    Counter `json:"in_errors"`
	OutErrors   Counter `json:"out_errors"`
	Speed       uint64  `json:"speed"`
}

func (this *Switch) walk(oid string) ([]byte, error) {
	return utils.RunCmdWithTimeout("snmpwalk", []string{"-v", this.Snmp.Version, "-c", this.Snmp.Community, this.IP, oid}, this.Snmp.Timeout)
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
	this.Cpu = make(map[string]uint64)
	this.Mem = make(map[string]uint64)
	this.initAllInterfaceData()
	this.CollectInterfaceName()
	this.GetHostname()
	this.GetVendor()
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
		for idx := range this.Cpu {
			buffer <- &MetricConfig{
				this.ID,
				TimeSeriesData{
					Metric:   "sw.cpu.used.percent",
					DataType: "GAUGE",
					Cycle:    this.CollectInterval,
					Tags:     map[string]string{"index": idx},
				},
			}
		}
		for idx := range this.Mem {
			buffer <- &MetricConfig{
				this.ID,
				TimeSeriesData{
					Metric:   "sw.mem.used.percent",
					DataType: "GAUGE",
					Cycle:    this.CollectInterval,
					Tags:     map[string]string{"index": idx},
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
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "recover ", this.IP, r)
		}
	}()
	for {
		ts := time.Now().Unix()
		ts = ts - (ts % 60)
		this.CollectTraffic()
		interval := uint64(this.CollectInterval)
		flag := 0
		for _, i := range this.Interfaces {
			if !this.IsLegalPrefix(i.Name) || i.InBytes.Prev == defaultCounterV {
				continue
			}
			flag++
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.InBytes",
				DataType:  "COUNTER",
				Value:     i.InBytes.getValue(interval),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "host": this.Hostname, "ifName": i.Name},
			}
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.OutBytes",
				DataType:  "COUNTER",
				Value:     i.OutBytes.getValue(interval),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "host": this.Hostname, "ifName": i.Name},
			}
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.InErrors",
				DataType:  "COUNTER",
				Value:     i.InErrors.getValue(interval),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "host": this.Hostname, "ifName": i.Name},
			}
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.OutErrors",
				DataType:  "COUNTER",
				Value:     i.OutErrors.getValue(interval),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "host": this.Hostname, "ifName": i.Name},
			}
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.InDiscards",
				DataType:  "COUNTER",
				Value:     i.InDiscards.getValue(interval),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "host": this.Hostname, "ifName": i.Name},
			}
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.OutDiscards",
				DataType:  "COUNTER",
				Value:     i.OutDiscards.getValue(interval),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "host": this.Hostname, "ifName": i.Name},
			}
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.InUsed.Percent",
				DataType:  "COUNTER",
				Value:     i.InBytes.getValue(interval) / float64(i.Speed) * 100,
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "host": this.Hostname, "ifName": i.Name},
			}
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.OutUsed.Percent",
				DataType:  "COUNTER",
				Value:     i.OutBytes.getValue(interval) / float64(i.Speed) * 100,
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "host": this.Hostname, "ifName": i.Name},
			}
			buffer <- &TimeSeriesData{
				Metric:    "sw.if.OperStatus",
				DataType:  "GAUGE",
				Value:     float64(i.OperStatus),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "host": this.Hostname, "ifName": i.Name},
			}
		}
		if flag != 0 {
			buffer <- &TimeSeriesData{
				Metric:    "agent.alive",
				DataType:  "GAUGE",
				Value:     1,
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "host": this.Hostname},
			}
		}
		for idx, val := range this.Mem {
			buffer <- &TimeSeriesData{
				Metric:    "sw.mem.used.percent",
				DataType:  "GAUGE",
				Value:     float64(val),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "host": this.Hostname, "index": idx},
			}
		}
		for idx, val := range this.Cpu {
			buffer <- &TimeSeriesData{
				Metric:    "sw.cpu.used.percent",
				DataType:  "GAUGE",
				Value:     float64(val),
				Timestamp: ts,
				Cycle:     this.CollectInterval,
				Tags:      map[string]string{"uuid": this.ID, "ip": this.IP, "host": this.Hostname, "index": idx},
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
func (this *Switch) collectCpu() error {
	if this.Cpu == nil {
		this.Cpu = make(map[string]uint64)
	}
	var cpuOID string
	switch this.Vendor {
	case "cisco":
		cpuOID = CISCO_CPU_OID
	case "h3c":
		cpuOID = H3C_CPU_OID
	case "huawei":
		cpuOID = HUAWEI_CPU_OID
	default:
		return fmt.Errorf("unsupport vendor")
	}
	output, err := this.walk(cpuOID)
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
		val, err := strconv.ParseUint(fields[len(fields)-1], 10, 64)
		if err != nil || val == 0 {
			continue
		}
		this.Cpu[index] = val
	}
	return nil
}

func (this *Switch) collectMem() error {
	if this.Mem == nil {
		this.Mem = make(map[string]uint64)
	}
	var memOID string
	switch this.Vendor {
	case "cisco":
		memOID = CISCO_MEM_OID
	case "h3c":
		memOID = H3C_MEM_OID
	case "huawei":
		memOID = HUAWEI_MEM_OID
	default:
		return fmt.Errorf("unsupport vendor")
	}
	output, err := this.walk(memOID)
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
		val, err := strconv.ParseUint(fields[len(fields)-1], 10, 64)
		if err != nil || val == 0 {
			continue
		}
		this.Mem[index] = val
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
	wg.Add(2)
	go func() {
		this.collectCpu()
		wg.Done()
	}()
	go func() {
		this.collectMem()
		wg.Done()
	}()

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
		fmt.Fprintln(os.Stderr, this.IP, err.Error())
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
		iface := this.Interfaces[index]
		switch oid {
		case "ifHCInOctets":
			iface.InBytes.Prev = iface.InBytes.Last
			iface.InBytes.Last = val
		case "ifHCOutOctets":
			iface.OutBytes.Prev = iface.OutBytes.Last
			iface.OutBytes.Last = val
		case "inDiscards":
			iface.InDiscards.Prev = iface.InDiscards.Last
			iface.InDiscards.Last = val
		case "outDiscards":
			iface.OutDiscards.Prev = iface.OutDiscards.Last
			iface.OutDiscards.Last = val
		case "inErrors":
			iface.InErrors.Prev = iface.InErrors.Last
			iface.InErrors.Last = val
		case "outErrors":
			iface.OutErrors.Prev = iface.OutErrors.Last
			iface.OutErrors.Last = val
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

func (this *Switch) GetHostname() {
	output, err := this.walk("1.3.6.1.2.1.1.5.0")
	if err != nil {
		this.Hostname = "Unknown"
	}
	fields := parseLine(string(output), " ")
	this.Hostname = fields[len(fields)-1]
}

func (this *Switch) GetVendor() {
	output, err := this.walk("1.3.6.1.2.1.1.1.0")
	if err != nil {
		fmt.Println("get sysdesc error ", err.Error())
		return
	}
	sysDescLowerString := strings.ToLower(string(output))
	if strings.Contains(sysDescLowerString, "cisco") {
		this.Vendor = "cisco"
	} else if strings.Contains(sysDescLowerString, "h3c") {
		this.Vendor = "h3c"
	} else if strings.Contains(sysDescLowerString, "huawei") {
		this.Vendor = "huawei"
	} else {
		this.Vendor = "Unknown"
	}
}

func (this *Switch) IsInit() bool {
	if len(this.Hostname) == 0 || this.Hostname == "Unknown" {
		return false
	}
	return true
}

func (this *Switch) initAllInterfaceData() {
	for _, i := range this.Interfaces {
		i.InBytes = newCounter()
		i.OutBytes = newCounter()
		i.InErrors = newCounter()
		i.OutErrors = newCounter()
		i.InDiscards = newCounter()
		i.OutDiscards = newCounter()
	}

}
