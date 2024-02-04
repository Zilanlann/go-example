package gee

import (
	"log"
	"net/http"
)

// router is a struct that holds a map of handlers, which associate a string key
// with a HandlerFunc (which is presumably a type defined elsewhere in your code).
type router struct {
	handlers map[string]HandlerFunc
}

// newRouter creates a new router instance with an initialized handlers map.
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// addRoute adds a new route to the router. It logs the method and pattern,
// and associates the handler function with a key in the handlers map.
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %s - %s", method, pattern) // Log the method and pattern
	key := method + "-" + pattern                // Create a key by concatenating method and pattern
	r.handlers[key] = handler                    // Map the key to the handler function
}

// handle is a method on router that takes a Context (presumably a struct that holds request info)
// and executes the appropriate handler based on the request method and path.
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path          // Create the key from the request method and path
	if handler, ok := r.handlers[key]; ok { // Check if a handler for the key exists
		handler(c) // If it exists, execute the handler with the context
	} else {
		// If no handler exists, respond with a 404 Not Found error
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
