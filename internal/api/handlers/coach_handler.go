package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"myapi/internal/models"
	"myapi/internal/pkg/pagination"
	"myapi/internal/pkg/sensitive"
	"myapi/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *CoachHandler) FetchVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": gin.H{
		"version": "2506102221",
	}})
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
		// Generate a random UUID and take first 6 characters
		uuid := uuid.New().String()
		// Remove hyphens and take first 6 characters
		cleanUUID := strings.ReplaceAll(uuid, "-", "")
		return cleanUUID[:6]
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
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
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
		Nickname  string `json:"nickname,omitempty" binding:"omitempty,min=1,max=10" label:"昵称"`
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

func (h *CoachHandler) RefreshCoachStats(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "请先登录", "data": nil})
		return
	}
	var existing models.Coach
	if err := h.db.Where("id = ?", uid).First(&existing).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "异常操作", "data": nil})
		return
	}

	// 获取所有训练天数
	var workout_days []models.WorkoutDay
	if err := h.db.Where("student_id = ?", uid).Find(&workout_days).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	// 计算最近7天的训练情况
	now := time.Now()
	seven_days_ago := now.AddDate(0, 0, -7)
	recent_workout_days := 0

	// 使用map来统计不重复的训练天数
	unique_workout_days := make(map[string]bool)
	recent_unique_workout_days := make(map[string]bool)

	for _, day := range workout_days {
		// 将时间转换为日期字符串（去掉时分秒）
		dateStr := day.CreatedAt.Format("2006-01-02")
		unique_workout_days[dateStr] = true
		if day.CreatedAt.After(seven_days_ago) {
			recent_unique_workout_days[dateStr] = true
		}
	}

	recent_workout_days = len(recent_unique_workout_days)

	type WorkoutStats struct {
		Version           string    `json:"v"`
		TotalWorkoutDays  int       `json:"total_workout_days"`
		TotalWorkoutTimes int       `json:"total_workout_times"`
		RecentWorkoutDays int       `json:"recent_workout_days"`
		CreatedAt         time.Time `json:"created_at"`
	}

	stats := WorkoutStats{
		Version:           "250608",
		TotalWorkoutDays:  len(unique_workout_days),
		TotalWorkoutTimes: len(workout_days),
		RecentWorkoutDays: recent_workout_days,
		CreatedAt:         now,
	}

	statsJSON, err := json.Marshal(stats)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to marshal stats", "data": nil})
		return
	}

	if err := h.db.Model(&existing).Update("workout_stats", string(statsJSON)).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update stats", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "统计数据获取成功",
		"data": stats,
	})
}

