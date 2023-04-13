package user

import "net/http"

type Request struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	ID          string `json:"-"`
	New         bool   `json:"-"`
}

func (c *Request) Bind(r *http.Request) error { return nil }
