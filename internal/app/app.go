package app

import (
	"log"
	"net/http"
	"time"

	"github.com/rxmy43/support-platform/internal/config"
	"github.com/rxmy43/support-platform/internal/db"
	"github.com/rxmy43/support-platform/internal/http/router"
	"github.com/rxmy43/support-platform/internal/socket"
)

type AppContext struct {
	Config *config.Config
	Router http.Handler
}

func InitApp() *AppContext {
	cfg := config.Load()

	DB, err := db.Connect(cfg.DB.DSN())
	if err != nil {
		log.Fatal("Database connection failed ", err)
	}
	defer db.Close()

	// if cfg.Env == "development" {
	// 	if err := db.SeedUsers(context.Background(), DB); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	hub := socket.NewHub()

	router := router.NewRouter(DB, hub)

	log.Println("Application bootstrap completed!")

	return &AppContext{
		Config: cfg,
		Router: router,
	}
}

func StartServer(appCtx *AppContext) {
	srv := &http.Server{
		Addr:         ":" + appCtx.Config.Port,
		Handler:      appCtx.Router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Server on http://localhost:" + appCtx.Config.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("server failed: ", err)
	}
}
