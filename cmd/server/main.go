package main

import (
	"fmt"
	"net/http"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"github.com/tomassantos99/shared-browser-ide/internal/handlers"
	"github.com/tomassantos99/shared-browser-ide/internal/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetReportCaller(true);
	var r *chi.Mux = chi.NewRouter();

	sessionStorage := storage.NewSessionStorage()
	
	go storage.MemoryCleanup(sessionStorage, 60, 300)

	handlers.Handlers(r, sessionStorage)
	
	fmt.Println("Starting GO IDE service...")
	
	handler := cors.Default().Handler(r) //TODO: actually handle this
	err := http.ListenAndServe(":8080", handler)
	
	if err != nil {
		logrus.Error(err)
	}
}
