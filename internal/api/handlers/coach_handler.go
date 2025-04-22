package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"myapi/internal/models"
	"myapi/pkg/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Coach represents a coach in the system
type Coach struct {
	Id        int    `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Config    string `json:"config"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CoachAccount represents an authentication method for a coach
type CoachAccount struct {
	ProviderType string    `json:"provider_type"`
	ProviderID   string    `json:"provider_id"`
	ProviderArg1 int       `json:"provider_arg1,omitempty"`
	ProviderArg2 string    `json:"provider_arg2,omitempty"`
	ProviderArg3 string    `json:"provider_arg3,omitempty"`
	UserID       string    `json:"coach_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// CoachRegisterRequest represents the request body for coach registration
type CoachRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Type     string `json:"provider_type"`
	Nickname string `json:"nickname,omitempty"`
}

// CoachLoginRequest represents the request body for coach login
type CoachLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Type     string `json:"provider_type"`
}

// CoachHandler handles coach-related requests
type CoachHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewCoachHandler creates a new coach handler
func NewCoachHandler(db *gorm.DB, logger *logger.Logger) *CoachHandler {
	return &CoachHandler{
		db:     db,
		logger: logger,
	}
}

// RegisterCoach handles coach registration
func (h *CoachHandler) RegisterCoach(c *gin.Context) {
	var req CoachRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// Validate required fields
	if req.Email == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Email is required", "data": nil})
		return
	}

	// Check if email is already registered
	var count int
	row := h.db.Raw("SELECT COUNT(*) FROM COACH_ACCOUNT WHERE provider_type = ? AND provider_id = ?", req.Type, req.Email).Row()
	err := row.Scan(&count)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Database error", "data": nil})
		return
	}

	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 409, "msg": "Email already registered", "data": nil})
		return
	}

	if req.Type == "email_password" {

	} else {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Either password or verification code is required", "data": nil})
		return
	}

	// Create account record
	// tx1 := h.db.Exec(
	// 	"INSERT INTO COACH_ACCOUNT (provider_type, provider_id, provider_arg1, created_at) VALUES (?, ?, ?, ?)",
	// 	req.Type, req.Email, 0, now,
	// )
	// if tx1.Error != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create coach account"})
	// 	return
	// }

	// Create coach record
	config := "{}"
	now := time.Now().Format(time.RFC3339)

	if req.Nickname == "" {
		req.Nickname = req.Email
	}

	// Create a Coach struct to insert
	the_coach_created := struct {
		Nickname  string
		Config    string
		CreatedAt string
		UpdatedAt string
	}{
		Nickname:  req.Nickname,
		Config:    config,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Insert the coach record and get the ID
	result := h.db.Exec(
		"INSERT INTO COACH (nickname, config, created_at, updated_at) VALUES (?, ?, ?, ?)",
		the_coach_created.Nickname, the_coach_created.Config, the_coach_created.CreatedAt, the_coach_created.UpdatedAt,
	)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create coach", "data": nil})
		return
	}

	// Get the last inserted ID
	var coach_id int
	err = h.db.Raw("SELECT id FROM COACH WHERE nickname = ? AND created_at = ? ORDER BY id DESC LIMIT 1",
		the_coach_created.Nickname, the_coach_created.CreatedAt).Row().Scan(&coach_id)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to retrieve coach ID", "data": nil})
		return
	}

	// Determine authentication method and create account
	// var providerType string
	// var providerArg2 string

	// Create coach account
	// tx = h.db.Exec(
	// 	"INSERT INTO COACH_ACCOUNT (provider_type, provider_id, provider_arg2, coach_id, created_at) VALUES (?, ?, ?, ?, ?)",
	// 	providerType, req.Email, providerArg2, coachId, now,
	// )
	// if tx.Error != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create coach account"})
	// 	return
	// }
	// Create Account type is email_password

	// Email + Password authentication
	hashed_pwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to hash password", "data": nil})
		return
	}
	tx := h.db.Exec(
		"INSERT INTO COACH_ACCOUNT (provider_type, provider_id, provider_arg1, created_at, coach_id) VALUES (?, ?, ?, ?, ?)",
		req.Type, req.Email, string(hashed_pwd), time.Now().Format(time.RFC3339), coach_id,
	)
	if tx.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create coach account", "data": nil})
		return
	}

	// Generate JWT token
	token, err := models.GenerateJWT(coach_id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to generate token", "data": nil})
		return
	}

	// Return response
	// coach := Coach{
	// 	Id:        coachId,
	// 	Nickname:  req.Nickname,
	// 	Config:    config,
	// 	CreatedAt: time.Now().Format(time.RFC3339),
	// 	UpdatedAt: time.Now().Format(time.RFC3339),
	// }

	response := models.AuthResponse{
		Token:  "Bearer " + token,
		Status: "success",
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Coach registered successfully", "data": response})
}

// LoginCoach handles coach login
func (h *CoachHandler) LoginCoach(c *gin.Context) {
	var req CoachLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	if req.Email == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Email is required", "data": nil})
		return
	}

	var coach_id int
	var provider_arg1 string
	var coach Coach

	// Determine authentication method
	if req.Type == "email_password" {
		if req.Password == "" {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Password is required", "data": nil})
			return
		}
		row := h.db.Raw(
			"SELECT ca.coach_id, ca.provider_arg1 FROM COACH_ACCOUNT ca WHERE ca.provider_type = ? AND ca.provider_id = ?",
			req.Type, req.Email,
		).Row()

		err := row.Scan(&coach_id, &provider_arg1)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Invalid email or password", "data": nil})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Database error", "data": nil})
			return
		}

		// Verify password
		err = bcrypt.CompareHashAndPassword([]byte(provider_arg1), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Invalid email or password", "data": nil})
			return
		}
	} else if req.Type == "email_code" {
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Either password or verification code is required", "data": nil})
		return
	}

	// Get coach details
	row := h.db.Raw(
		"SELECT id, nickname, avatar_url, config, created_at, updated_at FROM COACH WHERE id = ?",
		coach_id,
	).Row()

	// 添加日志以便调试
	h.logger.Info("Attempting to fetch coach with ID: %v", coach_id)

	err := row.Scan(&coach.Id, &coach.Nickname, &coach.AvatarURL, &coach.Config, &coach.CreatedAt, &coach.UpdatedAt)
	if err != nil {
		h.logger.Error("Database error when fetching coach: %v", err)

		// 检查是否是因为找不到记录
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Coach not found", "data": nil})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to get coach details: " + err.Error(), "data": nil})
		}
		return
	}

	// Generate JWT token
	token, err := models.GenerateJWT(coach_id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to generate token", "data": nil})
		return
	}

	// Return response
	response := models.AuthResponse{
		Token:  "Bearer " + token,
		Status: "success",
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Login successful", "data": response})
}

// GetCoachProfile retrieves the coach's profile
func (h *CoachHandler) GetCoachProfile(c *gin.Context) {
	// Get coach ID from context (set by auth middleware)
	coachId, exists := c.Get("coachId")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Unauthorized", "data": nil})
		return
	}

	var coach Coach
	row := h.db.Raw(
		"SELECT id, nickname, avatar_url, config, created_at, updated_at FROM COACH WHERE id = ?",
		coachId,
	).Row()

	err := row.Scan(&coach.Id, &coach.Nickname, &coach.AvatarURL, &coach.Config, &coach.CreatedAt, &coach.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Coach not found", "data": nil})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Database error", "data": nil})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": coach})
}

// UpdateCoachProfile updates the coach's profile
func (h *CoachHandler) UpdateCoachProfile(c *gin.Context) {
	// Get coach ID from context (set by auth middleware)
	coachId, exists := c.Get("coachId")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Unauthorized", "data": nil})
		return
	}

	var updateData struct {
		Nickname  string `json:"nickname,omitempty"`
		AvatarURL string `json:"avatar_url,omitempty"`
		Config    string `json:"config,omitempty"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// Build update query
	query := "UPDATE COACH SET updated_at = ? "
	args := []interface{}{time.Now().Format(time.RFC3339)}

	if updateData.Nickname != "" {
		query += ", nickname = ? "
		args = append(args, updateData.Nickname)
	}

	if updateData.AvatarURL != "" {
		query += ", avatar_url = ? "
		args = append(args, updateData.AvatarURL)
	}

	if updateData.Config != "" {
		query += ", config = ? "
		args = append(args, updateData.Config)
	}

	query += "WHERE id = ?"
	args = append(args, coachId)

	// Execute update
	tx := h.db.Exec(query, args...)
	if tx.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update profile", "data": nil})
		return
	}

	// Get updated coach profile
	var coach Coach
	row := h.db.Raw(
		"SELECT id, nickname, avatar_url, config, created_at, updated_at FROM COACH WHERE id = ?",
		coachId,
	).Row()

	err := row.Scan(&coach.Id, &coach.Nickname, &coach.AvatarURL, &coach.Config, &coach.CreatedAt, &coach.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to get updated profile", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Profile updated successfully", "data": coach})
}

// SendVerificationCode sends a verification code to the specified email
func (h *CoachHandler) SendVerificationCode(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	if req.Email == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Email is required", "data": nil})
		return
	}

	// In a real implementation, you would generate and send a verification code
	// For now, we'll just return a success response

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "success",
		Message: "Verification code sent",
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Verification code sent", "data": response})
}
