package models

import (
	"encoding/json"
	"fmt"
	"time"
)

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
// 废弃了，现在是用 JSON 实现
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

// func WorkoutPlanBodyDetailsToWorkoutPlanDetails(detail WorkoutPlanBodyDetailsJSON250627) WorkoutPlanDetails {
// 	if len(detail.Steps) == 0 {
// 		return WorkoutPlanBodyDetailsJSON250627{
// 			V:     "250627",
// 			Steps: make([]WorkoutPlanStepDetails, 0),
// 		}
// 	}
// 	var steps []WorkoutPlanStepDetails
// 	for _, step := range detail.Steps {
// 		var actions []WorkoutPlanActionDetails
// 		for _, v := range step.Actions {
// 			actions = append(actions, WorkoutPlanActionDetails{
// 				Action: struct {
// 					Id     int    `json:"id"`
// 					ZhName string `json:"zh_name"`
// 				}{
// 					Id:     v.Action.Id,
// 					ZhName: v.Action.ZhName,
// 				},
// 				Reps:         v.Reps,
// 				Weight:       v.Weight,
// 				RestDuration: v.RestDuration,
// 			})
// 		}
// 		d := WorkoutPlanStepDetails{
// 			SetType:         step.SetType,
// 			SetCount:        step.SetCount,
// 			SetRestDuration: step.SetRestDuration,
// 			SetWeight:       step.SetWeight,
// 			SetNote:         step.SetNote,
// 			SetTags:         step.SetTags,
// 			Actions:         actions,
// 		}
// 		steps = append(steps, d)
// 	}
// 	return WorkoutPlanBodyDetailsJSON250627{
// 		V:     "250627",
// 		Steps: steps,
// 	}
// }

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

type WorkoutPlanDetails interface {
	GetVersion() string
}

