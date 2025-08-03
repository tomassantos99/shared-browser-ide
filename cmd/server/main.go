package main

import (
	"fmt"
	"net/http"
	"github.com/go-chi/chi"
	"github.com/tomassantos99/shared-browser-ide/internal/handlers"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetReportCaller(true);
	var r *chi.Mux = chi.NewRouter();

	handlers.Handlers(r)
	
	fmt.Println("Starting GO IDE service...")

	err := http.ListenAndServe("localhost:8080", r)

	if err != nil {
		logrus.Error(err)
	}
}
