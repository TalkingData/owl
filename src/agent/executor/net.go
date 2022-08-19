package executor

import (
	"github.com/shirou/gopsutil/v3/net"
	"owl/common/logger"
	"owl/dto"
	"time"
)

func (e *Executor) ExecCollectNet(cycle int32) (res dto.TsDataArray) {
	e.logger.Info("Executor.ExecCollectNet called.")
	defer e.logger.Info("Executor.ExecCollectNet end.")

	ioCounters, err := net.IOCounters(true)
	if err != nil {
		e.logger.ErrorWithFields(logger.Fields{
			"cycle": cycle,
			"error": err,
		}, "An error occurred while Executor.ExecCollectNet.")
		return
	}

	currTs := time.Now().Unix()
	for _, v := range ioCounters {
		res = append(res,
			&dto.TsData{
				Metric:    "system.net.bytes",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.BytesSent),
				Timestamp: currTs,
				Cycle:     cycle,
				Tags:      map[string]string{"iface": v.Name, "direction": "out"},
			},
			&dto.TsData{
				Metric:    "system.net.bytes",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.BytesRecv),
				Timestamp: currTs,
				Cycle:     cycle,
				Tags:      map[string]string{"iface": v.Name, "direction": "in"},
			},
			&dto.TsData{
				Metric:    "system.net.packets",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.PacketsSent),
				Timestamp: currTs,
				Cycle:     cycle,
				Tags:      map[string]string{"iface": v.Name, "direction": "out"},
			},
			&dto.TsData{
				Metric:    "system.net.packets",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.PacketsRecv),
				Timestamp: currTs,
				Cycle:     cycle,
				Tags:      map[string]string{"iface": v.Name, "direction": "in"},
			},
			&dto.TsData{
				Metric:    "system.net.err",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.Errin),
				Timestamp: currTs,
				Cycle:     cycle,
				Tags:      map[string]string{"iface": v.Name, "direction": "in"},
			},
			&dto.TsData{
				Metric:    "system.net.err",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.Errout),
				Timestamp: currTs,
				Cycle:     cycle,
				Tags:      map[string]string{"iface": v.Name, "direction": "out"},
			},
			&dto.TsData{
				Metric:    "system.net.drop",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.Dropin),
				Timestamp: currTs,
				Cycle:     cycle,
				Tags:      map[string]string{"iface": v.Name, "direction": "in"},
			},
			&dto.TsData{
				Metric:    "system.net.drop",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.Dropout),
				Timestamp: currTs,
				Cycle:     cycle,
				Tags:      map[string]string{"iface": v.Name, "direction": "out"},
			},
			&dto.TsData{
				Metric:    "system.net.fifo",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.Fifoin),
				Timestamp: currTs,
				Cycle:     cycle,
				Tags:      map[string]string{"iface": v.Name, "direction": "in"},
			},
			&dto.TsData{
				Metric:    "system.net.fifo",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.Fifoout),
				Timestamp: currTs,
				Cycle:     cycle,
				Tags:      map[string]string{"iface": v.Name, "direction": "out"},
			},
		)
	}
	return
}
