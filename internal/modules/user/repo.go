package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rxmy43/support-platform/internal/repo"
)

type UserRepo struct {
	*repo.BaseRepo[User]
}

func NewUserRepo(DB *sqlx.DB) *UserRepo {
	return &UserRepo{
		BaseRepo: &repo.BaseRepo[User]{
			DB:        DB,
			TableName: "users",
		},
	}
}

func (r *UserRepo) FindOneByPhone(ctx context.Context, phone string) (*User, error) {
	var u User
	err := r.DB.GetContext(ctx, &u, "SELECT id, name, phone FROM users WHERE phone = $1 LIMIT 1", phone)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
