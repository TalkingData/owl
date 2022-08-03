package model

import "owl/common/utils"

type Metric struct {
	Id       uint64           `json:"id"`
	HostId   string           `json:"host_id"`
	Metric   string           `json:"metric"`
	Tags     string           `json:"tags"`
	Dt       string           `json:"dt"`
	Cycle    int              `json:"cycle"`
	CreateAt *utils.LocalTime `json:"create_at" gorm:"autoCreateTime"`
	UpdateAt *utils.LocalTime `json:"update_at" gorm:"autoUpdateTime"`
}

func (*Metric) TableName() string {
	return "metric"
}
