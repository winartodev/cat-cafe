package dto

import "github.com/winartodev/cat-cafe/internal/entities"

type FoodItemRequest struct {
	Slug           string                  `json:"slug"`
	Name           string                  `json:"name"`
	InitialCost    int64                   `json:"initial_cost"`
	InitialProfit  int64                   `json:"initial_profit"`
	CookingTime    float64                 `json:"cooking_time"`
	OverrideLevels []FoodItemOverrideLevel `json:"override_levels"`
}

type FoodItemResponse struct {
	ID            *int64  `json:"id,omitempty"`
	Slug          string  `json:"slug"`
	Name          string  `json:"name"`
	InitialCost   int64   `json:"initial_cost"`
	InitialProfit int64   `json:"initial_profit"`
	CookingTime   float64 `json:"cooking_time"`

	OverrideLevels []FoodItemOverrideLevel `json:"override_levels,omitempty"`
}

type FoodItemOverrideLevel struct {
	Level       int64   `json:"level"`
	Cost        int64   `json:"cost"`
	Profit      int64   `json:"profit"`
	CookingTime float64 `json:"cooking_time"`
}

func (r FoodItemRequest) ToEntity() (entities.FoodItem, []entities.FoodItemOverrideLevel) {
	foodItem := entities.FoodItem{
		Slug:          r.Slug,
		Name:          r.Name,
		InitialCost:   r.InitialCost,
		InitialProfit: r.InitialProfit,
		CookingTime:   r.CookingTime,
	}

	if len(r.OverrideLevels) == 0 {
		return foodItem, nil
	}

	var overrideLevels []entities.FoodItemOverrideLevel
	for _, overrideLevel := range r.OverrideLevels {
		overrideLevels = append(overrideLevels, overrideLevel.ToEntity())
	}

	return foodItem, overrideLevels
}

func (r FoodItemOverrideLevel) ToEntity() entities.FoodItemOverrideLevel {
	return entities.FoodItemOverrideLevel{
		Level:           r.Level,
		Cost:            r.Cost,
		Profit:          r.Profit,
		PreparationTime: r.CookingTime,
	}
}

func ToFoodItemResponse(data *entities.FoodItem, overrideLevels []entities.FoodItemOverrideLevel) FoodItemResponse {
	var id *int64
	if data.ID > 0 {
		id = &data.ID
	}

	var overrideLevelResponses []FoodItemOverrideLevel
	for _, overrideLevel := range overrideLevels {
		overrideLevelResponses = append(overrideLevelResponses, FoodItemOverrideLevel{
			Level:       overrideLevel.Level,
			Cost:        overrideLevel.Cost,
			Profit:      overrideLevel.Profit,
			CookingTime: overrideLevel.PreparationTime,
		})
	}

	return FoodItemResponse{
		ID:             id,
		Slug:           data.Slug,
		Name:           data.Name,
		InitialCost:    data.InitialCost,
		InitialProfit:  data.InitialProfit,
		CookingTime:    data.CookingTime,
		OverrideLevels: overrideLevelResponses,
	}
}

func ToFoodItemsResponse(data []entities.FoodItem) []FoodItemResponse {
	res := make([]FoodItemResponse, 0)
	for _, e := range data {
		res = append(res, ToFoodItemResponse(&e, nil))
	}

	return res
}
