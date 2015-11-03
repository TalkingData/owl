package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type db struct {
	Conn *sql.DB
}

func NewMysqlConnPool(cfg *Config) (*db, error) {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		cfg.MYSQL_USER, cfg.MYSQL_PASSWORD, cfg.MYSQL_BIND, cfg.MYSQL_DBNAME))
	if err != nil {
		return nil, err
	}
	row, err := conn.Query("select 1")
	if err != nil {
		return nil, err
	}
	defer row.Close()
	conn.SetMaxIdleConns(cfg.MYSQL_MAX_IDLE_CONN)
	conn.SetMaxOpenConns(cfg.MYSQL_MAX_CONN)
	return &db{conn}, nil
}

//－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－
//网络设备相关
//－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－

//func NewNetWorkDevice(ip, proxy, snmpVersion, snmpCommunity string, snmpPort, updateInterval,
//	checkInterval int, customoid map[string]*CustomOid) *NetDevice {

//	return &NetDevice{
//		IpAddr: ip,
//		Proxy:  proxy,
//		//		stopChan:         make(chan struct{}),
//		//		dataChan:         make(chan []byte, 1000),
//		SnmpCommunity:    snmpCommunity,
//		SnmpVersion:      snmpVersion,
//		SnmpPort:         snmpPort,
//		updateTicker:     time.NewTicker(time.Second * time.Duration(updateInterval)),
//		checkTicker:      time.NewTicker(time.Second * time.Duration(checkInterval)),
//		UpdateInterval:   updateInterval,
//		CheckInterval:    checkInterval,
//		CustomOids:       customoid,
//		DeviceInterfaces: make(map[string]*DeviceInterface),
//	}
//}

