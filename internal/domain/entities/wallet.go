package entities

type Wallet struct {
	ID      int64  `json:"id" db:"id"`
	Balance int64  `json:"balance" db:"balance"`
	Owner   string `json:"owner" db:"owner"`
}
