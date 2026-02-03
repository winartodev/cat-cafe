package repositories

import (
	"context"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/database"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

type TutorialRepository interface {
	WithTx(tx *sql.Tx) TutorialRepository
	WithTutorialTx(ctx context.Context, fn func(tx *sql.Tx) error) error

	CreateTutorialSequenceDB(ctx context.Context, data entities.TutorialSequencesEntity) (id *int64, err error)
	BulkInsertTutorialTranslationsDB(ctx context.Context, sequenceID int64, data []entities.TutorialSequenceTranslationsEntity) (ids []int64, err error)
	GetTutorialsDB(ctx context.Context, limit, offset int) (res []entities.TutorialSequencesEntity, err error)
	GetTotalTutorialsDB(ctx context.Context) (res int64, err error)

	GetTutorialsByKeyDB(ctx context.Context, key string, limit, offset int) (res []entities.TutorialSequencesEntity, err error)
	GetTotalSequenceByKeyDB(ctx context.Context, key string) (totalRows int64, err error)
	UpdateTutorialSequenceDB(ctx context.Context, data entities.TutorialSequencesEntity) error
	DeleteTranslationsBySequenceIDDB(ctx context.Context, sequenceID int64) error
}

type tutorialRepository struct {
	BaseRepository
}

func NewTutorialRepository(db *sql.DB, redis *redis.Client) TutorialRepository {
	return &tutorialRepository{
		BaseRepository{
			db:    db,
			redis: redis,
			pool:  db,
		},
	}
}

func (r *tutorialRepository) WithTx(tx *sql.Tx) TutorialRepository {
	if tx == nil {
		return r
	}

	return &tutorialRepository{
		BaseRepository: BaseRepository{
			db:    tx,
			pool:  r.pool,
			redis: r.redis,
		},
	}
}

func (r *tutorialRepository) WithTutorialTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := r.pool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *tutorialRepository) CreateTutorialSequenceDB(ctx context.Context, data entities.TutorialSequencesEntity) (id *int64, err error) {
	var lastInsertID int64
	now := helper.NowUTC()
	err = r.db.QueryRowContext(ctx, createTutorialSequenceQuery,
		data.TutorialKey,
		data.Location,
		data.Sequence,
		now,
		now,
	).Scan(&lastInsertID)
	if database.IsDuplicateError(err) {
		return nil, apperror.ErrConflict
	} else if err != nil {
		return nil, err
	}

	return &lastInsertID, nil
}

func (r *tutorialRepository) BulkInsertTutorialTranslationsDB(ctx context.Context, sequenceID int64, data []entities.TutorialSequenceTranslationsEntity) ([]int64, error) {
	if len(data) == 0 {
		return nil, nil
	}

	numFields := 6
	queryString := r.BuildBulkInsertQuery(createTutorialTranslations, len(data), numFields, "RETURNING id")

	args := make([]interface{}, 0, len(data)*numFields)
	now := helper.NowUTC()

	for _, item := range data {
		args = append(args,
			sequenceID,
			item.LanguageCode,
			item.Title,
			item.Description,
			now,
			now,
		)
	}

	rows, err := r.db.QueryContext(ctx, queryString, args...)
	if database.IsDuplicateError(err) {
		return nil, apperror.ErrConflict
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *tutorialRepository) GetTutorialsDB(ctx context.Context, limit, offset int) (res []entities.TutorialSequencesEntity, err error) {
	rows, err := r.db.QueryContext(ctx, getTutorialsQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var row entities.TutorialSequencesEntity

		err := rows.Scan(
			&row.TutorialKey,
			&row.Location,
		)

		if err != nil {
			return nil, err
		}

		res = append(res, row)
	}

	return res, err
}

func (r *tutorialRepository) GetTutorialsByKeyDB(ctx context.Context, key string, limit, offset int) (res []entities.TutorialSequencesEntity, err error) {
	rows, err := r.db.QueryContext(ctx, getDetailTutorialsQuery, key, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lookup := make(map[int]*entities.TutorialSequencesEntity)
	var keys []int

	for rows.Next() {
		var seq entities.TutorialSequencesEntity
		var trans entities.TutorialSequenceTranslationsEntity

		err := rows.Scan(
			&seq.ID,
			&seq.TutorialKey,
			&seq.Location,
			&seq.Sequence,
			&trans.LanguageCode,
			&trans.Title,
			&trans.Description,
		)
		if err != nil {
			return nil, err
		}

		// Check if we already created this sequence in our map
		if existingSeq, ok := lookup[seq.Sequence]; ok {
			// Sequence exists, just append the new translation
			existingSeq.Translations = append(existingSeq.Translations, trans)
		} else {
			// New sequence found: Initialize the Translations slice and store it
			seq.Translations = []entities.TutorialSequenceTranslationsEntity{trans}
			lookup[seq.Sequence] = &seq
			keys = append(keys, seq.Sequence) // Track order
		}
	}

	// Convert the map back into an ordered slice
	for _, k := range keys {
		res = append(res, *lookup[k])
	}

	return res, nil
}

func (r *tutorialRepository) GetTotalTutorialsDB(ctx context.Context) (totalRows int64, err error) {
	var total int64
	countQuery := `SELECT COUNT(DISTINCT (tutorial_key, location)) FROM tutorial_sequences;`
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *tutorialRepository) GetTotalSequenceByKeyDB(ctx context.Context, key string) (totalRows int64, err error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM tutorial_sequences WHERE tutorial_key = $1`
	err = r.db.QueryRowContext(ctx, countQuery, key).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *tutorialRepository) UpdateTutorialSequenceDB(ctx context.Context, data entities.TutorialSequencesEntity) error {
	now := helper.NowUTC()
	_, err := r.db.ExecContext(ctx, updateTutorialSequenceQuery,
		data.Location,
		data.Sequence,
		now,
		data.ID,
	)
	return err
}

func (r *tutorialRepository) DeleteTranslationsBySequenceIDDB(ctx context.Context, sequenceID int64) error {
	_, err := r.db.ExecContext(ctx, deleteTutorialTranslationsQuery, sequenceID)
	return err
}
