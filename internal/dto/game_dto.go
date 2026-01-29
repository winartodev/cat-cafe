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

type UserStageUpgradeResponse struct {
	Slug        string                   `json:"slug"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Cost        int64                    `json:"cost"`
	CostType    entities.UpgradeCostType `json:"cost_type"`

	IsPurchased bool `json:"is_purchased"`
}

type UserPurchasedStageUpgradeResponse struct {
	UpgradeEffectDTO `json:"upgrade_effect"`
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
		if v.Status == entities.GSStatusAvailable {
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
	if data == nil || config == nil {
		return nil
	}

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

func ToUserStageUpgradesResponse(items []entities.UserStageUpgrade) []UserStageUpgradeResponse {
	if items == nil || len(items) == 0 {
		return nil
	}

	var data []UserStageUpgradeResponse
	for _, item := range items {
		upgrade := item.Upgrade
		data = append(data, UserStageUpgradeResponse{
			Slug:        upgrade.Slug,
			Name:        upgrade.Name,
			Description: upgrade.Description,
			Cost:        upgrade.Cost,
			CostType:    upgrade.CostType,
			IsPurchased: item.IsPurchased,
		})
	}

	return data
}

func ToUserPurchasedStageUpgradeResponse(data *entities.Upgrade) *UserPurchasedStageUpgradeResponse {
	if data == nil {
		return nil
	}

	return &UserPurchasedStageUpgradeResponse{
		UpgradeEffectDTO{
			Type:       data.Effect.Type,
			Value:      data.Effect.Value,
			Unit:       data.Effect.Unit,
			Target:     data.Effect.Target,
			TargetName: data.Effect.TargetName,
		},
	}
}
