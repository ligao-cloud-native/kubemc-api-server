package handler

import (
	"k8s.io/klog/v2"
	"net/http"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	klog.Infof("Request: ", r.RemoteAddr, r.URL.Scheme, r.Method, r.URL.RequestURI(), r.Proto)

	ResError(w, NewError(errCodeNotFound))

}
