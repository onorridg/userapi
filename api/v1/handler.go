package v1

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"refactoring/interfaces/database"
	"refactoring/structs/http/code"
	customErr "refactoring/structs/http/response/error"
	"refactoring/structs/http/response/success"
	"refactoring/structs/user"
)

type ApiV1 struct{}

var db database.DBWorker

func Init() *ApiV1 {
	return &ApiV1{}
}

func (api ApiV1) Handler(r chi.Router, dbInterface database.DBWorker) {
	db = dbInterface
	r.Route("/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/", api.SearchUsers)
			r.Post("/", api.CreateUser)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", api.GetUser)
				r.Patch("/", api.UpdateUser)
				r.Delete("/", api.DeleteUser)
			})
		})
	})
}

func (api ApiV1) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, errCode, err := db.Get(id)
	if err != nil {
		customErr.RenderErr(w, r, err, errCode)
		return
	}

	resp := success.Response{StatusCode: code.OK, Data: user}
	resp.Render(w, r)
}

func (api ApiV1) SearchUsers(w http.ResponseWriter, r *http.Request) {
	userStore, errCode, err := db.Search()
	if err != nil {
		customErr.RenderErr(w, r, err, errCode)
		return
	}

	resp := success.Response{StatusCode: code.OK, Data: userStore}
	resp.Render(w, r)
}

func (api ApiV1) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := user.Request{}
	if err := render.Bind(r, &user); err != nil {
		customErr.RenderErr(w, r, err, code.InternalServerError)
		return
	}

	user.New = true
	id, _, _ := db.Save(&user)

	msg := fmt.Sprint("user_id: ", id)
	resp := success.Response{StatusCode: code.Created, Message: msg}
	resp.Render(w, r)
}

func (api ApiV1) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := user.Request{}
	if err := render.Bind(r, &user); err != nil {
		customErr.RenderErr(w, r, err, code.InternalServerError)
		return
	}

	user.New, user.ID = false, id
	_, errCode, err := db.Save(&user)
	if err != nil {
		customErr.RenderErr(w, r, err, errCode)
		return
	}

	msg := fmt.Sprint("user_id: ", id)
	resp := success.Response{StatusCode: code.OK, Message: msg}
	resp.Render(w, r)
}

func (api ApiV1) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	errCode, err := db.Delete(id)
	if err != nil {
		customErr.RenderErr(w, r, err, errCode)
		return
	}

	msg := fmt.Sprint("user_id: ", id)
	resp := success.Response{StatusCode: code.OK, Message: msg}
	resp.Render(w, r)
}
