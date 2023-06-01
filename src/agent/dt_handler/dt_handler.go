package dt_handler

import "owl/dto"

type DtHandlerFunc func(prevExist bool, curr, prev *dto.TsData) error

type DtHandlerMap map[string]DtHandlerFunc

// NewDtHandlerMap 创建一个新的DtHandlerMap
func NewDtHandlerMap() DtHandlerMap {
	return DtHandlerMap{
		dto.TsDataTypeGauge:   gaugeHandler,
		dto.TsDataTypeCounter: counterHandler,
		dto.TsDataTypeDerive:  deriveHandler,
		dto.TsDataTypeGrowth:  growthHandler,
	}
}

func (dtHMap DtHandlerMap) Get(dataType string) DtHandlerFunc {
	hFunc, ok := dtHMap[dataType]
	if !ok {
		return gaugeHandler
	}

	return hFunc
}
