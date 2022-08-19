package model

type Action struct {
	Id              uint   `json:"id"`
	StrategyId      uint64 `json:"strategy_id"`
	Type            uint   `json:"type"`
	Kind            uint   `json:"kind"`
	AlarmSubject    string `json:"alarm_subject"`
	AlarmTemplate   string `json:"alarm_template"`
	RestoreSubject  string `json:"restore_subject"`
	RestoreTemplate string `json:"restore_template"`
	ScriptId        uint   `json:"script_id"`
	BeginTime       string `json:"begin_time"`
	EndTime         string `json:"end_time"`
	TimePeriod      int    `json:"time_period"`
}

func (*Action) TableName() string {
	return "action"
}
