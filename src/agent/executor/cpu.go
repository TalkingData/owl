package executor

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"owl/dto"
	"time"
)

func (e *Executor) ExecCollectCpu(ts int64, cycle int32) (res dto.TsDataArray) {
	e.logger.Info("Executor.ExecCollectCpu called.")
	defer e.logger.Info("Executor.ExecCollectCpu end.")

	cts1 := getAllCpuTimesStat()
	time.Sleep(time.Second * 1)
	cts2 := getAllCpuTimesStat()

	for idx, data2 := range cts2 {
		res = append(res, processCpuTimeSeriesData(ts, cycle, &cts1[idx], &data2)...)
	}

	return
}

func getAllCpuTimesStat() []cpu.TimesStat {
	cts1, _ := cpu.Times(true)
	cts2, _ := cpu.Times(false)
	return append(cts1, cts2...)
}

func processCpuTimeSeriesData(currTs int64, cycle int32, data1, data2 *cpu.TimesStat) dto.TsDataArray {
	total := data2.Total() - data1.Total()
	return dto.TsDataArray{
		{
			Metric:    "system.cpu.user",
			DataType:  dto.TsDataTypeGauge,
			Value:     (data2.User - data1.User) / total * 100,
			Timestamp: currTs,
			Cycle:     cycle,
			Tags:      map[string]string{"cpu": data2.CPU},
		},
		{
			Metric:    "system.cpu.system",
			DataType:  dto.TsDataTypeGauge,
			Value:     (data2.System - data1.System) / total * 100,
			Timestamp: currTs,
			Cycle:     cycle,
			Tags:      map[string]string{"cpu": data2.CPU},
		},
		{
			Metric:    "system.cpu.idle",
			DataType:  dto.TsDataTypeGauge,
			Value:     (data2.Idle - data1.Idle) / total * 100,
			Timestamp: currTs,
			Cycle:     cycle,
			Tags:      map[string]string{"cpu": data2.CPU},
		},
		{
			Metric:    "system.cpu.nice",
			DataType:  dto.TsDataTypeGauge,
			Value:     (data2.Nice - data1.Nice) / total * 100,
			Timestamp: currTs,
			Cycle:     cycle,
			Tags:      map[string]string{"cpu": data2.CPU},
		},
		{
			Metric:    "system.cpu.iowait",
			DataType:  dto.TsDataTypeGauge,
			Value:     (data2.Iowait - data1.Iowait) / total * 100,
			Timestamp: currTs,
			Cycle:     cycle,
			Tags:      map[string]string{"cpu": data2.CPU},
		},
		{
			Metric:    "system.cpu.irq",
			DataType:  dto.TsDataTypeGauge,
			Value:     (data2.Irq - data1.Irq) / total * 100,
			Timestamp: currTs,
			Cycle:     cycle,
			Tags:      map[string]string{"cpu": data2.CPU},
		},
		{
			Metric:    "system.cpu.softirq",
			DataType:  dto.TsDataTypeGauge,
			Value:     (data2.Softirq - data1.Softirq) / total * 100,
			Timestamp: currTs,
			Cycle:     cycle,
			Tags:      map[string]string{"cpu": data2.CPU},
		},
		{
			Metric:    "system.cpu.steal",
			DataType:  dto.TsDataTypeGauge,
			Value:     (data2.Steal - data1.Steal) / total * 100,
			Timestamp: currTs,
			Cycle:     cycle,
			Tags:      map[string]string{"cpu": data2.CPU},
		},
		{
			Metric:    "system.cpu.guest",
			DataType:  dto.TsDataTypeGauge,
			Value:     (data2.Guest - data1.Guest) / total * 100,
			Timestamp: currTs,
			Cycle:     cycle,
			Tags:      map[string]string{"cpu": data2.CPU},
		},
		{
			Metric:    "system.cpu.guestnice",
			DataType:  dto.TsDataTypeGauge,
			Value:     (data2.GuestNice - data1.GuestNice) / total * 100,
			Timestamp: currTs,
			Cycle:     cycle,
			Tags:      map[string]string{"cpu": data2.CPU},
		},
	}
}
