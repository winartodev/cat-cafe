package dto

import "github.com/winartodev/cat-cafe/internal/entities"

type UserBalanceResponse struct {
	Coin int64 `json:"coin"`
	Gem  int64 `json:"gem"`
}

type UserResponse struct {
	ID          int64                `json:"id"`
	ExternalID  string               `json:"external_id"`
	Username    string               `json:"username"`
	Email       string               `json:"email"`
	IsActive    bool                 `json:"is_active"`
	UserBalance *UserBalanceResponse `json:"user_balance"`
}

func ToUserResponse(user *entities.User) *UserResponse {
	var userBalance *UserBalanceResponse

	if user.UserBalance != nil {
		userBalance = &UserBalanceResponse{
			Coin: user.UserBalance.Coin,
			Gem:  user.UserBalance.Gem,
		}
	}

	return &UserResponse{
		ID:          user.ID,
		ExternalID:  user.ExternalID,
		Username:    user.Username,
		Email:       user.Email,
		IsActive:    user.IsActive,
		UserBalance: userBalance,
	}
}
