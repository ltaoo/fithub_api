package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myapi/internal/models"
	"myapi/internal/pkg/pagination"
	"myapi/pkg/logger"
)

type QuizHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewQuizHandler(db *gorm.DB, logger *logger.Logger) *QuizHandler {
	return &QuizHandler{
		db:     db,
		logger: logger,
	}
}

func (h *QuizHandler) FetchPaperList(c *gin.Context) {
	var body struct {
		models.Pagination
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	query := h.db
	pb := pagination.NewPaginationBuilder[models.Paper](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")
	var list []models.Paper
	if err := pb.Build().Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch discount policies", "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list)
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

func (h *QuizHandler) FetchQuizList(c *gin.Context) {
	var body struct {
		models.Pagination
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	query := h.db
	pb := pagination.NewPaginationBuilder[models.Quiz](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")

	var list []models.Quiz
	if r := pb.Build().Find(&list); r.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch records", "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list)
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

func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Content    string `json:"content"`
		Overview   string `json:"overview"`
		Type       int    `json:"type"`
		Tags       string `json:"tags"`
		Difficulty int    `json:"difficulty"`
		Choices    string `json:"choices"`
		Answer     string `json:"answer"`
		Analysis   string `json:"analysis"`
		Medias     string `json:"medias"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
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
	paper := models.Quiz{
		Content:    body.Content,
		Overview:   body.Overview,
		Type:       body.Type,
		Tags:       body.Tags,
		Medias:     body.Medias,
		Difficulty: body.Difficulty,
		Analysis:   body.Analysis,
		Choices:    body.Choices,
		Answer:     body.Answer,
		CreatorId:  uid,
		CreatedAt:  time.Now(),
	}

	if err := tx.Create(&paper).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create subscription plan", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create subscription plan", "data": nil})
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
		"data": paper,
	})
}

func (h *QuizHandler) CreatePaper(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Name      string `json:"name"`
		Overview  string `json:"overview"`
		Tags      string `json:"tags"`
		PassScore int    `json:"pass_score"`
		Duration  int    `json:"duration"`
		QuizList  []struct {
			Id      int `json:"id"`
			Score   int `json:"score"`
			SortIdx int `json:"sort_idx"`
		} `json:"quiz_list"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// 开始事务
	tx := h.db.Begin()
	if tx.Error != nil {
		h.logger.Error("Failed to start transaction", tx.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		return
	}

	// 创建试卷
	paper := models.Paper{
		Name:      body.Name,
		Overview:  body.Overview,
		Tags:      body.Tags,
		QuizCount: len(body.QuizList),
		PassScore: body.PassScore,
		Duration:  body.Duration,
		CreatorId: uid,
		CreatedAt: time.Now(),
	}

	if err := tx.Create(&paper).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create paper", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create paper", "data": nil})
		return
	}

	// 创建试卷题目关联
	for _, quiz := range body.QuizList {
		paperQuiz := models.PaperQuiz{
			Visible: 1,
			Score:   quiz.Score,
			SortIdx: quiz.SortIdx,
			PaperId: paper.Id,
			QuizId:  quiz.Id,
		}
		if err := tx.Create(&paperQuiz).Error; err != nil {
			tx.Rollback()
			h.logger.Error("Failed to create paper quiz relation", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create paper quiz relation", "data": nil})
			return
		}
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
		"data": paper,
	})
}

func (h *QuizHandler) UpdatePaper(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Tags      string `json:"tags"`
		PassScore int    `json:"pass_score"`
		Duration  int    `json:"duration"`
		QuizList  []struct {
			RelationId int `json:"relation_id"`
			Id         int `json:"id"`
			Score      int `json:"score"`
			SortIdx    int `json:"sort_idx"`
		} `json:"quiz_list"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// 开始事务
	tx := h.db.Begin()
	if tx.Error != nil {
		h.logger.Error("Failed to start transaction", tx.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		return
	}

	// 检查试卷是否存在
	var paper models.Paper
	if err := tx.First(&paper, body.Id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Paper not found", "data": nil})
		} else {
			h.logger.Error("Failed to find paper", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to find paper", "data": nil})
		}
		return
	}

	// 检查权限
	if paper.CreatorId != uid {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 403, "msg": "Permission denied", "data": nil})
		return
	}

	// 更新试卷信息
	updates := map[string]interface{}{
		"name":       body.Name,
		"tags":       body.Tags,
		"pass_score": body.PassScore,
		"duration":   body.Duration,
	}
	if err := tx.Model(&paper).Updates(updates).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to update paper", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update paper", "data": nil})
		return
	}

	// 获取现有的试卷题目关联
	var existing_relations []models.PaperQuiz
	if err := tx.Where("paper_id = ?", body.Id).Find(&existing_relations).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to fetch existing paper quiz relations", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch existing relations", "data": nil})
		return
	}

	// 创建现有关联ID的映射，用于快速查找
	existing_map := make(map[int]models.PaperQuiz)
	for _, rel := range existing_relations {
		existing_map[rel.Id] = rel
	}

	// 创建新关联ID的映射，用于快速查找
	new_relation_map := make(map[int]struct{})
	for _, quiz := range body.QuizList {
		if quiz.RelationId > 0 {
			new_relation_map[quiz.RelationId] = struct{}{}
		}
	}

	// 处理每个题目关联
	for _, quiz := range body.QuizList {
		if quiz.RelationId > 0 {
			// 更新现有关联
			if _, exists := existing_map[quiz.RelationId]; exists {
				updates := map[string]interface{}{
					"score":    quiz.Score,
					"sort_idx": quiz.SortIdx,
				}
				if err := tx.Model(&models.PaperQuiz{}).Where("id = ?", quiz.RelationId).Updates(updates).Error; err != nil {
					tx.Rollback()
					h.logger.Error("Failed to update paper quiz relation", err)
					c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update relation", "data": nil})
					return
				}
			}
		} else {
			// 创建新关联
			paperQuiz := models.PaperQuiz{
				Visible: 1,
				Score:   quiz.Score,
				SortIdx: quiz.SortIdx,
				PaperId: paper.Id,
				QuizId:  quiz.Id,
			}
			if err := tx.Create(&paperQuiz).Error; err != nil {
				tx.Rollback()
				h.logger.Error("Failed to create paper quiz relation", err)
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create relation", "data": nil})
				return
			}
		}
	}

	// 删除不再需要的关联
	for _, existingRel := range existing_relations {
		if _, exists := new_relation_map[existingRel.Id]; !exists {
			if err := tx.Delete(&models.PaperQuiz{}, existingRel.Id).Error; err != nil {
				tx.Rollback()
				h.logger.Error("Failed to delete paper quiz relation", err)
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to delete relation", "data": nil})
				return
			}
		}
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
		"data": paper,
	})
}

func (h *QuizHandler) FetchPaperProfile(c *gin.Context) {
	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// 获取试卷信息
	var paper models.Paper
	if err := h.db.First(&paper, body.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Paper not found", "data": nil})
		} else {
			h.logger.Error("Failed to fetch paper", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch paper", "data": nil})
		}
		return
	}

	// 获取试卷关联的题目
	var paper_quizzes []models.PaperQuiz
	if err := h.db.Where("paper_id = ?", body.Id).Order("sort_idx asc").Preload("Quiz").Find(&paper_quizzes).Error; err != nil {
		h.logger.Error("Failed to fetch paper quizzes", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch paper quizzes", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": gin.H{
			"paper":   paper,
			"quizzes": paper_quizzes,
		},
	})
}

func (h *QuizHandler) StartExamWithPaper(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		PaperId int `json:"paper_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// 开始事务
	tx := h.db.Begin()
	if tx.Error != nil {
		h.logger.Error("Failed to start transaction", tx.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		return
	}

	// 检查试卷是否存在
	var paper models.Paper
	if err := tx.First(&paper, body.PaperId).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Paper not found", "data": nil})
		} else {
			h.logger.Error("Failed to find paper", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to find paper", "data": nil})
		}
		return
	}

	now := time.Now()
	// 创建考试记录
	exam := models.Exam{
		Status:    2, // 进行中
		PaperId:   body.PaperId,
		StudentId: uid,
		StartedAt: &now,
		CreatedAt: now,
	}

	if err := tx.Create(&exam).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create exam", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to start exam", "data": nil})
		return
	}

	// 获取试卷关联的题目
	var paper_quizzes []models.PaperQuiz
	if err := tx.Where("paper_id = ?", body.PaperId).Order("sort_idx asc").Find(&paper_quizzes).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to fetch paper quizzes", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch paper quizzes", "data": nil})
		return
	}

	// 为每个题目创建答题记录
	for _, pq := range paper_quizzes {
		quiz_answer := models.QuizAnswer{
			Status:    0, // 未作答
			StudentId: uid,
			QuizId:    pq.QuizId,
			ExamId:    exam.Id,
			PaperId:   body.PaperId,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&quiz_answer).Error; err != nil {
			tx.Rollback()
			h.logger.Error("Failed to create quiz answer", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to create quiz answer", "data": nil})
			return
		}
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
		"data": exam,
	})
}

