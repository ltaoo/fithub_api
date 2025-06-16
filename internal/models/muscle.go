package models

// Muscle represents a muscle in the fitness system
type Muscle struct {
	Id       int64  `json:"id"`
	D        int    `json:"d"`
	Name     string `json:"name"`
	ZhName   string `json:"zh_name"`
	Tags     string `json:"tags"`
	Overview string `json:"overview"`
	SortIdx  int    `json:"sort_idx"`
	Features string `json:"features"`
	Medias   string `json:"medias"`
}

func (Muscle) TableName() string {
	return "MUSCLE"
}
