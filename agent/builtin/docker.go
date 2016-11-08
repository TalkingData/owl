package builtin

import (
	"github.com/google/cadvisor/client"
	"github.com/google/cadvisor/info/v1"
	"owl/common/types"
	"time"
)

func DockerMetrics(cycle int, cadvisor_addr string) []*types.TimeSeriesData {
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
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.cpu.system",
				Value:     float64(c.Cpu.Usage.System-p.Cpu.Usage.System) / timeNao,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.cpu.user",
				Value:     float64(c.Cpu.Usage.User-p.Cpu.Usage.User) / timeNao,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.cpu.total",
				Value:     float64(c.Cpu.Usage.Total-p.Cpu.Usage.Total) / timeNao,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.mem.cache",
				Value:     float64(c.Memory.Cache),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.mem.usage",
				Value:     float64(c.Memory.Usage),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.mem.rss",
				Value:     float64(c.Memory.RSS),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.mem.failcnt",
				Value:     float64(c.Memory.Failcnt),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.mem.workingset",
				Value:     float64(c.Memory.WorkingSet),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.readbytes",
				Value:     float64(c.Network.RxBytes-p.Network.RxBytes) / timeSec,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "iface": c.Network.Name, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.writebytes",
				Value:     float64(c.Network.TxBytes-p.Network.TxBytes) / timeSec,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "iface": c.Network.Name, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.readpackets",
				Value:     float64(c.Network.RxPackets-p.Network.RxPackets) / timeSec,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "iface": c.Network.Name, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.writepackets",
				Value:     float64(c.Network.TxPackets-c.Network.TxPackets) / timeSec,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "iface": c.Network.Name, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.readerrors",
				Value:     float64(c.Network.RxErrors-c.Network.RxErrors) / timeSec,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "iface": c.Network.Name, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.writeerrors",
				Value:     float64(c.Network.TxErrors-c.Network.TxErrors) / timeSec,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "iface": c.Network.Name, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.readdropped",
				Value:     float64(c.Network.RxDropped-c.Network.RxDropped) / timeSec,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "iface": c.Network.Name, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.writedropped",
				Value:     float64(c.Network.TxDropped-c.Network.TxDropped) / timeSec,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "iface": c.Network.Name, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.close",
				Value:     float64(c.Network.Tcp.Close),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.closewait",
				Value:     float64(c.Network.Tcp.CloseWait),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.closeing",
				Value:     float64(c.Network.Tcp.Closing),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.established",
				Value:     float64(c.Network.Tcp.Established),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.finwait1",
				Value:     float64(c.Network.Tcp.FinWait1),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.finwait2",
				Value:     float64(c.Network.Tcp.FinWait2),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.lastack",
				Value:     float64(c.Network.Tcp.LastAck),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.listen",
				Value:     float64(c.Network.Tcp.Listen),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.synrecv",
				Value:     float64(c.Network.Tcp.SynRecv),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.synsent",
				Value:     float64(c.Network.Tcp.SynSent),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
			&types.TimeSeriesData{
				Metric:    "docker.net.tcp.timewait",
				Value:     float64(c.Network.Tcp.TimeWait),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"name": cName, "id": id},
			},
		)
	}
	return metrics
}
