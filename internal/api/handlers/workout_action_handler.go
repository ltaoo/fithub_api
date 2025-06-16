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

	query := h.db.Where("d IS NULL OR d = 0")
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
	pb := pagination.NewPaginationBuilder[models.WorkoutAction](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetNextMarker(body.NextMarker).
		SetOrderBy("sort_idx DESC")
	var list1 []models.WorkoutAction
	if err := pb.Build().Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout actions", "data": nil})
		return
	}
	list2, has_more, next_cursor := pb.ProcessResults(list1)
	list := make([]map[string]interface{}, 0, len(list2))
	for _, v := range list2 {
		list = append(list, map[string]interface{}{
			"id":            v.Id,
			"name":          v.Name,
			"zh_name":       v.ZhName,
			"idx":           v.SortIdx,
			"tags":          v.Tags1,
			"muscle_ids":    v.MuscleIds,
			"equipment_ids": v.EquipmentIds,
			"created_at":    v.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        list,
			"page_size":   pb.GetLimit(),
			"has_more":    has_more,
			"next_marker": next_cursor,
		},
	})
}

func (h *WorkoutActionHandler) FetchWorkoutActionListByIds(c *gin.Context) {

	var body struct {
		Ids []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	if len(body.Ids) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "ids 不能为空", "data": nil})
		return
	}
	query := h.db.Where("d IS NULL OR d = 0")
	query = query.Where("id IN (?)", body.Ids)
	var list1 []models.WorkoutAction
	if err := query.Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout actions", "data": nil})
		return
	}
	list := make([]map[string]interface{}, 0, len(list1))
	for _, v := range list1 {
		list = append(list, map[string]interface{}{
			"id":            v.Id,
			"name":          v.Name,
			"zh_name":       v.ZhName,
			"muscle_ids":    v.MuscleIds,
			"equipment_ids": v.EquipmentIds,
			"created_at":    v.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list": list,
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

type WorkoutActionBody struct {
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
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
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
		CoverURL:             body.CoverURL,
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
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": body})
}

func (h *WorkoutActionHandler) UpdateWorkoutActionProfile(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "请先登录", "data": nil})
		return
	}
	var body struct {
		WorkoutActionBody
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	var existing models.WorkoutAction
	if err := h.db.First(&existing, body.Id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "the record not found", "data": nil})
		return
	}
	now := time.Now().UTC()
	if err := h.db.Model(&existing).Updates(map[string]interface{}{
		"name":                   body.Name,
		"zh_name":                body.ZhName,
		"alias":                  body.Alias,
		"overview":               body.Overview,
		"cover_url":              body.CoverURL,
		"type":                   body.Type,
		"level":                  body.Level,
		"score":                  body.Score,
		"sort_idx":               body.SortIdx,
		"pattern":                body.Pattern,
		"tags1":                  body.Tags1,
		"tags2":                  body.Tags2,
		"details":                body.Details,
		"points":                 body.Points,
		"problems":               body.Problems,
		"equipment_ids":          body.EquipmentIds,
		"muscle_ids":             body.MuscleIds,
		"primary_muscle_ids":     body.PrimaryMuscleIds,
		"secondary_muscle_ids":   body.SecondaryMuscleIds,
		"alternative_action_ids": body.AlternativeActionIds,
		"advanced_action_ids":    body.AdvancedActionIds,
		"regressed_action_ids":   body.RegressedActionIds,
		"extra_config":           body.ExtraConfig,
		"owner_id":               uid,
		"updated_at":             now,
	}).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": body})
}

func (h *WorkoutActionHandler) UpdateWorkoutActionIdx(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "请先登录", "data": nil})
		return
	}
	var body struct {
		Id  int `json:"id"`
		Idx int `json:"idx"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	var existing models.WorkoutAction
	if err := h.db.First(&existing, body.Id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "the record not found", "data": nil})
		return
	}
	now := time.Now().UTC()
	if err := h.db.Model(&existing).Updates(map[string]interface{}{
		"sort_idx":   body.Idx,
		"updated_at": now,
	}).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": body})
}

func (h *WorkoutActionHandler) DeleteWorkoutAction(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "请先登录", "data": nil})
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
	var record models.WorkoutAction
	if err := h.db.Where("id = ?", body.Id).First(&record).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "the record not found", "data": nil})
		return
	}
	if err := h.db.Model(&record).Update("d", 1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to delete the record", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "操作成功", "data": nil})
}

func (h *WorkoutActionHandler) CreateContentWithWorkoutAction(c *gin.Context) {
	var body struct {
		ContentId       int `json:"content_id"`
		WorkoutActionId int `json:"workout_action_id"`
		StartPoint      int `json:"start_point"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	if body.ContentId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少 Content id 参数", "data": nil})
		return
	}
	if body.WorkoutActionId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少 WorkoutAction id 参数", "data": nil})
		return
	}
	content_with_act := models.CoachContentWithWorkoutAction{
		SortIdx:         0,
		WorkoutActionId: body.WorkoutActionId,
		CoachContentId:  body.ContentId,
		CreatedAt:       time.Now(),
	}
	if err := h.db.Create(&content_with_act).Error; err != nil {
		h.logger.Error("Failed to create CoachContentWithWorkoutAction", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "关联成功", "data": nil})
}

func (h *WorkoutActionHandler) FetchContentListOfWorkoutAction(c *gin.Context) {
	// uid := int(c.GetFloat64("id"))

	var body struct {
		models.Pagination
		WorkoutActionId int `json:"workout_action_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.WorkoutActionId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少 WorkoutActionId 参数", "data": nil})
		return
	}
	query := h.db.Where("d IS NULL OR d = 0")
	query = query.Where("workout_action_id = ?", body.WorkoutActionId)
	pb := pagination.NewPaginationBuilder[models.CoachContentWithWorkoutAction](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC").
		SetOrderBy("sort_idx DESC")

	var list1 []models.CoachContentWithWorkoutAction
	if err := pb.Build().Preload("Content").Preload("Content.Coach").Preload("Content.Coach.Profile1").Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)

	list := make([]map[string]interface{}, 0, len(list2))
	for _, v := range list2 {
		list = append(list, map[string]interface{}{
			"id":          v.Content.Id,
			"title":       v.Content.Title,
			"description": v.Content.Description,
			"video_url":   v.Content.VideoKey,
			"like_count":  v.Content.LikeCount,
			"creator": map[string]interface{}{
				"nickname":   v.Content.Coach.Profile1.Nickname,
				"avatar_url": v.Content.Coach.Profile1.AvatarURL,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        list,
			"page":        pb.GetLimit(),
			"page_size":   body.PageSize,
			"has_more":    has_more,
			"next_marker": next_marker,
		},
	})
}
