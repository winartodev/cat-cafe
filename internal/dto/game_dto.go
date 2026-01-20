package dto

import (
	"github.com/winartodev/cat-cafe/internal/entities"
	"time"
)

type SyncBalanceRequest struct {
	CoinsEarned  int64     `json:"coins_earned"`
	LastSyncTime time.Time `json:"last_sync_time"`
}

type SyncBalanceResponse struct {
	CurrentCoinBalance int64 `json:"current_coin_balance"`
	CurrentGemBalance  int64 `json:"current_gem_balance"`
}

type UserGameStage struct {
	Slug     string                   `json:"slug"`
	Name     string                   `json:"name"`
	Sequence int64                    `json:"sequence"`
	Status   entities.GameStageStatus `json:"status"`
}

func ToUserGameStageResponse(data *entities.UserGameStage) *UserGameStage {
	return &UserGameStage{
		Slug:     data.Slug,
		Name:     data.Name,
		Sequence: data.Sequence,
		Status:   data.Status,
	}
}

func ToUserGameStageResponses(data []entities.UserGameStage) []UserGameStage {
	var stages []UserGameStage
	for _, v := range data {
		stages = append(stages, *ToUserGameStageResponse(&v))
	}

	return stages
}
