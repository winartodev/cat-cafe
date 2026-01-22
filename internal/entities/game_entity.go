package entities

type GameStageStatus string

const (
	GSStatusCurrent  GameStageStatus = "current"
	GSStatusComplete GameStageStatus = "complete"
	GSStatusLocked   GameStageStatus = "locked"
)

type Game struct {
	DailyRewardAvailable bool         `json:"daily_reward_available"`
	UserBalance          *UserBalance `json:"user_balance"`
}

type UserGameStage struct {
	Slug     string          `json:"slug"`
	Name     string          `json:"name"`
	Sequence int64           `json:"sequence"`
	Status   GameStageStatus `json:"status"`
}

type UserNextGameStageInfo struct {
	Slug   string          `json:"slug"`
	Name   string          `json:"name"`
	Status GameStageStatus `json:"status"`
}
