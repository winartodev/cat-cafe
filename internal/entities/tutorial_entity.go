package entities

import "time"

type TutorialSequencesEntity struct {
	ID          *int64     `json:"id,omitempty"`
	TutorialKey string     `json:"tutorial_key"`
	Location    string     `json:"location"`
	Sequence    int        `json:"sequence"`
	CreatedAt   *time.Time `json:"-"`
	UpdatedAt   *time.Time `json:"-"`

	Translations []TutorialSequenceTranslationsEntity `json:"translations"`
}

type TutorialSequenceTranslationsEntity struct {
	ID                   *int64     `json:"id,omitempty"`
	TutorialSequencesKey string     `json:"tutorial_sequences_key"`
	LanguageCode         string     `json:"language_code"`
	Title                string     `json:"title"`
	Description          string     `json:"description"`
	CreatedAt            *time.Time `json:"-"`
	UpdatedAt            *time.Time `json:"-"`
}
