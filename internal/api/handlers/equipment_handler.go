package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myapi/internal/models"
	"myapi/pkg/logger"
)

// MuscleHandler handles HTTP requests for muscles
type EquipmentHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewWorkoutActionHandler creates a new workout action handler
func NewEquipmentHandler(db *gorm.DB, logger *logger.Logger) *EquipmentHandler {
	return &EquipmentHandler{
		db:     db,
		logger: logger,
	}
}

// GetMuscles retrieves all muscles
func (h *EquipmentHandler) FetchEquipmentList(c *gin.Context) {
	var equipments []models.Equipment

	// Get query parameters for filtering
	type RequestPayload struct {
		Ids []int `json:"ids"`
	}
	var request RequestPayload
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// Start with base query
	query := h.db

	if len(request.Ids) != 0 {
		query = query.Where("id IN (?)", request.Ids)
	}

	result := query.Find(&equipments)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch equipments" + result.Error.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": gin.H{"list": equipments, "total": len(equipments)}})
}

// GetEquipment retrieves a specific equipment by ID
func (h *EquipmentHandler) GetEquipment(c *gin.Context) {
	id := c.Param("id")

	var equipment models.Equipment
	result := h.db.First(&equipment, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Equipment not found", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": equipment})
}

// CreateEquipment creates a new equipment
func (h *EquipmentHandler) CreateEquipment(c *gin.Context) {
	var equipment models.Equipment

	if err := c.ShouldBindJSON(&equipment); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	result := h.db.Create(&equipment)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create equipment", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Equipment created successfully", "data": equipment})
}

// UpdateEquipment updates an existing equipment
func (h *EquipmentHandler) UpdateEquipment(c *gin.Context) {
	id := c.Param("id")

	var existingEquipment models.Equipment
	result := h.db.First(&existingEquipment, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Equipment not found", "data": nil})
		return
	}

	var updatedEquipment models.Equipment
	if err := c.ShouldBindJSON(&updatedEquipment); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// Ensure ID remains the same
	updatedEquipment.Id = existingEquipment.Id

	result = h.db.Save(&updatedEquipment)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update equipment", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Equipment updated successfully", "data": updatedEquipment})
}

// DeleteEquipment deletes an existing equipment
func (h *EquipmentHandler) DeleteEquipment(c *gin.Context) {
	id := c.Param("id")

	var equipment models.Equipment
	result := h.db.First(&equipment, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Equipment not found", "data": nil})
		return
	}

	result = h.db.Delete(&equipment)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to delete equipment", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Equipment deleted successfully", "data": nil})
}
