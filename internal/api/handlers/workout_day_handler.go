package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myapi/internal/models"
	"myapi/internal/pkg/pagination"
	"myapi/pkg/logger"
)

// WorkoutPlanHandler handles HTTP requests for workout plans
type WorkoutDayHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewWorkoutDayHandler creates a new workout day handler
func NewWorkoutDayHandler(db *gorm.DB, logger *logger.Logger) *WorkoutDayHandler {
	return &WorkoutDayHandler{
		db:     db,
		logger: logger,
	}
}

func (h *WorkoutDayHandler) CreateWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "非法操作", "data": nil})
		return
	}
	var body struct {
		StudentIds      []int  `json:"student_ids"`
		StartWhenCreate bool   `json:"start_when_create"`
		Title           string `json:"title"`
		Type            string `json:"type"`
		Details         string `json:"details"`
		WorkoutPlanId   int    `json:"workout_plan_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.WorkoutPlanId == 0 && body.Details == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少训练内容", "data": nil})
		return
	}
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	var workout_plan models.WorkoutPlan
	if body.WorkoutPlanId != 0 {
		if err := tx.Where("id = ?", body.WorkoutPlanId).First(&workout_plan).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 101, "msg": err.Error(), "data": nil})
			return
		}
	}
	if body.Title != "" && workout_plan.Title == "" {
		workout_plan.Title = body.Title
	}
	if body.Type != "" && workout_plan.Type == "" {
		workout_plan.Type = body.Type
	}
	var subscription models.Subscription
	if err := tx.Where("coach_id = ? AND step = 2", uid).First(&subscription).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 101, "msg": "该功能需订阅后才能使用", "data": nil})
		return
	}
	now := time.Now().UTC()
	group_no := strconv.FormatInt(now.Unix(), 10)
	var workout_day_ids []int
	if len(body.StudentIds) != 0 {
		for _, student_id := range body.StudentIds {
			is_coach_self := student_id == uid || student_id == 0
			if is_coach_self {
				workout_day := models.WorkoutDay{
					Title:  workout_plan.Title,
					Type:   workout_plan.Type,
					Status: int(models.WorkoutDayStatusPending),
					// 只有选了多人才有该字段
					GroupNo:       group_no,
					CreatedAt:     now,
					WorkoutPlanId: body.WorkoutPlanId,
					CoachId:       uid,
					StudentId:     uid,
				}
				if body.Details != "" {
					workout_day.UpdatedDetails = body.Details
				}
				if body.StartWhenCreate {
					workout_day.StartedAt = &now
					workout_day.Status = int(models.WorkoutDayStatusStarted)
				}
				if err := tx.Create(&workout_day).Error; err != nil {
					h.logger.Error("Failed to create workout plan", err)
					tx.Rollback()
					c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
					return
				}
				workout_day_ids = append(workout_day_ids, workout_day.Id)
				continue
			}
			var relation models.CoachRelationship
			if err := tx.Where("coach_id = ? AND student_id = ? OR coach_id = ? AND student_id = ?", uid, student_id, student_id, uid).First(&relation).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
					return
				}
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "非法操作", "data": nil})
				return
			}
			workout_day := models.WorkoutDay{
				Title:  workout_plan.Title,
				Type:   workout_plan.Type,
				Status: int(models.WorkoutDayStatusPending),
				// 只有选了多人才有该字段
				GroupNo:       group_no,
				CreatedAt:     now,
				WorkoutPlanId: body.WorkoutPlanId,
				CoachId:       uid,
				StudentId:     student_id,
			}
			if body.Details != "" {
				workout_day.UpdatedDetails = body.Details
			}
			if body.StartWhenCreate {
				workout_day.StartedAt = &now
				workout_day.Status = int(models.WorkoutDayStatusStarted)
			}
			if err := tx.Create(&workout_day).Error; err != nil {
				h.logger.Error("Failed to create workout plan", err)
				tx.Rollback()
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
				return
			}
			workout_day_ids = append(workout_day_ids, workout_day.Id)
		}
	} else {
		workout_day := models.WorkoutDay{
			Title:         workout_plan.Title,
			Type:          workout_plan.Type,
			Status:        int(models.WorkoutDayStatusPending),
			CreatedAt:     now,
			WorkoutPlanId: body.WorkoutPlanId,
			CoachId:       uid,
			StudentId:     uid,
		}
		if body.Details != "" {
			workout_day.UpdatedDetails = body.Details
		}
		if body.StartWhenCreate {
			workout_day.StartedAt = &now
			workout_day.Status = int(models.WorkoutDayStatusStarted)
		}
		if err := tx.Create(&workout_day).Error; err != nil {
			h.logger.Error("Failed to create workout plan", err)
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		workout_day_ids = append(workout_day_ids, workout_day.Id)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to commit transaction", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to commit transaction", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": gin.H{"ids": workout_day_ids}})
}

// 用于创建一个已经完成的训练，比如有氧
func (h *WorkoutDayHandler) CreateFreeWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "非法操作", "data": nil})
		return
	}
	var body struct {
		Type              string     `json:"type"`
		Title             string     `json:"title"`
		PendingSteps      string     `json:"pending_steps"`
		UpdatedDetails    string     `json:"updated_details"`
		StartAt           *time.Time `json:"start_at"`
		FinishedAt        *time.Time `json:"finished_at"`
		Remark            string     `json:"remark"` /** 备注 **/
		Medias            string     `json:"medias"`
		StartWhenCreate   bool       `json:"start_when_create"`
		FinishWhenCreated bool       `json:"finish_when_created"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Type == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少参数", "data": nil})
		return
	}
	if body.PendingSteps == "" || body.UpdatedDetails == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少参数", "data": nil})
		return
	}
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	var subscription models.Subscription
	if err := tx.Where("coach_id = ? AND step = 2", uid).First(&subscription).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 101, "msg": "该功能需订阅后才能使用", "data": nil})
		return
	}
	now := time.Now().UTC()

	workout_day := models.WorkoutDay{
		Title:          body.Title,
		Type:           body.Type,
		Status:         int(models.WorkoutDayStatusPending),
		PendingSteps:   body.PendingSteps,
		UpdatedDetails: body.UpdatedDetails,
		Remark:         body.Remark,
		Medias:         body.Medias,
		CreatedAt:      now,
		CoachId:        uid,
		StudentId:      uid,
	}
	if body.StartWhenCreate {
		workout_day.StartedAt = &now
		workout_day.Status = int(models.WorkoutDayStatusStarted)
	}
	if body.StartAt != nil {
		workout_day.StartedAt = body.StartAt
	}
	if body.FinishWhenCreated && body.FinishedAt != nil {
		workout_day.FinishedAt = body.FinishedAt
		workout_day.Status = int(models.WorkoutDayStatusFinished)
		// 计算训练时长（分钟）
		if body.Type == "cardio" && workout_day.StartedAt != nil {
			duration := int(body.FinishedAt.Sub(*workout_day.StartedAt).Minutes())
			workout_day.Duration = duration
		}
	}
	if err := tx.Create(&workout_day).Error; err != nil {
		h.logger.Error("Failed to create workout plan", err)
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	progress, err := models.ParseWorkoutDayProgress(body.PendingSteps)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	latest := models.ToWorkoutDayStepProgress(progress)
	total_volume := float64(0)
	for _, set := range latest.Sets {
		for _, act := range set.Actions {
			if act.Completed {
				history := models.WorkoutActionHistory{
					WorkoutDayId:    workout_day.Id,
					StudentId:       uid,
					WorkoutActionId: act.ActionId,
					Reps:            act.Reps,
					RepsUnit:        act.RepsUnit,
					Weight:          float64(act.Weight),
					WeightUnit:      act.WeightUnit,
					CreatedAt:       time.Unix(int64(act.CompletedAt), 0),
				}
				real_weight := float64(act.Weight)
				if act.WeightUnit == "磅" {
					real_weight = toFixed(real_weight*0.45, 1)
				}
				if act.RepsUnit == "次" {
					total_volume += float64(act.Reps) * real_weight
				}
				if err := tx.Create(&history).Error; err != nil {
					tx.Rollback()
					h.logger.Error("Failed to create workout action history", err)
				}
			}
		}
	}
	if total_volume != 0 {
		workout_day.TotalVolume = toFixed(total_volume, 1)
	}
	if err := tx.Save(&workout_day).Error; err != nil {
		h.logger.Error("Failed to save workout plan", err)
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to commit transaction", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to commit transaction", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": gin.H{
		"id": workout_day.Id,
	}})
}

