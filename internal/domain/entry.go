package domain

import "time"

type Entry struct {
	ID        int `json:"id"`
	AccountID int `json:"account_id"`
	// can be negative or positive
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateEntryParams struct {
	AccountID int `json:"account_id"`
	Amount    int `json:"amount"`
}

type ListEntriesParams struct {
	AccountID int `json:"account_id"`
	Limit     int `json:"limit"`
	Offset    int `json:"offset"`
}
