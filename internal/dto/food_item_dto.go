package dto

import "github.com/winartodev/cat-cafe/internal/entities"

type FoodItemRequest struct {
	Slug                string  `json:"slug"`
	Name                string  `json:"name"`
	StartingPrice       int64   `json:"starting_price"`
	StartingPreparation float64 `json:"starting_preparation"`
}

type FoodItemResponse struct {
	ID                  *int64  `json:"id,omitempty"`
	Slug                string  `json:"slug"`
	Name                string  `json:"name"`
	StartingPrice       int64   `json:"starting_price"`
	StartingPreparation float64 `json:"starting_preparation"`
}

func (r FoodItemRequest) ToEntity() entities.FoodItem {
	return entities.FoodItem{
		Slug:                r.Slug,
		Name:                r.Name,
		StartingPrice:       r.StartingPrice,
		StartingPreparation: r.StartingPreparation,
	}
}

func ToFoodItemResponse(data *entities.FoodItem) FoodItemResponse {
	var id *int64
	if data.ID > 0 {
		id = &data.ID
	}

	return FoodItemResponse{
		ID:                  id,
		Slug:                data.Slug,
		Name:                data.Name,
		StartingPrice:       data.StartingPrice,
		StartingPreparation: data.StartingPreparation,
	}
}

func ToFoodItemsResponse(data []entities.FoodItem) []FoodItemResponse {
	res := make([]FoodItemResponse, 0)
	for _, e := range data {
		res = append(res, ToFoodItemResponse(&e))
	}

	return res
}
