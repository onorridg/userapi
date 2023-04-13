package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	dbJSON "refactoring/database/json"
	"refactoring/interfaces/database"
	customErr "refactoring/structs/error"
	"refactoring/structs/user"
)

var db database.DBWorker

func main() {
	db = dbJSON.InitDB()

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
	user, err := db.Get(id)
	if err != nil {
		if errRender := render.Render(w, r, customErr.ErrInvalidRequest(err)); errRender != nil {
			log.Println(errRender)
		}
		return
	}
	render.JSON(w, r, user)
}

func searchUsers(w http.ResponseWriter, r *http.Request) {
	userStore, _ := db.Search() // error handler
	render.JSON(w, r, userStore)
}

func createUser(w http.ResponseWriter, r *http.Request) {
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

func updateUser(w http.ResponseWriter, r *http.Request) {
	user := user.Request{}
	if err := render.Bind(r, &user); err != nil {
		_ = render.Render(w, r, customErr.ErrInvalidRequest(err))
		return
	}
	id := chi.URLParam(r, "id")

	user.New = false
	user.ID = id
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

func deleteUser(w http.ResponseWriter, r *http.Request) {
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
