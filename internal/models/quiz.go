package models

import (
	"time"
)

// Quiz represents a quiz question in the system
type Quiz struct {
	Id         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Content    string    `json:"content" gorm:"not null;default:''"`
	Overview   string    `json:"overview" gorm:"not null;default:''"`
	Medias     string    `json:"medias" gorm:"not null;default:''"`
	Type       int       `json:"type" gorm:"not null;default:1"` // 1单选 2多选 3判断 4填空 5简答
	Difficulty int       `json:"difficulty" gorm:"not null;default:1"`
	Tags       string    `json:"tags" gorm:"not null;default:''"`
	Analysis   string    `json:"analysis" gorm:"not null;default:''"`
	Choices    string    `json:"choices" gorm:"not null;default:'{}'"`
	Answer     string    `json:"answer" gorm:"not null;default:'{}'"`
	CreatorId  int       `json:"creator_id" gorm:"not null;default:0"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null;"`
}

// TableName specifies the table name for Quiz
func (Quiz) TableName() string {
	return "QUIZ"
}

// Paper represents a quiz paper in the system
type Paper struct {
	Id        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"not null;default:''"`
	Overview  string    `json:"overview" gorm:"not null;default:''"`
	Tags      string    `json:"tags" gorm:"not null;default:''"`
	Duration  int       `json:"duration" gorm:"not null;default:0"` // 试卷答题时长，单位 分钟
	PassScore int       `json:"pass_score" gorm:"not null;default:0"`
	QuizCount int       `json:"quiz_count" gorm:"not null;default:0"`
	CreatorId int       `json:"creator_id" gorm:"not null;default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for Paper
func (Paper) TableName() string {
	return "PAPER"
}

// PaperQuiz represents the relationship between papers and quizzes
type PaperQuiz struct {
	Id      int `json:"id" gorm:"primaryKey;autoIncrement"`
	Score   int `json:"score" gorm:"not null;default:1"`    // 题目分数
	SortIdx int `json:"sort_idx" gorm:"not null;default:0"` // 排序
	Visible int `json:"visible" gorm:"not null;default:1"`  // 是否可见

	PaperId int   `json:"quiz_paper_id" gorm:"not null;default:0"` // 试卷id
	Paper   Paper `json:"paper"`
	QuizId  int   `json:"quiz_id" gorm:"not null;default:0"` // 题目id
	Quiz    Quiz  `json:"quiz"`
}

// TableName specifies the table name for PaperQuiz
func (PaperQuiz) TableName() string {
	return "PAPER_QUIZ"
}

// Exam represents an exam record in the system
type Exam struct {
	Id          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Status      int        `json:"status" gorm:"not null;default:1"` // 1待开始 2进行中 3已完成 4手动放弃
	CurQuizId   int        `json:"cur_quiz_id"`
	Score       int        `json:"score" gorm:"not null;default:0"`
	CorrectRate int        `json:"correct_rate" gorm:"not null;default:0"`
	Pass        int        `json:"pass" gorm:"not null;default:0"` // 0否 1是
	PaperId     int        `json:"paper_id" gorm:"not null;default:0"`
	StudentId   int        `json:"student_id" gorm:"not null;default:0"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"` // 完成时间 包含交卷、超时自动提交
	UpdatedAt   *time.Time `json:"updated_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for Exam
func (Exam) TableName() string {
	return "EXAM"
}

// QuizAnswer represents a quiz answer record in the system
type QuizAnswer struct {
	Id        int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Status    int        `json:"status" gorm:"not null;default:0"`     // 题目结果 0 1正确 2失败 3跳过
	Answer    string     `json:"answer" gorm:"not null;default:'{}'"`  // 答题内容，JSON 根据题目类型有不同结构
	Score     int        `json:"score" gorm:"not null;default:0"`      // 得分
	StudentId int        `json:"student_id" gorm:"not null;default:0"` // 答题人id
	UpdatedAt *time.Time `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	QuizId  int   `json:"quiz_id" gorm:"not null;default:0"` // 题目id
	Quiz    Quiz  `json:"quiz"`
	ExamId  int   `json:"exam_id" gorm:"not null;default:0"` // 考试id
	Exam    Exam  `json:"exam"`
	PaperId int   `json:"paper_id" gorm:"not null;default:0"` // 试卷id
	Paper   Paper `json:"paper"`
}

// TableName specifies the table name for QuizAnswer
func (QuizAnswer) TableName() string {
	return "QUIZ_ANSWER"
}
