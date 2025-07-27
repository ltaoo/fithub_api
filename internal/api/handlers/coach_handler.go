package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"myapi/config"
	"myapi/internal/models"
	"myapi/internal/pkg/pagination"
	"myapi/internal/pkg/sensitive"
	"myapi/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CoachHandler struct {
	db     *gorm.DB
	logger *logger.Logger
	config *config.Config
}

// NewCoachHandler creates a new coach handler
func NewCoachHandler(db *gorm.DB, logger *logger.Logger, config *config.Config) *CoachHandler {
	return &CoachHandler{
		db:     db,
		logger: logger,
		config: config,
	}
}

func (h *CoachHandler) FetchVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": gin.H{
		"version": "250707151622",
	}})
}

var AvatarPrefix = "//static.fithub.top/avatars/"
var DefaultAvatarURL = AvatarPrefix + "default1.jpeg"
var MobileSiteHostname = "https://h5.fithub.top"

func (h *CoachHandler) RegisterCoach(c *gin.Context) {
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
	if body.Password == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "请输入密码", "data": nil})
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

	var existing models.CoachAccount
	query := tx.Where("provider_type = ? AND provider_id = ?", models.AccountProviderTypeEmailWithPwd, body.Email)
	if err := query.First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
	}
	if existing.CoachId != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "该邮箱已被使用", "data": nil})
		return
	}

	nickname := func() string {
		// Generate a random UUID and take first 6 characters
		uuid := uuid.New().String()
		// Remove hyphens and take first 6 characters
		uid := strings.ReplaceAll(uuid, "-", "")
		return uid[:6]
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
	token, expires_at, err := models.GenerateJWT(the_coach.Id, h.config.TokenSecretKey)
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
		Token:     "Bearer " + token,
		ExpiresAt: expires_at.Unix(),
		Status:    "success",
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "注册成功", "data": response})
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
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "请输入邮箱", "data": nil})
		return
	}
	if body.Password == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "请输入密码", "data": nil})
		return
	}
	var account models.CoachAccount
	if err := h.db.Where("provider_type = ? AND provider_id = ?", models.AccountProviderTypeEmailWithPwd, body.Email).First(&account).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "帐号不存在", "data": nil})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(account.ProviderArg1), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "邮箱或密码错误", "data": nil})
		return
	}
	// Generate JWT token
	token, expires_at, err := models.GenerateJWT(account.CoachId, h.config.TokenSecretKey)
	if err != nil {
		h.logger.Error("Failed to generate JWT", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "登录失败", "data": nil})
		return
	}
	response := models.AuthResponse{
		Token:     "Bearer " + token,
		ExpiresAt: expires_at.Unix(),
		Status:    "success",
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "登录成功", "data": response})
}

func (h *CoachHandler) FetchCoachProfile(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
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
	var account models.CoachAccount
	if err := tx.Where("coach_id = ?", coach.Id).First(&account).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
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
		"count":      0,
	}

	// If we have a latest subscription
	if latest_subscription.Id != 0 {
		// If the latest is pending and we have an active one, use the active one
		if latest_subscription.Step != 2 && active_subscription.Id != 0 {
			subscription_resp = gin.H{
				"name":       active_subscription.SubscriptionPlan.Name,
				"status":     active_subscription.Step,
				"expired_at": active_subscription.ExpectExpiredAt,
				"count":      active_subscription.Count,
			}
		} else {
			// Otherwise use the latest one
			subscription_resp = gin.H{
				"name":       latest_subscription.SubscriptionPlan.Name,
				"status":     latest_subscription.Step,
				"expired_at": latest_subscription.ExpectExpiredAt,
				"count":      latest_subscription.Count,
			}
		}
		// 生效中
		if subscription_resp["status"] == 2 && subscription_resp["count"] == 9999 {
			subscription_resp["name"] = "终身VIP"
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": gin.H{
		"id":           coach.Id,
		"uid":          coach.Nickname,
		"nickname":     coach.Profile1.Nickname,
		"avatar_url":   coach.Profile1.AvatarURL,
		"subscription": subscription_resp,
		"no_account":   account.ProviderType == 0,
	}})
}

