package database

import (
	db "refactoring/database/json"
	"refactoring/structs/user"
)

type DBUserWorker interface {
	Get(id string) (*db.User, error)
	Search() (*db.UserStore, error)
	Save(user *user.Request) (uint64, error)
	Delete(id string) error
}

type DBWorker interface {
	DBUserWorker
}
