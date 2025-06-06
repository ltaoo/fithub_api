package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myapi/internal/models"
	"myapi/pkg/logger"
)

type ReportHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewReportHandler(db *gorm.DB, logger *logger.Logger) *ReportHandler {
	return &ReportHandler{
		db:     db,
		logger: logger,
	}
}

func (h *ReportHandler) FetchReportList(c *gin.Context) {
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
	query.Order("created_at DESC").Limit(limit + 1)
	var list []models.CoachReport
	r := query.Find(&list)
	if r.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch discount policies", "data": nil})
		return
	}

	has_more := false
	next_cursor := ""

	if len(list) > limit {
		has_more = true
		list = list[:limit]                               // Remove the extra item we fetched
		next_cursor = strconv.Itoa(int(list[limit-1].Id)) // Get the last item's ID as next cursor
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": gin.H{
			"list":        list,
			"page_size":   limit,
			"has_more":    has_more,
			"next_marker": next_cursor,
		},
	})
}

func (h *ReportHandler) FetchMineReportList(c *gin.Context) {
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
	query.Order("created_at DESC").Limit(limit + 1)
	var list []models.CoachReport
	r := query.Where("d != 0 AND coach_id = ?", uid).Find(&list)
	if r.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch discount policies", "data": nil})
		return
	}

	has_more := false
	next_cursor := ""

	if len(list) > limit {
		has_more = true
		list = list[:limit]                               // Remove the extra item we fetched
		next_cursor = strconv.Itoa(int(list[limit-1].Id)) // Get the last item's ID as next cursor
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": gin.H{
			"list":        list,
			"page_size":   limit,
			"has_more":    has_more,
			"next_marker": next_cursor,
		},
	})
}

func (h *ReportHandler) CreateReport(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Content    string  `json:"content"`
		ReasonType *string `json:"reason_type"`
		ReasonId   *int    `json:"reason_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	// 验证 Content 字段
	if body.Content == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Content cannot be empty", "data": nil})
		return
	}

	// 开始事务
	tx := h.db.Begin()
	if tx.Error != nil {
		h.logger.Error("Failed to start transaction", tx.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		return
	}

	// 创建
	record := models.CoachReport{
		Content:    body.Content,
		ReasonType: "common",
		ReasonId:   0,
		CoachId:    uid,
		CreatedAt:  time.Now(),
	}

	// 如果提供了 ReasonType，则设置它
	if body.ReasonType != nil {
		record.ReasonType = *body.ReasonType
	}

	// 如果提供了 ReasonId，则设置它
	if body.ReasonId != nil {
		record.ReasonId = *body.ReasonId
	}

	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create record", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create record", "data": nil})
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
		"data": record,
	})
}

func (h *ReportHandler) FetchReportProfile(c *gin.Context) {
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	var record models.CoachReport

	if r := h.db.Where("id = ?", body.Id).First(&record); r.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": r.Error.Error(), "data": nil})
		return
	}

	if record.D == 1 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "记录不存在", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": record,
	})
}
