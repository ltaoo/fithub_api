package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

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

// FetchWorkoutActionList retrieves all workout actions
func (h *WorkoutActionHandler) FetchWorkoutActionList(c *gin.Context) {

	var body struct {
		Type       string `json:"type"`
		Keyword    string `json:"keyword"`
		Tag        string `json:"tag"`
		Level      string `json:"level"`
		Muscle     string `json:"muscle"`
		NextMarker string `json:"next_marker"`
		PageSize   int    `json:"page_size"`
		Page       int    `json:"page"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// Get pagination parameters
	limit := 20 // default page size

	if body.PageSize != 0 {
		limit = body.PageSize
	}

	// Start with base query
	query := h.db

	// Apply cursor pagination if cursor is provided
	if body.NextMarker != "" {
		if cursor_id, err := strconv.Atoi(body.NextMarker); err == nil {
			query = query.Where("id < ?", cursor_id)
		}
	}
	if body.Page != 0 {
		query = query.Offset((body.Page - 1) * limit)
	}

	// Apply filters if provided
	if body.Type != "" {
		query = query.Where("type = ?", body.Type)
	}
	if body.Keyword != "" {
		query = query.Where("zh_name LIKE ? OR alias LIKE ?", "%"+body.Keyword+"%", "%"+body.Keyword+"%")
	}
	if body.Level != "" {
		levelInt, err := strconv.Atoi(body.Level)
		if err == nil {
			query = query.Where("level = ?", levelInt)
		}
	}
	if body.Tag != "" {
		query = query.Where("tags1 LIKE ?", "%"+body.Tag+"%")
	}
	if body.Muscle != "" {
		query = query.Where("target_muscle_ids LIKE ?", "%"+body.Muscle+"%")
	}
	// 按创建时间排序
	// query = query.Order("created_at DESC")
	// Add ordering and limit
	query = query.Order("sort_idx desc").Limit(limit + 1)

	var actions []models.WorkoutAction
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

func (h *WorkoutActionHandler) FetchWorkoutActionsByLevel(c *gin.Context) {
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

// 获取指定动作的 进阶、退阶、替代动作
func (h *WorkoutActionHandler) FetchRelatedWorkoutActions(c *gin.Context) {
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	var action models.WorkoutAction
	r := h.db.First(&action, body.Id)
	if r.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout action not found", "data": nil})
		return
	}
	// Get advanced actions
	var advanced_workout_actions []models.WorkoutAction
	if action.AdvancedActionIds != "" {
		advanced_ids := strings.Split(action.AdvancedActionIds, ",")
		h.db.Where("id IN ?", advanced_ids).Find(&advanced_workout_actions)
	}
	// Get regressed actions
	var regressed_actions []models.WorkoutAction
	if action.RegressedActionIds != "" {
		regressed_ids := strings.Split(action.RegressedActionIds, ",")
		h.db.Where("id IN ?", regressed_ids).Find(&regressed_actions)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"advanced":  advanced_workout_actions,
			"regressed": regressed_actions,
		},
	})
}

// 创建动作
func (h *WorkoutActionHandler) CreateWorkoutAction(c *gin.Context) {
	var body models.WorkoutAction
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	result := h.db.Create(&body)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "the record created successfully", "data": body})
}

// 更新动作详情
func (h *WorkoutActionHandler) UpdateWorkoutActionProfile(c *gin.Context) {
	var body models.WorkoutAction
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	var existing_action models.WorkoutAction
	result := h.db.First(&existing_action, body.Id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "the record not found", "data": nil})
		return
	}
	now := time.Now().UTC()
	result = h.db.Model(&existing_action).Updates(models.WorkoutAction{
		Name:                 body.Name,
		ZhName:               body.ZhName,
		Alias:                body.Alias,
		Overview:             body.Overview,
		SortIdx:              body.SortIdx,
		Type:                 body.Type,
		Level:                body.Level,
		Score:                body.Score,
		Tags1:                body.Tags1,
		Tags2:                body.Tags2,
		Pattern:              body.Pattern,
		Details:              body.Details,
		Points:               body.Points,
		Problems:             body.Problems,
		ExtraConfig:          body.ExtraConfig,
		UpdatedAt:            &now,
		MuscleIds:            body.MuscleIds,
		PrimaryMuscleIds:     body.PrimaryMuscleIds,
		SecondaryMuscleIds:   body.SecondaryMuscleIds,
		EquipmentIds:         body.EquipmentIds,
		AlternativeActionIds: body.AlternativeActionIds,
		AdvancedActionIds:    body.AdvancedActionIds,
		RegressedActionIds:   body.RegressedActionIds,
	})
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "updated successfully", "data": body})
}

// 删除动作
func (h *WorkoutActionHandler) DeleteWorkoutAction(c *gin.Context) {
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	var action models.WorkoutAction
	r := h.db.First(&action, body.Id)
	if r.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "the record not found", "data": nil})
		return
	}
	r = h.db.Delete(&action)
	if r.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to delete the record", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "the record deleted successfully", "data": nil})
}
