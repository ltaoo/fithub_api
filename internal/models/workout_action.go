package models

import (
	"time"
)

// WorkoutAction represents a fitness exercise or movement
type WorkoutAction struct {
	Id                   int        `json:"id" gorm:"primaryKey"`
	D                    int        `json:"d" gorm:"column:d;default:0"`
	Status               int        `json:"status" gorm:"column:status;default:1"`
	Name                 string     `json:"name" gorm:"column:name;not null"`
	ZhName               string     `json:"zh_name" gorm:"column:zh_name;not null"`
	Alias                string     `json:"alias" gorm:"column:alias;not null"`
	Overview             string     `json:"overview" gorm:"column:overview;not null"`
	CoverURL             string     `json:"cover_url"`
	SortIdx              int        `json:"sort_idx"`
	Type                 string     `json:"type" gorm:"column:type;not null"`
	Level                int        `json:"level" gorm:"column:level;default:1"`
	Score                int        `json:"score"`
	Tags1                string     `json:"tags1" gorm:"column:tags1;not null"`
	Tags2                string     `json:"tags2" gorm:"column:tags2;not null"`
	Pattern              string     `json:"pattern"`
	Details              string     `json:"details" gorm:"column:details;default:'{}';not null"`
	Points               string     `json:"points" gorm:"column:points;not null"`
	Problems             string     `json:"problems" gorm:"column:problems;not null"`
	ExtraConfig          string     `json:"extra_config"`
	CreatedAt            time.Time  `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt            *time.Time `json:"updated_at"`
	MuscleIds            string     `json:"muscle_ids" gorm:"column:muscle_ids"`
	PrimaryMuscleIds     string     `json:"primary_muscle_ids"`
	SecondaryMuscleIds   string     `json:"secondary_muscle_ids"`
	EquipmentIds         string     `json:"equipment_ids" gorm:"column:equipment_ids"`
	AlternativeActionIds string     `json:"alternative_action_ids" gorm:"column:alternative_action_ids"`
	AdvancedActionIds    string     `json:"advanced_action_ids" gorm:"column:advanced_action_ids"`
	RegressedActionIds   string     `json:"regressed_action_ids" gorm:"column:regressed_action_ids"`
	OwnerId              int        `json:"owner_id" gorm:"column:owner_id"`
}

// TableName specifies the table name for the WorkoutAction model
func (WorkoutAction) TableName() string {
	return "WORKOUT_ACTION"
}
