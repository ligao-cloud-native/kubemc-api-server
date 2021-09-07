package api

import (
	"fmt"
	"github.com/ligao-cloud-native/kubemc-api-server/api/handler"
	"github.com/ligao-cloud-native/kubemc-api-server/pkg/router"
	"k8s.io/klog/v2"
	"net/http"
	"strings"
)

func (s *APIServer) auth(h router.Handle) router.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps router.Params) {
		klog.Infof("Request: ", r.RemoteAddr, r.URL.Scheme, r.Method, r.URL.RequestURI(), r.Proto)

		token := r.Header.Get("Authorization")
		if ok := s.validateToken(token); ok {
			h(w, r, ps)
		} else {
			klog.Errorf("invalid token: %v", token)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}

	}
}

// generate api access token
func (s *APIServer) genAccessToken(w http.ResponseWriter, r *http.Request, ps router.Params) {
	klog.Infof("Request: ", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	user := r.PostFormValue("username")
	pwd := r.PostFormValue("password")
	if ok := handler.AuthAccess(user, pwd); !ok {
		handler.ResError(w)
	}

}

// getCluster GET xwc
func (s *APIServer) getCluster(w http.ResponseWriter, r *http.Request, ps router.Params) {
	mux, ok := s.XwcMux[ManagerCluster]
	if !ok {

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
