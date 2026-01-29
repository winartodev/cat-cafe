package dto

import (
	"time"

	"github.com/winartodev/cat-cafe/internal/entities"
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
	Slug        string                   `json:"slug"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Sequence    int64                    `json:"sequence"`
	Status      entities.GameStageStatus `json:"status"`
}

type UserGameStageResponse struct {
	CurrentStageIdx int             `json:"current_stage_idx"`
	Stages          []UserGameStage `json:"stages"`
}

type UserDetailGameStageResponse struct {
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	Description  string `json:"string"`
	StartingCoin int64  `json:"starting_coin"`
	StagePrize   int64  `json:"stage_prize"`
	IsActive     bool   `json:"is_active"`
	Sequence     int64  `json:"sequence"`

	Customer        *CustomerConfigDTO  `json:"customer_config,omitempty"`
	Staff           *StaffConfigDTO     `json:"staff_config,omitempty"`
	KitchenStations []KitchenStationDTO `json:"kitchen_stations,omitempty"`
	Kitchen         *KitchenConfigDTO   `json:"kitchen_config,omitempty"`
	Camera          *CameraConfigDTO    `json:"camera_config,omitempty"`

	NextStage *entities.UserNextGameStageInfo `json:"next_stage,omitempty"`
}

type UserUpgradeKitchenResponse struct {
	Name           string               `json:"name"`
	Slug           string               `json:"slug"`
	CurrentLevel   *currentStationLevel `json:"current_level,omitempty"`
	NextLevel      *nextStationLevel    `json:"next_level,omitempty"`
	GrantedRewards []kitchenPhaseReward `json:"granted_rewards,omitempty"`
}

type UserUnlockKitchenResponse struct {
	Name         string               `json:"name"`
	Slug         string               `json:"slug"`
	CurrentLevel *currentStationLevel `json:"current_level,omitempty"`
	NextLevel    *nextStationLevel    `json:"next_level,omitempty"`
}

type currentStationLevel struct {
	Level          int64               `json:"level"`
	Profit         int64               `json:"profit"`
	CookingTime    float64             `json:"cooking_time"`
	CompletedPhase int                 `json:"completed_phase"`
	Reward         *kitchenPhaseReward `json:"reward,omitempty"`
}

type nextStationLevel struct {
	Level  int64 `json:"level"`
	Cost   int64 `json:"cost"`
	Profit int64 `json:"profit"`
}

type kitchenPhaseReward struct {
	RewardType   string `json:"reward_type"`
	RewardName   string `json:"reward_name"`
	RewardAmount int64  `json:"reward_amount"`
}

func ToUserUpgradeKitchenResponse(data *entities.UpgradeKitchenStation) *UserUpgradeKitchenResponse {
	var grantedRewards []kitchenPhaseReward
	for _, v := range data.GrantedRewards {
		grantedRewards = append(grantedRewards, kitchenPhaseReward{
			RewardType:   v.RewardType,
			RewardName:   v.RewardName,
			RewardAmount: v.Amount,
		})
	}

	var nextLevel *nextStationLevel
	var rewards *kitchenPhaseReward
	if data.CurrentRewards != nil {
		rewards = &kitchenPhaseReward{
			RewardType:   data.CurrentRewards.RewardType,
			RewardName:   data.CurrentRewards.RewardName,
			RewardAmount: data.CurrentRewards.Amount,
		}
	}

	if !data.IsMaxLevel {
		nextLevel = &nextStationLevel{
			Level:  data.NextLevel,
			Cost:   data.NextCost,
			Profit: data.NextProfit,
		}
	}

	return &UserUpgradeKitchenResponse{
		Name: data.Name,
		Slug: data.Slug,
		CurrentLevel: &currentStationLevel{
			Level:          data.CurrentLevel,
			Profit:         data.CurrentProfit,
			CookingTime:    data.CurrentPrepTime,
			CompletedPhase: data.CompletedPhase,
			Reward:         rewards,
		},
		NextLevel:      nextLevel,
		GrantedRewards: grantedRewards,
	}
}

func ToUserUnlockKitchenResponse(data *entities.UnlockKitchenStation) *UserUnlockKitchenResponse {
	var rewards *kitchenPhaseReward
	if data.CurrentRewards != nil {
		rewards = &kitchenPhaseReward{
			RewardType:   data.CurrentRewards.RewardType,
			RewardName:   data.CurrentRewards.RewardName,
			RewardAmount: data.CurrentRewards.Amount,
		}
	}
	return &UserUnlockKitchenResponse{
		Name: data.Name,
		Slug: data.Slug,
		CurrentLevel: &currentStationLevel{
			Level:       data.CurrentLevel,
			Profit:      data.CurrentProfit,
			CookingTime: data.CurrentPrepTime,
			Reward:      rewards,
		},
		NextLevel: &nextStationLevel{
			Level:  data.NextLevel,
			Cost:   data.NextCost,
			Profit: data.NextProfit,
		},
	}
}

func ToUserGameStageResponse(data *entities.UserGameStage) *UserGameStage {
	return &UserGameStage{
		Slug:        data.Slug,
		Name:        data.Name,
		Sequence:    data.Sequence,
		Status:      data.Status,
		Description: data.Description,
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
		Slug:            data.Slug,
		Name:            data.Name,
		StartingCoin:    data.StartingCoin,
		StagePrize:      data.StagePrize,
		Description:     data.Description,
		IsActive:        data.IsActive,
		Sequence:        data.Sequence,
		Customer:        toCustomerConfigDTO(config.CustomerConfig),
		Staff:           toStaffConfigDTO(config.StaffConfig),
		KitchenStations: toKitchenStationsDTOWithProgress(config.KitchenStations, config.UserProgress),
		Kitchen:         toKitchenConfigDTO(config.KitchenConfig, config.KitchenPhaseReward),
		Camera:          toCameraConfigDTO(config.CameraConfig),
		NextStage:       nextStage,
	}
}
