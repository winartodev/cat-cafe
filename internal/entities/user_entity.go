package entities

import (
	"time"
)

type User struct {
	ID           int64        `json:"id"`
	ExternalID   string       `json:"external_id"`
	Username     string       `json:"username"`
	Email        string       `json:"email"`
	PasswordHash string       `json:"-"`
	IsActive     bool         `json:"is_active"`
	UserBalance  *UserBalance `json:"balance"`
	CreatedAt    time.Time    `json:"-"`
	UpdatedAt    time.Time    `json:"-"`
}

type UserBalance struct {
	Coin int64 `json:"coin"`
	Gem  int64 `json:"gem"`
}

type UserCache struct {
	UserID     int64             `json:"user_id"`
	ExternalID string            `json:"external_id"`
	Username   string            `json:"username"`
	Email      string            `json:"email"`
	IsActive   bool              `json:"is_active"`
	Balance    *UserBalanceCache `json:"balance,omitempty"`
}

type UserBalanceCache struct {
	Coin int64 `json:"coin"`
	Gem  int64 `json:"gem"`
}

type UserGameStageProgression struct {
	ID          int64           `json:"id"`
	UserID      int64           `json:"user_id"`
	StageID     int64           `json:"stage_id"`
	Status      GameStageStatus `json:"status"`
	IsComplete  bool            `json:"is_complete"`
	CompletedAt *time.Time      `json:"completed_at"`
}

type UserKitchenStageProgression struct {
	ID               int64                       `json:"id"`
	UserID           int64                       `json:"user_id"`
	StageID          int64                       `json:"stage_id"`
	StationLevels    map[string]UserStationLevel `json:"station_levels"`
	UnlockedStations []string                    `json:"unlocked_stations"`
	NextLevelStats   map[string]UserStationLevel `json:"next_level_stats,omitempty"` // Calculated, not stored in DB
}

type UserKitchenPhaseProgression struct {
	ID              int64   `json:"id"`
	UserID          int64   `json:"user_id"`
	KitchenConfigID int64   `json:"kitchen_config_id"`
	CurrentPhase    int64   `json:"current_phase"`
	CompletedPhases []int64 `json:"completed_phases"`
}

type UserKitchenPhaseRewardClaim struct {
	UserID          int64      `json:"user_id"`
	KitchenConfigID int64      `json:"kitchen_config_id"`
	CurrentPhase    int64      `json:"current_phase"`
	RewardID        int64      `json:"reward_id"`
	ClaimedAt       *time.Time `json:"claimed_at"`
}

type UserStationLevel struct {
	Level           int64   `json:"level"`
	Cost            int64   `json:"cost"`
	Profit          int64   `json:"profit"`
	PreparationTime float64 `json:"preparation_time"`

	Reward *Reward `json:"reward,omitempty"`
}

func (u *User) ToCache() *UserCache {
	cache := &UserCache{
		UserID:     u.ID,
		ExternalID: u.ExternalID,
		Username:   u.Username,
		Email:      u.Email,
		IsActive:   u.IsActive,
	}

	if u.UserBalance != nil {
		cache.Balance = &UserBalanceCache{
			Coin: u.UserBalance.Coin,
			Gem:  u.UserBalance.Gem,
		}
	}

	return cache
}

func UserFromCache(cache *UserCache) *User {
	user := &User{
		ID:         cache.UserID,
		ExternalID: cache.ExternalID,
		Username:   cache.Username,
		Email:      cache.Email,
		IsActive:   cache.IsActive,
	}

	if cache.Balance != nil {
		user.UserBalance = &UserBalance{
			Coin: cache.Balance.Coin,
			Gem:  cache.Balance.Gem,
		}
	}

	return user
}
