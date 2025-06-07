package models

import "time"

type GiftCardStatus int

const (
	// 0等待进行
	GiftCardStatusUnused GiftCardStatus = iota
	// 1已使用
	GiftCardStatusUsed
	// 2已过期
	GiftCardStatusExpired
	// 3已作废
	GiftCardStatusInvalid
)

type GiftCard struct {
	Id         int        `json:"id"`
	D          int        `json:"d"`
	Code       string     `json:"code"`
	Status     int        `json:"status"` //  0未使用 1已使用 2已过期 3已作废
	CreatorId  int        `json:"coach_id"`
	ConsumerId int        `json:"consumer_id"`
	UsedAt     *time.Time `json:"used_at"`
	ExpiredAt  *time.Time `json:"expired_at"`
	CreatedAt  time.Time  `json:"created_at"`

	GiftCardRewardId int            `json:"gift_card_reward_id"`
	GiftCardReward   GiftCardReward `json:"gift_card_reward" gorm:"foreignKey:GiftCardRewardId"`
}

func (*GiftCard) TableName() string {
	return "GIFT_CARD"
}

type GiftCardReward struct {
	Id        int        `json:"id"`
	D         int        `json:"d"`
	Name      string     `json:"name"`
	Overview  string     `json:"overview"`
	Status    int        `json:"status"` // 1生效中 2已过期 3已作废
	Details   string     `json:"details"`
	CreatorId int        `json:"coach_id"`
	ExpiredAt *time.Time `json:"expired_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func (*GiftCardReward) TableName() string {
	return "GIFT_CARD_REWARD"
}
