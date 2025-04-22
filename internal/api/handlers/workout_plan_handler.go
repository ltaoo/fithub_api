package handlers

import (
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
type WorkoutPlanHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewWorkoutActionHandler creates a new workout action handler
func NewWorkoutPlanHandler(db *gorm.DB, logger *logger.Logger) *WorkoutPlanHandler {
	return &WorkoutPlanHandler{
		db:     db,
		logger: logger,
	}
}

func (h *WorkoutPlanHandler) CreateWorkoutPlan(c *gin.Context) {
	id := int(c.GetFloat64("id"))

	var plan models.WorkoutPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	// Validate required fields
	if plan.Title == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Title and overview are required", "data": nil})
		return
	}

	// 使用 UTC 时间
	plan.CreatedAt = time.Now().UTC()
	plan.OwnerId = id

	// Insert into database
	result := h.db.Create(&plan)

	if result.Error != nil {
		// h.logger.Error("Failed to create workout plan", "error", result.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan created successfully", "data": gin.H{"id": plan.Id}})
}

type WorkoutPlanId struct {
	Id int `json:"id"`
}

func (h *WorkoutPlanHandler) GetWorkoutPlan(c *gin.Context) {
	var id WorkoutPlanId
	if err := c.ShouldBindJSON(&id); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	var plan models.WorkoutPlan
	// 使用 Preload 预加载关联的 Steps 和 Actions
	result := h.db.
		Preload("Steps").
		Preload("Steps.Actions").
		Preload("Steps.Actions.Action").
		Where("id = ?", id.Id).
		First(&plan)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout plan not found", "data": nil})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan retrieved successfully", "data": plan})
}

func (h *WorkoutPlanHandler) UpdateWorkoutPlan(c *gin.Context) {
	// 绑定更新的数据
	var body models.WorkoutPlan
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	// 先获取现有的计划
	var existing_plan models.WorkoutPlan
	if result := h.db.First(&existing_plan, body.Id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout plan not found", "data": nil})
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

	// 更新基本字段
	existing_plan.Title = body.Title
	existing_plan.Overview = body.Overview
	existing_plan.Level = body.Level
	existing_plan.Tags = body.Tags
	existing_plan.Details = body.Details
	existing_plan.Points = body.Points
	existing_plan.Suggestions = body.Suggestions
	existing_plan.EquipmentIds = body.EquipmentIds
	existing_plan.MuscleIds = body.MuscleIds
	existing_plan.EstimatedDuration = body.EstimatedDuration
	existing_plan.Status = body.Status
	now := time.Now().UTC()
	existing_plan.UpdatedAt = &now

	// 保存主记录
	if err := tx.Save(&existing_plan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update plan: " + err.Error(), "data": nil})
		return
	}

	// 获取现有的 steps 映射，用于后续判断是更新还是新增
	existing_steps := make(map[int]*models.WorkoutPlanStep)
	existing_actions := make(map[int]*models.WorkoutPlanAction)
	var current_steps []models.WorkoutPlanStep
	if err := tx.Where("workout_plan_id = ?", existing_plan.Id).Find(&current_steps).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch existing steps: " + err.Error(), "data": nil})
		return
	}
	for i := range current_steps {
		step := current_steps[i]
		existing_steps[step.Id] = &step
		actions := []models.WorkoutPlanAction{}
		if err := tx.Where("workout_plan_step_id = ?", step.Id).Find(&actions).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch existing actions: " + err.Error(), "data": nil})
			return
		}
		for j := range actions {
			existing_actions[actions[j].Id] = &actions[j]
		}
	}
	fmt.Println("existing_steps", existing_steps)
	fmt.Println("existing_actions", existing_actions)

	// 记录处理过的 step IDs，用于后续删除不再需要的 steps
	processed_step_ids := make(map[int]bool)

	// 记录处理过的 action IDs，用于后续删除不再需要的 actions
	// processed_action_ids := make(map[int]bool)

	// 处理每个 step
	for i := range body.Steps {
		step := body.Steps[i]
		step.WorkoutPlanId = existing_plan.Id

		if step.Id > 0 {
			// 更新现有的 step
			if existing_step, ok := existing_steps[step.Id]; ok {
				// 更新基本字段
				existing_step.Title = step.Title
				existing_step.Type = step.Type
				existing_step.SetType = step.SetType
				existing_step.SetCount = step.SetCount
				existing_step.SetRestDuration = step.SetRestDuration
				existing_step.Note = step.Note

				if err := tx.Save(existing_step).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update step: " + err.Error(), "data": nil})
					return
				}
				processed_step_ids[step.Id] = true
			}
		} else {
			// 创建新的 step
			if err := tx.Create(&step).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create step: " + err.Error(), "data": nil})
				return
			}
		}

		// 处理 actions
		for j := range step.Actions {
			action := step.Actions[j]
			action.WorkoutPlanStepId = step.Id
			if action.Id > 0 {
				if existing_action, ok := existing_actions[action.Id]; ok {
					// 更新现有的 action
					existing_action.ActionId = action.ActionId
					existing_action.SetIdx = action.SetIdx
					existing_action.Reps = action.Reps
					existing_action.RepsUnit = action.RepsUnit
					existing_action.Weight = action.Weight
					existing_action.Tempo = action.Tempo
					existing_action.RestDuration = action.RestDuration
					existing_action.Note = action.Note
					delete(existing_actions, action.Id)

					if err := tx.Save(existing_action).Error; err != nil {
						tx.Rollback()
						c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update action: " + err.Error(), "data": nil})
						return
					}
				} else {
					// fmt.Println("action.Id need delete", action.Id)
					// processed_action_ids[action.Id] = true
				}
			} else {
				// 创建新的 action
				if err := tx.Create(&action).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create action: " + err.Error(), "data": nil})
					return
				}
			}
		}
	}

	// 删除不再需要的 steps（及其关联的 actions，通过外键级联删除）
	for stepId := range existing_steps {
		if !processed_step_ids[stepId] {
			if err := tx.Delete(&models.WorkoutPlanStep{}, stepId).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to delete unused step: " + err.Error(), "data": nil})
				return
			}
		}
	}

	fmt.Println("existing_actions", existing_actions)
	// 删除不再需要的 actions
	for actionId := range existing_actions {
		if err := tx.Delete(&models.WorkoutPlanAction{}, actionId).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to delete unused action: " + err.Error(), "data": nil})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to commit transaction: " + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan updated successfully", "data": gin.H{"id": existing_plan.Id}})
}

func (h *WorkoutPlanHandler) DeleteWorkoutPlan(c *gin.Context) {
	var id WorkoutPlanId
	if err := c.ShouldBindJSON(&id); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	result := h.db.Where("id = ?", id.Id).Delete(&models.WorkoutPlan{})
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to delete workout plan", "data": nil})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout plan not found", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan deleted successfully", "data": nil})
}

func (h *WorkoutPlanHandler) FetchWorkoutPlanList(c *gin.Context) {
	var plans []models.WorkoutPlan

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
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout plans", "data": nil})
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

func (h *WorkoutPlanHandler) FetchMyWorkoutPlanList(c *gin.Context) {
	var plans []models.WorkoutPlan

	id := int(c.GetFloat64("id"))

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

	query = query.Where("owner_id = ?", id)

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
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout plans", "data": nil})
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
