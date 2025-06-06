package handlers

import (
	"fmt"
	"net/http"
	"strconv"
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
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "请输入邮箱", "data": nil})
		return
	}

	var existing models.CoachAccount
	query := h.db.Where("provider_type = ? AND provider_id = ?", models.AccountProviderTypeEmailWithPwd, body.Email)
	if err := query.First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
	}
	// Check if email is already registered
	if existing.CoachId != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "该邮箱已被使用", "data": nil})
		return
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()

	nickname := func() string {
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		b := make([]byte, 6)
		for i := range b {
			b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		}
		return string(b)
	}()

	now := time.Now()
	the_coach := models.Coach{
		Nickname:  nickname,
		Config:    "{}",
		CreatedAt: now,
	}

	if err := tx.Create(&the_coach).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create coach", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	// Email + Password authentication
	hashed_pwd, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 501, "msg": "注册失败", "data": nil})
		return
	}

	the_coach_account := models.CoachAccount{
		ProviderType: models.AccountProviderTypeEmailWithPwd,
		ProviderId:   body.Email,
		ProviderArg1: string(hashed_pwd),
		CreatedAt:    now,
		CoachId:      the_coach.Id,
	}
	if err := tx.Create(&the_coach_account).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create coach account", err)
		c.JSON(http.StatusOK, gin.H{"code": 502, "msg": "注册失败", "data": nil})
		return
	}
	fmt.Println("-----------------")
	fmt.Println("About to create CoachProfile1")
	profile1 := models.CoachProfile1{
		CoachId:   the_coach.Id,
		Nickname:  nickname,
		AvatarURL: "",
	}
	fmt.Println("Created profile1 struct:", profile1)
	if err := tx.Create(&profile1).Error; err != nil {
		fmt.Println("Error creating profile1:", err)
		tx.Rollback()
		h.logger.Error("Failed to create profile1", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create coach account", "data": nil})
		return
	}
	fmt.Println("Successfully created profile1")
	fmt.Println("===================")
	the_coach.Profile1Id = profile1.Id
	if err := tx.Save(&the_coach).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create coach account", "data": nil})
		return
	}

	// if err := tx.Create(models.CoachProfile2{
	// 	CoachId: the_coach.Id,
	// }).Error; err != nil {
	// 	tx.Rollback()
	// 	h.logger.Error("Failed to create profile2", err)
	// 	c.JSON(http.StatusOK, gin.H{"code": 503, "msg": "注册失败", "data": nil})
	// 	return
	// }

	// Generate JWT token
	token, err := models.GenerateJWT(the_coach.Id)
	if err != nil {
		tx.Rollback()
		h.logger.Error("Failed to generate JWT", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to generate token", "data": nil})
		return
	}

	if err := tx.Commit().Error; err != nil {
		// 提交失败（如数据库崩溃），可能需要处理补偿逻辑
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
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

	// Verify password
	err := bcrypt.CompareHashAndPassword([]byte(account.ProviderArg1), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Invalid email or password", "data": nil})
		return
	}

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
	uid := int(c.GetFloat64("id"))

	var coach models.Coach

	if err := h.db.Where("id = ?", uid).Preload("Profile1").First(&coach).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "the profile not found", "data": nil})
		} else {
			h.logger.Error("Failed to find paper", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to find profile", "data": nil})
		}
		return
	}

	var subscription models.Subscription
	if err := h.db.Where("coach_id = ?", uid).Preload("SubscriptionPlan").First(&subscription).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to find profile", "data": nil})
			return
		}
	}

	resp := gin.H{
		"nickname":   coach.Profile1.Nickname,
		"avatar_url": coach.Profile1.AvatarURL,
		"subscription": gin.H{
			"visible": false,
			"text":    "unknown",
		},
	}
	if subscription.Id != 0 {
		resp["subscription"] = gin.H{
			"text":    subscription.SubscriptionPlan.Name,
			"visible": true,
		}
	}
	if coach.Profile1.Nickname == "" {
		resp["nickname"] = coach.Nickname
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": resp})
}

