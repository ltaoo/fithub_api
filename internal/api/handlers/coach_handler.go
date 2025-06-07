package handlers

import (
	"fmt"
	"net/http"
	"time"

	"myapi/internal/models"
	"myapi/internal/pkg/pagination"
	"myapi/internal/pkg/sensitive"
	"myapi/pkg/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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

var AvatarPrefix = "//static.fithub.top/avatars/"
var DefaultAvatarURL = AvatarPrefix + "default1.jpeg"

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
	profile1 := models.CoachProfile1{
		CoachId:   the_coach.Id,
		Nickname:  nickname,
		AvatarURL: DefaultAvatarURL,
	}
	if err := tx.Create(&profile1).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create profile1", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create coach account", "data": nil})
		return
	}
	the_coach.Profile1Id = profile1.Id
	if err := tx.Save(&the_coach).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create coach account", "data": nil})
		return
	}
	var subscription_plan models.SubscriptionPlan
	if err := tx.First(&subscription_plan).Error; err == nil {
		day_count := 30
		expired_at := now.AddDate(0, 0, day_count)
		subscription := models.Subscription{
			Step:               2,
			Count:              30,
			ExpectExpiredAt:    &expired_at,
			Reason:             "注册赠送",
			SubscriptionPlanId: subscription_plan.Id,
			CoachId:            the_coach.Id,
			ActiveAt:           &now,
			CreatedAt:          now,
		}
		if err := tx.Create(&subscription).Error; err != nil {
			h.logger.Error("Failed to create subscription", err)
		}
	}

	// Generate JWT token
	token, err := models.GenerateJWT(the_coach.Id)
	if err != nil {
		tx.Rollback()
		h.logger.Error("Failed to generate JWT", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to generate token", "data": nil})
		return
	}

	if err := tx.Commit().Error; err != nil {
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

func (h *CoachHandler) FetchCoachProfile(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	var coach models.Coach
	if err := tx.
		Where("id = ?", uid).
		Preload("Profile1").
		First(&coach).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "没有找到记录", "data": nil})
		return
	}

	// First try to find the latest active subscription
	var active_subscription models.Subscription
	if err := tx.Where("coach_id = ? AND step = 2", uid).Order("created_at DESC").Preload("SubscriptionPlan").First(&active_subscription).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
	}
	// Check if active subscription is expired and handle status updates
	if active_subscription.Id != 0 {
		now := time.Now()
		if active_subscription.ExpectExpiredAt != nil && now.After(*active_subscription.ExpectExpiredAt) {
			// Update expired subscription status
			if err := tx.Model(&active_subscription).Update("step", 3).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
				return
			}

			// Find and activate the next pending subscription
			var next_subscription models.Subscription
			if err := tx.Where("coach_id = ? AND step = 1", uid).Order("created_at ASC").First(&next_subscription).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
					return
				}
			}
			if next_subscription.Id != 0 {
				// Update next subscription to active
				if err := tx.Model(&next_subscription).Update("step", 2).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
					return
				}
				// Update active_subscription to the newly activated one
				active_subscription = next_subscription
			}
		}
	}
	// Then find the latest subscription regardless of status
	var latest_subscription models.Subscription
	if err := h.db.Where("coach_id = ?", uid).Order("created_at DESC").Preload("SubscriptionPlan").First(&latest_subscription).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
	}
	subscription_resp := gin.H{
		"name":       "",
		"status":     0,
		"expired_at": "",
	}

	// If we have a latest subscription
	if latest_subscription.Id != 0 {
		// If the latest is pending and we have an active one, use the active one
		if latest_subscription.Step != 2 && active_subscription.Id != 0 {
			subscription_resp = gin.H{
				"name":       active_subscription.SubscriptionPlan.Name,
				"status":     active_subscription.Step,
				"expired_at": active_subscription.ExpectExpiredAt,
			}
		} else {
			// Otherwise use the latest one
			subscription_resp = gin.H{
				"name":       latest_subscription.SubscriptionPlan.Name,
				"status":     latest_subscription.Step,
				"expired_at": latest_subscription.ExpectExpiredAt,
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan retrieved successfully", "data": gin.H{
		"id":           coach.Id,
		"nickname":     coach.Profile1.Nickname,
		"avatar_url":   coach.Profile1.AvatarURL,
		"subscription": subscription_resp,
	}})
}

