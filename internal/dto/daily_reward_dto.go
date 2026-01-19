package dto

import (
	"github.com/winartodev/cat-cafe/internal/entities"
)

type CreateRewardTypeRequest struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type UpdateRewardTypeRequest struct {
	Name string `json:"name"`
}

type RewardTypeResponse struct {
	ID   *int64 `json:"id,omitempty"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type CreateRewardRequest struct {
	Slug       string `json:"slug"`
	Name       string `json:"name"`
	RewardType string `json:"reward_type"`
	Amount     int64  `json:"amount"`
	IsActive   bool   `json:"is_active"`
}

type UpdateRewardRequest struct {
	Name       string `json:"name"`
	RewardType string `json:"reward_type"`
	Amount     string `json:"amount"`
	IsActive   bool   `json:"is_active"`
}

type RewardResponse struct {
	ID         *int64               `json:"id,omitempty"`
	Slug       string               `json:"slug"`
	Name       string               `json:"name"`
	Amount     int64                `json:"amount"`
	IsActive   bool                 `json:"is_active"`
	RewardType *entities.RewardType `json:"reward_type"`
}

type DailyRewardRequest struct {
	DayNumber   int64  `json:"day_number"`
	Reward      string `json:"reward"`
	IsActive    bool   `json:"is_active"`
	Description string `json:"description"`
}

type DailyRewardResponse struct {
	ID          int64                 `json:"id"`
	DayNumber   int64                 `json:"day_number"`
	IsActive    bool                  `json:"is_active"`
	Status      entities.RewardStatus `json:"status,omitempty"`
	Description string                `json:"description"`
	Reward      *RewardResponse       `json:"reward,omitempty"`
}

type DailyRewardStatus struct {
	CurrentDailyRewardIdx int64                 `json:"current_daily_reward_idx"`
	IsNewDay              bool                  `json:"is_new_day"`
	Rewards               []DailyRewardResponse `json:"rewards"`
}

type ClaimDailyRewardResponse struct {
	Reward  *DailyRewardResponse `json:"reward"`
	Balance *UserBalanceResponse `json:"balance,omitempty"`
}

func (e *CreateRewardTypeRequest) ToEntity() *entities.RewardType {
	return &entities.RewardType{
		Slug: e.Slug,
		Name: e.Name,
	}
}

func (e *UpdateRewardTypeRequest) ToEntity() *entities.RewardType {
	return &entities.RewardType{
		Name: e.Name,
	}
}

func ToRewardTypeResponse(entity *entities.RewardType) *RewardTypeResponse {
	if entity == nil {
		return nil
	}

	return &RewardTypeResponse{
		ID:   entity.ID,
		Slug: entity.Slug,
		Name: entity.Name,
	}
}

func ToRewardTypeResponses(entities []entities.RewardType) []RewardTypeResponse {
	res := make([]RewardTypeResponse, 0)
	for _, e := range entities {
		res = append(res, *ToRewardTypeResponse(&e))
	}

	return res
}

func (e *CreateRewardRequest) ToEntity() entities.Reward {
	return entities.Reward{
		Slug:     e.Slug,
		Name:     e.Name,
		Amount:   e.Amount,
		IsActive: e.IsActive,
		RewardType: &entities.RewardType{
			Slug: e.RewardType,
		},
	}
}

func ToRewardResponse(data *entities.Reward) RewardResponse {
	var rewardType *entities.RewardType
	if data.RewardType != nil {
		rewardType = &entities.RewardType{
			ID:   nil,
			Slug: data.RewardType.Slug,
			Name: data.RewardType.Name,
		}
	}

	return RewardResponse{
		ID:         &data.ID,
		Slug:       data.Slug,
		Name:       data.Name,
		Amount:     data.Amount,
		IsActive:   data.IsActive,
		RewardType: rewardType,
	}
}

func ToRewardsResponse(data []entities.Reward) []RewardResponse {
	res := make([]RewardResponse, 0)
	for _, e := range data {
		res = append(res, ToRewardResponse(&e))
	}

	return res
}

func (e *DailyRewardRequest) ToEntity() *entities.DailyReward {
	return &entities.DailyReward{
		DayNumber:   e.DayNumber,
		IsActive:    e.IsActive,
		Description: e.Description,
	}
}

func ToDailyRewardResponse(dailyReward *entities.DailyReward) *DailyRewardResponse {
	if dailyReward == nil {
		return nil
	}

	reward := ToRewardResponse(dailyReward.Reward)
	reward.ID = nil

	return &DailyRewardResponse{
		ID:          dailyReward.ID,
		Reward:      &reward,
		DayNumber:   dailyReward.DayNumber,
		IsActive:    dailyReward.IsActive,
		Description: dailyReward.Description,
		Status:      dailyReward.Status,
	}
}

func ToDailyRewardResponses(entities []entities.DailyReward) []DailyRewardResponse {
	res := make([]DailyRewardResponse, 0)
	for _, e := range entities {
		res = append(res, *ToDailyRewardResponse(&e))
	}

	return res
}

func ToDailyRewardStatus(rewards []entities.DailyReward, dailyRewardIdx *int64, isNewDay *bool) DailyRewardStatus {
	return DailyRewardStatus{
		CurrentDailyRewardIdx: *dailyRewardIdx,
		IsNewDay:              *isNewDay,
		Rewards:               ToDailyRewardResponses(rewards),
	}
}

func ToClaimDailyRewardResponse(reward *entities.DailyReward, balance *entities.UserBalance) *ClaimDailyRewardResponse {
	if reward == nil {
		return nil
	}

	var userBalance *UserBalanceResponse
	if balance != nil {
		userBalance = &UserBalanceResponse{
			Coin: balance.Coin,
			Gem:  balance.Gem,
		}
	}

	return &ClaimDailyRewardResponse{
		Reward:  ToDailyRewardResponse(reward),
		Balance: userBalance,
	}
}
