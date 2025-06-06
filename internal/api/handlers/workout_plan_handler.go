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

func (h *WorkoutPlanHandler) CreateWorkoutPlan(c *gin.Context) {
	id := int(c.GetFloat64("id"))

	var plan models.WorkoutPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	// Validate required fields
	if plan.Title == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Title and overview are required", "data": nil})
		return
	}

	// 使用 UTC 时间
	plan.CreatedAt = time.Now().UTC()
	plan.OwnerId = id

	// Insert into database
	result := h.db.Create(&plan)

	if result.Error != nil {
		// h.logger.Error("Failed to create workout plan", "error", result.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan created successfully", "data": gin.H{"id": plan.Id}})
}

func (h *WorkoutPlanHandler) FetchWorkoutPlanProfile(c *gin.Context) {
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	var plan models.WorkoutPlan
	// 使用 Preload 预加载关联的 Steps 和 Actions
	result := h.db.
		Where("id = ?", body.Id).
		First(&plan)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout plan not found", "data": nil})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan retrieved successfully", "data": plan})
}

func (h *WorkoutPlanHandler) UpdateWorkoutPlan(c *gin.Context) {
	// 绑定更新的数据
	var body models.WorkoutPlan
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	// 先获取现有的计划
	var existing_plan models.WorkoutPlan
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
		}
	}()

	// 更新基本字段
	existing_plan.Title = body.Title
	existing_plan.Overview = body.Overview
	existing_plan.Level = body.Level
	existing_plan.Tags = body.Tags
	existing_plan.Details = body.Details
	existing_plan.Points = body.Points
	existing_plan.Suggestions = body.Suggestions
	existing_plan.EquipmentIds = body.EquipmentIds
	existing_plan.MuscleIds = body.MuscleIds
	existing_plan.EstimatedDuration = body.EstimatedDuration
	existing_plan.Status = body.Status
	now := time.Now().UTC()
	existing_plan.UpdatedAt = &now

	// 保存主记录
	if err := tx.Save(&existing_plan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update plan: " + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan updated successfully", "data": gin.H{"id": existing_plan.Id}})
}

func (h *WorkoutPlanHandler) DeleteWorkoutPlan(c *gin.Context) {
	var id struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&id); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	result := h.db.Where("id = ?", id.Id).Delete(&models.WorkoutPlan{})
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to delete workout plan", "data": nil})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout plan not found", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan deleted successfully", "data": nil})
}

func (h *WorkoutPlanHandler) FetchWorkoutPlanList(c *gin.Context) {
	var plans []models.WorkoutPlan

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

	has_more := false
	next_cursor := ""

	// Check if there are more results
	if len(plans) > limit {
		has_more = true
		plans = plans[:limit]                              // Remove the extra item we fetched
		next_cursor = strconv.Itoa(int(plans[limit-1].Id)) // Get the last item's ID as next cursor
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        plans,
			"page_size":   limit,
			"has_more":    has_more,
			"next_marker": next_cursor,
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

func (h *WorkoutPlanHandler) CreateWorkoutSchedule(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Title        string `json:"title"`
		Overview     string `json:"overview"`
		Level        int    `json:"level"`
		Type         int    `json:"type"`
		WorkoutPlans []struct {
			Weekday       int `json:"weekday"`
			Day           int `json:"day"`
			WorkoutPlanId int `json:"workout_plan_id"`
		} `json:"workout_plans"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Title == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少标题", "data": nil})
		return
	}
	if len(body.WorkoutPlans) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "至少选择一天配置训练计划", "data": nil})
		return
	}

	// 开始事务
	tx := h.db.Begin()
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		tx.Rollback()
	// 	}
	// }()
	// 创建训练计划集合
	record := models.WorkoutSchedule{
		Title:     body.Title,
		Overview:  body.Overview,
		Level:     body.Level,
		Type:      body.Type,
		CreatedAt: time.Now().UTC(),
		OwnerId:   uid,
	}
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "创建训练计划集合失败: " + err.Error(), "data": nil})
		return
	}
	// 创建训练计划集合中的训练计划
	for _, plan := range body.WorkoutPlans {
		plan_in_collection := models.WorkoutPlanInSchedule{
			WorkoutPlanCollectionId: record.Id,
			WorkoutPlanId:           plan.WorkoutPlanId,
			Weekday:                 plan.Weekday,
			Day:                     plan.Day,
		}
		if err := tx.Create(&plan_in_collection).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "创建训练计划关联失败: " + err.Error(), "data": nil})
			return
		}
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to commit transaction", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "提交事务失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建训练计划集合成功", "data": gin.H{"id": record.Id}})
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
			Weekday       int `json:"weekday"`
			Day           int `json:"day"`
			WorkoutPlanId int `json:"workout_plan_id"`
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
		plan_in_collection := models.WorkoutPlanInSchedule{
			WorkoutPlanCollectionId: body.Id,
			WorkoutPlanId:           plan.WorkoutPlanId,
			Weekday:                 plan.Weekday,
			Day:                     plan.Day,
		}
		if err := tx.Create(&plan_in_collection).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "创建训练计划关联失败，" + err.Error(), "data": nil})
			return
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "提交事务失败，" + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新周期计划成功", "data": nil})
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

	var record models.WorkoutSchedule

	result := h.db.
		Preload("WorkoutPlans").
		Preload("WorkoutPlans.WorkoutPlan").
		Where("id = ?", body.Id).
		First(&record)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Workout plan not found", "data": nil})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
		return
	}

	var record2 models.CoachWorkoutSchedule
	if err := h.db.Where("coach_id = ? AND workout_plan_collection_id = ?", uid, body.Id).First(&record2).Error; err != nil {
	}

	data := map[string]interface{}{
		"title":    record.Title,
		"overview": record.Overview,
		"level":    record.Level,
		"type":     record.Type,
	}
	schedules := []interface{}{}
	for _, schedule := range record.WorkoutPlans {
		schedules = append(schedules, map[string]interface{}{
			"day":          schedule.Day,
			"weekday":      schedule.Weekday,
			"workout_plan": schedule.WorkoutPlan,
		})
	}
	data["schedules"] = schedules
	if record2.Id != 0 {
		data["applied"] = record2.Status
		data["applied_in_interval"] = record2.Interval
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan retrieved successfully", "data": data})
}

func (h *WorkoutPlanHandler) FetchWorkoutScheduleList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		models.Pagination
		Level int    `json:"level"`
		Tag   string `json:"tag"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	// Get pagination parameters
	limit := 10 // default page size

	if body.PageSize != 0 {
		limit = body.PageSize
	}

	// Start with base query
	query := h.db
	query = query.Where("owner_id = ?", uid)
	// Apply filters if provided
	if body.Level != 0 {
		query = query.Where("level = ?", body.Level)
	}
	if body.Tag != "" {
		query = query.Where("tags LIKE ?", "%"+body.Tag+"%")
	}
	// query = query.Order("id asc").Limit(limit + 1)
	query = query.Order("created_at desc").Limit(limit + 1)

	// Execute the query
	var list []models.WorkoutSchedule
	result := query.Find(&list)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
		return
	}

	has_more := false
	next_cursor := ""

	// Check if there are more results
	if len(list) > limit {
		has_more = true
		list = list[:limit]                               // Remove the extra item we fetched
		next_cursor = strconv.Itoa(int(list[limit-1].Id)) // Get the last item's ID as next cursor
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        list,
			"page_size":   limit,
			"has_more":    has_more,
			"next_marker": next_cursor,
		},
	})
}

