package dt_handler

import (
	"errors"
	"owl/dto"
)

// counterHandler 处理counter，需要与上一次数据处理后才可发送
func counterHandler(prevExist bool, curr, prev *dto.TsData) error {
	if !prevExist {
		return errors.New("previous ts data not exist")
	}

	if curr.Cycle == 0 {
		return errors.New("ts data cycle is 0.")
	}

	rate := (curr.Value - prev.Value) / float64(curr.Cycle)
	if rate < 0 {
		rate = 0
	}

	curr.Value = rate
	return nil
}
