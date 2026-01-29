package dto

import (
	"github.com/winartodev/cat-cafe/internal/entities"
)

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
	FoodItemSlug string `json:"slug"`
	FoodName     string `json:"name"`
	AutoUnlock   bool   `json:"auto_unlock"`
	IsLocked     bool   `json:"is_locked"`

	// User progression data
	CurrentLevel *currentStationLevel `json:"current_level,omitempty"`
	NextLevel    *nextStationLevel    `json:"next_level,omitempty"`
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
		FoodItemSlug: data.FoodItemSlug,
		FoodName:     data.FoodName,
		AutoUnlock:   data.AutoUnlock,
		IsLocked:     true,
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

func toKitchenStationsDTOWithProgress(
	stations []entities.KitchenStation,
	userProgress *entities.UserKitchenStageProgression,
) []KitchenStationDTO {
	if len(stations) == 0 {
		return nil
	}

	kitchenStations := make([]KitchenStationDTO, 0)
	for _, station := range stations {
		dto := toKitchenStationDTO(&station)

		if userProgress == nil || len(userProgress.StationLevels) == 0 {
			continue
		}

		// Add user progression data if available
		if stationLevel, exists := userProgress.StationLevels[station.FoodItemSlug]; exists {
			if stationLevel.Level == 0 {
				dto.IsLocked = true
			} else {
				dto.IsLocked = false
			}

			// Current level data
			if stationLevel.Level > 0 {
				dto.CurrentLevel = &currentStationLevel{
					Level:       stationLevel.Level,
					Profit:      stationLevel.Profit,
					CookingTime: stationLevel.PreparationTime,
					Reward: &kitchenPhaseReward{
						RewardName:   stationLevel.Reward.Name,
						RewardType:   stationLevel.Reward.RewardType.Slug,
						RewardAmount: stationLevel.Reward.Amount,
					},
				}
			}

			// Add next level data if available
			if userProgress.NextLevelStats != nil {
				if nextLevel, exists := userProgress.NextLevelStats[station.FoodItemSlug]; exists {
					dto.NextLevel = &nextStationLevel{
						Level:  nextLevel.Level,
						Cost:   nextLevel.Cost,
						Profit: nextLevel.Profit,
					}
				}
			}
		}

		kitchenStations = append(kitchenStations, *dto)
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
		if reward.Reward == nil || reward.Reward.RewardType == nil {
			continue
		}

		kitchenPhaseRewards = append(kitchenPhaseRewards, KitchenPhaseCompletionRewardDTO{
			PhaseNumber: reward.PhaseNumber,
			Reward:      reward.Reward.Slug,
			RewardType:  reward.Reward.RewardType.Slug,
		})
	}

	return kitchenPhaseRewards
}
