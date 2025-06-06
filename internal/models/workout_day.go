package models

import (
	"time"
)

// WorkoutDay represents a training day record
type WorkoutDay struct {
	Id                int        `json:"id" db:"id"`                                 // Primary key
	Time              string     `json:"time" db:"time"`                             // Training time (YYYY-MM-DD HH:MM:SS)
	Status            int        `json:"status" db:"status"`                         // Status: 1=Pending 2=In Progress 3=Completed 4=Expired 5=Cancelled
	EstimatedDuration int        `json:"estimated_duration" db:"estimated_duration"` // Estimated duration
	PendingSteps      string     `json:"pending_steps" db:"pending_steps"`           // Execution records in JSON array
	UpdatedDetails    string     `json:"updated_details" db:"updated_details"`       // Execution records in JSON array
	Stats             string     `json:"stats" db:"stats"`                           // Statistics in JSON array
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`                 // Creation time
	StartedAt         *time.Time `json:"started_at,omitempty" db:"started_at"`       // Start time
	UpdatedAt         *time.Time `json:"updated_at,omitempty" db:"updated_at"`       // Update time
	FinishedAt        *time.Time `json:"finished_at,omitempty" db:"finished_at"`     // Finish time
	CoachId           int        `json:"coach_id" db:"coach_id" `                    // Coach ID

	WorkoutPlanId int         `json:"workout_plan_id" db:"workout_plan_id"` // Associated workout plan ID
	WorkoutPlan   WorkoutPlan `json:"workout_plan" gorm:"foreignKey:WorkoutPlanId"`
	StudentId     int         `json:"student_id" db:"student_id"`
	Student       Coach       `json:"student" gorm:"foreignKey:StudentId"`

	ActionHistories []WorkoutActionHistory `json:"action_histories" gorm:"foreignKey:WorkoutDayId"`
}

func (WorkoutDay) TableName() string {
	return "WORKOUT_DAY"
}

type WorkoutDayStatus int

const (
	// 0等待进行
	WorkoutDayStatusPending WorkoutDayStatus = iota
	// 1进行中
	WorkoutDayStatusStarted
	// 2已完成
	WorkoutDayStatusFinished
	// 3已过期
	WorkoutDayStatusExpired
	// 4手动作废
	WorkoutDayStatusCancelled
	// 5进行中又放弃了
	WorkoutDayStatusGiveUp
)
