package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myapi/internal/models"
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
	var day struct {
		StartWhenCreate bool  `json:"start_when_create"`
		WorkoutPlanId   int64 `json:"workout_plan_id"`
		StudentId       int64 `json:"student_id"`
	}
	id := int64(c.GetFloat64("id"))
	if err := c.ShouldBindJSON(&day); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	var r models.WorkoutDay
	r.Status = int(models.WorkoutDayStatusPending)
	if day.WorkoutPlanId != 0 {
		r.WorkoutPlanId = day.WorkoutPlanId
	}
	r.StudentId = id
	if day.StudentId != 0 {
		r.StudentId = day.StudentId
	}
	r.CoachId = id
	now := time.Now().UTC()
	if day.StartWhenCreate {
		r.StartedAt = &now
		r.Status = int(models.WorkoutDayStatusStarted)
	}
	r.CreatedAt = now

	// Insert into database
	result := h.db.Create(&r)

	if result.Error != nil {
		// h.logger.Error("Failed to create workout plan", "error", result.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout day created successfully", "data": gin.H{"id": r.Id}})
}

type WorkoutDayId struct {
	Id int `json:"id"`
}

func (h *WorkoutDayHandler) FetchRunningWorkoutDay(c *gin.Context) {
	var plan models.WorkoutDay
	// 使用 Preload 预加载关联的 Steps 和 Actions
	result := h.db.
		Where("status = ?", int(models.WorkoutDayStatusStarted)).
		First(&plan)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout day not found", "data": nil})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout day retrieved successfully", "data": plan})
}

func (h *WorkoutDayHandler) FetchWorkoutDayProfile(c *gin.Context) {
	var id WorkoutDayId
	if err := c.ShouldBindJSON(&id); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	var workout_day models.WorkoutDay
	// 使用 Preload 预加载关联的 Steps 和 Actions
	result := h.db.
		Preload("WorkoutPlan").
		Preload("WorkoutPlan.Steps").
		Preload("WorkoutPlan.Steps.Actions").
		Preload("WorkoutPlan.Steps.Actions.Action").
		Where("id = ?", id.Id).
		First(&workout_day)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout day not found", "data": nil})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
		return
	}

	body := gin.H{
		"id":            workout_day.Id,
		"status":        workout_day.Status,
		"started_at":    workout_day.StartedAt,
		"pending_steps": workout_day.PendingSteps,
		"steps":         workout_day.WorkoutPlan.Steps,
		"workout_plan":  workout_day.WorkoutPlan,
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout day retrieved successfully", "data": body})
}

func (h *WorkoutDayHandler) StartWorkoutDay(c *gin.Context) {
	var id WorkoutDayId
	if err := c.ShouldBindJSON(&id); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	var day models.WorkoutDay
	if result := h.db.First(&day, id.Id); result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout day not found", "data": nil})
		return
	}

	day.Status = int(models.WorkoutDayStatusStarted)
	now := time.Now().UTC()
	day.StartedAt = &now
	h.db.Save(&day)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout day started successfully", "data": gin.H{"id": day.Id}})
}

