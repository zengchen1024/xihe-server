package controller

import "github.com/opensourceways/xihe-server/domain/repository"

const (
	errorBadRequestBody    = "bad_request_body"
	errorBadRequestParam   = "bad_request_param"
	errorSystemError       = "system_error"
	errorDuplicateCreating = "duplicate_creating"
)

// responseData is the response data to client
type responseData struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func newResponseError(err error) responseData {
	code := errorSystemError

	switch err.(type) {
	case repository.ErrorDuplicateCreating:
		code = errorDuplicateCreating
	}

	return responseData{
		Code: code,
		Msg:  err.Error(),
	}
}

func newResponseData(data interface{}) responseData {
	return responseData{
		Data: data,
	}
}

func newResponseCodeError(code string, err error) responseData {
	return responseData{
		Code: code,
		Msg:  err.Error(),
	}
}

func newResponseCodeMsg(code, msg string) responseData {
	return responseData{
		Code: code,
		Msg:  msg,
	}
}
