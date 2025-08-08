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

	r.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Get("/session/create", CreateSession(storage)) // TODO: limit number of active sessions

		apiRouter.Route("/session/{id}/connect", func(connectionRouter chi.Router) {
			connectionRouter.Use(middleware.HttpVariables)

			connectionRouter.Use(middleware.SessionAuth(storage))

			connectionRouter.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) })

			connectionRouter.Get("/ws", WsUpgrade(storage))

		})
	})

	r.Get("/*", FrontendHandler)
}
