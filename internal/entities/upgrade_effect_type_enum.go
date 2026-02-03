package entities

import (
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"strings"
)

type UpgradeEffectType string

const (
	UpgradeEffectTypeAddHelper         UpgradeEffectType = "add_helper"
	UpgradeEffectTypeAddCustomer       UpgradeEffectType = "add_customer"
	UpgradeEffectTypeReduceCookingTime UpgradeEffectType = "reduce_cooking_time"
	UpgradeEffectTypeProfit            UpgradeEffectType = "profit"
	UpgradeEffectTypeUnlockRestaurant  UpgradeEffectType = "unlock_restaurant"
)

func (u UpgradeEffectType) String() string {
	return string(u)
}

func (u UpgradeEffectType) ToLower() string {
	return strings.ToLower(string(u))
}

func (u UpgradeEffectType) IsValid() bool {
	switch u {
	case UpgradeEffectTypeAddHelper,
		UpgradeEffectTypeAddCustomer,
		UpgradeEffectTypeReduceCookingTime,
		UpgradeEffectTypeProfit,
		UpgradeEffectTypeUnlockRestaurant:
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
		UpgradeEffectTypeReduceCookingTime,
		UpgradeEffectTypeProfit,
		UpgradeEffectTypeUnlockRestaurant,
	}
}
