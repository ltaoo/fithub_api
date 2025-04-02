package models

import (
	"time"
)

// WorkoutAction represents a fitness exercise or movement
type WorkoutAction struct {
	Id                   uint       `json:"id" gorm:"primaryKey"`
	Name                 string     `json:"name" gorm:"column:name;not null"`
	ZhName               string     `json:"zh_name" gorm:"column:zh_name;not null"`
	Alias                string     `json:"alias" gorm:"column:alias;not null"`
	Overview             string     `json:"overview" gorm:"column:overview;not null"`
	Type                 string     `json:"type" gorm:"column:type;not null"`
	Level                int        `json:"level" gorm:"column:level;default:1"`
	Tags1                string     `json:"tags1" gorm:"column:tags1;not null"`
	Tags2                string     `json:"tags2" gorm:"column:tags2;not null"`
	Details              string     `json:"details" gorm:"column:details;default:'{}';not null"`
	Points               string     `json:"points" gorm:"column:points;not null"`
	Problems             string     `json:"problems" gorm:"column:problems;not null"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            *time.Time `json:"updated_at"`
	MuscleIds            string     `json:"muscle_ids" gorm:"column:muscle_ids"`
	EquipmentIds         string     `json:"equipment_ids" gorm:"column:equipment_ids"`
	AlternativeActionIds string     `json:"alternative_action_ids" gorm:"column:alternative_action_ids"`
	AdvancedActionIds    string     `json:"advanced_action_ids" gorm:"column:advanced_action_ids"`
	RegressedActionIds   string     `json:"regressed_action_ids" gorm:"column:regressed_action_ids"`
}

// TableName specifies the table name for the WorkoutAction model
func (WorkoutAction) TableName() string {
	return "WORKOUT_ACTION"
}
