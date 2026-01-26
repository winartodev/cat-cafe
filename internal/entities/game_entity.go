package entities

type GameStageStatus string

const (
	GSStatusCurrent  GameStageStatus = "current"
	GSStatusComplete GameStageStatus = "complete"
	GSStatusLocked   GameStageStatus = "locked"
)

type Game struct {
	DailyRewardAvailable bool         `json:"daily_reward_available"`
	UserBalance          *UserBalance `json:"user_balance"`
}

type UserGameStage struct {
	Slug        string          `json:"slug"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Sequence    int64           `json:"sequence"`
	Status      GameStageStatus `json:"status"`
}

type UserNextGameStageInfo struct {
	Slug   string          `json:"slug"`
	Name   string          `json:"name"`
	Status GameStageStatus `json:"status"`
}

type UpgradeKitchenStation struct {
	NewLevel       int64 `json:"new_level"`
	IsMaxLevel     bool  `json:"is_max_level"`
	NewCoinBalance int64 `json:"new_coin_balance"`
	CoinsSpent     int64 `json:"coins_spent"`

	// Current values
	CurrentProfit   int64   `json:"current_profit"`
	CurrentPrepTime float64 `json:"current_prep_time"`
	ProfitPerSecond float64 `json:"profit_per_second"`

	// Phase info
	PhaseTransitioned      bool    `json:"phase_transitioned"`
	CurrentPhase           int64   `json:"current_phase"`
	CurrentPhaseStartLevel int64   `json:"current_phase_start_level"`
	CurrentPhaseLastLevel  int64   `json:"current_phase_last_level"`
	PhaseProfitMultiplier  float64 `json:"phase_profit_multiplier"`

	// Table count
	NewTableCount int64 `json:"new_table_count,omitempty"`

	// Rewards
	PhaseRewards    []PhaseRewardInfo `json:"phase_rewards,omitempty"`
	AllPhaseRewards []PhaseRewardInfo `json:"all_phase_rewards,omitempty"`
}

type UnlockKitchenStation struct {
	UnlockedSlug   string `json:"unlocked_slug"`
	NewCoinBalance int64  `json:"new_coin_balance"`
	CoinsSpent     int64  `json:"coins_spent"`
	StationName    string `json:"station_name"`
	StationLevel   int64  `json:"station_level"`
}