func (h *CoachHandler) RefreshWorkoutActionStats(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	// 获取所有训练动作历史
	var workout_action_histories []models.WorkoutActionHistory
	if err := h.db.Where("student_id = ?", uid).Find(&workout_action_histories).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	// 获取所有训练动作
	var workout_actions []models.WorkoutAction
	if err := h.db.Find(&workout_actions).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	// 获取所有训练计划
	var workout_plans []models.WorkoutPlan
	if err := h.db.Where("student_id = ?", uid).Find(&workout_plans).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	recent_workout_actions := 0
	now := time.Now()
	seven_days_ago := now.AddDate(0, 0, -7)

	// 用于排序的切片
	type ActionStat struct {
		ActionId    int
		TotalCount  int
		MaxWeight   int
		MaxReps     int
		TotalWeight int
		TotalReps   int
	}
	var actionStatsList []ActionStat

	// 统计每个动作的完成情况
	action_stats := make(map[int]map[string]interface{})
	for _, action := range workout_actions {
		action_stats[action.Id] = map[string]interface{}{
			"name":           action.Name,
			"total_count":    0,
			"recent_count":   0,
			"last_completed": nil,
			"max_weight":     0,
			"max_reps":       0,
			"total_weight":   0,
			"total_reps":     0,
			"avg_weight":     0.0,
			"avg_reps":       0.0,
		}
	}

	for _, history := range workout_action_histories {
		if stats, exists := action_stats[history.WorkoutActionId]; exists {
			// 更新总次数
			totalCount := stats["total_count"].(int) + 1
			stats["total_count"] = totalCount

			// 更新最近7天次数
			if history.CreatedAt.After(seven_days_ago) {
				recentCount := stats["recent_count"].(int) + 1
				stats["recent_count"] = recentCount
			}

			// 更新最后完成时间
			if stats["last_completed"] == nil || history.CreatedAt.After(stats["last_completed"].(time.Time)) {
				stats["last_completed"] = history.CreatedAt
			}

			// 更新重量相关统计
			if history.Weight > stats["max_weight"].(int) {
				stats["max_weight"] = history.Weight
			}
			totalWeight := stats["total_weight"].(int) + history.Weight
			stats["total_weight"] = totalWeight
			stats["avg_weight"] = float64(totalWeight) / float64(totalCount)

			// 更新次数相关统计
			if history.Reps > stats["max_reps"].(int) {
				stats["max_reps"] = history.Reps
			}
			totalReps := stats["total_reps"].(int) + history.Reps
			stats["total_reps"] = totalReps
			stats["avg_reps"] = float64(totalReps) / float64(totalCount)

			// 添加到排序列表
			actionStatsList = append(actionStatsList, ActionStat{
				ActionId:    history.WorkoutActionId,
				TotalCount:  totalCount,
				MaxWeight:   stats["max_weight"].(int),
				MaxReps:     stats["max_reps"].(int),
				TotalWeight: totalWeight,
				TotalReps:   totalReps,
			})
		}
	}

	for _, history := range workout_action_histories {
		if history.CreatedAt.After(seven_days_ago) {
			recent_workout_actions++
		}
	}
	// 按不同维度排序
	type TopActions struct {
		MostFrequent []map[string]interface{} `json:"most_frequent"` // 最常做的动作
		MaxWeight    []map[string]interface{} `json:"max_weight"`    // 最大重量
		MaxReps      []map[string]interface{} `json:"max_reps"`      // 最大次数
		MostProgress []map[string]interface{} `json:"most_progress"` // 进步最大的动作
	}

	topActions := TopActions{
		MostFrequent: make([]map[string]interface{}, 0, 3),
		MaxWeight:    make([]map[string]interface{}, 0, 3),
		MaxReps:      make([]map[string]interface{}, 0, 3),
		MostProgress: make([]map[string]interface{}, 0, 3),
	}

	// 按总次数排序
	sort.Slice(actionStatsList, func(i, j int) bool {
		return actionStatsList[i].TotalCount > actionStatsList[j].TotalCount
	})
	for i := 0; i < 3 && i < len(actionStatsList); i++ {
		actionId := actionStatsList[i].ActionId
		topActions.MostFrequent = append(topActions.MostFrequent, map[string]interface{}{
			"name":        action_stats[actionId]["name"],
			"total_count": actionStatsList[i].TotalCount,
			"max_weight":  actionStatsList[i].MaxWeight,
			"max_reps":    actionStatsList[i].MaxReps,
			"avg_weight":  action_stats[actionId]["avg_weight"],
			"avg_reps":    action_stats[actionId]["avg_reps"],
		})
	}

	// 按最大重量排序
	sort.Slice(actionStatsList, func(i, j int) bool {
		return actionStatsList[i].MaxWeight > actionStatsList[j].MaxWeight
	})
	for i := 0; i < 3 && i < len(actionStatsList); i++ {
		actionId := actionStatsList[i].ActionId
		topActions.MaxWeight = append(topActions.MaxWeight, map[string]interface{}{
			"name":        action_stats[actionId]["name"],
			"max_weight":  actionStatsList[i].MaxWeight,
			"total_count": actionStatsList[i].TotalCount,
		})
	}

	// 按最大次数排序
	sort.Slice(actionStatsList, func(i, j int) bool {
		return actionStatsList[i].MaxReps > actionStatsList[j].MaxReps
	})
	for i := 0; i < 3 && i < len(actionStatsList); i++ {
		actionId := actionStatsList[i].ActionId
		topActions.MaxReps = append(topActions.MaxReps, map[string]interface{}{
			"name":        action_stats[actionId]["name"],
			"max_reps":    actionStatsList[i].MaxReps,
			"total_count": actionStatsList[i].TotalCount,
		})
	}

}

func (h *CoachHandler) CreateStudent(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Name   string `json:"name" binding:"required,min=1,max=10"`
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
	if body.Name != "" && sensitive.ContainsSensitiveWord(body.Name) {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "名称包含敏感词", "data": nil})
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

func (h *CoachHandler) DeleteStudent(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少id参数", "data": nil})
		return
	}
	query := h.db.Where("d IS NULL OR d = 0")
	query = query.Where("coach_id = ? AND student_id = ?", uid, body.Id)
	var existing models.CoachRelationship
	if err := query.First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有找到记录", "data": nil})
		return
	}
	if err := h.db.Model(&existing).Update("d", 1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功", "data": nil})
}

func (h *CoachHandler) FetchStudentList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		models.Pagination
		Keyword string `json:"keyword"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	query := h.db.Where("d IS NULL OR d = 0")
	query = query.Where("coach_relationship.coach_id = ?", uid)
	if body.Keyword != "" {
		query = query.Joins("JOIN coach_profile1 ON coach_relationship.student_id = coach_profile1.coach_id").
			Where("coach_profile1.nickname LIKE ?", "%"+body.Keyword+"%")
	}
	pb := pagination.NewPaginationBuilder[models.CoachRelationship](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")

	var list1 []models.CoachRelationship
	if err := pb.Build().Preload("Student").Preload("Student.Profile1").Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
	list := make([]map[string]interface{}, 0, len(list2))
	for _, v := range list2 {
		list = append(list, map[string]interface{}{
			"id":         v.StudentId,
			"nickname":   v.Student.Profile1.Nickname,
			"avatar_url": v.Student.Profile1.AvatarURL,
			"age":        v.Student.Profile1.Age,
			"gender":     v.Student.Profile1.Gender,
			"created_at": v.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        list,
			"page_size":   pb.GetLimit(),
			"has_more":    has_more,
			"next_marker": next_marker,
		},
	})
}

func (h *CoachHandler) FetchStudentProfile(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少id参数", "data": nil})
		return
	}
	query := h.db.Where("d IS NULL OR d = 0")
	query = query.Where("id = ?", body.Id)
	var profile models.Coach
	if err := query.Preload("Profile1").First(&profile).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有找到该记录", "data": nil})
		return
	}
	var relation models.CoachRelationship
	if err := h.db.Where("d IS NULL OR d = 0").Where("coach_id = ? AND student_id = ?", uid, body.Id).First(&relation).Error; err != nil {
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
