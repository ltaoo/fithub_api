package models

// Muscle represents a muscle in the fitness system
type Equipment struct {
	Id       int64  `json:"id"`
	D        int    `json:"d"`
	Name     string `json:"name"`
	ZhName   string `json:"zh_name"`
	Alias    string `json:"alias"`
	Overview string `json:"overview"`
	Tags     string `json:"tags"`
	SortIdx  int    `json:"sort_idx"`
	Medias   string `json:"medias"`
}

func (Equipment) TableName() string {
	return "EQUIPMENT"
}
