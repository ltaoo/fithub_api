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

// FetchMuscleList retrieves all muscles
func (h *WorkoutActionHistoryHandler) FetchWorkoutActionHistoryList(c *gin.Context) {
	var histories []models.WorkoutActionHistory

	id, existing := c.Get("id")
	if !existing {
		c.JSON(http.StatusOK, gin.H{"code": 900, "msg": "Please login", "data": nil})
		return
	}
	type Body struct {
		models.Pagination
	}
	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	// Start with base query
	query := h.db.Preload("Action")

	query = query.Where("student_id = ?", id)
	query = query.Order("created_at desc").Limit(body.PageSize + 1)

	result := query.Find(&histories)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch muscles" + result.Error.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": gin.H{"list": histories, "total": len(histories)}})
}
