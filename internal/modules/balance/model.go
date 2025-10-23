package balance

import "github.com/shopspring/decimal"

type Balance struct {
	ID     uint            `db:"id"`
	Amount decimal.Decimal `db:"amount"`
	UserID uint            `db:"user_id"`
}
