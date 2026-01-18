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

// RequiresBalanceUpdate return true if the reward type requires balance update
func (e RewardTypeSlug) RequiresBalanceUpdate() bool {
	return e == RewardTypeCoin || e == RewardTypeGem
}

// IsSentExternally return true if the reward type is sent externally
func (e RewardTypeSlug) IsSentExternally() bool {
	return e == RewardTypeGoPayCoin
}

// ToUserBalance returns the user balance type
func (e RewardTypeSlug) ToUserBalance() UserBalanceType {
	switch e {
	case RewardTypeCoin:
		return BalanceTypeCoin
	case RewardTypeGem:
		return BalanceTypeGem
	default:
		return BalanceTypeCoin
	}
}
