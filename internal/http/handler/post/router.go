package post

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rxmy43/support-platform/internal/http/middleware"
	"github.com/rxmy43/support-platform/internal/modules/post"
	"github.com/rxmy43/support-platform/internal/modules/user"
)

func PostRoutes(r chi.Router, db *sqlx.DB) {
	postRepo := post.NewPostRepo(db)
	userRepo := user.NewUserRepo(db)

	postService := post.NewPostService(postRepo, userRepo)
	handler := NewPostHandler(postService)

	r.Route("/posts", func(r chi.Router) {
		r.Use(middleware.UserContext)
		r.Post("/", handler.Create)
		r.Get("/", handler.FindAll)
		r.Post("/ai-caption", handler.GenerateCaption)
	})
}
