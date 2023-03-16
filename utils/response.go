package utils

import (
	"errors"
)

type ResponseRepository interface {
	Ok(status string, content ...interface{}) (successResponse, error)
	Fail(status string, msg string) (errorResponse, error)
}

type successResponse struct {
	Status string      `json:"status"`
	Content   []interface{} `json:"content"`
}

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"content"`
}

type Response struct{}

var _ ResponseRepository = (*Response)(nil)

func NewResponsesUtils() ResponseRepository {
	return Response{}
}

func (Response) Fail(status string, msg string) (errorResponse, error) {
	var res errorResponse

	if status == "" {
		return res, errors.New("internal error: missing status")
	}
	
	if msg == "" {
		return res, errors.New("internal error: missing message")
	}

	res.Status = status
	res.Message = msg

	return res, nil
}

func (Response) Ok(status string, Content ...interface{}) (successResponse, error) {
	var res successResponse

	if Content[0] == nil && len(Content) <= 1 {
		return res, errors.New("internal error: Content is nil")
	}

	if status == "" {
		return res, errors.New("internal error: missing status")
	}

	res.Status = status
	res.Content = Content

	return res, nil
}