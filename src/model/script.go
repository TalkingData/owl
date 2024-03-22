package model

type Script struct {
	Id       uint32 `json:"id"`
	Name     string `json:"name"`
	FilePath string `json:"file_path"`
}

func (*Script) TableName() string {
	return "script"
}