func (h *CoachHandler) UpdateCoachProfile(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Nickname  string `json:"nickname,omitempty" binding:"omitempty,min=3,max=10" label:"昵称"`
		AvatarURL string `json:"avatar_url" label:"头像"`
		Config    string `json:"config"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.AvatarURL == "" && body.Nickname == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少修改参数", "data": nil})
		return
	}

	// Check for sensitive words in nickname
	if body.Nickname != "" && sensitive.ContainsSensitiveWord(body.Nickname) {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "昵称包含敏感词", "data": nil})
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
	var existing models.Coach
	if err := tx.Where("id = ?", uid).Preload("Profile1").First(&existing).Error; err != nil {
		h.logger.Error("Failed to find models.Coach", err)
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
	}
	if existing.Profile1Id == 0 {
		var existing_profile1 models.CoachProfile1
		if err := tx.Where("coach_id = ?", uid).First(&existing_profile1).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
				return
			}
			// 处理历史遗留问题 不存在 CoachProfile1
			existing_profile1.CoachId = existing.Id
			if err := tx.Create(&existing_profile1).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
				return
			}

		}
		// 处理历史遗留问题 CoachProfile1 没有关联上 Coach
		if err := tx.Model(&existing).Updates(map[string]interface{}{
			"profile1_id": existing_profile1.Id,
		}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		existing.Profile1 = existing_profile1
	}
	updates := map[string]interface{}{}
	if body.Nickname != "" {
		updates["nickname"] = body.Nickname
	}
	if body.AvatarURL != "" {
		updates["avatar_url"] = AvatarPrefix + body.AvatarURL
	}
	if err := tx.Model(&existing.Profile1).Updates(updates).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to update the record", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": nil})
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
		AvatarURL: DefaultAvatarURL,
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

func (h *CoachHandler) UpdateStudentProfile(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id        int    `json:"id"`
		Nickname  string `json:"nickname" binding:"min=3,max=10"`
		AvatarURL string `json:"avatar_url"`
		Age       int    `json:"age"`
		Gender    int    `json:"gender"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少id参数", "data": nil})
		return
	}
	if body.Nickname != "" && sensitive.ContainsSensitiveWord(body.Nickname) {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "昵称包含敏感词", "data": nil})
		return
	}
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	var relation models.CoachRelationship
	if err := tx.Where("coach_id = ? AND student_id = ?", uid, body.Id).First(&relation).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "异常操作", "data": nil})
		return
	}
	var existing models.Coach
	if err := tx.Where("id = ?", body.Id).Preload("Profile1").First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "记录不存在", "data": nil})
		return
	}
	if existing.Profile1Id == 0 {
		var existing_profile1 models.CoachProfile1
		if err := tx.Where("coach_id = ?", body.Id).First(&existing_profile1).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
				return
			}
			existing_profile1.CoachId = existing.Id
			// 处理历史遗留问题 不存在 CoachProfile1
			if err := tx.Create(&existing_profile1).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
				return
			}
		}
		// 处理历史遗留问题 CoachProfile1 没有关联上 Coach
		if err := tx.Model(&existing).Updates(map[string]interface{}{
			"profile1_id": existing_profile1.Id,
		}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		existing.Profile1 = existing_profile1
	}
	updates := map[string]interface{}{}
	if body.Nickname != "" {
		updates["nickname"] = body.Nickname
	}
	if body.AvatarURL != "" {
		updates["avatar_url"] = AvatarPrefix + body.AvatarURL
	}
	if body.Age != 0 {
		updates["age"] = body.Age
	}
	if body.Gender != 0 {
		updates["gender"] = body.Gender
	}
	if err := tx.Model(&existing.Profile1).Updates(&updates).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": nil})
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
	query = query.Where("coach_id = ?", uid)

	pb := pagination.NewPaginationBuilder[models.CoachRelationship](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")

	var list []models.CoachRelationship
	if err := pb.Build().Preload("Student").Preload("Student.Profile1").Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list)
	data_list := []interface{}{}
	for _, v := range list2 {
		data_list = append(data_list, map[string]interface{}{
			"id":         v.StudentId,
			"nickname":   v.Student.Profile1.Nickname,
			"avatar_url": v.Student.Profile1.AvatarURL,
			"age":        v.Student.Profile1.Age,
			"gender":     v.Student.Profile1.Gender,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        data_list,
			"page_size":   pb.GetLimit(),
			"has_more":    has_more,
			"next_marker": next_marker,
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

type PersonAvatar struct {
	Id  string
	URL string
}
type PersonWithAvatars struct {
	Id      int
	Avatars []PersonAvatar
}
type AvatarGroup struct {
	Gender  int
	Persons []PersonWithAvatars
}
