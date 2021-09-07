package api

import (
	"github.com/ligao-cloud-native/kubemc-api-server/api/servemux"
	crdwatcher "github.com/ligao-cloud-native/kubemc-api-server/crd/watcher"
	"github.com/ligao-cloud-native/kubemc-api-server/pkg/ttl"
	"io/ioutil"
	"k8s.io/client-go/tools/cache"
	"os"
	"strings"
	"sync"

	"k8s.io/klog/v2"
	"net/http"

	"time"
)

const (
	// token ttl
	MaxTokenTTL = 3600 * 24 // one day
	// token len
	randTokenLen = 42
	// manager cluster name
	ManagerCluster = "leader"
)

type crdName string

const (
	WorkloadClusters crdName = "crd"
	ClientToken      crdName = "token"
)

type clusterMuxer struct {
	http.Handler
	Host        string
	BearerToken string
}

type APIServer struct {
	// http server
	//Server *server.Server
	// access token
	Token *ttl.TokenTTL
	// cluster mux
	XwcMux map[string]*clusterMuxer
	// crd watcher
	CRDWatcher map[crdName]crdwatcher.WatcherInterface
	//xwc cache
	XWCCache cache.Store

	sync.Mutex
}

func NewAPIServer() *APIServer {
	watcher := map[crdName]crdwatcher.WatcherInterface{
		WorkloadClusters: crdwatcher.NewXwcWatcher(),
	}

	return &APIServer{
		XwcMux:     make(map[string]*clusterMuxer),
		Token:      ttl.New(MaxTokenTTL, deleteToken),
		CRDWatcher: watcher,
	}

}

func (s *APIServer) Run() {
	// add reverse proxy serveMux for manage
	// r cluster
	go s.addManageClusterMux()

	// start crd resource watcher
	for _, watcher := range s.CRDWatcher {
		go watcher.Start()
	}

	// start http server
	mux := NewHTTPMux(s)
	server := &http.Server{
		Addr:         ":" + "8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		klog.Fatal("http server startup err", err)
	}

}

func (s *APIServer) addManageClusterMux() {
	m := clusterMuxer{
		Handler:     servemux.NewReverseProxyServeMux(),
		BearerToken: getClusterToken(),
	}

	s.Lock()
	if _, ok := s.XwcMux[ManagerCluster]; !ok {
		s.XwcMux[ManagerCluster] = &m
	}
	s.Unlock()
}

func (s *APIServer) validateToken(token string) bool {
	t := strings.Split(token, " ")
	if len(t) == 2 && strings.ToLower(t[0]) == "bearer" {
		if len(t[1]) == randTokenLen {
			ok := s.Token.Get(t[1])
			return ok != nil
		}

	}

	return false
}

func deleteToken(token string) {}

func getClusterToken() string {
	const tokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"

	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		t := os.Getenv("XWC_CONTROL_PLANE_API_TOKEN")
		if t == "" {
			klog.Fatal("GET cluster token from env XWC_CONTROL_PLANE_API_TOKEN error")
		}

		return t
	}

	return string(token)
}
