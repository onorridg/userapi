package json

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"refactoring/structs/user"
	"refactoring/utils/convert"
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

var (
	uV UserStore
)

func (uS *UserStore) create(user *user.Request) (uint64, error) {
	uS.Increment++
	i := uS.Increment
	uS.List[i] = &User{
		DisplayName: user.DisplayName,
		Email:       user.Email,
		CreatedAt:   time.Now(),
	}
	return i, nil
}

func (uS *UserStore) update(user *user.Request) (uint64, error) {
	id, err := convert.StringToUInt64(user.ID)
	if err != nil {
		return 0, err
	}

	u, exist := uS.List[id]
	if !exist {
		return 0, UserNotFound
	}

	if user.Email != "" {
		u.Email = user.Email
	}
	if user.DisplayName != "" {
		u.DisplayName = user.DisplayName
	}
	uS.List[id] = u
	return id, nil
}
func (uS *UserStore) Save(user *user.Request) (uint64, error) {
	var id uint64
	var err error

	uS.Search()
	if user.New {
		id, err = uS.create(user)
	} else if id, err = uS.update(user); err != nil {
		return 0, err
	}
	if err := uS.saveJSONUserStore(); err != nil {
		return 0, err
	}
	return id, nil
}

func (uS *UserStore) Get(id string) (*User, error) {
	userID, err := convert.StringToUInt64(id)
	if err != nil {
		return nil, err
	}
	if err := uS.getJSONUserStore(); err != nil {
		return nil, err
	}
	user, exist := uS.List[userID]
	if !exist {
		err = UserNotFound
	}
	return user, err
}

func (uS *UserStore) Search() (*UserStore, error) {
	if err := uS.getJSONUserStore(); err != nil {
		return nil, err
	}
	return uS, nil
}

func (uS *UserStore) Delete(id string) error {
	uS.Search()
	userID, err := convert.StringToUInt64(id)
	if err != nil {
		return err
	}
	if _, exist := uS.List[userID]; !exist {
		return UserNotFound
	}
	delete(uS.List, userID)
	if err = uS.saveJSONUserStore(); err != nil {
		return err
	}
	return nil
}

func (uS *UserStore) saveJSONUserStore() error {
	b, err := json.Marshal(uS)
	if err != nil {
		return err
	}
	if err = os.WriteFile(store, b, fs.ModePerm); err != nil {
		return err
	}
	return nil
}

func (uS *UserStore) getJSONUserStore() error {
	f, err := os.ReadFile(store)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(f, uS); err != nil {
		return err
	}
	return nil
}

func InitDB() *UserStore {
	return &UserStore{}
}
