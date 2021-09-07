package handler

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	errCodeOK         = 20000
	errCodeBadRequest = 40000
	errCodeNotFound   = 40400
)

var errText = map[int]string{
	errCodeOK:         "OK",
	errCodeBadRequest: "Bad Request",
	errCodeNotFound:   "Not Found",
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

func NewError(code int) error {
	return &MessageError{
		Code:    code,
		Status:  Failure,
		Message: ErrText(code),
		Reason:  ErrText(code),
	}
}

func NewMessageError(code int, msg string) error {
	var e MessageError
	return e
}

type StatusError struct {
	ErrStatus metav1.Status
}

func (e *StatusError) Error() string {
	return e.ErrStatus.Message
}
