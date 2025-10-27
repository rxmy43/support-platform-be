package support

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rxmy43/support-platform/internal/http/middleware"
	"github.com/rxmy43/support-platform/internal/modules/balance"
	"github.com/rxmy43/support-platform/internal/modules/support"
	"github.com/rxmy43/support-platform/internal/modules/user"
	"github.com/rxmy43/support-platform/internal/socket"
)

func SupportRoutes(r chi.Router, db *sqlx.DB, hub *socket.Hub) {
	supportRepo := support.NewSupportRepo(db)
	userRepo := user.NewUserRepo(db)
	balanceRepo := balance.NewBalanceRepo(db)
	supportService := support.NewSupportService(supportRepo, userRepo, balanceRepo, hub)
	handler := NewSupportHandler(supportService)

	r.Post("/payment/callback", handler.PaymentCallback)

	r.Route("/supports", func(r chi.Router) {
		r.Use(middleware.UserContext)
		r.Post("/", handler.Donate)
		r.Get("/best", handler.GetBestSupporters)
		r.Get("/fan-spending", handler.GetFanSpending)
		r.Get("/fan-spending/history", handler.GetFanSpendingHistory)
	})
}
