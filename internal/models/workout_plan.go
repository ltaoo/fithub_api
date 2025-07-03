package models

import "time"

// WorkoutPlan 表示一个训练计划
type WorkoutPlan struct {
	Id                int        `json:"id"`
	D                 int        `json:"d"`
	Status            int        `json:"status" gorm:"default:2"` // 2仅自己可见
	Title             string     `json:"title"`
	Overview          string     `json:"overview"`
	CoverURL          string     `json:"cover_url"`
	Type              string     `json:"type"` //
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

	OwnerId int   `json:"owner_id"`
	Creator Coach `json:"creator" gorm:"foreignKey:OwnerId"`
}

func (WorkoutPlan) TableName() string {
	return "WORKOUT_PLAN"
}

type WorkoutPlanStatus int

const (
	// 0等待进行
	WorkoutPublishStatusPending WorkoutPlanStatus = iota
	// 1公开
	WorkoutPublishStatusPublic
	// 2仅自己可见
	WorkoutPublishStatusPrivate
	// 3禁用
	WorkoutPublishStatusDisabled
)

type WorkoutSchedule struct {
	Id        int        `json:"id"`
	D         int        `json:"d"`
	Title     string     `json:"title"`
	Overview  string     `json:"overview"`
	Status    int        `json:"status"`
	Level     int        `json:"level"`
	Type      int        `json:"type"`
	Details   string     `json:"details"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	WorkoutPlans []WorkoutPlanInSchedule `json:"workout_plans" gorm:"foreignKey:WorkoutPlanCollectionId"`

	OwnerId int   `json:"owner_id"`
	Creator Coach `json:"creator" gorm:"foreignKey:OwnerId"`
}

func (WorkoutSchedule) TableName() string {
	return "WORKOUT_PLAN_COLLECTION"
}

type WorkoutPlanInSchedule struct {
	Weekday int `json:"weekday"`
	Day     int `json:"day"`
	Idx     int `json:"idx"`

	WorkoutPlanId           int         `json:"workout_plan_id"`
	WorkoutPlan             WorkoutPlan `json:"workout_plan"`
	WorkoutPlanCollectionId int         `json:"workout_plan_collection_id"`
	// WorkoutPlanCollection   WorkoutSchedule `json:"workout_schedule" gorm:"foreignKey:WorkoutPlanCollectionId"`
}

func (WorkoutPlanInSchedule) TableName() string {
	return "WORKOUT_PLAN_IN_COLLECTION"
}

type CoachWorkoutSchedule struct {
	Id          int        `json:"id"`
	D           int        `json:"d"`
	Interval    int        `json:"interval"`
	Status      int        `json:"status"`
	StartDate   *time.Time `json:"start_date"`
	AppliedAt   time.Time  `json:"applied_at"`
	CancelledAt *time.Time `json:"cancelled_at"`

	WorkoutPlanCollectionId int             `json:"workout_plan_collection_id"`
	WorkoutPlanCollection   WorkoutSchedule `json:"workout_plan_collection"`
	CoachId                 int             `json:"coach_id"`
	Coach                   Coach           `json:"coach"`
}

func (CoachWorkoutSchedule) TableName() string {
	return "COACH_WORKOUT_PLAN_COLLECTION"
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

type WorkoutPlanSet struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Overview  string    `json:"overview"`
	IconURL   string    `json:"icon_url"`
	Idx       int       `json:"idx"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

func (WorkoutPlanSet) TableName() string {
	return "WORKOUT_PLAN_SET"
}
