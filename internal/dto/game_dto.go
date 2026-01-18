package dto

import "time"

type SyncBalanceRequest struct {
	CoinsEarned  int64     `json:"coins_earned"`
	LastSyncTime time.Time `json:"last_sync_time"`
}

type SyncBalanceResponse struct {
	CurrentCoinBalance int64 `json:"current_coin_balance"`
	CurrentGemBalance  int64 `json:"current_gem_balance"`
}
