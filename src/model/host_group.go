package model

import "owl/common/utils"

type HostGroup struct {
	Id          uint32           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	ProductId   uint32           `json:"product_id"`
	Creator     string           `json:"creator"`
	CreateAt    *utils.LocalTime `json:"create_at" gorm:"autoCreateTime"`
	UpdateAt    *utils.LocalTime `json:"update_at" gorm:"autoUpdateTime"`
}

func (*HostGroup) TableName() string {
	return "host_group"
}
