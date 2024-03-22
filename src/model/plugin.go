package model

import (
	"fmt"
	"owl/common/utils"
	"sort"
	"strings"
)

type Plugin struct {
	Id       uint32           `json:"id"`
	Name     string           `json:"name"`
	Path     string           `json:"path"`
	Args     string           `json:"args"`
	Interval int32            `json:"interval"`
	Timeout  int32            `json:"timeout"`
	Checksum string           `json:"checksum"`
	CreateAt *utils.LocalTime `json:"create_at" gorm:"autoCreateTime"`
	UpdateAt *utils.LocalTime `json:"update_at" gorm:"autoUpdateTime"`
	Creator  string           `json:"creator"`
	Comment  string           `json:"comment"`
}

func (plugin *Plugin) GenUniqueKey() string {
	argSlice := utils.ParseCommandArgs(
		fmt.Sprintf("%s %s", plugin.Path, plugin.Args),
	)
	sort.Strings(argSlice)
	uniqueKey := strings.Join(argSlice, "")
	return utils.Md5(uniqueKey)
}

func (*Plugin) TableName() string {
	return "plugin"
}
