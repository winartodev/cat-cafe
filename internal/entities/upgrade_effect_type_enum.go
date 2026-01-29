package entities

import "github.com/winartodev/cat-cafe/pkg/apperror"

type UpgradeEffectType string

const (
	UpgradeEffectTypeAddHelper     UpgradeEffectType = "add_helper"
	UpgradeEffectTypeAddCustomer   UpgradeEffectType = "add_customer"
	UpgradeEffectTypeCookingTime   UpgradeEffectType = "cooking_time"
	UpgradeEffectTypeProfit        UpgradeEffectType = "profit"
	UpgradeEffectTypeNewRestaurant UpgradeEffectType = "new_restaurant"
)

func (u UpgradeEffectType) String() string {
	return string(u)
}

func (u UpgradeEffectType) IsValid() bool {
	switch u {
	case UpgradeEffectTypeAddHelper,
		UpgradeEffectTypeAddCustomer,
		UpgradeEffectTypeCookingTime,
		UpgradeEffectTypeProfit,
		UpgradeEffectTypeNewRestaurant:
		return true
	}
	return false
}

func ParseUpgradeEffectType(s string) (UpgradeEffectType, error) {
	target := UpgradeEffectType(s)
	if !target.IsValid() {
		return "", apperror.ErrorInvalidRequest("upgrade effect type:", s)
	}
	return target, nil
}

func AllUpgradeEffectType() []UpgradeEffectType {
	return []UpgradeEffectType{
		UpgradeEffectTypeAddHelper,
		UpgradeEffectTypeAddCustomer,
		UpgradeEffectTypeCookingTime,
		UpgradeEffectTypeProfit,
		UpgradeEffectTypeNewRestaurant,
	}
}
