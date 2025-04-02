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

// WorkoutPlanFilter represents filter parameters for workout plans
type WorkoutPlanFilter struct {
	Title    string `form:"title"`
	Tags1    string `form:"tags1"`
	Tags2    string `form:"tags2"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

func (h *WorkoutPlanHandler) CreateWorkoutPlan(c *gin.Context) {
	var plan models.WorkoutPlan

	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	// Validate required fields
	if plan.Title == "" || plan.Overview == "" {
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

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan created successfully", "data": plan})
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
	// 使用 Preload 预加载关联的 Sets 和 Actions
	result := h.db.
		Preload("Sets").
		Preload("Sets.Actions").
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
	var id WorkoutPlanId
	if err := c.ShouldBindJSON(&id); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	var plan models.WorkoutPlan
	result := h.db.Where("id = ?", id.Id).First(&plan)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout plan not found", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan updated successfully", "data": plan})
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

func (h *WorkoutPlanHandler) ListWorkoutPlans(c *gin.Context) {
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
	query = query.Order("id asc").Limit(limit + 1)

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
