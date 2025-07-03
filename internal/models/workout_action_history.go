package models

import "time"

type WorkoutActionHistory struct {
	Id          int       `json:"id"`
	D           int       `json:"d" gorm:"column:d;default:0"`
	StepUid     int       `json:"step_uid"`
	SetUid      int       `json:"set_uid"`
	ActUid      int       `json:"act_uid"`
	Reps        int       `json:"reps"`
	RepsUnit    string    `json:"reps_unit"`
	Weight      float64   `json:"weight"`
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
