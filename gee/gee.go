package gee

import "net/http"

// HandlerFunc defines the request handler used by gee.
// It takes a pointer to Context as its only parameter.
type HandlerFunc func(*Context)

// Engine is the core of the gee web framework. It implements the interface of ServeHTTP
// and holds a pointer to a router that will be used to dispatch requests to the correct handler.
type Engine struct {
	router *router
}

// New is the constructor of gee.Engine. It initializes and returns a new Engine instance with a new router.
func New() *Engine {
	return &Engine{router: newRouter()}
}

// addRoute adds a route to the router with a specific method, pattern, and handler function.
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request handling with a specific pattern and handler.
// When a GET request matches the pattern, the handler will be invoked.
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request handling with a specific pattern and handler.
// When a POST request matches the pattern, the handler will be invoked.
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start an HTTP server on a specified address.
// It starts listening and serving HTTP requests using the Engine as the handler.
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP is the method that gets called by the HTTP server to handle each request.
// It creates a new Context for the current request, and delegates the request handling to the router.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