func (h *CoachHandler) UpdateCoachProfile(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Nickname  string `json:"nickname,omitempty" binding:"omitempty,min=1,max=18" label:"昵称"`
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

type LocalTime struct {
	time.Time
}

func (t *LocalTime) UnmarshalJSON(b []byte) error {
	// 去掉引号
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
	// 先按 RFC3339 解析
	tt, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	// 转为本地时区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t.Time = tt.In(loc)
	return nil
}

func (h *CoachHandler) RefreshCoachStats(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "请先登录", "data": nil})
		return
	}
	var body struct {
		RangeOfStart *LocalTime `json:"range_of_start"`
		RangeOfEnd   *LocalTime `json:"range_of_end"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.RangeOfStart == nil || body.RangeOfEnd == nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "缺少参数", "data": nil})
		return
	}
	var existing models.Coach
	if err := h.db.Where("id = ?", uid).First(&existing).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "异常操作", "data": nil})
		return
	}

	now := time.Now()

	// 1.1 获取所有训练日期（用于最长连续天数）
	var dateList []string
	h.db.Raw(`SELECT DISTINCT finished_at as date_str FROM WORKOUT_DAY WHERE student_id = ? AND status = 2 AND DATE(finished_at, '+8 hours') BETWEEN ? AND ? ORDER BY date_str ASC`, uid, body.RangeOfStart.Time, body.RangeOfEnd.Time).Scan(&dateList)

	// fmt.Println("dateList:", dateList)
	if len(dateList) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "本月没有任何训练记录", "data": nil})
		return
	}

	// 统计最长连续天数
	max_streak, cur_streak := 0, 0
	var last_date time.Time
	var cur_streak_start, max_streak_start, max_streak_end time.Time
	for i, dateStr := range dateList {
		t, _ := time.Parse("2006-01-02", dateStr)
		if i == 0 || t.Sub(last_date) == 24*time.Hour {
			cur_streak++
			if cur_streak == 1 {
				cur_streak_start = t
			}
		} else {
			cur_streak = 1
			cur_streak_start = t
		}
		if cur_streak > max_streak {
			max_streak = cur_streak
			max_streak_start = cur_streak_start
			max_streak_end = t
		}
		// fmt.Printf("[streak] i=%d, date=%s, cur_streak=%d, cur_streak_start=%s, max_streak=%d, max_streak_start=%s, max_streak_end=%s\n", i, dateStr, cur_streak, cur_streak_start.Format("2006-01-02"), max_streak, max_streak_start.Format("2006-01-02"), max_streak_end.Format("2006-01-02"))
		last_date = t
	}
	// fmt.Printf("[streak result] max_streak=%d, start=%s, end=%s\n", max_streak, max_streak_start.Format("2006-01-02"), max_streak_end.Format("2006-01-02"))

	// 2. 训练时长最长、容量最大的、最早开始、最晚完成的 WorkoutDay
	var max_duration_day, max_volume_day, earliest_start_day, latest_finish_day models.WorkoutDay
	h.db.Where("student_id = ? AND status = 2 AND DATE(finished_at, '+8 hours') BETWEEN ? AND ?", uid, body.RangeOfStart.Time, body.RangeOfEnd.Time).
		Order("duration DESC").Preload("WorkoutPlan").First(&max_duration_day)
	h.db.Where("student_id = ? AND status = 2 AND DATE(finished_at, '+8 hours') BETWEEN ? AND ?", uid, body.RangeOfStart.Time, body.RangeOfEnd.Time).
		Order("total_volume DESC").Preload("WorkoutPlan").First(&max_volume_day)
	h.db.Where("student_id = ? AND status = 2 AND DATE(finished_at, '+8 hours') BETWEEN ? AND ?", uid, body.RangeOfStart.Time, body.RangeOfEnd.Time).
		Order("strftime('%H', started_at), started_at ASC").Preload("WorkoutPlan").First(&earliest_start_day)
	h.db.Where("student_id = ? AND status = 2 AND DATE(finished_at, '+8 hours') BETWEEN ? AND ?", uid, body.RangeOfStart.Time, body.RangeOfEnd.Time).
		Order("strftime('%H', finished_at) DESC, finished_at DESC").Preload("WorkoutPlan").First(&latest_finish_day)

	// 3. 聚合 WorkoutPlan.type 下的所有 workout_day 及其 workout_plan
	type WorkoutDayWithPlan struct {
		WorkoutDayId int
		PlanId       int
		PlanTitle    string
		PlanType     string
		// 你可以根据需要加更多字段
	}
	var workout_days_with_plan []WorkoutDayWithPlan
	h.db.Raw(`
		SELECT 
			wd.id as workout_day_id, 
			wp.id as plan_id, 
			wp.title as plan_title, 
			wp.type as plan_type
		FROM WORKOUT_DAY wd 
		JOIN WORKOUT_PLAN wp ON wd.workout_plan_id = wp.id 
		WHERE wd.student_id = ? AND wd.status = 2 AND DATE(wd.finished_at, '+8 hours') BETWEEN ? AND ?
		ORDER BY wp.type, wd.id
	`, uid, body.RangeOfStart.Time, body.RangeOfEnd.Time).Scan(&workout_days_with_plan)

	// 按 type 分组
	workout_day_group_with_type := map[string][]map[string]interface{}{}
	for _, v := range workout_days_with_plan {
		workout_day_group_with_type[v.PlanType] = append(workout_day_group_with_type[v.PlanType], map[string]interface{}{
			"workout_day_id": v.WorkoutDayId,
			"workout_plan": map[string]interface{}{
				"id":    v.PlanId,
				"title": v.PlanTitle,
				// 你可以加更多字段
			},
		})
	}

	// 5. 对比同一训练计划最早和最晚一次训练的容量差距
	// type PlanStat struct {
	// 	PlanId      int     `json:"plan_id"`
	// 	PlanName    string  `json:"plan_name"`
	// 	FirstVolume float64 `json:"first_volume"`
	// 	LastVolume  float64 `json:"last_volume"`
	// 	Progress    float64 `json:"progress"`
	// }
	// var plan_stats []PlanStat
	// h.db.Raw(`
	// SELECT
	//   wd1.workout_plan_id as plan_id,
	//   wp.title as plan_name,
	//   wd1.total_volume as first_volume,
	//   wd2.total_volume as last_volume,
	//   wd2.total_volume - wd1.total_volume as progress
	// FROM (
	//   SELECT workout_plan_id, MIN(created_at) as first_date, MAX(created_at) as last_date
	//   FROM WORKOUT_DAY
	//   WHERE student_id = ? AND status = 2 AND finished_at BETWEEN ? AND ?
	//   GROUP BY workout_plan_id
	//   HAVING COUNT(*) >= 2
	// ) t
	// JOIN WORKOUT_DAY wd1 ON wd1.workout_plan_id = t.workout_plan_id AND wd1.created_at = t.first_date
	// JOIN WORKOUT_DAY wd2 ON wd2.workout_plan_id = t.workout_plan_id AND wd2.created_at = t.last_date
	// JOIN WORKOUT_PLAN wp ON wp.id = t.workout_plan_id
	// `, uid, body.RangeOfStart.Time, body.RangeOfEnd.Time).Scan(&plan_stats)

	// 总训练次数
	var totalWorkoutTimes int64
	h.db.Model(&models.WorkoutDay{}).
		Where("student_id = ? AND status = 2 AND DATE(finished_at, '+8 hours') BETWEEN ? AND ?", uid, body.RangeOfStart.Time, body.RangeOfEnd.Time).
		Count(&totalWorkoutTimes)

	type WorkoutStats struct {
		Version           string    `json:"v"`
		TotalWorkoutDays  int       `json:"total_workout_days"`
		TotalWorkoutTimes int64     `json:"total_workout_times"`
		CreatedAt         time.Time `json:"created_at"`
	}

	// 1. 统计不重复的训练天数
	var total_workout_days int
	h.db.Raw(`SELECT COUNT(DISTINCT DATE(finished_at, '+8 hours')) FROM WORKOUT_DAY WHERE student_id = ? AND status = 2 AND DATE(finished_at, '+8 hours') BETWEEN ? AND ?`, uid, body.RangeOfStart.Time, body.RangeOfEnd.Time).Scan(&total_workout_days)

	stats := WorkoutStats{
		Version:           "250608",
		TotalWorkoutDays:  total_workout_days,
		TotalWorkoutTimes: totalWorkoutTimes,
		CreatedAt:         now,
	}

	// 动作统计
	var action_histories []models.WorkoutActionHistory
	if err := h.db.
		Where("student_id = ? AND DATE(created_at, '+8 hours') BETWEEN ? AND ?", uid, body.RangeOfStart.Time, body.RangeOfEnd.Time).
		Preload("WorkoutAction").
		Find(&action_histories).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	// 1. 按动作分组
	type Record struct {
		Reps       int       `json:"reps"`
		RepsUnit   string    `json:"reps_unit"`
		Weight     float64   `json:"weight"`
		WeightUnit string    `json:"weight_unit"`
		CreatedAt  time.Time `json:"created_at"`
	}
	action_map := make(map[string][]Record)
	for _, history := range action_histories {
		action_name := history.WorkoutAction.ZhName // 或英文名
		rec := Record{
			Reps:       history.Reps,
			RepsUnit:   history.RepsUnit,
			Weight:     history.Weight,
			WeightUnit: history.WeightUnit,
			CreatedAt:  history.CreatedAt,
		}
		action_map[action_name] = append(action_map[action_name], rec)
	}

	// 2. 转为目标结构
	type ActionStat struct {
		Action  string   `json:"action"`
		Records []Record `json:"records"`
	}
	var action_stats []ActionStat
	for action, records := range action_map {
		action_stats = append(action_stats, ActionStat{
			Action:  action,
			Records: records,
		})
	}

	// 3. 可选：按动作名排序
	// sort.Slice(result, func(i, j int) bool { return result[i].Action < result[j].Action })

	// 4. 返回
	// c.JSON(http.StatusOK, gin.H{
	// 	"code": 200,
	// 	"msg":  "动作统计获取成功",
	// 	"data": result,
	// })

	stats_json, err := json.Marshal(stats)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to marshal stats", "data": nil})
		return
	}

	if err := h.db.Model(&existing).Update("workout_stats", string(stats_json)).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update stats", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"stats":      stats,
			"max_streak": max_streak,
			"max_streak_range": gin.H{
				"start": func() string {
					if max_streak > 0 {
						return max_streak_start.Format("2006-01-02")
					} else {
						return ""
					}
				}(),
				"end": func() string {
					if max_streak > 0 {
						return max_streak_end.Format("2006-01-02")
					} else {
						return ""
					}
				}(),
			},
			"max_duration_day": gin.H{
				"id":           max_duration_day.Id,
				"started_at":   max_duration_day.StartedAt,
				"finished_at":  max_duration_day.FinishedAt,
				"duration":     max_duration_day.Duration,
				"total_volume": max_duration_day.TotalVolume,
				"workout_plan": gin.H{
					"id":       max_duration_day.WorkoutPlan.Id,
					"title":    max_duration_day.WorkoutPlan.Title,
					"overview": max_duration_day.WorkoutPlan.Overview,
					"type":     max_duration_day.WorkoutPlan.Type,
				},
			},
			"earliest_start_day": gin.H{
				"id":           earliest_start_day.Id,
				"started_at":   earliest_start_day.StartedAt,
				"finished_at":  earliest_start_day.FinishedAt,
				"duration":     earliest_start_day.Duration,
				"total_volume": earliest_start_day.TotalVolume,
				"workout_plan": gin.H{
					"id":       earliest_start_day.WorkoutPlan.Id,
					"title":    earliest_start_day.WorkoutPlan.Title,
					"overview": earliest_start_day.WorkoutPlan.Overview,
					"type":     earliest_start_day.WorkoutPlan.Type,
				},
			},
			"latest_finish_day": gin.H{
				"id":           latest_finish_day.Id,
				"started_at":   latest_finish_day.StartedAt,
				"finished_at":  latest_finish_day.FinishedAt,
				"duration":     latest_finish_day.Duration,
				"total_volume": latest_finish_day.TotalVolume,
				"workout_plan": gin.H{
					"id":       latest_finish_day.WorkoutPlan.Id,
					"title":    latest_finish_day.WorkoutPlan.Title,
					"overview": latest_finish_day.WorkoutPlan.Overview,
					"type":     latest_finish_day.WorkoutPlan.Type,
				},
			},
			"type_plan_map": workout_day_group_with_type,
			"max_volume_day": gin.H{
				"id":           max_volume_day.Id,
				"started_at":   max_volume_day.StartedAt,
				"finished_at":  max_volume_day.FinishedAt,
				"duration":     max_volume_day.Duration,
				"total_volume": max_volume_day.TotalVolume,
				"workout_plan": gin.H{
					"id":       max_volume_day.WorkoutPlan.Id,
					"title":    max_volume_day.WorkoutPlan.Title,
					"overview": max_volume_day.WorkoutPlan.Overview,
					"type":     max_volume_day.WorkoutPlan.Type,
				},
			},
			// "plan_stats": plan_stats,
			"action_stats": action_stats,
		},
	})
}

