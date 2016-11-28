package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"owl/common/types"
	"time"

	"github.com/google/cadvisor/client"
	"github.com/google/cadvisor/info/v1"
)

func main() {
	CadvisorAddr := flag.String("cadvisor", "http://127.0.0.1:8080", "cadvisor api address")
	flag.Parse()
	b, err := json.MarshalIndent(DockerMetrics(*CadvisorAddr), "", "  ")
	res := string(b)
	if err != nil || res == "null" {
		fmt.Println("[]")
	} else {
		fmt.Println(string(res))
	}
}

func DockerMetrics(cadvisor_addr string) []*types.TimeSeriesData {
	client, err := client.NewClient(cadvisor_addr)
	if err != nil {
		return nil
	}
	request := v1.ContainerInfoRequest{NumStats: 2}
	info, err := client.AllDockerContainers(&request)
	if err != nil {
		return nil
	}
	metrics := []*types.TimeSeriesData{}
	var (
		cName string = "unknown"
	)
	ts := time.Now().Unix()
	for _, i := range info {
		if len(i.Aliases) > 0 {
			cName = i.Aliases[0]
		}
		id := i.Id[:12]
		c := i.Stats[1]
		p := i.Stats[0]
		timeNao := float64(c.Timestamp.Sub(p.Timestamp).Nanoseconds())
		timeSec := c.Timestamp.Sub(p.Timestamp).Seconds()
		metrics = append(metrics,
			&types.TimeSeriesData{
				Metric:    "docker.cpu.load",
				Value:     float64(c.Cpu.LoadAverage),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.cpu.system",
				Value:     float64(c.Cpu.Usage.System-p.Cpu.Usage.System) / timeNao,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.cpu.user",
				Value:     float64(c.Cpu.Usage.User-p.Cpu.Usage.User) / timeNao,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.cpu.total",
				Value:     float64(c.Cpu.Usage.Total-p.Cpu.Usage.Total) / timeNao,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.mem.cache",
				Value:     float64(c.Memory.Cache),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.mem.usage",
				Value:     float64(c.Memory.Usage),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.mem.rss",
				Value:     float64(c.Memory.RSS),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.mem.failcnt",
				Value:     float64(c.Memory.Failcnt),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.mem.workingset",
				Value:     float64(c.Memory.WorkingSet),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.RxBytes",
				Value:     float64(c.Network.RxBytes-p.Network.RxBytes) / timeSec,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "iface": c.Network.Name, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.TxBytes",
				Value:     float64(c.Network.TxBytes-p.Network.TxBytes) / timeSec,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "iface": c.Network.Name, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.RxPackets",
				Value:     float64(c.Network.RxPackets-p.Network.RxPackets) / timeSec,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "iface": c.Network.Name, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.TxPackets",
				Value:     float64(c.Network.TxPackets-c.Network.TxPackets) / timeSec,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "iface": c.Network.Name, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.RxErrors",
				Value:     float64(c.Network.RxErrors-c.Network.RxErrors) / timeSec,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "iface": c.Network.Name, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.TxErrors",
				Value:     float64(c.Network.TxErrors-c.Network.TxErrors) / timeSec,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "iface": c.Network.Name, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.RxDropped",
				Value:     float64(c.Network.RxDropped-c.Network.RxDropped) / timeSec,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "iface": c.Network.Name, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.TxDropped",
				Value:     float64(c.Network.TxDropped-c.Network.TxDropped) / timeSec,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "iface": c.Network.Name, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.Close",
				Value:     float64(c.Network.Tcp.Close),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.CloseWait",
				Value:     float64(c.Network.Tcp.CloseWait),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.Closing",
				Value:     float64(c.Network.Tcp.Closing),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.Established",
				Value:     float64(c.Network.Tcp.Established),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.FinWait1",
				Value:     float64(c.Network.Tcp.FinWait1),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.FinWait2",
				Value:     float64(c.Network.Tcp.FinWait2),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.LastAck",
				Value:     float64(c.Network.Tcp.LastAck),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.Listen",
				Value:     float64(c.Network.Tcp.Listen),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.SynRecv",
				Value:     float64(c.Network.Tcp.SynRecv),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.SynSent",
				Value:     float64(c.Network.Tcp.SynSent),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.TimeWait",
				Value:     float64(c.Network.Tcp.TimeWait),
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
		)
	}
	return metrics
}
