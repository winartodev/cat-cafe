package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

type DbTx interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type BaseRepository struct {
	db    DbTx
	pool  *sql.DB
	redis *redis.Client
}

func (r *BaseRepository) BuildBulkInsertQuery(baseQuery string, numItems int, numFields int, returningClause string) string {
	var queryBuilder strings.Builder
	queryBuilder.WriteString(baseQuery)

	placeholderCount := 1
	for i := 0; i < numItems; i++ {
		if i > 0 {
			queryBuilder.WriteString(", ")
		}
		queryBuilder.WriteString("(")
		for j := 0; j < numFields; j++ {
			queryBuilder.WriteString(fmt.Sprintf("$%d", placeholderCount))
			if j < numFields-1 {
				queryBuilder.WriteString(", ")
			}
			placeholderCount++
		}
		queryBuilder.WriteString(")")
	}

	if returningClause != "" {
		queryBuilder.WriteString(" " + returningClause)
	}

	queryBuilder.WriteString(";")
	return queryBuilder.String()
}

type Repository struct {
	RewardRepository              RewardRepository
	DailyRewardRepository         DailyRewardRepository
	UserRepository                UserRepository
	UserProgressionRepository     UserProgressionRepository
	GameStageRepository           GameStageRepository
	StageCustomerConfigRepository StageCustomerConfigRepository
	StageStaffConfigRepository    StageStaffConfigRepository
	StageKitchenConfigRepository  StageKitchenConfigRepository
	StageCameraConfigRepository   StageCameraConfigRepository
	FoodItemRepository            FoodItemRepository
	KitchenStationRepository      KitchenStationRepository
	UpgradeRepository             UpgradeRepository
	StageUpgradeRepository        StageUpgradeRepository
}

func SetupRepository(db *sql.DB, client *redis.Client) *Repository {
	return &Repository{
		RewardRepository:              NewRewardRepository(db, client),
		DailyRewardRepository:         NewDailyRewardsRepository(db, client),
		UserRepository:                NewUserRepository(db, client),
		UserProgressionRepository:     NewUserProgressionRepository(db, client),
		GameStageRepository:           NewGameStageRepository(db),
		StageCustomerConfigRepository: NewStageCustomerRepository(db),
		StageStaffConfigRepository:    NewStageStaffConfigRepository(db),
		StageKitchenConfigRepository:  NewStageKitchenConfigRepository(db),
		StageCameraConfigRepository:   NewStageCameraConfigRepository(db),
		FoodItemRepository:            NewFoodItemRepository(db),
		KitchenStationRepository:      NewKitchenStationRepository(db),
		UpgradeRepository:             NewUpgradeRepository(db),
		StageUpgradeRepository:        NewStageUpgradeRepository(db),
	}
}
