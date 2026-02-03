package dto

import "github.com/winartodev/cat-cafe/internal/entities"

type TutorialDTO struct {
	ID           *int64           `json:"id,omitempty"`
	TutorialKey  string           `json:"tutorial_key"`
	Location     string           `json:"location"`
	Sequence     int              `json:"sequence"`
	Translations []TranslationDTO `json:"translations,omitempty"`
}

type TranslationDTO struct {
	ID           *int64 `json:"id,omitempty"`
	LanguageCode string `json:"language_code"`
	Title        string `json:"title"`
	Description  string `json:"description"`
}

func (t *TutorialDTO) ToEntity() (entities.TutorialSequencesEntity, []entities.TutorialSequenceTranslationsEntity) {
	sequence := entities.TutorialSequencesEntity{
		ID:          t.ID,
		TutorialKey: t.TutorialKey,
		Location:    t.Location,
		Sequence:    t.Sequence,
	}

	var sequenceTranslations []entities.TutorialSequenceTranslationsEntity
	for _, translation := range t.Translations {
		sequenceTranslations = append(sequenceTranslations, entities.TutorialSequenceTranslationsEntity{
			LanguageCode: translation.LanguageCode,
			Title:        translation.Title,
			Description:  translation.Description,
		})
	}

	return sequence, sequenceTranslations
}

func ToCreateTutorialResponse(sequence *entities.TutorialSequencesEntity, translations []entities.TutorialSequenceTranslationsEntity) *TutorialDTO {
	if sequence == nil {
		return nil
	}

	var translationsDto []TranslationDTO
	for _, t := range translations {
		translationsDto = append(translationsDto, TranslationDTO{
			ID:           t.ID,
			LanguageCode: t.LanguageCode,
			Title:        t.Title,
			Description:  t.Description,
		})
	}

	return &TutorialDTO{
		ID:           sequence.ID,
		TutorialKey:  sequence.TutorialKey,
		Location:     sequence.Location,
		Sequence:     sequence.Sequence,
		Translations: translationsDto,
	}
}

func ToDetailTutorialsResponse(data []entities.TutorialSequencesEntity) []TutorialDTO {
	if data == nil {
		return nil
	}

	var tutorials []TutorialDTO
	for _, t := range data {
		tutorial := TutorialDTO{
			ID:          t.ID,
			TutorialKey: t.TutorialKey,
			Location:    t.Location,
			Sequence:    t.Sequence,
		}

		if t.Translations != nil {
			var translationsDto []TranslationDTO
			for _, tr := range t.Translations {
				translationsDto = append(translationsDto, TranslationDTO{
					ID:           tr.ID,
					LanguageCode: tr.LanguageCode,
					Title:        tr.Title,
					Description:  tr.Description,
				})
			}

			tutorial.Translations = translationsDto
		}

		tutorials = append(tutorials, tutorial)
	}

	return tutorials
}
