package usecase

import (
	"context"
	"database/sql"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/database"
	"github.com/winartodev/cat-cafe/pkg/helper"
	"time"
)

const (
	maxDailyRewardDays = 3
)

type DailyRewardUseCase interface {
	CreateDailyReward(ctx context.Context, data entities.DailyReward, rewardTypeSlug string) (res *entities.DailyReward, err error)
	GetDailyRewards(ctx context.Context, limit, offset int) (res []entities.DailyReward, totalRows int64, err error)
	GetDailyRewardID(ctx context.Context, id int64) (res *entities.DailyReward, err error)
	UpdateDailyReward(ctx context.Context, id int64, data entities.DailyReward, rewardTypeSlug string) (res *entities.DailyReward, err error)

	GetDailyRewardStatus(ctx context.Context) (rewards []entities.DailyReward, dailyRewardIdx *int64, isNewDay *bool, err error)
	ClaimDailyReward(ctx context.Context) (reward *entities.DailyReward, newBalance *entities.UserBalance, err error)
}

type dailyRewardUseCase struct {
	userUseCase     UserUseCase
	rewardUseCase   RewardUseCase
	dailyRewardRepo repositories.DailyRewardRepository
	userProgression repositories.UserProgressionRepository
	userRepo        repositories.UserRepository
}

func NewDailyRewardUseCase(
	dailyRewardRepo repositories.DailyRewardRepository,
	userProgression repositories.UserProgressionRepository,
	userRepo repositories.UserRepository,
	userUseCase UserUseCase,
	rewardUseCase RewardUseCase,
) DailyRewardUseCase {
	return &dailyRewardUseCase{
		userUseCase:     userUseCase,
		rewardUseCase:   rewardUseCase,
		dailyRewardRepo: dailyRewardRepo,
		userProgression: userProgression,
		userRepo:        userRepo,
	}
}

func (d *dailyRewardUseCase) CreateDailyReward(ctx context.Context, data entities.DailyReward, rewardSlug string) (res *entities.DailyReward, err error) {
	rewardType, err := d.rewardUseCase.GetRewardBySlug(ctx, rewardSlug)
	if err != nil {
		return nil, err
	}

	if rewardType == nil {
		return nil, apperror.ErrRecordNotFound
	}

	data.Reward = rewardType

	id, err := d.dailyRewardRepo.CreateDailyRewardDB(ctx, data)
	if database.IsDuplicateError(err) {
		return nil, apperror.ErrConflict
	} else if err != nil {
		return nil, err
	}

	data.ID = *id

	//_ = d.dailyRewardRepo.DeleteDailyRewardsRedis(ctx)

	return &data, err
}

func (d *dailyRewardUseCase) GetDailyRewards(ctx context.Context, limit, offset int) (res []entities.DailyReward, totalRows int64, err error) {
	dailyRewards, err := d.dailyRewardRepo.GetDailyRewardsWithPaginationDB(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	totalRows, err = d.dailyRewardRepo.CountDailyRewardsDB(ctx)
	if err != nil {
		return nil, 0, err
	}

	return dailyRewards, totalRows, nil
}

func (d *dailyRewardUseCase) GetDailyRewardID(ctx context.Context, id int64) (res *entities.DailyReward, err error) {
	res, err = d.dailyRewardRepo.GetDailyRewardByIDDB(ctx, id)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, apperror.ErrRecordNotFound
	}

	return res, err
}

func (d *dailyRewardUseCase) UpdateDailyReward(ctx context.Context, id int64, data entities.DailyReward, rewardSlug string) (res *entities.DailyReward, err error) {
	reward, err := d.rewardUseCase.GetRewardBySlug(ctx, rewardSlug)
	if err != nil {
		return nil, err
	}

	if reward == nil {
		return nil, apperror.ErrRecordNotFound
	}

	data.Reward = reward

	err = d.dailyRewardRepo.UpdateDailyRewardDB(ctx, id, data)
	if err != nil {
		return nil, err
	}

	res, err = d.GetDailyRewardID(ctx, id)
	if err != nil {
		return nil, err
	}

	//_ = d.dailyRewardRepo.DeleteDailyRewardsRedis(ctx)

	return res, err
}

