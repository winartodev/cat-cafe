package entities

import "time"

type FoodItem struct {
	ID                  int64      `json:"id"`
	Slug                string     `json:"slug"`
	Name                string     `json:"name"`
	StartingPrice       int64      `json:"starting_price"`
	StartingPreparation float64    `json:"starting_preparation"`
	CreatedAt           *time.Time `json:"-"`
	UpdatedAt           *time.Time `json:"-"`
}
