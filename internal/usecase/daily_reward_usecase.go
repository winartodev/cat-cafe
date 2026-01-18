package usecase

import (
	"context"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

const (
	maxDailyRewardDays = 7
)

type DailyRewardUseCase interface {
	CreateRewardType(ctx context.Context, data entities.RewardType) (res *entities.RewardType, err error)
	UpdateRewardTypes(ctx context.Context, id int64, data entities.RewardType) (res *entities.RewardType, err error)
	GetRewardTypes(ctx context.Context) (res []entities.RewardType, err error)
	GetRewardTypeByID(ctx context.Context, id int64) (res *entities.RewardType, err error)
	GetRewardTypeBySlug(ctx context.Context, slug string) (res *entities.RewardType, err error)

	CreateDailyReward(ctx context.Context, data entities.DailyReward, rewardTypeSlug string) (res *entities.DailyReward, err error)
	GetDailyRewards(ctx context.Context) (res []entities.DailyReward, err error)
	GetDailyRewardID(ctx context.Context, id int64) (res *entities.DailyReward, err error)
	UpdateDailyReward(ctx context.Context, id int64, data entities.DailyReward, rewardTypeSlug string) (res *entities.DailyReward, err error)

	GetRewardStatus(ctx context.Context) (rewards []entities.DailyReward, dailyRewardIdx *int64, isNewDay *bool, err error)
	ClaimReward(ctx context.Context) (reward *entities.DailyReward, newBalance *entities.UserBalance, err error)
}

type dailyRewardUseCase struct {
	userUseCase         UserUseCase
	dailyRewardRepo     repositories.DailyRewardRepository
	userDailyRewardRepo repositories.UserDailyRewardRepository
	userRepo            repositories.UserRepository
}

func NewDailyRewardUseCase(dailyRewardRepo repositories.DailyRewardRepository, userDailyRewardRepo repositories.UserDailyRewardRepository, userRepo repositories.UserRepository, userUseCase UserUseCase) DailyRewardUseCase {
	return &dailyRewardUseCase{
		userUseCase:         userUseCase,
		dailyRewardRepo:     dailyRewardRepo,
		userDailyRewardRepo: userDailyRewardRepo,
		userRepo:            userRepo,
	}
}

func (d *dailyRewardUseCase) CreateRewardType(ctx context.Context, data entities.RewardType) (res *entities.RewardType, err error) {
	id, err := d.dailyRewardRepo.CreateRewardTypeDB(ctx, data)
	if err != nil {
		return nil, err
	}

	if id == nil {
		return nil, apperror.ErrFailedRetrieveID
	}

	data.ID = *id

	_ = d.dailyRewardRepo.DeleteDailyRewardsRedis(ctx)

	return &data, err
}

func (d *dailyRewardUseCase) GetRewardTypes(ctx context.Context) (res []entities.RewardType, err error) {
	return d.dailyRewardRepo.GetRewardTypesDB(ctx)
}

func (d *dailyRewardUseCase) UpdateRewardTypes(ctx context.Context, id int64, data entities.RewardType) (res *entities.RewardType, err error) {
	err = d.dailyRewardRepo.UpdateRewardTypesDB(ctx, id, data)
	if err != nil {
		return nil, err
	}

	res, err = d.GetRewardTypeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = d.dailyRewardRepo.DeleteDailyRewardsRedis(ctx)

	return res, err
}

func (d *dailyRewardUseCase) GetRewardTypeByID(ctx context.Context, id int64) (res *entities.RewardType, err error) {
	res, err = d.dailyRewardRepo.GetRewardTypeByIDDB(ctx, id)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, apperror.ErrRecordNotFound
	}

	return res, err
}

func (d *dailyRewardUseCase) GetRewardTypeBySlug(ctx context.Context, slug string) (res *entities.RewardType, err error) {
	res, err = d.dailyRewardRepo.GetRewardTypeBySlugDB(ctx, slug)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, apperror.ErrRecordNotFound
	}

	return res, err
}

func (d *dailyRewardUseCase) CreateDailyReward(ctx context.Context, data entities.DailyReward, rewardTypeSlug string) (res *entities.DailyReward, err error) {
	rewardType, err := d.GetRewardTypeBySlug(ctx, rewardTypeSlug)
	if err != nil {
		return nil, err
	}

	if rewardType == nil {
		return nil, apperror.ErrRecordNotFound
	}

	data.RewardType = rewardType

	id, err := d.dailyRewardRepo.CreateDailyRewardDB(ctx, data)
	if err != nil {
		return nil, err
	}

	data.ID = *id

	_ = d.dailyRewardRepo.DeleteDailyRewardsRedis(ctx)

	return &data, err
}

