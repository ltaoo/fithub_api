package models

import "time"

// WorkoutPlan 表示一个训练计划
type WorkoutPlan struct {
	Id                int        `json:"id"`
	Status            int        `json:"status" gorm:"default:2"` // 2仅自己可见
	Title             string     `json:"title"`
	Overview          string     `json:"overview"`
	Level             int        `json:"level" gorm:"default:1"`
	Tags              string     `json:"tags"`
	EstimatedDuration int        `json:"estimated_duration" gorm:"default:60"`
	EquipmentIds      string     `json:"equipment_ids"`
	MuscleIds         string     `json:"muscle_ids"`
	Details           string     `json:"details"`
	Points            string     `json:"points"`
	Suggestions       string     `json:"suggestions"`
	CreatedAt         time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt         *time.Time `json:"updated_at"`
	OwnerId           int        `json:"owner_id"`

	Steps []WorkoutPlanStep `json:"steps" gorm:"foreignKey:WorkoutPlanId"`
}

func (WorkoutPlan) TableName() string {
	return "WORKOUT_PLAN"
}

// WorkoutPlanSet 表示训练计划中的一组
type WorkoutPlanStep struct {
	Id              int    `json:"id"`
	Title           string `json:"title"`
	Type            string `json:"type"` // warmup、strength、stretch、cool_down、cardio、heart、performance
	Idx             int    `json:"idx"`  // 第几个动作(阶段)
	SetType         string `json:"set_type"`
	SetCount        int    `json:"set_count"`
	SetRestDuration int    `json:"set_rest_duration"`
	SetWeight       string `json:"set_weight"`
	Note            string `json:"note"`
	WorkoutPlanId   int    `json:"workout_plan_id"`

	Actions []WorkoutPlanAction `json:"actions" gorm:"foreignKey:WorkoutPlanStepId"`
}

func (WorkoutPlanStep) TableName() string {
	return "WORKOUT_PLAN_STEP"
}

// WorkoutPlanAction 表示训练计划中的动作要求
type WorkoutPlanAction struct {
	Id                int    `json:"id"`
	Idx               int    `json:"idx"`
	ActionId          int    `json:"action_id"`
	SetIdx            int    `json:"set_idx"`
	Reps              int    `json:"reps" gorm:"default:12"`
	RepsUnit          string `json:"reps_unit" gorm:"default:'次'"`
	Weight            string `json:"weight"`
	Tempo             string `json:"tempo" gorm:"default:'4/1/2'"`
	RestDuration      int    `json:"rest_duration" gorm:"default:90"`
	Note              string `json:"note"`
	WorkoutPlanStepId int    `json:"workout_plan_step_id"`

	Action WorkoutAction `json:"action" gorm:"foreignKey:ActionId"`
}

func (WorkoutPlanAction) TableName() string {
	return "WORKOUT_PLAN_ACTION"
}