func (d *dailyRewardUseCase) getDailyRewardByDay(ctx context.Context, day int64) (res *entities.DailyReward, err error) {
	allRewards, _, err := d.GetDailyRewards(ctx, maxDailyRewardDays, 0)
	if err != nil {
		return nil, err
	}

	var reward *entities.DailyReward
	for _, r := range allRewards {
		if r.DayNumber == day {
			reward = &r
			break
		}
	}

	if reward == nil {
		return nil, apperror.ErrRecordNotFound
	}

	return reward, nil
}

func (d *dailyRewardUseCase) GetDailyRewardStatus(ctx context.Context) ([]entities.DailyReward, *int64, *bool, error) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil || userID <= 0 {
		return nil, nil, nil, apperror.ErrUnauthorized
	}

	progression, err := d.userUseCase.GetUserDailyRewardByID(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	rewards, _, err := d.GetDailyRewards(ctx, maxDailyRewardDays, 0)
	if err != nil {
		return nil, nil, nil, err
	}

	now := helper.NowUTC()
	today := now.Format("2006-01-02")

	// Default to true for new users. If LastClaimDate exists, check if today is a later date.
	isNewDay := true
	var dailyRewardIdx int64

	if progression == nil {
		// initialize an empty progress object.
		progression = &entities.UserDailyReward{LongestStreak: 0, CurrentStreak: 0, LastClaimDate: nil}
	} else if progression.LastClaimDate != nil {
		lastClaimDay := progression.LastClaimDate.Format("2006-01-02")
		isNewDay = today > lastClaimDay
	}

	// Calculate the current day index in the 7-day cycle (0-6)
	// LongestStreak tracks total days claimed (can exceed 7)
	// dailyRewardIdx maps this to a position in the 7-day cycle using modulo
	if isNewDay {
		// New day: user hasn't claimed yet today
		// Point to the next available day to claim
		// Example: LongestStreak=1 (claimed Day 1) → dailyRewardIdx=1 (Day 2 available)
		dailyRewardIdx = progression.LongestStreak % maxDailyRewardDays
	} else {
		// Same day: user already claimed earlier today
		// Point to the day they just claimed
		// Example: LongestStreak=1 (just claimed Day 1) → dailyRewardIdx=0 (Day 1 shown as claimed)
		dailyRewardIdx = (progression.LongestStreak - 1) % maxDailyRewardDays
		if progression.LongestStreak == 0 {
			dailyRewardIdx = 0
		}
	}

	// We do deep copy to prevent corruption while modifying data
	rewardCopy := make([]entities.DailyReward, len(rewards))
	copy(rewardCopy, rewards)

	for i := range rewardCopy {
		var status entities.RewardStatus

		rewardIdx := rewardCopy[i].DayNumber - 1

		if rewardIdx < dailyRewardIdx {
			// Days before current position = already claimed in previous days
			status = entities.StatusClaimed
		} else if rewardIdx == dailyRewardIdx {
			// Current day position = available to claim OR already claimed today
			if isNewDay {
				status = entities.StatusAvailable // Can claim today
			} else {
				status = entities.StatusClaimed // Already claimed today
			}
		} else {
			// Days after current position = locked until user progresses
			status = entities.StatusLocked
		}

		rewardCopy[i].Status = status

		// Calculate cycle number for display (Day 1-7, 8-14, 15-21, etc.)
		// Use LongestStreak-1 if already claimed today to prevent premature cycle increment
		claimedDays := progression.LongestStreak
		if !isNewDay {
			claimedDays = progression.LongestStreak - 1
		}

		// Determine which cycle (0 = Days 1-7, 1 = Days 8-14, etc.)
		cycleNumber := claimedDays / maxDailyRewardDays

		// Add cycle offset to base day number
		rewardCopy[i].DayNumber = rewardCopy[i].DayNumber + (cycleNumber * maxDailyRewardDays)
	}

	return rewardCopy, &dailyRewardIdx, &isNewDay, err
}

