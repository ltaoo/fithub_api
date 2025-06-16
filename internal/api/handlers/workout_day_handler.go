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
type WorkoutDayHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewWorkoutDayHandler creates a new workout day handler
func NewWorkoutDayHandler(db *gorm.DB, logger *logger.Logger) *WorkoutDayHandler {
	return &WorkoutDayHandler{
		db:     db,
		logger: logger,
	}
}

func (h *WorkoutDayHandler) CreateWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	if uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "非法操作", "data": nil})
		return
	}
	var body struct {
		StudentIds      []int `json:"student_ids"`
		StartWhenCreate bool  `json:"start_when_create"`
		WorkoutPlanId   int   `json:"workout_plan_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.WorkoutPlanId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少训练计划", "data": nil})
		return
	}
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()
	var subscription models.Subscription
	if err := tx.Where("coach_id = ? AND step = 2", uid).First(&subscription).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 101, "msg": "该功能需订阅后才能使用", "data": nil})
		return
	}
	now := time.Now().UTC()
	group_no := strconv.FormatInt(now.Unix(), 10)
	var workout_day_ids []int
	if len(body.StudentIds) != 0 {
		for _, student_id := range body.StudentIds {
			is_coach_self := student_id == uid || student_id == 0
			if is_coach_self {
				workout_day := models.WorkoutDay{
					Status: int(models.WorkoutDayStatusPending),
					// 只有选了多人才有该字段
					GroupNo:       group_no,
					CreatedAt:     now,
					WorkoutPlanId: body.WorkoutPlanId,
					CoachId:       uid,
					StudentId:     uid,
				}
				if body.StartWhenCreate {
					workout_day.StartedAt = &now
					workout_day.Status = int(models.WorkoutDayStatusStarted)
				}
				if err := tx.Create(&workout_day).Error; err != nil {
					h.logger.Error("Failed to create workout plan", err)
					tx.Rollback()
					c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
					return
				}
				workout_day_ids = append(workout_day_ids, workout_day.Id)
				continue
			}
			var relation models.CoachRelationship
			if err := tx.Where("coach_id = ? AND student_id = ?", uid, student_id).First(&relation).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
					return
				}
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "非法操作", "data": nil})
				return
			}
			workout_day := models.WorkoutDay{
				Status: int(models.WorkoutDayStatusPending),
				// 只有选了多人才有该字段
				GroupNo:       group_no,
				CreatedAt:     now,
				WorkoutPlanId: body.WorkoutPlanId,
				CoachId:       relation.CoachId,
				StudentId:     relation.StudentId,
			}
			if body.StartWhenCreate {
				workout_day.StartedAt = &now
				workout_day.Status = int(models.WorkoutDayStatusStarted)
			}
			if err := tx.Create(&workout_day).Error; err != nil {
				h.logger.Error("Failed to create workout plan", err)
				tx.Rollback()
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
				return
			}
			workout_day_ids = append(workout_day_ids, workout_day.Id)
		}
	} else {
		workout_day := models.WorkoutDay{
			Status:        int(models.WorkoutDayStatusPending),
			CreatedAt:     now,
			WorkoutPlanId: body.WorkoutPlanId,
			CoachId:       uid,
			StudentId:     uid,
		}
		if body.StartWhenCreate {
			workout_day.StartedAt = &now
			workout_day.Status = int(models.WorkoutDayStatusStarted)
		}
		if err := tx.Create(&workout_day).Error; err != nil {
			h.logger.Error("Failed to create workout plan", err)
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		workout_day_ids = append(workout_day_ids, workout_day.Id)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to commit transaction", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to commit transaction", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": gin.H{"ids": workout_day_ids}})
}

// 查看是否有进行中的训练，仅获取少量数据
func (h *WorkoutDayHandler) CheckHasStartedWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var list []models.WorkoutDay
	if err := h.db.
		Where("status = ?", int(models.WorkoutDayStatusStarted)).
		Where("coach_id = ?", uid).
		Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	data := []map[string]interface{}{}

	for _, v := range list {
		data = append(data, map[string]interface{}{
			"id":         v.Id,
			"status":     v.Status,
			"created_at": v.CreatedAt,
			"started_at": v.StartedAt,
			"student_id": v.StudentId,
			"coach_id":   v.CoachId,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "请求成功", "data": gin.H{
		"list": data,
	}})
}

func (h *WorkoutDayHandler) FetchStartedWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var list []models.WorkoutDay
	if err := h.db.
		Where("status = ?", int(models.WorkoutDayStatusStarted)).
		Where("coach_id = ?", uid).
		Where("started_at IS NOT NULL").
		Order("started_at DESC").
		Preload("WorkoutPlan").
		Preload("Student.Profile1").
		Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	data := []map[string]interface{}{}

	for _, v := range list {
		data = append(data, map[string]interface{}{
			"id":         v.Id,
			"status":     v.Status,
			"created_at": v.CreatedAt,
			"started_at": v.StartedAt,
			"coach_id":   v.CoachId,
			"group_no":   v.GroupNo,
			"student": map[string]interface{}{
				"id":         v.StudentId,
				"nickname":   v.Student.Profile1.Nickname,
				"avatar_url": v.Student.Profile1.AvatarURL,
			},
			"workout_plan": map[string]interface{}{
				"id":       v.WorkoutPlan.Id,
				"title":    v.WorkoutPlan.Title,
				"overview": v.WorkoutPlan.Overview,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "请求成功", "data": gin.H{
		"list": data,
	}})
}

// 获取训练计划记录 只能获取自己和自己学员的
func (h *WorkoutDayHandler) FetchWorkoutDayProfile(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
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
	var workout_day models.WorkoutDay
	if err := h.db.
		Where("id = ? AND coach_id = ?", body.Id, uid).
		Preload("WorkoutPlan").
		Preload("WorkoutPlan.Creator.Profile1").
		First(&workout_day).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "没有找到该记录", "data": nil})
		return
	}

	// Calculate the day number based on unique dates
	// var day_number int64
	// h.db.Model(&models.WorkoutDay{}).
	// 	Select("COUNT(DISTINCT DATE(created_at))").
	// 	Where("student_id = ? AND DATE(created_at) <= DATE(?)", uid, workout_day.CreatedAt).
	// 	Count(&day_number)

	data := gin.H{
		"id":     workout_day.Id,
		"status": workout_day.Status,
		// 训练记录
		"pending_steps": workout_day.PendingSteps,
		// 训练内容
		"updated_details": workout_day.UpdatedDetails,
		"student_id":      workout_day.StudentId,
		"is_self":         workout_day.StudentId == uid,
		"workout_plan": gin.H{
			"id":       workout_day.WorkoutPlan.Id,
			"title":    workout_day.WorkoutPlan.Title,
			"overview": workout_day.WorkoutPlan.Overview,
			"tags":     workout_day.WorkoutPlan.Tags,
			"details":  workout_day.WorkoutPlan.Details,
			"creator": gin.H{
				"nickname":   workout_day.WorkoutPlan.Creator.Profile1.Nickname,
				"avatar_url": workout_day.WorkoutPlan.Creator.Profile1.AvatarURL,
			},
		},
		// "day_number":  day_number,
		"started_at":  workout_day.StartedAt,
		"finished_at": workout_day.FinishedAt,
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": data})
}

func (h *WorkoutDayHandler) UpdateWorkoutDayPlanDetails(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id   int    `json:"id"`
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少id参数", "data": nil})
		return
	}
	var existing models.WorkoutDay
	if err := h.db.Where("id = ? AND coach_id = ?", body.Id, uid).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "没有找到记录", "data": nil})
		return
	}
	// existing.UpdatedDetails = body.Data
	// now := time.Now().UTC()
	// existing.UpdatedAt = &now
	if err := h.db.Model(&existing).Update("updated_details", body.Data).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": gin.H{"id": existing.Id}})
}

// 暂存的训练内容记录
func (h *WorkoutDayHandler) UpdateWorkoutDayStepProgress(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id   int    `json:"id"`
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少id参数", "data": nil})
		return
	}
	if body.Data == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少更新内容", "data": nil})
		return
	}
	var existing models.WorkoutDay
	if result := h.db.Where("id = ? AND coach_id = ?", body.Id, uid).First(&existing); result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": result.Error.Error(), "data": nil})
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

	if err := tx.Model(&existing).Update("pending_steps", body.Data).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to update day", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction: ", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "更新失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": gin.H{"id": existing.Id}})
}

// 基本上用不上，都是用 createWorkoutDay
func (h *WorkoutDayHandler) StartWorkoutDay(c *gin.Context) {
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
	var day models.WorkoutDay
	if err := h.db.Where("student_id = ?", uid).First(&day, body.Id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "找不到记录", "data": nil})
		return
	}

	day.Status = int(models.WorkoutDayStatusStarted)
	now := time.Now().UTC()
	day.StartedAt = &now
	h.db.Save(&day)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "操作成功", "data": gin.H{"id": day.Id}})
}

func (h *WorkoutDayHandler) FinishWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id   int    `json:"id"`
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	if body.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少id参数", "data": nil})
		return
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		}
	}()

	var existing models.WorkoutDay
	if err := tx.Where("id = ? AND coach_id = ?", body.Id, uid).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "没有找到记录", "data": nil})
		return
	}

	existing.Status = int(models.WorkoutDayStatusFinished)
	if body.Data != "" {
		existing.PendingSteps = body.Data
	}
	now := time.Now().UTC()
	existing.FinishedAt = &now
	if err := tx.Save(&existing).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	type WorkoutDayPendingAction struct {
		Idx         int     `json:"idx"`
		ActionId    int     `json:"action_id"`
		Reps        int     `json:"reps"`
		RepsUnit    string  `json:"reps_unit"`
		Weight      float64 `json:"weight"`
		WeightUnit  string  `json:"weight_unit"`
		Remark      string  `json:"remark"`
		Completed   bool    `json:"completed"`
		CompletedAt int     `json:"completed_at"`
		Time1       *int    `json:"time1"`
		Time2       *int    `json:"time2"`
		Time3       *int    `json:"time3"`
	}
	type WorkoutDayPendingSet struct {
		StepIdx int                       `json:"step_idx"`
		Idx     int                       `json:"idx"`
		Actions []WorkoutDayPendingAction `json:"actions"`
	}

	type WorkoutDayProgressJSON250424 struct {
		Version string                 `json:"v"`
		StepIdx int                    `json:"step_idx"`
		SetIdx  int                    `json:"set_idx"`
		Sets    []WorkoutDayPendingSet `json:"sets"`
	}
	// 训练进行中的记录
	var pending_data WorkoutDayProgressJSON250424
	if existing.PendingSteps != "" {
		// fmt.Println(existing.PendingSteps)
		if err := json.Unmarshal([]byte(existing.PendingSteps), &pending_data); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 600, "msg": "The PendingSteps has something wrong", "data": nil})
			return
		}
		for _, set := range pending_data.Sets {
			for _, act := range set.Actions {
				if act.Completed {
					// Create workout action history record
					history := models.WorkoutActionHistory{
						WorkoutDayId:    body.Id,
						StudentId:       int(existing.StudentId),
						WorkoutActionId: act.ActionId,
						Reps:            act.Reps,
						RepsUnit:        act.RepsUnit,
						Weight:          int(act.Weight),
						WeightUnit:      act.WeightUnit,
						CreatedAt:       time.Unix(int64(act.CompletedAt), 0),
					}

					if err := tx.Create(&history).Error; err != nil {
						tx.Rollback()
						h.logger.Error("Failed to create workout action history", err)
					}
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to commit transaction", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to commit transaction", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "操作成功", "data": gin.H{"id": existing.Id}})
}

func (h *WorkoutDayHandler) GiveUpWorkoutDay(c *gin.Context) {
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
	var existing models.WorkoutDay
	if err := h.db.Where("id = ? AND coach_id = ?", body.Id, uid).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "没有找到记录", "data": nil})
		return
	}
	existing.Status = int(models.WorkoutDayStatusGiveUp)
	h.db.Save(&existing)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": gin.H{"id": existing.Id}})
}

func (h *WorkoutDayHandler) ContinueWorkoutDay(c *gin.Context) {
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
	var existing models.WorkoutDay
	if err := h.db.Where("id = ? AND coach_id = ?", body.Id, uid).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "没有找到记录", "data": nil})
		return
	}
	existing.Status = int(models.WorkoutDayStatusStarted)
	h.db.Save(&existing)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功", "data": gin.H{"id": existing.Id}})
}

func (h *WorkoutDayHandler) DeleteWorkoutDay(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}
	var existing models.WorkoutDay
	if err := h.db.Where("id = ? AND student_id = ?", body.Id, uid).First(&existing).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	if err := h.db.Where("id = ?", body.Id).Update("d", 1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "操作成功", "data": nil})
}

func (h *WorkoutDayHandler) FetchWorkoutDayList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		models.Pagination
		Status          int        `json:"status"`
		FinishedAtStart *time.Time `json:"finished_at_start"`
		FinishedAtEnd   *time.Time `json:"finished_at_end"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	query := h.db.Where("d IS NULL OR d = 0")
	if body.Status != 0 {
		query = query.Where("status = ?", body.Status)
	}
	// Add finished time range filter
	if body.FinishedAtStart != nil {
		query = query.Where("finished_at >= ?", body.FinishedAtStart)
	}
	if body.FinishedAtEnd != nil {
		query = query.Where("finished_at <= ?", body.FinishedAtEnd)
	}
	// student_id 表示是自己训练的记录
	query = query.Where("student_id = ?", uid)
	pb := pagination.NewPaginationBuilder[models.WorkoutDay](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetNextMarker(body.NextMarker).
		SetOrderBy("created_at DESC")
	var list1 []models.WorkoutDay
	if err := pb.Build().Preload("WorkoutPlan").Find(&list1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list1)
	list := make([]map[string]interface{}, 0, len(list2))
	for _, v := range list2 {
		list = append(list, map[string]interface{}{
			"id":          v.Id,
			"status":      v.Status,
			"group_no":    v.GroupNo,
			"started_at":  v.StartedAt,
			"finished_at": v.FinishedAt,
			"workout_plan": map[string]interface{}{
				"id":       v.WorkoutPlan.Id,
				"title":    v.WorkoutPlan.Title,
				"overview": v.WorkoutPlan.Overview,
				"tags":     v.WorkoutPlan.Tags,
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

func (h *WorkoutDayHandler) FetchFinishedWorkoutDayList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		Id              int        `json:"id"`
		FinishedAtStart *time.Time `json:"finished_at_start"`
		FinishedAtEnd   *time.Time `json:"finished_at_end"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	var list []models.WorkoutDay

	query := h.db.Where("d IS NULL OR d = 0")
	// student_id 表示是自己训练的记录
	query = query.Where("student_id = ?", uid)

	// 添加开始时间范围筛选
	if body.FinishedAtStart != nil {
		query = query.Where("finished_at >= ?", body.FinishedAtStart)
	}
	if body.FinishedAtEnd != nil {
		query = query.Where("finished_at <= ?", body.FinishedAtEnd)
	}

	query = query.Where("status = ?", int(models.WorkoutDayStatusFinished))
	query = query.Order("created_at desc")

	if err := query.Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	data := []map[string]interface{}{}

	for _, v := range list {
		data = append(data, map[string]interface{}{
			"id":          v.Id,
			"status":      v.Status,
			"started_at":  v.StartedAt,
			"finished_at": v.FinishedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list": data,
		},
	})
}

func (h *WorkoutDayHandler) FetchMyStudentWorkoutDayList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		models.Pagination
		Id             int        `json:"id"`
		StartedAtStart *time.Time `json:"started_at_start"`
		StartedAtEnd   *time.Time `json:"started_at_end"`
		Status         *int       `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	// 确保是自己的学员
	var relation models.CoachRelationship
	if err := h.db.Where("coach_id = ? AND student_id = ?", uid, body.Id).First(&relation).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	query := h.db
	query = query.Where("student_id = ?", body.Id)
	// 添加开始时间范围筛选
	if body.StartedAtStart != nil {
		query = query.Where("started_at >= ?", body.StartedAtStart)
	}
	if body.StartedAtEnd != nil {
		query = query.Where("started_at <= ?", body.StartedAtEnd)
	}
	// 添加状态筛选
	if body.Status != nil {
		query = query.Where("status = ?", *body.Status)
	}
	pb := pagination.NewPaginationBuilder[models.WorkoutDay](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")

	var list []models.WorkoutDay
	if err := pb.Build().Preload("WorkoutPlan").Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list)
	data := []map[string]interface{}{}
	for _, v := range list2 {
		data = append(data, map[string]interface{}{
			"id":          v.Id,
			"status":      v.Status,
			"started_at":  v.StartedAt,
			"finished_at": v.FinishedAt,
			"workout_plan": map[string]interface{}{
				"id":       v.WorkoutPlan.Id,
				"title":    v.WorkoutPlan.Title,
				"overview": v.WorkoutPlan.Overview,
				"tags":     v.WorkoutPlan.Tags,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"list":        data,
			"page_size":   pb.GetLimit(),
			"has_more":    has_more,
			"next_marker": next_marker,
		},
	})
}
