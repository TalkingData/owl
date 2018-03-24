package types

type Trigger struct {
	ID          int64   `form:"id" json:"-" `
	StrategyID  int     `form:"strategy_id" json:"strategy_id"`
	Metric      string  `form:"metric" json:"metric"`
	Tags        string  `form:"tags" json:"tags"`
	Number      int     `form:"number" json:"number"`
	Index       string  `form:"index" json:"index"`
	Method      string  `form:"method" json:"method" `
	Symbol      string  `form:"symbol" json:"symbol" `
	Threshold   float64 `form:"threshold" json:"threshold" `
	Description string  `form:"description" json:"description"`
}

func (Trigger) TableName() string {
	return "trigger"
}
