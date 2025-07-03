package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myapi/internal/models"
	"myapi/internal/pkg/pagination"
	"myapi/pkg/logger"
)

// WorkoutPlanHandler handles HTTP requests for workout plans
type WorkoutPlanHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewWorkoutActionHandler creates a new workout action handler
func NewWorkoutPlanHandler(db *gorm.DB, logger *logger.Logger) *WorkoutPlanHandler {
	return &WorkoutPlanHandler{
		db:     db,
		logger: logger,
	}
}

type WorkoutPlanBody struct {
	Title             string `json:"title" binding:"required"`
	Type              string `json:"type"`
	Overview          string `json:"overview"`
	Level             int    `json:"level"`
	Status            int    `json:"status" binding:"required,oneof=1 2"`
	Details           string `json:"details" binding:"required"`
	Points            string `json:"points"`
	Suggestions       string `json:"suggestions"`
	EstimatedDuration int    `json:"estimated_duration"`
	Tags              string `json:"tags"`
	MuscleIds         string `json:"muscle_ids"`
	EquipmentIds      string `json:"equipment_ids"`
}

func (h *WorkoutPlanHandler) CreateWorkoutPlan(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		WorkoutPlanBody
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	record := models.WorkoutPlan{
		Title:             body.Title,
		Type:              body.Type,
		Overview:          body.Overview,
		Status:            body.Status,
		Level:             body.Level,
		EstimatedDuration: body.EstimatedDuration,
		Tags:              body.Tags,
		Details:           body.Details,
		Points:            body.Points,
		Suggestions:       body.Suggestions,
		MuscleIds:         body.MuscleIds,
		EquipmentIds:      body.EquipmentIds,
		CreatedAt:         time.Now(),
		OwnerId:           uid,
	}
	if record.Type == "" {
		record.Type = "strength"
	}
	if err := h.db.Create(&record).Error; err != nil {
		// h.logger.Error("Failed to create workout plan", "error", result.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": gin.H{"id": record.Id}})
}

func (h *WorkoutPlanHandler) UpdateWorkoutPlan(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		WorkoutPlanBody
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
	var existing models.WorkoutPlan
	if err := h.db.Where("id = ? AND owner_id = ?", body.Id, uid).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "没有找到记录", "data": nil})
		return
	}
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()

	existing.Title = body.Title
	existing.Overview = body.Overview
	existing.Level = body.Level
	existing.Tags = body.Tags
	existing.Details = body.Details
	// existing.Points = body.Points
	existing.Suggestions = body.Suggestions
	existing.EquipmentIds = body.EquipmentIds
	existing.MuscleIds = body.MuscleIds
	existing.EstimatedDuration = body.EstimatedDuration
	// existing.Status = body.Status
	now := time.Now().UTC()
	existing.UpdatedAt = &now

	if err := tx.Save(&existing).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update plan: " + err.Error(), "data": nil})
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to commit transaction", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "提交事务失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": gin.H{"id": existing.Id}})
}

func (h *WorkoutPlanHandler) FetchWorkoutPlanProfile(c *gin.Context) {
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
	var record models.WorkoutPlan
	if err := query.
		Preload("Creator.Profile1").
		First(&record).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有找到记录", "data": nil})
		return
	}
	if record.Status != int(models.WorkoutPublishStatusPublic) {
		if record.OwnerId != uid {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有权限查看", "data": nil})
			return
		}
	}
	data := map[string]interface{}{
		"id":                 record.Id,
		"title":              record.Title,
		"overview":           record.Overview,
		"level":              record.Level,
		"estimated_duration": record.EstimatedDuration,
		"suggestions":        record.Suggestions,
		"tags":               record.Tags,
		"cover_url":          record.CoverURL,
		"equipment_ids":      record.EquipmentIds,
		"muscle_ids":         record.MuscleIds,
		"details":            record.Details,
		"creator": map[string]interface{}{
			"nickname":   record.Creator.Profile1.Nickname,
			"avatar_url": record.Creator.Profile1.AvatarURL,
			"is_self":    record.Creator.Id == uid,
		},
		"created_at": record.CreatedAt,
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "请求成功", "data": data})
}

func (h *WorkoutPlanHandler) DeleteWorkoutPlan(c *gin.Context) {
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少 id 参数", "data": nil})
		return
	}
	var existing models.WorkoutPlan
	if err := h.db.Where("id = ?", body.Id).First(&existing).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if err := h.db.Model(&existing).Update("d", 1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功", "data": nil})
}