func (h *WorkoutDayHandler) UpdateWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "非法操作", "data": nil})
		return
	}
	var body struct {
		Id             int        `json:"id"`
		Type           string     `json:"type"`
		Title          string     `json:"title"`
		PendingSteps   string     `json:"pending_steps"`
		UpdatedDetails string     `json:"updated_details"`
		StartAt        *time.Time `json:"start_at"`
		FinishedAt     *time.Time `json:"finished_at"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少参数", "data": nil})
		return
	}
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	var existing models.WorkoutDay
	if err := tx.Where("id = ? AND student_id = ?", body.Id, uid).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有找到记录", "data": nil})
		return
	}
	if body.Title != "" {
		existing.Title = body.Title
	}
	if body.Type != "" {
		existing.Type = body.Type
	}
	if body.UpdatedDetails != "" {
		existing.UpdatedDetails = body.UpdatedDetails
	}
	if body.PendingSteps != "" {
		existing.PendingSteps = body.PendingSteps
		if err := tx.Model(&models.WorkoutActionHistory{}).Where("workout_day_id = ?", existing.Id).Update("d", 1).Error; err != nil {
			tx.Rollback()
			h.logger.Error("Failed to update WorkoutActionHistory d field", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		progress, err := models.ParseWorkoutDayProgress(body.PendingSteps)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		latest := models.ToWorkoutDayStepProgress(progress)
		total_volume := float64(0)
		for _, set := range latest.Sets {
			for _, act := range set.Actions {
				fmt.Println("the action", act.Completed)
				if act.Completed {
					history := models.WorkoutActionHistory{
						WorkoutDayId:    existing.Id,
						StudentId:       uid,
						WorkoutActionId: act.ActionId,
						Reps:            act.Reps,
						RepsUnit:        act.RepsUnit,
						Weight:          float64(act.Weight),
						WeightUnit:      act.WeightUnit,
						CreatedAt:       time.Unix(int64(act.CompletedAt), 0),
					}
					real_weight := float64(act.Weight)
					if act.WeightUnit == "磅" {
						real_weight = toFixed(real_weight*0.45, 1)
					}
					if act.RepsUnit == "次" {
						total_volume += float64(act.Reps) * real_weight
					}
					if err := tx.Create(&history).Error; err != nil {
						tx.Rollback()
						h.logger.Error("Failed to create workout action history", err)
					}
				}
			}
		}
		if total_volume != 0 {
			existing.TotalVolume = toFixed(total_volume, 1)
		}
	}
	if err := tx.Save(&existing).Error; err != nil {
		h.logger.Error("Failed to save workout plan", err)
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to commit transaction", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to commit transaction", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "编辑成功", "data": gin.H{
		"id": existing.Id,
	}})
}

// 查看是否有进行中的训练，仅获取少量数据
func (h *WorkoutDayHandler) CheckHasStartedWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var list []models.WorkoutDay
	if err := h.db.
		Where("status = ?", int(models.WorkoutDayStatusStarted)).
		Where("coach_id = ?", uid).
		Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	data := []map[string]interface{}{}

	for _, v := range list {
		data = append(data, map[string]interface{}{
			"id":         v.Id,
			"status":     v.Status,
			"created_at": v.CreatedAt,
			"started_at": v.StartedAt,
			"student_id": v.StudentId,
			"coach_id":   v.CoachId,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "请求成功", "data": gin.H{
		"list": data,
	}})
}

func (h *WorkoutDayHandler) FetchStartedWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var list []models.WorkoutDay
	if err := h.db.
		Where("status = ?", int(models.WorkoutDayStatusStarted)).
		Where("coach_id = ? OR student_id = ?", uid, uid).
		Where("started_at IS NOT NULL").
		Order("started_at DESC").
		Preload("WorkoutPlan").
		Preload("Student.Profile1").
		Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	data := []map[string]interface{}{}

	for _, v := range list {
		vv := map[string]interface{}{
			"id":         v.Id,
			"status":     v.Status,
			"title":      v.Title,
			"type":       v.Type,
			"created_at": v.CreatedAt,
			"started_at": v.StartedAt,
			"coach_id":   v.CoachId,
			"group_no":   v.GroupNo,
			"student": map[string]interface{}{
				"id":         v.StudentId,
				"nickname":   v.Student.Profile1.Nickname,
				"avatar_url": v.Student.Profile1.AvatarURL,
			},
			"workout_plan": nil,
		}
		if v.WorkoutPlanId != 0 {
			vv["workout_plan"] = map[string]interface{}{
				"id":       v.WorkoutPlan.Id,
				"title":    v.WorkoutPlan.Title,
				"overview": v.WorkoutPlan.Overview,
			}
		}
		data = append(data, vv)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "请求成功", "data": gin.H{
		"list": data,
	}})
}

// 获取训练计划记录 获取自己的和当时一起训练好友、学员的
func (h *WorkoutDayHandler) FetchWorkoutDayProfile(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少 id 参数", "data": nil})
		return
	}
	var existing_workout_day models.WorkoutDay
	if err := h.db.
		Where("id = ?", body.Id).
		Preload("WorkoutPlan").
		Preload("WorkoutPlan.Creator.Profile1").
		Preload("Student.Profile1").
		First(&existing_workout_day).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有找到该记录", "data": nil})
		return
	}
	var existing_relation models.CoachRelationship
	if err := h.db.
		Where("(coach_id = ? AND student_id = ?) OR (coach_id = ? AND student_id = ?)", uid, existing_workout_day.StudentId, existing_workout_day.StudentId, uid).
		First(&existing_relation).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
	}
	if existing_workout_day.StudentId != uid && existing_relation.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "异常操作", "data": nil})
		return
	}
	data := gin.H{
		"id":           existing_workout_day.Id,
		"title":        existing_workout_day.Title,
		"type":         existing_workout_day.Type,
		"status":       existing_workout_day.Status,
		"duration":     existing_workout_day.Duration,
		"total_volume": existing_workout_day.TotalVolume,
		"remark":       existing_workout_day.Remark,
		"medias":       existing_workout_day.Medias,
		// 训练记录
		"pending_steps": existing_workout_day.PendingSteps,
		// 训练内容
		"updated_details": existing_workout_day.UpdatedDetails,
		"student_id":      existing_workout_day.StudentId,
		"student": gin.H{
			"id":         existing_workout_day.Student.Id,
			"nickname":   existing_workout_day.Student.Profile1.Nickname,
			"avatar_url": existing_workout_day.Student.Profile1.AvatarURL,
		},
		"is_self":      existing_workout_day.StudentId == uid,
		"started_at":   existing_workout_day.StartedAt,
		"finished_at":  existing_workout_day.FinishedAt,
		"workout_plan": nil,
	}
	if existing_workout_day.WorkoutPlanId != 0 {
		data["workout_plan"] = gin.H{
			"id":       existing_workout_day.WorkoutPlan.Id,
			"title":    existing_workout_day.WorkoutPlan.Title,
			"overview": existing_workout_day.WorkoutPlan.Overview,
			"tags":     existing_workout_day.WorkoutPlan.Tags,
			"details":  existing_workout_day.WorkoutPlan.Details,
			"creator": gin.H{
				"nickname":   existing_workout_day.WorkoutPlan.Creator.Profile1.Nickname,
				"avatar_url": existing_workout_day.WorkoutPlan.Creator.Profile1.AvatarURL,
			},
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": data})
}

// 获取训练计划记录结果 只能获取自己的
func (h *WorkoutDayHandler) FetchWorkoutDayResult(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少 id 参数", "data": nil})
		return
	}
	var workout_day models.WorkoutDay
	if err := h.db.
		Where("id = ? AND student_id = ?", body.Id, uid).
		Preload("WorkoutPlan").
		Preload("WorkoutPlan.Creator.Profile1").
		First(&workout_day).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有找到该记录", "data": nil})
		return
	}

	// Calculate the day number based on unique dates
	// var day_number int64
	// h.db.Model(&models.WorkoutDay{}).
	// 	Select("COUNT(DISTINCT DATE(created_at))").
	// 	Where("student_id = ? AND DATE(created_at) <= DATE(?)", uid, workout_day.CreatedAt).
	// 	Count(&day_number)
	result, err := models.BuildResultFromWorkoutDay(workout_day, h.db)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	data := gin.H{
		"id":     workout_day.Id,
		"status": workout_day.Status,
		// 训练记录
		// "pending_steps": workout_day.PendingSteps,
		// 训练内容
		// "updated_details": workout_day.UpdatedDetails,
		"steps":        result.List,
		"set_count":    result.SetCount,
		"duration":     result.DurationCount,
		"total_volume": result.TotalVolume,
		"tags":         result.Tags,
		"student_id":   workout_day.StudentId,
		"is_self":      workout_day.StudentId == uid,
		"workout_plan": nil,
		// "day_number":  day_number,
		"started_at":  workout_day.StartedAt,
		"finished_at": workout_day.FinishedAt,
	}
	if workout_day.WorkoutPlanId != 0 {
		data["workout_plan"] = gin.H{
			"id":       workout_day.WorkoutPlan.Id,
			"title":    workout_day.WorkoutPlan.Title,
			"overview": workout_day.WorkoutPlan.Overview,
			"tags":     workout_day.WorkoutPlan.Tags,
			// "details":  workout_day.WorkoutPlan.Details,
			"creator": gin.H{
				"nickname":   workout_day.WorkoutPlan.Creator.Profile1.Nickname,
				"avatar_url": workout_day.WorkoutPlan.Creator.Profile1.AvatarURL,
			},
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": data})
}

func (h *WorkoutDayHandler) FetchStudentWorkoutDayProfile(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id        int `json:"id"`
		StudentId int `json:"student_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少 id 参数", "data": nil})
		return
	}
	var existing_relation models.CoachRelationship
	if err := h.db.Where("(coach_id = ? AND student_id = ?) OR (coach_id = ? AND student_id = ?)", uid, body.StudentId, body.StudentId, uid).First(&existing_relation).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	query := h.db.Where("d IS NULL OR d = 0")
	query = query.Where("id = ? AND (coach_id = ? OR student_id = ?)", body.Id, body.StudentId, body.StudentId)
	var workout_day models.WorkoutDay
	if err := query.
		Preload("WorkoutPlan").
		Preload("WorkoutPlan.Creator.Profile1").
		Preload("Student.Profile1").
		First(&workout_day).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有找到该记录", "data": nil})
		return
	}

	// Calculate the day number based on unique dates
	// var day_number int64
	// h.db.Model(&models.WorkoutDay{}).
	// 	Select("COUNT(DISTINCT DATE(created_at))").
	// 	Where("student_id = ? AND DATE(created_at) <= DATE(?)", uid, workout_day.CreatedAt).
	// 	Count(&day_number)

	data := gin.H{
		"id":           workout_day.Id,
		"title":        workout_day.Title,
		"type":         workout_day.Type,
		"duration":     workout_day.Duration,
		"total_volume": workout_day.TotalVolume,
		"status":       workout_day.Status,
		"remark":       workout_day.Remark,
		"medias":       workout_day.Medias,
		// 训练记录
		"pending_steps": workout_day.PendingSteps,
		// 训练内容
		"updated_details": workout_day.UpdatedDetails,
		"student_id":      workout_day.StudentId,
		"student": gin.H{
			"id":         workout_day.Student.Id,
			"nickname":   workout_day.Student.Profile1.Nickname,
			"avatar_url": workout_day.Student.Profile1.AvatarURL,
		},
		"is_self": workout_day.StudentId == uid,
		// "day_number":  day_number,
		"started_at":   workout_day.StartedAt,
		"finished_at":  workout_day.FinishedAt,
		"workout_plan": nil,
	}
	if workout_day.WorkoutPlanId != 0 {
		data["workout_plan"] = gin.H{
			"id":       workout_day.WorkoutPlan.Id,
			"title":    workout_day.WorkoutPlan.Title,
			"overview": workout_day.WorkoutPlan.Overview,
			"tags":     workout_day.WorkoutPlan.Tags,
			"details":  workout_day.WorkoutPlan.Details,
			"creator": gin.H{
				"nickname":   workout_day.WorkoutPlan.Creator.Profile1.Nickname,
				"avatar_url": workout_day.WorkoutPlan.Creator.Profile1.AvatarURL,
			},
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": data})
}

// 获取学员/好友训练记录
// 可以不传 student_id，但是会匹配好友关系
func (h *WorkoutDayHandler) FetchStudentWorkoutDayResult(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id int `json:"id"`
		// StudentId int `json:"student_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少 id 参数", "data": nil})
		return
	}
	var workout_day models.WorkoutDay
	if err := h.db.Where("d IS NULL OR d = 0").Where("id = ?", body.Id).
		Preload("WorkoutPlan").
		Preload("WorkoutPlan.Creator.Profile1").
		Preload("Student.Profile1").
		First(&workout_day).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有找到该记录", "data": nil})
		return
	}

	fmt.Println(uid, workout_day.StudentId)
	var existing_relation models.CoachRelationship
	if err := h.db.Where("(coach_id = ? AND student_id = ?) OR (coach_id = ? AND student_id = ?)", uid, workout_day.StudentId, workout_day.StudentId, uid).First(&existing_relation).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "无法查看", "data": nil})
		return
	}
	result, err := models.BuildResultFromWorkoutDay(workout_day, h.db)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	data := gin.H{
		"id":           workout_day.Id,
		"status":       workout_day.Status,
		"steps":        result.List,
		"set_count":    result.SetCount,
		"duration":     result.DurationCount,
		"total_volume": result.TotalVolume,
		"tags":         result.Tags,
		"student_id":   workout_day.StudentId,
		"is_self":      workout_day.StudentId == uid,
		"workout_plan": nil,
		"started_at":   workout_day.StartedAt,
		"finished_at":  workout_day.FinishedAt,
	}
	if workout_day.WorkoutPlanId != 0 {
		data["workout_plan"] = gin.H{
			"id":       workout_day.WorkoutPlan.Id,
			"title":    workout_day.WorkoutPlan.Title,
			"overview": workout_day.WorkoutPlan.Overview,
			"tags":     workout_day.WorkoutPlan.Tags,
			"creator": gin.H{
				"nickname":   workout_day.WorkoutPlan.Creator.Profile1.Nickname,
				"avatar_url": workout_day.WorkoutPlan.Creator.Profile1.AvatarURL,
			},
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": data})
}

func (h *WorkoutDayHandler) UpdateWorkoutDayPlanDetails(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id   int    `json:"id"`
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少id参数", "data": nil})
		return
	}
	var existing models.WorkoutDay
	if err := h.db.Where("id = ? AND (coach_id = ? OR student_id = ?)", body.Id, uid, uid).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "没有找到记录", "data": nil})
		return
	}
	// existing.UpdatedDetails = body.Data
	// now := time.Now().UTC()
	// existing.UpdatedAt = &now
	if err := h.db.Model(&existing).Update("updated_details", body.Data).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": gin.H{"id": existing.Id}})
}

