package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myapi/internal/models"
	"myapi/pkg/logger"
)

// MuscleHandler handles HTTP requests for muscles
type MuscleHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewWorkoutActionHandler creates a new workout action handler
func NewMuscleHandler(db *gorm.DB, logger *logger.Logger) *MuscleHandler {
	return &MuscleHandler{
		db:     db,
		logger: logger,
	}
}

// GetMuscles retrieves all muscles
func (h *MuscleHandler) GetMuscles(c *gin.Context) {
	var muscles []models.Muscle

	// Get query parameters for filtering

	// Start with base query
	query := h.db

	result := query.Find(&muscles)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch muscles" + result.Error.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": gin.H{"list": muscles, "total": len(muscles)}})
}

// GetMuscle retrieves a specific muscle by ID
func (h *MuscleHandler) GetMuscle(c *gin.Context) {
	id := c.Param("id")

	var muscle models.Muscle
	result := h.db.First(&muscle, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Muscle not found", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": muscle})
}

// GetActionsByMuscle retrieves workout actions targeting a specific muscle
func (h *MuscleHandler) GetActionsByMuscle(c *gin.Context) {
	muscleID := c.Param("muscleId")

	var muscles []models.Muscle
	result := h.db.Where("target_muscle_ids LIKE ?", "%"+muscleID+"%").Find(&muscles)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout actions", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": muscles})
}

// CreateAction creates a new workout action
func (h *MuscleHandler) CreateMuscle(c *gin.Context) {
	var muscle models.Muscle

	if err := c.ShouldBindJSON(&muscle); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	result := h.db.Create(&muscle)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create muscle", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Muscle created successfully", "data": muscle})
}

// UpdateAction updates an existing workout action
func (h *MuscleHandler) UpdateMuscle(c *gin.Context) {
	id := c.Param("id")

	var existingMuscle models.Muscle
	result := h.db.First(&existingMuscle, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Muscle not found", "data": nil})
		return
	}

	var updatedMuscle models.Muscle
	if err := c.ShouldBindJSON(&updatedMuscle); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// Ensure ID remains the same
	updatedMuscle.Id = existingMuscle.Id

	result = h.db.Save(&updatedMuscle)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update muscle", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Muscle updated successfully", "data": updatedMuscle})
}

// DeleteAction deletes a workout action
func (h *MuscleHandler) DeleteMuscle(c *gin.Context) {
	id := c.Param("id")

	var action models.WorkoutAction
	result := h.db.First(&action, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Muscle not found", "data": nil})
		return
	}

	result = h.db.Delete(&action)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to delete muscle", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Muscle deleted successfully", "data": nil})
}
