package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
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

type WorkoutDayProgress interface {
	GetVersion() string
}

func ParseWorkoutDayProgress(data string) (WorkoutDayProgress, error) {
	var version struct {
		V string `json:"v"`
	}
	if err := json.Unmarshal([]byte(data), &version); err != nil {
		return nil, err
	}
	switch version.V {
	case "250424":
		var v WorkoutDayProgressJSON250424
		if err := json.Unmarshal([]byte(data), &v); err != nil {
			return nil, err
		}
		return v, nil
	case "250531":
		var v WorkoutDayStepProgressJSON250531
		if err := json.Unmarshal([]byte(data), &v); err != nil {
			return nil, err
		}
		return v, nil
	case "250616":
		var v WorkoutDayStepProgressJSON250616
		if err := json.Unmarshal([]byte(data), &v); err != nil {
			return nil, err
		}
		return v, nil
	case "250629":
		var v WorkoutDayStepProgressJSON250629
		if err := json.Unmarshal([]byte(data), &v); err != nil {
			return nil, err
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unknown version: %s", version.V)
	}
}

type WorkoutDayProgressJSON250424 struct {
	V             string                            `json:"v"`
	StepIdx       int                               `json:"step_idx"`
	SetIdx        int                               `json:"set_idx"`
	ActIdx        int                               `json:"act_idx"`
	TouchedSetIdx []string                          `json:"touched_set_idx"`
	Sets          []WorkoutDayStepProgressSet250424 `json:"sets"`
}

func (w WorkoutDayProgressJSON250424) GetVersion() string { return w.V }

type WorkoutDayStepProgressSet250424 struct {
	StepIdx int                                  `json:"step_idx"`
	Idx     int                                  `json:"idx"`
	Actions []WorkoutDayStepProgressAction250424 `json:"actions"`
}

type WorkoutDayStepProgressAction250424 struct {
	Idx         int     `json:"idx"`
	ActionId    int     `json:"action_id"`
	Reps        int     `json:"reps"`
	RepsUnit    string  `json:"reps_unit"`
	Weight      float64 `json:"weight"`
	WeightUnit  string  `json:"weight_unit"`
	Remark      string  `json:"remark"`
	Completed   bool    `json:"completed"`
	CompletedAt int     `json:"completed_at"`
	Time1       float64 `json:"time1"`
	Time2       float64 `json:"time2"`
	Time3       float64 `json:"time3"`
}

type WorkoutDayStepProgressJSON250531 struct {
	V             string                            `json:"v"`
	StepIdx       int                               `json:"step_idx"`
	SetIdx        int                               `json:"set_idx"`
	ActIdx        int                               `json:"act_idx"`
	TouchedSetIdx []string                          `json:"touched_set_idx"`
	Sets          []WorkoutDayStepProgressSet250531 `json:"sets"`
}

func (w WorkoutDayStepProgressJSON250531) GetVersion() string { return w.V }

type WorkoutDayStepProgressSet250531 struct {
	StepIdx       int                                  `json:"step_idx"`
	Idx           int                                  `json:"idx"`
	Actions       []WorkoutDayStepProgressAction250531 `json:"actions"`
	RemainingTime float64                              `json:"remaining_time"`
	ExceedTime    float64                              `json:"exceed_time"`
	Completed     bool                                 `json:"completed"`
	Remark        string                               `json:"remark"`
}

type WorkoutDayStepProgressAction250531 struct {
	Idx         int         `json:"idx"`
	ActionId    interface{} `json:"action_id"` // int or string
	Reps        int         `json:"reps"`
	RepsUnit    string      `json:"reps_unit"`
	Weight      float64     `json:"weight"`
	WeightUnit  string      `json:"weight_unit"`
	Completed   bool        `json:"completed"`
	CompletedAt int         `json:"completed_at"`
	Time1       float64     `json:"time1"`
	Time2       float64     `json:"time2"`
	Time3       float64     `json:"time3"`
}

type WorkoutDayStepProgressJSON250616 struct {
	V             string                            `json:"v"`
	StepIdx       int                               `json:"step_idx"`
	SetIdx        int                               `json:"set_idx"`
	ActIdx        int                               `json:"act_idx"`
	TouchedSetUid []string                          `json:"touched_set_uid"`
	Sets          []WorkoutDayStepProgressSet250616 `json:"sets"`
}

func (w WorkoutDayStepProgressJSON250616) GetVersion() string { return w.V }

type WorkoutDayStepProgressSet250616 struct {
	StepUid       int                                  `json:"step_uid"`
	Uid           int                                  `json:"uid"`
	Actions       []WorkoutDayStepProgressAction250616 `json:"actions"`
	RemainingTime float64                              `json:"remaining_time"`
	ExceedTime    float64                              `json:"exceed_time"`
	Completed     bool                                 `json:"completed"`
	Remark        string                               `json:"remark"`
}

type WorkoutDayStepProgressAction250616 struct {
	Uid         int     `json:"uid"`
	ActionId    int     `json:"action_id"`
	Reps        int     `json:"reps"`
	RepsUnit    string  `json:"reps_unit"`
	Weight      float64 `json:"weight"`
	WeightUnit  string  `json:"weight_unit"`
	Completed   bool    `json:"completed"`
	CompletedAt int     `json:"completed_at"`
	Time1       float64 `json:"time1"`
	Time2       float64 `json:"time2"`
	Time3       float64 `json:"time3"`
}

type WorkoutDayStepProgressJSON250629 struct {
	V             string                            `json:"v"`
	StepIdx       int                               `json:"step_idx"`
	SetIdx        int                               `json:"set_idx"`
	ActIdx        int                               `json:"act_idx"`
	TouchedSetUid []string                          `json:"touched_set_uid"`
	Sets          []WorkoutDayStepProgressSet250629 `json:"sets"`
}

func (w WorkoutDayStepProgressJSON250629) GetVersion() string { return w.V }

type WorkoutDayStepProgressSet250629 struct {
	StepUid       int                                  `json:"step_uid"`
	Uid           int                                  `json:"uid"`
	Actions       []WorkoutDayStepProgressAction250629 `json:"actions"`
	StartAt       int                                  `json:"start_at"`
	FinishedAt    int                                  `json:"finished_at"`
	RemainingTime float64                              `json:"remaining_time"`
	ExceedTime    float64                              `json:"exceed_time"`
	Completed     bool                                 `json:"completed"`
	Remark        string                               `json:"remark"`
}

type WorkoutDayStepProgressAction250629 struct {
	Uid         int     `json:"uid"`
	ActionId    int     `json:"action_id"`
	ActionName  string  `json:"action_name"`
	Reps        int     `json:"reps"`
	RepsUnit    string  `json:"reps_unit"`
	Weight      float64 `json:"weight"`
	WeightUnit  string  `json:"weight_unit"`
	Completed   bool    `json:"completed"`
	CompletedAt int     `json:"completed_at"`
	StartAt1    int     `json:"start_at1"`
	StartAt2    int     `json:"start_at2"`
	StartAt3    int     `json:"start_at3"`
	FinishedAt1 int     `json:"finished_at1"`
	FinishedAt2 int     `json:"finished_at2"`
	FinishedAt3 int     `json:"finished_at3"`
	Time1       float64 `json:"time1"`
	Time2       float64 `json:"time2"`
	Time3       float64 `json:"time3"`
}

func ToWorkoutDayStepProgress(progress WorkoutDayProgress) WorkoutDayStepProgressJSON250629 {
	switch v := progress.(type) {
	case WorkoutDayProgressJSON250424:
		sets := make([]WorkoutDayStepProgressSet250629, len(v.Sets))
		for i, set := range v.Sets {
			actions := make([]WorkoutDayStepProgressAction250629, len(set.Actions))
			for j, act := range set.Actions {
				actions[j] = WorkoutDayStepProgressAction250629{
					Uid:         0, // 旧版无此字段，补0
					ActionId:    act.ActionId,
					ActionName:  "", // 旧版无此字段，补空
					Reps:        act.Reps,
					RepsUnit:    act.RepsUnit,
					Weight:      act.Weight,
					WeightUnit:  act.WeightUnit,
					Completed:   act.Completed,
					CompletedAt: act.CompletedAt,
					StartAt1:    0,
					StartAt2:    0,
					StartAt3:    0,
					FinishedAt1: 0,
					FinishedAt2: 0,
					FinishedAt3: 0,
					Time1:       act.Time1,
					Time2:       act.Time2,
					Time3:       act.Time3,
				}
			}
			sets[i] = WorkoutDayStepProgressSet250629{
				StepUid:       0, // 旧版无此字段
				Uid:           0,
				Actions:       actions,
				StartAt:       0,
				FinishedAt:    0,
				RemainingTime: 0,
				ExceedTime:    0,
				Completed:     false,
				Remark:        "",
			}
		}
		return WorkoutDayStepProgressJSON250629{
			V:             "250629",
			StepIdx:       v.StepIdx,
			SetIdx:        v.SetIdx,
			ActIdx:        v.ActIdx,
			TouchedSetUid: v.TouchedSetIdx, // 旧版叫 TouchedSetIdx，类型一样
			Sets:          sets,
		}
	case WorkoutDayStepProgressJSON250531:
		sets := make([]WorkoutDayStepProgressSet250629, len(v.Sets))
		for i, set := range v.Sets {
			actions := make([]WorkoutDayStepProgressAction250629, len(set.Actions))
			for j, act := range set.Actions {
				actionId := 0
				switch id := act.ActionId.(type) {
				case int:
					actionId = id
				case float64:
					actionId = int(id)
				case string:
					// 可选：尝试转成 int
				}
				actions[j] = WorkoutDayStepProgressAction250629{
					Uid:         0,
					ActionId:    actionId,
					ActionName:  "",
					Reps:        act.Reps,
					RepsUnit:    act.RepsUnit,
					Weight:      act.Weight,
					WeightUnit:  act.WeightUnit,
					Completed:   act.Completed,
					CompletedAt: act.CompletedAt,
					StartAt1:    0,
					StartAt2:    0,
					StartAt3:    0,
					FinishedAt1: 0,
					FinishedAt2: 0,
					FinishedAt3: 0,
					Time1:       act.Time1,
					Time2:       act.Time2,
					Time3:       act.Time3,
				}
			}
			sets[i] = WorkoutDayStepProgressSet250629{
				StepUid:       0,
				Uid:           0,
				Actions:       actions,
				StartAt:       0,
				FinishedAt:    0,
				RemainingTime: set.RemainingTime,
				ExceedTime:    set.ExceedTime,
				Completed:     set.Completed,
				Remark:        set.Remark,
			}
		}
		return WorkoutDayStepProgressJSON250629{
			V:             "250629",
			StepIdx:       v.StepIdx,
			SetIdx:        v.SetIdx,
			ActIdx:        v.ActIdx,
			TouchedSetUid: v.TouchedSetIdx,
			Sets:          sets,
		}
	case WorkoutDayStepProgressJSON250616:
		sets := make([]WorkoutDayStepProgressSet250629, len(v.Sets))
		for i, set := range v.Sets {
			actions := make([]WorkoutDayStepProgressAction250629, len(set.Actions))
			for j, act := range set.Actions {
				actions[j] = WorkoutDayStepProgressAction250629{
					Uid:         act.Uid,
					ActionId:    act.ActionId,
					ActionName:  "",
					Reps:        act.Reps,
					RepsUnit:    act.RepsUnit,
					Weight:      act.Weight,
					WeightUnit:  act.WeightUnit,
					Completed:   act.Completed,
					CompletedAt: act.CompletedAt,
					StartAt1:    0,
					StartAt2:    0,
					StartAt3:    0,
					FinishedAt1: 0,
					FinishedAt2: 0,
					FinishedAt3: 0,
					Time1:       float64(act.Time1),
					Time2:       float64(act.Time2),
					Time3:       act.Time3,
				}
			}
			sets[i] = WorkoutDayStepProgressSet250629{
				StepUid:       set.StepUid,
				Uid:           set.Uid,
				Actions:       actions,
				StartAt:       0,
				FinishedAt:    0,
				RemainingTime: set.RemainingTime,
				ExceedTime:    set.ExceedTime,
				Completed:     set.Completed,
				Remark:        set.Remark,
			}
		}
		return WorkoutDayStepProgressJSON250629{
			V:             "250629",
			StepIdx:       v.StepIdx,
			SetIdx:        v.SetIdx,
			ActIdx:        v.ActIdx,
			TouchedSetUid: v.TouchedSetUid,
			Sets:          sets,
		}
	case WorkoutDayStepProgressJSON250629:
		return v
	default:
		return WorkoutDayStepProgressJSON250629{}
	}
}

type WorkoutDayStepDetails interface {
	GetVersion() string
}

func ParseWorkoutDayStepDetails(data string) (WorkoutDayStepDetails, error) {
	var version struct {
		V string `json:"v"`
	}
	if err := json.Unmarshal([]byte(data), &version); err != nil {
		return nil, err
	}
	switch version.V {
	case "250424":
		var v WorkoutDayStepDetailsJSON250424
		if err := json.Unmarshal([]byte(data), &v); err != nil {
			return nil, err
		}
		return v, nil
	case "250616":
		var v WorkoutDayStepDetailsJSON250616
		if err := json.Unmarshal([]byte(data), &v); err != nil {
			return nil, err
		}
		return v, nil
	case "250629":
		var v WorkoutDayStepDetailsJSON250629
		if err := json.Unmarshal([]byte(data), &v); err != nil {
			return nil, err
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unknown version: %s", version.V)
	}
}

func WorkoutDayStepDetailsToWorkoutPlanBodyDetails(details WorkoutDayStepDetails) WorkoutPlanBodyDetailsJSON250627 {
	switch v := details.(type) {
	case WorkoutDayStepDetailsJSON250424:
		steps := make([]WorkoutPlanBodyStepJSON250627, len(v.Steps))
		for step_uid, step := range v.Steps {
			if len(step.Sets) == 0 {
				continue
			}
			first_set := step.Sets[0]
			if len(first_set.Actions) == 0 {
				continue
			}
			actions := make([]WorkoutPlanBodyStepActionJSON250627, len(first_set.Actions))
			for act_uid, act := range first_set.Actions {
				actions[act_uid] = WorkoutPlanBodyStepActionJSON250627{
					Action: WorkoutActionJSON250627{
						Id:     act.Id,
						ZhName: act.ZhName,
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
				SetType:  first_set.Type,
				SetCount: len(step.Sets),
				SetRestDuration: WorkoutRestDuration{
					Num:  first_set.RestDuration,
					Unit: "秒",
				},
				SetWeight: WorkoutWeight{
					Num:  first_set.Weight,
					Unit: "RPE",
				},
				SetNote: "",
				Actions: actions,
			}
		}
		return WorkoutPlanBodyDetailsJSON250627{
			V:     "250627",
			Steps: steps,
		}
	case WorkoutDayStepDetailsJSON250616:
		steps := make([]WorkoutPlanBodyStepJSON250627, len(v.Steps))
		for step_uid, step := range v.Steps {
			if len(step.Sets) == 0 {
				continue
			}
			first_set := step.Sets[0]
			if len(first_set.Actions) == 0 {
				continue
			}
			actions := make([]WorkoutPlanBodyStepActionJSON250627, len(first_set.Actions))
			for act_uid, act := range first_set.Actions {
				actions[act_uid] = WorkoutPlanBodyStepActionJSON250627{
					Action: WorkoutActionJSON250627{
						Id:     act.Id,
						ZhName: act.ZhName,
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
				SetType:  first_set.Type,
				SetCount: len(step.Sets),
				SetRestDuration: WorkoutRestDuration{
					Num:  first_set.RestDuration,
					Unit: "秒",
				},
				SetWeight: WorkoutWeight{
					Num:  first_set.Weight,
					Unit: "RPE",
				},
				SetNote: "",
				Actions: actions,
			}
		}
		return WorkoutPlanBodyDetailsJSON250627{
			V:     "250627",
			Steps: steps,
		}
	case WorkoutDayStepDetailsJSON250629:
		steps := make([]WorkoutPlanBodyStepJSON250627, len(v.Steps))
		for step_uid, step := range v.Steps {
			if len(step.Sets) == 0 {
				continue
			}
			first_set := step.Sets[0]
			if len(first_set.Actions) == 0 {
				continue
			}
			actions := make([]WorkoutPlanBodyStepActionJSON250627, len(first_set.Actions))
			for act_uid, act := range first_set.Actions {
				actions[act_uid] = WorkoutPlanBodyStepActionJSON250627{
					Action: WorkoutActionJSON250627{
						Id:     act.Id,
						ZhName: act.ZhName,
					},
					Reps:         act.Reps,
					Weight:       act.Weight,
					RestDuration: act.RestDuration,
				}
			}
			steps[step_uid] = WorkoutPlanBodyStepJSON250627{
				StepUid:         step.Uid,
				SetType:         first_set.Type,
				SetCount:        len(step.Sets),
				SetRestDuration: first_set.RestDuration,
				SetWeight:       first_set.Weight,
				SetNote:         "",
				Actions:         actions,
			}
		}
		return WorkoutPlanBodyDetailsJSON250627{
			V:     "250627",
			Steps: steps,
		}
	default:
		return WorkoutPlanBodyDetailsJSON250627{}
	}
}

type WorkoutDayUpdatedDetails interface {
	GetVersion() string
}

func ParseWorkoutDayUpdatedDetails(data string) (WorkoutDayUpdatedDetails, error) {
	var versionHolder struct {
		V string `json:"v"`
	}
	if err := json.Unmarshal([]byte(data), &versionHolder); err != nil {
		return nil, err
	}
	switch versionHolder.V {
	case "250424":
		var v WorkoutDayStepDetailsJSON250424
		if err := json.Unmarshal([]byte(data), &v); err != nil {
			return nil, err
		}
		return v, nil
	case "250616":
		var v WorkoutDayStepDetailsJSON250616
		if err := json.Unmarshal([]byte(data), &v); err != nil {
			return nil, err
		}
		return v, nil
	case "250629":
		var v WorkoutDayStepDetailsJSON250629
		if err := json.Unmarshal([]byte(data), &v); err != nil {
			return nil, err
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unknown version: %s", versionHolder.V)
	}
}

func ToWorkoutDayStepDetails(progress WorkoutDayUpdatedDetails) WorkoutDayStepDetailsJSON250629 {
	switch v := progress.(type) {
	case WorkoutDayStepDetailsJSON250424:
		steps := make([]WorkoutDayStepDetailsStep250629, len(v.Steps))
		for step_uid, step := range v.Steps {
			sets := make([]WorkoutDayStepDetailsSet250629, len(step.Sets))
			for set_uid, set := range step.Sets {
				acts := make([]WorkoutDayStepDetailsAction250629, len(set.Actions))
				for act_uid, act := range set.Actions {
					acts = append(acts, WorkoutDayStepDetailsAction250629{
						Uid: act_uid,
						Id:  act.Id,
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
					})
				}
				sets[set_uid] = WorkoutDayStepDetailsSet250629{
					Uid:     set_uid,
					Type:    set.Type,
					Actions: acts,
					RestDuration: WorkoutRestDuration{
						Num:  set.RestDuration,
						Unit: "秒",
					},
					Weight: WorkoutWeight{
						Num:  set.Weight,
						Unit: "RPE",
					},
				}
			}
			steps[step_uid] = WorkoutDayStepDetailsStep250629{
				Uid:  step_uid,
				Sets: sets,
				Note: step.Note,
			}
		}
		return WorkoutDayStepDetailsJSON250629{
			V:     "250629",
			Steps: steps,
		}
	case WorkoutDayStepDetailsJSON250616:
		steps := make([]WorkoutDayStepDetailsStep250629, len(v.Steps))
		for step_uid, step := range v.Steps {
			sets := make([]WorkoutDayStepDetailsSet250629, len(step.Sets))
			for set_uid, set := range step.Sets {
				acts := make([]WorkoutDayStepDetailsAction250629, len(set.Actions))
				for act_uid, act := range set.Actions {
					acts = append(acts, WorkoutDayStepDetailsAction250629{
						Uid: act_uid,
						Id:  act.Id,
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
					})
				}
				sets[set_uid] = WorkoutDayStepDetailsSet250629{
					Uid:     set_uid,
					Type:    set.Type,
					Actions: acts,
					RestDuration: WorkoutRestDuration{
						Num:  set.RestDuration,
						Unit: "秒",
					},
					Weight: WorkoutWeight{
						Num:  set.Weight,
						Unit: "RPE",
					},
				}
			}
			steps[step_uid] = WorkoutDayStepDetailsStep250629{
				Uid:  step_uid,
				Sets: sets,
				Note: step.Note,
			}
		}
		return WorkoutDayStepDetailsJSON250629{
			V:     "250629",
			Steps: steps,
		}
	case WorkoutDayStepDetailsJSON250629:
		return v
	default:
		return WorkoutDayStepDetailsJSON250629{}
	}
}

type WorkoutDayStepDetailsJSON250424 struct {
	V     string                            `json:"v"`
	Steps []WorkoutDayStepDetailsStep250424 `json:"steps"`
}

func (w WorkoutDayStepDetailsJSON250424) GetVersion() string { return w.V }

type WorkoutDayStepDetailsStep250424 struct {
	Idx  int                              `json:"idx"`
	Sets []WorkoutDayStepDetailsSet250424 `json:"sets"`
	Note string                           `json:"note"`
}

type WorkoutDayStepDetailsSet250424 struct {
	Idx          int                                 `json:"idx"`
	Type         string                              `json:"type"` // normal super increasing decreasing hiit
	Actions      []WorkoutDayStepDetailsAction250424 `json:"actions"`
	RestDuration int                                 `json:"rest_duration"`
	Weight       string                              `json:"weight"`
}

type WorkoutDayStepDetailsAction250424 struct {
	Id           int    `json:"id"`
	ZhName       string `json:"zh_name"`
	Reps         int    `json:"reps"`
	RepsUnit     string `json:"reps_unit"` // SetValueUnit
	Weight       string `json:"weight"`
	RestDuration int    `json:"rest_duration"`
}
type WorkoutDayStepDetailsJSON250616 struct {
	V     string                            `json:"v"`
	Steps []WorkoutDayStepDetailsStep250616 `json:"steps"`
}

func (w WorkoutDayStepDetailsJSON250616) GetVersion() string { return w.V }

type WorkoutDayStepDetailsStep250616 struct {
	Uid  int                              `json:"uid"`
	Sets []WorkoutDayStepDetailsSet250616 `json:"sets"`
	Note string                           `json:"note"`
}

type WorkoutDayStepDetailsSet250616 struct {
	Uid          int                                 `json:"uid"`
	Type         string                              `json:"type"` // WorkoutPlanSetType
	Actions      []WorkoutDayStepDetailsAction250616 `json:"actions"`
	RestDuration int                                 `json:"rest_duration"`
	Weight       string                              `json:"weight"`
}

type WorkoutDayStepDetailsAction250616 struct {
	Uid          int    `json:"uid"`
	Id           int    `json:"id"`
	ZhName       string `json:"zh_name"`
	Reps         int    `json:"reps"`
	RepsUnit     string `json:"reps_unit"` // SetValueUnit
	Weight       string `json:"weight"`
	RestDuration int    `json:"rest_duration"`
}

type WorkoutDayStepDetailsJSON250629 struct {
	V     string                            `json:"v"`
	Steps []WorkoutDayStepDetailsStep250629 `json:"steps"`
}

func (w WorkoutDayStepDetailsJSON250629) GetVersion() string { return w.V }

type WorkoutDayStepDetailsStep250629 struct {
	Uid  int                              `json:"uid"`
	Sets []WorkoutDayStepDetailsSet250629 `json:"sets"`
	Note string                           `json:"note"`
}

type WorkoutDayStepDetailsSet250629 struct {
	Uid          int                                 `json:"uid"`
	Type         string                              `json:"type"` // WorkoutPlanSetType
	Actions      []WorkoutDayStepDetailsAction250629 `json:"actions"`
	RestDuration WorkoutRestDuration                 `json:"rest_duration"`
	Weight       WorkoutWeight                       `json:"weight"`
}

type WorkoutDayStepDetailsAction250629 struct {
	Uid          int                 `json:"uid"`
	Id           int                 `json:"id"`
	ZhName       string              `json:"zh_name"`
	Reps         WorkoutReps         `json:"reps"`
	Weight       WorkoutWeight       `json:"weight"`
	RestDuration WorkoutRestDuration `json:"rest_duration"`
}

type TodayWorkoutAction struct {
	ActionId   int     `json:"action_id"`
	ActionName string  `json:"action_name"`
	Reps       int     `json:"reps"`
	RepsUnit   string  `json:"reps_unit"`
	Weight     float64 `json:"weight"`
	WeightUnit string  `json:"weight_unit"`
}
type TodayWorkoutActionGroupSet struct {
	Idx     int                  `json:"idx"`
	Texts   []string             `json:"texts"`
	Actions []TodayWorkoutAction `json:"actions"`
}
type TodayWorkoutActionGroup struct {
	Title       string                       `json:"title"`
	Type        string                       `json:"type"`
	TotalVolume float64                      `json:"total_volume"`
	Duration    int                          `json:"duration"`
	Sets        []TodayWorkoutActionGroupSet `json:"sets"`
}

type BuildedWorkoutDayResult struct {
	TotalVolume   float64                   `json:"float64"`
	SetCount      int                       `json:"set_count"`
	DurationCount int                       `json:"duration_count"`
	List          []TodayWorkoutActionGroup `json:"list"`
	Tags          []string                  `json:"tags"`
}

func BuildResultFromWorkoutDay(v WorkoutDay, db *gorm.DB) (*BuildedWorkoutDayResult, error) {
	if v.Status != 2 {
		return &BuildedWorkoutDayResult{}, nil
	}
	list := make([]TodayWorkoutActionGroup, 0)
	tags := make([]string, 0)
	total_volume := float64(0)
	set_count := 0
	duration_count := v.Duration
	tmp_pending_steps, err := ParseWorkoutDayProgress(v.PendingSteps)
	if err != nil {
		return nil, err
	}
	var workout_plan_details WorkoutPlanBodyDetailsJSON250627
	if v.UpdatedDetails != "" {
		tmp_details, err := ParseWorkoutDayUpdatedDetails(v.UpdatedDetails)
		if err != nil {
			// error_msg = append(error_msg, err.Error())
			return nil, err
		}
		workout_plan_details = WorkoutDayStepDetailsToWorkoutPlanBodyDetails(tmp_details)
	} else {
		var existing_workout_plan WorkoutPlan
		if err := db.Where("id = ?", v.WorkoutPlanId).First(&existing_workout_plan).Error; err != nil {
			// error_msg = append(error_msg, err.Error())
			// continue
			return nil, err
		}
		// fmt.Println("the profile", existing_workout_plan.Details)
		tmp_details, err := ParseWorkoutPlanDetail(existing_workout_plan.Details)
		if err != nil {
			// error_msg = append(error_msg, err.Error())
			// continue
			return nil, err
		}
		workout_plan_details = ToWorkoutPlanBodyDetails(tmp_details)
	}
	if len(workout_plan_details.Steps) == 0 {
		// error_msg = append(error_msg, "没有解析出数据1")
		// continue
		return nil, err
	}
	pending_steps := ToWorkoutDayStepProgress(tmp_pending_steps)
	if len(pending_steps.Sets) == 0 {
		// error_msg = append(error_msg, "没有解析出数据2")
		// continue
		return nil, err
	}
	day_total_volume := float64(0)
	for _, step := range workout_plan_details.Steps {
		var pending_sets []WorkoutDayStepProgressSet250629
		for _, vv := range pending_steps.Sets {
			if vv.StepUid == step.StepUid {
				pending_sets = append(pending_sets, vv)
			}
		}
		if len(pending_sets) != 0 {
			first_set := pending_sets[0]
			var sets []TodayWorkoutActionGroupSet
			var action_names []string
			var title string
			if step.SetType == "super" {
				for _, s := range first_set.Actions {
					action_names = append(action_names, s.ActionName)
				}
				title = strings.Join(action_names, " + ") + " 超级组"
			}
			if step.SetType == "hiit" {
				_, ok := lo.Find(tags, func(i string) bool {
					return i == v.Type
				})
				if !ok {
					tags = append(tags, "hiit")
				}
				for _, s := range first_set.Actions {
					action_names = append(action_names, s.ActionName)
				}
				title = strings.Join(action_names, " + ") + " HIIT"
			}
			if step.SetType == "decreasing" && len(first_set.Actions) > 0 {
				action_names = append(action_names, first_set.Actions[0].ActionName)
				title = strings.Join(action_names, " + ") + " 递减组"
			}
			if step.SetType == "normal" && len(first_set.Actions) > 0 {
				action_names = append(action_names, first_set.Actions[0].ActionName)
				title = strings.Join(action_names, "")
			}
			for idx, ss := range pending_sets {
				if ss.Completed {
					set_count += 1
					var texts []string
					var actions []TodayWorkoutAction
					for _, act := range ss.Actions {
						// act_id := act.ActionId
						act_name := act.ActionName
						reps := act.Reps
						reps_unit := act.RepsUnit
						weight := act.Weight
						weight_unit := act.WeightUnit
						if weight_unit == "公斤" && reps_unit == "次" {
							day_total_volume += float64(reps) * weight
						}
						if weight_unit == "磅" && reps_unit == "次" {
							day_total_volume += float64(reps) * (weight * 0.45)
						}
						// weight_text := act.WeightUnit == "自重"
						// fmt.Println("title", act.Weight)
						// weight_text := fmt.Sprintf("%#g", act.Weight)
						weight_text := strconv.FormatFloat(act.Weight, 'g', -1, 64) + act.WeightUnit
						if act.WeightUnit == "自重" {
							weight_text = act.WeightUnit
						}
						reps_text := strconv.Itoa(act.Reps) + act.RepsUnit
						t := weight_text + "x" + reps_text
						if step.SetType != "normal" && step.SetType != "decreasing" {
							t = act_name + " " + t
						}
						texts = append(texts, t)
						actions = append(actions, TodayWorkoutAction{
							ActionId:   act.ActionId,
							ActionName: act.ActionName,
							Reps:       act.Reps,
							RepsUnit:   act.RepsUnit,
							Weight:     act.Weight,
							WeightUnit: act.WeightUnit,
						})
					}
					sets = append(sets, TodayWorkoutActionGroupSet{
						Idx:     idx + 1,
						Texts:   texts,
						Actions: actions,
					})
				}
			}
			if len(sets) != 0 {
				dd := TodayWorkoutActionGroup{
					Title:       title,
					Type:        step.SetType,
					Sets:        sets,
					Duration:    v.Duration,
					TotalVolume: day_total_volume,
				}
				if v.Type == "cardio" {
					dd = TodayWorkoutActionGroup{
						Title:       title,
						Type:        "cardio",
						Sets:        make([]TodayWorkoutActionGroupSet, 0),
						Duration:    v.Duration,
						TotalVolume: 0,
					}
				}
				list = append(list, dd)
			}
		}
	}
	total_volume += day_total_volume
	_, ok := lo.Find(tags, func(i string) bool {
		return i == v.Type
	})
	if !ok && v.Type != "" {
		tags = append(tags, v.Type)
	}
	return &BuildedWorkoutDayResult{
		SetCount:      set_count,
		DurationCount: duration_count,
		TotalVolume:   total_volume,
		List:          list,
		Tags:          tags,
	}, nil
}