// 刷新当日的训练总结
func (h *CoachHandler) RefreshTodayWorkoutStats(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "请先登录", "data": nil})
		return
	}
	var body struct {
		RangeOfStart *time.Time `json:"range_of_start"`
		RangeOfEnd   *time.Time `json:"range_of_end"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.RangeOfStart == nil || body.RangeOfEnd == nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "缺少参数", "data": nil})
		return
	}
	var existing_coach models.Coach
	if err := h.db.Where("id = ?", uid).First(&existing_coach).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "异常操作", "data": nil})
		return
	}

	var existing_workout_days []models.WorkoutDay
	if err := h.db.Where("student_id = ? AND status = 2 AND finished_at BETWEEN ? AND ?", uid, body.RangeOfStart, body.RangeOfEnd).Find(&existing_workout_days).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	error_msg := make([]string, 0)
	list := make([]models.TodayWorkoutActionGroup, 0)
	tags := make([]string, 0)
	total_volume := float64(0)
	set_count := 0
	duration_count := 0
	for _, v := range existing_workout_days {
		if v.Status != 2 {
			continue
		}
		result, err := models.BuildResultFromWorkoutDay(v, h.db)
		if err != nil {
			error_msg = append(error_msg, err.Error())
		}
		set_count += result.SetCount
		total_volume += result.TotalVolume
		duration_count += result.DurationCount
		list = append(list, result.List...)
		tags = append(tags, result.Tags...)
	}
	tags = lo.Uniq(tags)
	tags = lo.Filter(tags, func(x string, idx int) bool {
		return x != ""
	})
	if len(error_msg) != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  strings.Join(error_msg, "\n"),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"workout_steps":  list,
			"times":          len(existing_workout_days),
			"set_count":      set_count,
			"duration_count": duration_count,
			"volume_count":   total_volume,
			"tags":           tags,
		},
	})
}

func (h *CoachHandler) RefreshWorkoutActionStats(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		RangeOfStart *time.Time `json:"range_of_start"`
		RangeOfEnd   *time.Time `json:"range_of_end"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.RangeOfStart == nil || body.RangeOfEnd == nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "缺少参数", "data": nil})
		return
	}

	query := h.db.Where("student_id = ? AND created_at BETWEEN ? AND ?", uid, body.RangeOfStart, body.RangeOfEnd)
	var workout_action_histories []models.WorkoutActionHistory
	if err := query.Preload("WorkoutAction").Find(&workout_action_histories).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	// 根据 WorkoutAction 的 tags1 统计每个部位的组数，tags1 是 胸,肩 这样的字符串，请将其按照 , 分割后统计
	body_part_stats := make(map[string]int) // 用于统计每个部位的组数

	for _, history := range workout_action_histories {
		if history.WorkoutAction.Tags1 != "" {
			// 将 tags1 按逗号分割
			tags := strings.Split(history.WorkoutAction.Tags1, ",")
			for _, tag := range tags {
				// 去除空格
				tag = strings.TrimSpace(tag)
				if tag != "" {
					// 统计每个部位的组数
					body_part_stats[tag]++
				}
			}
		}
	}

	// 将统计结果转换为排序后的切片
	type BodyPartStat struct {
		BodyPart string `json:"body_part"`
		Sets     int    `json:"sets"`
	}

	var body_part_stats_list []BodyPartStat
	for body_part, sets := range body_part_stats {
		body_part_stats_list = append(body_part_stats_list, BodyPartStat{
			BodyPart: body_part,
			Sets:     sets,
		})
	}

	// 按组数降序排序
	sort.Slice(body_part_stats_list, func(i, j int) bool {
		return body_part_stats_list[i].Sets > body_part_stats_list[j].Sets
	})

	// 计算总体统计
	total_sets := 0
	for _, stat := range body_part_stats_list {
		total_sets += stat.Sets
	}

	// 构建响应数据
	response := gin.H{
		"body_part_stats": body_part_stats_list,
		"total_sets":      total_sets,
		"total_actions":   len(workout_action_histories),
		"range_start":     body.RangeOfStart,
		"range_end":       body.RangeOfEnd,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "统计数据获取成功",
		"data": response,
	})
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
	nickname := func() string {
		// Generate a random UUID and take first 6 characters
		uuid := uuid.New().String()
		// Remove hyphens and take first 6 characters
		uid := strings.ReplaceAll(uuid, "-", "")
		return uid[:6]
	}()
	student := models.Coach{
		Nickname:  nickname,
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
		StudentId: student.Id,
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

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": gin.H{"id": student.Id}})
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

func (h *CoachHandler) BuildStudentAuthURL(c *gin.Context) {
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
	token, _, err := models.GenerateJWT(existing.StudentId, h.config.TokenSecretKey)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	data := gin.H{
		"url": "/home/index?token=" + token,
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": data})
}

func (h *CoachHandler) BuildCoachAuthURLInAdmin(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid != 1 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "没有权限", "data": nil})
		return
	}
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
	var existing models.Coach
	if err := h.db.Where("id = ?", body.Id).First(&existing).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	token, _, err := models.GenerateJWT(existing.Id, h.config.TokenSecretKey)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	data := gin.H{
		"url": MobileSiteHostname + "/home/index?token=" + token,
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": data})
}

