package support

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/rxmy43/support-platform/internal/repo"
	"github.com/shopspring/decimal"
)

type SupportRepo struct {
	*repo.BaseRepo[Support]
}

func NewSupportRepo(DB *sqlx.DB) *SupportRepo {
	return &SupportRepo{
		BaseRepo: &repo.BaseRepo[Support]{
			DB:        DB,
			TableName: "supports",
		},
	}
}

type supportTimestampCreatorFan struct {
	PaymentTimestamp int64  `db:"payment_timestamp"`
	FanID            uint   `db:"fan_id"`
	CreatorID        uint   `db:"creator_id"`
	FanName          string `db:"fan_name"`
	CreatorName      string `db:"creator_name"`
}

func (r *SupportRepo) GetPaymentTimestamp(ctx context.Context, reference, supportID string) (*supportTimestampCreatorFan, error) {
	var result supportTimestampCreatorFan

	query := `
		SELECT 
			s.payment_timestamp, 
			s.creator_id, 
			s.fan_id, 
			f.name AS fan_name,
			c.name AS creator_name
		FROM supports s
		JOIN users f ON s.fan_id = f.id
		JOIN users c ON s.creator_id = c.id
		WHERE s.reference_code = $1
		AND s.support_id = $2
	`

	err := r.DB.GetContext(ctx, &result, query, reference, supportID)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *SupportRepo) GetCreatorSupporters(ctx context.Context, cursor *uint, creatorID uint) ([]BestSupporters, *uint, error) {
	supporters := []BestSupporters{}

	queryBase := `
		SELECT 
			s.id,
			f.name AS fan_name,
			s.amount,
			s.sent_at
		FROM supports s
		JOIN users f ON f.id = s.fan_id
		WHERE s.creator_id = $1
		AND s.status = 'paid'
	`

	args := []any{creatorID}
	if cursor != nil {
		queryBase += " AND s.id < $2"
		args = append(args, *cursor)
	}

	query := queryBase + `
		ORDER BY s.id DESC
		LIMIT 10
	`

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s BestSupporters
		var amountNumeric string
		if err := rows.Scan(new(uint), &s.FanName, &amountNumeric, &s.SentAt); err != nil {
			return nil, nil, err
		}

		if strings.Contains(amountNumeric, ".") {
			parts := strings.Split(amountNumeric, ".")
			amountNumeric = parts[0]
		}
		amount, err := strconv.ParseInt(amountNumeric, 10, 64)
		if err != nil {
			return nil, nil, err
		}
		s.Amount = amount
		supporters = append(supporters, s)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	var nextCursor *uint
	if len(supporters) > 0 {
		lastID := supporters[len(supporters)-1].ID
		nextCursor = &lastID
	}

	return supporters, nextCursor, nil
}

func (r *SupportRepo) GetFanSpendingAmount(ctx context.Context, fanID uint) (int64, error) {
	var amount decimal.Decimal

	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM supports
		WHERE fan_id = $1
		AND status = 'paid'
	`

	err := r.DB.QueryRowContext(ctx, query, fanID).Scan(&amount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return amount.IntPart(), nil
}

func (r *SupportRepo) GetFanSupportHistory(ctx context.Context, cursor, fanID *uint) ([]FanSupportHistory, *uint, error) {
	histories := []FanSupportHistory{}
	var err error

	queryBase := `
		SELECT
			s.id,
			s.amount,
			s.sent_at,
			c.name AS creator_name
		FROM supports s
		JOIN users c ON c.id = s.creator_id
		WHERE s.fan_id = $1
		AND s.status = 'paid'
	`

	args := []any{*fanID}
	if cursor != nil {
		queryBase += " AND s.id < $2"
		args = append(args, *cursor)
	}

	query := `
		` + queryBase + `
		ORDER BY s.id DESC
		LIMIT 10
	`

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var h FanSupportHistory
		var amountNumeric string
		if err := rows.Scan(&h.ID, &amountNumeric, &h.SentAt, &h.CreatorName); err != nil {
			return nil, nil, err
		}

		// konversi dari numeric(15,2) → int64
		if strings.Contains(amountNumeric, ".") {
			parts := strings.Split(amountNumeric, ".")
			amountNumeric = parts[0]
		}
		amount, err := strconv.ParseInt(amountNumeric, 10, 64)
		if err != nil {
			return nil, nil, err
		}
		h.Amount = amount
		histories = append(histories, h)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	var nextCursor *uint
	if len(histories) > 0 {
		lastID := histories[len(histories)-1].ID
		nextCursor = &lastID
	}

	return histories, nextCursor, nil
}
