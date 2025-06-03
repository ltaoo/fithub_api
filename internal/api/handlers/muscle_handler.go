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

// FetchMuscleList retrieves all muscles
func (h *MuscleHandler) FetchMuscleList(c *gin.Context) {

	var request struct {
		Ids []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// Start with base query
	query := h.db

	if len(request.Ids) != 0 {
		query = query.Where("id IN (?)", request.Ids)
	}

	var muscles []models.Muscle
	result := query.Find(&muscles).Order("created_at DESC")
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch muscles" + result.Error.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": gin.H{"list": muscles, "total": len(muscles)}})
}

// GetMuscleProfile retrieves a specific muscle by ID
func (h *MuscleHandler) GetMuscleProfile(c *gin.Context) {
	id := c.Param("id")

	var muscle models.Muscle
	result := h.db.First(&muscle, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Muscle not found", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": muscle})
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
	var body models.Muscle

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	var existing_muscle models.Muscle
	result := h.db.First(&existing_muscle, body.Id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Muscle not found", "data": nil})
		return
	}

	// Ensure ID remains the same
	body.Id = existing_muscle.Id

	payload := models.Muscle{
		Id:       body.Id,
		Name:     body.Name,
		ZhName:   body.ZhName,
		Tags:     body.Tags,
		Overview: body.Overview,
		Features: body.Features,
	}
	result = h.db.Save(&payload)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update muscle", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Muscle updated successfully", "data": nil})
}

// DeleteAction deletes a workout action
func (h *MuscleHandler) DeleteMuscle(c *gin.Context) {
	var body struct {
		Id int `json:"id"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	var muscle models.Muscle
	result := h.db.First(&muscle, body.Id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Muscle not found", "data": nil})
		return
	}

	result = h.db.Delete(&muscle)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to delete muscle", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Muscle deleted successfully", "data": nil})
}
