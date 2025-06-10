package routes

import (
	"github.com/Ricky004/watchdata/pkg/api/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/v1/logs", handlers.GetLogs)
	return r
}
