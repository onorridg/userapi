package main

import (
	"log"
	"net/http"
	"refactoring/api/v1"
	customHttp "refactoring/interfaces/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	dbJSON "refactoring/database/json"
	"refactoring/interfaces/database"
)

var db database.DBWorker

var v1Router customHttp.Handler

func main() {
	db = dbJSON.InitDB()
	v1Router = v1.Init()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	//r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})

	r.Route("/api", func(r chi.Router) {
		v1Router.Handler(r, db)
	})

	err := http.ListenAndServe(":3333", r)
	if err != nil {
		log.Println(err)
	}
}
