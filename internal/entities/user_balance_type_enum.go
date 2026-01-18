package entities

type UserBalanceType string

const (
	BalanceTypeCoin UserBalanceType = "COIN"
	BalanceTypeGem  UserBalanceType = "GEM"
)

func (e UserBalanceType) String() string {
	return string(e)
}
