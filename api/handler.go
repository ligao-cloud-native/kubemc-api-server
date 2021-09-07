package api

import (
	"fmt"
	"github.com/ligao-cloud-native/kubemc-api-server/api/handler"
	"github.com/ligao-cloud-native/kubemc-api-server/pkg/router"
	"k8s.io/klog/v2"
	"net/http"
	"strings"
)


// generate api access token
func (s *APIServer) genAccessToken(w http.ResponseWriter, r *http.Request, ps router.Params) {
	klog.Infof("Request: %v %v %v", r.RemoteAddr, r.Method, r.URL.RequestURI())

	user := r.PostFormValue("username")
	//pwd := r.PostFormValue("password")
	//if ok := handler.AuthAccess(user, pwd); !ok {
	//	handler.ResError(w, handler.ErrorMsg(handler.ErrCodeUnauthorized, "invalid username or password"))
	//	return
	//}

	token := handler.GetAccessToken(randTokenLen, MaxTokenTTL)

	// 每次login都需要更新token
	s.Token.PUT(token.Token, user)

	//TODO: create client token

	handler.ResOK(w, token)

}

func (s *APIServer) auth(h router.Handle) router.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps router.Params) {
		klog.Infof("Request: %v %v %v", r.RemoteAddr, r.Method, r.URL.RequestURI())

		token := r.Header.Get("Authorization")
		t := strings.Split(token, " ")
		if len(t) == 2 && strings.ToLower(t[0]) == "bearer" {
			if ok := s.Token.Get(t[1]); ok != nil {
				h(w, r, ps)
				return
			}
		}

		klog.Errorf("invalid token: %v", token)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}


// getCluster GET xwc
func (s *APIServer) getCluster(w http.ResponseWriter, r *http.Request, ps router.Params) {
	mux, ok := s.XwcMux[ManagerCluster]
	if !ok {
		handler.ResError(w, handler.ErrorMsg(handler.ErrCodeServiceUnavailable,""))
		return
	}

	if strings.HasPrefix(r.URL.Path, "/clusters") {
		r.URL.Path = strings.Replace(r.URL.Path,
			"/clusters", "apis/xwc.kubemc.io/v1/workloadclusters", 1)
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mux.BearerToken))
		mux.ServeHTTP(w, r)
	} else {
		handler.NotFound(w, r)
	}
}
