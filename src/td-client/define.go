package main

import "time"

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

type Host struct {
	ID       int        `json:"id"`
	UUID     string     `json:"uuid"`
	AssestId int        `json:"assest_id"`
	Ip       string     `json:"ip"`
	Group    string     `json:"group"`
	Proxy    string     `json:"proxy"`
	Services []*Service `json:"services"`
	Status   int
}

type HeartBeat struct {
	UUID string `json:"uuid"`
}

type Service struct {
	Id           int              `json:"id"`
	Name         string           `json:"name"`
	Plugin       string           `json:"plugin"`
	Args         string           `json:"args"`
	ExecInterval int              `json:"exec_interval"`
	Items        map[string]*Item `json:"items"`
	LastCheck    time.Time
}

type Item struct {
	Key string `json:"key"`
	//value [2]float64 `json:"-"`
	DT string `json:"data_type"` //GAUGE|COUNTER|DERIVE
}

type Port struct {
	Id       int    `json:"id"`
	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	ProcName string `json:"proc_name"`
	Status   int    `json:"status"`
}

type Cpu struct {
	ID       int    `json:"id"`
	AssestId int    `json:"assest_id"`
	Name     string `json:"name"`
	SN       string `json:"sn"`
	Model    string `json:"model"`
	Vender   string `json:"vender"`
}

type Memory struct {
	ID       int    `json:"id"`
	AssestId int    `json:"assest_id"`
	SN       string `json:"sn"`
	Size     string `json:"size"`
	Speed    string `json:"speed"`
	Locator  string `json:"locator"`
	Vender   string `json:"vender"`
}

type Disk struct {
	ID             int    `json:"id"`
	AssestId       int    `json:"assest_id"`
	SN             string `json:"sn"`
	Vender         string `json:"vender"`
	ProductId      string `json:"product_id"`
	Slot           string `json:"slot"`
	ProductionDate string `json:"production_date"`
	Capacity       string `json:"capacity"`
	Speed          string `json:"speed"`
	Bus            string `json:"bus"`
	Media          string `json:"media"`
}

type Nic struct {
	ID          int    `json:"id"`
	AssestId    int    `json:"assest_id"`
	Name        string `json:"name"`
	Mac         string `json:"mac"`
	Vender      string `json:"vender"`
	Description string `json:"description"`
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

type HostConfigReq struct {
	UUID      string `json:"uuid"`
	Proxy     string `json:"proxy"`
	IP        string `json:"ip"`
	HostName  string `json:"hostname"`
	Kernel    string `json:"kernel"`
	OS        string `json:"os"`
	IDRACAddr string `json:"idrac_addr"`
}

type PortDiscovery struct {
	UUID  string  `json:"uuid"`
	Ports []*Port `json:"ports"`
}

//时间序列数据结构体
type SeriesData struct {
	UUID        string  `json:"uuid"`
	IP          string  `json:"ip"`
	Group       string  `json:"group"`
	TimeStamp   int64   `json:"timestamp"`
	ServiceName string  `json:"service_name"`
	Metric      string  `json:"metric"`
	Value       float64 `json:"value"`
	Flag        int     `json:"flag"`
}
