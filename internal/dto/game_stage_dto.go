package dto

import (
	"errors"
	"github.com/winartodev/cat-cafe/internal/entities"
)

type BaseGameStageRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	StartingCoin int64  `json:"starting_coin"`
	StagePrize   int64  `json:"stage_prize"`
	IsActive     bool   `json:"is_active"`
	Sequence     int64  `json:"sequence"`

	Customer        *CustomerConfigDTO    `json:"customer_config"`
	Staff           *StaffConfigDTO       `json:"staff_config"`
	KitchenStations []KitchenStationDTO   `json:"kitchen_stations"`
	KitchenConfig   *KitchenConfigRequest `json:"kitchen_config"`
	Camera          *CameraConfigDTO      `json:"camera_config"`
}

type CreateGameStageRequest struct {
	Slug string `json:"slug"`
	BaseGameStageRequest
}

type UpdateGameStageRequest struct {
	BaseGameStageRequest
}

type UpdateGameStageResponse struct {
	BaseGameStageRequest
}

type PhaseRewardRequest struct {
	PhaseNumber int64    `json:"phase_number"`
	RewardSlugs []string `json:"reward_slugs"` // Array of slugs
}

type GameStageDetailResponse struct {
	ID           int64  `json:"id"`
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	StartingCoin int64  `json:"starting_coin"`
	StagePrize   int64  `json:"stage_prize"`
	IsActive     bool   `json:"is_active"`
	Sequence     int64  `json:"sequence"`

	Customer       *CustomerConfigDTO  `json:"customer_config,omitempty"`
	Staff          *StaffConfigDTO     `json:"staff_config,omitempty"`
	KitchenStation []KitchenStationDTO `json:"kitchen_station"`
	KitchenConfig  *KitchenConfigDTO   `json:"kitchen_config,omitempty"`
	Camera         *CameraConfigDTO    `json:"camera_config,omitempty"`
}

type GameStageResponse struct {
	ID           int64  `json:"id"`
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	StartingCoin int64  `json:"starting_coin"`
	StagePrize   int64  `json:"stage_prize"`
	IsActive     bool   `json:"is_active"`
	Sequence     int64  `json:"sequence"`
}

type CustomerConfigDTO struct {
	CustomerSpawnTime       float64 `json:"customer_spawn_time"`
	MaxCustomerOrderCount   int64   `json:"max_customer_order_count"`
	MaxCustomerOrderVariant int64   `json:"max_customer_order_variant"`
	StartingOrderTableCount int64   `json:"starting_order_table_count"`
}

type StaffConfigDTO struct {
	StartingStaffManager string `json:"starting_staff_manager"`
	StartingStaffHelper  string `json:"starting_staff_helper"`
}

type CameraConfigDTO struct {
	ZoomSize  float64 `json:"zoom_size"`
	MinBoundX float64 `json:"min_bound_x"`
	MinBoundY float64 `json:"min_bound_y"`
	MaxBoundX float64 `json:"max_bound_x"`
	MaxBoundY float64 `json:"max_bound_y"`
}

func (d *BaseGameStageRequest) ValidateConfig() error {
	if d.Customer == nil || d.Staff == nil || d.KitchenConfig == nil || d.Camera == nil || d.KitchenStations == nil {
		return errors.New("invalid request")
	}

	return nil
}

func (d *CreateGameStageRequest) ValidateConfig() error {
	return d.BaseGameStageRequest.ValidateConfig()
}

func (d *UpdateGameStageRequest) ValidateConfig() error {
	return d.BaseGameStageRequest.ValidateConfig()
}

func toCustomerConfigDTO(data *entities.StageCustomerConfig) *CustomerConfigDTO {
	if data == nil {
		return nil
	}

	return &CustomerConfigDTO{
		CustomerSpawnTime:       data.CustomerSpawnTime,
		MaxCustomerOrderCount:   data.MaxCustomerOrderCount,
		MaxCustomerOrderVariant: data.MaxCustomerOrderVariant,
		StartingOrderTableCount: data.StartingOrderTableCount,
	}
}

func toStaffConfigDTO(data *entities.StageStaffConfig) *StaffConfigDTO {
	if data == nil {
		return nil
	}

	return &StaffConfigDTO{
		StartingStaffManager: data.StartingStaffManager,
		StartingStaffHelper:  data.StartingStaffHelper,
	}
}

func toCameraConfigDTO(data *entities.StageCameraConfig) *CameraConfigDTO {
	return &CameraConfigDTO{
		ZoomSize:  data.ZoomSize,
		MinBoundX: data.MinBoundX,
		MaxBoundY: data.MaxBoundY,
		MaxBoundX: data.MaxBoundX,
		MinBoundY: data.MinBoundY,
	}
}

func ToGameStageResponse(data entities.GameStage) GameStageResponse {
	return GameStageResponse{
		ID:           data.ID,
		Slug:         data.Slug,
		Name:         data.Name,
		Description:  data.Description,
		StartingCoin: data.StartingCoin,
		StagePrize:   data.StagePrize,
		IsActive:     data.IsActive,
		Sequence:     data.Sequence,
	}
}

