package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

//网络设备结构体定义
type NetDevice struct {
	ID     int    `json:"id"`
	UUID   string `json:"uuid"`
	IpAddr string `json:"ip"`    //设备IP地址
	Proxy  string `json:"proxy"` //设备代理地址,同server一个网段该地址为空

	SnmpVersion   string `json:"snmp_version"`   //snmp版本, 目前只支持2c
	SnmpCommunity string `json:"snmp_community"` //snmp团体名
	SnmpPort      int    `json:"snmp_port"`      //snmp端口号

	//lock     *sync.Mutex
	stopChan       chan struct{} //退出信号
	UpdateInterval int           `json:"update_interval"` //设备硬件配置(端口、速率、状态等)更新间隔
	CheckInterval  int           `json:"check_interval"`  //设备性能指标采集间隔
	updateTicker   *time.Ticker  //硬件配置更新定时器
	checkTicker    *time.Ticker  //性能指标采集定时器

	DeviceInterfaces map[string]*DeviceInterface `json:"interfaces"` //设备端口列表

	CustomOids map[string]*CustomOid `json:"custom_oids"` //用户自定义oid,key为oid号
}

type DeviceInterface struct {
	ID     int    `json:"id"`
	Index  int    `json:"index"`  //接口索引号
	Name   string `json:"name"`   //接口名称
	Status string `json:"status"` //接口状态
	Mac    string `json:"mac"`    //mac 地址
	Speed  string `json:"speed"`  //接口速率

	inOctets, outOctets         uint64 //接口当前时间周期接收/发送字节数
	prevInOctets, prevOutOctets uint64 //上一次采样数值

	inUcastPkts, outUcastPkts         uint64 //当前时间周期单播包数目
	prevInUcastPkts, prevOutUcastPkts uint64 //上次采样单播包数目

	inDiscards, outDiscards         uint64 //丢弃包数
	prevInDiscards, prevOutDiscards uint64 //上一次采样丢弃包数

	inErrors, outErrors         uint64 //当前时间周期接口错误包数
	prevInErrors, prevOutErrors uint64 //接口错误包数

	inUnknownProtos     uint64 //未知协议
	prevInUnknownProtos uint64 //上一次采样值
}

//自定义oid数据采集，一般用于内存，cpu
type CustomOid struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	OID  string `json:"oid"` //自定义oid号
	val  uint64 //值
}

//获取网络设备ip地址
func (this *NetDevice) GetIpAddr() string {
	return this.IpAddr
}

//获取网络设备代理地址
func (this *NetDevice) GetProxy() string {
	return this.Proxy
}

func (this *NetDevice) GetInterfaces() map[string]*DeviceInterface {
	if len(this.DeviceInterfaces) == 0 {
		return nil
	}
	return this.DeviceInterfaces
}

//执行硬件信息采集以及性能数据采集
func (this *NetDevice) Run() {
	nlog.Info("run() %s", this.IpAddr)
	this.collectHardwareData()
	for {
		select {
		case <-this.stopChan: //退出信号
			nlog.Info("%s receive stop single, exit.", this.IpAddr)
			this.checkTicker.Stop()
			this.updateTicker.Stop()
			return
		case <-this.checkTicker.C:
			nlog.Info("%s collectPerformanceData() ", this.IpAddr)
			this.collectPerformanceData()
		case <-this.updateTicker.C:
			nlog.Info("%s collectHardwareData() ", this.IpAddr)
			go this.collectHardwareData()
		}
	}
}

//发送停止信号，将停止性能数据采集和硬件信息
func (this *NetDevice) Stop() {

	this.stopChan <- struct{}{}
}

