package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tomassantos99/shared-browser-ide/internal/storage"
)

func SessionAuth(sessionStorage *storage.SessionStorage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionId := chi.URLParam(r, "id")
			sessionUUID, _ := uuid.Parse(sessionId)
			clientSessionPassword := strings.TrimSpace(r.URL.Query().Get("password"))

			var session, ok = sessionStorage.GetSession(sessionUUID)
			if !ok {
				errorMessage := fmt.Sprintf("Session with ID %s not found", sessionId)
				logrus.Error(errorMessage)
				http.Error(w, errorMessage, http.StatusNotFound)
				return
			}

			if session.Password != clientSessionPassword {
				errorMessage := fmt.Sprintf("Wrong password for session ID %s", sessionId)
				logrus.Error(errorMessage)
				http.Error(w, errorMessage, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