func (h *QuizHandler) FetchRunningExam(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	// 获取用户正在进行的考试
	var exams []models.Exam
	if err := h.db.Where("student_id = ? AND status = ?", uid, 2).Find(&exams).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": nil})
			return
		}
		h.logger.Error("Failed to fetch running exam", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch running exam", "data": nil})
		return
	}

	// // 获取试卷信息
	// var paper models.Paper
	// if err := h.db.First(&paper, exam.PaperId).Error; err != nil {
	// 	h.logger.Error("Failed to fetch paper", err)
	// 	c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch paper", "data": nil})
	// 	return
	// }

	// // 获取答题记录
	// var quiz_answers []models.QuizAnswer
	// if err := h.db.Where("exam_id = ?", exam.Id).Preload("Quiz").Order("id asc").Find(&quiz_answers).Error; err != nil {
	// 	h.logger.Error("Failed to fetch quiz answers", err)
	// 	c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch quiz answers", "data": nil})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": gin.H{
			"list": exams,
		},
	})
}

func (h *QuizHandler) FetchExamList(c *gin.Context) {
	uid := int(c.GetFloat64("id"))
	var body struct {
		models.Pagination
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}
	query := h.db
	pb := pagination.NewPaginationBuilder[models.Exam](query).
		SetLimit(body.PageSize).
		SetPage(body.Page).
		SetOrderBy("created_at DESC")
	var list []models.Exam
	if err := pb.Build().Where("student_id = ?", uid).Preload("Paper").Find(&list).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": nil})
			return
		}
		h.logger.Error("Failed to fetch exam list", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch exam", "data": nil})
		return
	}
	list2, has_more, next_marker := pb.ProcessResults(list)
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

