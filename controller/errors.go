package controller

import "fmt"

type HttpError struct {
	StatusCode int
	Title      string
	Detail     string
}

var InvalidRequest = HttpError{StatusCode: 400, Title: "Invalid Request"}
var InternalServerError = HttpError{StatusCode: 500, Title: "Internal Server Error"}
var Forbidden = HttpError{StatusCode: 403, Title: "Forbidden"}
var PageNotFound = HttpError{StatusCode: 404, Title: "Page Not Found"}

func (e *HttpError) Error() string {
	return fmt.Sprintf("code:%d, title:%s, detail:%s", e.StatusCode, e.Title, e.Detail)
}
