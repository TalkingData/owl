package model

import "owl/common/utils"

type StrategyEventRecord struct {
	StrategyEventId uint64           `json:"strategy_event_id"`
	Count           uint32           `json:"count" gorm:"type:uint"`
	StrategyId      uint64           `json:"strategy_id"`
	StrategyName    string           `json:"strategy_name"`
	Priority        uint32           `json:"priority" gorm:"type:uint"`
	Cycle           uint32           `json:"cycle" gorm:"type:uint"`
	AlarmCount      uint32           `json:"alarm_count" gorm:"type:uint"`
	Expression      string           `json:"expression"`
	CreateTime      *utils.LocalTime `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime      *utils.LocalTime `json:"update_time" gorm:"autoUpdateTime"`
	AwareEndTime    *utils.LocalTime `json:"aware_end_time"`
	Status          uint32           `json:"status" gorm:"type:uint"`
	HostId          string           `json:"host_id"`
	HostName        string           `json:"host_name"`
	Ip              string           `json:"ip"`
}

func (*StrategyEventRecord) TableName() string {
	return "strategy_event_record"
}
