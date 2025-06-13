package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myapi/internal/models"
	"myapi/internal/pkg/pagination"
	"myapi/pkg/logger"
)

// 设备列表，可以参考这个
// https://repfitness.com/collections/best-sellers?page=2

type EquipmentHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewWorkoutActionHandler creates a new workout action handler
func NewEquipmentHandler(db *gorm.DB, logger *logger.Logger) *EquipmentHandler {
	return &EquipmentHandler{
		db:     db,
		logger: logger,
	}
}

func (h *EquipmentHandler) FetchEquipmentList(c *gin.Context) {
	var body struct {
		models.Pagination
		Ids []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	query := h.db
	if len(body.Ids) != 0 {
		query = query.Where("id IN (?)", body.Ids)
	}
	pb := pagination.NewPaginationBuilder[models.Equipment](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("sort_idx DESC")
	var list1 []models.Equipment
	if err := query.Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch equipments" + err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        list2,
			"page_size":   pb.GetLimit(),
			"has_more":    has_more,
			"next_marker": next_marker,
		},
	})
}

func (h *EquipmentHandler) FetchEquipment(c *gin.Context) {
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	var equipment models.Equipment
	if err := h.db.First(&equipment, body.Id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "没有找到记录", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": equipment})
}

func (h *EquipmentHandler) CreateEquipment(c *gin.Context) {
	var body struct {
		ZhName   string `json:"zh_name"`
		Name     string `json:"name"`
		Alias    string `json:"alias"`
		Overview string `json:"overview"`
		Tags     string `json:"tags"`
		SortIdx  int    `json:"sort_idx"`
		Medias   string `json:"medias"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	record := models.Equipment{
		Name:     body.Name,
		ZhName:   body.ZhName,
		Alias:    body.Alias,
		Overview: body.Overview,
		// Tags:     body.Tags,
		// SortIdx:  body.SortIdx,
		Medias: body.Medias,
	}
	if err := h.db.Create(&record).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": nil})
}

func (h *EquipmentHandler) UpdateEquipment(c *gin.Context) {
	var body struct {
		Id       int    `json:"id"`
		ZhName   string `json:"zh_name"`
		Name     string `json:"name"`
		Alias    string `json:"alias"`
		Overview string `json:"overview"`
		Tags     string `json:"tags"`
		SortIdx  int    `json:"sort_idx"`
		Medias   string `json:"medias"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	var existing models.Equipment
	if err := h.db.First(&existing, body.Id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "没有找到记录", "data": nil})
		return
	}
	updates := map[string]interface{}{}
	if body.ZhName != "" {
		updates["zh_name"] = body.ZhName
	}
	if body.Name != "" {
		updates["name"] = body.Name
	}
	if body.Alias != "" {
		updates["alias"] = body.Alias
	}
	if body.Overview != "" {
		updates["overview"] = body.Overview
	}
	if body.SortIdx != 0 {
		updates["sort_idx"] = body.SortIdx
	}
	if body.Medias != "" {
		updates["medias"] = body.Medias
	}
	if err := h.db.Model(&existing).Updates(&updates).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": nil})
}

func (h *EquipmentHandler) DeleteEquipment(c *gin.Context) {
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	var existing models.Equipment
	if err := h.db.First(&existing, body.Id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "没有找到记录", "data": nil})
		return
	}
	if err := h.db.Model(&existing).Update("d", 1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功", "data": nil})
}
