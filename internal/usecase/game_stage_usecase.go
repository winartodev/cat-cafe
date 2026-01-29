package usecase

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
)

type GameStageUseCase interface {
	CreateGameStage(ctx context.Context, data *entities.GameStage, config *entities.GameStageConfig) (*entities.GameStage, error)
	UpdateGameStage(ctx context.Context, data *entities.GameStage, config *entities.GameStageConfig) (*entities.GameStage, error)
	GetGameStages(ctx context.Context, limit, offset int) ([]entities.GameStage, int64, error)
	GetGameStageByID(ctx context.Context, id int64) (*entities.GameStage, *entities.GameStageConfig, error)

	CreateStageUpgrade(ctx context.Context, stageSlug string, upgradeTypes []string) error
	GetStageUpgrades(ctx context.Context, stageSlug string, limit, offset int) ([]entities.StageUpgrade, int64, error)
	UpdateStageUpgrades(ctx context.Context, stageSlug string, upgradeTypes []string) error
}

type gameStageUseCase struct {
	gameStageRepo      repositories.GameStageRepository
	customerConfigRepo repositories.StageCustomerConfigRepository
	staffConfigRepo    repositories.StageStaffConfigRepository
	kitchenConfigRepo  repositories.StageKitchenConfigRepository
	cameraConfigRepo   repositories.StageCameraConfigRepository
	rewardRepo         repositories.RewardRepository
	kitchenStationRepo repositories.KitchenStationRepository
	foodItemRepo       repositories.FoodItemRepository
	upgradeRepo        repositories.UpgradeRepository
	stageUpgradeRepo   repositories.StageUpgradeRepository
}

func NewGameStageUseCase(
	gameStageRepo repositories.GameStageRepository,
	customerConfigRepo repositories.StageCustomerConfigRepository,
	staffConfigRepo repositories.StageStaffConfigRepository,
	kitchenConfigRepo repositories.StageKitchenConfigRepository,
	cameraConfigRepo repositories.StageCameraConfigRepository,
	rewardRepo repositories.RewardRepository,
	kitchenStationRepo repositories.KitchenStationRepository,
	foodItemRepo repositories.FoodItemRepository,
	upgradeRepo repositories.UpgradeRepository,
	stageUpgradeRepo repositories.StageUpgradeRepository,
) GameStageUseCase {
	return &gameStageUseCase{
		gameStageRepo:      gameStageRepo,
		customerConfigRepo: customerConfigRepo,
		staffConfigRepo:    staffConfigRepo,
		kitchenConfigRepo:  kitchenConfigRepo,
		cameraConfigRepo:   cameraConfigRepo,
		rewardRepo:         rewardRepo,
		kitchenStationRepo: kitchenStationRepo,
		foodItemRepo:       foodItemRepo,
		upgradeRepo:        upgradeRepo,
		stageUpgradeRepo:   stageUpgradeRepo,
	}
}