// 暂存的训练内容记录
func (h *WorkoutDayHandler) UpdateWorkoutDayStepProgress(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id   int    `json:"id"`
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少id参数", "data": nil})
		return
	}
	if body.Data == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少更新内容", "data": nil})
		return
	}
	var existing models.WorkoutDay
	if result := h.db.Where("id = ? AND (coach_id = ? OR student_id = ?)", body.Id, uid, uid).First(&existing); result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "没有找到记录", "data": nil})
		return
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()

	if err := tx.Model(&existing).Update("pending_steps", body.Data).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to update day", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction: ", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "更新失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": gin.H{"id": existing.Id}})
}

// 基本上用不上，都是用 createWorkoutDay
func (h *WorkoutDayHandler) StartWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少id参数", "data": nil})
		return
	}
	var day models.WorkoutDay
	if err := h.db.Where("student_id = ?", uid).First(&day, body.Id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "找不到记录", "data": nil})
		return
	}

	day.Status = int(models.WorkoutDayStatusStarted)
	now := time.Now().UTC()
	day.StartedAt = &now
	h.db.Save(&day)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "操作成功", "data": gin.H{"id": day.Id}})
}

func (h *WorkoutDayHandler) FinishWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id             int    `json:"id"`
		PendingSteps   string `json:"pending_steps"`
		UpdatedDetails string `json:"updated_details"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少id参数", "data": nil})
		return
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()

	var existing models.WorkoutDay
	if err := tx.Where("id = ? AND (coach_id = ? OR student_id = ?)", body.Id, uid, uid).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "没有找到记录", "data": nil})
		return
	}
	if body.PendingSteps != "" {
		existing.PendingSteps = body.PendingSteps
	}
	if body.UpdatedDetails != "" {
		existing.UpdatedDetails = body.UpdatedDetails
	}
	// 训练进行中的记录
	progress, err := models.ParseWorkoutDayProgress(existing.PendingSteps)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	latest := models.ToWorkoutDayStepProgress(progress)
	total_volume := float64(0)
	for _, set := range latest.Sets {
		for _, act := range set.Actions {
			if act.Completed {
				history := models.WorkoutActionHistory{
					WorkoutDayId:    body.Id,
					StudentId:       int(existing.StudentId),
					WorkoutActionId: act.ActionId,
					StepUid:         set.StepUid,
					SetUid:          set.Uid,
					ActUid:          act.Uid,
					Reps:            act.Reps,
					RepsUnit:        act.RepsUnit,
					Weight:          float64(act.Weight),
					WeightUnit:      act.WeightUnit,
					CreatedAt:       time.Unix(int64(act.CompletedAt), 0),
				}
				real_weight := float64(act.Weight)
				if act.WeightUnit == "磅" {
					real_weight = toFixed(real_weight*0.45, 1)
				}
				if act.RepsUnit == "次" {
					total_volume += float64(act.Reps) * real_weight
				}
				if err := tx.Create(&history).Error; err != nil {
					tx.Rollback()
					h.logger.Error("Failed to create workout action history", err)
				}
			}
		}
	}
	now := time.Now().UTC()
	now_trunc := now.Truncate(time.Minute)
	started_at_trunc := existing.StartedAt.Truncate(time.Minute)
	existing.Duration = int(now_trunc.Sub(started_at_trunc).Minutes())
	existing.FinishedAt = &now
	existing.Status = int(models.WorkoutDayStatusFinished)
	existing.TotalVolume = toFixed(total_volume, 1)
	if err := tx.Save(&existing).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to commit transaction", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to commit transaction", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "操作成功", "data": gin.H{"id": existing.Id}})
}

func (h *WorkoutDayHandler) GiveUpWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少id参数", "data": nil})
		return
	}
	var existing models.WorkoutDay
	if err := h.db.Where("id = ? AND (coach_id = ? OR student_id = ?)", body.Id, uid, uid).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "没有找到记录", "data": nil})
		return
	}
	existing.Status = int(models.WorkoutDayStatusGiveUp)
	h.db.Save(&existing)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": gin.H{"id": existing.Id}})
}

func (h *WorkoutDayHandler) ContinueWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少id参数", "data": nil})
		return
	}
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	var existing models.WorkoutDay
	if err := tx.Where("id = ? AND coach_id = ?", body.Id, uid).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "没有找到记录", "data": nil})
		return
	}
	existing.Status = int(models.WorkoutDayStatusStarted)
	if err := tx.Save(&existing).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	// 将 WorkoutActionHistory 所有 workout_day_id = existing.id 的记录 d 设置为 1
	if err := tx.Model(&models.WorkoutActionHistory{}).Where("workout_day_id = ?", existing.Id).Update("d", 1).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to update WorkoutActionHistory d field", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to commit transaction", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to commit transaction", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": gin.H{"id": existing.Id}})
}

func (h *WorkoutDayHandler) DeleteWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	var existing models.WorkoutDay
	if err := h.db.Where("id = ? AND student_id = ?", body.Id, uid).First(&existing).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if err := h.db.Where("id = ?", body.Id).Update("d", 1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "操作成功", "data": nil})
}

func (h *WorkoutDayHandler) FetchWorkoutDayList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		models.Pagination
		Status          int        `json:"status"`
		FinishedAtStart *time.Time `json:"finished_at_start"`
		FinishedAtEnd   *time.Time `json:"finished_at_end"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	query := h.db.Where("d IS NULL OR d = 0")
	if body.Status != 0 {
		query = query.Where("status = ?", body.Status)
	}
	// Add finished time range filter
	if body.FinishedAtStart != nil {
		query = query.Where("finished_at >= ?", body.FinishedAtStart)
	}
	if body.FinishedAtEnd != nil {
		query = query.Where("finished_at <= ?", body.FinishedAtEnd)
	}
	// student_id 表示是自己训练的记录
	query = query.Where("student_id = ?", uid)
	pb := pagination.NewPaginationBuilder[models.WorkoutDay](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetNextMarker(body.NextMarker).
		SetOrderBy("finished_at DESC")
	var list1 []models.WorkoutDay
	if err := pb.Build().Preload("WorkoutPlan").Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
	list := make([]map[string]interface{}, 0, len(list2))
	for _, v := range list2 {
		d := map[string]interface{}{
			"id":           v.Id,
			"status":       v.Status,
			"title":        v.Title,
			"type":         v.Type,
			"group_no":     v.GroupNo,
			"workout_plan": nil,
			"started_at":   v.StartedAt,
			"finished_at":  v.FinishedAt,
		}
		if v.WorkoutPlanId != 0 {
			d["workout_plan"] = map[string]interface{}{
				"id":       v.WorkoutPlan.Id,
				"title":    v.WorkoutPlan.Title,
				"overview": v.WorkoutPlan.Overview,
				"tags":     v.WorkoutPlan.Tags,
			}
		}
		list = append(list, d)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        list,
			"page_size":   pb.GetLimit(),
			"has_more":    has_more,
			"next_marker": next_marker,
		},
	})
}

