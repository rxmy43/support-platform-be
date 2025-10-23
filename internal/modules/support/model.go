package support

import "github.com/shopspring/decimal"

type Support struct {
	ID        uint            `db:"id"`
	FanID     uint            `db:"fan_id"`
	CreatorID uint            `db:"creator_id"`
	Amount    decimal.Decimal `db:"amount"`
	Status    string          `db:"status"`
}