// 应用某个周期计划
func (h *WorkoutPlanHandler) ApplyWorkoutSchedule(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Id       int `json:"id"`
		Interval int `json:"interval"`
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

func (h *WorkoutPlanHandler) FetchMyWorkoutScheduleList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var list []models.CoachWorkoutSchedule

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
				"tags":            schedule.WorkoutPlan.Tags,
			})
		}

		data = append(data, map[string]interface{}{
			"status":              relation.Status,
			"workout_schedule_id": relation.WorkoutPlanCollection.Id,
			"title":               relation.WorkoutPlanCollection.Title,
			"overview":            relation.WorkoutPlanCollection.Overview,
			"type":                relation.WorkoutPlanCollection.Type,
			"level":               relation.WorkoutPlanCollection.Level,
			"schedules":           schedules,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": gin.H{
		"list": data,
	}})
}

// 训练计划合集，用来在前端按组展示
func (h *WorkoutPlanHandler) FetchWorkoutPlanSetList(c *gin.Context) {
	var plans []models.WorkoutPlanSet

	cursor := c.Query("cursor")
	pageSize := c.Query("page_size")

	// Get pagination parameters
	limit := 100 // default page size

	if pageSize != "" {
		if parsedSize, err := strconv.Atoi(pageSize); err == nil && parsedSize > 0 {
			limit = parsedSize
		}
	}

	// Start with base query
	query := h.db

	// Apply cursor pagination if cursor is provided
	if cursor != "" {
		if cursorID, err := strconv.Atoi(cursor); err == nil {
			query = query.Where("id > ?", cursorID)
		}
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

	has_more := false
	next_cursor := ""

	// Check if there are more results
	if len(plans) > limit {
		has_more = true
		plans = plans[:limit]                              // Remove the extra item we fetched
		next_cursor = strconv.Itoa(int(plans[limit-1].Id)) // Get the last item's ID as next cursor
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        plans,
			"page_size":   limit,
			"has_more":    has_more,
			"next_marker": next_cursor,
		},
	})
}

func (h *WorkoutPlanHandler) CreateWorkoutPlanSet(c *gin.Context) {
	// id := int(c.GetFloat64("id"))

	var body models.WorkoutPlanSet
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	// Validate required fields
	if body.Title == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Title and overview are required", "data": nil})
		return
	}

	// 使用 UTC 时间
	body.CreatedAt = time.Now().UTC()

	// Insert into database
	result := h.db.Create(&body)

	if result.Error != nil {
		// h.logger.Error("Failed to create workout plan", "error", result.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan created successfully", "data": gin.H{"id": body.Id}})
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

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Workout plan updated successfully", "data": gin.H{"id": existing_plan.Id}})
}
