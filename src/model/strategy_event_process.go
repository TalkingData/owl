package model

import "owl/common/utils"

type StrategyEventProcess struct {
	StrategyEventId uint64           `json:"strategy_event_id"`
	ProcessStatus   int32            `json:"process_status"`
	ProcessUser     string           `json:"process_user"`
	ProcessComments string           `json:"process_comments"`
	ProcessTime     *utils.LocalTime `json:"process_time" gorm:"autoCreateTime"`
}

func (*StrategyEventProcess) TableName() string {
	return "strategy_event_process"
}
