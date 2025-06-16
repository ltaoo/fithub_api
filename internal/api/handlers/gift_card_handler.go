package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myapi/internal/models"
	"myapi/internal/pkg/pagination"
	"myapi/pkg/logger"
)

const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // 无歧义字符集

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateGiftCode 生成8位无歧义礼品码
func GenerateGiftCode() string {
	code := make([]byte, 8) // 预分配8字节切片，提升性能
	for i := range code {
		// 生成 [0, len(charset)-1] 范围内的随机索引
		randomIndex := rnd.Intn(len(charset))
		code[i] = charset[randomIndex]
	}
	return string(code)
}

type GiftCardHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewGiftCardHandler(db *gorm.DB, logger *logger.Logger) *GiftCardHandler {
	return &GiftCardHandler{
		db:     db,
		logger: logger,
	}
}

func (h *GiftCardHandler) FetchGiftCardList(c *gin.Context) {
	var body struct {
		models.Pagination
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	query := h.db
	pb := pagination.NewPaginationBuilder[models.GiftCard](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")
	var list1 []models.GiftCard
	if err := pb.Build().Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch discount policies", "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
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

func (h *GiftCardHandler) FetchGiftCardRewardList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		models.Pagination
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	query := h.db
	query = query.Where("d != 1 AND creator_id = ?", uid)
	pb := pagination.NewPaginationBuilder[models.GiftCardReward](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")
	var list1 []models.GiftCardReward
	if err := pb.Build().Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
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

func (h *GiftCardHandler) CreateGiftCard(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Num              int `json:"num"`
		GiftCardRewardId int `json:"gift_card_reward_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	// 验证参数
	if body.GiftCardRewardId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少参数", "data": nil})
		return
	}
	if body.Num <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "数量必须大于0", "data": nil})
		return
	}
	// 开始事务
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	if tx.Error != nil {
		h.logger.Error("Failed to start transaction", tx.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		return
	}

	// 批量创建记录
	records := make([]models.GiftCard, body.Num)
	now := time.Now()
	for i := 0; i < body.Num; i++ {
		records[i] = models.GiftCard{
			Code:             GenerateGiftCode(),
			Status:           int(models.GiftCardStatusUnused),
			CreatorId:        uid,
			GiftCardRewardId: body.GiftCardRewardId,
			CreatedAt:        now,
		}
	}

	if err := tx.Create(&records).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create records", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create records", "data": nil})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to commit transaction", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to commit transaction", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": nil,
	})
}

func (h *GiftCardHandler) CreateGiftCardReward(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Name     string `json:"name"`
		Overview string `json:"overview"`
		Details  string `json:"details"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	// 验证参数
	if body.Name == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少参数", "data": nil})
		return
	}
	if body.Details == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少参数", "data": nil})
		return
	}
	// @todo 验证 Details 的有效性
	// 开始事务
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	if tx.Error != nil {
		h.logger.Error("Failed to start transaction", tx.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		return
	}
	now := time.Now()
	record := models.GiftCardReward{
		Name:      body.Name,
		Overview:  body.Name,
		Details:   body.Details,
		Status:    1,
		CreatorId: uid,
		CreatedAt: now,
	}
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create records", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create records", "data": nil})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to commit transaction", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to commit transaction", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": nil,
	})
}

type GiftCardRewardDetailsJSON250607 struct {
	SubscriptionPlanId int `json:"subscription_plan_id"`
	DayCount           int `json:"day_count"`
}

func (h *GiftCardHandler) FetchGiftCardProfile(c *gin.Context) {
	var body struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	var record models.GiftCard
	if r := h.db.Where("code = ? AND d != 1", body.Code).Preload("GiftCardReward").First(&record); r.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": r.Error.Error(), "data": nil})
		return
	}
	if record.GiftCardReward.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "异常数据", "data": nil})
		return
	}
	var details GiftCardRewardDetailsJSON250607
	if err := json.Unmarshal([]byte(record.GiftCardReward.Details), &details); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 600, "msg": "异常数据", "data": nil})
		return
	}
	if details.SubscriptionPlanId == 0 || details.DayCount == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 600, "msg": "异常数据", "data": nil})
		return
	}
	var subscription_plan models.SubscriptionPlan
	if err := h.db.Where("id = ?", details.SubscriptionPlanId).First(&subscription_plan).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 600, "msg": "异常数据", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"name":   record.GiftCardReward.Name,
			"status": record.Status,
		},
	})
}

func (h *GiftCardHandler) UsingGiftCard(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	if tx.Error != nil {
		h.logger.Error("Failed to start transaction", tx.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		return
	}
	var card models.GiftCard
	if r := tx.Where("code = ? AND d != 1", body.Code).Preload("GiftCardReward").First(&card); r.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": r.Error.Error(), "data": nil})
		return
	}
	if card.Status == int(models.GiftCardStatusUsed) {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "兑换码已被使用", "data": nil})
		return
	}
	if card.GiftCardReward.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 600, "msg": "异常数据", "data": nil})
		return
	}
	var details GiftCardRewardDetailsJSON250607
	if err := json.Unmarshal([]byte(card.GiftCardReward.Details), &details); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 601, "msg": "异常数据", "data": nil})
		return
	}
	if details.SubscriptionPlanId == 0 || details.DayCount == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 602, "msg": "异常数据", "data": nil})
		return
	}
	var subscription_plan models.SubscriptionPlan
	if err := tx.Where("id = ?", details.SubscriptionPlanId).First(&subscription_plan).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 603, "msg": "异常数据", "data": nil})
		return
	}
	now := time.Now()

	// 检查是否有生效中的订阅
	var active_subscription models.Subscription
	has_active_subscription := tx.Where("coach_id = ? AND step = 2 AND expired_at IS NULL", uid).First(&active_subscription).Error == nil

	// 创建新订阅
	subscription := models.Subscription{
		Step:               1,
		Count:              details.DayCount,
		Reason:             "使用礼品卡兑换",
		CreatedAt:          now,
		SubscriptionPlanId: subscription_plan.Id,
		CoachId:            uid,
	}

	// 如果没有生效中的订阅，则新订阅立即生效
	if !has_active_subscription {
		subscription.Step = 2
		subscription.ActiveAt = &now
		expired_at := now.AddDate(0, 0, details.DayCount)
		subscription.ExpectExpiredAt = &expired_at
	}

	if err := tx.Create(&subscription).Error; err != nil {
		tx.Rollback()
		h.logger.Error("create subscription failed", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "兑换失败", "data": nil})
		return
	}

	updates := map[string]interface{}{
		"status":      int(models.GiftCardStatusUsed),
		"consumer_id": uid,
		"used_at":     now,
	}
	if err := tx.Model(&card).Updates(updates).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to update gift card", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to commit transaction", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to commit transaction", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "兑换成功",
		"data": nil,
	})
}
