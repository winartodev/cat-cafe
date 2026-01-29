package entities

import "time"

type KitchenStation struct {
	ID         int64 `json:"id"`
	StageID    int64 `json:"stage_id"`
	FoodItemID int64 `json:"food_item_id"`
	AutoUnlock bool  `json:"auto_unlock"`

	// Additional field that didn't store into db
	FoodItemSlug  string  `json:"food_item_slug"`
	FoodName      string  `json:"food_name"`
	InitialCost   int64   `json:"initial_cost"`
	InitialProfit int64   `json:"initial_profit"`
	CookingTime   float64 `json:"cooking_time"`

	FoodItem *FoodItem `json:"food_item,omitempty"`
}

type StageKitchenConfig struct {
	ID                          int64     `json:"id"`
	StageID                     int64     `json:"stage_id"`
	MaxLevel                    int64     `json:"max_level"`
	UpgradeProfitMultiply       int64     `json:"upgrade_profit_multiply"`
	UpgradeCostMultiply         int64     `json:"upgrade_cost_multiply"`
	TransitionPhaseLevels       []int64   `json:"transition_phase_levels"`
	PhaseProfitMultipliers      []float64 `json:"phase_profit_multipliers"`
	PhaseUpgradeCostMultipliers []float64 `json:"phase_upgrade_cost_multipliers"`
	TableCountPerPhases         []int64   `json:"table_count_per_phases"`
	CreatedAt                   time.Time `json:"-"`
	UpdatedAt                   time.Time `json:"-"`
}

type KitchenPhaseCompletionRewards struct {
	ID              int64     `json:"id"`
	KitchenConfigID int64     `json:"kitchen_config_id"`
	PhaseNumber     int64     `json:"phase_number"`
	RewardID        int64     `json:"reward_id"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`

	Reward *Reward `json:"reward,omitempty"`
}

type PhaseRewardInfo struct {
	PhaseNumber int    `json:"phase_number"`
	RewardType  string `json:"reward_type"`
	RewardSlug  string `json:"reward_slug"`
	RewardName  string `json:"reward_name"`
	Amount      int64  `json:"amount"`
}
