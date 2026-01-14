package entities

import "time"

type User struct {
	ID           int64
	ExternalID   string
	Username     string
	Email        string
	PasswordHash string
	IsActive     bool
	UserBalance  *UserBalance
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserBalance struct {
	Coin int64
	Gem  int64
}
