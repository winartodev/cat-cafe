package entities

import "time"

type GameStage struct {
	ID           int64     `json:"id"`
	Slug         string    `json:"slug"`
	Name         string    `json:"name"`
	StartingCoin int64     `json:"starting_coin"`
	StagePrize   int64     `json:"stage_prize"`
	IsActive     bool      `json:"is_active"`
	Sequence     int64     `json:"sequence"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type StageCustomerConfig struct {
	ID                      int64     `json:"id"`
	StageID                 int64     `json:"stage_id"`
	CustomerSpawnTime       float64   `json:"customer_spawn_time"`
	MaxCustomerOrderCount   int64     `json:"max_customer_order_count"`
	MaxCustomerOrderVariant int64     `json:"max_customer_order_variant"`
	StartingOrderTableCount int64     `json:"starting_order_table_count"`
	CreatedAt               time.Time `json:"-"`
	UpdatedAt               time.Time `json:"-"`
}

type StageStaffConfig struct {
	ID                   int64     `json:"id"`
	StageID              int64     `json:"stage_id"`
	StartingStaffManager string    `json:"starting_staff_manager"`
	StartingStaffHelper  string    `json:"starting_staff_helper"`
	CreatedAt            time.Time `json:"-"`
	UpdatedAt            time.Time `json:"-"`
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

type StageCameraConfig struct {
	ID        int64     `json:"id"`
	StageID   int64     `json:"stage_id"`
	ZoomSize  float64   `json:"zoom_size"`
	MinBoundX float64   `json:"min_bound_x"`
	MinBoundY float64   `json:"min_bound_y"`
	MaxBoundX float64   `json:"max_bound_x"`
	MaxBoundY float64   `json:"max_bound_y"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type KitchenPhaseCompletionRewards struct {
	ID              int64     `json:"id"`
	KitchenConfigID int64     `json:"kitchen_config_id"`
	PhaseNumber     int64     `json:"phase_number"`
	RewardID        int64     `json:"reward_id"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`

	// Additional field that didn't store into db
	RewardSlug string `json:"reward_slug"`
}

type GameStageConfig struct {
	CustomerConfig     *StageCustomerConfig
	StaffConfig        *StageStaffConfig
	KitchenConfig      *StageKitchenConfig
	CameraConfig       *StageCameraConfig
	KitchenPhaseReward []KitchenPhaseCompletionRewards
}
