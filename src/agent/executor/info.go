package executor

import (
	"io/ioutil"
	"net"
	"os"
	"owl/common/logger"
	"owl/common/utils"
	"runtime"
	"strconv"
	"strings"
)

const (
	uuidFile   = "/sys/devices/virtual/dmi/id/product_uuid"
	uptimeFile = "/proc/uptime"
)

func (e *Executor) GetHostID() string {
	e.logger.Info("Executor.GetHostID called.")
	defer e.logger.Info("Executor.GetHostID end.")

	var id string

	if data, err := ioutil.ReadFile(uuidFile); err == nil {
		if len(data) > 0 {
			dataStr := strings.TrimSpace(string(data))
			id = utils.Md5(dataStr)
			e.logger.DebugWithFields(logger.Fields{
				"uuid_file": uuidFile,
				"data_str":  dataStr,
				"id":        id,
			}, "Success got uuid from uuid file.")
		}
	} else {
		hostname := e.GetHostname()
		id = utils.Md5(hostname)
		e.logger.DebugWithFields(logger.Fields{
			"hostname": hostname,
			"id":       id,
		}, "Generated uuid by hostname.")
	}

	if len(id) > 16 {
		return id[:16]
	}
	return id
}

func (e *Executor) GetHostname() string {
	if hostname, err := os.Hostname(); err == nil {
		return hostname
	}
	return "unknown"
}

func (e *Executor) GetHostUptimeAndIdle() (float64, float64) {
	e.logger.Info("Executor.GetHostUptime called.")
	defer e.logger.Info("Executor.GetHostUptime end.")

	var uptimeStr, idleStr string

	data, err := ioutil.ReadFile(uptimeFile)
	if err != nil {
		e.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while Executor.GetHostUptime.")
		return 0, 0
	}

	fields := strings.Split(strings.TrimSpace(string(data)), " ")
	if len(fields) == 2 {
		uptimeStr = fields[0]
		idleStr = fields[1]
	}

	uptime, _ := strconv.ParseFloat(uptimeStr, 64)
	idle, _ := strconv.ParseFloat(idleStr, 64)
	if uptime != 0 {
		idle = (idle / (uptime * float64(runtime.NumCPU()))) * 100
	}

	return uptime, idle
}

func (e *Executor) GetLocalIp(tcpAddr string) string {
	e.logger.Info("Executor.GetLocalIp called.")
	defer e.logger.Info("Executor.GetLocalIp end.")

	// 通过Socket获取本机IP
	conn, err := net.Dial("tcp", tcpAddr)
	if err != nil {
		e.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while Executor.GetLocalIp.")
		return ""
	}

	defer func() {
		_ = conn.Close()
	}()

	return strings.Split(conn.LocalAddr().String(), ":")[0]
}
