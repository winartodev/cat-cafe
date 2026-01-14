package repositories

import "database/sql"

type BaseRepository struct {
	db *sql.DB
	tx *sql.Tx
}

func (r *BaseRepository) GetTx() *sql.Tx {
	return r.tx
}

type Repository struct {
	DailyRewardRepository     DailyRewardRepository
	UserRepository            UserRepository
	UserDailyRewardRepository UserDailyRewardRepository
}

func SetupRepository(db *sql.DB) *Repository {
	return &Repository{
		DailyRewardRepository:     NewDailyRewardsRepository(db),
		UserRepository:            NewUserRepository(db),
		UserDailyRewardRepository: NewUserDailyRewardRepository(db),
	}
}
