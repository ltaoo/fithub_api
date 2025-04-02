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
	Id              uint       `json:"id" gorm:"primaryKey"`
	Nickname        string     `json:"nickname" gorm:"not null"`
	AvatarURL       string     `json:"avatar_url"`
	Bio             string     `json:"bio" gorm:"default:''"`             // 个人简介
	Specialties     string     `json:"specialties" gorm:"default:''"`     // 专长领域
	Certification   string     `json:"certification" gorm:"default:'{}'"` // 认证信息
	ExperienceYears int        `json:"experience_years" gorm:"default:0"` // 从业年限
	CoachType       int        `json:"coach_type" gorm:"default:1"`       // 教练类型
	Status          int        `json:"status" gorm:"default:1"`           // 状态
	Config          string     `json:"config" gorm:"default:'{}'"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`

	// 关联
	Students []CoachRelationship `json:"students" gorm:"foreignKey:CoachId"`  // 作为教练的关系
	Coaches  []CoachRelationship `json:"coaches" gorm:"foreignKey:StudentId"` // 作为学员的关系
}

func (Coach) TableName() string {
	return "COACH"
}

// CoachAccount 教练账号模型
type CoachAccount struct {
	ProviderType string    `json:"provider_type" gorm:"primaryKey"`
	ProviderId   string    `json:"provider_id" gorm:"primaryKey"`
	ProviderArg1 int       `json:"provider_arg1"`
	ProviderArg2 string    `json:"provider_arg2"`
	ProviderArg3 string    `json:"provider_arg3"`
	CoachId      uint      `json:"coach_id" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
}

func (CoachAccount) TableName() string {
	return "COACH_ACCOUNT"
}

// CoachRelationship 教练-学员关系模型
type CoachRelationship struct {
	Id        uint       `json:"id" gorm:"primaryKey"`
	CoachId   uint       `json:"coach_id" gorm:"not null"`   // 教练ID
	StudentId uint       `json:"student_id" gorm:"not null"` // 学员ID
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
	Coach  Coach  `json:"coach"`
	Status string `json:"status"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey     string
	TokenDuration time.Duration
}

// Claims represents the JWT claims
type Claims struct {
	CoachId string `json:"coach_id"`
	jwt.RegisteredClaims
}

// Default JWT configuration
var DefaultJWTConfig = JWTConfig{
	SecretKey:     "your-secret-key-change-in-production", // Should be loaded from environment variables
	TokenDuration: 24 * time.Hour,                         // 24 hours
}

// Helper function to generate JWT token
func GenerateJWT(coachId string) (string, error) {
	expirationTime := time.Now().Add(DefaultJWTConfig.TokenDuration)

	claims := &Claims{
		CoachId: coachId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "myapi",
			Subject:   coachId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(DefaultJWTConfig.SecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJWT parses and validates a JWT token
func ParseJWT(tokenString string) (*Claims, error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(DefaultJWTConfig.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
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
		c.Set("coachId", claims.CoachId)
		c.Next()
	}
}
