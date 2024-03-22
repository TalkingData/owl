package model

import "owl/common/utils"

type Product struct {
	Id          uint32           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Creator     string           `json:"creator"`
	CreateAt    *utils.LocalTime `json:"create_at" gorm:"autoCreateTime"`
	IsDelete    bool             `json:"is_delete"`
}

func (*Product) TableName() string {
	return "product"
}
