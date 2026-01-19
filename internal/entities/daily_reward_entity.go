package entities

import (
	"time"
)

type RewardType struct {
	ID        *int64    `json:"id,omitempty"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Reward struct {
	ID         int64       `json:"id,omitempty"`
	Slug       string      `json:"slug"`
	Name       string      `json:"name"`
	Amount     int64       `json:"amount"`
	IsActive   bool        `json:"is_active"`
	CreatedAt  time.Time   `json:"-"`
	UpdatedAt  time.Time   `json:"-"`
	RewardType *RewardType `json:"reward_type"`
}

type DailyReward struct {
	ID          int64        `json:"id"`
	DayNumber   int64        `json:"day_number"`
	IsActive    bool         `json:"is_active"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"-"`
	UpdatedAt   time.Time    `json:"-"`
	Reward      *Reward      `json:"reward"`
	Status      RewardStatus `json:"status"`
}

type UserDailyReward struct {
	ID            int64      `json:"id"`
	UserID        int64      `json:"user_id"`
	LongestStreak int64      `json:"longest_streak"`
	CurrentStreak int64      `json:"current_streak"`
	LastClaimDate *time.Time `json:"last_claim_date"`
	CreatedAt     time.Time  `json:"-"`
	UpdatedAt     time.Time  `json:"-"`
}
