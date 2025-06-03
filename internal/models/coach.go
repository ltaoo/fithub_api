package models

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Coach 教练模型
type Coach struct {
	Id        int        `json:"id" gorm:"primaryKey"`
	Nickname  string     `json:"nickname"`
	Bio       string     `json:"bio" gorm:"default:''"`       // 个人简介
	CoachType int        `json:"coach_type" gorm:"default:1"` // 教练类型
	Status    int        `json:"status" gorm:"default:1"`     // 状态
	Config    string     `json:"config" gorm:"default:'{}'"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	Students   []CoachRelationship `json:"students" gorm:"foreignKey:CoachId"`  // 作为教练的关系
	Coaches    []CoachRelationship `json:"coaches" gorm:"foreignKey:StudentId"` // 作为学员的关系
	Profile1Id int                 `json:"profile1_id"`
	Profile1   CoachProfile1       `json:"profile1" gorm:"foreignKey:Profile1Id"`
	Profile2Id int                 `json:"profile2_id"`
	Profile2   CoachProfile2       `json:"profile2" gorm:"foreignKey:Profile2Id"`
}

func (Coach) TableName() string {
	return "COACH"
}

// CoachAccount 教练账号模型
type CoachAccount struct {
	ProviderType int       `json:"provider_type"`
	ProviderId   string    `json:"provider_id"`
	ProviderArg1 string    `json:"provider_arg1"`
	ProviderArg2 string    `json:"provider_arg2"`
	ProviderArg3 string    `json:"provider_arg3"`
	CreatedAt    time.Time `json:"created_at"`

	CoachId int   `json:"coach_id"`
	Coach   Coach `json:"coach"`
}

func (CoachAccount) TableName() string {
	return "COACH_ACCOUNT"
}

type CoachProfile1 struct {
	Id                  int    `json:"id"`
	CoachId             int    `json:"coach_id"`
	Nickname            string `json:"nickname"`
	AvatarURL           string `json:"avatar_url"`
	Age                 int    `json:"age"`
	Gender              int    `json:"gender"`
	BodyType            int    `json:"body_type"`
	Height              int    `json:"height"`
	Weight              int    `json:"weight"`
	BodyFatPercent      int    `json:"body_fat_percent"`
	RiskScreenings      string `json:"risk_screenings"`
	TrainingGoals       string `json:"training_goals"`
	TrainingFrequency   int    `json:"training_frequency"`
	TrainingPreferences string `json:"training_preferences"`
	DietPreferences     string `json:"diet_preferences"`
}

func (CoachProfile1) TableName() string {
	return "COACH_PROFILE1"
}

type CoachProfile2 struct {
	Id              int    `json:"id"`
	CoachId         int    `json:"coach_id"`
	Specialties     string `json:"specialties"`
	Certification   string `json:"certification"`
	ExperienceYears string `json:"experience_years"`
}

func (CoachProfile2) TableName() string {
	return "COACH_PROFILE2"
}

// CoachRelationship 教练-学员关系模型
type CoachRelationship struct {
	Id        int        `json:"id" gorm:"primaryKey"`
	CoachId   int        `json:"coach_id" gorm:"not null"`   // 教练ID
	StudentId int        `json:"student_id" gorm:"not null"` // 学员ID
	Status    int        `json:"status" gorm:"default:1"`    // 关系状态
	Role      int        `json:"role" gorm:"default:1"`      // 关系角色
	Note      string     `json:"note" gorm:"default:''"`     // 备注信息
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	// 关联
	Coach   Coach `json:"coach" gorm:"foreignKey:CoachId"`     // 教练信息
	Student Coach `json:"student" gorm:"foreignKey:StudentId"` // 学员信息
}

func (CoachRelationship) TableName() string {
	return "COACH_RELATIONSHIP"
}

// 常量定义
const (
	AccountProviderTypeEmailWithPwd = 1

	// CoachType 教练类型
	CoachTypePersonal = 1 // 私教
	CoachTypeGroup    = 2 // 团课
	CoachTypeOnline   = 3 // 在线

	// CoachStatus 教练状态
	CoachStatusNormal = 1 // 正常
	CoachStatusPaused = 2 // 暂停服务
	CoachStatusBanned = 3 // 封禁

	// RelationshipStatus 关系状态
	RelationPending   = 1 // 待确认
	RelationConfirmed = 2 // 已确认
	RelationRejected  = 3 // 已拒绝
	RelationDismissed = 4 // 已解除

	// RelationshipRole 关系角色
	RoleCoachStudent = 1 // 教练-学员
	RolePartner      = 2 // 合作教练
)

// AuthResponse represents the response for authentication operations
type AuthResponse struct {
	Token  string `json:"token"`
	Status string `json:"status"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey     []byte
	TokenDuration time.Duration
}

// Claims represents the JWT claims
type Claims struct {
	Id        float64 `json:"id"`
	ExpiresAt float64 `json:"expires_at"`
	Issuer    string  `json:"issuer"`
}

// Default JWT configuration
var DefaultJWTConfig = JWTConfig{
	SecretKey:     []byte("your-secret-key-change-in-production"), // Should be loaded from environment variables
	TokenDuration: 24 * time.Hour,                                 // 24 hours
}

// Helper function to generate JWT token
func GenerateJWT(coach_id int) (string, error) {
	expiration_time := time.Now().Add(DefaultJWTConfig.TokenDuration)

	fmt.Println("before generate claims")
	fmt.Println(coach_id)
	claims := jwt.MapClaims{
		"id":         coach_id,
		"expires_at": jwt.NewNumericDate(expiration_time),
		"issuer":     "top.fithub",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString(DefaultJWTConfig.SecretKey)

	if err != nil {
		return "", err
	}

	return str, nil
}

// ParseJWT parses and validates a JWT token
func ParseJWT(str string) (*Claims, error) {
	// Remove "Bearer " prefix if present
	str = strings.TrimPrefix(str, "Bearer ")
	token, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return DefaultJWTConfig.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	v, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}
	claims := &Claims{
		Id:        v["id"].(float64),
		ExpiresAt: v["expires_at"].(float64),
		Issuer:    v["issuer"].(string),
	}
	return claims, nil
}

// AuthMiddleware is a middleware to authenticate requests
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Authorization header is required", "data": nil})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Invalid or expired token", "data": nil})
			c.Abort()
			return
		}

		// Set coach ID in context
		c.Set("coachId", claims.Id)
		c.Next()
	}
}
