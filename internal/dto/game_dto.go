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

type UserGameStageResponse struct {
	CurrentStageIdx int             `json:"current_stage_idx"`
	Stages          []UserGameStage `json:"stages"`
}

type UserDetailGameStageResponse struct {
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	StartingCoin int64  `json:"starting_coin"`
	StagePrize   int64  `json:"stage_prize"`
	IsActive     bool   `json:"is_active"`
	Sequence     int64  `json:"sequence"`

	Customer *CustomerConfigDTO `json:"customer_config,omitempty"`
	Staff    *StaffConfigDTO    `json:"staff_config,omitempty"`
	Kitchen  *KitchenConfigDTO  `json:"kitchen_config,omitempty"`
	Camera   *CameraConfigDTO   `json:"camera_config,omitempty"`

	NextStage *entities.UserNextGameStageInfo `json:"next_stage,omitempty"`
}

func ToUserGameStageResponse(data *entities.UserGameStage) *UserGameStage {
	return &UserGameStage{
		Slug:     data.Slug,
		Name:     data.Name,
		Sequence: data.Sequence,
		Status:   data.Status,
	}
}

func ToUserGameStageResponses(data []entities.UserGameStage) *UserGameStageResponse {
	var stages []UserGameStage
	var currentStageIdx int
	for i, v := range data {
		if v.Status == entities.GSStatusCurrent {
			currentStageIdx = i
		}
		stages = append(stages, *ToUserGameStageResponse(&v))
	}

	return &UserGameStageResponse{
		CurrentStageIdx: currentStageIdx,
		Stages:          stages,
	}
}

func ToUserDetailGameStageResponse(
	data *entities.GameStage,
	config *entities.GameStageConfig,
	nextStage *entities.UserNextGameStageInfo,
) *UserDetailGameStageResponse {
	return &UserDetailGameStageResponse{
		Slug:         data.Slug,
		Name:         data.Name,
		StartingCoin: data.StartingCoin,
		StagePrize:   data.StagePrize,
		IsActive:     data.IsActive,
		Sequence:     data.Sequence,
		Customer:     toCustomerConfigDTO(config.CustomerConfig),
		Staff:        toStaffConfigDTO(config.StaffConfig),
		Kitchen:      toKitchenConfigDTO(config.KitchenConfig, config.KitchenPhaseReward),
		Camera:       toCameraConfigDTO(config.CameraConfig),
		NextStage:    nextStage,
	}
}
