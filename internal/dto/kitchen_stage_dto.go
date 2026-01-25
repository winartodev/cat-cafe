package dto

import "github.com/winartodev/cat-cafe/internal/entities"

type KitchenConfigRequest struct {
	MaxLevel                    int64                `json:"max_level"`
	UpgradeProfitMultiply       int64                `json:"upgrade_profit_multiply"`
	UpgradeCostMultiply         int64                `json:"upgrade_cost_multiply"`
	TransitionPhaseLevels       []int64              `json:"transition_phase_levels"`
	PhaseProfitMultipliers      []float64            `json:"phase_profit_multipliers"`
	PhaseUpgradeCostMultipliers []float64            `json:"phase_upgrade_cost_multipliers"`
	TableCountPerPhases         []int64              `json:"table_count_per_phases"`
	PhaseRewards                []PhaseRewardRequest `json:"phase_rewards"`
}

type KitchenStationDTO struct {
	FoodItemSlug        string  `json:"food_item_slug"`
	FoodName            string  `json:"food_name"`
	AutoUnlock          bool    `json:"auto_unlock"`
	StartingPrice       int64   `json:"starting_price"`
	StartingPreparation float64 `json:"starting_preparation"`
}

type KitchenConfigDTO struct {
	MaxLevel                    int64     `json:"max_level"`
	UpgradeProfitMultiply       int64     `json:"upgrade_profit_multiply"`
	UpgradeCostMultiply         int64     `json:"upgrade_cost_multiply"`
	TransitionPhaseLevels       []int64   `json:"transition_phase_levels"`
	PhaseProfitMultipliers      []float64 `json:"phase_profit_multipliers"`
	PhaseUpgradeCostMultipliers []float64 `json:"phase_upgrade_cost_multipliers"`
	TableCountPerPhases         []int64   `json:"table_count_per_phases"`

	PhaseRewards []KitchenPhaseCompletionRewardDTO `json:"phase_rewards,omitempty"`
}

type KitchenPhaseCompletionRewardDTO struct {
	PhaseNumber int64  `json:"phase_number"`
	Reward      string `json:"reward"`
	RewardType  string `json:"reward_type"`
}

func toKitchenStationDTO(data *entities.KitchenStation) *KitchenStationDTO {
	if data == nil {
		return nil
	}

	return &KitchenStationDTO{
		FoodItemSlug:        data.FoodItemSlug,
		FoodName:            data.FoodName,
		AutoUnlock:          data.AutoUnlock,
		StartingPrice:       data.StartingPrice,
		StartingPreparation: data.StartingPreparation,
	}
}

func toKitchenStationsDTO(data []entities.KitchenStation) []KitchenStationDTO {
	if len(data) == 0 {
		return nil
	}

	kitchenStations := make([]KitchenStationDTO, 0)
	for _, item := range data {
		kitchenStations = append(kitchenStations, *toKitchenStationDTO(&item))
	}

	return kitchenStations
}

func toKitchenConfigDTO(data *entities.StageKitchenConfig, kitchenPhaseReward []entities.KitchenPhaseCompletionRewards) *KitchenConfigDTO {
	if data == nil {
		return nil
	}
	return &KitchenConfigDTO{
		MaxLevel:                    data.MaxLevel,
		UpgradeProfitMultiply:       data.UpgradeProfitMultiply,
		UpgradeCostMultiply:         data.UpgradeCostMultiply,
		TransitionPhaseLevels:       data.TransitionPhaseLevels,
		PhaseProfitMultipliers:      data.PhaseProfitMultipliers,
		PhaseUpgradeCostMultipliers: data.PhaseUpgradeCostMultipliers,
		TableCountPerPhases:         data.TableCountPerPhases,
		PhaseRewards:                toKitchenPhaseRewards(kitchenPhaseReward),
	}
}

func toKitchenPhaseRewards(rewards []entities.KitchenPhaseCompletionRewards) []KitchenPhaseCompletionRewardDTO {
	if len(rewards) == 0 {
		return nil
	}
	kitchenPhaseRewards := make([]KitchenPhaseCompletionRewardDTO, 0, len(rewards))
	for _, reward := range rewards {
		kitchenPhaseRewards = append(kitchenPhaseRewards, KitchenPhaseCompletionRewardDTO{
			PhaseNumber: reward.PhaseNumber,
			Reward:      reward.RewardSlug,
			RewardType:  reward.RewardType,
		})
	}

	return kitchenPhaseRewards
}