func (h *QuizHandler) FetchExamProfile(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// 获取考试信息
	var exam models.Exam
	if err := h.db.First(&exam, body.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Exam not found", "data": nil})
		} else {
			h.logger.Error("Failed to fetch exam", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch exam", "data": nil})
		}
		return
	}

	// 检查权限
	if exam.StudentId != uid {
		c.JSON(http.StatusOK, gin.H{"code": 403, "msg": "Permission denied", "data": nil})
		return
	}

	// 获取试卷信息
	var paper models.Paper
	if err := h.db.First(&paper, exam.PaperId).Error; err != nil {
		h.logger.Error("Failed to fetch paper", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch paper", "data": nil})
		return
	}

	// 获取答题记录
	var quiz_answers []models.QuizAnswer
	if err := h.db.Where("exam_id = ?", exam.Id).Preload("Quiz").Order("id asc").Find(&quiz_answers).Error; err != nil {
		h.logger.Error("Failed to fetch quiz answers", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch quiz answers", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": gin.H{
			"exam":         exam,
			"paper":        paper,
			"quiz_answers": quiz_answers,
		},
	})
}

func (h *QuizHandler) UpdateQuizAnswer(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Id      int    `json:"id"`
		ExamId  int    `json:"exam_id"`
		QuizId  int    `json:"quiz_id"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// 开始事务
	tx := h.db.Begin()
	if tx.Error != nil {
		h.logger.Error("Failed to start transaction", tx.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		return
	}

	// 检查考试是否存在且属于当前用户
	var exam models.Exam
	if err := tx.First(&exam, body.ExamId).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Exam not found", "data": nil})
		} else {
			h.logger.Error("Failed to find exam", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to find exam", "data": nil})
		}
		return
	}

	if exam.StudentId != uid {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 403, "msg": "Permission denied", "data": nil})
		return
	}

	// 检查考试状态
	if exam.Status != 2 { // 2表示进行中
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Exam is not in progress", "data": nil})
		return
	}
	// 更新考试进度
	exam_updates := map[string]interface{}{
		"cur_quiz_id": body.QuizId,
	}
	if err := tx.Model(&exam).Updates(exam_updates).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to update exam", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update exam", "data": nil})
		return
	}

	// 获取答题记录
	var quiz_answer models.QuizAnswer
	if err := tx.Where("exam_id = ? AND quiz_id = ?", body.ExamId, body.QuizId).First(&quiz_answer).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Quiz answer not found", "data": nil})
		} else {
			h.logger.Error("Failed to find quiz answer", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to find quiz answer", "data": nil})
		}
		return
	}

	// 获取题目信息
	var quiz models.Quiz
	if err := tx.First(&quiz, body.QuizId).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to find quiz", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to find quiz", "data": nil})
		return
	}

	// 获取试卷题目关联信息
	var paper_quiz models.PaperQuiz
	if err := tx.Where("paper_id = ? AND quiz_id = ?", exam.PaperId, body.QuizId).First(&paper_quiz).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to find paper quiz relation", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to find paper quiz relation", "data": nil})
		return
	}

	// 解析用户答案
	var user_answer struct {
		Choices []int  `json:"choices"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal([]byte(body.Content), &user_answer); err != nil {
		tx.Rollback()
		h.logger.Error("Failed to parse user answer", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to parse user answer", "data": nil})
		return
	}

	// 解析正确答案
	var correct_answer struct {
		Value []int `json:"value"`
	}
	if err := json.Unmarshal([]byte(quiz.Answer), &correct_answer); err != nil {
		tx.Rollback()
		h.logger.Error("Failed to parse correct answer", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to parse correct answer", "data": nil})
		return
	}

	// 根据题目类型判断答案是否正确
	var status int
	var score int

	switch quiz.Type {
	case 1: // 单选
		fmt.Println("check the answer is correct1")
		fmt.Println(user_answer.Choices)
		fmt.Println(correct_answer.Value)
		fmt.Println(quiz.Answer)
		if len(user_answer.Choices) == len(correct_answer.Value) {
			// 创建一个map来存储正确答案中的元素
			correct_map := make(map[int]bool)
			for _, choice := range correct_answer.Value {
				correct_map[choice] = true
			}

			// 检查用户答案中的每个元素是否都在正确答案中
			all_correct := true
			for _, choice := range user_answer.Choices {
				if !correct_map[choice] {
					all_correct = false
					break
				}
			}
			if all_correct {
				status = 1 // 正确
				score = paper_quiz.Score
			} else {
				status = 2 // 错误
				score = 0
			}
		} else {
			status = 2 // 错误
			score = 0
		}
	case 2: // 多选
		fmt.Println("check the answer is correct1")
		fmt.Println(user_answer.Choices)
		fmt.Println(correct_answer.Value)
		// 检查两个数组是否包含相同的元素（顺序不重要）
		if len(user_answer.Choices) == len(correct_answer.Value) {
			// 创建一个map来存储正确答案中的元素
			correct_map := make(map[int]bool)
			for _, choice := range correct_answer.Value {
				correct_map[choice] = true
			}
			// 检查用户答案中的每个元素是否都在正确答案中
			all_correct := true
			for _, choice := range user_answer.Choices {
				if !correct_map[choice] {
					all_correct = false
					break
				}
			}
			if all_correct {
				status = 1 // 正确
				score = paper_quiz.Score
			} else {
				status = 2 // 错误
				score = 0
			}
		} else {
			status = 2 // 错误
			score = 0
		}
	case 3: // 判断
		if body.Content == quiz.Answer {
			status = 1
			score = paper_quiz.Score
		} else {
			status = 2
			score = 0
		}
	case 4: // 填空
		// TODO: 实现填空题答案判断逻辑
		status = 0
		score = 0
	case 5: // 简答
		// TODO: 实现简答题答案判断逻辑
		status = 0
		score = 0
	default:
		status = 0
		score = 0
	}

	updates := map[string]interface{}{
		"answer":     body.Content,
		"status":     status,
		"score":      score,
		"updated_at": time.Now(),
	}

	if err := tx.Model(&quiz_answer).Updates(updates).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to update quiz answer", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update quiz answer", "data": nil})
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
		"data": quiz_answer,
	})
}

