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
	TargetName string
}
