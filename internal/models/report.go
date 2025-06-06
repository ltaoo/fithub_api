package models

import "time"

// CoachReport represents a coach's report/feedback in the system
type CoachReport struct {
	Id           int       `json:"id" db:"id"`
	Type         int       `json:"type" db:"type"`                   // Type of the report
	Status       int       `json:"status" db:"status"`               // Status: 1=pending, 2=completed, 3=ignored, 4=revoked
	D            int       `json:"d" db:"d"`                         // Soft delete flag: 0=no, 1=yes
	Content      string    `json:"content" db:"content"`             // Report content
	ReplyContent string    `json:"reply_content" db:"reply_content"` // Admin's reply
	ReasonType   string    `json:"reason_type" db:"reason_type"`     // Type of the reported item (e.g., workout, plan, quiz)
	ReasonId     int       `json:"reason_id" db:"reason_id"`         // ID of the reported item
	CoachId      int       `json:"coach_id" db:"coach_id"`           // ID of the reporting coach
	CreatedAt    time.Time `json:"created_at" db:"created_at"`       // Creation timestamp
}

func (*CoachReport) TableName() string {
	return "COACH_REPORT"
}
