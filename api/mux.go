package api

import (
	"github.com/ligao-cloud-native/kubemc-api-server/pkg/router"
)

func NewHTTPMux(s *APIServer) *router.Router {
	r := router.New()

	// 获取api访问token
	// curl -XPOST -H "Content-Type:application/x-www-form-urlencoded"
	// -d 'username=admin;password=xxx'
	// http://<xwc-apiserver-ip>:<port>/login
	r.POST("/login", s.genAccessToken)

	// 托管集群管理 (名称为master),
	r.POST("/clusters", s.auth(s.getCluster))
	//r.GET("/clusters",)
	//r.GET("/clusters/:cluster",)
	//r.DELETE("/clusters/:cluster",)
	//
	//
	//// 转发到k8s边缘集群的rest api （边缘集群名称为xwc名）
	//r.Handle("clusters/:cluster/api", )
	//r.Handle("clusters/:cluster/api/*path", )
	//r.Handle("clusters/:cluster/apis", )
	//r.Handle("clusters/:cluster/apis/*path", )

	return r

}
