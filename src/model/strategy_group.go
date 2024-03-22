package model

type StrategyGroup struct {
	Id         uint64 `json:"id"`
	StrategyId uint64 `json:"strategy_id"`
	GroupId    uint32 `json:"group_id"`
}

func (*StrategyGroup) TableName() string {
	return "strategy_group"
}