//获取所有网络设备信息
func (this *db) GetAllNetDevice() ([]*NetDevice, error) {
	rows, err := this.Conn.Query("select d.id, d.uuid, d.ip, d.snmp_version, d.snmp_community, d.snmp_port," +
		"d.config_update_interval, d.check_interval, p.ip from network_device as d left join system_proxy as p on d.proxy_id = p.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	devices := []*NetDevice{}
	for rows.Next() {
		device := &NetDevice{}
		var proxy sql.NullString
		if err := rows.Scan(&device.ID, &device.UUID, &device.IpAddr, &device.SnmpVersion,
			&device.SnmpCommunity, &device.SnmpPort, &device.UpdateInterval,
			&device.CheckInterval, &proxy); err != nil {
			continue
		}
		device.Proxy = proxy.String
		device.stopChan = make(chan struct{})
		device.updateTicker = time.NewTicker(time.Second * time.Duration(device.UpdateInterval))
		device.checkTicker = time.NewTicker(time.Second * time.Duration(device.CheckInterval))
		device.DeviceInterfaces = make(map[string]*DeviceInterface)
		//获取oid
		if oids, err := this.GetCustomOidByDeviceID(device.ID); err == nil {
			device.CustomOids = oids
		}
		//获取端口
		ifts := this.GetInterfacesByDeviceId(device.ID)
		for _, ift := range ifts {
			device.DeviceInterfaces[string(ift.Index)] = ift
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (this *db) GetDeviceByProxy(proxy string) ([]*NetDevice, error) {
	devices, err := this.GetAllNetDevice()
	if err != nil {
		return nil, err
	}
	devs := []*NetDevice{}
	for _, dev := range devices {
		if dev.Proxy == proxy {
			devs = append(devs, dev)
		}
	}
	return devs, nil

}

//根据网络设备id获取端口信息
func (this *db) GetInterfacesByDeviceId(id int) []*DeviceInterface {
	rows, err := this.Conn.Query("select id, index, name, mac, speed, status from network_interface where device_id=?", id)
	if err != nil {
		return nil
	}
	defer rows.Close()
	ifts := []*DeviceInterface{}
	for rows.Next() {
		ift := &DeviceInterface{
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
		if err := rows.Scan(&ift.ID, &ift.Index, &ift.Name, &ift.Mac, &ift.Speed, &ift.Status); err != nil {
			continue
		}
		ifts = append(ifts, ift)
	}
	return ifts
}

//根据网络设备id获取自定义oid信息
func (this *db) GetCustomOidByDeviceID(id int) (map[string]*CustomOid, error) {
	rows, err := this.Conn.Query("select name, oid from network_oid where device_id=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	oids := make(map[string]*CustomOid)
	for rows.Next() {
		oid := &CustomOid{}
		if err := rows.Scan(&oid.Name, &oid.OID); err != nil {
			continue
		}
		oids[oid.OID] = oid
	}
	return oids, nil
}

//根据id号查询设备
func (this *db) GetDeviceByID(id int) *NetDevice {
	devices, err := this.GetAllNetDevice()
	if err != nil {
		return nil
	}
	for _, dev := range devices {
		if dev.ID == id {
			return dev
		}
	}
	return nil
}

//端口是否存在数据库
func (this *db) InterfaceExist(i *DeviceInterface, device *NetDevice) bool {
	var id int
	if err := this.Conn.QueryRow("select id from network_interface where device_id=? and `index`=?",
		device.ID, i.Index).Scan(&id); err != nil {
		return false
	}
	return true
}

//数据库中创建端口
func (this *db) CreateInterface(i *DeviceInterface, device *NetDevice) error {
	row, err := this.Conn.Exec("insert into network_interface(device_id, `index`, name, mac, speed, status,alarm) values(?,?,?,?,?,?,0)",
		device.ID, i.Index, i.Name, i.Mac, i.Speed, i.Status)
	if err != nil {
		return err
	}
	id, err := row.LastInsertId()
	if err != nil {
		return err
	}
	i.ID = int(id)
	//TODO: auto add monitor item
	return nil
}

//更新网络设备接口状态
func (this *db) UpdateInterfaces(dev *NetDevice, i *DeviceInterface) error {
	_, err := this.Conn.Exec("update network_interface set name=?, mac=?, speed=?, status=? where `index`=? and device_id=?",
		i.Name, i.Mac, i.Speed, i.Status, i.Index, dev.ID)
	if err != nil {
		return err
	}
	return nil
}

func (this *db) GetHostConfigByUUID(uuid string) (*Host, error) {
	host, err := this.GetHostByUUID(uuid)
	if err != nil {
		return nil, err
	}
	this.GetGroupByHost(host)
	this.GetServiceByHost(host)
	for _, s := range host.Services {
		this.GetItemByService(s)
	}
	return host, nil
}

func (this *db) GetHostByUUID(uuid string) (*Host, error) {
	host := &Host{}
	var (
		proxy     sql.NullString
		assest_id sql.NullInt64
	)

	if err := this.Conn.QueryRow("select h.id, h.server_id, h.uuid, h.ip, h.status, p.ip from "+
		"host_host as h left join system_proxy as p on h.proxy_id = p.id where h.uuid=?", uuid).Scan(
		&host.ID, &assest_id, &host.UUID, &host.Ip, &host.Status, &proxy); err != nil {
		return nil, err
	}
	host.AssestId = int(assest_id.Int64)
	host.Proxy = proxy.String
	return host, nil
}

func (this *db) GetGroupByHost(host *Host) {
	if host.ID == 0 {
		return
	}
	rows, err := this.Conn.Query("select name from host_group where id in(select group_id from host_host_group where host_id=?)", host.ID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var group string
		if err := rows.Scan(&group); err != nil {
			continue
		}
		host.Group = host.Group + group + ","
	}
	host.Group = strings.Trim(host.Group, ",")
}

func (this *db) GetAllHosts() []*Host {
	rows, err := this.Conn.Query("select h.id, h.server_id, h.uuid, h.ip, h.status, h.last_check, p.ip " +
		"from host_host as h left join system_proxy as p on h.proxy_id = p.id ")
	if err != nil {
		slog.Error("GetAllHosts() error(%s)", err.Error())
		return nil
	}
	defer rows.Close()
	hosts := []*Host{}
	for rows.Next() {
		host := &Host{}
		var proxy sql.NullString
		var assestid sql.NullInt64
		if err := rows.Scan(&host.ID, &assestid, &host.UUID, &host.Ip,
			&host.Status, &host.last_check, &proxy); err != nil {
			slog.Error("GetAllHosts() Scan error %s", err.Error())
			continue
		}
		host.Proxy = proxy.String
		if assestid.Valid {
			host.AssestId = int(assestid.Int64)
		}
		hosts = append(hosts, host)
	}
	return hosts
}

func (this *db) GetHostsByProxy(proxy string) []*Host {
	var (
		rows *sql.Rows
		err  error
	)
	if len(proxy) == 0 {
		rows, err = this.Conn.Query("select h.id, h.server_id, h.uuid, h.ip, h.status, h.last_check, p.ip from host_host as h " +
			"left join system_proxy as p on h.proxy_id = p.id where p.ip is NULL")
	} else {
		rows, err = this.Conn.Query("select h.id, h.server_id, h.uuid, h.ip, h.status, h.last_check, p.ip from host_host as h "+
			"left join system_proxy as p on h.proxy_id = p.id where p.ip=?", proxy)
	}
	if err != nil {
		return nil
	}
	defer rows.Close()
	hosts := []*Host{}
	for rows.Next() {
		host := &Host{}
		var proxy sql.NullString
		var assestid sql.NullInt64
		if err := rows.Scan(&host.ID, &assestid, &host.UUID, &host.Ip, &host.Status, &host.last_check, &proxy); err != nil {
			continue
		}
		host.Proxy = proxy.String
		if assestid.Valid {
			host.AssestId = int(assestid.Int64)
		}
		hosts = append(hosts, host)
	}
	return hosts
}

func (this *db) CreateHost(config *HostConfigReq) error {
	now := time.Now().Format("2006-01-02 15:03:04")
	if len(config.Proxy) == 0 {
		_, err := this.Conn.Exec("insert into host_host(uuid,ip, status, alarm, proxy_id, c_time, last_check, hostname, os, kernel, idrac) values(?,?,?,?,NULL,?,?,?,?,?,?)",
			config.UUID, config.IP, 0, 0, now, now, config.HostName, config.OS, config.Kernel, config.IDRACAddr)
		return err
	} else {
		id, err := this.GetProxyID(config.Proxy)
		if err != nil {
			return err
		}
		_, err = this.Conn.Exec("insert into host_host(uuid,ip, status, alarm, proxy_id, c_time, last_check, hostname, os, kernel, idrac) values(?,?,?,?,?,?,?,?,?,?,?)",
			config.UUID, config.IP, 0, 0, id, now, now, config.HostName, config.OS, config.Kernel, config.IDRACAddr)
		return err
	}
	return nil
}

func (this *db) UpdateHostLastCheck(host *Host) error {
	last_check := time.Now().Format("2006-01-02 15:04:05")
	_, err := this.Conn.Exec("update host_host set last_check=? where id=?", last_check, host.ID)
	return err
}

func (this *db) GetProxyID(proxy string) (int, error) {
	var id int
	err := this.Conn.QueryRow("select id from system_proxy where ip=?", proxy).Scan(&id)
	return id, err
}

func (this *db) UpdateHostStatus(host *Host) error {
	_, err := this.Conn.Exec("update host_host set status=? where uuid=?", host.Status, host.UUID)
	return err
}

func (this *db) UpdateHostInfo(cfg *HostConfigReq) error {
	_, err := this.Conn.Exec("update host_host set ip=?, hostname=?, os=?, kernel=?, idrac=? where uuid=?",
		cfg.IP, cfg.HostName, cfg.OS, cfg.Kernel, cfg.IDRACAddr, cfg.UUID)
	return err
}

//获取主机下的服务
//主机ID是必须的
func (this *db) GetServiceByHost(host *Host) error {
	rows, err := this.Conn.Query("select id, name, plugin, args, exec_interval from host_service where host_id=?", host.ID)
	if err != nil {
		return err
	}
	//host.Services = make(map[string]*Service)
	host.Services = []*Service{}
	defer rows.Close()
	for rows.Next() {
		s := &Service{}
		if err := rows.Scan(&s.Id, &s.Name, &s.Plugin, &s.Args, &s.ExecInterval); err != nil {
			continue
		}
		//host.Services[s.Name] = s
		host.Services = append(host.Services, s)
	}
	return nil
}

func (this *db) GetPortsByHost(host *Host) []*Port {
	rows, err := this.Conn.Query("select `id`, `port`, `proc_name`, `status` from host_port where host_id=?", host.ID)
	if err != nil {
		return nil
	}
	ports := make([]*Port, 0)
	defer rows.Close()
	for rows.Next() {
		port := &Port{}
		if err := rows.Scan(&port.Id, &port.Port, &port.ProcName, &port.Status); err != nil {
			continue
		}
		port.Ip = host.Ip
		ports = append(ports, port)
	}
	return ports
}

func (this *db) GetItemByService(s *Service) error {
	rows, err := this.Conn.Query("select `key`, `dt`, `attempt`, `symbol`, `method`, `threshold`, `cycle`, `drawing`, `alarm` from host_item where service_id=?", s.Id)
	if err != nil {
		return err
	}
	s.Items = make(map[string]*Item)
	defer rows.Close()
	for rows.Next() {
		i := &Item{}
		if err := rows.Scan(&i.Key, &i.DT, &i.attempt, &i.symbol, &i.method, &i.threshold, &i.cycle, &i.drawing, &i.alarm); err != nil {
			continue
		}
		s.Items[i.Key] = i
	}
	return nil
}

func (this *db) GetVersion() (version string) {

	_ = this.Conn.QueryRow("SELECT version FROM host_agent order by timestramp desc limit 1").Scan(&version)
	return version
}

func (this *db) CreateItemByService(s *Service, i *Item) error {
	var name string
	if err := this.Conn.QueryRow("select `id` from host_item where service_id = ? and `key`=?", s.Id, i.Key).Scan(&name); err == nil {
		return errors.New("item is exists.")
	}
	timestamp := time.Now().Format("2006-01-02 15:03:04")
	_, err := this.Conn.Exec("insert into host_item(service_id, `key`, last_check,"+
		"duration, attempt, counter, symbol, method, threshold, units, current, cycle, dt, alarm, drawing) "+
		"values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", s.Id, i.Key, timestamp, timestamp, i.attempt, i.counter,
		i.symbol, i.method, i.threshold, "", 0, i.cycle, i.DT, i.alarm, i.drawing)
	return err
}

func (this *db) GetPortsByProxy(proxy string) []*Port {
	hosts := this.GetHostsByProxy(proxy)
	if len(hosts) == 0 {
		return nil
	}
	ports := make([]*Port, 0)
	for _, host := range hosts {
		ports = append(ports, this.GetPortsByHost(host)...)
	}
	return ports
}

func (this *db) PortIsExists(host_id, port int) bool {
	var id int
	err := this.Conn.QueryRow("select id from host_port where host_id=? and port=?", host_id, port).Scan(&id)
	if err != nil {
		return false
	}
	return true
}

func (this *db) CreatePort(host_id int, port *Port) error {
	_, err := this.Conn.Exec("insert into host_port(`host_id`, `port`, `proc_name`,`status`, `alarm`) values(?,?,?,?,0)",
		host_id, port.Port, port.ProcName, port.Status)
	return err
}

func (this *db) UpdatePortStatus(port *Port) error {
	_, err := this.Conn.Exec("update host_port set status=? where id=?", port.Status, port.Id)
	return err
}

func (this *db) GetAssestByID(id int) *Assest {
	assest := &Assest{}
	err := this.Conn.QueryRow("select id, sn,vender, model from assest_server where id =?", id).Scan(
		&assest.ID, &assest.SN, &assest.Vender, &assest.Model)
	if err != nil {
		return nil
	}
	//assest.Cpus = make(map[string]*Cpu)
	assest.Cpus = this.GetCpusByAssestID(assest.ID)
	//assest.Memorys := make(map[string]*Memory)
	assest.Memorys = this.GetMemorysByAssestID(assest.ID)
	//assest.Nics := make(map[string]*Nic)
	assest.Nics = this.GetNicsByAssestID(assest.ID)
	//assest.Disks := make(map[string]*Disk)
	assest.Disks = this.GetDisksByAssestID(assest.ID)
	return assest
}

func (this *db) GetCpusByAssestID(assest_id int) map[string]*Cpu {
	rows, err := this.Conn.Query("select id, name, sn ,model, vender from assest_cpu where server_id=?", assest_id)
	if err != nil {
		return nil
	}
	defer rows.Close()
	cpus := make(map[string]*Cpu)
	for rows.Next() {
		cpu := &Cpu{}
		if err := rows.Scan(&cpu.ID, &cpu.Name, &cpu.SN, &cpu.Model, &cpu.Vender); err != nil {
			continue
		}
		cpu.AssestId = assest_id
		cpus[cpu.Name] = cpu
	}
	return cpus
}

func (this *db) GetMemorysByAssestID(assest_id int) map[string]*Memory {
	rows, err := this.Conn.Query("select id, sn, size, speed, locator, vender from assest_memory where server_id=?", assest_id)
	if err != nil {
		return nil
	}
	defer rows.Close()
	mems := make(map[string]*Memory)
	for rows.Next() {
		mem := &Memory{}
		if err := rows.Scan(&mem.ID, &mem.SN, &mem.Size, &mem.Speed, &mem.Locator, &mem.Vender); err != nil {
			continue
		}
		mem.AssestId = assest_id
		mems[mem.Locator] = mem
	}
	return mems
}

func (this *db) GetDisksByAssestID(assest_id int) map[string]*Disk {
	rows, err := this.Conn.Query("select id, sn, vender, product_id, slot, "+
		"production_date, capacity, speed, bus, media from assest_disk where server_id=?", assest_id)
	if err != nil {
		return nil
	}
	defer rows.Close()
	disks := make(map[string]*Disk)
	for rows.Next() {
		disk := &Disk{}
		if err := rows.Scan(&disk.ID, &disk.SN, &disk.Vender, &disk.ProductId, &disk.Slot,
			&disk.ProductionDate, &disk.Capacity, &disk.Speed, &disk.Bus, &disk.Media); err != nil {
			continue
		}
		disk.AssestId = assest_id
		disks[disk.Slot] = disk
	}
	return disks
}

func (this *db) GetNicsByAssestID(assest_id int) map[string]*Nic {
	rows, err := this.Conn.Query("select id, name, mac, vender, description from assest_nic where server_id=?", assest_id)
	if err != nil {
		return nil
	}
	defer rows.Close()
	nics := make(map[string]*Nic)
	for rows.Next() {
		nic := &Nic{}
		if err := rows.Scan(&nic.ID, &nic.Name, &nic.Mac, &nic.Vender, &nic.Description); err != nil {
			continue
		}
		nic.AssestId = assest_id
		nics[nic.Name] = nic
	}
	return nics
}

func (this *db) CreateAssest(assest *AssestPacket) (int, error) {
	res, err := this.Conn.Exec("insert into assest_server(`sn`, `vender`, `model`,`location`,`unit`) values(?,?,?,0,0)", assest.SN, assest.Vender, assest.Model)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	aid := int(id)
	for _, cpu := range assest.Cpus {
		this.CreateCpu(aid, cpu)
	}
	for _, disk := range assest.Disks {
		this.CreateDisk(aid, disk)
	}
	for _, memory := range assest.Memorys {
		this.CreateMemory(aid, memory)
	}
	for _, nic := range assest.Nics {
		this.CreateNic(aid, nic)
	}
	return aid, nil
}

func (this *db) UpdateAssest(assest *Assest) error {
	_, err := this.Conn.Exec("update assest_server set sn=?, vender=?, model=? where id=?",
		assest.SN, assest.Vender, assest.Model, assest.ID)
	return err
}

func (this *db) CreateCpu(assest_id int, cpu *Cpu) error {
	res, err := this.Conn.Exec("insert into assest_cpu(`server_id`, `name`, `sn`, `model`, `vender`) values(?,?,?,?,?)",
		assest_id, cpu.Name, cpu.SN, cpu.Model, cpu.Vender)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	cpu.ID = int(id)
	return nil
}

func (this *db) UpdateCpu(cpu *Cpu) error {
	_, err := this.Conn.Exec("update assest_cpu set name=?, sn=?, model=?,vender=? where id=?", cpu.Name, cpu.SN, cpu.Model, cpu.Vender, cpu.ID)
	return err
}

func (this *db) CreateDisk(assest_id int, disk *Disk) error {
	res, err := this.Conn.Exec("insert into assest_disk(server_id, sn, vender,"+
		"product_id, slot, production_date, capacity, speed, bus, media) values(?,?,?,?,?,?,?,?,?,?)",
		assest_id, disk.SN, disk.Vender, disk.ProductId, disk.Slot, disk.ProductionDate,
		disk.Capacity, disk.Speed, disk.Bus, disk.Media)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	disk.ID = int(id)
	return nil
}

func (this *db) UpdateDisk(disk *Disk) error {
	if disk.ID == 0 {
		return errors.New("disk id is required")
	}
	_, err := this.Conn.Exec("update assest_disk set sn=?, vender=?, product_id=?,"+
		"slot=?,production_date=?, capacity=?, speed=?, bus=?,media=? where id=?",
		disk.SN, disk.Vender, disk.ProductId, disk.Slot, disk.ProductionDate, disk.Capacity, disk.Speed, disk.Bus, disk.Media, disk.ID)
	return err
}

func (this *db) CreateMemory(assest_id int, memory *Memory) error {
	res, err := this.Conn.Exec("insert into assest_memory(server_id, sn, size, speed, locator, vender) values(?,?,?,?,?,?)",
		assest_id, memory.SN, memory.Size, memory.Speed, memory.Locator, memory.Vender)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	memory.ID = int(id)
	return nil
}

func (this *db) UpdateMemory(memory *Memory) error {
	_, err := this.Conn.Exec("update assest_memory set sn=?, size=?, speed=?, locator=?, vender=? where id=?",
		memory.SN, memory.Size, memory.Speed, memory.Locator, memory.Vender, memory.ID)
	return err
}

func (this *db) CreateNic(assest_id int, nic *Nic) error {
	res, err := this.Conn.Exec("insert into assest_nic(`server_id`, `name`, `mac`, `vender`, `description`) values(?,?,?,?,?)",
		assest_id, nic.Name, nic.Mac, nic.Vender, nic.Description)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	nic.ID = int(id)
	return nil
}

func (this *db) UpdateNic(nic *Nic) error {
	_, err := this.Conn.Exec("update assest_nic set name=?, mac=?, vender=?, description=? where id=?",
		nic.Name, nic.Mac, nic.Vender, nic.Description, nic.ID)
	return err
}

func (this *db) UpdateHostAssest(host_id, assest_id int) error {
	_, err := this.Conn.Exec("update host_host set server_id =? where id=?", assest_id, host_id)
	return err
}

func (this *db) GetServiceTempates() []*ServiceTemplate {
	templates := []*ServiceTemplate{}
	rows, err := this.Conn.Query("select id, name from host_template")
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		template := &ServiceTemplate{}
		if err := rows.Scan(&template.ID, &template.Name); err != nil {
			slog.Error("GetServiceTempates() Scan template error %s", err.Error())
		}
		//获取模板下的服务项
		template.Services = this.GetServiceByTemplateID(template.ID)
		template.Hosts = this.GetHostsByTemplateID(template.ID)
		templates = append(templates, template)
	}
	return templates
}

func (this *db) GetServiceByTemplateID(id int) []*Service {
	rows, err := this.Conn.Query("select id,name,plugin, args, exec_interval"+
		" from host_service where template_id = ?", id)
	if err != nil {
		slog.Error("GetServiceByTemplateID() error %s", err.Error())
		return nil
	}
	defer rows.Close()
	services := []*Service{}
	for rows.Next() {
		service := &Service{}
		if err := rows.Scan(&service.Id, &service.Name, &service.Plugin,
			&service.Args, &service.ExecInterval); err != nil {
			slog.Error("GetServiceByTemplateID() Scan error template_id(%d) %s", id, err.Error())
			continue
		}
		this.GetItemByService(service)
		services = append(services, service)
	}
	return services
}

func (this *db) GetHostsByTemplateID(id int) []*Host {
	rows, err := this.Conn.Query(" select h.uuid from host_host h,host_group_template gt,host_host_group hg  "+
		"where h.id=hg.host_id and hg.group_id=gt.group_id  and  gt.template_id= ?  "+
		"union select h.uuid from host_host h,host_host_template ht where h.id=ht.host_id and ht.template_id =? ", id, id)
	if err != nil {
		return nil
	}
	defer rows.Close()
	hosts := []*Host{}
	for rows.Next() {
		var uuid string
		if err := rows.Scan(&uuid); err != nil {
			slog.Error("GetHostsByTemplateID() Scan error %s", err.Error())
			continue
		}
		host, err := this.GetHostByUUID(uuid)
		if err != nil {
			slog.Error("GetHostsByTemplateID() GetHostByUUID() error uuid(%s) %s", uuid, err.Error())
			continue
		}
		this.GetServiceByHost(host)
		for _, s := range host.Services {
			this.GetItemByService(s)
		}
		hosts = append(hosts, host)
	}
	return hosts
}

func (this *db) CreateServiceByHostID(host_id int, service *Service) error {
	row, err := this.Conn.Exec("insert into host_service(host_id, name, plugin, args, exec_interval, status, alarm) values(?,?,?,?,?,?,?)",
		host_id, service.Name, service.Plugin, service.Args, service.ExecInterval, 0, 0)
	id, err := row.LastInsertId()
	if err != nil {
		return err
	}

	groupIDS := this.GetContactGroupIdByServiceID(service.Id)
	if len(groupIDS) == 0 || groupIDS == nil {
		return nil
	}

	for _, groupID := range groupIDS {
		if _, err := this.Conn.Exec("insert into host_service_group(service_id,group_id) values(?,?)", int(id), groupID); err != nil {
			return err
		}
	}
	service.Id = int(id)
	return nil
}

func (this *db) GetContactGroupIdByServiceID(id int) []int {
	rows, err := this.Conn.Query("select group_id from host_service_group where service_id=?", id)
	if err != nil {
		return nil
	}
	defer rows.Close()
	ids := []int{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}
