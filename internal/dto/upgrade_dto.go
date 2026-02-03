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

type UpdateUpgradeDTO struct {
	Slug        string                   `json:"slug"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Cost        int64                    `json:"cost"`
	CostType    entities.UpgradeCostType `json:"cost_type"`
	IsActive    bool                     `json:"is_active"`
	Sequence    int64                    `json:"sequence"`
	Effect      UpgradeEffectDTO         `json:"effect"`
}

type BaseUpgradeResponseDTO struct {
	ID          int64                    `json:"id"`
	Slug        string                   `json:"slug"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Cost        int64                    `json:"cost"`
	CostType    entities.UpgradeCostType `json:"cost_type"`
	IsActive    bool                     `json:"is_active"`
	Sequence    int64                    `json:"sequence"`
	Effect      *UpgradeEffectDTO        `json:"effect,omitempty"`
}

type BaseStageUpgradeResponseDTO struct {
	Stage    string                   `json:"stage"`
	Upgrades []BaseUpgradeResponseDTO `json:"upgrades"`
}

type UpgradeEffectDTO struct {
	Type       entities.UpgradeEffectType   `json:"type"`
	Value      float64                      `json:"value"`
	Unit       entities.UpgradeEffectUnit   `json:"unit"`
	Target     entities.UpgradeEffectTarget `json:"target"`
	TargetID   *int64                       `json:"target_id,omitempty"`
	TargetName string                       `json:"target_name"`
}

type DetailUpgradeResponseDTO struct {
	UpgradeResponseDTO
	Effect *UpgradeEffectDTO `json:"effect,omitempty"`
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

type UpgradeResponseDTO struct {
	ID          int64  `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	Sequence    int64  `json:"sequence"`
}

func ToCreateUpgradeResponseDTO(e *entities.Upgrade) *BaseUpgradeResponseDTO {
	return &BaseUpgradeResponseDTO{
		ID:          e.ID,
		Slug:        e.Slug,
		Name:        e.Name,
		Description: e.Description,
		Cost:        e.Cost,
		CostType:    e.CostType,
		IsActive:    e.IsActive,
		Sequence:    e.Sequence,
		Effect:      ToUpgradeEffectDTO(&e.Effect),
	}
}

func ToGetUpgradesResponseDTO(e []entities.Upgrade) []UpgradeResponseDTO {
	var res []UpgradeResponseDTO
	for _, v := range e {
		res = append(res, UpgradeResponseDTO{
			ID:          v.ID,
			Slug:        v.Slug,
			Name:        v.Name,
			Description: v.Description,
			IsActive:    v.IsActive,
			Sequence:    v.Sequence,
		})
	}
	return res
}

func ToDetailUpgradeResponseDTO(e *entities.Upgrade) *BaseUpgradeResponseDTO {
	return &BaseUpgradeResponseDTO{
		ID:          e.ID,
		Slug:        e.Slug,
		Name:        e.Name,
		Description: e.Description,
		Cost:        e.Cost,
		CostType:    e.CostType,
		IsActive:    e.IsActive,
		Sequence:    e.Sequence,
		Effect:      ToUpgradeEffectDTO(&e.Effect),
	}
}

func ToStageUpgradesResponseDTO(stageSlug string, e []entities.StageUpgrade) BaseStageUpgradeResponseDTO {
	var upgrades []BaseUpgradeResponseDTO
	for _, v := range e {
		upgrade := v.Upgrade
		upgrades = append(upgrades, BaseUpgradeResponseDTO{
			ID:          upgrade.ID,
			Slug:        upgrade.Slug,
			Name:        upgrade.Name,
			Description: upgrade.Description,
			Cost:        upgrade.Cost,
			CostType:    upgrade.CostType,
			IsActive:    upgrade.IsActive,
			Sequence:    upgrade.Sequence,
			Effect:      nil,
		})
	}
	return BaseStageUpgradeResponseDTO{
		Stage:    stageSlug,
		Upgrades: upgrades,
	}
}

func (u *UpdateUpgradeDTO) ToEntity() entities.Upgrade {
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

func (u *UpdateUpgradeDTO) ValidateRequest() error {
	if !u.CostType.IsValid() {
		return apperror.ErrorInvalidRequest("cost type:", u.CostType.String())
	}

	if err := u.Effect.ValidateRequest(); err != nil {
		return err
	}

	return nil
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

func ToUpgradeEffectDTO(e *entities.UpgradeEffect) *UpgradeEffectDTO {
	return &UpgradeEffectDTO{
		Type:       e.Type,
		Value:      e.Value,
		Unit:       e.Unit,
		Target:     e.Target,
		TargetID:   &e.TargetID,
		TargetName: e.TargetName,
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