// 似乎废弃了，使用 FetchWorkoutDayList 替代
func (h *WorkoutDayHandler) FetchFinishedWorkoutDayList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id              int        `json:"id"`
		FinishedAtStart *time.Time `json:"finished_at_start"`
		FinishedAtEnd   *time.Time `json:"finished_at_end"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	var list []models.WorkoutDay

	query := h.db.Where("d IS NULL OR d = 0")
	// student_id 表示是自己训练的记录
	query = query.Where("student_id = ?", uid)

	// 添加开始时间范围筛选
	if body.FinishedAtStart != nil {
		query = query.Where("finished_at >= ?", body.FinishedAtStart)
	}
	if body.FinishedAtEnd != nil {
		query = query.Where("finished_at <= ?", body.FinishedAtEnd)
	}

	query = query.Where("status = ?", int(models.WorkoutDayStatusFinished))
	query = query.Order("created_at desc")

	if err := query.Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	data := []map[string]interface{}{}

	for _, v := range list {
		data = append(data, map[string]interface{}{
			"id":          v.Id,
			"title":       v.Title,
			"type":        v.Type,
			"status":      v.Status,
			"started_at":  v.StartedAt,
			"finished_at": v.FinishedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list": data,
		},
	})
}

func (h *WorkoutDayHandler) FetchMyStudentWorkoutDayList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		models.Pagination
		Id             int        `json:"id"`
		StartedAtStart *time.Time `json:"started_at_start"`
		StartedAtEnd   *time.Time `json:"started_at_end"`
		Status         *int       `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	// 确保是自己的学员
	var relation models.CoachRelationship
	if err := h.db.Where("coach_id = ? AND student_id = ? OR coach_id = ? AND student_id = ?", uid, body.Id, body.Id, uid).First(&relation).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	query := h.db
	query = query.Where("student_id = ?", body.Id)
	// 添加开始时间范围筛选
	if body.StartedAtStart != nil {
		query = query.Where("started_at >= ?", body.StartedAtStart)
	}
	if body.StartedAtEnd != nil {
		query = query.Where("started_at <= ?", body.StartedAtEnd)
	}
	// 添加状态筛选
	if body.Status != nil {
		query = query.Where("status = ?", *body.Status)
	}
	pb := pagination.NewPaginationBuilder[models.WorkoutDay](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")

	var list []models.WorkoutDay
	if err := pb.Build().Preload("WorkoutPlan").Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list)
	data := []map[string]interface{}{}
	for _, v := range list2 {
		d := map[string]interface{}{
			"id":           v.Id,
			"status":       v.Status,
			"title":        v.Title,
			"type":         v.Type,
			"workout_plan": nil,
			"started_at":   v.StartedAt,
			"finished_at":  v.FinishedAt,
		}
		if v.WorkoutPlanId != 0 {
			d["workout_plan"] = map[string]interface{}{
				"id":       v.WorkoutPlan.Id,
				"title":    v.WorkoutPlan.Title,
				"overview": v.WorkoutPlan.Overview,
				"tags":     v.WorkoutPlan.Tags,
			}
		}
		data = append(data, d)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        data,
			"page_size":   pb.GetLimit(),
			"has_more":    has_more,
			"next_marker": next_marker,
		},
	})
}