func (d *dailyRewardUseCase) GetDailyRewards(ctx context.Context) (res []entities.DailyReward, err error) {
	res, err = d.dailyRewardRepo.GetDailyRewardsRedis(ctx)
	if err == nil && res != nil {
		return res, nil
	}

	rewards, err := d.dailyRewardRepo.GetDailyRewardsDB(ctx)
	if err != nil {
		return nil, err
	}

	go func(r []entities.DailyReward) {
		_ = d.dailyRewardRepo.SetDailyRewardsRedis(context.Background(), r)
	}(rewards)

	return rewards, err
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

func (d *dailyRewardUseCase) UpdateDailyReward(ctx context.Context, id int64, data entities.DailyReward, rewardTypeSlug string) (res *entities.DailyReward, err error) {
	rewardType, err := d.GetRewardTypeBySlug(ctx, rewardTypeSlug)
	if err != nil {
		return nil, err
	}

	if rewardType == nil {
		return nil, apperror.ErrRecordNotFound
	}

	data.RewardType = rewardType

	err = d.dailyRewardRepo.UpdateDailyRewardDB(ctx, id, data)
	if err != nil {
		return nil, err
	}

	res, err = d.GetDailyRewardID(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = d.dailyRewardRepo.DeleteDailyRewardsRedis(ctx)

	return res, err
}

func (d *dailyRewardUseCase) getDailyRewardByDay(ctx context.Context, day int64) (res *entities.DailyReward, err error) {
	allRewards, err := d.GetDailyRewards(ctx)
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

func (d *dailyRewardUseCase) GetRewardStatus(ctx context.Context) ([]entities.DailyReward, *int64, *bool, error) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil || userID <= 0 {
		return nil, nil, nil, apperror.ErrUnauthorized
	}

	progression, err := d.userUseCase.GetUserDailyRewardByID(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	rewards, err := d.GetDailyRewards(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	// initialize an empty progress object.
	if progression == nil {
		progression = &entities.UserDailyReward{CurrentStreak: 0, LastClaimDate: nil}
	}

	now := helper.NowUTC()
	today := now.Format("2006-01-02")

	// Default to true for new users. If LastClaimDate exists, check if today is a later date.
	isNewDay := true
	var dailyRewardIdx int64

	if progression.LastClaimDate != nil {
		lastClaimDay := progression.LastClaimDate.Format("2006-01-02")
		isNewDay = today > lastClaimDay
	}

	// Calculate the current day index in the 7-day cycle (0-6)
	// CurrentStreak tracks total days claimed (can exceed 7)
	// dailyRewardIdx maps this to a position in the 7-day cycle using modulo
	if isNewDay {
		// New day: user hasn't claimed yet today
		// Point to the next available day to claim
		// Example: CurrentStreak=1 (claimed Day 1) → dailyRewardIdx=1 (Day 2 available)
		dailyRewardIdx = progression.CurrentStreak % maxDailyRewardDays
	} else {
		// Same day: user already claimed earlier today
		// Point to the day they just claimed
		// Example: CurrentStreak=1 (just claimed Day 1) → dailyRewardIdx=0 (Day 1 shown as claimed)
		dailyRewardIdx = (progression.CurrentStreak - 1) % maxDailyRewardDays
		if progression.CurrentStreak == 0 {
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
		// Use CurrentStreak-1 if already claimed today to prevent premature cycle increment
		claimedDays := progression.CurrentStreak
		if !isNewDay {
			claimedDays = progression.CurrentStreak - 1
		}

		// Determine which cycle (0 = Days 1-7, 1 = Days 8-14, etc.)
		cycleNumber := claimedDays / maxDailyRewardDays

		// Add cycle offset to base day number
		rewardCopy[i].DayNumber = rewardCopy[i].DayNumber + (cycleNumber * maxDailyRewardDays)
	}

	return rewardCopy, &dailyRewardIdx, &isNewDay, err
}

func (d *dailyRewardUseCase) ClaimReward(ctx context.Context) (reward *entities.DailyReward, newBalance *entities.UserBalance, err error) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil || userID <= 0 {
		return nil, nil, apperror.ErrUnauthorized
	}

	progression, err := d.userUseCase.GetUserDailyRewardByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	now := helper.NowUTC()
	today := now.Format("2006-01-02")

	// Check if user already claimed today
	if progression != nil && progression.LastClaimDate != nil {
		lastClaimDate := progression.LastClaimDate.UTC().Format("2006-01-02")
		if today <= lastClaimDate {
			return nil, nil, apperror.ErrAlreadyClaimed
		}
	}

	// Initialize progression for new users who haven't claimed any rewards yet
	if progression == nil {
		progression = &entities.UserDailyReward{CurrentStreak: 0, LastClaimDate: nil}
	}

	// Calculate which day to claim in the 7-day cycle
	// CurrentStreak tracks total claims, rewardIdx maps to cycle position (0-6)
	rewardIdx := progression.CurrentStreak % maxDailyRewardDays
	dayToClaim := rewardIdx + 1

	reward, err = d.getDailyRewardByDay(ctx, dayToClaim)
	if err != nil {
		return nil, nil, err
	}

	// Update user progression
	err = d.dailyRewardRepo.DailyRewardWithTx(ctx, func(txRepo repositories.DailyRewardRepository) error {
		// Increment the streak counter after successful claim
		newStreak := progression.CurrentStreak + 1
		rawTx := txRepo.GetTx()

		userDailyRewardTx := d.userDailyRewardRepo.WithTx(rawTx)
		userRepoTx := d.userRepo.WithTx(rawTx)

		// Update user's streak and last claim date
		err = userDailyRewardTx.UpsertUserProgressionWithTx(ctx, userID, newStreak, now)
		if err != nil {
			return err
		}

		// Handle reward based on type
		if reward.RewardType.Slug.RequiresBalanceUpdate() {
			// For COIN and GEM, update user balance in database
			balanceType := reward.RewardType.Slug.ToUserBalance()
			err = userRepoTx.UpdateUserBalanceWithTx(ctx, userID, balanceType, reward.RewardAmount)
			if err != nil {
				return err
			}
		} else if reward.RewardType.Slug.IsSentExternally() {
			// TODO: Call External API to give Gopay Coin to player
		} else {
			return apperror.ErrUnknownRewardType
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	// Clear the cache so that the next request retrieves the latest progress data from the database
	_ = d.userDailyRewardRepo.DeleteUserDailyRewardRedis(ctx, userID)

	// Get updated balance for rewards that are stored in DB
	var userBalance *entities.UserBalance
	if reward.RewardType.Slug.RequiresBalanceUpdate() {
		userBalance, err = d.userUseCase.GetUserBalance(ctx, userID)
		if err != nil {
			return nil, nil, err
		}
	}

	return reward, userBalance, err
}
