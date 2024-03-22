package model

type HostPlugin struct {
	Id       uint32 `json:"id"`
	HostId   string `json:"host_id"`
	PluginId uint32 `json:"plugin_id"`
	Args     string `json:"args"`
	Interval int32  `json:"interval"`
	Timeout  int32  `json:"timeout"`
	Comment  string `json:"comment"`
}

func (*HostPlugin) TableName() string {
	return "host_plugin"
}
