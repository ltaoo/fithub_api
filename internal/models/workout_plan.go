package models

import "time"

// WorkoutPlan 表示一个训练计划
type WorkoutPlan struct {
	Id                int              `json:"id"`
	Title             string           `json:"title"`
	Overview          string           `json:"overview"`
	Level             int              `json:"level"`
	Tags              string           `json:"tags"`
	EstimatedDuration int              `json:"estimated_duration"`
	EquipmentIds      string           `json:"equipment_ids"`
	MuscleIds         string           `json:"muscle_ids"`
	Sets              []WorkoutPlanSet `json:"sets"`
	Details           string           `json:"details"`
	Points            string           `json:"points"`
	Suggestions       string           `json:"suggestions"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         *time.Time       `json:"updated_at"`
	OwnerId           int              `json:"owner_id"`
}

func (WorkoutPlan) TableName() string {
	return "WORKOUT_PLAN"
}

// WorkoutPlanSet 表示训练计划中的一组
type WorkoutPlanSet struct {
	Id            int                 `json:"id"`
	Name          string              `json:"name"`
	StepType      string              `json:"step_type"` // warmup、strength、stretch、cool_down、cardio、heart、performance
	StepIdx       int                 `json:"step_idx"`  // 第几个动作(阶段)
	SetType       string              `json:"set_type"`  // normal、combo、free
	SetIdx        int                 `json:"set_idx"`   // 组内第几组
	SetCount      int                 `json:"set_count"`
	Actions       []WorkoutPlanAction `json:"actions"`
	Note          string              `json:"note"`
	WorkoutPlanId int                 `json:"workout_plan_id"`
}

func (WorkoutPlanSet) TableName() string {
	return "WORKOUT_PLAN_SET"
}

// WorkoutPlanAction 表示训练计划中的动作要求
type WorkoutPlanAction struct {
	Id               int    `json:"id"`
	Idx              int    `json:"idx"`
	ActionId         int    `json:"action_id"`
	Reps             int    `json:"reps"`
	Unit             string `json:"unit"`
	Weight           string `json:"weight"`
	Tempo            string `json:"tempo"`
	RestInterval     int    `json:"rest_interval"`
	Note             string `json:"note"`
	WorkoutPlanSetId int    `json:"workout_plan_set_id"`
}

func (WorkoutPlanAction) TableName() string {
	return "WORKOUT_PLAN_ACTION"
}
