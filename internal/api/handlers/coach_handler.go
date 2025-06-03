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
// type Coach struct {
// 	Id        int    `json:"id"`
// 	Nickname  string `json:"nickname"`
// 	AvatarURL string `json:"avatar_url,omitempty"`
// 	Config    string `json:"config"`
// 	CreatedAt string `json:"created_at"`
// 	UpdatedAt string `json:"updated_at"`
// }

// CoachAccount represents an authentication method for a coach
// type CoachAccount struct {
// 	ProviderType string    `json:"provider_type"`
// 	ProviderId   string    `json:"provider_id"`
// 	ProviderArg1 int       `json:"provider_arg1,omitempty"`
// 	ProviderArg2 string    `json:"provider_arg2,omitempty"`
// 	ProviderArg3 string    `json:"provider_arg3,omitempty"`
// 	UserId       string    `json:"coach_id"`
// 	CreatedAt    time.Time `json:"created_at"`
// }

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

	// CoachRegisterBody represents the request body for coach registration
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	if body.Email == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Email is required", "data": nil})
		return
	}

	var existing models.CoachAccount
	query := h.db.Where("provider_type = ? AND provider_id = ?", models.AccountProviderTypeEmailWithPwd, body.Email)
	r := query.First(&existing)
	if r.Error != nil {
		if r.Error != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": r.Error.Error(), "data": nil})
			return
		}
	}
	// Check if email is already registered
	if existing.CoachId != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "The account is existing", "data": nil})
		return
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()
	the_coach := models.Coach{
		Nickname:  body.Email,
		Config:    "{}",
		CreatedAt: now,
	}

	r2 := tx.Create(&the_coach)
	if r2.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": r2.Error.Error(), "data": nil})
		return
	}

	// Email + Password authentication
	hashed_pwd, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to hash password", "data": nil})
		return
	}

	the_coach_account := models.CoachAccount{
		ProviderType: models.AccountProviderTypeEmailWithPwd,
		ProviderId:   body.Email,
		ProviderArg1: string(hashed_pwd),
		CreatedAt:    now,
		CoachId:      the_coach.Id,
	}
	r3 := tx.Create(the_coach_account)
	if r3.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create coach account", "data": nil})
		return
	}

	// profile1 := models.CoachProfile1{
	// 	CoachId:  the_coach.Id,
	// 	Nickname: body.Email,
	// }
	// r4 := tx.Create(profile1)
	// if r4.Error != nil {
	// 	tx.Rollback()
	// 	c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create coach account", "data": nil})
	// 	return
	// }
	// profile2 := models.CoachProfile2{
	// 	CoachId: the_coach.Id,
	// }
	// r5 := tx.Create(profile2)
	// if r5.Error != nil {
	// 	tx.Rollback()
	// 	c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create coach account", "data": nil})
	// 	return
	// }

	if err := tx.Commit().Error; err != nil {
		// 提交失败（如数据库崩溃），可能需要处理补偿逻辑
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	// Generate JWT token
	token, err := models.GenerateJWT(the_coach.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to generate token", "data": nil})
		return
	}

	response := models.AuthResponse{
		Token:  "Bearer " + token,
		Status: "success",
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Coach registered successfully", "data": response})
}

// LoginCoach handles coach login
func (h *CoachHandler) LoginCoach(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	if body.Email == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Email is required", "data": nil})
		return
	}
	if body.Password == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Password is required", "data": nil})
		return
	}

	var account models.CoachAccount

	result := h.db.Where("provider_type = ? AND provider_id = ?", models.AccountProviderTypeEmailWithPwd, body.Email).First(&account)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "There no matched record", "data": nil})
		return
	}
	if account.CoachId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "There no matched record", "data": nil})
		return
	}

	// row := h.db.Raw(
	// 	"SELECT ca.coach_id, ca.provider_arg1 FROM COACH_ACCOUNT ca WHERE ca.provider_type = ? AND ca.provider_id = ?",
	// 	models.AccountProviderTypeEmailWithPwd, body.Email,
	// ).Row()

	// err := row.Scan(&coach_id, &provider_arg1)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Invalid email or password", "data": nil})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Database error", "data": nil})
	// 	return
	// }

	// Verify password
	err := bcrypt.CompareHashAndPassword([]byte(account.ProviderArg1), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Invalid email or password", "data": nil})
		return
	}

	// Get coach details
	// row = h.db.Raw(
	// 	"SELECT id, nickname, avatar_url, config, created_at, updated_at FROM COACH WHERE id = ?",
	// 	coach_id,
	// ).Row()

	// 添加日志以便调试
	// h.logger.Info("Attempting to fetch coach with ID: %v", coach_id)

	// err = row.Scan(&coach.Id, &coach.Nickname, &coach.AvatarURL, &coach.Config, &coach.CreatedAt, &coach.UpdatedAt)
	// if err != nil {
	// 	h.logger.Error("Database error when fetching coach: %v", err)

	// 	// 检查是否是因为找不到记录
	// 	if err == sql.ErrNoRows {
	// 		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Coach not found", "data": nil})
	// 	} else {
	// 		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to get coach details: " + err.Error(), "data": nil})
	// 	}
	// 	return
	// }

	// Generate JWT token
	token, err := models.GenerateJWT(account.CoachId)
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

	var coach models.Coach
	row := h.db.Raw(
		"SELECT id, nickname, avatar_url, config, created_at, updated_at FROM COACH WHERE id = ?",
		coachId,
	).Row()

	err := row.Scan(&coach.Id, &coach.Nickname, &coach.Config, &coach.CreatedAt, &coach.UpdatedAt)
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
	var coach models.Coach
	row := h.db.Raw(
		"SELECT id, nickname, avatar_url, config, created_at, updated_at FROM COACH WHERE id = ?",
		coachId,
	).Row()

	err := row.Scan(&coach.Id, &coach.Nickname, &coach.Config, &coach.CreatedAt, &coach.UpdatedAt)
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

