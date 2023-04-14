package success

import (
	"github.com/go-chi/render"
	"net/http"
)

type Data struct {
}

type Response struct {
	StatusCode int         `json:"-"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, resp.StatusCode)
	render.JSON(w, r, resp)
	return nil
}
