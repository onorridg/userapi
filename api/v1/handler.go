package v1

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"refactoring/interfaces/database"
	customErr "refactoring/structs/error"
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
	user, err := db.Get(id)
	if err != nil {
		if errRender := render.Render(w, r, customErr.ErrInvalidRequest(err)); errRender != nil {
			log.Println(errRender)
		}
		return
	}
	render.JSON(w, r, user)
}

func (api ApiV1) SearchUsers(w http.ResponseWriter, r *http.Request) {
	userStore, _ := db.Search() // error handler
	render.JSON(w, r, userStore)
}

func (api ApiV1) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := user.Request{}
	if err := render.Bind(r, &user); err != nil {
		if err = render.Render(w, r, customErr.ErrInvalidRequest(err)); err != nil {
			log.Println(err)
		}
		return
	}
	user.New = true
	id, _ := db.Save(&user)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": id,
	})
}

func (api ApiV1) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := user.Request{}
	if err := render.Bind(r, &user); err != nil {
		_ = render.Render(w, r, customErr.ErrInvalidRequest(err))
		return
	}
	id := chi.URLParam(r, "id")
	user.New, user.ID = false, id
	_, err := db.Save(&user)
	if err != nil {
		if err = render.Render(w, r, customErr.ErrInvalidRequest(err)); err != nil {
			log.Println(err)
		}
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"user_id": id,
	})
}

func (api ApiV1) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := db.Delete(id)
	if err != nil {
		if err = render.Render(w, r, customErr.ErrInvalidRequest(err)); err != nil {
			log.Println(err)
		}
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"message": "deleted",
		"user_id": id,
	})
}
