package entities

import (
	"time"
)

type RewardType struct {
	ID        int64
	Slug      RewardTypeSlug
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DailyReward struct {
	ID           int64
	DayNumber    int64
	RewardAmount int64
	IsActive     bool
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	RewardType   *RewardType
	Status       RewardStatus
}

type UserDailyReward struct {
	ID            int64
	UserID        int64
	CurrentStreak int64
	LastClaimDate *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
