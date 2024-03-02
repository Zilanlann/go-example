package gee

import (
	"log"
	"net/http"
	"path"
	"strings"
	"text/template"
)

// HandlerFunc defines the request handler used by gee.
// It takes a pointer to Context as its only parameter.
type HandlerFunc func(*Context)

// Engine is the core component of the gee web framework. It manages routing,
// middleware, and the overall HTTP request/response lifecycle.
type Engine struct {
	*RouterGroup                     // Embeds the root RouterGroup for top-level routes
	router        *router            // The underlying router used for efficient route matching
	groups        []*RouterGroup     // Stores all defined RouterGroups for organization
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

// RouterGroup represents a group of routes that share a common prefix and middleware.
// It enables organization of routes within the web framework.
type RouterGroup struct {
	parent      *RouterGroup  // Pointer to the parent group (for nested groups)
	engine      *Engine       // Reference to the main Engine instance
	prefix      string        // The common prefix for all routes within this group
	middlewares []HandlerFunc // Middleware functions to apply to routes in this group
}

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

// Run defines the method to start an HTTP server on a specified address.
// It starts listening and serving HTTP requests using the Engine as the handler.
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// SetFuncMap sets the template function map for the engine.
// The funcMap parameter is a mapping of names to functions that can be called from within templates.
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// LoadHTMLGlob loads HTML templates from the specified pattern.
// It parses all the files that match the given pattern and adds them to the engine's HTML template collection.
// The pattern can include wildcard characters to match multiple files.
// The loaded templates can be rendered using the Render method.
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// ServeHTTP handles the HTTP request by executing the registered middlewares and routing the request.
// It implements the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}
