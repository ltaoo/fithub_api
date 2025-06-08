package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"gorm.io/gorm"

	"myapi/internal/models"
	"myapi/internal/pkg/pagination"
	"myapi/pkg/logger"
)

type MediaResourceHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewMediaResourceHandler(db *gorm.DB, logger *logger.Logger) *MediaResourceHandler {
	return &MediaResourceHandler{
		db:     db,
		logger: logger,
	}
}

// qiniu.region.z0: 代表华东区域
// qiniu.region.z1: 代表华北区域
// qiniu.region.z2: 代表华南区域
// qiniu.region.na0: 代表北美区域
// qiniu.region.as0: 代表新加坡区域
func (h *MediaResourceHandler) BuildQiniuToken(c *gin.Context) {
	access_key := "HriJcPNneVwy8gZWB_QAB-hxswRtk1zFWSbYVlUu"
	secret_key := "Ku0DYLqYWO6GLhKiixvTLFeluU0hSc7itJvc6eIJ"
	mac := credentials.NewCredentials(access_key, secret_key)
	bucket := "fithub"
	putPolicy, err := uptoken.NewPutPolicy(bucket, time.Now().Add(1*time.Hour))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	token, err := uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": gin.H{
		"token": token,
	}})
}

func (h *MediaResourceHandler) FetchMediaResourceList(c *gin.Context) {
	var body struct {
		models.Pagination
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	query := h.db
	pb := pagination.NewPaginationBuilder[models.MediaResource](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")
	var list1 []models.MediaResource
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

func (h *MediaResourceHandler) CreateMediaResource(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Type int `json:"type"`
		// 文件信息
		Width    int    `json:"width"`
		Height   int    `json:"height"`
		Duration int    `json:"duration"`
		Size     int    `json:"size"`
		Filetype string `json:"filetype"`
		Filename string `json:"filename"`
		Hash     string `json:"hash"`
		Key      string `json:"key"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Type == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少 type", "data": nil})
		return
	}
	record := models.MediaResource{
		URL:       "//static.fithub.top/" + body.Key,
		MediaType: body.Type,
		Width:     body.Width,
		Height:    body.Height,
		Size:      body.Size,
		Hash:      body.Hash,
		Filetype:  body.Filetype,
		Filename:  body.Filename,
		CreatedAt: time.Now(),
		CreatorId: uid,
	}
	// 验证 Content 字段
	if body.Type == 1 {
	}
	if body.Type == 2 {
		record.Duration = body.Duration
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
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create record", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create record", "data": nil})
		return
	}
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

func (h *MediaResourceHandler) DeleteMediaResource(c *gin.Context) {
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	var existing models.MediaResource
	if err := h.db.Where("id = ?", body.Id).First(&existing).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if err := h.db.Delete(&existing).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功", "data": nil})
}
