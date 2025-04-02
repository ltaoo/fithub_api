package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myapi/internal/models"
	"myapi/pkg/logger"
)

// UserHandler 处理用户相关请求
type UserHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewUserHandler 创建新的用户处理器
func NewUserHandler(db *gorm.DB, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		db:     db,
		logger: logger,
	}
}

// GetUsers 获取所有用户
func (h *UserHandler) GetUsers(c *gin.Context) {
	var users []models.User
	result := h.db.Find(&users)
	if result.Error != nil {
		h.logger.Error("Failed to fetch users", result.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch users", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": users})
}

// GetUser 获取单个用户
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid user ID", "data": nil})
		return
	}

	var user models.User
	result := h.db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "User not found", "data": nil})
			return
		}
		h.logger.Error("Failed to fetch user", result.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch user", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": user})
}

// CreateUser 创建新用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	result := h.db.Create(&user)
	if result.Error != nil {
		h.logger.Error("Failed to create user", result.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create user", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "User created successfully", "data": user})
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid user ID", "data": nil})
		return
	}

	var user models.User
	if result := h.db.First(&user, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "User not found", "data": nil})
			return
		}
		h.logger.Error("Failed to fetch user for update", result.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch user", "data": nil})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	result := h.db.Save(&user)
	if result.Error != nil {
		h.logger.Error("Failed to update user", result.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update user", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "User updated successfully", "data": user})
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid user ID", "data": nil})
		return
	}

	result := h.db.Delete(&models.User{}, id)
	if result.Error != nil {
		h.logger.Error("Failed to delete user", result.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to delete user", "data": nil})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "User not found", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "User deleted successfully", "data": nil})
}
