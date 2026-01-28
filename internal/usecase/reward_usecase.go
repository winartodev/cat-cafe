package usecase

import (
	"context"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
)

type RewardUseCase interface {
	CreateRewardType(ctx context.Context, data entities.RewardType) (res *entities.RewardType, err error)
	UpdateRewardTypes(ctx context.Context, id int64, data entities.RewardType) (res *entities.RewardType, err error)
	GetRewardTypeBySlug(ctx context.Context, slug string) (res *entities.RewardType, err error)
	GetRewardTypes(ctx context.Context) (res []entities.RewardType, err error)
	GetRewardTypeByID(ctx context.Context, id int64) (res *entities.RewardType, err error)

	CreateReward(ctx context.Context, data entities.Reward) (res *entities.Reward, err error)
	GetRewards(ctx context.Context, limit, offset int) (res []entities.Reward, totalRow int64, err error)
	GetRewardBySlug(ctx context.Context, slug string) (res *entities.Reward, err error)
}

type rewardUseCase struct {
	rewardRepo repositories.RewardRepository
}

func NewRewardUseCase(rewardRepo repositories.RewardRepository) RewardUseCase {
	return &rewardUseCase{
		rewardRepo: rewardRepo,
	}
}

func (r *rewardUseCase) CreateRewardType(ctx context.Context, data entities.RewardType) (res *entities.RewardType, err error) {
	id, err := r.rewardRepo.CreateRewardTypeDB(ctx, data)
	if err != nil {
		return nil, err
	}

	if id == nil {
		return nil, apperror.ErrFailedRetrieveID
	}

	data.ID = id

	return &data, err
}

func (r *rewardUseCase) UpdateRewardTypes(ctx context.Context, id int64, data entities.RewardType) (res *entities.RewardType, err error) {
	err = r.rewardRepo.UpdateRewardTypesDB(ctx, id, data)
	if err != nil {
		return nil, err
	}

	res, err = r.GetRewardTypeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (r *rewardUseCase) GetRewardTypes(ctx context.Context) (res []entities.RewardType, err error) {
	return r.rewardRepo.GetRewardTypesDB(ctx)
}

func (r *rewardUseCase) GetRewardTypeByID(ctx context.Context, id int64) (res *entities.RewardType, err error) {
	res, err = r.rewardRepo.GetRewardTypeByIDDB(ctx, id)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, apperror.ErrRecordNotFound
	}

	return res, err
}

func (r *rewardUseCase) GetRewardTypeBySlug(ctx context.Context, slug string) (res *entities.RewardType, err error) {
	res, err = r.rewardRepo.GetRewardTypeBySlugDB(ctx, slug)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, apperror.ErrRecordNotFound
	}

	return res, err
}

func (r *rewardUseCase) CreateReward(ctx context.Context, data entities.Reward) (res *entities.Reward, err error) {
	rewardType, err := r.GetRewardTypeBySlug(ctx, data.RewardType.Slug)
	if err != nil {
		return nil, err
	}

	if rewardType == nil {
		return nil, apperror.ErrRecordNotFound
	}

	data.RewardType = rewardType

	id, err := r.rewardRepo.CreateRewardDB(ctx, data)
	if err != nil {
		return nil, err
	}

	if id == nil {
		return nil, apperror.ErrFailedRetrieveID
	}

	data.ID = *id

	return &data, nil
}

func (r *rewardUseCase) GetRewardBySlug(ctx context.Context, slug string) (res *entities.Reward, err error) {
	reward, err := r.rewardRepo.GetRewardBySlugDB(ctx, slug)
	if err != nil {
		return nil, err
	}

	if reward == nil {
		return nil, apperror.ErrRecordNotFound
	}

	return reward, nil
}

func (r *rewardUseCase) GetRewards(ctx context.Context, limit, offset int) (res []entities.Reward, totalRow int64, err error) {
	res, err = r.rewardRepo.GetRewardsDB(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	totalRow, err = r.rewardRepo.CountRewardDB(ctx)
	if err != nil {
		return nil, 0, err
	}

	return res, totalRow, err
}