//通过给定oid获取snmp信息
func (this *NetDevice) getSnmp(oid string) (string, error) {

	cmd := exec.Command("/usr/bin/snmpwalk", "-v", this.SnmpVersion, "-c", this.SnmpCommunity, this.IpAddr, oid)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

//采集硬件信息
func (this *NetDevice) collectHardwareData() {
	this.collectIndex()
	if len(this.DeviceInterfaces) == 0 {
		return
	}
	go this.collectInterfaceName()
	go this.collectInterfaceMac()
	go this.collectInterfaceSpeed()
	go this.collectInterfaceStatus()
}

//采集性能数据并计算每秒平均值
func (this *NetDevice) collectPerformanceData() {
	defer func() {
		recover()
	}()
	if len(this.DeviceInterfaces) == 0 {
		nlog.Info("device %s interfaces is null, return", this.IpAddr)
		return
	}
	//ifHCOutOctets,ifHCInOctets 64位计数器
	//inOctets,outOctets   32位计数器

	oids := [9]string{"ifHCOutOctets", "ifHCInOctets",
		"inUcastPkts", "outUcastPkts",
		"inDiscards", "outDiscards",
		"inErrors", "outErrors",
		"inUnknownProtos",
	}
	var wg sync.WaitGroup
	for _, oid := range oids {
		wg.Add(1)
		go func(oid string) {
			this.getSNMPData(oid)
			wg.Done()
		}(oid)
	}
	for _, customoid := range this.CustomOids {
		wg.Add(1)
		go func(oid string) {
			this.getSNMPData(oid)
			wg.Done()
		}(customoid.OID)
	}
	wg.Wait()
	timeStamp := time.Now().Unix()
	rfcTimeStamp := time.Unix(timeStamp, 0).Format(time.RFC3339)
	interval := uint64(this.CheckInterval)
	for _, i := range this.DeviceInterfaces {
		if i.prevInOctets != 1234567890 {
			DataBuffer <- []byte(fmt.Sprintf("put %s.inOctets %v %v uuid=%s",
				i.Name, timeStamp, (i.inOctets-i.prevInOctets)/interval, this.UUID))
			if cfg.ENABLE_REDIS {
				RedisBuffer <- fmt.Sprintf(`{"ip":"%s","metric":"%s.inOctets","@timestamp":"%v","value":%v}`,
					this.IpAddr, i.Name, rfcTimeStamp, (i.inOctets-i.prevInOctets)/interval)
			}
		}
		if i.prevOutOctets != 1234567890 {
			DataBuffer <- []byte(fmt.Sprintf("put %s.outOctets %v %v uuid=%s",
				i.Name, timeStamp, (i.outOctets-i.prevOutOctets)/interval, this.UUID))
			if cfg.ENABLE_REDIS {
				RedisBuffer <- fmt.Sprintf(`{"ip":"%s","metric":"%s.outOctets","@timestamp":"%v","value":%v}`,
					this.IpAddr, i.Name, rfcTimeStamp, (i.outOctets-i.prevOutOctets)/interval)
			}
		}
		if i.prevInUcastPkts != 1234567890 {
			DataBuffer <- []byte(fmt.Sprintf("put %s.inUcastPkts %v %v uuid=%s",
				i.Name, timeStamp, (i.inUcastPkts-i.prevInUcastPkts)/interval, this.UUID))
			if cfg.ENABLE_REDIS {
				RedisBuffer <- fmt.Sprintf(`{"ip":"%s","metric":"%s.inUcastPkts","@timestamp":"%v","value":%v}`,
					this.IpAddr, i.Name, rfcTimeStamp, (i.inUcastPkts-i.prevInUcastPkts)/interval)
			}
		}
		if i.prevOutUcastPkts != 1234567890 {
			DataBuffer <- []byte(fmt.Sprintf("put %s.outUcastPkts %v %v uuid=%s",
				i.Name, timeStamp, (i.outUcastPkts-i.prevOutUcastPkts)/interval, this.UUID))
			if cfg.ENABLE_REDIS {
				RedisBuffer <- fmt.Sprintf(`{"ip":"%s","metric":"%s.outUcastPkts","@timestamp":"%v","value":%v}`,
					this.IpAddr, i.Name, rfcTimeStamp, (i.outUcastPkts-i.prevOutUcastPkts)/interval)
			}
		}
		if i.prevInDiscards != 1234567890 {
			DataBuffer <- []byte(fmt.Sprintf("put %s.inDiscards %v %v uuid=%s",
				i.Name, timeStamp, (i.inDiscards-i.prevInDiscards)/interval, this.UUID))
			if cfg.ENABLE_REDIS {
				RedisBuffer <- fmt.Sprintf(`{"ip":"%s","metric":"%s.inDiscards","@timestamp":"%v","value":%v}`,
					this.IpAddr, i.Name, rfcTimeStamp, (i.inDiscards-i.prevInDiscards)/interval)
			}
		}
		if i.prevOutDiscards != 1234567890 {
			DataBuffer <- []byte(fmt.Sprintf("put %s.outDiscards %v %v uuid=%s",
				i.Name, timeStamp, (i.outDiscards-i.prevOutDiscards)/interval, this.UUID))
			if cfg.ENABLE_REDIS {
				RedisBuffer <- fmt.Sprintf(`{"ip":"%s","metric":"%s.outDiscards","@timestamp":"%v","value":%v}`,
					this.IpAddr, i.Name, rfcTimeStamp, (i.outDiscards-i.prevOutDiscards)/interval)
			}
		}
		if i.prevInErrors != 1234567890 {
			DataBuffer <- []byte(fmt.Sprintf("put %s.inErrors %v %v uuid=%s",
				i.Name, timeStamp, (i.inErrors-i.prevInErrors)/interval, this.UUID))
			if cfg.ENABLE_REDIS {
				RedisBuffer <- fmt.Sprintf(`{"ip":"%s","metric":"%s.inErrors","@timestamp":"%v","value":%v}`,
					this.IpAddr, i.Name, rfcTimeStamp, (i.inErrors-i.prevInErrors)/interval)
			}
		}
		if i.prevOutErrors != 1234567890 {
			DataBuffer <- []byte(fmt.Sprintf("put %s.outErrors %v %v uuid=%s",
				i.Name, timeStamp, (i.outErrors-i.prevOutErrors)/interval, this.UUID))
			if cfg.ENABLE_REDIS {
				RedisBuffer <- fmt.Sprintf(`{"ip":"%s","metric":"%s.outErrors","@timestamp":"%v","value":%v}`,
					this.IpAddr, i.Name, rfcTimeStamp, (i.outErrors-i.prevOutErrors)/interval)
			}
		}
		if i.prevInUnknownProtos != 1234567890 {
			DataBuffer <- []byte(fmt.Sprintf("put %s.inUnknownProtos %v %v uuid=%s",
				i.Name, timeStamp, (i.inUnknownProtos-i.prevInUnknownProtos)/interval, this.UUID))
			if cfg.ENABLE_REDIS {
				RedisBuffer <- fmt.Sprintf(`{"ip":"%s","metric":"%s.inUnknownProtos","@timestamp":"%v","value":%v}`,
					this.IpAddr, i.Name, rfcTimeStamp, (i.inUnknownProtos-i.prevInUnknownProtos)/interval)
			}
		}

	}

	//自定义oid不支持计算每秒值
	for _, customoid := range this.CustomOids {
		DataBuffer <- []byte(fmt.Sprintf("put %s %v %v uuid=%s", customoid.Name, timeStamp, customoid.val, this.UUID))
		if cfg.ENABLE_REDIS {
			RedisBuffer <- fmt.Sprintf(`{"ip":"%s","metric":"%s.inUnknownProtos","@timestamp":"%v","value":%v}`,
				this.IpAddr, customoid.Name, rfcTimeStamp, customoid.val)
		}
	}

}

//通过snmp获取网络设备性能数据并循环获取每个接口的当前counter值
func (this *NetDevice) getSNMPData(oid string) {
	defer func() {
		recover()
	}()
	output, err := this.getSnmp(oid)
	if err != nil {
		nlog.Error("%s getSNMPData() %s error %s", this.IpAddr, oid, err)
		return
	}
	buf := bufio.NewReader(strings.NewReader(string(output)))
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}

		//获取接口索引
		index := getIndex(line)
		if len(index) == 0 {
			continue
		}
		//过滤 vlan,null,down 接口
		t := strings.Split(line, " ")
		res := strings.Trim(t[len(t)-1], "\n")
		val, err := strconv.ParseUint(res, 10, 64)
		if err != nil {
			continue
		}
		switch oid {
		//case "inOctets":
		case "ifHCInOctets":
			this.DeviceInterfaces[index].prevInOctets = this.DeviceInterfaces[index].inOctets
			this.DeviceInterfaces[index].inOctets = val
		//case "outOctets":
		case "ifHCOutOctets":
			this.DeviceInterfaces[index].prevOutOctets = this.DeviceInterfaces[index].outOctets
			this.DeviceInterfaces[index].outOctets = val
		case "inUcastPkts":
			this.DeviceInterfaces[index].prevInUcastPkts = this.DeviceInterfaces[index].inUcastPkts
			this.DeviceInterfaces[index].inUcastPkts = val
		case "outUcastPkts":
			this.DeviceInterfaces[index].prevOutUcastPkts = this.DeviceInterfaces[index].outUcastPkts
			this.DeviceInterfaces[index].outUcastPkts = val
		case "inDiscards":
			this.DeviceInterfaces[index].prevInDiscards = this.DeviceInterfaces[index].inDiscards
			this.DeviceInterfaces[index].inDiscards = val
		case "outDiscards":
			this.DeviceInterfaces[index].prevOutDiscards = this.DeviceInterfaces[index].outDiscards
			this.DeviceInterfaces[index].outDiscards = val
		case "inErrors":
			this.DeviceInterfaces[index].prevInErrors = this.DeviceInterfaces[index].inErrors
			this.DeviceInterfaces[index].inErrors = val
		case "outErrors":
			this.DeviceInterfaces[index].prevOutErrors = this.DeviceInterfaces[index].outErrors
			this.DeviceInterfaces[index].outErrors = val
		case "inUnknownProtos":
			this.DeviceInterfaces[index].prevInUnknownProtos = this.DeviceInterfaces[index].inUnknownProtos
			this.DeviceInterfaces[index].inUnknownProtos = val
		default:
			if _, ok := this.CustomOids[oid]; ok {
				this.CustomOids[oid].val = val
			}
		}
	}
}