func ToGameStageDetailResponse(data *entities.GameStage, gameConfig *entities.GameStageConfig) *GameStageDetailResponse {
	if data == nil {
		return nil
	}

	return &GameStageDetailResponse{
		ID:             data.ID,
		Slug:           data.Slug,
		Name:           data.Name,
		Description:    data.Description,
		StartingCoin:   data.StartingCoin,
		StagePrize:     data.StagePrize,
		IsActive:       data.IsActive,
		Sequence:       data.Sequence,
		Customer:       toCustomerConfigDTO(gameConfig.CustomerConfig),
		Staff:          toStaffConfigDTO(gameConfig.StaffConfig),
		KitchenStation: toKitchenStationsDTO(gameConfig.KitchenStations),
		KitchenConfig:  toKitchenConfigDTO(gameConfig.KitchenConfig, gameConfig.KitchenPhaseReward),
		Camera:         toCameraConfigDTO(gameConfig.CameraConfig),
	}
}

func ToGameStageResponses(stages []entities.GameStage) []GameStageResponse {
	if len(stages) == 0 {
		return nil
	}
	responses := make([]GameStageResponse, 0, len(stages))
	for _, stage := range stages {
		responses = append(responses, ToGameStageResponse(stage))
	}

	return responses
}

func (d *BaseGameStageRequest) toEntitiesCommon() (*entities.GameStageConfig, error) {
	if err := d.ValidateConfig(); err != nil {
		return nil, err
	}

	var gameStageConfig entities.GameStageConfig

	gameStageConfig.CustomerConfig = &entities.StageCustomerConfig{
		CustomerSpawnTime:       d.Customer.CustomerSpawnTime,
		MaxCustomerOrderCount:   d.Customer.MaxCustomerOrderCount,
		MaxCustomerOrderVariant: d.Customer.MaxCustomerOrderVariant,
		StartingOrderTableCount: d.Customer.StartingOrderTableCount,
	}

	gameStageConfig.StaffConfig = &entities.StageStaffConfig{
		StartingStaffManager: d.Staff.StartingStaffManager,
		StartingStaffHelper:  d.Staff.StartingStaffHelper,
	}

	gameStageConfig.KitchenConfig = &entities.StageKitchenConfig{
		MaxLevel:                    d.KitchenConfig.MaxLevel,
		UpgradeProfitMultiply:       d.KitchenConfig.UpgradeProfitMultiply,
		UpgradeCostMultiply:         d.KitchenConfig.UpgradeCostMultiply,
		TransitionPhaseLevels:       d.KitchenConfig.TransitionPhaseLevels,
		PhaseProfitMultipliers:      d.KitchenConfig.PhaseProfitMultipliers,
		PhaseUpgradeCostMultipliers: d.KitchenConfig.PhaseUpgradeCostMultipliers,
		TableCountPerPhases:         d.KitchenConfig.TableCountPerPhases,
	}

	gameStageConfig.CameraConfig = &entities.StageCameraConfig{
		ZoomSize:  d.Camera.ZoomSize,
		MinBoundX: d.Camera.MinBoundX,
		MinBoundY: d.Camera.MinBoundY,
		MaxBoundX: d.Camera.MaxBoundX,
		MaxBoundY: d.Camera.MaxBoundY,
	}

	var kitchenPhaseRewards []entities.KitchenPhaseCompletionRewards
	for _, phaseData := range d.KitchenConfig.PhaseRewards {
		for _, slug := range phaseData.RewardSlugs {
			rewardEntry := entities.KitchenPhaseCompletionRewards{
				PhaseNumber: phaseData.PhaseNumber,
				RewardSlug:  slug,
			}
			kitchenPhaseRewards = append(kitchenPhaseRewards, rewardEntry)
		}
	}

	gameStageConfig.KitchenPhaseReward = kitchenPhaseRewards

	var kitchenStations []entities.KitchenStation
	for _, kitchenStation := range d.KitchenStations {
		data := entities.KitchenStation{
			FoodItemSlug: kitchenStation.FoodItemSlug,
			AutoUnlock:   kitchenStation.AutoUnlock,
		}

		kitchenStations = append(kitchenStations, data)
	}

	gameStageConfig.KitchenStations = kitchenStations

	return &gameStageConfig, nil
}

func (d *CreateGameStageRequest) ToEntities() (
	*entities.GameStage,
	*entities.GameStageConfig,
	error,
) {
	config, err := d.toEntitiesCommon()
	if err != nil {
		return nil, nil, err
	}

	gameStage := &entities.GameStage{
		Slug:         d.Slug,
		Name:         d.Name,
		Description:  d.Description,
		StartingCoin: d.StartingCoin,
		StagePrize:   d.StagePrize,
		IsActive:     d.IsActive,
		Sequence:     d.Sequence,
	}

	return gameStage, config, nil
}

func (d *UpdateGameStageRequest) ToEntities(id int64) (
	*entities.GameStage,
	*entities.GameStageConfig,
	error,
) {
	config, err := d.toEntitiesCommon()
	if err != nil {
		return nil, nil, err
	}

	gameStage := &entities.GameStage{
		ID:           id,
		Name:         d.Name,
		Description:  d.Description,
		StartingCoin: d.StartingCoin,
		StagePrize:   d.StagePrize,
		IsActive:     d.IsActive,
		Sequence:     d.Sequence,
	}

	return gameStage, config, nil
}
