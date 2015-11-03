package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"strings"
	"tcpserver"
	"time"
)

func (host *Host) SaveToFile() error {
	res, err := json.MarshalIndent(host, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile("./var/host.json", res, 0644)
}

func (host *Host) LoadFromFile() error {
	res, err := ioutil.ReadFile("./var/host.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(res, &host)
}

//循环执行插件获取数据放入缓存channel中
func (host *Host) Monitor(buffer chan []byte) {
	for {
		if len(host.Services) == 0 || host.ID == 0 {
			log.Warn("no services, wait 30 second retry..")
			time.Sleep(time.Second * 30)
			continue
		}
		if len(DataBuffer) == cfg.BUFFER_SIZE {
			log.Critical("buffer is full , sleep 5 minute")
			time.Sleep(time.Minute * 5)
			continue
		}
		for _, s := range host.Services {
			now := time.Now().Sub(s.LastCheck).Seconds()
			if now < float64(s.ExecInterval) {
				continue
			}
			s.LastCheck = time.Now()
			go func(s *Service) {
				TimeStamp := time.Now().Unix()
				if res, err := s.Exec(); err != nil {
					log.Error("exec plugin error plugin(%s) args(%s) error(%s)", s.Plugin, s.Args, err)
					return
				} else {

					for k, v := range res {
						data := &SeriesData{
							UUID:        host.UUID,
							IP:          host.Ip,
							Group:       host.Group,
							TimeStamp:   TimeStamp,
							ServiceName: s.Name,
							Metric:      k,
						}
						key := s.Name + "." + k
						if cache, ok := MetricCache[key]; ok {
							cache[1] = cache[0]
							cache[0] = v
						} else {
							MetricCache[key] = [2]float64{v, 0}
							continue
						}

						if i, ok := s.Items[k]; ok {
							cache := MetricCache[key]
							cache[1] = cache[0]
							cache[0] = v
							MetricCache[key] = cache
							log.Debug("plugin(%s) args(%s) key:%s raw data:[%v->%v]", s.Plugin, s.Args, s.Items[k].Key, v, cache)
							switch i.DT {
							case "GAUGE": //原始值
								data.Value = cache[0]
							case "COUNTER": //平均值
								if val := (cache[0] - cache[1]) / float64(s.ExecInterval); val < 0 {
									continue
								} else {
									data.Value = val
								}
							case "DERIVE": //差值
								if val := (cache[0] - cache[1]); val < 0 {
									continue
								} else {
									data.Value = val
								}
							default:
								log.Warn("unknown items data type item(%s), type(%s) ", s.Items[k].Key, s.Items[k].DT)
								continue
							}

						} else {
							s.Items[k] = &Item{Key: k, DT: "GAUGE"}
							data.Flag = 1
							data.Value = v
						}
						if buf, err := data.Encode(); err != nil {
							log.Error("SeriesData encode error data(%v) error(%v)", data, err)
						} else {
							buffer <- buf
						}
					}
				}

			}(s)
		}
		time.Sleep(1e9)
	}
}

func (this *Service) Exec() (map[string]float64, error) {
	command := "./plugins/" + this.Plugin
	args := strings.Split(this.Args, " ")
	tmp := map[string]float64{}
	result, err := exec.Command(command, args...).Output()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(result, &tmp); err != nil {
		return nil, err
	}
	return tmp, nil
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

//定时更新配置文件
func (this *Host) Loop(server *tcpserver.Server) {
	//主机配置更新定时器
	HostUpdateTicker := time.NewTicker(time.Minute * time.Duration(cfg.HOSTUPDATEINTERVAL))
	AssestCollectTicker := time.NewTicker(time.Minute * time.Duration(cfg.ASSESTINTERVAL))
	PortCollectTicker := time.NewTicker(time.Minute * time.Duration(cfg.PORTINTERVAL))
	go func() {
		for {
			if this.Ip == "" {
				log.Warn("host ip is blank, wait 30 second continue")
				time.Sleep(time.Second * 30)
				continue
			}
			select {
			case <-HostUpdateTicker.C:
				//获取最新监控配置信息
				DataBuffer <- this.GetConfigPacket()
				DataBuffer <- this.GetVersionPacket()

			case <-AssestCollectTicker.C:
				req, err := this.CollectAssest(cfg.ASSECTPLUGIN)
				if err != nil {
					log.Error("collect assect info error(%s)", err.Error())
					continue
				}
				DataBuffer <- req
			case <-PortCollectTicker.C:
				req, err := this.CollectPort(cfg.PORTPLUGIN)
				if err != nil {
					log.Error("collect port info error(%s)", err.Error())
					continue
				}
				DataBuffer <- req
			}
		}
	}()
START:
	conn, err := server.Connect(cfg.TCPSERVER)
	if err != nil {
		log.Error("connect to server(%s) error(%s)", cfg.TCPSERVER, err.Error())
		time.Sleep(time.Duration(cfg.RECONNECTINTERVAL) * time.Second)
		goto START
	} else {
		log.Info("connect to server(%s) done", cfg.TCPSERVER)
	}

	go conn.Run()
	this.Ip = conn.GetLocalIp()
	if len(DataBuffer) < cfg.BUFFER_SIZE/2 {
		DataBuffer <- this.GetConfigPacket()
		DataBuffer <- this.GetVersionPacket()
	}
	req, err := this.CollectAssest(cfg.ASSECTPLUGIN)
	if err == nil && len(DataBuffer) < cfg.BUFFER_SIZE/2 {
		DataBuffer <- req
	}
	for {
		select {
		case data, ok := <-DataBuffer:
			if ok {
				if err := conn.AsyncWriteData(data); err != nil {
					log.Critical("write data to server error message(%v)", err)
					DataBuffer <- data
					goto START
				}
				log.Info("send data: %s|%s", PROTOCOLTYPE[data[4]], string(data[5:]))
			}
		}
	}
}

func (this *Host) GetConfigPacket() []byte {
	var buf bytes.Buffer
	head := make([]byte, 4)
	req := &HostConfigReq{
		UUID:      this.UUID,
		Proxy:     this.Proxy,
		IP:        this.Ip,
		OS:        GetOs(),
		Kernel:    GetKernel(),
		IDRACAddr: GetIDRACAddr(),
		HostName:  GetHostName(),
	}
	b, _ := json.Marshal(req)
	binary.BigEndian.PutUint32(head, uint32(len(b)+1))
	binary.Write(&buf, binary.BigEndian, head)
	binary.Write(&buf, binary.BigEndian, byte(HOSTCONFIG))
	binary.Write(&buf, binary.BigEndian, b)
	return buf.Bytes()
}

func (this *Host) GetVersionPacket() []byte {
	var buf bytes.Buffer
	head := make([]byte, 4)
	binary.BigEndian.PutUint32(head, uint32(1))
	binary.Write(&buf, binary.BigEndian, head)
	binary.Write(&buf, binary.BigEndian, byte(CLIENTVERSION))
	return buf.Bytes()
}

func (this *Host) CollectAssest(plugin string) ([]byte, error) {
	command := "./plugins/" + plugin
	result, err := exec.Command(command).Output()
	assest := &AssestPacket{
		UUID: this.UUID,
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(result, &assest); err != nil {
		return nil, err
	}
	res, _ := json.Marshal(assest)
	var buf bytes.Buffer
	head := make([]byte, 4)
	binary.BigEndian.PutUint32(head, uint32(len(res)+1))
	binary.Write(&buf, binary.BigEndian, head)
	binary.Write(&buf, binary.BigEndian, byte(ASSESTDISCOVERY))
	binary.Write(&buf, binary.BigEndian, res)
	return buf.Bytes(), nil
}

func (this *Host) CollectPort(plugin string) ([]byte, error) {
	command := "./plugins/" + plugin
	result, err := exec.Command(command).Output()
	if err != nil {
		return nil, err
	}
	tmp := make(map[string][]int)
	if err := json.Unmarshal(result, &tmp); err != nil {
		return nil, err
	}
	req := &PortDiscovery{
		UUID:  this.UUID,
		Ports: []*Port{},
	}
	for k, v := range tmp {
		for _, p := range v {
			req.Ports = append(req.Ports, &Port{
				ProcName: k,
				Port:     p,
			})
		}
	}
	res, _ := json.Marshal(req)
	var buf bytes.Buffer
	head := make([]byte, 4)
	binary.BigEndian.PutUint32(head, uint32(len(res)+1))
	binary.Write(&buf, binary.BigEndian, head)
	binary.Write(&buf, binary.BigEndian, byte(PORTDISCOVERYR))
	binary.Write(&buf, binary.BigEndian, res)
	return buf.Bytes(), nil
}

func (this *Host) ServerHB() {
	for {
		if len(DataBuffer) > 20 {
			time.Sleep(time.Second * 30)
			continue
		}
		var buf bytes.Buffer
		head := make([]byte, 4)
		p := HeartBeat{
			UUID: this.UUID,
		}
		res, _ := json.Marshal(p)
		binary.BigEndian.PutUint32(head, uint32(len(res)+1))
		binary.Write(&buf, binary.BigEndian, head)
		binary.Write(&buf, binary.BigEndian, byte(HOSTHB))
		binary.Write(&buf, binary.BigEndian, res)
		DataBuffer <- buf.Bytes()
		time.Sleep(time.Second * 30)
	}
}