func (h *CoachHandler) CreateCoach(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid != 1 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "没有权限", "data": nil})
	}
	var body struct {
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
		Bio       string `json:"bio"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	if body.Nickname == "" {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "缺少昵称", "data": nil})
		return
	}
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()

	now := time.Now()
	the_coach := models.Coach{
		Nickname:  body.Nickname,
		Config:    "{}",
		Bio:       body.Bio,
		CreatedAt: now,
	}

	if err := tx.Create(&the_coach).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create coach", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	profile1 := models.CoachProfile1{
		Nickname:  body.Nickname,
		AvatarURL: body.AvatarURL,
		CoachId:   the_coach.Id,
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
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": nil})
}

func (h *CoachHandler) FetchCoachList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid != 1 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "没有权限", "data": nil})
	}
	var body struct {
		models.Pagination
		Keyword string `json:"keyword"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	query := h.db.Where("d IS NULL OR d = 0")
	if body.Keyword != "" {
		query = query.Where("coach.profile1.nickname LIKE ?", "%"+body.Keyword+"%")
	}
	pb := pagination.NewPaginationBuilder[models.Coach](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")

	var list1 []models.Coach
	if err := pb.Build().Preload("Profile1").Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
	list := make([]map[string]interface{}, 0, len(list2))
	for _, v := range list2 {
		list = append(list, map[string]interface{}{
			"id":         v.Id,
			"nickname":   v.Profile1.Nickname,
			"avatar_url": v.Profile1.AvatarURL,
			"age":        v.Profile1.Age,
			"gender":     v.Profile1.Gender,
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
	query = query.Where("coach_relationship.coach_id = ? OR coach_relationship.student_id = ?", uid, uid)
	if body.Keyword != "" {
		query = query.Joins("JOIN coach_profile1 ON coach_relationship.student_id = coach_profile1.coach_id").
			Where("coach_profile1.nickname LIKE ?", "%"+body.Keyword+"%")
	}
	pb := pagination.NewPaginationBuilder[models.CoachRelationship](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")

	var list1 []models.CoachRelationship
	if err := pb.Build().Preload("Coach").Preload("Coach.Profile1").Preload("Student").Preload("Student.Profile1").Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
	list := make([]map[string]interface{}, 0, len(list2))
	for _, v := range list2 {
		if v.Role == models.RoleCoachStudent {
			if uid == v.StudentId {
				data := map[string]interface{}{
					"id":         v.CoachId,
					"nickname":   v.Coach.Profile1.Nickname,
					"avatar_url": v.Coach.Profile1.AvatarURL,
					"age":        v.Coach.Profile1.Age,
					"gender":     v.Coach.Profile1.Gender,
					"role":       v.Role,
					"role_text":  "教练",
					"status":     v.Status,
					"created_at": v.CreatedAt,
				}
				list = append(list, data)
				if v.Role == int(models.RoleCoachStudent) {
					data["role"] = int(models.RoleStudentCoach)
				}
				if v.Role == int(models.RoleCoachAndStudentHasAccount) {
					data["role"] = int(models.RoleStudentHasAccountAndCoach)
				}
			}
			if uid == v.CoachId {
				list = append(list, map[string]interface{}{
					"id":         v.StudentId,
					"nickname":   v.Student.Profile1.Nickname,
					"avatar_url": v.Student.Profile1.AvatarURL,
					"age":        v.Student.Profile1.Age,
					"gender":     v.Student.Profile1.Gender,
					"role":       v.Role,
					"role_text":  "学员",
					"status":     v.Status,
					"created_at": v.CreatedAt,
				})
			}
		}
		if v.Role == models.RoleFriendAndFriend {
			if uid == v.StudentId {
				list = append(list, map[string]interface{}{
					"id":         v.CoachId,
					"nickname":   v.Coach.Profile1.Nickname,
					"avatar_url": v.Coach.Profile1.AvatarURL,
					"age":        v.Coach.Profile1.Age,
					"gender":     v.Coach.Profile1.Gender,
					"role":       v.Role,
					"role_text":  "好友",
					"status":     v.Status,
					"created_at": v.CreatedAt,
				})
			}
			if uid == v.CoachId {
				list = append(list, map[string]interface{}{
					"id":         v.StudentId,
					"nickname":   v.Student.Profile1.Nickname,
					"avatar_url": v.Student.Profile1.AvatarURL,
					"age":        v.Student.Profile1.Age,
					"gender":     v.Student.Profile1.Gender,
					"role":       v.Role,
					"role_text":  "好友",
					"status":     v.Status,
					"created_at": v.CreatedAt,
				})
			}

		}

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
	if err := h.db.Where("d IS NULL OR d = 0").Where("coach_id = ? AND student_id = ? OR coach_id = ? AND student_id = ?", uid, body.Id, body.Id, uid).First(&relation).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 501, "msg": "没有找到该记录", "data": nil})
		return
	}
	data := gin.H{
		"id":         profile.Id,
		"nickname":   profile.Profile1.Nickname,
		"avatar_url": profile.Profile1.AvatarURL,
		"gender":     profile.Profile1.Gender,
		"age":        profile.Profile1.Age,
		"status":     relation.Status,
		"role":       relation.Role,
	}
	if relation.StudentId == uid {
		if relation.Role == int(models.RoleCoachStudent) {
			data["role"] = int(models.RoleStudentCoach)
		}
		if relation.Role == int(models.RoleCoachAndStudentHasAccount) {
			data["role"] = int(models.RoleStudentHasAccountAndCoach)
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "请求成功", "data": data})
}

func (h *CoachHandler) RefreshToken(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	// Generate new JWT token
	token, expires_at, err := models.GenerateJWT(uid, h.config.TokenSecretKey)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to generate token", "data": nil})
		return
	}

	// Return response
	response := models.AuthResponse{
		Token:     "Bearer " + token,
		ExpiresAt: expires_at.Unix(),
		Status:    "success",
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "请求成功", "data": response})
}

// 通过授权链接访问，并且补全登录信息
func (h *CoachHandler) CreateAccount(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

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
	if body.Password == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "请输入密码", "data": nil})
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

	var coach models.Coach
	if err := tx.Where("id = ?", uid).First(&coach).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "异常操作", "data": nil})
		return
	}
	var relation models.CoachRelationship
	// 只有学员才需要补全帐号
	if err := tx.Where("student_id = ?", uid).First(&relation).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "异常操作", "data": nil})
		return
	}
	var existing models.CoachAccount
	query := tx.Where("provider_type = ? AND provider_id = ?", models.AccountProviderTypeEmailWithPwd, body.Email)
	if err := query.First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "已经补全过帐号了", "data": nil})
		return
	}
	if existing.CoachId != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "该邮箱已被使用", "data": nil})
		return
	}

	now := time.Now()

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
		CoachId:      coach.Id,
	}
	if err := tx.Create(&the_coach_account).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create coach account", err)
		c.JSON(http.StatusOK, gin.H{"code": 502, "msg": "操作失败", "data": nil})
		return
	}
	// 补全后，更新教练和学员的关系，让教练不能再管理学员了
	if err := tx.Model(&relation).Update("role", int(models.RoleCoachAndStudentHasAccount)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 502, "msg": "操作失败", "data": nil})
		return
	}

	token, expires_at, err := models.GenerateJWT(coach.Id, h.config.TokenSecretKey)
	if err != nil {
		tx.Rollback()
		h.logger.Error("Failed to generate JWT", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "操作失败", "data": nil})
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	response := models.AuthResponse{
		Token:     "Bearer " + token,
		ExpiresAt: expires_at.Unix(),
		Status:    "success",
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "操作成功", "data": response})
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

