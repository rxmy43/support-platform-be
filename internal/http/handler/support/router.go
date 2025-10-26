package support

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rxmy43/support-platform/internal/http/middleware"
	"github.com/rxmy43/support-platform/internal/modules/support"
	"github.com/rxmy43/support-platform/internal/modules/user"
)

func SupportRoutes(r chi.Router, db *sqlx.DB) {
	supportRepo := support.NewSupportRepo(db)
	userRepo := user.NewUserRepo(db)
	supportService := support.NewSupportService(supportRepo, userRepo)
	handler := NewSupportHandler(supportService)

	r.Route("/supports", func(r chi.Router) {
		r.Use(middleware.UserContext)
		r.Post("/", handler.Donate)
		r.Post("/callback", handler.PaymentCallback)
	})
}
