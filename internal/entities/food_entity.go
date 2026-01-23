package entities

import "time"

type FoodItem struct {
	ID                  int64
	Slug                string
	Name                string
	StartingPrice       int64
	StartingPreparation float64
	CreatedAt           *time.Time
	UpdatedAt           *time.Time
}
