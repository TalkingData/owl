package model

type Action struct {
	Id              uint32 `json:"id"`
	StrategyId      uint64 `json:"strategy_id"`
	Type            uint32 `json:"type"`
	Kind            int32  `json:"kind"`
	AlarmSubject    string `json:"alarm_subject"`
	AlarmTemplate   string `json:"alarm_template"`
	RestoreSubject  string `json:"restore_subject"`
	RestoreTemplate string `json:"restore_template"`
	ScriptId        uint32 `json:"script_id"`
	BeginTime       string `json:"begin_time"`
	EndTime         string `json:"end_time"`
	TimePeriod      int32  `json:"time_period"`
}

func (*Action) TableName() string {
	return "action"
}
