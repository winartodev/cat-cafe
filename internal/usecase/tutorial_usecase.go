package usecase

import (
	"context"
	"database/sql"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
)

type TutorialUseCase interface {
	CreateTutorial(ctx context.Context, sequence entities.TutorialSequencesEntity, tutorials []entities.TutorialSequenceTranslationsEntity) (*entities.TutorialSequencesEntity, []entities.TutorialSequenceTranslationsEntity, error)
	GetTutorials(ctx context.Context, limit, offset int) (res []entities.TutorialSequencesEntity, totalRows int64, err error)
	UpdateTutorial(ctx context.Context, sequence entities.TutorialSequencesEntity, tutorials []entities.TutorialSequenceTranslationsEntity) error

	GetTranslationsByTutorialKey(ctx context.Context, key string, limit, offset int) (res []entities.TutorialSequencesEntity, totalRows int64, err error)
	GetTranslationByID(ctx context.Context, key string, id int64) (res *entities.TutorialSequencesEntity, err error)
}

type tutorialUseCase struct {
	tutorialRepo repositories.TutorialRepository
}

func NewTutorialUseCase(tutorialRepo repositories.TutorialRepository) TutorialUseCase {
	return &tutorialUseCase{
		tutorialRepo: tutorialRepo,
	}
}

func (t *tutorialUseCase) CreateTutorial(ctx context.Context, sequence entities.TutorialSequencesEntity, tutorials []entities.TutorialSequenceTranslationsEntity) (*entities.TutorialSequencesEntity, []entities.TutorialSequenceTranslationsEntity, error) {
	err := t.tutorialRepo.WithTutorialTx(ctx, func(tx *sql.Tx) error {
		tutorialRepoTx := t.tutorialRepo.WithTx(tx)
		seqID, err := tutorialRepoTx.CreateTutorialSequenceDB(ctx, sequence)
		if err != nil {
			return err
		}

		if seqID == nil {
			return apperror.ErrFailedRetrieveID
		}

		sequence.ID = seqID
		_, err = tutorialRepoTx.BulkInsertTutorialTranslationsDB(ctx, *seqID, tutorials)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return &sequence, tutorials, nil
}

func (t *tutorialUseCase) GetTutorials(ctx context.Context, limit, offset int) (res []entities.TutorialSequencesEntity, totalRows int64, err error) {
	res, err = t.tutorialRepo.GetTutorialsDB(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	totalRows, err = t.tutorialRepo.GetTotalTutorialsDB(ctx)
	if err != nil {
		return nil, 0, err
	}

	return res, totalRows, nil
}

func (t *tutorialUseCase) GetTranslationsByTutorialKey(ctx context.Context, key string, limit, offset int) (res []entities.TutorialSequencesEntity, totalRows int64, err error) {
	res, err = t.tutorialRepo.GetTutorialsByKeyDB(ctx, key, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	totalRows, err = t.tutorialRepo.GetTotalSequenceByKeyDB(ctx, key)
	if err != nil {
		return nil, 0, err
	}

	return res, totalRows, nil
}

func (t *tutorialUseCase) UpdateTutorial(ctx context.Context, sequence entities.TutorialSequencesEntity, tutorials []entities.TutorialSequenceTranslationsEntity) error {
	// We wrap everything in a transaction to ensure data integrity
	return t.tutorialRepo.WithTutorialTx(ctx, func(tx *sql.Tx) error {
		tutorialRepoTx := t.tutorialRepo.WithTx(tx)

		// 1. Update the parent sequence (location, sequence order, etc.)
		// Ensure sequence.ID is populated from the request body
		err := tutorialRepoTx.UpdateTutorialSequenceDB(ctx, sequence)
		if err != nil {
			return err
		}

		// 2. Delete all existing translations for this specific sequence ID
		// This is cleaner than trying to match/update individual translation IDs
		err = tutorialRepoTx.DeleteTranslationsBySequenceIDDB(ctx, *sequence.ID)
		if err != nil {
			return err
		}

		// 3. Re-insert the fresh set of translations
		_, err = tutorialRepoTx.BulkInsertTutorialTranslationsDB(ctx, *sequence.ID, tutorials)
		if err != nil {
			return err
		}

		return nil
	})
}

func (t *tutorialUseCase) GetTranslationByID(ctx context.Context, key string, id int64) (res *entities.TutorialSequencesEntity, err error) {
	return nil, nil
}
