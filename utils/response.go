package utils

import (
	"log"
	"net/http"
	"strings"

	"example.com/m/dto"
	"example.com/m/errs"
	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Payload    any    `json:"payload,omitempty"`
	Meta       any    `json:"meta,omitempty"`
	Error      string `json:"error,omitempty"`
}

type LoginResponse struct {
	User  dto.PublicUser `json:"user"`
	Token string         `json:"token"`
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

func (r *Response) SetPayload(data any) *Response {
	r.Payload = data
	return r
}

func (r *Response) SetError(err string) *Response {
	r.Error = err
	return r
}

func (r *Response) SetMeta(meta any) *Response {
	r.Meta = meta
	return r
}

func (r *Response) Send(c *gin.Context) {
	c.JSON(r.StatusCode, r)
}

func GetSafeErrorMessage(err error, fallback string) string {
	if err == nil || strings.TrimSpace(err.Error()) == "" {
		return fallback
	}
	return err.Error()
}

func GetSafeStatusCode(statusCode int) int {
	if statusCode >= 200 && statusCode <= 503 {
		return statusCode
	}
	return http.StatusInternalServerError
}

func PrintErrorResponse(resp *Response, err error, message string) *Response {
	log.Printf("[ERROR] %v", err)
	if httpErr, ok := err.(*errs.HTTPError); ok {
		resp.SetStatus(httpErr.StatusCode).
			SetMessage(message).
			SetError(httpErr.Message)
		return resp
	}
	resp.SetStatus(http.StatusInternalServerError).
		SetMessage(message).
		SetError(GetSafeErrorMessage(err, "Unknown error occurred"))
	return resp
}