func (h *CoachHandler) CreateStudent(c *gin.Context) {
	id := int(c.GetFloat64("id"))
	var body struct {
		Name   string `json:"name"`
		Gender int    `json:"gender"`
		Age    int    `json:"age"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	// Validate required fields
	if body.Name == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Nickname are required", "data": nil})
		return
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	now := time.Now().UTC()
	student := models.Coach{
		Nickname:  body.Name,
		CreatedAt: now,
		Profile1: models.CoachProfile1{
			Nickname: body.Name,
			Age:      body.Age,
			Gender:   body.Gender,
		},
	}

	r := tx.Create(&student)

	if r.Error != nil {
		// h.logger.Error("Failed to create workout plan", "error", result.Error)
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": r.Error.Error(), "data": nil})
		return
	}

	relationship := models.CoachRelationship{
		CoachId:   id,
		StudentId: int(student.Id),
		Status:    models.RelationPending,
		Role:      models.RoleCoachStudent,
		CreatedAt: now,
	}
	r2 := tx.Create(&relationship)
	if r2.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": r2.Error.Error(), "data": nil})
		return
	}

	if err := tx.Commit().Error; err != nil {
		// 提交失败（如数据库崩溃），可能需要处理补偿逻辑
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan created successfully", "data": gin.H{"id": student.Id}})
}

func (h *CoachHandler) FetchMyStudentList(c *gin.Context) {
	id := int(c.GetFloat64("id"))

	var body struct {
		models.Pagination
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	query := h.db
	query = query.Where("coach_id = ?", id)
	query = query.Order("created_at desc").Limit(body.PageSize + 1)
	query = query.Preload("Student").Preload("Student.Profile1")

	var students []models.CoachRelationship
	r := query.Find(&students)

	if r.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": r.Error.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        students,
			"page_size":   body.PageSize,
			"has_more":    false,
			"next_marker": "",
		},
	})
}

func (h *CoachHandler) FetchMyStudentProfile(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	var coach models.Coach
	result := h.db.
		Where("id = ?", body.Id).
		Preload("Profile1").
		Preload("Profile2").
		First(&coach)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "the record not found", "data": nil})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
		return
	}
	if body.Id != uid {
		var relation models.CoachRelationship
		r2 := h.db.Where("coach_id = ? AND student_id = ?", uid, body.Id).First(&relation)
		if r2.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "the record not found", "data": nil})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan retrieved successfully", "data": coach})
}

func (h *CoachHandler) FetchMyProfile(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var coach models.Coach
	r := h.db.
		Where("id = ?", uid).
		Preload("Profile1").
		First(&coach)
	if r.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "the record not found", "data": nil})
		return
	}
	if r.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": r.Error.Error(), "data": nil})
		return
	}

	var subscription models.Subscription
	r2 := h.db.Where("coach_id = ?", uid).Preload("SubscriptionPlan").First(&subscription)
	if r2.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": r2.Error.Error(), "data": nil})
		return
	}
	subscription_resp := gin.H{
		"visible": false,
		"text":    "unknown",
	}
	if subscription.Id != 0 {
		subscription_resp = gin.H{
			"visible": true,
			"text":    subscription.SubscriptionPlan.Name,
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan retrieved successfully", "data": gin.H{
		"nickname":     coach.Profile1.Nickname,
		"avatar_url":   coach.Profile1.AvatarURL,
		"subscription": subscription_resp,
	}})
}
