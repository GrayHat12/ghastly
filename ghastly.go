package ghastly

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"
)

// A Server defines parameters for running an HTTP server.
// The zero value for Server is a valid configuration.
type Server struct {
	// Addr optionally specifies the TCP address for the server to listen on,
	// in the form "host:port". If empty, ":http" (port 80) is used.
	// The service names are defined in RFC 6335 and assigned by IANA.
	// See net.Dial for details of the address format.
	Addr string

	// DisableGeneralOptionsHandler, if true, passes "OPTIONS *" requests to the Handler,
	// otherwise responds with 200 OK and Content-Length: 0.
	DisableGeneralOptionsHandler bool

	// TLSConfig optionally provides a TLS configuration for use
	// by ServeTLS and ListenAndServeTLS. Note that this value is
	// cloned by ServeTLS and ListenAndServeTLS, so it's not
	// possible to modify the configuration with methods like
	// tls.Config.SetSessionTicketKeys. To use
	// SetSessionTicketKeys, use Server.Serve with a TLS Listener
	// instead.
	TLSConfig *tls.Config

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body. A zero or negative value means
	// there will be no timeout.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	ReadTimeout time.Duration

	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body. If zero, the value of
	// ReadTimeout is used. If negative, or if zero and ReadTimeout
	// is zero or negative, there is no timeout.
	ReadHeaderTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	// A zero or negative value means there will be no timeout.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If zero, the value
	// of ReadTimeout is used. If negative, or if zero and ReadTimeout
	// is zero or negative, there is no timeout.
	IdleTimeout time.Duration

	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line. It does not limit the
	// size of the request body.
	// If zero, DefaultMaxHeaderBytes is used.
	MaxHeaderBytes int

	// TLSNextProto optionally specifies a function to take over
	// ownership of the provided TLS connection when an ALPN
	// protocol upgrade has occurred. The map key is the protocol
	// name negotiated. The Handler argument should be used to
	// handle HTTP requests and will initialize the Request's TLS
	// and RemoteAddr if not already set. The connection is
	// automatically closed when the function returns.
	// If TLSNextProto is not nil, HTTP/2 support is not enabled
	// automatically.
	TLSNextProto map[string]func(*http.Server, *tls.Conn, http.Handler)

	// ConnState specifies an optional callback function that is
	// called when a client connection changes state. See the
	// ConnState type and associated constants for details.
	ConnState func(net.Conn, http.ConnState)

	// ErrorLog specifies an optional logger for errors accepting
	// connections, unexpected behavior from handlers, and
	// underlying FileSystem errors.
	// If nil, logging is done via the log package's standard logger.
	ErrorLog *log.Logger

	// BaseContext optionally specifies a function that returns
	// the base context for incoming requests on this server.
	// The provided Listener is the specific Listener that's
	// about to start accepting requests.
	// If BaseContext is nil, the default is context.Background().
	// If non-nil, it must return a non-nil context.
	BaseContext func(net.Listener) context.Context

	// ConnContext optionally specifies a function that modifies
	// the context used for a new connection c. The provided ctx
	// is derived from the base context and has a ServerContextKey
	// value.
	ConnContext func(ctx context.Context, c net.Conn) context.Context
}

type Ghastly struct {
	Server  *http.Server
	handler *GhastlyHandler
}

func NewGhastly(server Server) *Ghastly {
	ghastly := &Ghastly{
		Server: &http.Server{
			Addr:                         server.Addr,
			DisableGeneralOptionsHandler: server.DisableGeneralOptionsHandler,
			TLSConfig:                    server.TLSConfig,
			ReadTimeout:                  server.ReadTimeout,
			ReadHeaderTimeout:            server.ReadHeaderTimeout,
			WriteTimeout:                 server.WriteTimeout,
			IdleTimeout:                  server.IdleTimeout,
			MaxHeaderBytes:               server.MaxHeaderBytes,
			TLSNextProto:                 server.TLSNextProto,
			ConnState:                    server.ConnState,
			ErrorLog:                     server.ErrorLog,
			BaseContext:                  server.BaseContext,
			ConnContext:                  server.ConnContext,
		},
	}
	ghastly.handler = NewGhastlyHandler(ghastly)
	ghastly.Server.Handler = ghastly.handler
	return ghastly
}

func (server *Ghastly) Request(method string, path string, middlewares []Middleware, function ApiHandler) {
	server.handler.Request(method, path, middlewares, function)
}

func (server *Ghastly) Get(path string, middlewares []Middleware, function ApiHandler) {
	server.handler.Request("GET", path, middlewares, function)
}

func (server *Ghastly) Head(path string, middlewares []Middleware, function ApiHandler) {
	server.handler.Request("HEAD", path, middlewares, function)
}

func (server *Ghastly) Options(path string, middlewares []Middleware, function ApiHandler) {
	server.handler.Request("OPTIONS", path, middlewares, function)
}

func (server *Ghastly) PUT(path string, middlewares []Middleware, function ApiHandler) {
	server.handler.Request("PUT", path, middlewares, function)
}

func (server *Ghastly) POST(path string, middlewares []Middleware, function ApiHandler) {
	server.handler.Request("POST", path, middlewares, function)
}

func (server *Ghastly) PATCH(path string, middlewares []Middleware, function ApiHandler) {
	server.handler.Request("PATCH", path, middlewares, function)
}

func (server *Ghastly) Close() {
	server.Server.Close()
}

func (server *Ghastly) ListenAndServe() error {
	return server.Server.ListenAndServe()
}

func (server *Ghastly) ListenAndServeTLS(certFile string, keyFile string) error {
	return server.Server.ListenAndServeTLS(certFile, keyFile)
}

func (server *Ghastly) Serve(l net.Listener) error {
	return server.Server.Serve(l)
}

func (server *Ghastly) ServeTLS(l net.Listener, certFile string, keyFile string) error {
	return server.Server.ServeTLS(l, certFile, keyFile)
}

func (server *Ghastly) Shutdown(ctx context.Context) error {
	return server.Server.Shutdown(ctx)
}
