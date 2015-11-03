package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"tcpserver"
	"time"
)

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

var PROTOCOLTYPE = map[byte]string{
	HOSTCONFIG:        "HOSTCONFIG",
	HOSTCONFIGRESP:    "HOSTCONFIGRESP",
	CLIENTVERSION:     "CLIENTVERSION",
	CLIENTVERSIONRESP: "CLIENTVERSIONRESP",
	SERIESDATA:        "SERIESDATA",
	GETDEVICES:        "GETDEVICES",
	GETDEVICESRESP:    "GETDEVICESRESP",
	GETPORTS:          "GETPORTS",
	GETPORTSRESP:      "GETPORTSRESP",
	UPDATEDEVICES:     "UPDATEDEVICES",
	UPDATEPORTS:       "UPDATEPORTS",
	PORTDISCOVERYR:    "PORTDISCOVERYR",
	ASSESTDISCOVERY:   "ASSESTDISCOVERY",
	HOSTHB:            "HOSTHB",
	GUARDHB:           "GUARDHB",
}

type HostConfigReq struct {
	UUID     string `json:"uuid"`
	Proxy    string `json:"proxy"`
	IP       string `json:"ip"`
	HostName string `json:"hostname"`
	Kernel   string `json:"kernel"`
	OS       string `json:"os"`
}

