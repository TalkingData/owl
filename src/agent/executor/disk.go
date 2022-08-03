package executor

import (
	"github.com/shirou/gopsutil/v3/disk"
	"owl/common/logger"
	"owl/dto"
	"time"
)

func (e *Executor) ExecCollectDisk(cycle int32) (res []*dto.TsData) {
	e.logger.Info("Executor.ExecCollectDisk called.")
	defer e.logger.Info("Executor.ExecCollectDisk end.")

	res = append(res, e.getAllDiskPartition(cycle)...)
	return append(res, e.getAllDiskIo(cycle)...)
}

func (e *Executor) getAllDiskIo(cycle int32) (res []*dto.TsData) {
	e.logger.Info("Executor.getAllDiskIo called.")
	defer e.logger.Info("Executor.getAllDiskIo end.")

	ioCounters, err := disk.IOCounters()
	if err != nil {
		e.logger.ErrorWithFields(logger.Fields{
			"cycle": cycle,
			"error": err,
		}, "An error occurred while Executor.getAllDiskIo.")
		return
	}

	currTs := time.Now().Unix()
	for _, v := range ioCounters {
		res = append(res,
			&dto.TsData{
				Metric:    "system.disk.bytes",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.ReadBytes),
				Timestamp: currTs,
				Cycle:     cycle,
				Tags:      map[string]string{"device": v.Name, "direction": "in"},
			},
			&dto.TsData{
				Metric:    "system.disk.count",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.ReadCount),
				Timestamp: currTs,
				Cycle:     cycle,
				Tags:      map[string]string{"device": v.Name, "direction": "in"},
			},
			&dto.TsData{
				Metric:    "system.disk.time",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.ReadTime),
				Timestamp: currTs,
				Cycle:     cycle,
				Tags:      map[string]string{"device": v.Name, "direction": "in"},
			},
			&dto.TsData{
				Metric:    "system.disk.bytes",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.WriteBytes),
				Timestamp: currTs,
				Cycle:     cycle,
				Tags:      map[string]string{"device": v.Name, "direction": "out"},
			},
			&dto.TsData{
				Metric:    "system.disk.count",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.WriteCount),
				Timestamp: currTs,
				Cycle:     cycle,
				Tags:      map[string]string{"device": v.Name, "direction": "out"},
			},
			&dto.TsData{
				Metric:    "system.disk.time",
				DataType:  dto.TsDataTypeCounter,
				Value:     float64(v.WriteTime),
				Timestamp: currTs,
				Cycle:     cycle,
				Tags:      map[string]string{"device": v.Name, "direction": "out"},
			},
		)
	}
	return
}

func (e *Executor) getAllDiskPartition(cycle int32) (res []*dto.TsData) {
	e.logger.Info("Executor.getAllDiskPartition called.")
	defer e.logger.Info("Executor.getAllDiskPartition end.")

	pts, err := disk.Partitions(false)
	if err != nil {
		return
	}

	currTd := time.Now().Unix()
	currTd = currTd - (currTd % int64(cycle))
	for _, p := range pts {
		usageStat, err := disk.Usage(p.Mountpoint)
		if err != nil {
			e.logger.ErrorWithFields(logger.Fields{
				"device":      p.Device,
				"mount_point": p.Mountpoint,
				"fs_type":     p.Fstype,
				"cycle":       cycle,
				"error":       err,
			}, "An error occurred while Executor.getAllDiskPartition, Skipped this one.")
			continue
		}
		res = append(res,
			&dto.TsData{
				Metric:    "system.disk.free",
				DataType:  dto.TsDataTypeGauge,
				Value:     float64(usageStat.Free),
				Timestamp: currTd,
				Cycle:     cycle,
				Tags:      map[string]string{"path": usageStat.Path, "fstype": usageStat.Fstype},
			},
			&dto.TsData{
				Metric:    "system.disk.total",
				DataType:  dto.TsDataTypeGauge,
				Value:     float64(usageStat.Total),
				Timestamp: currTd,
				Cycle:     cycle,
				Tags:      map[string]string{"path": usageStat.Path, "fstype": usageStat.Fstype},
			},
			&dto.TsData{
				Metric:    "system.disk.used",
				DataType:  dto.TsDataTypeGauge,
				Value:     float64(usageStat.Used),
				Timestamp: currTd,
				Cycle:     cycle,
				Tags:      map[string]string{"path": usageStat.Path, "fstype": usageStat.Fstype},
			},
			&dto.TsData{
				Metric:    "system.disk.used_pct",
				DataType:  dto.TsDataTypeGauge,
				Value:     usageStat.UsedPercent,
				Timestamp: currTd,
				Cycle:     cycle,
				Tags:      map[string]string{"path": usageStat.Path, "fstype": usageStat.Fstype},
			},
			&dto.TsData{
				Metric:    "system.disk.inodes.free",
				DataType:  dto.TsDataTypeGauge,
				Value:     float64(usageStat.InodesFree),
				Timestamp: currTd,
				Cycle:     cycle,
				Tags:      map[string]string{"path": usageStat.Path, "fstype": usageStat.Fstype},
			},
			&dto.TsData{
				Metric:    "system.disk.inodes.total",
				DataType:  dto.TsDataTypeGauge,
				Value:     float64(usageStat.InodesTotal),
				Timestamp: currTd,
				Cycle:     cycle,
				Tags:      map[string]string{"path": usageStat.Path, "fstype": usageStat.Fstype},
			},
			&dto.TsData{
				Metric:    "system.disk.inodes.used",
				DataType:  dto.TsDataTypeGauge,
				Value:     float64(usageStat.InodesUsed),
				Timestamp: currTd,
				Cycle:     cycle,
				Tags:      map[string]string{"path": usageStat.Path, "fstype": usageStat.Fstype},
			},
			&dto.TsData{
				Metric:    "system.disk.inodes.used_pct",
				DataType:  dto.TsDataTypeGauge,
				Value:     float64(usageStat.InodesUsedPercent),
				Timestamp: currTd,
				Cycle:     cycle,
				Tags:      map[string]string{"path": usageStat.Path, "fstype": usageStat.Fstype},
			},
		)
	}
	return
}