func (h *WorkoutPlanHandler) FetchWorkoutPlanList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		models.Pagination
		Ids     []int  `json:"ids"`
		Keyword string `json:"keyword"`
		Level   int    `json:"level"`
		Tag     string `json:"tag"`
		Muscle  string `json:"muscle"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	query := h.db.Where("d IS NULL OR d = 0")
	if uid != 0 {
		query = query.Where("(status = 1) OR (status = 2 AND owner_id = ?)", uid)
	} else {
		query = query.Where("status = 1")
	}
	if body.Keyword != "" {
		query = query.Where("title LIKE ?", "%"+body.Keyword+"%")
	}
	if body.Level != 0 {
		query = query.Where("level = ?", body.Level)
	}
	if body.Tag != "" {
		query = query.Where("tags LIKE ?", "%"+body.Tag+"%")
	}
	if body.Muscle != "" {
		query = query.Where("target_muscle_ids LIKE ?", "%"+body.Muscle+"%")
	}
	if len(body.Ids) != 0 {
		query = query.Where("id IN (?)", body.Ids)
	}
	pb := pagination.NewPaginationBuilder[models.WorkoutPlan](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")

	var list1 []models.WorkoutPlan
	if err := pb.Build().Preload("Creator.Profile1").Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
	list := make([]map[string]interface{}, 0, len(list2))
	for _, v := range list2 {
		list = append(list, map[string]interface{}{
			"id":                 v.Id,
			"title":              v.Title,
			"overview":           v.Overview,
			"level":              v.Level,
			"tags":               v.Tags,
			"estimated_duration": v.EstimatedDuration,
			"creator": map[string]interface{}{
				"nickname":   v.Creator.Profile1.Nickname,
				"avatar_url": v.Creator.Profile1.AvatarURL,
			},
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

func (h *WorkoutPlanHandler) FetchMyWorkoutPlanList(c *gin.Context) {
	var plans []models.WorkoutPlan

	id := int(c.GetFloat64("id"))

	// Get query parameters for filtering
	level := c.Query("level")
	tag := c.Query("tag")
	muscle := c.Query("muscle")
	cursor := c.Query("cursor")
	pageSize := c.Query("page_size")

	// Get pagination parameters
	limit := 10 // default page size

	if pageSize != "" {
		if parsedSize, err := strconv.Atoi(pageSize); err == nil && parsedSize > 0 {
			limit = parsedSize
		}
	}

	// Start with base query
	query := h.db

	query = query.Where("owner_id = ?", id)

	// Apply cursor pagination if cursor is provided
	if cursor != "" {
		if cursorID, err := strconv.Atoi(cursor); err == nil {
			query = query.Where("id > ?", cursorID)
		}
	}

	// Apply filters if provided
	if level != "" {
		levelInt, err := strconv.Atoi(level)
		if err == nil {
			query = query.Where("level = ?", levelInt)
		}
	}

	if tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}

	if muscle != "" {
		query = query.Where("target_muscle_ids LIKE ?", "%"+muscle+"%")
	}

	// Add ordering and limit
	// query = query.Order("id asc").Limit(limit + 1)
	query = query.Order("created_at desc").Limit(limit + 1)

	// Execute the query
	result := query.Find(&plans)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout plans", "data": nil})
		return
	}

	hasMore := false
	nextCursor := ""

	// Check if there are more results
	if len(plans) > limit {
		hasMore = true
		plans = plans[:limit]                             // Remove the extra item we fetched
		nextCursor = strconv.Itoa(int(plans[limit-1].Id)) // Get the last item's ID as next cursor
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        plans,
			"page_size":   limit,
			"has_more":    hasMore,
			"next_marker": nextCursor,
		},
	})
}

func (h *WorkoutPlanHandler) FetchContentListOfWorkoutPlan(c *gin.Context) {
	// uid := int(c.GetFloat64("id"))

	var body struct {
		models.Pagination
		WorkoutPlanId int    `json:"workout_plan_id"`
		Level         int    `json:"level"`
		Tag           string `json:"tag"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if body.WorkoutPlanId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "缺少 WorkoutPlanId 参数", "data": nil})
		return
	}
	query := h.db.Where("d IS NULL OR d = 0")
	query = query.Where("workout_plan_id = ?", body.WorkoutPlanId)
	if body.Level != 0 {
		query = query.Where("level = ?", body.Level)
	}
	if body.Tag != "" {
		query = query.Where("tags LIKE ?", "%"+body.Tag+"%")
	}
	pb := pagination.NewPaginationBuilder[models.CoachContentWithWorkoutPlan](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")

	var list1 []models.CoachContentWithWorkoutPlan
	if err := pb.Build().Preload("Content").Preload("Content.Coach.Profile1").Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
	list := make([]map[string]interface{}, 0, len(list2))
	for _, v := range list2 {
		list = append(list, map[string]interface{}{
			"id":        v.Id,
			"title":     v.Content.Title,
			"overview":  v.Content.Description,
			"video_key": v.Content.VideoKey,
			"details":   v.Details,
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
			"page_size":   pb.GetLimit(),
			"has_more":    has_more,
			"next_marker": next_marker,
		},
	})
}

