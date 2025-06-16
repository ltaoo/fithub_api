package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Coach 教练模型
type Coach struct {
	Id           int        `json:"id" gorm:"primaryKey"`
	D            int        `json:"d"`
	Nickname     string     `json:"nickname"`
	Bio          string     `json:"bio" gorm:"default:''"`       // 个人简介
	CoachType    int        `json:"coach_type" gorm:"default:1"` // 教练类型
	Status       int        `json:"status" gorm:"default:1"`     // 状态
	Config       string     `json:"config" gorm:"default:'{}'"`
	WorkoutStats string     `json:"workout_stats"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`

	Profile1Id int                 `json:"profile1_id"`
	Profile1   CoachProfile1       `json:"profile1" gorm:"foreignKey:Profile1Id"`
	Profile2Id int                 `json:"profile2_id"`
	Profile2   CoachProfile2       `json:"profile2" gorm:"foreignKey:Profile2Id"`
	Students   []CoachRelationship `json:"students" gorm:"foreignKey:CoachId"`  // 作为教练的关系
	Coaches    []CoachRelationship `json:"coaches" gorm:"foreignKey:StudentId"` // 作为学员的关系
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
	Status    int        `json:"status" gorm:"default:1"` // 关系状态
	Role      int        `json:"role" gorm:"default:1"`   // 关系角色
	Note      string     `json:"note" gorm:"default:''"`  // 备注信息
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	CoachId   int   `json:"coach_id" gorm:"not null"`            // 教练ID
	Coach     Coach `json:"coach" gorm:"foreignKey:CoachId"`     // 教练信息
	StudentId int   `json:"student_id" gorm:"not null"`          // 学员ID
	Student   Coach `json:"student" gorm:"foreignKey:StudentId"` // 学员信息
}

func (CoachRelationship) TableName() string {
	return "COACH_RELATIONSHIP"
}

type CoachContent struct {
	Id            int        `json:"id" gorm:"primaryKey"`
	ContentType   int        `json:"content_type" gorm:"default:0"` //
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	ContentURL    string     `json:"content_url"`
	CoverImageURL string     `json:"cover_image_url"`
	VideoKey      string     `json:"video_key"`
	ImageKeys     string     `json:"image_keys"`
	LikeCount     int        `json:"like_count"`
	Status        int        `json:"status"`  // 1审核通过
	Publish       int        `json:"publish"` // 1公开 2私有
	PublishedAt   *time.Time `json:"published_at"`
	D             int        `json:"d"`
	CreatedAt     time.Time  `json:"created_at"`

	CoachId int   `json:"coach_id"`
	Coach   Coach `json:"coach"  gorm:"foreignKey:CoachId"`
}

func (CoachContent) TableName() string {
	return "COACH_CONTENT"
}

type CoachContentWithWorkoutAction struct {
	Id         int       `json:"id" gorm:"primaryKey"`
	SortIdx    int       `json:"sort_idx"`
	StartPoint int       `json:"start_point"`
	D          int       `json:"d"`
	Details    string    `json:"details"`
	CreatedAt  time.Time `json:"created_at"`

	CoachContentId  int           `json:"coach_content_id"`
	Content         CoachContent  `json:"content" gorm:"foreignKey:CoachContentId"`
	WorkoutActionId int           `json:"workout_action_id"`
	WorkoutAction   WorkoutAction `json:"workout_action" gorm:"foreignKey:WorkoutActionId"`
}

func (CoachContentWithWorkoutAction) TableName() string {
	return "COACH_CONTENT_WITH_WORKOUT_ACTION"
}

type CoachContentWithWorkoutPlan struct {
	Id        int       `json:"id" gorm:"primaryKey"`
	SortIdx   int       `json:"sort_idx"`
	D         int       `json:"d"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`

	CoachContentId int          `json:"coach_content_id"`
	Content        CoachContent `json:"content" gorm:"foreignKey:CoachContentId"`
	WorkoutPlanId  int          `json:"workout_plan_id"`
	WorkoutPlan    WorkoutPlan  `json:"workout_plan" gorm:"foreignKey:WorkoutPlanId"`
}

func (CoachContentWithWorkoutPlan) TableName() string {
	return "COACH_CONTENT_WITH_WORKOUT_PLAN"
}

type CoachFollow struct {
	Id        int        `json:"id" gorm:"primaryKey"`
	Status    int        `json:"status"` // 1关注中 2取消关注
	UpdatedAt *time.Time `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`

	FollowingId int   `json:"following_id"`
	Following   Coach `json:"following" gorm:"foreignKey:FollowingId"`
	FollowerId  int   `json:"follower_id"`
	Follower    Coach `json:"follower" gorm:"foreignKey:FollowingId"`
}

func (CoachFollow) TableName() string {
	return "COACH_FOLLOW"
}

type MediaSocialPlatform struct {
	Id          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"image_keys"`
	LogoURL     string    `json:"logo_url"`
	HomepageURL string    `json:"homepage_url"`
	CreatedAt   time.Time `json:"created_at"`
}

func (MediaSocialPlatform) TableName() string {
	return "INFLUENCER_PLATFORM"
}

type CoachMediaSocialAccount struct {
	Id             int        `json:"id" gorm:"primaryKey"`
	D              int        `json:"d"`
	Status         int        `json:"status"`
	Nickname       string     `json:"nickname"`
	NicknameUsed   string     `json:"nickname_used"`
	AvatarURL      string     `json:"avatar_url"`
	Handle         string     `json:"handle"`
	AccountURL     string     `json:"account_url"`
	FollowersCount int        `json:"followers_count"`
	LogoURL        string     `json:"logo_url"`
	HomepageURL    string     `json:"homepage_url"`
	UpdatedAt      *time.Time `json:"updated_at"`
	CreatedAt      time.Time  `json:"created_at"`

	CoachId    int                 `json:"coach_id"`
	Coach      Coach               `json:"coach"  gorm:"foreignKey:CoachId"`
	PlatformId int                 `json:"platform_id"`
	Platform   MediaSocialPlatform `json:"platform"  gorm:"foreignKey:PlatformId"`
}

func (CoachMediaSocialAccount) TableName() string {
	return "INFLUENCER_ACCOUNT_IN_PLATFORM"
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
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	Status    string `json:"status"`
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
	TokenDuration: 48 * time.Hour,                                 // 24 hours
	// TokenDuration: 5 * time.Minute, // 5分钟，测试用
}

// Helper function to generate JWT token
func GenerateJWT(coach_id int) (string, time.Time, error) {
	expiration_time := time.Now().Add(DefaultJWTConfig.TokenDuration)

	claims := jwt.MapClaims{
		"id":         coach_id,
		"expires_at": jwt.NewNumericDate(expiration_time),
		"issuer":     "top.fithub",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString(DefaultJWTConfig.SecretKey)

	if err != nil {
		return "", time.Now(), err
	}

	return str, expiration_time, nil
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