// CreateGameStage creates a new game stage with transaction
func (u *gameStageUseCase) CreateGameStage(ctx context.Context, data *entities.GameStage, config *entities.GameStageConfig) (*entities.GameStage, error) {
	gameStage, err := u.gameStageRepo.GetGameStageBySlugDB(ctx, data.Slug)
	if err != nil {
		return nil, err
	}

	if gameStage != nil {
		return nil, apperror.ErrorAlreadyExists("game stage", "slug", data.Slug)
	}

	// Create game stage with transaction
	err = u.gameStageRepo.GameStageWithTx(ctx, func(tx *sql.Tx) error {
		stageRepo := u.gameStageRepo.WithTx(tx)
		customerRepo := u.customerConfigRepo.WithTx(tx)
		staffRepo := u.staffConfigRepo.WithTx(tx)
		kitchenConfigRepo := u.kitchenConfigRepo.WithTx(tx)
		cameraConfigRepo := u.cameraConfigRepo.WithTx(tx)
		kitchenStationRepo := u.kitchenStationRepo.WithTX(tx)
		foodItemRepo := u.foodItemRepo.WithTx(tx)

		// Create game stage
		stageID, err := stageRepo.CreateGameStageWithTxDB(ctx, data)
		if err != nil {
			return err
		}

		data.ID = *stageID

		// Create customer config
		_, err = customerRepo.CreateCustomerConfigWithTxDB(ctx, data.ID, config.CustomerConfig)
		if err != nil {
			return err
		}

		// Create staff config
		_, err = staffRepo.CreateStaffConfigWithTxDB(ctx, data.ID, config.StaffConfig)
		if err != nil {
			return err
		}

		var slugs []string
		for _, ks := range config.KitchenStations {
			slugs = append(slugs, ks.FoodItemSlug)
		}

		// Get food item IDs by slugs
		foodItemMap, err := foodItemRepo.GetFoodItemIDsBySlugsDB(ctx, slugs)
		if err != nil {
			return err
		}

		// Set food item IDs and stage IDs
		for i := range config.KitchenStations {
			// Check if food item slug exists
			id, ok := foodItemMap[config.KitchenStations[i].FoodItemSlug]
			if !ok {
				return apperror.ErrorNotFound("food item", config.KitchenStations[i].FoodItemSlug)
			}

			config.KitchenStations[i].FoodItemID = id
			config.KitchenStations[i].StageID = data.ID
		}

		// Create kitchen stations
		_, err = kitchenStationRepo.CreateKitchenStationsWithTxDB(ctx, data.ID, config.KitchenStations)
		if err != nil {
			return err
		}

		// Create kitchen config
		kitchenConfigID, err := kitchenConfigRepo.CreateKitchenConfigWithTxDB(ctx, data.ID, config.KitchenConfig)
		if err != nil {
			return err
		}

		// Create kitchen complete reward
		err = u.createKitchenCompleteReward(ctx, tx, *kitchenConfigID, config.KitchenPhaseReward)
		if err != nil {
			return err
		}

		// Create camera config
		_, err = cameraConfigRepo.CreateStageCameraDB(ctx, data.ID, config.CameraConfig)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UpdateGameStage updates an existing game stage with transaction
func (u *gameStageUseCase) UpdateGameStage(ctx context.Context, data *entities.GameStage, config *entities.GameStageConfig) (*entities.GameStage, error) {
	gameStage, err := u.gameStageRepo.GetGameStageByIDDB(ctx, data.ID)
	if err != nil {
		return nil, err
	}

	if gameStage == nil {
		return nil, apperror.ErrorNotFound(fmt.Sprintf("game stage id %d", data.ID))
	}

	// Update game stage with transaction
	err = u.gameStageRepo.GameStageWithTx(ctx, func(tx *sql.Tx) error {
		stageRepo := u.gameStageRepo.WithTx(tx)
		customerRepo := u.customerConfigRepo.WithTx(tx)
		staffRepo := u.staffConfigRepo.WithTx(tx)
		kitchenConfigRepo := u.kitchenConfigRepo.WithTx(tx)
		cameraConfigRepo := u.cameraConfigRepo.WithTx(tx)
		kitchenStationRepo := u.kitchenStationRepo.WithTX(tx)
		foodItemRepo := u.foodItemRepo.WithTx(tx)

		// Update game stage
		err := stageRepo.UpdateGameStageWithTxDB(ctx, data)
		if err != nil {
			return err
		}

		// Update customer config
		err = customerRepo.UpdateCustomerConfigWithTxDB(ctx, data.ID, config.CustomerConfig)
		if err != nil {
			return err
		}

		// Update staff config
		err = staffRepo.UpdateStaffConfigWithTxDB(ctx, data.ID, config.StaffConfig)
		if err != nil {
			return err
		}

		// Get food item IDs by slugs
		var slugs []string
		for _, ks := range config.KitchenStations {
			slugs = append(slugs, ks.FoodItemSlug)
		}

		// Get food item IDs by slugs
		foodItemMap, err := foodItemRepo.GetFoodItemIDsBySlugsDB(ctx, slugs)
		if err != nil {
			return err
		}

		// Set food item IDs and stage IDs
		for i := range config.KitchenStations {
			// Check if food item slug exists
			id, ok := foodItemMap[config.KitchenStations[i].FoodItemSlug]
			if !ok {
				return apperror.ErrorNotFound("food item", config.KitchenStations[i].FoodItemSlug)
			}
			config.KitchenStations[i].FoodItemID = id
			config.KitchenStations[i].StageID = data.ID
		}

		// Delete kitchen stations
		err = kitchenStationRepo.DeleteKitchenStationDB(ctx, data.ID)
		if err != nil {
			return err
		}

		// Create kitchen stations
		_, err = kitchenStationRepo.CreateKitchenStationsWithTxDB(ctx, data.ID, config.KitchenStations)
		if err != nil {
			return err
		}

		// Update kitchen config
		kitchenConfigID, err := kitchenConfigRepo.UpdateKitchenConfigWithTxDB(ctx, data.ID, config.KitchenConfig)
		if err != nil {
			return err
		}

		// Delete kitchen complete reward
		err = kitchenConfigRepo.DeleteKitchenCompletionRewardDB(ctx, *kitchenConfigID)
		if err != nil {
			return err
		}

		// Create kitchen complete reward
		err = u.createKitchenCompleteReward(ctx, tx, *kitchenConfigID, config.KitchenPhaseReward)
		if err != nil {
			return err
		}

		// Update camera config
		err = cameraConfigRepo.UpdateStageCameraDB(ctx, data.ID, config.CameraConfig)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetGameStages gets all game stages with limit and offset
func (u *gameStageUseCase) GetGameStages(ctx context.Context, limit, offset int) ([]entities.GameStage, int64, error) {
	return u.gameStageRepo.GetGameStagesDB(ctx, limit, offset)
}

// GetGameStageByID gets a game stage by ID
func (u *gameStageUseCase) GetGameStageByID(ctx context.Context, id int64) (*entities.GameStage, *entities.GameStageConfig, error) {
	gameStage, err := u.gameStageRepo.GetGameStageByIDDB(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	if gameStage == nil {
		return nil, nil, apperror.ErrorNotFound(fmt.Sprintf("game stage id %d", id))
	}

	gameConfig, err := u.getGameConfig(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return gameStage, gameConfig, nil
}

// getGameConfig gets the game config for a specific stage
func (u *gameStageUseCase) getGameConfig(ctx context.Context, stageID int64) (*entities.GameStageConfig, error) {
	gameConfig, err := u.gameStageRepo.GetGameConfigByIDDB(ctx, stageID)
	if err != nil {
		return nil, err
	}

	if gameConfig == nil {
		return nil, apperror.ErrorNotFound(fmt.Sprintf("game config id %d", stageID))
	}

	return gameConfig, nil
}

// createKitchenCompleteReward creates kitchen complete reward with transaction
func (u *gameStageUseCase) createKitchenCompleteReward(ctx context.Context, tx *sql.Tx, kitchenConfigID int64, phaseRewards []entities.KitchenPhaseCompletionRewards) error {
	kitchenConfigRepo := u.kitchenConfigRepo.WithTx(tx)
	rewardRepo := u.rewardRepo.WithTx(tx)
	for _, phaseReward := range phaseRewards {
		reward, err := rewardRepo.GetRewardBySlugDB(ctx, phaseReward.Reward.Slug)
		if err != nil {
			return err
		}

		if reward == nil {
			return apperror.ErrorNotFound("reward", phaseReward.Reward.Slug)
		}

		phaseReward.RewardID = reward.ID
		_, err = kitchenConfigRepo.CreateKitchenCompletionRewardDB(ctx, kitchenConfigID, &phaseReward)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateStageUpgrade creates a new stage upgrade with transaction
func (u *gameStageUseCase) CreateStageUpgrade(ctx context.Context, stageSlug string, upgradeTypes []string) error {
	stage, err := u.getStageUpgradeBySlug(ctx, stageSlug)
	if err != nil {
		return err
	}

	stageUpgrades, err := u.buildStageUpgrades(ctx, stage.ID, upgradeTypes)
	if err != nil {
		return err
	}

	err = u.stageUpgradeRepo.StageUpgradeWithTx(ctx, func(tx *sql.Tx) error {
		stageUpgradeTx := u.stageUpgradeRepo.WithTx(tx)
		return stageUpgradeTx.BulkCreateStageUpgradesDB(ctx, stageUpgrades)
	})

	return err
}

// GetStageUpgrades gets all stage upgrades with limit and offset
func (u *gameStageUseCase) GetStageUpgrades(ctx context.Context, stageSlug string, limit, offset int) ([]entities.StageUpgrade, int64, error) {
	stage, err := u.getStageUpgradeBySlug(ctx, stageSlug)
	if err != nil {
		return nil, 0, err
	}

	upgrades, err := u.stageUpgradeRepo.GetStageUpgradesDB(ctx, stage.ID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := u.stageUpgradeRepo.CountStageUpgradesDB(ctx, stage.ID)
	if err != nil {
		return nil, 0, err
	}

	return upgrades, total, nil
}

// UpdateStageUpgrades updates stage upgrades with transaction
func (u *gameStageUseCase) UpdateStageUpgrades(ctx context.Context, stageSlug string, upgradeTypes []string) error {
	stage, err := u.gameStageRepo.GetGameStageBySlugDB(ctx, stageSlug)
	if err != nil {
		return err
	}

	stageUpgrades, err := u.buildStageUpgrades(ctx, stage.ID, upgradeTypes)
	if err != nil {
		return err
	}

	err = u.stageUpgradeRepo.StageUpgradeWithTx(ctx, func(tx *sql.Tx) error {
		stageUpgradeTx := u.stageUpgradeRepo.WithTx(tx)

		err = stageUpgradeTx.DeleteStageUpgradeDB(ctx, stage.ID)
		if err != nil {
			return err
		}

		err = stageUpgradeTx.BulkCreateStageUpgradesDB(ctx, stageUpgrades)
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}

// buildStageUpgrades builds stage upgrades with transaction
func (u *gameStageUseCase) buildStageUpgrades(ctx context.Context, stageID int64, upgradeTypes []string) ([]entities.StageUpgrade, error) {
	if len(upgradeTypes) == 0 {
		return []entities.StageUpgrade{}, nil
	}

	// Get upgrades by slugs
	upgrades, err := u.upgradeRepo.GetUpgradesBySlugsDB(ctx, upgradeTypes)
	if err != nil {
		return nil, err
	}

	// Validate all slugs exist
	if len(upgrades) != len(upgradeTypes) {
		return nil, apperror.ErrorInvalidRequest("some upgrade slugs are invalid")
	}

	// Build StageUpgrade entities
	stageUpgrades := make([]entities.StageUpgrade, len(upgrades))
	for i, upgrade := range upgrades {
		stageUpgrades[i] = entities.StageUpgrade{
			StageID:   stageID,
			UpgradeID: upgrade.ID,
		}
	}

	return stageUpgrades, nil
}

// getStageUpgradeBySlug gets stage upgrade by slug
func (u *gameStageUseCase) getStageUpgradeBySlug(ctx context.Context, stageSlug string) (*entities.GameStage, error) {
	stage, err := u.gameStageRepo.GetGameStageBySlugDB(ctx, stageSlug)
	if err != nil {
		return nil, err
	}

	if stage == nil {
		return nil, apperror.ErrorNotFound("game stage", stageSlug)
	}

	return stage, nil
}
