package models

import (
	"time"
)

// WorkoutDay represents a training day record
type WorkoutDay struct {
	Id                int        `json:"id" db:"id"`                           // Primary key
	Title             string     `json:"title"`                                // 标题
	Type              string     `json:"type"`                                 // 类型
	Time              string     `json:"time" db:"time"`                       // Training time (YYYY-MM-DD HH:MM:SS)
	Status            int        `json:"status" db:"status"`                   // 0等待进行 1进行中 2已完成 3已过期 4手动作废
	GroupNo           string     `json:"group_no"`                             // 一起训练的标记
	PendingSteps      string     `json:"pending_steps" db:"pending_steps"`     // Execution records in JSON array
	UpdatedDetails    string     `json:"updated_details" db:"updated_details"` // Execution records in JSON array
	Stats             string     `json:"stats" db:"stats"`                     // Statistics in JSON array
	Remark            string     `json:"remark"`
	Medias            string     `json:"medias"`
	EstimatedDuration int        `json:"estimated_duration" db:"estimated_duration"` // 训练预计时长
	Duration          int        `json:"duration"`                                   // 本次训练实际时长 单位 分
	TotalVolume       float64    `json:"total_volume"`                               // 总容量 单位 公斤
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