func (h *WorkoutDayHandler) FinishWorkoutDay(c *gin.Context) {
	// var id WorkoutDayId
	var body struct {
		Id   int    `json:"id"`
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	var day models.WorkoutDay
	if result := h.db.First(&day, body.Id); result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout day not found", "data": nil})
		return
	}

	day.Status = int(models.WorkoutDayStatusFinished)
	if body.Data != "" {
		day.PendingSteps = body.Data
	}
	now := time.Now().UTC()
	day.FinishedAt = &now
	h.db.Save(&day)

	type WorkoutDayPendingAction struct {
		StepIdx     int     `json:"step_idx"`
		SetIdx      int     `json:"set_idx"`
		ActIdx      int     `json:"act_idx"`
		ActionId    int     `json:"action_id"`
		Reps        int     `json:"reps"`
		RepsUnit    string  `json:"reps_unit"`
		Weight      float64 `json:"weight"`
		WeightUnit  string  `json:"weight_unit"`
		Remark      string  `json:"remark"`
		Completed   bool    `json:"completed"`
		CompletedAt int     `json:"completed_at"`
	}

	type WorkoutDayPendingData struct {
		StepIdx int                       `json:"step_idx"`
		SetIdx  int                       `json:"set_idx"`
		Data    []WorkoutDayPendingAction `json:"data"`
	}

	var pending_data WorkoutDayPendingData
	if day.PendingSteps != "" {
		fmt.Println(day.PendingSteps)
		if err := json.Unmarshal([]byte(day.PendingSteps), &pending_data); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 600, "msg": "The PendingSteps has something wrong", "data": nil})
			return
		}
		for _, action := range pending_data.Data {
			if action.Completed {
				// Create workout action history record
				history := models.WorkoutActionHistory{
					WorkoutDayId: body.Id,
					StudentId:    int(day.StudentId),
					ActionId:     action.ActionId,
					Reps:         action.Reps,
					RepsUnit:     action.RepsUnit,
					Weight:       int(action.Weight),
					WeightUnit:   action.WeightUnit,
					Remark:       action.Remark,
					CreatedAt:    time.Unix(int64(action.CompletedAt), 0),
				}

				if result := h.db.Create(&history); result.Error != nil {
					h.logger.Error("Failed to create workout action history", result.Error)
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout day finished successfully", "data": gin.H{"id": day.Id}})
}

func (h *WorkoutDayHandler) GiveUpWorkoutDay(c *gin.Context) {
	var id WorkoutDayId
	if err := c.ShouldBindJSON(&id); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	var day models.WorkoutDay
	if result := h.db.First(&day, id.Id); result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout day not found", "data": nil})
		return
	}

	day.Status = int(models.WorkoutDayStatusGiveUp)
	h.db.Save(&day)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout day finished successfully", "data": gin.H{"id": day.Id}})
}

func (h *WorkoutDayHandler) UpdateWorkoutDaySteps(c *gin.Context) {
	// 绑定更新的数据
	var body struct {
		Id   int    `json:"id"`
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	// 先获取现有的计划
	var existing_day models.WorkoutDay
	if result := h.db.First(&existing_day, body.Id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout day not found", "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
		return
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	existing_day.PendingSteps = body.Data
	now := time.Now().UTC()
	existing_day.UpdatedAt = &now

	// 保存主记录
	if err := tx.Save(&existing_day).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update day: " + err.Error(), "data": nil})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to commit transaction: " + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout day updated successfully", "data": gin.H{"id": existing_day.Id}})
}

func (h *WorkoutDayHandler) DeleteWorkoutDay(c *gin.Context) {
	var id WorkoutDayId
	if err := c.ShouldBindJSON(&id); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	result := h.db.Where("id = ?", id.Id).Delete(&models.WorkoutDay{})
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to delete workout day", "data": nil})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout day not found", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout day deleted successfully", "data": nil})
}

// WorkoutDayFilter represents filter parameters for workout days
type WorkoutDayFilter struct {
	Title    string `form:"title"`
	Tags1    string `form:"tags1"`
	Tags2    string `form:"tags2"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

func (h *WorkoutDayHandler) FetchWorkoutDayList(c *gin.Context) {
	var plans []models.WorkoutDay

	// Get query parameters for filtering
	level := c.Query("level")
	tag := c.Query("tag")
	muscle := c.Query("muscle")
	cursor := c.Query("cursor")
	pageSize := c.Query("page_size")

	// Get pagination parameters
	limit := 10 // default page size

	if pageSize != "" {
		if parsedSize, err := strconv.Atoi(pageSize); err == nil && parsedSize > 0 {
			limit = parsedSize
		}
	}

	// Start with base query
	query := h.db

	// Apply cursor pagination if cursor is provided
	if cursor != "" {
		if cursorID, err := strconv.Atoi(cursor); err == nil {
			query = query.Where("id > ?", cursorID)
		}
	}

	// Apply filters if provided
	if level != "" {
		levelInt, err := strconv.Atoi(level)
		if err == nil {
			query = query.Where("level = ?", levelInt)
		}
	}

	if tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}

	if muscle != "" {
		query = query.Where("target_muscle_ids LIKE ?", "%"+muscle+"%")
	}

	// Add ordering and limit
	// query = query.Order("id asc").Limit(limit + 1)
	query = query.Order("created_at desc").Limit(limit + 1)

	// Execute the query
	result := query.Find(&plans)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout days", "data": nil})
		return
	}

	hasMore := false
	nextCursor := ""

	// Check if there are more results
	if len(plans) > limit {
		hasMore = true
		plans = plans[:limit]                             // Remove the extra item we fetched
		nextCursor = strconv.Itoa(int(plans[limit-1].Id)) // Get the last item's ID as next cursor
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        plans,
			"page_size":   limit,
			"has_more":    hasMore,
			"next_marker": nextCursor,
		},
	})
}
