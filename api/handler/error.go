package handler

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ErrCodeOK           = 20000
	ErrCodeBadRequest   = 40000
	ErrCodeUnauthorized = 40100
	ErrCodeNotFound     = 40400
)

var errText = map[int]string{
	ErrCodeOK:           "OK",
	ErrCodeBadRequest:   "Bad Request",
	ErrCodeNotFound:     "Not Found",
	ErrCodeUnauthorized: "Unauthorized",
}

// ErrText returns a text for the status code. It returns the empty
// string if the code is unknown.
func ErrText(code int) string {
	return errText[code]
}

// 自定义error, 实现error.Error接口
type MessageError struct {
	Code    int
	Status  StatusType
	Message string
	Reason  string
}

func (e *MessageError) Error() string {
	return fmt.Sprintf("[%v] %v", e.Code, e.Message)
}

func NewMsgError(code int, msg string) error {
	return &MessageError{
		Code:    code,
		Status:  Failure,
		Message: msg,
		Reason:  ErrText(code),
	}
}

type StatusError struct {
	ErrStatus metav1.Status
}

func (e *StatusError) Error() string {
	return e.ErrStatus.Message
}
