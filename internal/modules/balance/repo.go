package balance

import (
	"github.com/jmoiron/sqlx"
	"github.com/rxmy43/support-platform/internal/repo"
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
