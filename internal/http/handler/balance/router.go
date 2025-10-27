package balance

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rxmy43/support-platform/internal/http/middleware"
	"github.com/rxmy43/support-platform/internal/modules/balance"
	"github.com/rxmy43/support-platform/internal/modules/user"
)

func BalanceRoutes(r chi.Router, db *sqlx.DB) {
	balanceRepo := balance.NewBalanceRepo(db)
	userRepo := user.NewUserRepo(db)

	balanceService := balance.NewBalanceService(balanceRepo, userRepo)
	handler := NewBalanceHandler(balanceService)

	r.Route("/balances", func(r chi.Router) {
		r.Use(middleware.UserContext)
		r.Get("/creator", handler.GetCreatorBalance)
	})
}
