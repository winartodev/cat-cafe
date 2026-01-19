package repositories

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
)

type BaseRepository struct {
	db    *sql.DB
	tx    *sql.Tx
	redis *redis.Client
}

func (r *BaseRepository) GetTx() *sql.Tx {
	return r.tx
}

type Repository struct {
	RewardRepository          RewardRepository
	DailyRewardRepository     DailyRewardRepository
	UserRepository            UserRepository
	UserDailyRewardRepository UserDailyRewardRepository
}

func SetupRepository(db *sql.DB, client *redis.Client) *Repository {
	return &Repository{
		RewardRepository:          NewRewardRepository(db, client),
		DailyRewardRepository:     NewDailyRewardsRepository(db, client),
		UserRepository:            NewUserRepository(db, client),
		UserDailyRewardRepository: NewUserDailyRewardRepository(db, client),
	}
}