type ProxyReq struct {
	Proxy string `json:"proxy"`
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

func (this *ProxyReq) ToJson() []byte {
	res, _ := json.Marshal(this)
	return res
}

type handle struct {
}

func (this *handle) Connect(conn *tcpserver.Conn) {
	log.Info("%s new connection ", conn.RemoteAddr())
}

func (this *handle) Disconnect(conn *tcpserver.Conn) {
	log.Info("%s disconnect ", conn.RemoteAddr())
}

//数据包逻辑处理
func (this *handle) HandlerMessage(conn *tcpserver.Conn, data []byte) {
	defer func() {
		recover()
	}()
	log.Trace("receive data: %s|%s", PROTOCOLTYPE[data[0]], string(data[1:]))
	switch data[0] {
	case GETDEVICESRESP:
		devices := []*NetDevice{}
		err := json.Unmarshal(data[1:], &devices)
		if err != nil {
			log.Error("HandlerMessage() marshal netdevices error %s data:%s", err.Error(), string(data[1:]))
			return
		}
		for _, dev := range devices {
			if _, ok := devicemap[dev.ID]; ok {
				if dev.IpAddr != devicemap[dev.ID].IpAddr {
					log.Info("change device ipaddr old:%v new:%v", devicemap[dev.ID].IpAddr, dev.IpAddr)
					devicemap[dev.ID].IpAddr = dev.IpAddr
				}
				if dev.SnmpVersion != devicemap[dev.ID].SnmpVersion {
					log.Info("change device %v SnmpVersion old:%v new:%v", dev.IpAddr, devicemap[dev.ID].SnmpVersion, dev.SnmpVersion)
					devicemap[dev.ID].SnmpVersion = dev.SnmpVersion
				}
				if dev.SnmpCommunity != devicemap[dev.ID].SnmpCommunity {
					log.Info("change device %v SnmpCommunity old:%v new:%v", dev.IpAddr, devicemap[dev.ID].SnmpCommunity, dev.SnmpCommunity)
					devicemap[dev.ID].SnmpCommunity = dev.SnmpCommunity
				}
				if dev.SnmpPort != devicemap[dev.ID].SnmpPort {
					log.Info("change device %v SnmpPort old:%v new:%v", dev.IpAddr, devicemap[dev.ID].SnmpPort, dev.SnmpPort)
					devicemap[dev.ID].SnmpPort = dev.SnmpPort
				}
				if dev.UpdateInterval != devicemap[dev.ID].UpdateInterval {
					log.Info("change device %v UpdateInterval old:%v new:%v", dev.IpAddr, devicemap[dev.ID].UpdateInterval, dev.UpdateInterval)
					devicemap[dev.ID].UpdateInterval = dev.UpdateInterval
					devicemap[dev.ID].updateTicker = time.NewTicker(time.Second * time.Duration(dev.UpdateInterval))
				}
				if dev.CheckInterval != devicemap[dev.ID].CheckInterval {
					log.Info("change device %v CheckInterval old:%v new:%v", dev.IpAddr, devicemap[dev.ID].CheckInterval, dev.CheckInterval)
					devicemap[dev.ID].CheckInterval = dev.CheckInterval
					devicemap[dev.ID].checkTicker = time.NewTicker(time.Second * time.Duration(dev.CheckInterval))
				}
			} else {
				dev.stopChan = make(chan struct{})
				dev.updateTicker = time.NewTicker(time.Second * time.Duration(dev.UpdateInterval))
				dev.checkTicker = time.NewTicker(time.Second * time.Duration(dev.CheckInterval))
				go dev.Run()
				devicemap[dev.ID] = dev
			}
		}
		for _, dev := range devicemap {
			flag := -1
			for _, d := range devices {
				if dev.ID == d.ID {
					flag = 0
				}
			}
			if flag == -1 {
				dev.Stop()
				delete(devicemap, dev.ID)
			}
		}

	case GETPORTSRESP:
		ports := []*Port{}
		fmt.Println(string(data[1:]))
		err := json.Unmarshal(data[1:], &ports)
		if err != nil {
			log.Error("HandlerMessage() marshal ports error %s data:%s", err.Error(), string(data[1:]))
			return
		}
		nports := []*Port{}
		for _, port := range ports {
			if err := port.StatusCheck(); err != nil {
				if port.Status == 0 {
					port.Status = 1
					nports = append(nports, port)
				}
			} else {
				if port.Status == 1 {
					port.Status = 0
					nports = append(nports, port)
				}
			}
		}
		if len(nports) == 0 {
			return
		}
		resp, err := json.Marshal(nports)
		if err != nil {
			return
		}
		DataBuffer <- NewPacket(UPDATEPORTS, resp)

	case HOSTCONFIG:
		req := &HostConfigReq{}
		if err := json.Unmarshal(data[1:], &req); err != nil {
			log.Error("HandlerMessage() Parse HostConfigReq error(%s), data(%s)", err.Error(), string(data[1:]))
			return
		}
		if len(req.UUID) == 0 {
			return
		}
		if host, ok := hostmap[req.UUID]; ok {
			resp, err := json.Marshal(host)
			if err != nil {
				log.Error("marshal host error %s", err.Error())
				return
			}
			packet := NewPacket(HOSTCONFIGRESP, resp)
			conn.AsyncWriteData(packet)
		}
		buf := make([]byte, 4+len(data))
		binary.BigEndian.PutUint32(buf, uint32(len(data)))
		buf = append(buf[:4], data...)
		DataBuffer <- buf

	case HOSTCONFIGRESP:
		host := &Host{}
		if err := json.Unmarshal(data[1:], &host); err != nil {
			log.Error("HandlerMessage() Unmarshal() HOSTCONFIGRESP error %s data:%s", err.Error(), string(data[1:]))
			return
		}
		log.Info("update hostconfig %s:%s", host.Ip, host.UUID)
		hostmap[host.UUID] = host

	case CLIENTVERSIONRESP:
		var vp map[string]string
		if err := json.Unmarshal(data[1:], &vp); err != nil {
			fmt.Println("unmarshal error %s", err.Error())
			return
		}
		if v, ok := vp["version"]; ok {
			log.Info("update client version %s->%s", version, v)
			version = v
		}

	case CLIENTVERSION:
		if len(version) == 0 {
			return
		}
		vp := map[string]string{"version": version}
		resp, _ := json.Marshal(vp)
		packet := NewPacket(CLIENTVERSIONRESP, resp)
		if err := conn.AsyncWriteData(packet); err != nil {
			log.Error("HandlerMessage() AsyncWriteData() version packet error %s %s", err.Error(), conn.RemoteAddr())
			return
		}
	case SERIESDATA, PORTDISCOVERYR, HOSTHB, ASSESTDISCOVERY:
		buf := make([]byte, 4+len(data))
		binary.BigEndian.PutUint32(buf, uint32(len(data)))
		buf = append(buf[:4], data...)
		DataBuffer <- buf
	default:
		log.Warn("unkown protocol type:%d data %s", data[0], string(data[1:]))
	}
}

func NewPacket(t byte, data []byte) []byte {
	var buf bytes.Buffer
	head := make([]byte, 4)
	binary.BigEndian.PutUint32(head, uint32(len(data)+1))
	binary.Write(&buf, binary.BigEndian, head)
	binary.Write(&buf, binary.BigEndian, t)
	binary.Write(&buf, binary.BigEndian, data)
	return buf.Bytes()
}
