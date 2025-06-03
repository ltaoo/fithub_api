package models

import (
	"time"
)

// WorkoutDay represents a training day record
type WorkoutDay struct {
	Id                int64      `db:"id" json:"id"`                                 // Primary key
	Time              string     `db:"time" json:"time"`                             // Training time (YYYY-MM-DD HH:MM:SS)
	Status            int        `db:"status" json:"status"`                         // Status: 1=Pending 2=In Progress 3=Completed 4=Expired 5=Cancelled
	EstimatedDuration int        `db:"estimated_duration" json:"estimated_duration"` // Estimated duration
	PendingSteps      string     `db:"pending_steps" json:"pending_steps"`           // Execution records in JSON array
	UpdatedDetails    string     `db:"updated_details" json:"updated_details"`       // Execution records in JSON array
	Stats             string     `db:"stats" json:"stats"`                           // Statistics in JSON array
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`                 // Creation time
	StartedAt         *time.Time `db:"started_at" json:"started_at,omitempty"`       // Start time
	UpdatedAt         *time.Time `db:"updated_at" json:"updated_at,omitempty"`       // Update time
	FinishedAt        *time.Time `db:"finished_at" json:"finished_at,omitempty"`     // Finish time
	CoachId           int64      `db:"coach_id" json:"coach_id"`                     // Coach ID
	StudentId         int64      `db:"student_id" json:"student_id"`                 // Student ID
	WorkoutPlanId     int64      `db:"workout_plan_id" json:"workout_plan_id"`       // Associated workout plan ID

	WorkoutPlan     WorkoutPlan            `json:"workout_plan" gorm:"foreignKey:WorkoutPlanId"`
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
