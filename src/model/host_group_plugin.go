package model

type HostGroupPlugin struct {
	Id       uint   `json:"id"`
	GroupId  uint   `json:"group_id"`
	PluginId uint   `json:"plugin_id"`
	Args     string `json:"args"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
	Comment  string `json:"comment"`
}

func (*HostGroupPlugin) TableName() string {
	return "host_group_plugin"
}