func (h *WorkoutPlanHandler) CreateContentWithWorkoutPlan(c *gin.Context) {
	var body struct {
		Details       string `json:"details"`
		ContentId     int    `json:"content_id"`
		WorkoutPlanId int    `json:"workout_plan_id"`
		// Points        []struct {
		// 	Time              int    `json:"time"`
		// 	TimeText          string `json:"time_text"`
		// 	WorkoutActionId   int    `json:"workout_action_id"`
		// 	WorkoutActionName string `json:"workout_action_name"`
		// } `json:"points"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	if body.ContentId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少 Content id 参数", "data": nil})
		return
	}
	if body.WorkoutPlanId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少 WorkoutPlan id 参数", "data": nil})
		return
	}
	content_with_plan := models.CoachContentWithWorkoutPlan{
		SortIdx:        0,
		Details:        body.Details,
		CreatedAt:      time.Now(),
		CoachContentId: body.ContentId,
		WorkoutPlanId:  body.WorkoutPlanId,
	}
	if err := h.db.Create(&content_with_plan).Error; err != nil {
		h.logger.Error("Failed to create CoachContentWithWorkoutPlan", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "关联成功", "data": nil})
}

func (h *WorkoutPlanHandler) FetchContentProfileOfWorkoutPlan(c *gin.Context) {
	// uid := int(c.GetFloat64("id"))

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
	var record models.CoachContentWithWorkoutPlan
	if err := h.db.Where("id = ?", body.Id).Preload("Content").Preload("Content.Coach.Profile1").First(&record).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有找到记录", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"id":           record.Id,
			"content_type": record.Content.ContentType,
			"title":        record.Content.Title,
			"overview":     record.Content.Description,
			"details":      record.Details,
			"video_url":    record.Content.VideoKey,
			"creator": map[string]interface{}{
				"nickname":   record.Content.Coach.Profile1.Nickname,
				"avatar_url": record.Content.Coach.Profile1.AvatarURL,
			},
		},
	})
}