func ParseWorkoutPlanDetail(data string) (WorkoutPlanDetails, error) {
	var version struct {
		V string `json:"v"`
	}
	if err := json.Unmarshal([]byte(data), &version); err != nil {
		return nil, err
	}
	switch version.V {
	case "250424":
		var v WorkoutPlanBodyDetailsJSON250424
		if err := json.Unmarshal([]byte(data), &v); err != nil {
			return nil, err
		}
		return v, nil
	case "250627":
		var v WorkoutPlanBodyDetailsJSON250627
		if err := json.Unmarshal([]byte(data), &v); err != nil {
			return nil, err
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unknown version: %s", version.V)
	}
}

func ToWorkoutPlanBodyDetails(details WorkoutPlanDetails) WorkoutPlanBodyDetailsJSON250627 {
	switch v := details.(type) {
	case WorkoutPlanBodyDetailsJSON250424:
		steps := make([]WorkoutPlanBodyStepJSON250627, len(v.Steps))
		for step_uid, step := range v.Steps {
			actions := make([]WorkoutPlanBodyStepActionJSON250627, len(step.Actions))
			for act_uid, act := range step.Actions {
				actions[act_uid] = WorkoutPlanBodyStepActionJSON250627{
					Action: WorkoutActionJSON250627{
						Id:     act.Action.Id,
						ZhName: act.Action.ZhName,
					},
					Reps: WorkoutReps{
						Num:  act.Reps,
						Unit: act.RepsUnit,
					},
					Weight: WorkoutWeight{
						Num:  act.Weight,
						Unit: "RM",
					},
					RestDuration: WorkoutRestDuration{
						Num:  act.RestDuration,
						Unit: "秒",
					},
				}
			}
			steps[step_uid] = WorkoutPlanBodyStepJSON250627{
				StepUid:  step_uid + 1,
				SetType:  step.SetType,
				SetCount: step.SetCount,
				SetRestDuration: WorkoutRestDuration{
					Num:  step.SetRestDuration,
					Unit: "秒",
				},
				SetWeight: WorkoutWeight{
					Num:  step.SetWeight,
					Unit: "RPE",
				},
				SetNote: step.SetNote,
			}
		}
		return WorkoutPlanBodyDetailsJSON250627{
			V:     "250627",
			Steps: steps,
		}
	case WorkoutPlanBodyDetailsJSON250627:
		return WorkoutPlanBodyDetailsJSON250627{
			V:     "250627",
			Steps: v.Steps,
		}
	default:
		return WorkoutPlanBodyDetailsJSON250627{}
	}
}

type WorkoutPlanStepDetails struct {
	SetType         string                     `json:"set_type"`
	SetCount        int                        `json:"set_count"`
	SetRestDuration WorkoutRestDuration        `json:"set_rest_duration"`
	SetWeight       WorkoutWeight              `json:"set_weight"`
	SetNote         string                     `json:"set_note"`
	SetTags         string                     `json:"set_tags"`
	Actions         []WorkoutPlanActionDetails `json:"actions"`
}
type WorkoutPlanActionDetails struct {
	Action struct {
		Id     int    `json:"id"`
		ZhName string `json:"zh_name"`
	} `json:"action"`
	Reps         WorkoutReps         `json:"reps"`
	Weight       WorkoutWeight       `json:"weight"`
	RestDuration WorkoutRestDuration `json:"rest_duration"`
}

type WorkoutReps struct {
	Num  int    `json:"num"`
	Unit string `json:"unit"` // SetValueUnit
}
type WorkoutWeight struct {
	Num  string `json:"num"`
	Unit string `json:"unit"` // SetValueUnit
}
type WorkoutRestDuration struct {
	Num  int    `json:"num"`
	Unit string `json:"unit"` // SetValueUnit
}

type WorkoutPlanBodyDetailsJSON250424 struct {
	V     string                          `json:"v"`
	Steps []WorkoutPlanBodyStepJSON250424 `json:"steps"`
}

func (w WorkoutPlanBodyDetailsJSON250424) GetVersion() string { return w.V }

type WorkoutPlanBodyStepJSON250424 struct {
	SetType         string                                `json:"set_type"` // WorkoutPlanSetType
	Actions         []WorkoutPlanBodyStepActionJSON250424 `json:"actions"`
	SetCount        int                                   `json:"set_count"`
	SetRestDuration int                                   `json:"set_rest_duration"`
	SetWeight       string                                `json:"set_weight"`
	SetNote         string                                `json:"set_note"`
}

type WorkoutPlanBodyStepActionJSON250424 struct {
	ActionId     int                     `json:"action_id"`
	Action       WorkoutActionJSON250424 `json:"action"`
	Reps         int                     `json:"reps"`
	RepsUnit     string                  `json:"reps_unit"` // SetValueUnit
	Weight       string                  `json:"weight"`
	RestDuration int                     `json:"rest_duration"`
}

type WorkoutActionJSON250424 struct {
	Id     int    `json:"id"`
	ZhName string `json:"zh_name"`
}

type WorkoutPlanBodyStepJSON250607 struct {
	Title           string                                `json:"title"`
	Type            string                                `json:"type"` // WorkoutPlanStepType
	Idx             int                                   `json:"idx"`
	SetType         string                                `json:"set_type"` // WorkoutPlanSetType
	SetCount        int                                   `json:"set_count"`
	SetRestDuration int                                   `json:"set_rest_duration"`
	SetWeight       string                                `json:"set_weight"`
	Actions         []WorkoutPlanBodyStepActionJSON250607 `json:"actions"`
	SetNote         string                                `json:"set_note"`
}
type WorkoutPlanBodyStepActionJSON250607 struct {
	ActionId     int                     `json:"action_id"`
	Action       WorkoutActionJSON250607 `json:"action"`
	Idx          int                     `json:"idx"`
	SetIdx       int                     `json:"set_idx"`
	Reps         int                     `json:"reps"`
	RepsUnit     string                  `json:"reps_unit"` // SetValueUnit
	Weight       string                  `json:"weight"`
	RestDuration int                     `json:"rest_duration"`
	Note         string                  `json:"note"`
}

type WorkoutActionJSON250607 struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	ZhName string `json:"zh_name"`
}

// 250627 版本
type WorkoutPlanBodyDetailsJSON250627 struct {
	V     string                          `json:"v"`
	Steps []WorkoutPlanBodyStepJSON250627 `json:"steps"`
}

func (w WorkoutPlanBodyDetailsJSON250627) GetVersion() string { return w.V }

type WorkoutPlanBodyStepJSON250627 struct {
	StepUid         int                                   `json:"step_uid"`
	SetType         string                                `json:"set_type"` // WorkoutPlanSetType
	SetCount        int                                   `json:"set_count"`
	SetRestDuration WorkoutRestDuration                   `json:"set_rest_duration"`
	SetWeight       WorkoutWeight                         `json:"set_weight"`
	SetNote         string                                `json:"set_note"`
	SetTags         string                                `json:"set_tags"`
	Actions         []WorkoutPlanBodyStepActionJSON250627 `json:"actions"`
}

type WorkoutPlanBodyStepActionJSON250627 struct {
	Action       WorkoutActionJSON250627 `json:"action"`
	Reps         WorkoutReps             `json:"reps"`
	Weight       WorkoutWeight           `json:"weight"`
	RestDuration WorkoutRestDuration     `json:"rest_duration"`
}

type WorkoutActionJSON250627 struct {
	Id     int    `json:"id"`
	ZhName string `json:"zh_name"`
}
