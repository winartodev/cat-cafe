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
