package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
)

//数据类型
const (
	HOSTCONFIG byte = iota //主机配置信息
	HOSTCONFIGRESP
	CLIENTVERSION //客户端版本号
	CLIENTVERSIONRESP
	SERIESDATA //监控数据

	GETDEVICES
	GETDEVICESRESP
	GETPORTS
	GETPORTSRESP
	UPDATEDEVICES //代理更新网络设备信息
	UPDATEPORTS   //代理更新主机端口状态

	PORTDISCOVERYR  //端口自动发现数据
	ASSESTDISCOVERY //资产自动发现

	HOSTHB
	GUARDHB
)

var PROTOCOLTYPE = []string{
	"HOSTCONFIG",
	"HOSTCONFIGRESP",
	"CLIENTVERSION",
	"CLIENTVERSIONRESP",
	"SERIESDATA",
	"GETDEVICES",
	"GETDEVICESRESP",
	"GETPORTS",
	"GETPORTSRESP",
	"UPDATEDEVICES",
	"UPDATEPORTS",
	"PORTDISCOVERYR",
	"ASSESTDISCOVERY",
	"HOSTHB",
	"GUARDHB",
}

//时间序列数据结构体
type SeriesData struct {
	UUID        string  `json:"uuid"`
	IP          string  `json:"ip"`
	Group       string  `json:"group"`
	TimeStamp   int64   `json:"timestamp"`
	ServiceName string  `json:"service_name"`
	Metric      string  `json:"metric"`
	Value       float32 `json:"value"`
	Flag        int     `json:"flag"`
}

type HostConfigReq struct {
	UUID      string `json:"uuid"`
	Proxy     string `json:"proxy"`
	IP        string `json:"ip"`
	IDRACAddr string `json:"idrac_addr"`
	HostName  string `json:"hostname"`
	Kernel    string `json:"kernel"`
	OS        string `json:"os"`
}

//如果版本不相等,返回版本信息,相等忽略
type GetVersionReq struct {
	Version string `json:"version"`
}

type PortDiscovery struct {
	UUID  string  `json:"uuid"`
	Ports []*Port `json:"ports"`
}

type HeartBeat struct {
	UUID string `json:"uuid"`
}

type ProxyReq struct {
	Proxy string `json:"proxy"`
}

type AssestPacket struct {
	UUID    string    `json:"uuid"`
	SN      string    `json:"sn"`
	Unit    string    `json:"unit"`
	Vender  string    `json:"vender"`
	Model   string    `json:"model"`
	Cpus    []*Cpu    `json:"cpus"`    //key is cpu.name
	Memorys []*Memory `json:"memorys"` //key is memory.locator
	Disks   []*Disk   `json:"disks"`   //key is disk.slot
	Nics    []*Nic    `json:"nics"`    //key is nic.name
}

func (this *SeriesData) Encode() ([]byte, error) {
	var buf bytes.Buffer
	res, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}
	head := make([]byte, 4)
	binary.BigEndian.PutUint32(head, uint32(1+len(res)))
	binary.Write(&buf, binary.BigEndian, head)
	binary.Write(&buf, binary.BigEndian, SERIESDATA)
	binary.Write(&buf, binary.BigEndian, res)
	return buf.Bytes(), nil
}
