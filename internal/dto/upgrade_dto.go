package dto

import (
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
)

type CreateUpgradeDTO struct {
	Slug        string                   `json:"slug"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Cost        int64                    `json:"cost"`
	CostType    entities.UpgradeCostType `json:"cost_type"`
	IsActive    bool                     `json:"is_active"`
	Sequence    int64                    `json:"sequence"`
	Effect      UpgradeEffectDTO         `json:"effect"`
}

type UpgradeEffectDTO struct {
	Type       entities.UpgradeEffectType   `json:"type"`
	Value      float64                      `json:"value"`
	Unit       entities.UpgradeEffectUnit   `json:"unit"`
	Target     entities.UpgradeEffectTarget `json:"target"`
	TargetName string                       `json:"target_name"`
}

func (u *CreateUpgradeDTO) ToEntity() entities.Upgrade {
	return entities.Upgrade{
		Slug:        u.Slug,
		Name:        u.Name,
		Description: u.Description,
		Cost:        u.Cost,
		CostType:    u.CostType,
		IsActive:    u.IsActive,
		Sequence:    u.Sequence,
		Effect:      u.Effect.ToEntity(),
	}
}

func (u *CreateUpgradeDTO) ValidateRequest() error {
	if !u.CostType.IsValid() {
		return apperror.ErrorInvalidRequest("cost type:", u.CostType.String())
	}

	if err := u.Effect.ValidateRequest(); err != nil {
		return err
	}

	return nil
}

func (u *UpgradeEffectDTO) ToEntity() entities.UpgradeEffect {
	return entities.UpgradeEffect{
		Type:       u.Type,
		Value:      u.Value,
		Unit:       u.Unit,
		Target:     u.Target,
		TargetName: u.TargetName,
	}
}

func (u *UpgradeEffectDTO) ValidateRequest() error {
	if !u.Type.IsValid() {
		return apperror.ErrorInvalidRequest("effect type:", u.Type.String())
	}

	if !u.Unit.IsValid() {
		return apperror.ErrorInvalidRequest("effect unit:", u.Unit.String())
	}

	if !u.Target.IsValid() {
		return apperror.ErrorInvalidRequest("effect target:", u.Target.String())
	}

	return nil
}
