package entities

type Game struct {
	DailyRewardAvailable bool         `json:"daily_reward_available"`
	UserBalance          *UserBalance `json:"user_balance"`
}
