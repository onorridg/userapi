package json

import (
	"encoding/json"
	"io/fs"
	"os"
	"refactoring/structs/http/code"
	customErr "refactoring/structs/http/response/error"
	"time"

	"refactoring/structs/user"
	"refactoring/utils/convert"
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

func (uS *UserStore) create(user *user.Request) uint64 {
	uS.Increment++
	i := uS.Increment
	uS.List[i] = &User{
		DisplayName: user.DisplayName,
		Email:       user.Email,
		CreatedAt:   time.Now(),
	}
	return i
}

func (uS *UserStore) update(user *user.Request) (uint64, int, error) {
	id, err := convert.StringToUInt64(user.ID)
	if err != nil {
		return 0, code.BadRequest, err
	}

	u, exist := uS.List[id]
	if !exist {
		return 0, code.BadRequest, customErr.UserNotFound
	}

	if user.Email != "" {
		u.Email = user.Email
	}
	if user.DisplayName != "" {
		u.DisplayName = user.DisplayName
	}
	uS.List[id] = u
	return id, 0, nil
}
func (uS *UserStore) Save(user *user.Request) (uint64, int, error) {
	var id uint64
	var err error
	var errCode int

	if _, errCode, err = uS.Search(); err != nil {
		return 0, errCode, err
	}

	if user.New {
		id = uS.create(user)
	} else if id, errCode, err = uS.update(user); err != nil {
		return 0, errCode, err
	}
	if err := uS.saveJSONUserStore(); err != nil {
		return 0, errCode, err
	}
	return id, 0, nil
}

func (uS *UserStore) Get(id string) (*User, int, error) {
	userID, err := convert.StringToUInt64(id)
	if err != nil {
		return nil, code.BadRequest, err
	}

	if err := uS.getJSONUserStore(); err != nil {
		return nil, code.InternalServerError, err
	}

	user, exist := uS.List[userID]
	if !exist {
		return nil, code.BadRequest, customErr.UserNotFound
	}
	return user, 0, nil
}

func (uS *UserStore) Search() (*UserStore, int, error) {
	if err := uS.getJSONUserStore(); err != nil {
		return nil, code.InternalServerError, err
	}
	return uS, 0, nil
}

func (uS *UserStore) Delete(id string) (int, error) {
	if _, errCode, err := uS.Search(); err != nil {
		return errCode, err
	}

	userID, err := convert.StringToUInt64(id)
	if err != nil {
		return code.BadRequest, err
	}

	if _, exist := uS.List[userID]; !exist {
		return code.BadRequest, customErr.UserNotFound
	}

	delete(uS.List, userID)
	if err = uS.saveJSONUserStore(); err != nil {
		return code.InternalServerError, err
	}
	return 0, nil
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
