package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       // Represents the roots for different HTTP methods (GET, POST, etc.)
	handlers map[string]HandlerFunc // Stores handler functions associated with specific route patterns
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

// newRouter creates a new router instance with initialized data structures
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern splits a route pattern into individual parts, simplifies, and validates wildcard usage
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' { // Ensure only a single wildcard (*) is present
				break
			}
		}
	}
	return parts
}

// addRoute registers a route pattern, its handler function, and the supported HTTP method
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern // Creates a unique key for the handler map

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{} // Create a new root node for the method if it doesn't exist
	}
	r.roots[method].insert(pattern, parts, 0) // Insert the pattern into the routing tree
	r.handlers[key] = handler                 // Store the handler function
}

// getRoute attempts to find a matching route and extract parameters for a given HTTP method and path
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)

	root, ok := r.roots[method] // Retrieve the root node for the corresponding HTTP method
	if !ok {
		return nil, nil // No routes exist for the method
	}

	n := root.search(searchParts, 0) // Search for a matching pattern in the tree

	if n != nil {
		// If a match is found, extract parameters from the route pattern
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' { // Dynamic parameter (e.g., :id)
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 { // Catch-all parameter
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil // No matching route found
}

// handle is a method on router that takes a Context (a struct that holds request info)
// and executes the appropriate handler based on the request method and path.
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