func (h *CoachHandler) FetchArticleList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		models.Pagination
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	query := h.db.Where("d IS NULL OR d = 0")
	if uid != 0 {
		query = query.Where("(publish = 1) OR (publish = 2 AND coach_id = ?)", uid)
	} else {
		query = query.Where("publish = 1")
	}
	pb := pagination.NewPaginationBuilder[models.CoachContent](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")
	var list1 []models.CoachContent
	if err := pb.Build().Preload("Coach.Profile1").Find(&list1).Error; err != nil {
		h.logger.Error("Failed to fetch coach content list", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
	fmt.Println("the", len(list2))
	list := make([]map[string]interface{}, 0, len(list2))
	for _, v := range list2 {
		list = append(list, gin.H{
			"id":        v.Id,
			"title":     v.Title,
			"overview":  v.Description,
			"type":      v.ContentType,
			"video_url": v.VideoKey,
			"creator": gin.H{
				"nickname":   v.Coach.Profile1.Nickname,
				"avatar_url": v.Coach.Profile1.AvatarURL,
			},
			"created_at": v.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": gin.H{
			"list":        list,
			"page_size":   pb.GetLimit(),
			"has_more":    has_more,
			"next_marker": next_marker,
		},
	})
}

func (h *CoachHandler) FetchArticleProfile(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少参数", "data": nil})
		return
	}
	query := h.db.Where("d IS NULL OR d == 0")
	if uid != 0 {
		query = query.Where("(publish = 1) OR (publish = 2 AND coach_id = ?)", uid)
	} else {
		query = query.Where("publish = 1")
	}
	query = query.Where("id = ?", body.Id)
	var existing models.CoachContent
	if err := query.Preload("Coach.Profile1").First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有找到记录", "data": nil})
		return
	}
	var the_points1 []models.CoachContentWithWorkoutAction
	query2 := h.db.Where("d IS NULL OR d = 0")
	query2 = query2.Where("coach_content_id = ?", existing.Id)
	if err := query2.Preload("WorkoutAction").Find(&the_points1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	points2 := make([]map[string]interface{}, 0, len(the_points1))
	for _, v := range the_points1 {
		var act map[string]interface{}
		if v.WorkoutActionId != 0 {
			act = make(map[string]interface{})
			act["id"] = v.WorkoutAction.Id
			act["zh_name"] = v.WorkoutAction.ZhName
			act["score"] = v.WorkoutAction.Score
		}
		points2 = append(points2, gin.H{
			"id":             v.Id,
			"text":           v.Text,
			"time":           v.StartPoint,
			"workout_action": act,
		})
	}
	data := gin.H{
		"id":          existing.Id,
		"title":       existing.Title,
		"overview":    existing.Description,
		"type":        existing.ContentType,
		"video_url":   existing.VideoKey,
		"created_at":  existing.CreatedAt,
		"time_points": points2,
		"is_author":   existing.CoachId == uid,
		"creator": gin.H{
			"nickname":   existing.Coach.Profile1.Nickname,
			"avatar_url": existing.Coach.Profile1.AvatarURL,
		},
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": data})
}

func (h *CoachHandler) CreateArticle(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Title    string `json:"title"`
		Overview string `json:"overview"`
		Type     int    `json:"type"`
		Status   int    `json:"status"`
		// CoverImageURL string `json:"cover_image_url"`
		// ContentURL    string `json:"content_url"`
		VideoURL   string `json:"video_url"`
		TimePoints []struct {
			Time            int    `json:"time"`
			Text            string `json:"text"`
			WorkoutActionId int    `json:"workout_action_id"`
		} `json:"time_points"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	if body.Title == "" {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "缺少标题", "data": nil})
		return
	}
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	// var existing models.CoachContent
	// if err := tx.Where("video_key = ?", body.VideoURL).First(&existing).Error; err == nil {
	// 	c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "该记录已存在", "data": nil})
	// 	return
	// }
	now := time.Now()
	the_content := models.CoachContent{
		Title:         body.Title,
		ContentType:   body.Type,
		Description:   body.Overview,
		ContentURL:    "",
		CoverImageURL: "",
		VideoKey:      body.VideoURL,
		Status:        1,
		Publish:       body.Status,
		LikeCount:     0,
		CreatedAt:     now,
		CoachId:       uid,
	}
	if err := tx.Create(&the_content).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create content", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	for _, point := range body.TimePoints {
		the_content_with_action := models.CoachContentWithWorkoutAction{
			WorkoutActionId: point.WorkoutActionId,
			CoachContentId:  the_content.Id,
			StartPoint:      point.Time,
			Status:          body.Status,
			Text:            point.Text,
			CreatedAt:       now,
		}
		if err := tx.Create(&the_content_with_action).Error; err != nil {
			tx.Rollback()
			h.logger.Error("Failed to create point", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}

	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": gin.H{
		"id": the_content.Id,
	}})
}

func (h *CoachHandler) UpdateArticle(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id         int    `json:"id"`
		Title      string `json:"title"`
		Overview   string `json:"overview"`
		Type       int    `json:"type"`
		VideoURL   string `json:"video_url"`
		TimePoints []struct {
			Id              int    `json:"id"`
			Time            int    `json:"time"`
			Text            string `json:"text"`
			WorkoutActionId int    `json:"workout_action_id"`
		} `json:"time_points"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "缺少 id", "data": nil})
		return
	}
	if body.Title == "" {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "缺少标题", "data": nil})
		return
	}
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	var existing models.CoachContent
	if err := tx.Where("id = ? AND coach_id = ?", body.Id, uid).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有找到记录", "data": nil})
		return
	}
	now := time.Now()
	updates := map[string]interface{}{}
	if body.Title != "" {
		updates["title"] = body.Title
	}
	if body.Overview != "" {
		updates["description"] = body.Overview
	}
	if body.Type != 0 {
		updates["content_type"] = body.Type
	}
	if body.VideoURL != "" {
		updates["video_key"] = body.VideoURL
	}
	if err := tx.Model(&existing).Updates(&updates).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create content", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	// 1. 查询当前 content 的所有 time_points
	var existing_points []models.CoachContentWithWorkoutAction
	if err := tx.Where("d IS NULL OR d = 0").Where("coach_content_id = ?", existing.Id).Find(&existing_points).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	existing_point_map := make(map[int]models.CoachContentWithWorkoutAction)
	for _, p := range existing_points {
		existing_point_map[p.Id] = p
	}

	// 2. 遍历 body.TimePoints，做新增/更新，并记录前端传来的 id
	new_point_id_set := make(map[int]bool)
	for _, point := range body.TimePoints {
		if point.Id > 0 {
			// update
			if _, ok := existing_point_map[point.Id]; ok {
				update_fields := map[string]interface{}{
					"start_point":       point.Time,
					"text":              point.Text,
					"workout_action_id": point.WorkoutActionId,
				}
				if err := tx.Model(&models.CoachContentWithWorkoutAction{}).Where("id = ?", point.Id).Updates(update_fields).Error; err != nil {
					tx.Rollback()
					h.logger.Error("Failed to update point", err)
					c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
					return
				}
				new_point_id_set[point.Id] = true
			}
		} else {
			// create
			newPoint := models.CoachContentWithWorkoutAction{
				CoachContentId:  existing.Id,
				StartPoint:      point.Time,
				Text:            point.Text,
				WorkoutActionId: point.WorkoutActionId,
				CreatedAt:       now,
				Status:          1, // 默认公开
			}
			if err := tx.Create(&newPoint).Error; err != nil {
				tx.Rollback()
				h.logger.Error("Failed to create point", err)
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
				return
			}
			new_point_id_set[newPoint.Id] = true
		}
	}

	// 3. 删除未出现在前端的点（软删除）
	for id := range existing_point_map {
		if !new_point_id_set[id] {
			if err := tx.Model(&models.CoachContentWithWorkoutAction{}).Where("id = ?", id).Update("d", 1).Error; err != nil {
				tx.Rollback()
				h.logger.Error("Failed to delete point", err)
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": gin.H{
		"id": existing.Id,
	}})
}

