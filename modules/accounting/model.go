package accounting

import "time"

type TransferRequest struct {
	BankAccount string `json:"bank_account"`
	Amount      int64  `json:"amount"`
}

type TransferItem struct {
	ID              uint64 `json:"id"`
	FromBankAccount string `json:"from_bank_account"`
	ToBankAccount   string `json:"to_bank_account"`
	Amount          int64  `json:"amount"`
	CreatedAt       string `json:"created_at"`
	Direction       string `json:"direction"`
}

type Transfer struct {
	ID         uint64    `json:"id"`
	FromUserID uint64    `json:"from_user_id"`
	ToUserID   uint64    `json:"to_user_id"`
	Amount     int64     `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}
