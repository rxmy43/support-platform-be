package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rxmy43/support-platform/internal/modules/user"
)

func SeedUsers(ctx context.Context, db *sqlx.DB) error {
	users := []user.User{
		{Name: "RinaArt", Phone: "+62811111111", Role: "creator"},
		{Name: "BambangSketch", Phone: "+62822222222", Role: "creator"},
		{Name: "UjangFan", Phone: "+62833333333", Role: "fan"},
	}

	query := `
		INSERT INTO users (name, phone, role)
		VALUES (:name, :phone, :role)
		ON CONFLICT (phone) DO NOTHING;
	`

	for _, u := range users {
		_, err := db.NamedExecContext(ctx, query, u)
		if err != nil {
			return fmt.Errorf("failed seeding user %s: %w", u.Name, err)
		}
	}

	fmt.Println("Seeded users successfully")
	return nil
}
