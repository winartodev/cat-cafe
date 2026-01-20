package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/winartodev/cat-cafe/internal/dto"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

type GameUseCase interface {
	UpdateUserBalance(ctx context.Context, coinEarned int64, userID int64) (res *dto.UserBalanceResponse, err error)
	GetUserGameData(ctx context.Context, userID int64) (res *entities.Game, err error)
	GetGameStages(ctx context.Context, userID int64) (res []entities.UserGameStage, err error)
}

type gameUseCase struct {
	userUseCase         UserUseCase
	userRepo            repositories.UserRepository
	userProgressionRepo repositories.UserProgressionRepository
	gameStageRepo       repositories.GameStageRepository
}

func NewGameUseCase(
	userUc UserUseCase,
	userRepo repositories.UserRepository,
	userProgressionRepo repositories.UserProgressionRepository,
	gameStageRepo repositories.GameStageRepository,
) GameUseCase {
	return &gameUseCase{
		userUseCase:         userUc,
		userRepo:            userRepo,
		userProgressionRepo: userProgressionRepo,
		gameStageRepo:       gameStageRepo,
	}
}

func (g *gameUseCase) UpdateUserBalance(ctx context.Context, coinEarned int64, userID int64) (res *dto.UserBalanceResponse, err error) {
	// TODO: should we validate earning rate ???

	err = g.userRepo.BalanceWithTx(ctx, func(tx *sql.Tx) error {
		txRepo := g.userRepo.WithTx(tx)

		if err := txRepo.UpdateUserBalanceWithTx(ctx, userID, entities.BalanceTypeCoin, coinEarned); err != nil {
			return err
		}

		if err := txRepo.UpdateLastSyncBalanceWithTx(ctx, userID, helper.NowUTC()); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	_ = g.userRepo.DeleteUserRedis(ctx, userID)

	user, err := g.userUseCase.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserBalanceResponse{
		Coin: user.UserBalance.Coin,
		Gem:  user.UserBalance.Gem,
	}, nil
}

func (g *gameUseCase) GetUserGameData(ctx context.Context, userID int64) (res *entities.Game, err error) {
	isDailyRewardAvailable, err := g.userUseCase.IsDailyRewardAvailable(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &entities.Game{
		DailyRewardAvailable: isDailyRewardAvailable,
	}, nil
}

func (g *gameUseCase) GetGameStages(ctx context.Context, userID int64) (res []entities.UserGameStage, err error) {
	gameStages, err := g.gameStageRepo.GetActiveGameStagesDB(ctx)
	if err != nil || len(gameStages) == 0 {
		return nil, err
	}

	lastProgress, err := g.userProgressionRepo.GetGameStageProgressionDB(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create new one if player haven't  last progress
	if lastProgress == nil {
		firstStage := gameStages[0]
		helper.PrettyPrint(firstStage)
		fmt.Println(firstStage.ID)
		_, err = g.userProgressionRepo.CreateGameStageProgressionDB(ctx, userID, firstStage.ID)
		if err != nil {
			return nil, err
		}

		// recursive
		return g.GetGameStages(ctx, userID)
	}

	if lastProgress.IsComplete {
		var nextStageID int64
		for i, stage := range gameStages {
			if stage.ID == lastProgress.StageID {
				if i+1 < len(gameStages) {
					nextStageID = gameStages[i+1].ID
				}
				break
			}
		}

		if nextStageID != 0 {
			_, err = g.userProgressionRepo.CreateGameStageProgressionDB(ctx, userID, nextStageID)
			if err != nil {
				return nil, err
			}

			return g.GetGameStages(ctx, userID)
		}
	}

	return g.mapToUserGameStage(gameStages, lastProgress), nil
}

func (g *gameUseCase) mapToUserGameStage(stages []entities.GameStage, lastProgress *entities.UserGameStageProgression) []entities.UserGameStage {
	isFoundCurrent := false
	var res = make([]entities.UserGameStage, len(stages))
	for i, stage := range stages {
		currentStage := entities.UserGameStage{
			Slug:     stage.Slug,
			Name:     stage.Name,
			Sequence: stage.Sequence,
		}

		// User hasn't started any stage yet
		if lastProgress == nil {
			// First stage should be current/available
			if i == 0 {
				currentStage.Status = entities.GSStatusCurrent
				isFoundCurrent = true
			} else {
				// All other stages are locked
				currentStage.Status = entities.GSStatusLocked
			}
		} else {
			// User has some progression

			// This is the stage user is currently on or has completed
			if lastProgress.StageID == stage.ID {
				if lastProgress.IsComplete {
					currentStage.Status = entities.GSStatusComplete
				} else {
					currentStage.Status = entities.GSStatusCurrent
				}
			} else if currentStage.Sequence < g.getSequenceByID(stages, lastProgress.StageID) {
				// This stage comes before the user's last progress, so it's completed
				currentStage.Status = entities.GSStatusComplete
			} else {
				// This stage comes after the user's last progress

				// If last progress is complete, and we haven't set a current stage yet,
				// this is the next available stage
				if lastProgress.IsComplete && !isFoundCurrent {
					currentStage.Status = entities.GSStatusCurrent
					isFoundCurrent = true
				} else {
					// Stage is still locked (either already found current or last progress incomplete)
					currentStage.Status = entities.GSStatusLocked

				}
			}
		}

		res[i] = currentStage
	}

	return res
}

func (g *gameUseCase) getSequenceByID(stages []entities.GameStage, id int64) int64 {
	for _, s := range stages {
		if s.ID == id {
			return s.Sequence
		}
	}
	return 0
}
