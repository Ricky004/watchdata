package routes

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/v1/logs", nil)
	return r
}
