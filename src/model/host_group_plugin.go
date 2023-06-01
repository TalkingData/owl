package model

type HostGroupPlugin struct {
	Id       uint32 `json:"id"`
	GroupId  uint32 `json:"group_id"`
	PluginId uint32 `json:"plugin_id"`
	Args     string `json:"args"`
	Interval int32  `json:"interval"`
	Timeout  int32  `json:"timeout"`
	Comment  string `json:"comment"`
}

func (*HostGroupPlugin) TableName() string {
	return "host_group_plugin"
}
