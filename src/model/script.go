package model

type Script struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	FilePath string `json:"file_path"`
}

func (*Script) TableName() string {
	return "script"
}
