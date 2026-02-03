package entities

import "github.com/winartodev/cat-cafe/pkg/apperror"

type UpgradeEffectUnit string

const (
	UpgradeEffectUnitCount      UpgradeEffectUnit = "count"
	UpgradeEffectUnitMultiplier UpgradeEffectUnit = "multiplier"
	UpgradeEffectUnitSeconds    UpgradeEffectUnit = "seconds"
	UpgradeEffectUnitPercentage UpgradeEffectUnit = "percent"
)

func (u UpgradeEffectUnit) String() string {
	return string(u)
}

func (u UpgradeEffectUnit) IsValid() bool {
	switch u {
	case UpgradeEffectUnitCount,
		UpgradeEffectUnitMultiplier,
		UpgradeEffectUnitSeconds,
		UpgradeEffectUnitPercentage:
		return true
	}
	return false
}

func ParseUpgradeEffectUnit(s string) (UpgradeEffectUnit, error) {
	target := UpgradeEffectUnit(s)
	if !target.IsValid() {
		return "", apperror.ErrorInvalidRequest("upgrade effect unit:", s)
	}
	return target, nil
}

func AllUpgradeEffectUnit() []UpgradeEffectUnit {
	return []UpgradeEffectUnit{
		UpgradeEffectUnitCount,
		UpgradeEffectUnitMultiplier,
		UpgradeEffectUnitSeconds,
		UpgradeEffectUnitPercentage,
	}
}
