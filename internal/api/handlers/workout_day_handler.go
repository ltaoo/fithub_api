package handlers

import (
	"encoding/json"
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
		StudentIds      []int `json:"student_ids"`
		StartWhenCreate bool  `json:"start_when_create"`
		WorkoutPlanId   int   `json:"workout_plan_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.WorkoutPlanId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少训练计划", "data": nil})
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
	if err := tx.Where("id = ?", body.WorkoutPlanId).First(&workout_plan).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 101, "msg": err.Error(), "data": nil})
		return
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
			if err := tx.Where("coach_id = ? AND student_id = ?", uid, student_id).First(&relation).Error; err != nil {
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
				CoachId:       relation.CoachId,
				StudentId:     relation.StudentId,
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
	}
	if err := tx.Create(&workout_day).Error; err != nil {
		h.logger.Error("Failed to create workout plan", err)
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	progress, err := ParseWorkoutDayProgress(body.PendingSteps)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	latest := ConvertToLatestProgress(progress)
	total_volume := float64(0)
	for _, set := range latest.Sets {
		for _, act := range set.Actions {
			fmt.Println("the action", act.Completed)
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
		progress, err := ParseWorkoutDayProgress(body.PendingSteps)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		latest := ConvertToLatestProgress(progress)
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
		Where("coach_id = ?", uid).
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

// 获取训练计划记录 只能获取自己的
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

	data := gin.H{
		"id":           workout_day.Id,
		"title":        workout_day.Title,
		"type":         workout_day.Type,
		"status":       workout_day.Status,
		"duration":     workout_day.Duration,
		"total_volume": workout_day.TotalVolume,
		"remark":       workout_day.Remark,
		"medias":       workout_day.Medias,
		// 训练记录
		"pending_steps": workout_day.PendingSteps,
		// 训练内容
		"updated_details": workout_day.UpdatedDetails,
		"student_id":      workout_day.StudentId,
		"is_self":         workout_day.StudentId == uid,
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

	data := gin.H{
		"id":     workout_day.Id,
		"status": workout_day.Status,
		// 训练记录
		"pending_steps": workout_day.PendingSteps,
		// 训练内容
		"updated_details": workout_day.UpdatedDetails,
		"student_id":      workout_day.StudentId,
		"is_self":         workout_day.StudentId == uid,
		"workout_plan": gin.H{
			"id":       workout_day.WorkoutPlan.Id,
			"title":    workout_day.WorkoutPlan.Title,
			"overview": workout_day.WorkoutPlan.Overview,
			"tags":     workout_day.WorkoutPlan.Tags,
			"details":  workout_day.WorkoutPlan.Details,
			"creator": gin.H{
				"nickname":   workout_day.WorkoutPlan.Creator.Profile1.Nickname,
				"avatar_url": workout_day.WorkoutPlan.Creator.Profile1.AvatarURL,
			},
		},
		// "day_number":  day_number,
		"started_at":  workout_day.StartedAt,
		"finished_at": workout_day.FinishedAt,
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
	query := h.db
	query = query.Where("id = ? AND coach_id = ?", body.Id, uid)
	if body.StudentId != 0 {
		query = query.Where("student_id = ?", body.StudentId)
	}
	var workout_day models.WorkoutDay
	if err := query.
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
		"is_self":         workout_day.StudentId == uid,
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
	if err := h.db.Where("id = ? AND coach_id = ?", body.Id, uid).First(&existing).Error; err != nil {
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
	if result := h.db.Where("id = ? AND coach_id = ?", body.Id, uid).First(&existing); result.Error != nil {
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
	if err := tx.Where("id = ? AND coach_id = ?", body.Id, uid).First(&existing).Error; err != nil {
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
	progress, err := ParseWorkoutDayProgress(existing.PendingSteps)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	latest := ConvertToLatestProgress(progress)
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
	if err := h.db.Where("id = ? AND coach_id = ?", body.Id, uid).First(&existing).Error; err != nil {
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
		list = append(list, map[string]interface{}{
			"id":          v.Id,
			"status":      v.Status,
			"title":       v.Title,
			"type":        v.Type,
			"group_no":    v.GroupNo,
			"started_at":  v.StartedAt,
			"finished_at": v.FinishedAt,
			"workout_plan": map[string]interface{}{
				"id":       v.WorkoutPlan.Id,
				"title":    v.WorkoutPlan.Title,
				"overview": v.WorkoutPlan.Overview,
				"tags":     v.WorkoutPlan.Tags,
			},
		})
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
	if err := h.db.Where("coach_id = ? AND student_id = ?", uid, body.Id).First(&relation).Error; err != nil {
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
		data = append(data, map[string]interface{}{
			"id":          v.Id,
			"status":      v.Status,
			"title":       v.Title,
			"type":        v.Type,
			"started_at":  v.StartedAt,
			"finished_at": v.FinishedAt,
			"workout_plan": map[string]interface{}{
				"id":       v.WorkoutPlan.Id,
				"title":    v.WorkoutPlan.Title,
				"overview": v.WorkoutPlan.Overview,
				"tags":     v.WorkoutPlan.Tags,
			},
		})
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
	query1 = query1.Where("id = ? AND coach_id = ?", body.WorkoutDayId, uid)
	if body.StudentId != 0 {
		query1 = query1.Where("student_id = ?", body.StudentId)
	}
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
	err := tx.Where("status = ?", int(models.WorkoutDayStatusFinished)).Find(&days).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "查询失败: " + err.Error(), "data": nil})
		return
	}
	updated := 0
	for _, day := range days {
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
		progress, err := ParseWorkoutDayProgress(day.PendingSteps)
		if err == nil {
			total_volume := float64(0)
			latest := ConvertToLatestProgress(progress)
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

type WorkoutDayProgress interface {
	GetVersion() string
}

func ParseWorkoutDayProgress(data string) (WorkoutDayProgress, error) {
	var versionHolder struct {
		V string `json:"v"`
	}
	if err := json.Unmarshal([]byte(data), &versionHolder); err != nil {
		return nil, err
	}
	switch versionHolder.V {
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
		return nil, fmt.Errorf("unknown version: %s", versionHolder.V)
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

func ConvertToLatestProgress(progress WorkoutDayProgress) WorkoutDayStepProgressJSON250629 {
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
