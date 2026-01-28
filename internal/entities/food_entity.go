package entities

import "time"

type FoodItem struct {
	ID            int64      `json:"id"`
	Slug          string     `json:"slug"`
	Name          string     `json:"name"`
	InitialCost   int64      `json:"initial_cost"`
	InitialProfit int64      `json:"initial_profit"`
	CookingTime   float64    `json:"cooking_time"`
	CreatedAt     *time.Time `json:"-"`
	UpdatedAt     *time.Time `json:"-"`
	UseOverride   bool       `json:"use_override"`
}

type FoodItemOverrideLevel struct {
	ID              int64      `json:"id"`
	FoodItemID      int64      `json:"food_item_id"`
	Level           int64      `json:"level"`
	Cost            int64      `json:"cost"`
	Profit          int64      `json:"profit"`
	PreparationTime float64    `json:"preparation_time"`
	CreatedAt       *time.Time `json:"-"`
	UpdatedAt       *time.Time `json:"-"`
}
