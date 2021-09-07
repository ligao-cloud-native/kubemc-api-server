/*
反向代理对客户端是透明的，也就是说客户端一般不知道代理的存在，认为自己是直接和服务器通信。
我们大部分访问的网站就是反向代理服务器，反向代理服务器会转发到真正的服务器，
一般在反向代理这一层实现负载均衡和高可用的功能。
*/

package servemux

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ReverseProxyHandler struct {
	HTTPReverseProxy      http.Handler
	WebSocketReverseProxy http.Handler
}

func NewReverseProxyServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	k8sHandler := ReverseProxyHandler{}
	k8sHandler.NewReverseProxy("https://kubernetes.default.svc:443")

	metricHandler := ReverseProxyHandler{}
	metricHandler.NewReverseProxy("http://prometheus.kube-system:9090")

	// kubernetes service api
	mux.Handle("/", &k8sHandler)
	// metric api
	mux.Handle("/api/v1/query", &metricHandler)

	return mux
}

func (h *ReverseProxyHandler) NewReverseProxy(target string) {
	reqUrl, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	h.HTTPReverseProxy = newHttpReverseProxy(reqUrl)
	h.WebSocketReverseProxy = newWebsocketReverseProxy(reqUrl)
}

func (h *ReverseProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.WebSocketReverseProxy != nil &&
		strings.ToLower(r.Header.Get("Connection")) == "upgrade" &&
		strings.ToLower(r.Header.Get("Upgrade")) == "websocket" {
		h.WebSocketReverseProxy.ServeHTTP(w, r)
	} else {
		h.HTTPReverseProxy.ServeHTTP(w, r)
	}
}

// newHttpReverseProxy to new a http ReverseProxy
func newHttpReverseProxy(url *url.URL) http.Handler {
	// new a HTTP ReverseProxy
	httpRP := httputil.NewSingleHostReverseProxy(url)

	// modify the request into a new request
	httpRP.Director = func(req *http.Request) {
		// Add 追加到存在的值后面；Set 重置，替换存在的值
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", req.Host)
		// Host header 不能通过 req.Header.Add("Host", req.Host)方式设置
		req.Host = url.Host

		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
		req.URL.Path = SingleJoiningSlash(url.Path, req.URL.Path)

		targetQuery := url.RawQuery
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + targetQuery
		}
	}

	if url.Scheme == "https" {
		httpRP.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         url.Host,
			},
		}
	}

	return httpRP

}

// newWebsocketReverseProxy to new a Websocket ReverseProxy
func newWebsocketReverseProxy(target *url.URL) http.Handler {
	wsScheme := "ws" + strings.TrimLeft(target.Scheme, "http")
	wsRP := NewSingleHostWsReverseProxy(&url.URL{
		Scheme: wsScheme,
		Host:   target.Host,
	})
	if wsScheme == "wss" {
		wsRP.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         target.Host,
		}
	}

	return wsRP
}
