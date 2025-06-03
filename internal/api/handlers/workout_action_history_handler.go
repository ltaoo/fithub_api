package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myapi/internal/models"
	"myapi/pkg/logger"
)

// MuscleHandler handles HTTP requests for muscles
type WorkoutActionHistoryHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewWorkoutActionHandler creates a new workout action handler
func NewWorkoutActionHistoryHandler(db *gorm.DB, logger *logger.Logger) *WorkoutActionHistoryHandler {
	return &WorkoutActionHistoryHandler{
		db:     db,
		logger: logger,
	}
}

// 获取健身动作历史记录
func (h *WorkoutActionHistoryHandler) FetchWorkoutActionHistoryList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	type Body struct {
		models.Pagination
		WorkoutDayId    int    `json:"workout_day_id"`
		WorkoutActionId int    `json:"workout_action_id"`
		OrderBy         string `json:"order_by"`
	}
	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	offset := (body.Page - 1) * body.PageSize

	// Start with base query
	query := h.db.Preload("WorkoutAction")
	if body.WorkoutDayId != 0 {
		query = query.Where("workout_day_id = ?", body.WorkoutDayId)
	}
	if body.WorkoutActionId != 0 {
		query = query.Where("action_id = ?", body.WorkoutActionId)
	}

	query = query.Where("student_id = ?", uid)
	if body.WorkoutDayId != 0 {
		query = query.Order("created_at desc")
	}
	if body.WorkoutActionId != 0 {
		query = query.Order("weight desc")
	}
	query = query.Offset(offset).Limit(body.PageSize)

	var histories []models.WorkoutActionHistory
	if err := query.Find(&histories).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout history: " + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":      histories,
			"page":      body.Page,
			"page_size": body.PageSize,
			"has_more":  len(histories) >= body.PageSize,
		},
	})
}
