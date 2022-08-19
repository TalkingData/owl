package model

import "owl/common/utils"

type User struct {
	Id          uint             `json:"id"`
	Username    string           `json:"username"`
	DisplayName string           `json:"display_name"`
	Password    string           `json:"-"`
	Role        bool             `json:"role"`
	Phone       string           `json:"phone"`
	Email       string           `json:"email"`
	Wechat      string           `json:"wechat"`
	Type        string           `json:"type"`
	Status      bool             `json:"status"`
	CreateAt    *utils.LocalTime `json:"create_at" gorm:"autoCreateTime"`
	UpdateAt    *utils.LocalTime `json:"update_at" gorm:"autoUpdateTime"`
}

func (*User) TableName() string {
	return "user"
}
