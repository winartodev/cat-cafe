package entities

import "time"

type Upgrade struct {
	ID          int64
	Slug        string
	Name        string
	Description string
	Cost        int64
	CostType    UpgradeCostType
	Effect      UpgradeEffect
	IsActive    bool
	Sequence    int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

type UpgradeEffect struct {
	Type       UpgradeEffectType
	Value      float64
	Unit       UpgradeEffectUnit
	Target     UpgradeEffectTarget
	TargetID   int64
	TargetName string
}

func (e *UpgradeEffect) CalculateNewValue(current float64) float64 {
	isReduction := e.Type == UpgradeEffectTypeReduceCookingTime

	switch e.Unit {
	case UpgradeEffectUnitPercentage:
		if isReduction {
			return current * (1.0 - (e.Value / 100.0))
		}
		return current * (1.0 + (e.Value / 100.0))

	case UpgradeEffectUnitMultiplier:
		return current * e.Value

	case UpgradeEffectUnitSeconds:
		if isReduction {
			return current - e.Value
		}
		return current + e.Value

	default:
		return current + e.Value
	}
}
