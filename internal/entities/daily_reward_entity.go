package entities

import (
	"time"
)

type RewardType struct {
	ID        int64          `json:"id"`
	Slug      RewardTypeSlug `json:"slug"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
}

type DailyReward struct {
	ID           int64        `json:"id"`
	DayNumber    int64        `json:"day_number"`
	RewardAmount int64        `json:"reward_amount"`
	IsActive     bool         `json:"is_active"`
	Description  string       `json:"description"`
	CreatedAt    time.Time    `json:"-"`
	UpdatedAt    time.Time    `json:"-"`
	RewardType   *RewardType  `json:"reward_type"`
	Status       RewardStatus `json:"status"`
}

type UserDailyReward struct {
	ID            int64      `json:"id"`
	UserID        int64      `json:"user_id"`
	CurrentStreak int64      `json:"current_streak"`
	LastClaimDate *time.Time `json:"last_claim_date"`
	CreatedAt     time.Time  `json:"-"`
	UpdatedAt     time.Time  `json:"-"`
}
