package types

import (
	"encoding/json"
	"fmt"
)

type Plugin struct {
	// plugin 在数据库中的id号
	ID int `json:"id"`
	// plugin 名称
	Name string `json:"name"`
	// plugin 在磁盘上的路径
	Path string `json:"path"`
	// 执行参数
	Args string `json:"args"`
	// 执行间隔
	Interval int `json:"interval"`
	// 执行超时时间
	Timeout int `json:"timeout"`
	// 校验和，用于插件同步比对
	Checksum string `json:"checksum"`
}

func (plugin *Plugin) Encode() []byte {
	data, _ := json.Marshal(plugin)
	return data
}

func (plugin *Plugin) Decode(data []byte) error {
	return json.Unmarshal(data, &plugin)
}

func (plugin Plugin) String() string {
	return fmt.Sprintf("{id:%d, name:%s, path:%s, args:%s, interval:%d, timeout:%d, checksum:%s}",
		plugin.ID,
		plugin.Name,
		plugin.Path,
		plugin.Args,
		plugin.Interval,
		plugin.Timeout,
		plugin.Checksum,
	)
}

func (plugin *Plugin) Equal(p Plugin) bool {
	if plugin.ID != p.ID ||
		plugin.Name != p.Name ||
		plugin.Path != p.Path ||
		plugin.Args != p.Args ||
		plugin.Interval != p.Interval ||
		plugin.Checksum != p.Checksum {
		return false
	}
	return true
}

func (Plugin) TableName() string {
	return "plugin"
}
