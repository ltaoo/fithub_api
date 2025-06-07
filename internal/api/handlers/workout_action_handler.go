package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myapi/internal/models"
	"myapi/internal/pkg/pagination"
	"myapi/pkg/logger"
)

// WorkoutActionHandler handles HTTP requests for workout actions
type WorkoutActionHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewWorkoutActionHandler creates a new workout action handler
func NewWorkoutActionHandler(db *gorm.DB, logger *logger.Logger) *WorkoutActionHandler {
	return &WorkoutActionHandler{
		db:     db,
		logger: logger,
	}
}

// FetchWorkoutActionList retrieves all workout actions
func (h *WorkoutActionHandler) FetchWorkoutActionList(c *gin.Context) {
	var body struct {
		models.Pagination
		Type    string `json:"type"`
		Keyword string `json:"keyword"`
		Tag     string `json:"tag"`
		Level   string `json:"level"`
		Muscle  string `json:"muscle"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// Start with base query
	query := h.db

	// Apply filters if provided
	if body.Type != "" {
		query = query.Where("type = ?", body.Type)
	}
	if body.Keyword != "" {
		query = query.Where("zh_name LIKE ? OR alias LIKE ?", "%"+body.Keyword+"%", "%"+body.Keyword+"%")
	}
	if body.Level != "" {
		levelInt, err := strconv.Atoi(body.Level)
		if err == nil {
			query = query.Where("level = ?", levelInt)
		}
	}
	if body.Tag != "" {
		query = query.Where("tags1 LIKE ?", "%"+body.Tag+"%")
	}
	if body.Muscle != "" {
		query = query.Where("target_muscle_ids LIKE ?", "%"+body.Muscle+"%")
	}

	// Build paginated query
	pb := pagination.NewPaginationBuilder[models.WorkoutAction](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetNextMarker(body.NextMarker).
		SetOrderBy("sort_idx DESC")

	// Execute the query
	var list []models.WorkoutAction
	if err := pb.Build().Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout actions", "data": nil})
		return
	}

	// Process results
	list_data, has_more, next_cursor := pb.ProcessResults(list)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        list_data,
			"page_size":   pb.GetLimit(),
			"has_more":    has_more,
			"next_marker": next_cursor,
		},
	})
}

// GetWorkoutActionList retrieves all workout actions
func (h *WorkoutActionHandler) GetWorkoutActionListByIds(c *gin.Context) {
	var actions []models.WorkoutAction

	type WorkoutActionListRequest struct {
		Ids []int `json:"ids"`
	}

	var request WorkoutActionListRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	// Start with base query
	query := h.db

	// 按创建时间排序
	query = query.Order("created_at DESC")

	// Apply filters if provided
	if len(request.Ids) > 0 {
		query = query.Where("id IN (?)", request.Ids)
	}
	// Execute the query
	result := query.Find(&actions)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout actions", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list": actions,
		},
	})
}

// GetWorkoutAction retrieves a specific workout action by ID
func (h *WorkoutActionHandler) GetWorkoutAction(c *gin.Context) {
	// id := c.Param("id")
	var request struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	var action models.WorkoutAction
	result := h.db.First(&action, request.Id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout action not found", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": action})
}

// GetActionsByMuscle retrieves workout actions targeting a specific muscle
func (h *WorkoutActionHandler) GetActionsByMuscle(c *gin.Context) {
	muscleId := c.Param("muscleId")

	var actions []models.WorkoutAction
	result := h.db.Where("target_muscle_ids LIKE ?", "%"+muscleId+"%").Find(&actions)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout actions", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": actions})
}

func (h *WorkoutActionHandler) FetchWorkoutActionsByLevel(c *gin.Context) {
	levelStr := c.Param("level")
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid level parameter", "data": nil})
		return
	}

	var actions []models.WorkoutAction
	result := h.db.Where("level = ?", level).Find(&actions)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout actions", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": actions})
}

// 获取指定动作的 进阶、退阶、替代动作
func (h *WorkoutActionHandler) FetchRelatedWorkoutActions(c *gin.Context) {
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	var action models.WorkoutAction
	r := h.db.First(&action, body.Id)
	if r.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout action not found", "data": nil})
		return
	}
	// Get advanced actions
	var advanced_workout_actions []models.WorkoutAction
	if action.AdvancedActionIds != "" {
		advanced_ids := strings.Split(action.AdvancedActionIds, ",")
		h.db.Where("id IN ?", advanced_ids).Find(&advanced_workout_actions)
	}
	// Get regressed actions
	var regressed_actions []models.WorkoutAction
	if action.RegressedActionIds != "" {
		regressed_ids := strings.Split(action.RegressedActionIds, ",")
		h.db.Where("id IN ?", regressed_ids).Find(&regressed_actions)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"advanced":  advanced_workout_actions,
			"regressed": regressed_actions,
		},
	})
}

// WorkoutActionBody represents the request body for workout action operations
type WorkoutActionBody struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name" bindings:"min=3"`
	ZhName               string `json:"zh_name"`
	Alias                string `json:"alias"`
	Overview             string `json:"overview"`
	CoverURL             string `json:"cover_url"`
	Type                 string `json:"type"`
	Level                int    `json:"level"`
	SortIdx              int    `json:"sort_idx"`
	Pattern              string `json:"pattern"`
	Score                int    `json:"score"`
	Tags1                string `json:"tags1"`
	Tags2                string `json:"tags2"`
	Details              string `json:"details"`
	Points               string `json:"points"`
	Problems             string `json:"problems"`
	ExtraConfig          string `json:"extra_config"`
	EquipmentIds         string `json:"equipment_ids"`
	MuscleIds            string `json:"muscle_ids"`
	PrimaryMuscleIds     string `json:"primary_muscle_ids"`
	SecondaryMuscleIds   string `json:"secondary_muscle_ids"`
	AlternativeActionIds string `json:"alternative_action_ids"`
	AdvancedActionIds    string `json:"advanced_action_ids"`
	RegressedActionIds   string `json:"regressed_action_ids"`
}

func (h *WorkoutActionHandler) CreateWorkoutAction(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "请先登录", "data": nil})
		return
	}
	var body WorkoutActionBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	var existing models.WorkoutAction
	if err := h.db.Where("zh_name = ?", body.ZhName).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
			return
		}
		if existing.Id != 0 {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "已存在同名动作", "data": gin.H{
				"id": existing.Id,
			}})
			return
		}
	}
	data := models.WorkoutAction{
		Name:                 body.Name,
		ZhName:               body.ZhName,
		Alias:                body.Alias,
		Overview:             body.Overview,
		Type:                 body.Type,
		Level:                body.Level,
		Status:               1,
		Score:                body.Score,
		SortIdx:              body.SortIdx,
		Pattern:              body.Pattern,
		Tags1:                body.Tags1,
		Tags2:                body.Tags2,
		Details:              body.Details,
		Points:               body.Points,
		Problems:             body.Problems,
		EquipmentIds:         body.EquipmentIds,
		MuscleIds:            body.MuscleIds,
		PrimaryMuscleIds:     body.PrimaryMuscleIds,
		SecondaryMuscleIds:   body.SecondaryMuscleIds,
		AlternativeActionIds: body.AlternativeActionIds,
		AdvancedActionIds:    body.AdvancedActionIds,
		RegressedActionIds:   body.RegressedActionIds,
		ExtraConfig:          body.ExtraConfig,
		OwnerId:              uid,
		CreatedAt:            time.Now(),
	}
	if err := h.db.Create(&data).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "the record created successfully", "data": body})
}

