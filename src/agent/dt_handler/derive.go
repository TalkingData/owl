package dt_handler

import (
	"errors"
	"owl/dto"
)

// deriveHandler 处理derive，需要与上一次数据处理后才可发送
func deriveHandler(prevExist bool, curr, prev *dto.TsData) error {
	if !prevExist {
		return errors.New("previous ts data not exist")
	}

	curr.Value = curr.Value - prev.Value
	return nil
}
