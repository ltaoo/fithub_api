package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myapi/internal/models"
	"myapi/pkg/logger"
)

// WorkoutActionHandler handles HTTP requests for workout actions
type WorkoutActionHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewWorkoutActionHandler creates a new workout action handler
func NewWorkoutActionHandler(db *gorm.DB, logger *logger.Logger) *WorkoutActionHandler {
	return &WorkoutActionHandler{
		db:     db,
		logger: logger,
	}
}

// GetWorkoutActionList retrieves all workout actions
func (h *WorkoutActionHandler) GetWorkoutActionList(c *gin.Context) {
	var actions []models.WorkoutAction

	type WorkoutActionListRequest struct {
		Type       string `json:"type"`
		Keyword    string `json:"keyword"`
		Level      string `json:"level"`
		Tag        string `json:"tag"`
		Muscle     string `json:"muscle"`
		NextMarker string `json:"next_marker"`
		PageSize   int    `json:"page_size"`
		Page       int    `json:"page"`
	}

	var request WorkoutActionListRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// Get pagination parameters
	limit := 10 // default page size

	if request.PageSize != 0 {
		limit = request.PageSize
	}

	// Start with base query
	query := h.db

	// 按创建时间排序
	query = query.Order("created_at DESC")

	// Apply cursor pagination if cursor is provided
	if request.NextMarker != "" {
		if cursorId, err := strconv.Atoi(request.NextMarker); err == nil {
			fmt.Println("cursorId", cursorId)
			query = query.Where("id < ?", cursorId)
		}
	}
	if request.Page != 0 {
		query = query.Offset((request.Page - 1) * limit)
	}

	// Apply filters if provided
	if request.Type != "" {
		query = query.Where("type = ?", request.Type)
	}

	if request.Keyword != "" {
		query = query.Where("zh_name LIKE ? OR alias LIKE ?", "%"+request.Keyword+"%", "%"+request.Keyword+"%")
	}
	if request.Level != "" {
		levelInt, err := strconv.Atoi(request.Level)
		if err == nil {
			query = query.Where("level = ?", levelInt)
		}
	}

	if request.Tag != "" {
		query = query.Where("tags LIKE ?", "%"+request.Tag+"%")
	}

	if request.Muscle != "" {
		query = query.Where("target_muscle_ids LIKE ?", "%"+request.Muscle+"%")
	}

	// Add ordering and limit
	// query = query.Order("id asc").Limit(limit + 1)

	// Execute the query
	result := query.Find(&actions)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout actions", "data": nil})
		return
	}

	has_more := false
	next_cursor := ""

	// Check if there are more results
	if len(actions) > limit {
		has_more = true
		actions = actions[:limit]                            // Remove the extra item we fetched
		next_cursor = strconv.Itoa(int(actions[limit-1].Id)) // Get the last item's ID as next cursor
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        actions,
			"page_size":   limit,
			"has_more":    has_more,
			"next_marker": next_cursor,
		},
	})
}

// GetWorkoutActionList retrieves all workout actions
func (h *WorkoutActionHandler) GetWorkoutActionListByIds(c *gin.Context) {
	var actions []models.WorkoutAction

	type WorkoutActionListRequest struct {
		Ids []int `json:"ids"`
	}

	var request WorkoutActionListRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	// Start with base query
	query := h.db

	// 按创建时间排序
	query = query.Order("created_at DESC")

	// Apply filters if provided
	if len(request.Ids) > 0 {
		query = query.Where("id IN (?)", request.Ids)
	}
	// Execute the query
	result := query.Find(&actions)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout actions", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list": actions,
		},
	})
}

// GetWorkoutAction retrieves a specific workout action by ID
func (h *WorkoutActionHandler) GetWorkoutAction(c *gin.Context) {
	// id := c.Param("id")
	var request struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	var action models.WorkoutAction
	result := h.db.First(&action, request.Id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout action not found", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": action})
}

// GetActionsByMuscle retrieves workout actions targeting a specific muscle
func (h *WorkoutActionHandler) GetActionsByMuscle(c *gin.Context) {
	muscleID := c.Param("muscleId")

	var actions []models.WorkoutAction
	result := h.db.Where("target_muscle_ids LIKE ?", "%"+muscleID+"%").Find(&actions)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout actions", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": actions})
}

// GetActionsByLevel retrieves workout actions by difficulty level
func (h *WorkoutActionHandler) GetActionsByLevel(c *gin.Context) {
	levelStr := c.Param("level")
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid level parameter", "data": nil})
		return
	}

	var actions []models.WorkoutAction
	result := h.db.Where("level = ?", level).Find(&actions)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout actions", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": actions})
}

// GetRelatedActions retrieves advanced and regressed actions for a specific action
func (h *WorkoutActionHandler) GetRelatedActions(c *gin.Context) {
	id := c.Param("id")

	var action models.WorkoutAction
	result := h.db.First(&action, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout action not found", "data": nil})
		return
	}

	// Get advanced actions
	var advancedActions []models.WorkoutAction
	if action.AdvancedActionIds != "" {
		advancedIDs := strings.Split(action.AdvancedActionIds, ",")
		h.db.Where("id IN ?", advancedIDs).Find(&advancedActions)
	}

	// Get regressed actions
	var regressedActions []models.WorkoutAction
	if action.RegressedActionIds != "" {
		regressedIDs := strings.Split(action.RegressedActionIds, ",")
		h.db.Where("id IN ?", regressedIDs).Find(&regressedActions)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"advanced":  advancedActions,
			"regressed": regressedActions,
		},
	})
}

// CreateAction creates a new workout action
func (h *WorkoutActionHandler) CreateAction(c *gin.Context) {
	var action models.WorkoutAction

	if err := c.ShouldBindJSON(&action); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	result := h.db.Create(&action)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create workout action", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout action created successfully", "data": action})
}

// UpdateAction updates an existing workout action
func (h *WorkoutActionHandler) UpdateAction(c *gin.Context) {
	id := c.Param("id")

	var existingAction models.WorkoutAction
	result := h.db.First(&existingAction, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout action not found", "data": nil})
		return
	}

	var updatedAction models.WorkoutAction
	if err := c.ShouldBindJSON(&updatedAction); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// Ensure ID remains the same
	updatedAction.Id = existingAction.Id

	result = h.db.Save(&updatedAction)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update workout action", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout action updated successfully", "data": updatedAction})
}

// DeleteAction deletes a workout action
func (h *WorkoutActionHandler) DeleteAction(c *gin.Context) {
	id := c.Param("id")

	var action models.WorkoutAction
	result := h.db.First(&action, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout action not found", "data": nil})
		return
	}

	result = h.db.Delete(&action)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to delete workout action", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout action deleted successfully", "data": nil})
}
