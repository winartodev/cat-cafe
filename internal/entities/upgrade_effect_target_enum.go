package entities

import "github.com/winartodev/cat-cafe/pkg/apperror"

type UpgradeEffectTarget string

const (
	UpgradeEffectTargetAllFood    UpgradeEffectTarget = "all_food"
	UpgradeEffectTargetFood       UpgradeEffectTarget = "food"
	UpgradeEffectTargetHelper     UpgradeEffectTarget = "helper"
	UpgradeEffectTargetCustomer   UpgradeEffectTarget = "customer"
	UpgradeEffectTargetRestaurant UpgradeEffectTarget = "restaurant"
)

func (u UpgradeEffectTarget) String() string {
	return string(u)
}

func (u UpgradeEffectTarget) IsValid() bool {
	switch u {
	case UpgradeEffectTargetAllFood,
		UpgradeEffectTargetFood,
		UpgradeEffectTargetHelper,
		UpgradeEffectTargetCustomer,
		UpgradeEffectTargetRestaurant:
		return true
	}
	return false
}

func ParseUpgradeEffectTarget(s string) (UpgradeEffectTarget, error) {
	target := UpgradeEffectTarget(s)
	if !target.IsValid() {
		return "", apperror.ErrorInvalidRequest("upgrade target type:", s)
	}
	return target, nil
}

func AllUpgradeEffectTarget() []UpgradeEffectTarget {
	return []UpgradeEffectTarget{
		UpgradeEffectTargetAllFood,
		UpgradeEffectTargetFood,
		UpgradeEffectTargetHelper,
		UpgradeEffectTargetCustomer,
		UpgradeEffectTargetRestaurant,
	}
}
