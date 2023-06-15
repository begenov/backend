package domain

import "time"

type Transfer struct {
	ID            int `json:"id"`
	FromAccountID int `json:"from_account_id"`
	ToAccountID   int `json:"to_account_id"`
	// must be positive
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTransferParams struct {
	FromAccountID int `json:"from_account_id"`
	ToAccountID   int `json:"to_account_id"`
	Amount        int `json:"amount"`
}

type ListTransfersParams struct {
	FromAccountID int `json:"from_account_id"`
	ToAccountID   int `json:"to_account_id"`
	Limit         int `json:"limit"`
	Offset        int `json:"offset"`
}