func (h *QuizHandler) CompleteExam(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// 开始事务
	tx := h.db.Begin()
	if tx.Error != nil {
		h.logger.Error("Failed to start transaction", tx.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		return
	}

	// 检查考试是否存在且属于当前用户
	var exam models.Exam
	if err := tx.First(&exam, body.Id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Exam not found", "data": nil})
		} else {
			h.logger.Error("Failed to find exam", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to find exam", "data": nil})
		}
		return
	}

	if exam.StudentId != uid {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 403, "msg": "Permission denied", "data": nil})
		return
	}

	// 检查考试状态
	if exam.Status != 2 { // 2表示进行中
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Exam is not in progress", "data": nil})
		return
	}

	// 获取试卷信息
	var paper models.Paper
	if err := tx.First(&paper, exam.PaperId).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to find paper", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to find paper", "data": nil})
		return
	}

	// 获取所有答题记录
	var quiz_answers []models.QuizAnswer
	if err := tx.Where("exam_id = ?", exam.Id).Find(&quiz_answers).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to fetch quiz answers", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch quiz answers", "data": nil})
		return
	}

	// 计算总分
	var total_score int
	for _, answer := range quiz_answers {
		total_score += answer.Score
	}

	status := 3
	// 判断是否通过
	pass := 0
	if total_score >= paper.PassScore {
		pass = 1
	}

	// 更新考试状态
	now := time.Now()
	updates := map[string]interface{}{
		"status":       status,
		"pass":         pass,
		"score":        total_score,
		"completed_at": now,
	}

	if err := tx.Model(&exam).Updates(updates).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to update exam", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update exam", "data": nil})
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
		"data": gin.H{
			"exam":        exam,
			"total_score": total_score,
			"is_passed":   pass,
		},
	})
}

