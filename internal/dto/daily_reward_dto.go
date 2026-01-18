package dto

import "github.com/winartodev/cat-cafe/internal/entities"

type CreateRewardTypeRequest struct {
	Slug entities.RewardTypeSlug `json:"slug"`
	Name string                  `json:"name"`
}

type UpdateRewardTypeRequest struct {
	Name string `json:"name"`
}

type RewardTypeResponse struct {
	ID   *int64                  `json:"id,omitempty"`
	Slug entities.RewardTypeSlug `json:"slug"`
	Name string                  `json:"name"`
}

type DailyRewardRequest struct {
	DayNumber    int64  `json:"day_number"`
	RewardType   string `json:"reward_type"`
	RewardAmount int64  `json:"reward_amount"`
	IsActive     bool   `json:"is_active"`
	Description  string `json:"description"`
}

type DailyRewardResponse struct {
	ID           int64                 `json:"id"`
	DayNumber    int64                 `json:"day_number"`
	RewardAmount int64                 `json:"reward_amount"`
	IsActive     bool                  `json:"is_active"`
	Status       entities.RewardStatus `json:"status,omitempty"`
	Description  string                `json:"description"`
	RewardType   *RewardTypeResponse   `json:"reward_type,omitempty"`
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
		ID:   &entity.ID,
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

func (e *DailyRewardRequest) ToEntity() *entities.DailyReward {
	return &entities.DailyReward{
		DayNumber:    e.DayNumber,
		RewardAmount: e.RewardAmount,
		IsActive:     e.IsActive,
		Description:  e.Description,
	}
}

func ToDailyRewardResponse(dailyReward *entities.DailyReward) *DailyRewardResponse {
	if dailyReward == nil {
		return nil
	}

	rewardTypeResp := ToRewardTypeResponse(dailyReward.RewardType)
	rewardTypeResp.ID = nil

	return &DailyRewardResponse{
		ID:           dailyReward.ID,
		RewardType:   rewardTypeResp,
		DayNumber:    dailyReward.DayNumber,
		RewardAmount: dailyReward.RewardAmount,
		IsActive:     dailyReward.IsActive,
		Description:  dailyReward.Description,
		Status:       dailyReward.Status,
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
