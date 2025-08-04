package utils

import (
	"log"
	"net/http"

	"example.com/m/errs"
)

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

func ContainsString(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}
