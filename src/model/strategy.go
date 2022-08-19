package model

type Strategy struct {
	Id          uint64 `json:"id"`
	ProductId   uint   `json:"product_id"`
	Name        string `json:"name"`
	Priority    uint   `json:"priority"`
	AlarmCount  uint   `json:"alarm_count"`
	Cycle       uint   `json:"cycle"`
	Expression  string `json:"expression"`
	Description string `json:"description"`
	Enable      bool   `json:"enable"`
	UserId      uint   `json:"user_id"`
}

func (*Strategy) TableName() string {
	return "strategy"
}
