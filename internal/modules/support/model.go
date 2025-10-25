package support

import (
	"time"

	"github.com/shopspring/decimal"
)

type Support struct {
	ID        uint            `db:"id"`
	FanID     uint            `db:"fan_id"`
	CreatorID uint            `db:"creator_id"`
	Amount    decimal.Decimal `db:"amount"`
	Status    string          `db:"status"`
	SupportID string          `db:"support_id"`
	SentAt    time.Time       `db:"sent_at"`
}