func (d *dailyRewardUseCase) ClaimDailyReward(ctx context.Context) (dailyReward *entities.DailyReward, newBalance *entities.UserBalance, err error) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil || userID <= 0 {
		return nil, nil, apperror.ErrUnauthorized
	}

	// Start transaction early to lock the user row
	err = d.dailyRewardRepo.DailyRewardWithTx(ctx, func(tx *sql.Tx) error {
		userRepoTx := d.userRepo.WithTx(tx)
		userProgressionTx := d.userProgression.WithTx(tx)

		// Lock user row to prevent concurrent claims
		_, err = userRepoTx.GetUserByIDForUpdateDB(ctx, userID)
		if err != nil {
			return err
		}

		// Re-fetch progression inside the transaction/lock
		progression, err := userProgressionTx.GetUserDailyRewardByIDDB(ctx, userID)
		if err != nil {
			return err
		}

		now := helper.NowUTC()
		today := now.Truncate(24 * time.Hour)

		// Initialize progression for new users who haven't claimed any rewards yet
		if progression == nil {
			progression = &entities.UserDailyReward{LongestStreak: 0, CurrentStreak: 0, LastClaimDate: nil}
		} else if progression.LastClaimDate != nil {
			// Check if user already claimed today
			lastClaim := progression.LastClaimDate.UTC().Truncate(24 * time.Hour)
			if today.Equal(lastClaim) {
				return apperror.ErrAlreadyClaimed
			}

			diffDay := today.Sub(lastClaim).Hours() / 24
			if diffDay > 1 {
				progression.CurrentStreak = 0
			}
		}

		// Calculate which day to claim in the 7-day cycle
		// LongestStreak tracks total claims, rewardIdx maps to cycle position (0-6)
		rewardIdx := progression.LongestStreak % maxDailyRewardDays
		dayToClaim := rewardIdx + 1

		dailyReward, err = d.getDailyRewardByDay(ctx, dayToClaim)
		if err != nil {
			return err
		}

		rewardType := dailyReward.Reward.RewardType
		rewardTypeEnum, err := entities.ToRewardType(rewardType.Slug)
		if err != nil {
			return err
		}

		// Increment the streak counter after successful claim
		newLongestStreak := progression.LongestStreak + 1
		newCurrentStreak := progression.CurrentStreak + 1

		// Update longest streak if current streak exceeds it
		if newCurrentStreak > newLongestStreak {
			newLongestStreak = newCurrentStreak
		}

		// Update user's streak and last claim date
		err = userProgressionTx.UpsertDailyRewardProgressionDB(ctx, userID, newLongestStreak, newCurrentStreak, now)
		if err != nil {
			return err
		}

		// Handle dailyReward based on type
		if rewardTypeEnum.RequiresBalanceUpdate() {
			// For COIN and GEM, update user balance in database
			balanceType := rewardTypeEnum.ToUserBalance()
			err = userRepoTx.UpdateUserBalanceWithTx(ctx, userID, balanceType, dailyReward.Reward.Amount)
			if err != nil {
				return err
			}
		} else if rewardTypeEnum.IsSentExternally() {
			// TODO: Call External API to give GoPay Coin to player
		} else {
			return apperror.ErrUnknownRewardType
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	// Clear the cache so that the next request retrieves the latest progress data from the database
	//_ = d.userDailyRewardRepo.DeleteUserDailyRewardRedis(ctx, userID)

	if dailyReward != nil {
		rewardTypeEnum, _ := entities.ToRewardType(dailyReward.Reward.RewardType.Slug)
		if rewardTypeEnum.RequiresBalanceUpdate() {
			newBalance, err = d.userUseCase.GetUserBalance(ctx, userID)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	return dailyReward, newBalance, err
}
