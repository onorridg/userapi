package error

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

var (
	UserNotFound = errors.New("user_not_found")
)

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

func ErrRequest(err error, code int) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: code,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func RenderErr(w http.ResponseWriter, r *http.Request, err error, code int) error {
	if errRender := render.Render(w, r, ErrRequest(err, code)); errRender != nil {
		log.Println(errRender)
	}
	return nil
}