func (h *QuizHandler) GiveUpExam(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// 开始事务
	tx := h.db.Begin()
	if tx.Error != nil {
		h.logger.Error("Failed to start transaction", tx.Error)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Internal server error", "data": nil})
		return
	}

	// 检查考试是否存在且属于当前用户
	var exam models.Exam
	if err := tx.First(&exam, body.Id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Exam not found", "data": nil})
		} else {
			h.logger.Error("Failed to find exam", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to find exam", "data": nil})
		}
		return
	}

	if exam.StudentId != uid {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 403, "msg": "Permission denied", "data": nil})
		return
	}

	// 检查考试状态
	if exam.Status != 2 { // 2表示进行中
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Exam is not in progress", "data": nil})
		return
	}

	// 更新考试状态为已放弃
	now := time.Now()
	updates := map[string]interface{}{
		"status":     4, // 4表示已放弃
		"give_up_at": now,
	}

	if err := tx.Model(&exam).Updates(updates).Error; err != nil {
		tx.Rollback()
		h.logger.Error("Failed to update exam", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to update exam", "data": nil})
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
		"data": exam,
	})
}

func (h *QuizHandler) FetchExamResult(c *gin.Context) {
	uid := int(c.GetFloat64("id"))

	var body struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body", "data": nil})
		return
	}

	// 获取考试信息
	var exam models.Exam
	if err := h.db.First(&exam, body.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Exam not found", "data": nil})
		} else {
			h.logger.Error("Failed to fetch exam", err)
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch exam", "data": nil})
		}
		return
	}

	// 检查权限
	if exam.StudentId != uid {
		c.JSON(http.StatusOK, gin.H{"code": 403, "msg": "Permission denied", "data": nil})
		return
	}

	// 获取试卷信息
	var paper models.Paper
	if err := h.db.First(&paper, exam.PaperId).Error; err != nil {
		h.logger.Error("Failed to fetch paper", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch paper", "data": nil})
		return
	}

	// 获取答题记录
	var quiz_answers []models.QuizAnswer
	if err := h.db.Where("exam_id = ?", exam.Id).Preload("Quiz").Order("id asc").Find(&quiz_answers).Error; err != nil {
		h.logger.Error("Failed to fetch quiz answers", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Failed to fetch quiz answers", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": gin.H{
			"exam":         exam,
			"paper":        paper,
			"quiz_answers": quiz_answers,
		},
	})
}
