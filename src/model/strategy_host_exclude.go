package model

type StrategyHostExclude struct {
	Id         uint64 `json:"id"`
	StrategyId uint64 `json:"strategy_id"`
	HostId     string `json:"host_id"`
}

func (*StrategyHostExclude) TableName() string {
	return "strategy_host_exclude"
}
