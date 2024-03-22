package executor

import (
	"github.com/shirou/gopsutil/v3/disk"
	"owl/common/logger"
	"owl/dto"
)

func (e *Executor) ExecCollectDisk(ts int64, cycle int32) (res dto.TsDataArray) {
	e.logger.Info("Executor.ExecCollectDisk called.")
	defer e.logger.Info("Executor.ExecCollectDisk end.")

	res = append(res, e.getAllDiskPartition(ts, cycle)...)
	return append(res, e.getAllDiskIo(ts, cycle)...)
}

func (e *Executor) getAllDiskIo(ts int64, cycle int32) (res dto.TsDataArray) {
	e.logger.Info("Executor.getAllDiskIo called.")
	defer e.logger.Info("Executor.getAllDiskIo end.")

	ioCounters, err := disk.IOCounters()
	if err != nil {
		e.logger.ErrorWithFields(logger.Fields{
			"cycle": cycle,
			"error": err,
		}, "An error occurred while calling Executor.getAllDiskIo.")
		return
	}

	for _, v := range ioCounters {
		res = append(res,
			&dto.TsData{
				Metric:    "system.disk.bytes",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.ReadBytes),
				Timestamp: ts,
				Cycle:     cycle,
				Tags:      map[string]string{"device": v.Name, "direction": "in"},
			},
			&dto.TsData{
				Metric:    "system.disk.count",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.ReadCount),
				Timestamp: ts,
				Cycle:     cycle,
				Tags:      map[string]string{"device": v.Name, "direction": "in"},
			},
			&dto.TsData{
				Metric:    "system.disk.time",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.ReadTime),
				Timestamp: ts,
				Cycle:     cycle,
				Tags:      map[string]string{"device": v.Name, "direction": "in"},
			},
			&dto.TsData{
				Metric:    "system.disk.bytes",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.WriteBytes),
				Timestamp: ts,
				Cycle:     cycle,
				Tags:      map[string]string{"device": v.Name, "direction": "out"},
			},
			&dto.TsData{
				Metric:    "system.disk.count",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.WriteCount),
				Timestamp: ts,
				Cycle:     cycle,
				Tags:      map[string]string{"device": v.Name, "direction": "out"},
			},
			&dto.TsData{
				Metric:    "system.disk.time",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.WriteTime),
				Timestamp: ts,
				Cycle:     cycle,
				Tags:      map[string]string{"device": v.Name, "direction": "out"},
			},
		)
	}
	return
}

func (e *Executor) getAllDiskPartition(ts int64, cycle int32) (res dto.TsDataArray) {
	e.logger.Info("Executor.getAllDiskPartition called.")
	defer e.logger.Info("Executor.getAllDiskPartition end.")

	pts, err := disk.Partitions(false)
	if err != nil {
		return
	}

	for _, p := range pts {
		usageStat, err := disk.Usage(p.Mountpoint)
		if err != nil {
			e.logger.ErrorWithFields(logger.Fields{
				"device":      p.Device,
				"mount_point": p.Mountpoint,
				"fs_type":     p.Fstype,
				"cycle":       cycle,
				"error":       err,
			}, "An error occurred while calling Executor.getAllDiskPartition, Skipped this one.")
			continue
		}
		res = append(res,
			&dto.TsData{
				Metric:    "system.disk.free",
				DataType:  dto.TsDataTypeGauge,
				Value:     float64(usageStat.Free),
				Timestamp: ts,
				Cycle:     cycle,
				Tags:      map[string]string{"path": usageStat.Path, "fstype": usageStat.Fstype},
			},
			&dto.TsData{
				Metric:    "system.disk.total",
				DataType:  dto.TsDataTypeGauge,
				Value:     float64(usageStat.Total),
				Timestamp: ts,
				Cycle:     cycle,
				Tags:      map[string]string{"path": usageStat.Path, "fstype": usageStat.Fstype},
			},
			&dto.TsData{
				Metric:    "system.disk.used",
				DataType:  dto.TsDataTypeGauge,
				Value:     float64(usageStat.Used),
				Timestamp: ts,
				Cycle:     cycle,
				Tags:      map[string]string{"path": usageStat.Path, "fstype": usageStat.Fstype},
			},
			&dto.TsData{
				Metric:    "system.disk.used_pct",
				DataType:  dto.TsDataTypeGauge,
				Value:     usageStat.UsedPercent,
				Timestamp: ts,
				Cycle:     cycle,
				Tags:      map[string]string{"path": usageStat.Path, "fstype": usageStat.Fstype},
			},
			&dto.TsData{
				Metric:    "system.disk.inodes.free",
				DataType:  dto.TsDataTypeGauge,
				Value:     float64(usageStat.InodesFree),
				Timestamp: ts,
				Cycle:     cycle,
				Tags:      map[string]string{"path": usageStat.Path, "fstype": usageStat.Fstype},
			},
			&dto.TsData{
				Metric:    "system.disk.inodes.total",
				DataType:  dto.TsDataTypeGauge,
				Value:     float64(usageStat.InodesTotal),
				Timestamp: ts,
				Cycle:     cycle,
				Tags:      map[string]string{"path": usageStat.Path, "fstype": usageStat.Fstype},
			},
			&dto.TsData{
				Metric:    "system.disk.inodes.used",
				DataType:  dto.TsDataTypeGauge,
				Value:     float64(usageStat.InodesUsed),
				Timestamp: ts,
				Cycle:     cycle,
				Tags:      map[string]string{"path": usageStat.Path, "fstype": usageStat.Fstype},
			},
			&dto.TsData{
				Metric:    "system.disk.inodes.used_pct",
				DataType:  dto.TsDataTypeGauge,
				Value:     float64(usageStat.InodesUsedPercent),
				Timestamp: ts,
				Cycle:     cycle,
				Tags:      map[string]string{"path": usageStat.Path, "fstype": usageStat.Fstype},
			},
		)
	}
	return
}
