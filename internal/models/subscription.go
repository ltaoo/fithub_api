package models

import (
	"time"
)

// SubscriptionPlan 订阅计划
type SubscriptionPlan struct {
	Id        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Details   string    `json:"details" db:"details"`
	UnitPrice int       `json:"unit_price" db:"unit_price"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (SubscriptionPlan) TableName() string {
	return "SUBSCRIPTION_PLAN"
}

// CoachPermission 教练权限
type CoachPermission struct {
	Id      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Details string `json:"details" db:"details"`
	SortIdx int    `json:"sort_idx" db:"sort_idx"`
}

func (CoachPermission) TableName() string {
	return "COACH_PERMISSION"
}

// SubscriptionPlanCoachPermission 订阅计划和权限的关联
type SubscriptionPlanCoachPermission struct {
	Id      int `json:"id" db:"id"`
	Checked int `json:"checked" db:"checked"`

	SubscriptionPlanId int `json:"subscription_plan_id" db:"subscription_plan_id"`
	PermissionId       int `json:"permission_id" db:"permission_id"`
}

func (SubscriptionPlanCoachPermission) TableName() string {
	return "SUBSCRIPTION_PLAN_COACH_PERMISSION"
}

// DiscountPolicy 折扣政策
type DiscountPolicy struct {
	Id           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Rate         int       `json:"rate" db:"rate"`
	CountRequire int       `json:"count_require" db:"count_require"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

func (DiscountPolicy) TableName() string {
	return "DISCOUNT_POLICY"
}

// SubscriptionPlanDiscountPolicy 订阅计划和折扣政策的关联
type SubscriptionPlanDiscountPolicy struct {
	Id      int `json:"id" db:"id"`
	Enabled int `json:"enabled" db:"enabled"`

	SubscriptionPlanId int              `json:"subscription_plan_id" db:"subscription_plan_id"`
	SubscriptionPlan   SubscriptionPlan `json:"subscription_plan"`
	DiscountPolicyId   int              `json:"discount_policy_id"`
	DiscountPolicy     DiscountPolicy   `json:"discount_policy"`
}

func (SubscriptionPlanDiscountPolicy) TableName() string {
	return "SUBSCRIPTION_PLAN_DISCOUNT_POLICY"
}

// SubscriptionOrder 订阅订单
type SubscriptionOrder struct {
	Id              int    `json:"id" db:"id"`
	Amount          int    `json:"amount" db:"amount"`
	Discount        int    `json:"discount" db:"discount"`
	DiscountDetails string `json:"discount_details" db:"discount_details"`

	SubscriptionPlanId int `json:"subscription_plan_id" db:"subscription_plan_id"`
	InvoiceId          int `json:"invoice_id" db:"invoice_id"`
	CoachId            int `json:"coach_id" db:"coach_id"`
}

func (SubscriptionOrder) TableName() string {
	return "SUBSCRIPTION_ORDER"
}

// Invoice 账单
type Invoice struct {
	Id           int        `json:"id" db:"id"`
	Status       int        `json:"status" db:"status"`
	Amount       int        `json:"amount" db:"amount"`
	CurrencyUnit int        `json:"currency_unit" db:"currency_unit"`
	OrderType    int        `json:"order_type" db:"order_type"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	PaidAt       *time.Time `json:"paid_at" db:"paid_at"`
	CanceledAt   *time.Time `json:"canceled_at" db:"canceled_at"`
	CancelReason string     `json:"cancel_reason" db:"cancel_reason"`

	OrderId int `json:"order_id" db:"order_id"`
	CoachId int `json:"coach_id" db:"coach_id"`
}

func (Invoice) TableName() string {
	return "INVOICE"
}

// Subscription 订阅
type Subscription struct {
	Id              int        `json:"id" db:"id"`
	Step            int        `json:"step" db:"step"`
	Count           int        `json:"count" db:"count"`
	Reason          string     `json:"reason"`
	ExpectExpiredAt *time.Time `json:"expect_expired_at" db:"expect_expired_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	ActiveAt        *time.Time `json:"active_at" db:"active_at"`
	ExpiredAt       *time.Time `json:"expired_at" db:"expired_at"`
	PausedAt        *time.Time `json:"paused_at" db:"paused_at"`
	PauseCount      int        `json:"pause_count" db:"pause_count"`
	InvalidAt       *time.Time `json:"invalid_at" db:"invalid_at"`

	CoachId            int              `json:"coach_id" db:"coach_id"`
	Coach              Coach            `json:"coach"`
	SubscriptionPlanId int              `json:"subscription_plan_id" db:"subscription_plan_id"`
	SubscriptionPlan   SubscriptionPlan `json:"subscription_plan"`
}

func (Subscription) TableName() string {
	return "SUBSCRIPTION"
}
