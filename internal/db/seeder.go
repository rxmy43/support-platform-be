package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rxmy43/support-platform/internal/modules/user"
)

func SeedUsers(ctx context.Context, db *sqlx.DB) error {
	users := []user.User{
		{Name: "NovaArtemis", Phone: "+62811100001", Role: "creator"},
		{Name: "LumenKai", Phone: "+62811100002", Role: "creator"},
		{Name: "OrionVex", Phone: "+62811100003", Role: "creator"},
		{Name: "AstraNova", Phone: "+62811100004", Role: "creator"},
		{Name: "VegaSol", Phone: "+62811100005", Role: "creator"},
		{Name: "EonLyra", Phone: "+62811100006", Role: "creator"},
		{Name: "CyraZen", Phone: "+62811100007", Role: "creator"},
		{Name: "KaiRen", Phone: "+62811100008", Role: "fan"},
		{Name: "MiraSol", Phone: "+62811100009", Role: "fan"},
		{Name: "LeoNix", Phone: "+62811100010", Role: "fan"},
	}

	userQuery := `
		INSERT INTO users (name, phone, role)
		VALUES (:name, :phone, :role)
		ON CONFLICT (phone) DO NOTHING
		RETURNING id;
	`

	balanceQuery := `
		INSERT INTO balances (user_id, amount)
		VALUES ($1, 0.00)
		ON CONFLICT (user_id) DO NOTHING;
	`

	for _, u := range users {
		var userID int64
		err := db.QueryRowxContext(ctx, userQuery, u.Name, u.Phone, u.Role).Scan(&userID)
		if err != nil {
			row := db.QueryRowxContext(ctx, "SELECT id FROM users WHERE phone=$1", u.Phone)
			if err := row.Scan(&userID); err != nil {
				return fmt.Errorf("failed to get user id for %s: %w", u.Name, err)
			}
		}

		if _, err := db.ExecContext(ctx, balanceQuery, userID); err != nil {
			return fmt.Errorf("failed to seed balance for user %s: %w", u.Name, err)
		}
	}

	fmt.Println("Seeded users and balances successfully")
	return nil
}
