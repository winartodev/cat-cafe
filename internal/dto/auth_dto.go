package dto

import "github.com/winartodev/cat-cafe/internal/entities"

type LoginRequest struct {
	AuthCode string `json:"auth_code"`
}

type LoginResponse struct {
	User      *UserResponse `json:"user"`
	AuthToken *string       `json:"auth_token"`
}

func ToLoginResponse(authToken *string, user *entities.User) *LoginResponse {
	return &LoginResponse{
		User:      ToUserResponse(user),
		AuthToken: authToken,
	}
}
