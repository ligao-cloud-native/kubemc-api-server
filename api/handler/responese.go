package handler

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Code    int
	Status  StatusType
	Message string
	Reason  string
}

type StatusType string

var (
	Failure StatusType = "Failure"
	Success StatusType = "Success"
)

func ResError(w http.ResponseWriter, err error) {
	res := genResJson(err)

	if body, err := json.MarshalIndent(res, "", " "); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(res.Code)
		w.Write(body)
	}
}

func genResJson(err error) *APIResponse {
	resp := new(APIResponse)

	if err == nil {
		resp.Code = http.StatusOK
		resp.Status = Success
	} else {
		if e, ok := err.(*MessageError); ok {
			resp.Code = e.Code
			resp.Status = Failure
			resp.Reason = e.Reason
			resp.Message = e.Message
		} else if e, ok := err.(*StatusError); ok {
			resp.Code = int(e.ErrStatus.Code)
			resp.Status = Success
			resp.Message = e.ErrStatus.Message
			resp.Reason = string(e.ErrStatus.Reason)

		} else {
			resp.Code = errCodeBadRequest
			resp.Status = Failure
			resp.Message = err.Error()
			resp.Reason = ErrText(errCodeBadRequest)
		}

	}

	return resp
}
