package database

import (
	db "refactoring/database/json"
	"refactoring/structs/user"
)

type DBUserWorker interface {
	Get(id string) (*db.User, int, error)
	Search() (*db.UserStore, int, error)
	Save(user *user.Request) (uint64, int, error)
	Delete(id string) (int, error)
}

type DBWorker interface {
	DBUserWorker
}