func (h *WorkoutActionHandler) UpdateWorkoutActionProfile(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "请先登录", "data": nil})
		return
	}
	var body WorkoutActionBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	var existing models.WorkoutAction
	if err := h.db.First(&existing, body.Id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "the record not found", "data": nil})
		return
	}
	now := time.Now().UTC()
	if err := h.db.Model(&existing).Updates(models.WorkoutAction{
		Name:                 body.Name,
		ZhName:               body.ZhName,
		Alias:                body.Alias,
		Overview:             body.Overview,
		Type:                 body.Type,
		Level:                body.Level,
		Score:                body.Score,
		SortIdx:              body.SortIdx,
		Pattern:              body.Pattern,
		Tags1:                body.Tags1,
		Tags2:                body.Tags2,
		Details:              body.Details,
		Points:               body.Points,
		Problems:             body.Problems,
		EquipmentIds:         body.EquipmentIds,
		MuscleIds:            body.MuscleIds,
		PrimaryMuscleIds:     body.PrimaryMuscleIds,
		SecondaryMuscleIds:   body.SecondaryMuscleIds,
		AlternativeActionIds: body.AlternativeActionIds,
		AdvancedActionIds:    body.AdvancedActionIds,
		RegressedActionIds:   body.RegressedActionIds,
		ExtraConfig:          body.ExtraConfig,
		OwnerId:              uid,
		UpdatedAt:            &now,
	}).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "updated successfully", "data": body})
}

// 删除动作
func (h *WorkoutActionHandler) DeleteWorkoutAction(c *gin.Context) {
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	var action models.WorkoutAction
	r := h.db.First(&action, body.Id)
	if r.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "the record not found", "data": nil})
		return
	}
	r = h.db.Delete(&action)
	if r.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to delete the record", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "the record deleted successfully", "data": nil})
}
