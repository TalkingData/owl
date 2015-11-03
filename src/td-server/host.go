package main

import (
	"fmt"
	"net"
	"time"
)

type Host struct {
	ID         int    `json:"id"`
	UUID       string `json:"uuid"`
	AssestId   int    `json:"assest_id"`
	Ip         string `json:"ip"`
	Group      string `json:"group"`
	Status     int    `json:"status"`
	Proxy      string `json:"proxy"`
	last_check string
	Services   []*Service `json:"services"`
}

type Service struct {
	Id           int              `json:"id"`
	Name         string           `json:"name"`
	Plugin       string           `json:"plugin"`
	Args         string           `json:"args"`
	ExecInterval int              `json:"exec_interval"`
	Items        map[string]*Item `json:"items"`
	//lastcheck    time.Time
}

type Item struct {
	Key       string `json:"key"`
	DT        string `json:"data_type"` //GAUGE|COUNTER|DERIVE
	alarm     int
	attempt   int
	symbol    string
	method    string
	threshold float64
	cycle     int
	drawing   int
	counter   int
	duration  int
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

type Assest struct {
	ID      int                `json:"id"`
	UUID    string             `json:"uuid"`
	SN      string             `json:"sn"`
	Unit    string             `json:"unit"`
	Vender  string             `json:"vender"`
	Model   string             `json:"model"`
	Cpus    map[string]*Cpu    `json:"cpus"`    //key is cpu.name
	Memorys map[string]*Memory `json:"memorys"` //key is memory.locator
	Disks   map[string]*Disk   `json:"disks"`   //key is disk.slot
	Nics    map[string]*Nic    `json:"nics"`    //key is nic.name
}

type ServiceTemplate struct {
	ID       int
	Name     string
	Services []*Service
	Hosts    []*Host
}

//检查端口状态是否正常
//超时时间5秒
func (this *Port) StatusCheck() error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port), time.Second*5)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func PortCheck() {
	for {
		for _, port := range mysql.GetPortsByProxy("") {
			go func(port *Port) {
				if err := port.StatusCheck(); err != nil {
					if port.Status == 0 {
						port.Status = 1
						if err := mysql.UpdatePortStatus(port); err != nil {
							slog.Error("update port status error port(%v) message(%v)", port, err)
						} else {
							slog.Info("update port status done port(%v) ", port)
						}
					}
				} else {
					if port.Status == 1 {
						port.Status = 0
						if err := mysql.UpdatePortStatus(port); err != nil {
							slog.Error("update port status error port(%v) message(%v)", port, err)
						} else {
							slog.Info("update port status done port(%v) ", port)
						}
					}
				}
			}(port)
		}
		time.Sleep(time.Minute * 2)
	}
}

func HostLoop() {
	for {
		time.Sleep(time.Second * 30)

		for _, host := range mysql.GetAllHosts() {

			last_check, _ := time.ParseInLocation("2006-01-02 15:04:05", host.last_check, time.Local)
			interval := time.Now().Sub(last_check).Seconds()
			flag := 0
			if interval > 60 && host.Status == 0 {
				host.Status = 1
				flag = 1
			} else if interval < 60 && host.Status == 1 {
				host.Status = 0
				flag = 158
			}
			if flag != 0 {
				if err := mysql.UpdateHostStatus(host); err != nil {
					slog.Error("update host status error host(%#v) error(%s)", host, err.Error())
				} else {
					slog.Info("update host status success host(%#v)", host)
				}
			}
		}
		go ServiceTemplateHandle()
	}
}

func ServiceTemplateHandle() {
	for _, tpl := range mysql.GetServiceTempates() {
		for _, service := range tpl.Services {
			for _, host := range tpl.Hosts {
				flag := 0
				for _, s := range host.Services {
					if service.Name == s.Name {
						flag = 1
						continue
					}
				}
				if flag == 0 {
					if err := mysql.CreateServiceByHostID(host.ID, service); err != nil {
						slog.Error("ServiceTemplateHandle() CreateServiceByHostID() error host(%s:%s) service(%s) %s",
							host.UUID, host.Ip, service.Name, err.Error())
					} else {
						for _, item := range service.Items {
							if err := mysql.CreateItemByService(service, item); err != nil {
								slog.Error("ServiceTemplateHandle() CreateItemByService error host(%s:%s) service(%s) item(%s) %s",
									host.UUID, host.Ip, service.Name, item.Key, err.Error())
							}
						}
					}
				}
			}
		}
	}
}