//获取所有接口索引,并生成接口对象,需要在其他采集函数前调用
func (this *NetDevice) collectIndex() {
	output, err := this.getSnmp("ifIndex")
	if err != nil {
		nlog.Error("%s getSnmp('ifIndex') error %s  ", this.IpAddr, err)
		return
	}
	buf := bufio.NewReader(strings.NewReader(string(output)))
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			return
		}
		index := getIndex(line)
		if len(index) == 0 {
			continue
		}
		if _, ok := this.DeviceInterfaces[index]; !ok {
			i, err := strconv.Atoi(index)
			if err != nil {
				continue
			}
			this.DeviceInterfaces[index] = &DeviceInterface{
				Index:           i,
				inOctets:        1234567890,
				outOctets:       1234567890,
				inUcastPkts:     1234567890,
				outUcastPkts:    1234567890,
				inDiscards:      1234567890,
				outDiscards:     1234567890,
				inErrors:        1234567890,
				outErrors:       1234567890,
				inUnknownProtos: 1234567890,
			}
		}
	}
}

//获取交换机接口名称，需要在调用collectIndex方法后调用
func (this *NetDevice) collectInterfaceName() {
	output, err := this.getSnmp("ifDesc")
	if err != nil {
		nlog.Error("%s getSnmp('ifDesc') error %s  ", this.IpAddr, err)
		return
	}
	buf := bufio.NewReader(strings.NewReader(string(output)))
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			return
		}
		index := getIndex(line)
		if len(index) == 0 {
			continue
		}
		data := strings.Split(line, ":")
		name := strings.Trim(strings.Trim(data[len(data)-1], " "), "\n")
		if this.DeviceInterfaces[index].Name != name {
			this.DeviceInterfaces[index].Name = name
		}
	}
}

