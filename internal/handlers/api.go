package handlers

import (
	"github.com/go-chi/chi"
	"github.com/tomassantos99/shared-browser-ide/internal/middleware"
	"github.com/tomassantos99/shared-browser-ide/internal/storage"
	chimiddle "github.com/go-chi/chi/middleware"
)

func Handlers(r *chi.Mux, storage *storage.SessionStorage) {
	r.Use(chimiddle.StripSlashes)

	r.Route("/session/create", func(router chi.Router) {
		router.Get("/", CreateSession(storage))
	})

	r.Route("/session/{id}", func(router chi.Router) {

		router.Use(middleware.SessionAuth)

		router.Get("/ws", WsUpgrade(storage))
	})
}
