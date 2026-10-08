package handlers

import (
	"github.com/gin-gonic/gin"
	"myapi/internal/models"
	"net/http"
)

func (h *WorkoutPlanHandler) can_read_plan(ctx *gin.Context, plan_id int) bool {
	var record models.WorkoutPlan
	if err := h.db.Where("id = ? AND (d IS NULL OR d = 0)", plan_id).First(&record).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 404, "msg": "找不到训练计划", "data": nil})
		return false
	}
	user_id := int(ctx.GetFloat64("id"))
	if record.Status != int(models.WorkoutPublishStatusPublic) && (user_id == 0 || record.OwnerId != user_id) {
		ctx.JSON(http.StatusOK, gin.H{"code": 403, "msg": "没有权限查看训练计划", "data": nil})
		return false
	}
	return true
}
