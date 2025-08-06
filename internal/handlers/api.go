package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	"github.com/tomassantos99/shared-browser-ide/internal/middleware"
	"github.com/tomassantos99/shared-browser-ide/internal/storage"
)

func Handlers(r *chi.Mux, storage *storage.SessionStorage) {
	r.Use(chimiddle.StripSlashes)

	r.Route("/session/create", func(router chi.Router) {
		router.Get("/", CreateSession(storage))
	})

	r.Route("/session/{id}/connect", func(router chi.Router) {
		router.Use(middleware.HttpVariables)

		router.Use(middleware.SessionAuth(storage))

		router.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) })

		router.Get("/ws", WsUpgrade(storage))
	})
}
