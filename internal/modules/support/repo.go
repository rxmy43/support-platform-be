package support

import (
	"github.com/jmoiron/sqlx"
	"github.com/rxmy43/support-platform/internal/repo"
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
