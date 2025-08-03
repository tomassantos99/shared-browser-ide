package handlers

import (
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	"github.com/tomassantos99/shared-browser-ide/internal/middleware"
)

func Handlers(r *chi.Mux) {
	r.Use(chimiddle.StripSlashes)

	r.Route("/session/create", func(router chi.Router) {
		router.Get("/", CreateSession)
	})

	r.Route("/session/{id}", func(router chi.Router) {

		router.Use(middleware.SessionAuth)

		// TODO: create page and route to it
		// router.Get("/", BrowserIDE)

		router.Get("/ws", WsUpgrade)
	})
}
