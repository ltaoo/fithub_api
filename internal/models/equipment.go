package models

// Muscle represents a muscle in the fitness system
type Equipment struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	ZhName   string `json:"zh_name"`
	Alias    string `json:"alias"`
	Overview string `json:"overview"`
	Medias   string `json:"medias"`
}

func (Equipment) TableName() string {
	return "EQUIPMENT"
}
