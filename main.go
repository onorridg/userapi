package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	db "refactoring/database"
)

var (
	UserNotFound = errors.New("user_not_found")
)

func main() {
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
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", searchUsers)
				r.Post("/", createUser)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", getUser)
					r.Patch("/", updateUser)
					r.Delete("/", deleteUser)
				})
			})
		})
	})

	err := http.ListenAndServe(":3333", r)
	if err != nil {
		log.Println(err)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := db.UserStore{}.Get(id)
	if err != nil {
		if errRender := render.Render(w, r, ErrInvalidRequest(err)); errRender != nil {
			log.Println(errRender)
		}
		return
	}
	render.JSON(w, r, user)
}

func searchUsers(w http.ResponseWriter, r *http.Request) {
	userStore := db.UserStore{}.Search()
	render.JSON(w, r, userStore)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	request := UserRequest{}
	if err := render.Bind(r, &request); err != nil {
		if err = render.Render(w, r, ErrInvalidRequest(err)); err != nil {
			log.Println(err)
		}
		return
	}

	userStore := db.UserStore{}.Search()
	user := &db.User{
		DisplayName: request.DisplayName,
		Email:       request.Email,
		New:         true,
	}
	id, _ := userStore.Save(user)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": id,
	})
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	request := UserRequest{}
	if err := render.Bind(r, &request); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	id := chi.URLParam(r, "id")
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	userStore := db.UserStore{}.Search()
	user := &db.User{
		DisplayName: request.DisplayName,
		Email:       request.Email,
		New:         false,
		ID:          userID,
	}
	_, err = userStore.Save(user)
	if err != nil {
		if err = render.Render(w, r, ErrInvalidRequest(err)); err != nil {
			log.Println(err)
		}
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"user_id": userID,
	})
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	userStore := db.UserStore{}.Search()
	err := userStore.Delete(id)
	if err != nil {
		if err = render.Render(w, r, ErrInvalidRequest(err)); err != nil {
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

type UserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func (c *UserRequest) Bind(r *http.Request) error { return nil }

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}
