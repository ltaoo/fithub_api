package handlers

import (
	"fmt"
	"myapi/internal/models"
	"myapi/internal/pkg/pagination"
	"myapi/pkg/logger"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewSubscriptionHandler(db *gorm.DB, logger *logger.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		db:     db,
		logger: logger,
	}
}

func (h *SubscriptionHandler) FetchSubscriptionPlanList(c *gin.Context) {
	var body struct {
		models.Pagination
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	query := h.db
	var plans []models.SubscriptionPlan
	if err := query.Find(&plans).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch discount policies", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": gin.H{
			"list": plans,
		},
	})
}

func (h *SubscriptionHandler) CreateSubscriptionPlan(c *gin.Context) {
	var body struct {
		Name             string `json:"name"`
		Details          string `json:"details"`
		UnitPrice        int    `json:"unit_price"`
		DiscountPolicies []struct {
			Name         string `json:"name"`
			Rate         int    `json:"rate"`
			CountRequire int    `json:"count_require"`
			Enabled      int    `json:"enabled"`
		} `json:"discount_policies"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// 开始事务
	tx := h.db.Begin()
	if tx.Error != nil {
		h.logger.Error("Failed to start transaction", tx.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		return
	}

	// 创建订阅计划
	subscription_plan := models.SubscriptionPlan{
		Name:      body.Name,
		Details:   body.Details,
		UnitPrice: body.UnitPrice,
		CreatedAt: time.Now(),
	}

	if err := tx.Create(&subscription_plan).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create subscription plan", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create subscription plan", "data": nil})
		return
	}

	// 创建折扣政策并关联
	for _, policy := range body.DiscountPolicies {
		// 创建折扣政策
		discount_policy := models.DiscountPolicy{
			Name:         policy.Name,
			Rate:         policy.Rate,
			CountRequire: policy.CountRequire,
			CreatedAt:    time.Now(),
		}

		if err := tx.Create(&discount_policy).Error; err != nil {
			tx.Rollback()
			h.logger.Error("Failed to create discount policy", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create discount policy", "data": nil})
			return
		}

		// 创建订阅计划和折扣政策的关联
		plan_policy := models.SubscriptionPlanDiscountPolicy{
			SubscriptionPlanId: subscription_plan.Id,
			DiscountPolicyId:   discount_policy.Id,
			Enabled:            policy.Enabled,
		}

		if err := tx.Create(&plan_policy).Error; err != nil {
			tx.Rollback()
			h.logger.Error("Failed to create subscription plan discount policy", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create subscription plan discount policy", "data": nil})
			return
		}
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
		"data": subscription_plan,
	})
}

func (h *SubscriptionHandler) UpdateSubscriptionPlan(c *gin.Context) {
	var body struct {
		Id               int    `json:"id"`
		Name             string `json:"name"`
		Details          string `json:"details"`
		UnitPrice        int    `json:"unit_price"`
		DiscountPolicies []struct {
			Id           int    `json:"id"`
			Name         string `json:"name"`
			Rate         int    `json:"rate"`
			CountRequire int    `json:"count_require"`
			Enabled      int    `json:"enabled"`
		} `json:"discount_policies"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// 开始事务
	tx := h.db.Begin()
	if tx.Error != nil {
		h.logger.Error("Failed to start transaction", tx.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		return
	}

	// 更新订阅计划
	subscription_plan := models.SubscriptionPlan{
		Id:        body.Id,
		Name:      body.Name,
		Details:   body.Details,
		UnitPrice: body.UnitPrice,
	}

	if err := tx.Model(&models.SubscriptionPlan{}).Where("id = ?", body.Id).Updates(subscription_plan).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to update subscription plan", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update subscription plan", "data": nil})
		return
	}

	// 删除旧的折扣政策关联
	if err := tx.Where("subscription_plan_id = ?", body.Id).Delete(&models.SubscriptionPlanDiscountPolicy{}).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to delete old discount policies", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to delete old discount policies", "data": nil})
		return
	}

	// 更新折扣政策
	for _, policy := range body.DiscountPolicies {
		// 更新或创建折扣政策
		discount_policy := models.DiscountPolicy{
			Id:           policy.Id,
			Name:         policy.Name,
			Rate:         policy.Rate,
			CountRequire: policy.CountRequire,
		}

		if err := tx.Save(&discount_policy).Error; err != nil {
			tx.Rollback()
			h.logger.Error("Failed to update discount policy", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update discount policy", "data": nil})
			return
		}

		// 创建新的订阅计划和折扣政策的关联
		plan_policy := models.SubscriptionPlanDiscountPolicy{
			SubscriptionPlanId: body.Id,
			DiscountPolicyId:   discount_policy.Id,
			Enabled:            policy.Enabled,
		}

		if err := tx.Create(&plan_policy).Error; err != nil {
			tx.Rollback()
			h.logger.Error("Failed to create subscription plan discount policy", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create subscription plan discount policy", "data": nil})
			return
		}
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
		"data": subscription_plan,
	})
}

func (h *SubscriptionHandler) CalcSubscriptionOrderAmount(c *gin.Context) {
	var body struct {
		SubscriptionPlanId int    `json:"subscription_plan_id"`
		Type               string `json:"type"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	count := 30
	if body.Type == "season" {
		count = 30 * 4
	}
	if body.Type == "year" {
		count = 30 * 12
	}

	var plan models.SubscriptionPlan
	if err := h.db.First(&plan, body.SubscriptionPlanId).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch subscription plan", "data": nil})
		return
	}

	// Get discount policies
	var discount_policies []models.SubscriptionPlanDiscountPolicy
	if err := h.db.Where("subscription_plan_id = ? AND enabled = 1", body.SubscriptionPlanId).
		Preload("DiscountPolicy").
		Find(&discount_policies).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch discount policies", "data": nil})
		return
	}

	// Sort discount policies by count_require in descending order
	sort.Slice(discount_policies, func(i, j int) bool {
		return discount_policies[i].DiscountPolicy.CountRequire > discount_policies[j].DiscountPolicy.CountRequire
	})

	// Find applicable discount rate
	rate := 100
	for _, policy := range discount_policies {
		if count >= policy.DiscountPolicy.CountRequire {
			rate = policy.DiscountPolicy.Rate
			break
		}
	}

	// Calculate prices
	total_amount := count * plan.UnitPrice
	amount := total_amount * rate / 100
	discount := total_amount - amount

	// Generate text
	month := float64(count) / 30
	text := fmt.Sprintf("购买%.1f年共计%.2f元", float64(count)/360, float64(amount)/10)
	if month < 12 {
		text = fmt.Sprintf("购买%.0f个月共计%.2f元", month, float64(amount)/10)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": gin.H{
			"total_amount": total_amount,
			"amount":       amount,
			"discount_texts": []gin.H{
				{
					"value": discount,
					"text":  fmt.Sprintf("满%.1f年打%.2f折", float64(count)/360, float64(rate)/100),
				},
			},
			"text": text,
		},
	})
}

func (h *SubscriptionHandler) FetchSubscriptionList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		models.Pagination
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	query := h.db
	query = query.Where("coach_id = ?", uid).Preload("SubscriptionPlan")
	pb := pagination.NewPaginationBuilder[models.Subscription](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")
	var list []models.Subscription
	if err := pb.Build().Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list)
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
