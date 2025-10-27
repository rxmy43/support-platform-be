package balance

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/rxmy43/support-platform/internal/repo"
	"github.com/shopspring/decimal"
)

type BalanceRepo struct {
	*repo.BaseRepo[Balance]
}

func NewBalanceRepo(db *sqlx.DB) *BalanceRepo {
	return &BalanceRepo{
		BaseRepo: &repo.BaseRepo[Balance]{
			DB:        db,
			TableName: "balances",
		},
	}
}

func (r *BalanceRepo) GetBalanceAmountByUserID(ctx context.Context, userID uint) (int64, error) {
	var amountStr string

	query := `
		SELECT COALESCE(amount, 0)
		FROM balances
		WHERE user_id = $1
	`

	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&amountStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	// parse aman dari decimal Postgres
	dec, err := decimal.NewFromString(amountStr)
	if err != nil {
		return 0, err
	}

	return dec.IntPart(), nil
}
