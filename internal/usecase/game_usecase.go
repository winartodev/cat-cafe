package usecase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/winartodev/cat-cafe/internal/dto"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

type GameUseCase interface {
	UpdateUserBalance(ctx context.Context, coinEarned int64, userID int64) (res *dto.UserBalanceResponse, err error)
	GetUserGameData(ctx context.Context, userID int64) (res *entities.Game, err error)
	GetGameStages(ctx context.Context, userID int64) (res []entities.UserGameStage, nextStage *entities.UserNextGameStageInfo, err error)
	StartGameStage(ctx context.Context, userID int64, slug string) (stage *entities.GameStage, config *entities.GameStageConfig, nextStage *entities.UserNextGameStageInfo, err error)
	CompleteGameStage(ctx context.Context, userID int64, slug string) (err error)
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

func (g *gameUseCase) GetGameStages(ctx context.Context, userID int64) (res []entities.UserGameStage, nextStage *entities.UserNextGameStageInfo, err error) {
	gameStages, err := g.gameStageRepo.GetActiveGameStagesDB(ctx)
	if err != nil || len(gameStages) == 0 {
		return nil, nil, err
	}

	latestProgress, err := g.userProgressionRepo.GetLatestGameStageProgressionDB(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	// Create new progression if player doesn't have any
	if latestProgress == nil {
		firstStage := gameStages[0]
		_, err = g.userProgressionRepo.CreateGameStageProgressionDB(ctx, userID, firstStage.ID)
		if err != nil {
			return nil, nil, err
		}

		// recursive
		return g.GetGameStages(ctx, userID)
	}

	stages, nextStage := g.mapToUserGameStage(gameStages, latestProgress)

	return stages, nextStage, nil
}

func (g *gameUseCase) StartGameStage(ctx context.Context, userID int64, slug string) (stage *entities.GameStage, config *entities.GameStageConfig, nextStage *entities.UserNextGameStageInfo, err error) {
	stage, err = g.gameStageRepo.GetGameStageBySlugDB(ctx, slug)
	if err != nil {
		return nil, nil, nil, err
	}

	// Get user's available stages
	userStages, nextStage, err := g.GetGameStages(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Check if user can access this stage
	var canAccess bool
	for _, userStage := range userStages {
		if userStage.Slug == slug {
			// Only allow if status is Current or Complete (if you allow replay)
			if userStage.Status == entities.GSStatusCurrent {
				canAccess = true
			} else if userStage.Status == entities.GSStatusComplete {
				// Decide if you want to allow replaying completed stages
				return nil, nil, nil, apperror.ErrStageAlreadyCompleted
			}
			break
		}
	}

	if !canAccess {
		return nil, nil, nil, apperror.ErrStageNotUnlocked
	}

	config, err = g.gameStageRepo.GetGameConfigByIDDB(ctx, stage.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	return stage, config, nextStage, nil
}

func (g *gameUseCase) mapToUserGameStage(stages []entities.GameStage, lastProgress *entities.UserGameStageProgression) ([]entities.UserGameStage, *entities.UserNextGameStageInfo) {
	isFoundCurrent := false

	var res = make([]entities.UserGameStage, len(stages))
	var nextStage *entities.UserNextGameStageInfo

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

				nextStage = g.setNextStageIfExists(stages, i+1, entities.GSStatusLocked)
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

					nextStage = g.setNextStageIfExists(stages, i+1, entities.GSStatusCurrent)
				} else {
					currentStage.Status = entities.GSStatusCurrent

					nextStage = g.setNextStageIfExists(stages, i+1, entities.GSStatusLocked)
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

					nextStage = g.setNextStageIfExists(stages, i+1, entities.GSStatusLocked)
				} else {
					// Stage is still locked (either already found current or last progress incomplete)
					currentStage.Status = entities.GSStatusLocked
				}
			}
		}

		res[i] = currentStage
	}

	return res, nextStage
}

func (g *gameUseCase) CompleteGameStage(ctx context.Context, userID int64, slug string) (err error) {
	stage, err := g.gameStageRepo.GetGameStageBySlugDB(ctx, slug)
	if err != nil {
		return err
	}

	// TODO: FIX THIS EITHER USE LATEST GAME STAGE PROGRESSION OR BY STAGE GAME
	userStageProgress, err := g.userProgressionRepo.GetGameStageProgressionDB(ctx, userID, stage.ID)
	if err != nil {
		return err
	}

	if err = g.validateLastProgression(userStageProgress, stage); err != nil {
		return
	}

	err = g.userProgressionRepo.MarkStageAsCompleteDB(ctx, userID, stage.ID)
	if err != nil {
		return err
	}

	gameStages, err := g.gameStageRepo.GetActiveGameStagesDB(ctx)
	if err != nil {
		return err
	}

	for i, s := range gameStages {
		if s.ID == stage.ID && i+1 < len(gameStages) {
			nextStageID := gameStages[i+1].ID
			_, err = g.userProgressionRepo.CreateGameStageProgressionDB(ctx, userID, nextStageID)
			if err != nil && !errors.Is(err, apperror.ErrConflict) {
				return err
			}
			break
		}
	}

	return nil
}

func (g *gameUseCase) validateLastProgression(lastProgression *entities.UserGameStageProgression, stage *entities.GameStage) error {
	if lastProgression == nil || stage == nil {
		return apperror.ErrStageNotUnlocked
	}

	if lastProgression.StageID != stage.ID {
		return apperror.ErrStageNotUnlocked
	}

	if lastProgression.IsComplete {
		return apperror.ErrStageAlreadyCompleted
	}

	return nil
}

func (g *gameUseCase) getSequenceByID(stages []entities.GameStage, id int64) int64 {
	for _, s := range stages {
		if s.ID == id {
			return s.Sequence
		}
	}
	return 0
}

func (g *gameUseCase) setNextStageIfExists(stages []entities.GameStage, index int, status entities.GameStageStatus) *entities.UserNextGameStageInfo {
	if index >= len(stages) {
		return nil
	}

	return &entities.UserNextGameStageInfo{
		Slug:   stages[index].Slug,
		Name:   stages[index].Name,
		Status: status,
	}
}
