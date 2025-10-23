package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rxmy43/support-platform/internal/modules/auth"
	"github.com/rxmy43/support-platform/internal/modules/user"
)

func AuthRoutes(r chi.Router, db *sqlx.DB) {
	userRepo := user.NewUserRepo(db)
	authService := auth.NewAuthService(userRepo)
	handler := NewAuthHandler(authService)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/generate-otp", handler.GenerateOTP)
		r.Post("/verify-otp", handler.VerifyOTP)
	})
}
