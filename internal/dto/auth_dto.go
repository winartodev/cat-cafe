package dto

import "github.com/winartodev/cat-cafe/internal/entities"

type LoginRequest struct {
	AuthCode string `json:"auth_code"`
}

type LoginResponse struct {
	User      *UserResponse  `json:"user"`
	GameData  *entities.Game `json:"game_data,omitempty"`
	AuthToken *string        `json:"auth_token"`
}

func ToLoginResponse(authToken *string, user *entities.User, gameData *entities.Game) *LoginResponse {
	return &LoginResponse{
		User:      ToUserResponse(user),
		GameData:  gameData,
		AuthToken: authToken,
	}
}
