package ghastly

import "net/http"

type ApiHandler func(*map[string]string, http.ResponseWriter, *http.Request)
type Middleware func(*map[string]string, http.ResponseWriter, *http.Request, func())

type RegisteredFunction struct {
	Server      *Ghastly
	Method      string
	Endpoint    string
	FullPath    string
	Middlewares []Middleware
	Function    ApiHandler
}

func run(function *RegisteredFunction, responseWriter http.ResponseWriter, request *http.Request, context *map[string]string, index int) {
	if index < len(function.Middlewares) {
		function.Middlewares[index](context, responseWriter, request, func() {
			run(function, responseWriter, request, context, index+1)
		})
	} else {
		function.Function(context, responseWriter, request)
	}
}

func (function *RegisteredFunction) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	context := map[string]string{}
	run(function, responseWriter, request, &context, 0)
}
