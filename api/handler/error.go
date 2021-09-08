package handler

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ErrCodeOK           = 20000
	ErrCodeBadRequest   = 40000
	ErrCodeInvalidParam =  40010
	ErrCodeActionNotSupport = 40030
	ErrCodeUnauthorized = 40100
	ErrCodeNotFound     = 40400
	ErrCodeUnprocessable = 42200
	ErrCodeServiceUnavailable = 50300
	ErrCodeClusterNotFound = 40430


)

var errText = map[int]string{
	ErrCodeOK:           "OK",
	ErrCodeBadRequest:   "Bad Request",
	ErrCodeNotFound:     "Not Found",
	ErrCodeUnauthorized: "Unauthorized",
	ErrCodeServiceUnavailable: "ServiceUnavailable",
	ErrCodeUnprocessable: "Unprocessable",
	ErrCodeClusterNotFound: "ClusterNotExist",
	ErrCodeInvalidParam: "InvalidParam",
	ErrCodeActionNotSupport: "ActionNotSupport",
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
	Reason  string
	Message string
}

func (e *MessageError) Error() string {
	return fmt.Sprintf("[%v] %v", e.Code, e.Message)
}


func (e *MessageError) New(code int, msg string) error {
	e.Code = code
	e.Status = Failure
	e.Reason = ErrText(code)
	if msg != "" {
		e.Message = msg
	} else {
		e.Message = ErrText(code)
	}

	return e
}


func ErrorMsg(code int, msg string) error {
	var e MessageError
	return e.New(code, msg)
}

type StatusError struct {
	ErrStatus metav1.Status
}

func (e *StatusError) Error() string {
	return e.ErrStatus.Message
}
