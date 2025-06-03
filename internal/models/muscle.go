package models

// Muscle represents a muscle in the fitness system
type Muscle struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	ZhName   string `json:"zh_name"`
	Tags     string `json:"tags"`
	Overview string `json:"overview"`
	Features string `json:"features"`
	Pics     string `json:"pics"`
}

func (Muscle) TableName() string {
	return "MUSCLE"
}