func (h *CoachHandler) CreateCoachContent(c *gin.Context) {
	// uid := int(c.GetFloat64("id"))
	var body struct {
		Title         string `json:"title"`
		Description   string `json:"description"`
		ContentType   int    `json:"content_type"`
		CoverImageURL string `json:"cover_image_url"`
		ContentURL    string `json:"content_url"`
		VideoKey      string `json:"video_key"`
		CoachId       int    `json:"coach_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	if body.Title == "" {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "缺少标题", "data": nil})
		return
	}
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	var existing models.CoachContent
	if err := tx.Where("content_url = ?", body.ContentURL).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "该记录已存在", "data": nil})
		return
	}
	now := time.Now()
	the_content := models.CoachContent{
		Title:         body.Title,
		ContentType:   body.ContentType,
		Description:   body.Description,
		ContentURL:    body.ContentURL,
		CoverImageURL: body.ContentURL,
		VideoKey:      body.VideoKey,
		Status:        1,
		Publish:       1,
		LikeCount:     0,
		CreatedAt:     now,
		CoachId:       body.CoachId,
	}
	if err := tx.Create(&the_content).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create coach", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": nil})
}

// 添加好友
func (h *CoachHandler) AddFriend(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		// 添加好友只能通过 uid
		UID string `json:"uid"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if body.UID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "缺少参数", "data": nil})
		return
	}
	var friend models.Coach
	if err := h.db.Where("nickname = ?", body.UID).First(&friend).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "UID 错误", "data": nil})
		return
	}
	var existing models.CoachRelationship
	if err := h.db.Where("coach_id = ? AND student_id = ? OR coach_id = ? AND student_id = ?", uid, friend.Id, friend.Id, uid).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		created := models.CoachRelationship{
			Role:      3,
			CoachId:   uid,
			StudentId: friend.Id,
		}
		if err := h.db.Create(&created).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "添加成功", "data": nil})
		return
	}
	if existing.Role == 1 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "无法添加为好友", "data": nil})
		return
	}
	if existing.Role == 2 {
		c.JSON(http.StatusOK, gin.H{"code": 201, "msg": "无法添加为好友", "data": nil})
		return
	}
	if existing.Role == 3 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "已经添加过了", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "异常操作", "data": nil})
}

