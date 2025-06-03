package handlers

import (
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

func (h *WorkoutPlanHandler) FetchWorkoutPlanProfile(c *gin.Context) {
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	var plan models.WorkoutPlan
	// 使用 Preload 预加载关联的 Steps 和 Actions
	result := h.db.
		Where("id = ?", body.Id).
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

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan updated successfully", "data": gin.H{"id": existing_plan.Id}})
}

func (h *WorkoutPlanHandler) DeleteWorkoutPlan(c *gin.Context) {
	var id struct {
		Id int `json:"id"`
	}
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

	has_more := false
	next_cursor := ""

	// Check if there are more results
	if len(plans) > limit {
		has_more = true
		plans = plans[:limit]                              // Remove the extra item we fetched
		next_cursor = strconv.Itoa(int(plans[limit-1].Id)) // Get the last item's ID as next cursor
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        plans,
			"page_size":   limit,
			"has_more":    has_more,
			"next_marker": next_cursor,
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

func (h *WorkoutPlanHandler) CreateWorkoutPlanCollection(c *gin.Context) {
	id := int(c.GetFloat64("id"))

	var collection models.WorkoutPlanCollection
	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	// Validate required fields
	if collection.Title == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Title are required", "data": nil})
		return
	}

	// 使用 UTC 时间
	collection.CreatedAt = time.Now().UTC()
	collection.OwnerId = id

	// Insert into database
	result := h.db.Create(&collection)

	if result.Error != nil {
		// h.logger.Error("Failed to create workout plan", "error", result.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan created successfully", "data": gin.H{"id": collection.Id}})
}

func (h *WorkoutPlanHandler) FetchWorkoutPlanCollectionProfile(c *gin.Context) {
	var id struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&id); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	var plan models.WorkoutPlanCollection

	result := h.db.
		Preload("WorkoutPlans").
		Preload("WorkoutPlans.WorkoutPlan").
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

func (h *WorkoutPlanHandler) FetchWorkoutPlanSetList(c *gin.Context) {
	var plans []models.WorkoutPlanSet

	cursor := c.Query("cursor")
	pageSize := c.Query("page_size")

	// Get pagination parameters
	limit := 100 // default page size

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

	// Add ordering and limit
	// query = query.Order("id asc").Limit(limit + 1)
	query = query.Order("created_at desc").Limit(limit + 1)

	// Execute the query
	result := query.Find(&plans)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout plans", "data": nil})
		return
	}

	has_more := false
	next_cursor := ""

	// Check if there are more results
	if len(plans) > limit {
		has_more = true
		plans = plans[:limit]                              // Remove the extra item we fetched
		next_cursor = strconv.Itoa(int(plans[limit-1].Id)) // Get the last item's ID as next cursor
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        plans,
			"page_size":   limit,
			"has_more":    has_more,
			"next_marker": next_cursor,
		},
	})
}

func (h *WorkoutPlanHandler) CreateWorkoutPlanSet(c *gin.Context) {
	// id := int(c.GetFloat64("id"))

	var plan models.WorkoutPlanSet
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

	// Insert into database
	result := h.db.Create(&plan)

	if result.Error != nil {
		// h.logger.Error("Failed to create workout plan", "error", result.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan created successfully", "data": gin.H{"id": plan.Id}})
}

func (h *WorkoutPlanHandler) UpdateWorkoutPlanSet(c *gin.Context) {
	// 绑定更新的数据
	var body models.WorkoutPlanSet
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	// 先获取现有的计划
	var existing_plan models.WorkoutPlanSet
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
	existing_plan.Idx = body.Idx
	existing_plan.Details = body.Details

	// 保存主记录
	if err := tx.Save(&existing_plan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update plan: " + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan updated successfully", "data": gin.H{"id": existing_plan.Id}})
}
