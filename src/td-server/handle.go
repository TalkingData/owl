package main

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"tcpserver"
	"time"
)

type handle struct {
}

func (this *handle) Connect(conn *tcpserver.Conn) {
	slog.Info("%s new connection ", conn.RemoteAddr())
}

func (this *handle) Disconnect(conn *tcpserver.Conn) {
	slog.Info("%s disconnect ", conn.RemoteAddr())
}

//数据包逻辑处理
func (this *handle) HandlerMessage(conn *tcpserver.Conn, data []byte) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered in HandlerMessage", r)
		}
	}()
	slog.Debug("receive data: %s|%s", PROTOCOLTYPE[data[0]], string(data[1:]))
	switch data[0] {
	case HOSTCONFIG: //get host config
		req := &HostConfigReq{}
		if err := json.Unmarshal(data[1:], &req); err != nil {
			slog.Error("HandlerMessage() Parse HostConfigReq error(%s), data(%s)", err.Error(), string(data[1:]))
			return
		}
		if len(req.UUID) == 0 {
			return
		}
		host, err := mysql.GetHostConfigByUUID(req.UUID)
		if err != nil {
			if err == sql.ErrNoRows {
				proxy := conn.GetLocalIp()
				if req.Proxy == proxy {
					req.Proxy = ""
				}
				err = mysql.CreateHost(req)
				if err != nil {
					slog.Error("add host failed uuid(%s) ip(%s) proxy(%s) error(%s) ", req.UUID, req.IP, req.Proxy, err.Error())
					return
				} else {
					slog.Info("add host uuid(%s) proxy(%s) ip(%s) host(%v)", req.UUID, req.Proxy, req.IP, host)
				}
			}
			slog.Error("get host config error uuid(%s) ip(%s) proxy(%s) error(%s)", req.UUID, req.IP, req.Proxy, err.Error())
			return
		}

		if err := mysql.UpdateHostInfo(req); err != nil {
			slog.Error("update host info error %s", req.IP)
		}
		host.Ip = req.IP
		resp, err := json.Marshal(host)
		if err != nil {
			slog.Error("marshal host error %s", err.Error())
			return
		}
		packet := NewPacket(HOSTCONFIGRESP, resp)
		if err := conn.AsyncWriteData(packet); err != nil {
			slog.Error("HandlerMessage() AsyncWriteData() error(%s)", err.Error())
		}

	case CLIENTVERSION: //获取当前最新版本
		version := mysql.GetVersion()
		if len(version) == 0 {
			return
		}
		vp := map[string]string{"version": version}
		resp, _ := json.Marshal(vp)
		packet := NewPacket(CLIENTVERSIONRESP, resp)
		if err := conn.AsyncWriteData(packet); err != nil {
			slog.Error("AsyncWriteData error %s", err.Error())
		}
	case SERIESDATA:
		sdata := &SeriesData{}
		if err := json.Unmarshal(data[1:], &sdata); err != nil {
			slog.Error("parse monior data error(%s) data(%s)", err.Error(), string(data))
			return
		}
		//添加指标
		if sdata.Flag == 1 {
			host, err := mysql.GetHostConfigByUUID(sdata.UUID)
			if err != nil {
				return
			}
			for _, s := range host.Services {
				if s.Name == sdata.ServiceName {
					item := &Item{Key: sdata.Metric, DT: "GAUGE"}
					if err := mysql.CreateItemByService(s, item); err != nil {
						slog.Error("create item error(%s)", err.Error())
					} else {
						slog.Info("create item done service(%v) item(%v)", s.Name, sdata)
					}
				}
			}
		}
		DataBuffer <- []byte(fmt.Sprintf("put %s.%s %v %v uuid=%s ip=%s", sdata.ServiceName, sdata.Metric, sdata.TimeStamp, sdata.Value, sdata.UUID, sdata.IP))
		if cfg.ENABLE_REDIS {
			RedisBuffer <- fmt.Sprintf(`{"ip":"%s","metric":"%s.%s","@timestamp":"%v","value":%v}`,
				sdata.IP, sdata.ServiceName, sdata.Metric, time.Unix(sdata.TimeStamp, 0).Format(time.RFC3339), sdata.Value)
		}
	case PORTDISCOVERYR:
		req := &PortDiscovery{}
		if err := json.Unmarshal(data[1:], &req); err != nil {
			slog.Error("HandlerMessage() Parse PortDiscovery error(%s), data(%s)", err.Error(), string(data[1:]))
			return
		}
		host, err := mysql.GetHostByUUID(req.UUID)
		if err != nil {
			return
		}
		for _, port := range req.Ports {
			if !mysql.PortIsExists(host.ID, port.Port) {
				if err := mysql.CreatePort(host.ID, port); err != nil {
					slog.Error("create host port error(%s), host(%#v), port(%#v)", err.Error(), host, port)
					continue
				}
				slog.Info("create port, host(%s|%s) port(%v)", host.UUID, host.Ip, port)
			} else {
				slog.Info("port exists host(%s|%s) port(%d)", host.UUID, host.Ip, port.Port)
			}
		}

	case GETDEVICES: //获取网络设备
		req := &ProxyReq{}
		if err := json.Unmarshal(data[1:], &req); err != nil {
			slog.Error("HandlerMessage() Parse GetDevicesReq error(%s), data(%s)", err.Error(), string(data[1:]))
			return
		}
		if len(req.Proxy) == 0 {
			return
		}
		devices, err := mysql.GetDeviceByProxy(req.Proxy)
		if err != nil {
			slog.Error("HandlerMessage() Parse GetDeviceByProxy error(%s), proxy(%s)", err.Error(), req.Proxy)
			return
		}
		if len(devices) == 0 {
			return
		}
		resp, err := json.Marshal(devices)
		if err != nil {
			slog.Error("HandlerMessage() Marshal devices error(%s), proxy(%s), devices(%#v)", err.Error(), req.Proxy, devices)
			return
		}
		if err := conn.AsyncWriteData(NewPacket(GETDEVICESRESP, resp)); err != nil {
			slog.Error("HandlerMessage() send devices error %s proxy(%s)", err.Error(), req.Proxy)
		} else {
			slog.Info("HandleMessage() send devices to proxy(%s) done.", req.Proxy)
		}

	case GETPORTS:
		req := &ProxyReq{}
		if err := json.Unmarshal(data[1:], &req); err != nil {
			slog.Error("HandlerMessage() Parse GetDevicesReq error(%s), data(%s)", err.Error(), string(data[1:]))
			return
		}
		if len(req.Proxy) == 0 {
			return
		}
		ports := mysql.GetPortsByProxy(req.Proxy)
		if len(ports) == 0 {
			slog.Info("GetPortsByProxy no ports %s", req.Proxy)
			return
		}
		resp, err := json.Marshal(ports)
		if err != nil {
			slog.Error("HandlerMessage() Marshal ports error(%s), proxy(%s) ports(%#v)", err.Error(), req.Proxy, ports)
			return
		}
		if err := conn.AsyncWriteData(NewPacket(GETPORTSRESP, resp)); err != nil {
			slog.Error("HandlerMessage() send ports error %s proxy(%s)", err.Error(), req.Proxy)
		} else {
			slog.Info("HandleMessage() send prots to proxy(%s) done.", req.Proxy)
		}
	case UPDATEPORTS:
		ports := make([]*Port, 0)
		if err := json.Unmarshal(data[1:], &ports); err != nil {
			slog.Error("HandlerMessage() Parse UPDATEPORTS error(%s), data(%s)", err.Error(), string(data[1:]))
			return
		}
		for _, port := range ports {
			if err := mysql.UpdatePortStatus(port); err != nil {
				slog.Error("update port status error(%s) port(%v)", err.Error(), port)
				continue
			}
			slog.Info("update port status done, port(%v)", port)
		}

	case UPDATEDEVICES:

		devices := make([]*NetDevice, 0)
		if err := json.Unmarshal(data[1:], &devices); err != nil {
			slog.Error("HandlerMessage() Parse UPDATEDEVICES error(%s), data(%s)", err.Error(), string(data[1:]))
			return
		}
		for _, device := range devices {
			for _, i := range device.DeviceInterfaces {
				name := strings.ToLower(i.Name)
				if strings.Contains(name, "vlan") ||
					strings.Contains(name, "aux") ||
					strings.Contains(name, "loop") ||
					strings.Contains(name, "null") {
					continue
				}
				if mysql.InterfaceExist(i, device) {
					mysql.UpdateInterfaces(device, i)
				} else {
					slog.Info("device %s have a new interface %s, insert to mysql", device.IpAddr, i.Name)
					if err := mysql.CreateInterface(i, device); err != nil {
						slog.Info("device %s interface %s insert to mysql error %s", device.IpAddr, i.Name, err)
						continue
					} else {
						slog.Info("device %s interface %s insert to mysql done", device.IpAddr, i.Name)
					}
				}
			}
		}

	case ASSESTDISCOVERY:
		assest := &AssestPacket{}
		if err := json.Unmarshal(data[1:], &assest); err != nil {
			slog.Error("HandlerMessage() Parse ASSESTDISCOVERY error(%s), data(%s)", err.Error(), string(data[1:]))
			return
		}
		host, err := mysql.GetHostByUUID(assest.UUID)
		if err != nil {
			slog.Error("HandlerMessage() ASSESTDISCOVERY get host by uuid error(%s) assest(%v)", err.Error(), assest)
			return
		}
		//有资产信息
		if host.AssestId != 0 {
			oldassest := mysql.GetAssestByID(host.AssestId)
			flag := 0
			if oldassest.SN != assest.SN {
				oldassest.SN = assest.SN
				flag = 1
			}
			if oldassest.Vender != assest.Vender {
				oldassest.Vender = assest.Vender
				flag = 1
			}
			if oldassest.Model != assest.Model {
				oldassest.Model = assest.Model
				flag = 1
			}
			if flag == 1 {
				if err := mysql.UpdateAssest(oldassest); err != nil {
					slog.Error("update assest error(%s) assest(%v)", err.Error(), oldassest)
					return
				}
			}
			for _, cpu := range assest.Cpus {
				if c, ok := oldassest.Cpus[cpu.Name]; ok {
					cpu.ID = c.ID
					mysql.UpdateCpu(cpu)

				} else {
					mysql.CreateCpu(oldassest.ID, cpu)
				}

			}
			for _, disk := range assest.Disks {
				if d, ok := oldassest.Disks[disk.Slot]; ok {
					disk.ID = d.ID
					mysql.UpdateDisk(disk)
				} else {
					mysql.CreateDisk(oldassest.ID, disk)
				}

			}
			for _, nic := range assest.Nics {
				if n, ok := oldassest.Nics[nic.Name]; ok {
					nic.ID = n.ID
					mysql.UpdateNic(nic)
				} else {
					mysql.CreateNic(oldassest.ID, nic)
				}

			}
			for _, memory := range assest.Memorys {
				if m, ok := oldassest.Memorys[memory.Locator]; ok {
					memory.ID = m.ID
					mysql.UpdateMemory(memory)
				} else {
					mysql.CreateMemory(oldassest.ID, memory)
				}
			}
			return
		}
		//创建资产
		id, err := mysql.CreateAssest(assest)
		if err != nil {
			slog.Error("add assest error(%s) assest(%v)", err.Error(), assest)
			return
		}
		if err := mysql.UpdateHostAssest(host.ID, id); err != nil {
			slog.Error("linked host and assest error (%s) host(%v) assest(%v)", err.Error(), host, assest)
			return
		}
		slog.Info("add assest done, assest(%v)", assest)
	case HOSTHB:
		p := HeartBeat{}
		if err := json.Unmarshal(data[1:], &p); err != nil {
			slog.Error("HandlerMessage() Parse HOSTHB error(%s), data(%s)", err.Error(), string(data[1:]))
			return
		}
		host, err := mysql.GetHostByUUID(p.UUID)
		if err != nil {
			slog.Error("HandlerMessage() GetHostByUUID() error message(%s)", err.Error())
			return
		}
		if err := mysql.UpdateHostLastCheck(host); err != nil {
			slog.Error("update host last_check eror host(%v) error(%s)", host, err.Error())
		} else {
			slog.Info("update host last_check, host(%s:%s)", host.UUID, host.Ip)
		}
	default:
		slog.Warn("unknown protocol type, ip:%s,  data %s", conn.RemoteAddr(), string(data))
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
