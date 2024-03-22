package model

type Strategy struct {
	Id          uint64 `json:"id"`
	ProductId   uint32 `json:"product_id"`
	Name        string `json:"name"`
	Priority    uint32 `json:"priority"`
	AlarmCount  uint32 `json:"alarm_count"`
	Cycle       uint32 `json:"cycle"`
	Expression  string `json:"expression"`
	Description string `json:"description"`
	Enable      bool   `json:"enable"`
	UserId      uint32 `json:"user_id"`
}

func (*Strategy) TableName() string {
	return "strategy"
}