func (h *CoachHandler) StudentToFriend(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "缺少参数", "data": nil})
		return
	}
	var friend models.Coach
	if err := h.db.Where("id = ?", body.Id).First(&friend).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "用户不存在", "data": nil})
		return
	}
	var existing models.CoachRelationship
	if err := h.db.Where("coach_id = ? AND student_id = ?", uid, friend.Id).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "记录不存在", "data": nil})
		return
	}
	if existing.Role == int(models.RoleCoachStudent) {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "无法添加为好友", "data": nil})
		return
	}
	if existing.Role == int(models.RoleFriendAndFriend) {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "已经是好友了", "data": nil})
		return
	}
	if existing.Role == int(models.RoleCoachAndStudentHasAccount) {
		if err := h.db.Model(&existing).Update("role", int(models.RoleFriendAndFriend)).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "操作成功", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "异常操作", "data": nil})
}

// 我关注别人
func (h *CoachHandler) FollowCoach(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		FollowingId int `json:"following_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if body.FollowingId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "缺少关注人信息", "data": nil})
		return
	}
	var existing models.CoachFollow
	if err := h.db.Where("following_id = ? AND follower_id = ?", body.FollowingId, uid).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		the_created := models.CoachFollow{
			Status:      1,
			FollowingId: body.FollowingId,
			FollowerId:  uid,
		}
		if err := h.db.Create(&the_created).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "关注成功", "data": nil})
		return
	}
	if existing.Status == 1 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "已经关注了", "data": nil})
		return
	}
	if err := h.db.Model(&existing).Update("status", 1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "关注成功", "data": nil})
}

// 我取消关注别人
func (h *CoachHandler) UnFollowCoach(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		FollowingId int `json:"following_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if body.FollowingId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "缺少关注人信息", "data": nil})
		return
	}
	var existing models.CoachFollow
	if err := h.db.Where("following_id = ? AND follower_id = ?", body.FollowingId, uid).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "没有关注记录", "data": nil})
		return
	}
	if existing.Status == 2 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "已经取消关注了", "data": nil})
		return
	}
	if err := h.db.Model(&existing).Update("status", 2).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "取消关注成功", "data": nil})
}