func (h *WorkoutPlanHandler) CreateWorkoutSchedule(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Title    string `json:"title" binding:"required"`
		Overview string `json:"overview"`
		Status   int    `json:"status" binding:"required,oneof=1 2"`
		Level    int    `json:"level"`
		Type     int    `json:"type"`
		// Details []struct {
		// 	Weekday        int   `json:"weekday"`
		// 	Day            int   `json:"day"`
		// 	Idx            int   `json:"idx"`
		// 	WorkoutPlanIds []int `json:"workout_plan_ids"`
		// } `json:"schedules"`
		Details string `json:"details"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Title == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少标题", "data": nil})
		return
	}
	// if len(body.Schedules) == 0 {
	// 	c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "至少选择一天配置训练计划", "data": nil})
	// 	return
	// }
	// 开始事务
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	// 创建训练计划集合
	record := models.WorkoutSchedule{
		Title:     body.Title,
		Overview:  body.Overview,
		Status:    body.Status,
		Level:     body.Level,
		Type:      body.Type,
		Details:   body.Details,
		CreatedAt: time.Now().UTC(),
		OwnerId:   uid,
	}
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "创建训练计划集合失败: " + err.Error(), "data": nil})
		return
	}
	// 创建训练计划集合中的训练计划
	// for _, plan := range body.Schedules {
	// 	for _, vv := range plan.WorkoutPlanIds {
	// 		plan_in_collection := models.WorkoutPlanInSchedule{
	// 			WorkoutPlanCollectionId: record.Id,
	// 			WorkoutPlanId:           vv,
	// 			Weekday:                 plan.Weekday,
	// 			Day:                     plan.Day,
	// 			Idx:                     plan.Idx,
	// 		}
	// 		if err := tx.Create(&plan_in_collection).Error; err != nil {
	// 			tx.Rollback()
	// 			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
	// 			return
	// 		}
	// 	}
	// }
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to commit transaction", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": gin.H{"id": record.Id}})
}

func (h *WorkoutPlanHandler) UpdateWorkoutSchedule(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Id           int    `json:"id"`
		Title        string `json:"title"`
		Overview     string `json:"overview"`
		Level        int    `json:"level"`
		Type         int    `json:"type"`
		WorkoutPlans []struct {
			Weekday        int   `json:"weekday"`
			Day            int   `json:"day"`
			WorkoutPlanIds []int `json:"workout_plan_ids"`
		} `json:"workout_plans"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少 id 参数", "data": nil})
		return
	}

	var existing models.WorkoutSchedule
	if err := h.db.Where("id = ? AND owner_id = ?", body.Id, uid).First(&existing).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if body.Title != "" {
		updates["title"] = body.Title
	}
	if body.Overview != "" {
		updates["overview"] = body.Overview
	}
	if body.Level != 0 {
		updates["level"] = body.Level
	}
	if body.Type != 0 {
		updates["type"] = body.Type
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	if err := tx.Model(&existing).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "更新周期计划失败，" + err.Error(), "data": nil})
		return
	}

	if err := tx.Where("workout_plan_collection_id = ?", body.Id).Delete(&models.WorkoutPlanInSchedule{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "删除现有训练计划关联失败，" + err.Error(), "data": nil})
		return
	}
	// Create new workout plans in collection
	for _, plan := range body.WorkoutPlans {
		for _, vv := range plan.WorkoutPlanIds {
			plan_in_collection := models.WorkoutPlanInSchedule{
				WorkoutPlanCollectionId: body.Id,
				WorkoutPlanId:           vv,
				Weekday:                 plan.Weekday,
				Day:                     plan.Day,
			}
			if err := tx.Create(&plan_in_collection).Error; err != nil {
				tx.Rollback()
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
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": nil})
}

func (h *WorkoutPlanHandler) FetchWorkoutScheduleProfile(c *gin.Context) {
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
	var record models.WorkoutSchedule
	query := h.db.Where("d IS NULL OR d = 0")
	query = query.Where("id = ?", body.Id)
	if err := query.
		Preload("WorkoutPlans").
		Preload("WorkoutPlans.WorkoutPlan").
		Preload("Creator.Profile1").
		First(&record).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "找不到记录", "data": nil})
		return
	}
	if record.Status != int(models.WorkoutPublishStatusPublic) {
		if record.OwnerId != uid {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有权限查看", "data": nil})
			return
		}
	}
	data := map[string]interface{}{
		"id":       record.Id,
		"title":    record.Title,
		"overview": record.Overview,
		"level":    record.Level,
		"type":     record.Type,
		"details":  record.Details,
		"creator": map[string]interface{}{
			"nickname":   record.Creator.Profile1.Nickname,
			"avatar_url": record.Creator.Profile1.AvatarURL,
			"is_self":    record.Creator.Id == uid,
		},
		"created_at": record.CreatedAt,
	}
	schedules := []interface{}{}
	for _, schedule := range record.WorkoutPlans {
		schedules = append(schedules, map[string]interface{}{
			"idx":     schedule.Idx,
			"day":     schedule.Day,
			"weekday": schedule.Weekday,
			"workout_plan": gin.H{
				"id":                 schedule.WorkoutPlan.Id,
				"title":              schedule.WorkoutPlan.Title,
				"overview":           schedule.WorkoutPlan.Overview,
				"tags":               schedule.WorkoutPlan.Tags,
				"estimated_duration": schedule.WorkoutPlan.EstimatedDuration,
			},
		})
	}
	data["schedules"] = schedules
	if uid != 0 {
		var record2 models.CoachWorkoutSchedule
		h.db.Where("coach_id = ? AND workout_plan_collection_id = ?", uid, body.Id).First(&record2)
		if record2.Id != 0 {
			data["applied"] = record2.Status
			data["applied_in_interval"] = record2.Interval
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": data})
}

func (h *WorkoutPlanHandler) FetchWorkoutScheduleList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		models.Pagination
		Keyword string `json:"keyword"`
		Level   int    `json:"level"`
		Tag     string `json:"tag"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	query := h.db.Where("d IS NULL OR d = 0")
	if uid != 0 {
		query = query.Where("(status = 1) OR (status = 2 AND owner_id = ?)", uid)
	} else {
		query = query.Where("status = 1")
	}
	if body.Level != 0 {
		query = query.Where("level = ?", body.Level)
	}
	if body.Keyword != "" {
		query = query.Where("title LIKE ?", "%"+body.Keyword+"%")
	}
	if body.Tag != "" {
		query = query.Where("tags LIKE ?", "%"+body.Tag+"%")
	}
	pb := pagination.NewPaginationBuilder[models.WorkoutSchedule](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")

	var list1 []models.WorkoutSchedule
	if err := pb.Build().Preload("Creator.Profile1").Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
	list := make([]map[string]interface{}, 0, len(list2))
	for _, v := range list2 {
		list = append(list, map[string]interface{}{
			"id":       v.Id,
			"type":     v.Type,
			"title":    v.Title,
			"overview": v.Overview,
			"level":    v.Level,
			"creator": map[string]interface{}{
				"nickname":   v.Creator.Profile1.Nickname,
				"avatar_url": v.Creator.Profile1.AvatarURL,
			},
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

// 应用某个周期计划
func (h *WorkoutPlanHandler) ApplyWorkoutSchedule(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Id       int `json:"id"`
		Interval int `json:"interval"`
		// 只有 天循环 会需要？
		StartDate time.Time `json:"start_date"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "缺少周期计划 id", "data": nil})
		return
	}

	var existing models.CoachWorkoutSchedule
	if err := h.db.Where("coach_id = ? AND workout_plan_collection_id = ?", uid, body.Id).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		// 创建
		record := models.CoachWorkoutSchedule{
			WorkoutPlanCollectionId: body.Id,
			CoachId:                 uid,
			StartDate:               &body.StartDate,
			Status:                  1,
			AppliedAt:               time.Now(),
		}
		if err := h.db.Create(&record).Error; err != nil {
			h.logger.Error("Failed to create record", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create record", "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "应用周期计划成功", "data": nil})
		return
	}

	if existing.Status == 1 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "该周期计划已应用", "data": nil})
		return
	}
	// 更新
	updates := map[string]interface{}{
		"status":     1,
		"start_date": body.StartDate,
		"applied_at": time.Now(),
	}
	if err := h.db.Model(&existing).Updates(updates).Error; err != nil {
		h.logger.Error("Failed to update", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "应用失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "应用周期计划成功", "data": nil})

}

// 取消应用某个周期计划
func (h *WorkoutPlanHandler) CancelWorkoutSchedule(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "缺少周期计划 id", "data": nil})
		return
	}

	var existing models.CoachWorkoutSchedule
	if err := h.db.Where("coach_id = ? AND workout_plan_collection_id = ?", uid, body.Id).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "周期计划不存在", "data": nil})
		return
	}
	if existing.CoachId != uid {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Workout plan not found", "data": nil})
		return
	}

	// 更新
	updates := map[string]interface{}{
		"status":       2,
		"cancelled_at": time.Now(),
	}
	if err := h.db.Model(&existing).Updates(updates).Error; err != nil {
		h.logger.Error("Failed to update", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "取消失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "取消周期计划成功", "data": nil})
}

// 获取当前应用中的周期计划
func (h *WorkoutPlanHandler) FetchAppliedWorkoutScheduleList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var list []models.CoachWorkoutSchedule
	// @todo 周期计划，历史数据处理完了后，这里就改成 Preload("WorkoutPlanCollection")
	if err := h.db.Where("status = 1 AND coach_id = ?", uid).Preload("WorkoutPlanCollection.WorkoutPlans.WorkoutPlan").Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	data := []interface{}{}
	for _, relation := range list {
		schedules := []interface{}{}
		for _, schedule := range relation.WorkoutPlanCollection.WorkoutPlans {
			schedules = append(schedules, map[string]interface{}{
				"day":             schedule.Day,
				"weekday":         schedule.Weekday,
				"workout_plan_id": schedule.WorkoutPlanId,
				"title":           schedule.WorkoutPlan.Title,
				"overview":        schedule.WorkoutPlan.Overview,
				"tags":            schedule.WorkoutPlan.Tags,
			})
		}
		d := map[string]interface{}{
			"id":                  relation.Id,
			"status":              relation.Status,
			"workout_schedule_id": relation.WorkoutPlanCollection.Id,
			"title":               relation.WorkoutPlanCollection.Title,
			"overview":            relation.WorkoutPlanCollection.Overview,
			"type":                relation.WorkoutPlanCollection.Type,
			"level":               relation.WorkoutPlanCollection.Level,
			"details":             relation.WorkoutPlanCollection.Details,
			"schedules":           schedules,
			"applied_at":          relation.AppliedAt,
			"start_date":          nil,
		}
		if relation.StartDate != nil {
			d["start_date"] = relation.StartDate
		}
		data = append(data, d)
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": gin.H{
		"list": data,
	}})
}

