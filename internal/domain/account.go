package domain

import "time"

type Account struct {
	ID        int       `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int       `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateAccountParams struct {
	Owner    string `json:"owner"`
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
}

type ListAccountsParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type UpdateAccountParams struct {
	ID      int `json:"id"`
	Balance int `json:"balance"`
}

type AddAccountBalanceParams struct {
	Amount int `json:"amount"`
	ID     int `json:"id"`
}