//获取网络设备接口速率
func (this *NetDevice) collectInterfaceSpeed() {
	output, err := this.getSnmp("ifSpeed")
	if err != nil {
		nlog.Error("%s getSnmp('ifSpeed') error  %s  ", this.IpAddr, err)
		return
	}
	buf := bufio.NewReader(strings.NewReader(string(output)))
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			return
		}
		index := getIndex(line)
		if len(index) == 0 {
			continue
		}
		data := strings.Split(line, ":")
		speed := strings.Trim(strings.Trim(data[len(data)-1], " "), "\n")
		if this.DeviceInterfaces[index].Speed != speed {
			this.DeviceInterfaces[index].Speed = speed
		}
	}
}

//获取网络设备接口的mac地址
func (this *NetDevice) collectInterfaceMac() {
	output, err := this.getSnmp("ifPhysAddress")
	if err != nil {
		nlog.Error("%s getSnmp('ifPhysAddress') error %s  ", this.IpAddr, err)
		return
	}
	buf := bufio.NewReader(strings.NewReader(string(output)))
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			return
		}
		index := getIndex(line)
		if len(index) == 0 {
			continue
		}
		data := strings.Split(line, " ")
		mac := strings.Trim(strings.Trim(data[len(data)-1], " "), "\n")
		if this.DeviceInterfaces[index].Mac != mac {
			this.DeviceInterfaces[index].Mac = mac
		}
	}
}

//通过snmp获取网络设备所有接口状态并赋值给对应的interface
//IF-MIB::ifOperStatus.84 = INTEGER: up(1)		取得 up(1)
func (this *NetDevice) collectInterfaceStatus() {
	output, err := this.getSnmp("ifOperStatus")
	if err != nil {
		nlog.Error("%s getSnmp('ifOperStatus') error %s  ", this.IpAddr, err)
		return
	}
	buf := bufio.NewReader(strings.NewReader(string(output)))
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			return
		}
		index := getIndex(line)
		if len(index) == 0 {
			continue
		}
		data := strings.Split(line, ":")
		status := strings.Trim(strings.Trim(data[len(data)-1], " "), "\n")
		if this.DeviceInterfaces[index].Status != status {
			this.DeviceInterfaces[index].Status = status
		}
	}
}

//从snmp返回单行数据中截取index信息
//IF-MIB::ifIndex.34 = INTEGER: 34 将返回34
func getIndex(s string) string {
	data := strings.Split(strings.Split(strings.Trim(s, "\n"), " ")[0], ".")
	return data[len(data)-1]
}

