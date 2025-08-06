package middleware

import (
	"net/http"
	"github.com/go-chi/chi"
	"strings"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func HttpVariables(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		sessionId := chi.URLParam(r, "id")
		if sessionId == "" {
			errorMessage := "Missing session ID"
			logrus.Error(errorMessage)
			http.Error(w, errorMessage, http.StatusBadRequest)
			return
		}

		_, err := uuid.Parse(sessionId)
		if err != nil {
			errorMessage := "Invalid session ID"
			logrus.Error(errorMessage)
			http.Error(w, errorMessage, http.StatusBadRequest)
			return
		}

		var clientName = strings.TrimSpace(r.URL.Query().Get("name"))
		if clientName == "" {
			errorMessage := "Missing client name"
			logrus.Error(errorMessage)
			http.Error(w, errorMessage, http.StatusBadRequest)
			return
		}

		var sessionPassword = strings.TrimSpace(r.URL.Query().Get("password"))
		if sessionPassword == "" {
			errorMessage := "Missing session password"
			logrus.Error(errorMessage)
			http.Error(w, errorMessage, http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})

}
