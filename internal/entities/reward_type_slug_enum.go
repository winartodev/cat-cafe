package entities

type RewardTypeSlug string

const (
	RewardTypeGoPayCoin RewardTypeSlug = "GOPAY_COIN"
	RewardTypeCoin      RewardTypeSlug = "COIN"
	RewardTypeGem       RewardTypeSlug = "GEM"
)

// IsValid checks if the reward type is valid
func (e RewardTypeSlug) IsValid() bool {
	switch e {
	case RewardTypeGoPayCoin, RewardTypeCoin, RewardTypeGem:
		return true
	default:
		return false
	}
}

// String returns string representation
func (e RewardTypeSlug) String() string {
	return string(e)
}

func (e RewardTypeSlug) RequiresBalanceUpdate() bool {
	return e == RewardTypeCoin || e == RewardTypeGem
}

func (e RewardTypeSlug) IsSentExternally() bool {
	return e == RewardTypeGoPayCoin
}