func (h *WorkoutPlanHandler) FetchWorkoutPlanSetList(c *gin.Context) {
	var body struct {
		models.Pagination
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	query := h.db.Where("d IS NULL OR d = 0")
	pb := pagination.NewPaginationBuilder[models.WorkoutPlanSet](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("idx DESC")
	var list1 []models.WorkoutPlanSet
	if err := pb.Build().Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch workout plans", "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
	list := make([]map[string]interface{}, 0, len(list2))
	type WorkoutPlanSetJSON250607 struct {
		Type     int    `json:"type"`
		Id       int    `json:"id"`
		Title    string `json:"title"`
		Overview string `json:"overview"`
		Tags     string `json:"tags"`
		Creator  struct {
			Nickname  string `json:"nickname"`
			AvatarURL string `json:"avatar_url"`
		} `json:"creator"`
	}
	for _, v := range list2 {
		var details []WorkoutPlanSetJSON250607
		if err := json.Unmarshal([]byte(v.Details), &details); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 600, "msg": err.Error(), "data": nil})
			return
		}
		list = append(list, map[string]interface{}{
			"id":      v.Id,
			"title":   v.Title,
			"idx":     v.Idx,
			"details": details,
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

func (h *WorkoutPlanHandler) CreateWorkoutPlanSet(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "请先登录", "data": nil})
		return
	}
	var body struct {
		Title    string `json:"title"`
		Overview string `json:"overview"`
		IconURL  string `json:"icon_url"`
		Idx      int    `json:"idx"`
		Details  []struct {
			Type int `json:"type"`
			Id   int `json:"id"`
		} `json:"details"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Title == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少标题", "data": nil})
		return
	}
	if len(body.Details) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少内容", "data": nil})
		return
	}
	var details []map[string]interface{}
	for _, v := range body.Details {
		if v.Type == 1 {
			var existing models.WorkoutPlan
			if err := h.db.Where("id = ?", v.Id).Preload("Creator.Profile1").First(&existing).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "数据异常", "data": nil})
				return
			}
			details = append(details, map[string]interface{}{
				"type":     v.Type,
				"id":       v.Id,
				"title":    existing.Title,
				"overview": existing.Overview,
				"tags":     existing.Tags,
				"creator": map[string]interface{}{
					"nickname":   existing.Creator.Profile1.Nickname,
					"avatar_url": existing.Creator.Profile1.AvatarURL,
				},
			})
		}
		if v.Type == 2 {
			var existing models.WorkoutSchedule
			if err := h.db.Where("id = ?", v.Id).Preload("Creator.Profile1").First(&existing).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "数据异常", "data": nil})
				return
			}
			details = append(details, map[string]interface{}{
				"type":     v.Type,
				"id":       v.Id,
				"title":    existing.Title,
				"overview": existing.Overview,
				"tags":     "",
				"creator": map[string]interface{}{
					"nickname":   existing.Creator.Profile1.Nickname,
					"avatar_url": existing.Creator.Profile1.AvatarURL,
				},
			})
		}
	}

	detailsJSON, err := json.Marshal(details)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to serialize details: " + err.Error(), "data": nil})
		return
	}
	record := models.WorkoutPlanSet{
		Title:     body.Title,
		Overview:  body.Overview,
		IconURL:   body.IconURL,
		Idx:       body.Idx,
		Details:   string(detailsJSON),
		CreatedAt: time.Now().UTC(),
	}

	if err := h.db.Create(&record).Error; err != nil {
		// h.logger.Error("Failed to create workout plan", "error", result.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": gin.H{"id": record.Id}})
}

func (h *WorkoutPlanHandler) UpdateWorkoutPlanSet(c *gin.Context) {
	// 绑定更新的数据
	var body models.WorkoutPlanSet
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	// 先获取现有的计划
	var existing_plan models.WorkoutPlanSet
	if result := h.db.First(&existing_plan, body.Id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout plan not found", "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
		return
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()

	// 更新基本字段
	existing_plan.Title = body.Title
	existing_plan.Overview = body.Overview
	existing_plan.Idx = body.Idx
	existing_plan.Details = body.Details

	// 保存主记录
	if err := tx.Save(&existing_plan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update plan: " + err.Error(), "data": nil})
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": gin.H{"id": existing_plan.Id}})
}
