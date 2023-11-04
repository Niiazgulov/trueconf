package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

// Объединненная структура для 2 обработчиков - создания и редактирования пользователя.
type UserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email,omitempty"`
}

func (c *UserRequest) Bind(r *http.Request) error {
	return nil
}

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
