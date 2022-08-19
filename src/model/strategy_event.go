package model

import "owl/common/utils"

const (
	StrategyEventStatusNew = iota + 1
	StrategyEventStatusAware
	StrategyEventStatusClosed
	StrategyEventStatusNoData
	StrategyEventStatusUnknown
)

type StrategyEvent struct {
	Id           uint64           `json:"id"`
	ProductId    uint             `json:"product_id"`
	StrategyId   uint64           `json:"strategy_id"`
	StrategyName string           `json:"strategy_name"`
	Priority     uint             `json:"priority"`
	Cycle        uint             `json:"cycle"`
	AlarmCount   uint             `json:"alarm_count"`
	Expression   string           `json:"expression"`
	CreateTime   *utils.LocalTime `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime   *utils.LocalTime `json:"update_time" gorm:"autoUpdateTime"`
	AwareEndTime *utils.LocalTime `json:"aware_end_time"`
	Count        uint             `json:"count"`
	Status       uint             `json:"status"`
	HostId       string           `json:"host_id"`
	HostName     string           `json:"host_name"` // 实际数据库内字段名为：host_name
	Ip           string           `json:"ip"`
}

func (*StrategyEvent) TableName() string {
	return "strategy_event"
}
