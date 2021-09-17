package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ligao-cloud-native/kubemc-api-server/api/handler"
	"github.com/ligao-cloud-native/kubemc-api-server/pkg/router"
	xwcv1 "github.com/ligao-cloud-native/kubemc/pkg/apis/xwc/v1"
	"io/ioutil"
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
	if !strings.HasPrefix(r.URL.Path, "/clusters") {
		handler.NotFound(w, r)
		return
	}

	r.URL.Path = "/apis/xwc.kubemc.io/v1/workloadclusters"
	clusterName := ps.ByName("cluster")
	if clusterName != "" {
		r.URL.Path = r.URL.Path + "/" + clusterName
	}
	klog.Info(r.URL.Path)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mux.BearerToken))
	mux.ServeHTTP(w, r)
}


func (s *APIServer) createCluster(w http.ResponseWriter, r *http.Request, ps router.Params) {
	mux, ok := s.XwcMux[ManagerCluster]
	if !ok {
		handler.ResError(w, handler.ErrorMsg(handler.ErrCodeServiceUnavailable,""))
		return
	}
	if !strings.HasPrefix(r.URL.Path, "/clusters") {
		handler.NotFound(w, r)
		return
	}

	wc := new(xwcv1.WorkloadCluster)
	buf, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	rc := ioutil.NopCloser(bytes.NewBuffer(buf))
	if err := json.Unmarshal(buf, wc); err != nil {
		klog.Errorf("Unmarshal body err: %s %v", string(buf), err)
		handler.ResError(w, handler.ErrorMsg(handler.ErrCodeBadRequest, err.Error()))
		return
	}
	r.Body = rc

	if err := handler.ValidateCluster(wc); err != nil {
		handler.ResError(w, handler.ErrorMsg(handler.ErrCodeUnprocessable, err.Error()))
		return
	}

	r.URL.Path = strings.Replace(r.URL.Path, "/clusters", "/apis/xwc.kubemc.io/v1/workloadclusters", 1)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mux.BearerToken))
	mux.ServeHTTP(w, r)
}

func (s *APIServer) scaleCluster(w http.ResponseWriter, r *http.Request, ps router.Params) {
	cluster := ps.ByName("cluster")
	if cluster == "" {
		handler.ResError(w, handler.ErrorMsg(handler.ErrCodeClusterNotFound, ""))
		return
	}
	var forceScale bool
	force, ok := r.URL.Query()["force"]
	if ok && force[0] == "true" {
		forceScale = true
	}


	scaleBody := new(handler.ClusterScale)
	if err := handler.ParseRequestBody(r, scaleBody); err != nil {
		klog.Error(err)
		handler.ResError(w, handler.ErrorMsg(handler.ErrCodeBadRequest, err.Error()))
		return
	}

	if (scaleBody.Scale != handler.ClusterScaleTypeUp && scaleBody.Scale != handler.ClusterScaleTypeDown) ||
		len(scaleBody.Workers) == 0 {
		handler.ResError(w, handler.ErrorMsg(handler.ErrCodeInvalidParam, ""))
		return
	}

	// TODO: check cluster is exist

	// TODO: get xwc cluster info
	wc := new(xwcv1.WorkloadCluster)

	// if allow scale
	var allowScale bool
	if wc.Status.Phase == xwcv1.WorkloadClusterSuccess ||
		(forceScale && wc.Status.Phase == xwcv1.WorkloadClusterFailed ){
		allowScale = true
	}
	if !allowScale {
		handler.ResError(w, handler.ErrorMsg(handler.ErrCodeActionNotSupport, "not support scale"))
		return
	}

	// TODO: 集群其他异常状态也不允许扩缩容

	// scale up/no worker nodes
	//workers := wc.Status.Cluster.Workers
	//workers := wc.Spec.Cluster.Workers
	if scaleBody.Scale == handler.ClusterScaleTypeUp {
		// TODO: add node
	} else {
		// TODO: remove node
	}

	//TODO: update wc

	handler.ResOK(w, wc)

}

func (s *APIServer) deleteCluster(w http.ResponseWriter, r *http.Request, ps router.Params) {

}