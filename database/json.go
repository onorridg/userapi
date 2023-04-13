package database

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"strconv"
	"time"
)

const store = `users.json`

type (
	User struct {
		CreatedAt   time.Time `json:"created_at"`
		DisplayName string    `json:"display_name"`
		Email       string    `json:"email"`
		New         bool      `json:"-"`
		ID          uint64    `json:"-"`
	}
	UserList  map[uint64]*User
	UserStore struct {
		Increment uint64   `json:"increment"`
		List      UserList `json:"list"`
	}
)

var (
	UserNotFound = errors.New("user_not_found")
)

func (uS *UserStore) Save(user *User) (uint64, error) {
	if user.New {
		uS.Increment++
		user.ID = uS.Increment
		user.CreatedAt = time.Now()
		uS.List[user.ID] = user
	} else {
		u, exist := uS.List[user.ID]
		if !exist {
			return 0, UserNotFound
		}
		if user.Email != "" {
			u.Email = user.Email
		}
		if user.DisplayName != "" {
			u.DisplayName = user.DisplayName
		}
		uS.List[user.ID] = u
	}
	b, _ := json.Marshal(uS)
	_ = os.WriteFile(store, b, fs.ModePerm)
	return user.ID, nil
}

func (uS UserStore) Get(id string) (*User, error) {
	userID, err := parseID(id)
	if err != nil {
		return nil, err
	}
	data := getJSONUserData()
	user, exist := data.List[userID]
	if !exist {
		err = UserNotFound
	}
	return user, err
}

func (uS UserStore) Search() *UserStore {
	userStore := getJSONUserData()
	return userStore
}

func (uS *UserStore) Delete(id string) error {
	userID, err := parseID(id)
	if err != nil {
		return err
	}
	if _, exist := uS.List[userID]; !exist {
		return UserNotFound
	}
	delete(uS.List, userID)
	return nil
}

func parseID(id string) (uint64, error) {
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func getJSONUserData() *UserStore {
	f, _ := os.ReadFile(store)
	s := UserStore{}
	_ = json.Unmarshal(f, &s)
	return &s
}
