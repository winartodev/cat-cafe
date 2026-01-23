package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
)

type GameStageUseCase interface {
	CreateGameStage(ctx context.Context, data *entities.GameStage, config *entities.GameStageConfig) (*entities.GameStage, error)
	UpdateGameStage(ctx context.Context, data *entities.GameStage, config *entities.GameStageConfig) (*entities.GameStage, error)
	GetGameStages(ctx context.Context, limit, offset int) ([]entities.GameStage, int64, error)
	GetGameStageByID(ctx context.Context, id int64) (*entities.GameStage, *entities.GameStageConfig, error)
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
	}
}

func (u *gameStageUseCase) CreateGameStage(ctx context.Context, data *entities.GameStage, config *entities.GameStageConfig) (*entities.GameStage, error) {
	gameStage, err := u.gameStageRepo.GetGameStageBySlugDB(ctx, data.Slug)
	if err != nil {
		return nil, err
	}

	if gameStage != nil {
		return nil, fmt.Errorf("game stage %s already exists", data.Slug)
	}

	err = u.gameStageRepo.GameStageWithTx(ctx, func(tx *sql.Tx) error {
		stageRepo := u.gameStageRepo.WithTx(tx)
		customerRepo := u.customerConfigRepo.WithTx(tx)
		staffRepo := u.staffConfigRepo.WithTx(tx)
		kitchenConfigRepo := u.kitchenConfigRepo.WithTx(tx)
		cameraConfigRepo := u.cameraConfigRepo.WithTx(tx)
		kitchenStationRepo := u.kitchenStationRepo.WithTX(tx)
		foodItemRepo := u.foodItemRepo.WithTx(tx)

		stageID, err := stageRepo.CreateGameStageWithTxDB(ctx, data)
		if err != nil {
			return err
		}

		data.ID = *stageID

		_, err = customerRepo.CreateCustomerConfigWithTxDB(ctx, data.ID, config.CustomerConfig)
		if err != nil {
			return err
		}

		_, err = staffRepo.CreateStaffConfigWithTxDB(ctx, data.ID, config.StaffConfig)
		if err != nil {
			return err
		}

		var slugs []string
		for _, ks := range config.KitchenStations {
			slugs = append(slugs, ks.FoodItemSlug)
		}

		foodItemMap, err := foodItemRepo.GetFoodItemIDsBySlugsDB(ctx, slugs)
		if err != nil {
			return err
		}

		for i := range config.KitchenStations {
			id, ok := foodItemMap[config.KitchenStations[i].FoodItemSlug]
			if !ok {
				return fmt.Errorf("food item slug %s not found", config.KitchenStations[i].FoodItemSlug)
			}

			config.KitchenStations[i].FoodItemID = id
			config.KitchenStations[i].StageID = data.ID
		}

		_, err = kitchenStationRepo.CreateKitchenStationsWithTxDB(ctx, data.ID, config.KitchenStations)
		if err != nil {
			return err
		}

		kitchenConfigID, err := kitchenConfigRepo.CreateKitchenConfigWithTxDB(ctx, data.ID, config.KitchenConfig)
		if err != nil {
			return err
		}

		err = u.createKitchenCompleteReward(ctx, tx, *kitchenConfigID, config.KitchenPhaseReward)
		if err != nil {
			return err
		}

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

func (u *gameStageUseCase) UpdateGameStage(ctx context.Context, data *entities.GameStage, config *entities.GameStageConfig) (*entities.GameStage, error) {
	gameStage, err := u.gameStageRepo.GetGameStageByIDDB(ctx, data.ID)
	if err != nil {
		return nil, err
	}

	if gameStage == nil {
		return nil, fmt.Errorf("game stage %d not found", data.ID)
	}

	err = u.gameStageRepo.GameStageWithTx(ctx, func(tx *sql.Tx) error {
		stageRepo := u.gameStageRepo.WithTx(tx)
		customerRepo := u.customerConfigRepo.WithTx(tx)
		staffRepo := u.staffConfigRepo.WithTx(tx)
		kitchenConfigRepo := u.kitchenConfigRepo.WithTx(tx)
		cameraConfigRepo := u.cameraConfigRepo.WithTx(tx)
		kitchenStationRepo := u.kitchenStationRepo.WithTX(tx)
		foodItemRepo := u.foodItemRepo.WithTx(tx)

		err := stageRepo.UpdateGameStageWithTxDB(ctx, data)
		if err != nil {
			return err
		}

		err = customerRepo.UpdateCustomerConfigWithTxDB(ctx, data.ID, config.CustomerConfig)
		if err != nil {
			return err
		}

		err = staffRepo.UpdateStaffConfigWithTxDB(ctx, data.ID, config.StaffConfig)
		if err != nil {
			return err
		}

		var slugs []string
		for _, ks := range config.KitchenStations {
			slugs = append(slugs, ks.FoodItemSlug)
		}

		foodItemMap, err := foodItemRepo.GetFoodItemIDsBySlugsDB(ctx, slugs)
		if err != nil {
			return err
		}

		for i := range config.KitchenStations {
			id, ok := foodItemMap[config.KitchenStations[i].FoodItemSlug]
			if !ok {
				return fmt.Errorf("food item slug %s not found", config.KitchenStations[i].FoodItemSlug)
			}
			config.KitchenStations[i].FoodItemID = id
			config.KitchenStations[i].StageID = data.ID
		}

		err = kitchenStationRepo.DeleteKitchenStationDB(ctx, data.ID)
		if err != nil {
			return err
		}

		_, err = kitchenStationRepo.CreateKitchenStationsWithTxDB(ctx, data.ID, config.KitchenStations)
		if err != nil {
			return err
		}

		kitchenConfigID, err := kitchenConfigRepo.UpdateKitchenConfigWithTxDB(ctx, data.ID, config.KitchenConfig)
		if err != nil {
			return err
		}

		err = kitchenConfigRepo.DeleteKitchenCompletionReward(ctx, *kitchenConfigID)
		if err != nil {
			return err
		}

		err = u.createKitchenCompleteReward(ctx, tx, *kitchenConfigID, config.KitchenPhaseReward)
		if err != nil {
			return err
		}

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

func (u *gameStageUseCase) GetGameStages(ctx context.Context, limit, offset int) ([]entities.GameStage, int64, error) {
	return u.gameStageRepo.GetGameStagesDB(ctx, limit, offset)
}

func (u *gameStageUseCase) GetGameStageByID(ctx context.Context, id int64) (*entities.GameStage, *entities.GameStageConfig, error) {
	gameStage, err := u.gameStageRepo.GetGameStageByIDDB(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	gameConfig, err := u.getGameConfig(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return gameStage, gameConfig, nil
}

func (u *gameStageUseCase) getGameConfig(ctx context.Context, stageID int64) (*entities.GameStageConfig, error) {
	gameConfig, err := u.gameStageRepo.GetGameConfigByIDDB(ctx, stageID)
	if err != nil {
		return nil, err
	}

	if gameConfig == nil {
		return nil, nil
	}

	return gameConfig, nil
}

func (u *gameStageUseCase) createKitchenCompleteReward(ctx context.Context, tx *sql.Tx, kitchenConfigID int64, phaseRewards []entities.KitchenPhaseCompletionRewards) error {
	kitchenConfigRepo := u.kitchenConfigRepo.WithTx(tx)
	rewardRepo := u.rewardRepo.WithTx(tx)
	for _, phaseReward := range phaseRewards {
		reward, err := rewardRepo.GetRewardBySlugDB(ctx, phaseReward.RewardSlug)
		if err != nil {
			return err
		}

		if reward == nil {
			return fmt.Errorf("reward %s not exist", phaseReward.RewardSlug)
		}

		phaseReward.RewardID = reward.ID
		_, err = kitchenConfigRepo.CreateKitchenCompletionReward(ctx, kitchenConfigID, &phaseReward)
		if err != nil {
			return err
		}
	}

	return nil
}
