package entities

import "fmt"

type RewardStatus string

const (
	StatusClaimed   RewardStatus = "claimed"
	StatusAvailable RewardStatus = "available"
	StatusLocked    RewardStatus = "locked"
)

func (e RewardStatus) IsValid() bool {
	switch e {
	case StatusClaimed, StatusAvailable, StatusLocked:
		return true
	default:
		return false
	}
}

func (e RewardStatus) String() string {
	return string(e)
}

func ParseRewardStatus(s string) (RewardStatus, error) {
	rs := RewardStatus(s)
	if !rs.IsValid() {
		return "", fmt.Errorf("%s is not a valid RewardStatus", s)
	}

	return rs, nil
}
