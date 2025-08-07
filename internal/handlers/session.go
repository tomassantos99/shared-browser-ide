package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/sirupsen/logrus"
	"github.com/tomassantos99/shared-browser-ide/api"
	"github.com/tomassantos99/shared-browser-ide/internal/storage"
	"github.com/tomassantos99/shared-browser-ide/internal/ws"
)

func CreateSession(sessionStorage *storage.SessionStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var session *ws.Session = ws.NewSession(sessionStorage.RemoveSession)
		sessionStorage.SaveSession(session)
		
		sessionResponse := api.Session{
			Id: session.Id.String(),
			Password: session.Password,
		}

		w.Header().Set("Content-Type", "application/json")

		data, err := json.Marshal(sessionResponse)
		if err != nil {
			logrus.Error(err)
			return
		}

		go session.Run()

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
