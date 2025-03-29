package ghastly

import (
	"fmt"
	"net/http"
)

type GhastlyHandler struct {
	server *Ghastly
	mux    *http.ServeMux
}

func NewGhastlyHandler(server *Ghastly) *GhastlyHandler {
	return &GhastlyHandler{
		server: server,
		mux:    http.NewServeMux(),
	}
}

func (handler *GhastlyHandler) Request(method string, path string, middlewares []Middleware, function ApiHandler) {
	constructed_path := fmt.Sprintf("%s %s", method, path)
	if method == "*" {
		constructed_path = path
	}
	registeredFunction := RegisteredFunction{
		Server:      handler.server,
		Method:      method,
		Endpoint:    path,
		FullPath:    constructed_path,
		Middlewares: middlewares,
		Function:    function,
	}
	handler.mux.HandleFunc(constructed_path, registeredFunction.ServeHTTP)
}

func (handler *GhastlyHandler) Get(path string, middlewares []Middleware, function ApiHandler) {
	handler.Request("GET", path, middlewares, function)
}

func (handler *GhastlyHandler) Head(path string, middlewares []Middleware, function ApiHandler) {
	handler.Request("HEAD", path, middlewares, function)
}

func (handler *GhastlyHandler) Options(path string, middlewares []Middleware, function ApiHandler) {
	handler.Request("OPTIONS", path, middlewares, function)
}

func (handler *GhastlyHandler) PUT(path string, middlewares []Middleware, function ApiHandler) {
	handler.Request("PUT", path, middlewares, function)
}

func (handler *GhastlyHandler) POST(path string, middlewares []Middleware, function ApiHandler) {
	handler.Request("POST", path, middlewares, function)
}

func (handler *GhastlyHandler) PATCH(path string, middlewares []Middleware, function ApiHandler) {
	handler.Request("PATCH", path, middlewares, function)
}

func (handler *GhastlyHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	handler.mux.ServeHTTP(responseWriter, request)
}
