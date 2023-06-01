package dt_handler

import (
	"errors"
	"owl/dto"
)

// growthHandler 处理growth，需要与上一次数据处理后才可发送
func growthHandler(prevExist bool, curr, prev *dto.TsData) error {
	if !prevExist {
		return errors.New("previous ts data not exist")
	}

	tmpVal := curr.Value - prev.Value
	if tmpVal < 0 {
		return errors.New("growth value is negative")
	}

	curr.Value = tmpVal
	return nil
}
