package model

type Trigger struct {
	Id          uint64  `json:"id"`
	StrategyId  uint64  `json:"strategy_id"`
	Metric      string  `json:"metric"`
	Tags        string  `json:"tags"`
	Number      uint32  `json:"number"`
	Index       string  `json:"index"`
	Method      string  `json:"method"`
	Symbol      string  `json:"symbol"`
	Threshold   float64 `json:"threshold"`
	Description string  `json:"description"`
}

func (*Trigger) TableName() string {
	return "trigger"
}
