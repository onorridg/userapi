package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"refactoring/interfaces/database"
)

type V1UserHandler interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	SearchUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type V1Handler interface {
	V1UserHandler
	Handler(r chi.Router, dbInterface database.DBWorker)
}

type Handler interface {
	V1Handler
}
