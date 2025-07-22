package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Payload    interface{} `json:"payload,omitempty"`
	Meta       interface{} `json:"meta,omitempty"`
	Error      string      `json:"error,omitempty"`
}

func NewResponse() *Response {
	return &Response{
		StatusCode: http.StatusOK,
		Message:    "Success",
	}
}

func (r *Response) SetStatus(code int) *Response {
	r.StatusCode = code
	return r
}

func (r *Response) SetMessage(msg string) *Response {
	r.Message = msg
	return r
}

func (r *Response) SetPayload(data interface{}) *Response {
	r.Payload = data
	return r
}

func (r *Response) SetError(err string) *Response {
	r.Error = err
	return r
}

func (r *Response) SetMeta(meta interface{}) *Response {
	r.Meta = meta
	return r
}

func (r *Response) Send(c *gin.Context) {
	c.JSON(r.StatusCode, r)
}
