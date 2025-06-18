package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myapi/internal/models"
	"myapi/internal/pkg/pagination"
	"myapi/pkg/logger"
)

// MuscleHandler handles HTTP requests for muscles
type WorkoutActionHistoryHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewWorkoutActionHandler creates a new workout action handler
func NewWorkoutActionHistoryHandler(db *gorm.DB, logger *logger.Logger) *WorkoutActionHistoryHandler {
	return &WorkoutActionHistoryHandler{
		db:     db,
		logger: logger,
	}
}

func (h *WorkoutActionHistoryHandler) CreateWorkoutHistory(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		WorkoutActionId int    `json:"workout_action_id"`
		Reps            int    `json:"reps"`
		RepsUnit        string `json:"reps_unit"`
		Weight          int    `json:"weight"`
		WeightUnit      string `json:"weight_unit"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.WorkoutActionId == 0 || body.Reps == 0 || body.Weight == 0 || body.RepsUnit == "" || body.WeightUnit == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少动作参数", "data": nil})
		return
	}
	if uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "参数错误", "data": nil})
		return
	}

	history := models.WorkoutActionHistory{
		StudentId:       uid,
		WorkoutActionId: body.WorkoutActionId,
		Reps:            body.Reps,
		RepsUnit:        body.RepsUnit,
		Weight:          body.Weight,
		WeightUnit:      body.WeightUnit,
		CreatedAt:       time.Now(),
	}

	if err := h.db.Create(&history).Error; err != nil {
		h.logger.Error("Failed to create workout action history", err)
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": nil})
}

func (h *WorkoutActionHistoryHandler) FetchWorkoutActionHistoryListOfWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		models.Pagination
		WorkoutDayId int    `json:"workout_day_id"`
		OrderBy      string `json:"order_by"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.WorkoutDayId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少参数", "data": nil})
		return
	}
	var d models.WorkoutDay
	if err := h.db.Where("id = ? AND student_id = ?", body.WorkoutDayId, uid).First(&d).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有找到记录", "data": nil})
		return
	}

	query := h.db
	query = query.Where("workout_day_id = ?", body.WorkoutDayId)
	pb := pagination.NewPaginationBuilder[models.WorkoutActionHistory](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")

	var list1 []models.WorkoutActionHistory
	if err := pb.Build().Preload("WorkoutAction").Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout history: " + err.Error(), "data": nil})
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

// 获取健身动作历史记录
func (h *WorkoutActionHistoryHandler) FetchWorkoutActionHistoryListOfWorkoutAction(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		models.Pagination
		WorkoutActionId int    `json:"workout_action_id"`
		StudentId       int    `json:"student_id"`
		OrderBy         string `json:"order_by"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.WorkoutActionId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少参数", "data": nil})
		return
	}
	if body.StudentId == 0 {
		body.StudentId = uid
	}
	if uid != body.StudentId {
		var relation models.CoachRelationship
		if err := h.db.Where("coach_id = ? AND student_id = ?", uid, body.StudentId).First(&relation).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
			return
		}
	}
	query := h.db
	query = query.Where("action_id = ? AND student_id = ?", body.WorkoutActionId, body.StudentId)
	pb := pagination.NewPaginationBuilder[models.WorkoutActionHistory](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")
	if body.OrderBy != "" {
		pb.SetOrderBy(body.OrderBy)
	}
	var list1 []models.WorkoutActionHistory
	if err := pb.Build().Preload("WorkoutAction").Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout history: " + err.Error(), "data": nil})
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
