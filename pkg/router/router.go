package router

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Handle is a function that can be registered to a route to handle HTTP
// requests.
type Handle func(http.ResponseWriter, *http.Request, Params)

// Param is a URL parameter
type Params httprouter.Params

// ByName returns the value of the first Param which key matches the given name.
func (ps Params) ByName(name string) string {
	return (httprouter.Params)(ps).ByName(name)
}

type Router struct {
	router *httprouter.Router
}

func New() *Router {
	router := httprouter.New()
	return &Router{
		router: router,
	}
}
func (r *Router) handle(h Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		h(w, req, Params(ps))
	}
}

func (r *Router) GET(path string, handle Handle) {
	h := r.handle(handle)
	r.router.Handle("GET", path, h)
}

func (r *Router) HEAD(path string, handle Handle) {
	h := r.handle(handle)
	r.router.Handle("HEAD", path, h)
}

func (r *Router) OPTIONS(path string, handle Handle) {
	h := r.handle(handle)
	r.router.Handle("OPTIONS", path, h)
}

func (r *Router) POST(path string, handle Handle) {
	h := r.handle(handle)
	r.router.Handle("POST", path, h)
}

func (r *Router) PUT(path string, handle Handle) {
	h := r.handle(handle)
	r.router.Handle("PUT", path, h)
}

func (r *Router) PATCH(path string, handle Handle) {
	h := r.handle(handle)
	r.router.Handle("PATCH", path, h)
}

func (r *Router) DELETE(path string, handle Handle) {
	h := r.handle(handle)
	r.router.Handle("DELETE", path, h)
}

func (r *Router) Handle(path string, handle Handle) {
	h := r.handle(handle)
	r.router.Handle("GET", path, h)
	r.router.Handle("HEAD", path, h)
	r.router.Handle("OPTIONS", path, h)
	r.router.Handle("POST", path, h)
	r.router.Handle("PUT", path, h)
	r.router.Handle("PATCH", path, h)
	r.router.Handle("DELETE", path, h)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) serveFiles(path, basedir string, wildcard bool) {
	fileServer := http.FileServer(http.Dir(basedir))

	r.GET(path, func(w http.ResponseWriter, req *http.Request, ps Params) {
		if wildcard {
			path = ps.ByName("path")
		}
		req.URL.Path = path
		fileServer.ServeHTTP(w, req)
	})
}