//状态维护

func Contains(devices []*NetDevice, device *NetDevice) bool {
	for _, dev := range devices {
		if dev.ID == device.ID {
			return true
		}
	}
	return false
}

func HandleDevicesLoop(devmap map[int]*NetDevice) {
	for {
		time.Sleep(time.Minute * 5)
		//获取数据库中网络设备信息
		nlog.Info("start HandleDevicesLoop()")
		devs, err := mysql.GetDeviceByProxy("")
		if err != nil {
			continue
		}
		for _, dev := range devs {
			if _, ok := devmap[dev.ID]; !ok {
				nlog.Info("new device %s, add to devmap", dev.IpAddr)
				go dev.Run()
				devmap[dev.ID] = dev
			} else {
				if dev.IpAddr != devmap[dev.ID].IpAddr {
					nlog.Info("change device ipaddr old:%v new:%v", devmap[dev.ID].IpAddr, dev.IpAddr)
					devmap[dev.ID].IpAddr = dev.IpAddr
				}
				if dev.SnmpVersion != devmap[dev.ID].SnmpVersion {
					nlog.Info("change device %v SnmpVersion old:%v new:%v", dev.IpAddr, devmap[dev.ID].SnmpVersion, dev.SnmpVersion)
					devmap[dev.ID].SnmpVersion = dev.SnmpVersion
				}
				if dev.SnmpCommunity != devmap[dev.ID].SnmpCommunity {
					nlog.Info("change device %v SnmpCommunity old:%v new:%v", dev.IpAddr, devmap[dev.ID].SnmpCommunity, dev.SnmpCommunity)
					devmap[dev.ID].SnmpCommunity = dev.SnmpCommunity
				}
				if dev.SnmpPort != devmap[dev.ID].SnmpPort {
					nlog.Info("change device %v SnmpPort old:%v new:%v", dev.IpAddr, devmap[dev.ID].SnmpPort, dev.SnmpPort)
					devmap[dev.ID].SnmpPort = dev.SnmpPort
				}
				if dev.UpdateInterval != devmap[dev.ID].UpdateInterval {
					nlog.Info("change device %v UpdateInterval old:%v new:%v", dev.IpAddr, devmap[dev.ID].UpdateInterval, dev.UpdateInterval)
					devmap[dev.ID].UpdateInterval = dev.UpdateInterval
					devmap[dev.ID].updateTicker = time.NewTicker(time.Second * time.Duration(dev.UpdateInterval))
				}
				if dev.CheckInterval != devmap[dev.ID].CheckInterval {
					nlog.Info("change device %v CheckInterval old:%v new:%v", dev.IpAddr, devmap[dev.ID].CheckInterval, dev.CheckInterval)
					devmap[dev.ID].CheckInterval = dev.CheckInterval
					devmap[dev.ID].checkTicker = time.NewTicker(time.Second * time.Duration(dev.CheckInterval))
				}
				for oid, customoid := range dev.CustomOids {
					if _, ok := devmap[dev.ID].CustomOids[oid]; !ok {
						devmap[dev.ID].CustomOids[oid] = customoid
					}
				}
				for oid, _ := range devmap[dev.ID].CustomOids {
					if _, ok := dev.CustomOids[oid]; !ok {
						delete(devmap[dev.ID].CustomOids, oid)
					}
				}
			}
		}
		for _, dev := range devmap {
			flag := -1
			for _, odev := range devs {
				if dev.ID == odev.ID {
					flag = 0
				}
			}
			if flag == -1 {
				nlog.Info("device %s not in the db, remove from devmap.", dev.IpAddr)
				dev.Stop()
				delete(devmap, dev.ID)
			}
		}
		for _, dev := range devmap {
			for _, i := range dev.DeviceInterfaces {
				name := strings.ToLower(i.Name)
				if strings.Contains(name, "vlan") ||
					strings.Contains(name, "aux") ||
					strings.Contains(name, "loop") ||
					strings.Contains(name, "null") {
					continue
				}
				if mysql.InterfaceExist(i, dev) {
					mysql.UpdateInterfaces(dev, i)
				} else {
					nlog.Info("device %s have a new interface %s, insert to mysql", dev.IpAddr, i.Name)
					if err := mysql.CreateInterface(i, dev); err != nil {
						nlog.Info("device %s interface %s insert to mysql error %s", dev.IpAddr, i.Name, err)
						continue
					} else {
						nlog.Info("device %s interface %s insert to mysql done", dev.IpAddr, i.Name)
					}
				}
			}
		}

	}
}
