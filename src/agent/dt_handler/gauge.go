package dt_handler

import "owl/dto"

// gaugeHandler 不需要任何处理，直接跳过
func gaugeHandler(_ bool, _, _ *dto.TsData) error {
	return nil
}