func (h *CoachHandler) UpdateCoachProfile(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
		Config    string `json:"config"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	if body.AvatarURL == "" && body.Nickname == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少修改参数", "data": nil})
		return
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	var user models.CoachProfile1
	if err := tx.Where("coach_id = ?", uid).First(&user).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to find record", err)
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "记录不存在", "data": nil})
		return
	}
	updates := map[string]interface{}{
		"nickname":   body.Nickname,
		"avatar_url": body.AvatarURL,
	}
	if body.Nickname != "" {
		updates["nickname"] = body.Nickname
	}
	if body.AvatarURL != "" {
		updates["avatar_url"] = body.AvatarURL
	}
	if err := h.db.Model(&user).Updates(updates).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to update the record", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update the record", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Profile updated successfully", "data": nil})
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
	uid := int(c.GetFloat64("id"))
	var body struct {
		Name   string `json:"name"`
		Gender int    `json:"gender"`
		Age    int    `json:"age"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	if body.Name == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Nickname are required", "data": nil})
		return
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered from panic:", r)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	now := time.Now().UTC()
	student := models.Coach{
		Nickname:  body.Name,
		CreatedAt: now,
	}
	if err := tx.Create(&student).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	profile1 := models.CoachProfile1{
		CoachId:   student.Id,
		Nickname:  body.Name,
		AvatarURL: "",
		Age:       body.Age,
		Gender:    body.Gender,
	}
	if err := tx.Create(&profile1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	student.Profile1Id = profile1.Id
	if err := tx.Save(&student).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	relationship := models.CoachRelationship{
		CoachId:   uid,
		StudentId: int(student.Id),
		Status:    models.RelationPending,
		Role:      models.RoleCoachStudent,
		CreatedAt: now,
	}
	if err := tx.Create(&relationship).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan created successfully", "data": gin.H{"id": student.Id}})
}

func (h *CoachHandler) DeleteMyStudent(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	var existing models.CoachRelationship
	if err := h.db.Where("coach_id = ? AND student_id = ?", uid, body.Id).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有找到记录", "data": nil})
		return
	}
	if err := h.db.Delete(&existing).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功", "data": nil})
}

func (h *CoachHandler) FetchMyStudentList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		models.Pagination
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	query := h.db
	limit := 20
	if body.PageSize != 0 {
		limit = body.PageSize
	}
	if body.Page != 0 {
		query = query.Offset((body.Page - 1) * limit)
	}
	query = query.Where("coach_id = ?", uid)
	query = query.Order("created_at desc").Limit(body.PageSize + 1)

	var list []models.CoachRelationship
	if err := query.Preload("Student").Preload("Student.Profile1").Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	has_more := false
	next_cursor := ""
	if len(list) > limit {
		has_more = true
		list = list[:limit]
		next_cursor = strconv.Itoa(int(list[limit-1].Id))
	}
	data_list := []interface{}{}
	for _, v := range list {
		data_list = append(data_list, map[string]interface{}{
			"id":         v.StudentId,
			"nickname":   v.Student.Profile1.Nickname,
			"avatar_url": v.Student.Profile1.AvatarURL,
			"age":        v.Student.Profile1.Age,
			"gender":     v.Student.Profile1.Gender,
		})
	}
	data := map[string]interface{}{
		"list":        data_list,
		"page_size":   limit,
		"has_more":    has_more,
		"next_marker": next_cursor,
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": data,
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
	var profile models.Coach
	if err := h.db.
		Where("id = ?", body.Id).
		Preload("Profile1").
		First(&profile).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有找到该记录", "data": nil})
		return
	}
	// var p models.CoachProfile1
	// if err := h.db.
	// 	Where("coach_id = ?", body.Id).
	// 	First(&p).Error; err != nil {
	// 	if err != gorm.ErrRecordNotFound {
	// 		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有找到该记录", "data": nil})
	// 	return
	// }
	var relation models.CoachRelationship
	if err := h.db.Where("coach_id = ? AND student_id = ?", uid, body.Id).First(&relation).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 501, "msg": "没有找到该记录", "data": nil})
		return
	}

	data := map[string]interface{}{
		"id":         profile.Id,
		"nickname":   profile.Profile1.Nickname,
		"avatar_url": profile.Profile1.AvatarURL,
		"gender":     profile.Profile1.Gender,
		"age":        profile.Profile1.Age,
	}
	// if profile.Profile1.Nickname == "" {
	// 	data["nickname"] = profile.Nickname
	// }

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan retrieved successfully", "data": data})
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

func (h *CoachHandler) RefreshToken(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	// Generate new JWT token
	token, err := models.GenerateJWT(uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to generate token", "data": nil})
		return
	}

	// Return response
	response := models.AuthResponse{
		Token:  "Bearer " + token,
		Status: "success",
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Token refreshed successfully", "data": response})
}
