package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"scp-parser/pkg/config"
	"scp-parser/server/internal/handlers"
	"scp-parser/server/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()
	r := chi.NewRouter()

	scpRoute, err := handlers.NewSCPHandler(ctx, cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("Error when creating SCPRoute: %v", err))
		return
	}

	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Mount("/scp", scpRoute.Routes())
		})
	})

	slog.Info(fmt.Sprintf("Server has been started on %s:%s", cfg.API.Host, cfg.API.Port))
	http.ListenAndServe(cfg.API.Host+":"+cfg.API.Port, r)
}
