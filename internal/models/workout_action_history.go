package models

import "time"

type WorkoutActionHistory struct {
	Id          int       `json:"id"`
	Reps        int       `json:"reps"`
	RepsUnit    string    `json:"reps_unit"`
	Weight      int       `json:"weight"`
	WeightUnit  string    `json:"weight_unit"`
	Remark      string    `json:"remark"`
	ExtraMedias string    `json:"extra_medias"`
	CreatedAt   time.Time `json:"created_at"`

	WorkoutDayId    int           `json:"workout_day_id"`
	StudentId       int           `json:"student_id"`
	WorkoutActionId int           `json:"action_id" gorm:"column:action_id"`
	WorkoutAction   WorkoutAction `json:"action" gorm:"foreignKey:WorkoutActionId"`
}

func (WorkoutActionHistory) TableName() string {
	return "WORKOUT_ACTION_HISTORY"
}