func (h *WorkoutActionHistoryHandler) FetchStudentWorkoutActionHistoryListOfWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		models.Pagination
		WorkoutDayId int    `json:"workout_day_id"`
		StudentId    int    `json:"student_id"`
		OrderBy      string `json:"order_by"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	if body.WorkoutDayId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少参数", "data": nil})
		return
	}
	var d models.WorkoutDay
	query1 := h.db
	query1 = query1.Where("id = ? AND (coach_id = ? OR student_id = ?)", body.WorkoutDayId, uid, body.StudentId)
	if err := query1.First(&d).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	query := h.db.Where("d IS NULL OR d = 0")
	query = query.Where("workout_day_id = ?", body.WorkoutDayId)
	pb := pagination.NewPaginationBuilder[models.WorkoutActionHistory](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")

	var list1 []models.WorkoutActionHistory
	if err := pb.Build().Preload("WorkoutAction").Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout history: " + err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": gin.H{
			"list":        list2,
			"page_size":   pb.GetLimit(),
			"has_more":    has_more,
			"next_marker": next_marker,
		},
	})
}

func (h *WorkoutDayHandler) RefreshWorkoutDayRecords250630(c *gin.Context) {
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	// 查询所有 status=2（已完成）的 WorkoutDay 记录
	var days []models.WorkoutDay
	err := tx.Where("status = ?", int(models.WorkoutDayStatusFinished)).Preload("WorkoutPlan").Find(&days).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "查询失败: " + err.Error(), "data": nil})
		return
	}
	updated := 0
	for _, day := range days {
		if day.Title == "" {
			if err := tx.Model(&models.WorkoutDay{}).Where("id = ?", day.Id).Update("title", day.WorkoutPlan.Title).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
				return
			}
		}
		if day.Type == "" {
			if err := tx.Model(&models.WorkoutDay{}).Where("id = ?", day.Id).Update("type", day.WorkoutPlan.Type).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
				return
			}
		}
		if day.StartedAt != nil && day.FinishedAt != nil && day.FinishedAt.After(*day.StartedAt) {
			dur_sec := int(day.FinishedAt.Sub(*day.StartedAt).Seconds())
			// Duration 字段单位为分，四舍五入
			dur_min := (dur_sec + 30) / 60
			if day.Duration != dur_min {
				if err := tx.Model(&models.WorkoutDay{}).Where("id = ?", day.Id).Update("duration", dur_min).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
					return
				}
				updated++
			}
		}
		progress, err := models.ParseWorkoutDayProgress(day.PendingSteps)
		if err == nil {
			total_volume := float64(0)
			latest := models.ToWorkoutDayStepProgress(progress)
			for _, set := range latest.Sets {
				for _, act := range set.Actions {
					if act.Completed {
						weight := act.Weight
						if act.RepsUnit == "次" {
							if act.WeightUnit == "磅" {
								weight = act.Weight * 0.45
							}
							total_volume += float64(act.Reps) * weight
						}
					}
				}
			}
			if day.TotalVolume != float64(total_volume) {
				if err := tx.Model(&models.WorkoutDay{}).Where("id = ?", day.Id).Update("total_volume", total_volume).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
					return
				}
			}
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to commit transaction", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to commit transaction", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "刷新完成", "data": gin.H{"updated": updated, "total": len(days)}})
}

// 辅助函数
func derefInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func derefFloat(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

// toFixed 保留 n 位小数
func toFixed(num float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Round(num*pow) / pow
}
