package main

import (
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/wuyingsong/utils"
)

const (
	uuid_file   = "/sys/devices/virtual/dmi/id/product_uuid"
	uptime_file = "/proc/uptime"
)

func getHostname() string {
	if hostname, err := os.Hostname(); err == nil {
		return hostname
	}
	return "unknown"
}

func getHostUptime() (float64, float64) {
	var (
		uptimeStr string
		idleStr   string
	)
	if data, err := ioutil.ReadFile(uptime_file); err == nil {
		fields := strings.Split(strings.TrimSpace(string(data)), " ")
		if len(fields) == 2 {
			uptimeStr = fields[0]
			idleStr = fields[1]
		}
	}
	uptime, _ := strconv.ParseFloat(uptimeStr, 64)
	idle, _ := strconv.ParseFloat(idleStr, 64)
	if uptime != 0 {
		idle = (idle / (uptime * float64(runtime.NumCPU()))) * 100
	}
	return uptime, idle
}

func getHostID() string {
	var id string
	if data, err := ioutil.ReadFile(uuid_file); err == nil {
		if len(data) > 0 {
			dataStr := strings.TrimSpace(string(data))
			id = utils.Md5(dataStr)
			lg.Info("generate id from %s:%s, result:%s", uuid_file, dataStr, id)
		}
	} else {
		id = utils.Md5(getHostname())
		lg.Info("generate id from hostname:%s, result:%s", getHostname(), id)
	}
	if len(id) > 16 {
		return id[:16]
	}
	return id
}
