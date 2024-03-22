package model

type TriggerEvent struct {
	StrategyEventId  uint64  `json:"strategy_event_id"`
	Index            string  `json:"index"`
	Metric           string  `json:"metric"`
	Tags             string  `json:"tags"`
	Number           uint32  `json:"number"`
	AggregateTags    string  `json:"aggregate_tags"`
	CurrentThreshold float64 `json:"current_threshold"`
	Method           string  `json:"method"`
	Symbol           string  `json:"symbol"`
	Threshold        float64 `json:"threshold"`
	Triggered        bool    `json:"triggered"`
}

func (*TriggerEvent) TableName() string {
	return "trigger_event"
}