// 获取我的关注者列表
func (h *CoachHandler) FetchMyFollowerList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		models.Pagination
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	query := h.db
	query = query.Where("following_id = ?", uid)
	pb := pagination.NewPaginationBuilder[models.CoachFollow](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")
	var list1 []models.CoachFollow
	if err := pb.Build().Preload("Follower.Profile1").Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch discount policies", "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
	list := make([]map[string]interface{}, 0, len(list2))
	for _, v := range list2 {
		list = append(list, map[string]interface{}{
			"id":         v.FollowerId,
			"nickname":   v.Follower.Profile1.Nickname,
			"avatar_url": v.Follower.Profile1.AvatarURL,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": gin.H{
			"list":        list,
			"page_size":   pb.GetLimit(),
			"has_more":    has_more,
			"next_marker": next_marker,
		},
	})
}

// 获取我关注的人列表
func (h *CoachHandler) FetchMyFollowingList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		models.Pagination
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	query := h.db
	query = query.Where("follower_id = ?", uid)
	pb := pagination.NewPaginationBuilder[models.CoachFollow](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")
	var list1 []models.CoachFollow
	if err := pb.Build().Preload("Following.Profile1").Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch discount policies", "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
	list := make([]map[string]interface{}, 0, len(list2))
	for _, v := range list2 {
		list = append(list, map[string]interface{}{
			"id":         v.FollowingId,
			"nickname":   v.Following.Profile1.Nickname,
			"avatar_url": v.Following.Profile1.AvatarURL,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": gin.H{
			"list":        list,
			"page_size":   pb.GetLimit(),
			"has_more":    has_more,
			"next_marker": next_marker,
		},
	})
}

func (h *CoachHandler) FetchCoachProfileInWechat(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "请先登录", "data": nil})
		return
	}
	var body struct {
		UID string `json:"uid"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if body.UID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "缺少参数", "data": nil})
		return
	}
	var coach models.Coach
	if err := h.db.
		Where("nickname = ?", body.UID).
		Preload("Profile1").
		First(&coach).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "没有找到记录", "data": nil})
		return
	}
	var platform_accounts []models.CoachMediaSocialAccount
	if err := h.db.Where("coach_id = ?", coach.Id).Find(&platform_accounts).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": gin.H{
		"id":         coach.Id,
		"uid":        coach.Nickname,
		"nickname":   coach.Profile1.Nickname,
		"avatar_url": coach.Profile1.AvatarURL,
		"accounts":   platform_accounts,
	}})
}

func (h *CoachHandler) FetchCoachProfileInAdmin(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid != 1 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有权限", "data": nil})
		return
	}
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "缺少 id 参数", "data": nil})
		return
	}
	var coach models.Coach
	if err := h.db.
		Where("id = ?", body.Id).
		Preload("Profile1").
		First(&coach).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "没有找到记录", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": gin.H{
		"id":         coach.Id,
		"nickname":   coach.Profile1.Nickname,
		"avatar_url": coach.Profile1.AvatarURL,
	}})
}

func (h *CoachHandler) FetchCoachContentList(c *gin.Context) {
	// uid := int(c.GetFloat64("id"))
	var body struct {
		models.Pagination
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	query := h.db
	pb := pagination.NewPaginationBuilder[models.CoachContent](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")
	var list1 []models.CoachContent
	if err := pb.Build().Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch discount policies", "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
	// list := make([]map[string]interface{}, 0, len(list2))
	// for _, v := range list2 {
	// 	list = append(list, map[string]interface{}{
	// 		"id":         v.FollowingId,
	// 		"nickname":   v.Following.Profile1.Nickname,
	// 		"avatar_url": v.Following.Profile1.AvatarURL,
	// 	})
	// }
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": gin.H{
			"list":        list2,
			"page_size":   pb.GetLimit(),
			"has_more":    has_more,
			"next_marker": next_marker,
		},
	})
}
