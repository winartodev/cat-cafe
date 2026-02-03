package entities

import "github.com/winartodev/cat-cafe/pkg/apperror"

type UpgradeCostType string

const (
	UpgradeCostTypeCoin UpgradeCostType = "coin"
	UpgradeCostTypeGem  UpgradeCostType = "gem"
)

func (u UpgradeCostType) String() string {
	return string(u)
}

func (u UpgradeCostType) IsValid() bool {
	switch u {
	case UpgradeCostTypeCoin,
		UpgradeCostTypeGem:
		return true
	}
	return false
}

func ParseUpgradeCostType(s string) (UpgradeCostType, error) {
	costType := UpgradeCostType(s)
	if !costType.IsValid() {
		return "", apperror.ErrorInvalidRequest("upgrade cost type:", s)
	}
	return costType, nil
}

func AllUpgradeCostType() []UpgradeCostType {
	return []UpgradeCostType{
		UpgradeCostTypeCoin,
		UpgradeCostTypeGem,
	}
}
