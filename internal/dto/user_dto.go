package dto

type UserBalanceResponse struct {
	Coin int64 `json:"coin"`
	Gem  int64 `json:"gem"`
}
