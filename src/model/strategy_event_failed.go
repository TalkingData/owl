package model

import "owl/common/utils"

type StrategyEventFailed struct {
	StrategyId uint64           `json:"strategy_id"`
	HostId     string           `json:"host_id"`
	Status     int32            `json:"status"`
	Message    string           `json:"message"`
	CreateTime *utils.LocalTime `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime *utils.LocalTime `json:"update_time" gorm:"autoUpdateTime"`
}

func (*StrategyEventFailed) TableName() string {
	return "strategy_event_failed"
}
