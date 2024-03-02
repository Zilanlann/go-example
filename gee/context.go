package gee

import (
	"encoding/json" // Import the package for encoding and decoding JSON
	"fmt"           // Import the package for formatted I/O functions
	"net/http"      // Import the package for HTTP server and client
)

// H is a convenient alias for constructing JSON data
type H map[string]interface{}

// Context represents the context of a single HTTP request.
type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Params     map[string]string
	Path       string
	Method     string
	handlers   []HandlerFunc
	StatusCode int
	index      int
	engine     *Engine
}

// newContext creates a new Context object with the provided http.ResponseWriter and *http.Request.
// It initializes the Path, Method, Req, Writer, and index fields of the Context.
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:   req.URL.Path,
		Method: req.Method,
		Req:    req,
		Writer: w,
		index:  -1,
	}
}

// Next advances the context to the next middleware/handler in the chain.
// It calls the next middleware/handler in the chain by incrementing the index
// and invoking the corresponding handler function.
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Fail is a method of the Context struct that is used to handle a failed request.
// It sets the index of the handlers to the length of the handlers slice and
// responds with a JSON object containing the specified error message and status code.
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// Param returns the value of the specified parameter key from the Context's Params map.
func (c *Context) Param(key string) string {
	value := c.Params[key]
	return value
}

// PostForm retrieves data from a submitted form
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query retrieves the value of a URL query parameter
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status sets the status code for the response
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader sets a header field for the response
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String sends a plain text response
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")              // Set response type as plain text
	c.Status(code)                                         // Set response status code
	c.Writer.Write([]byte(fmt.Sprintf(format, values...))) // Send formatted response text
}

// TODO: need to handle json error

// JSON sends a response in JSON format
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json") // Set response type as JSON
	c.Status(code)                                  // Set response status code
	encoder := json.NewEncoder(c.Writer)            // Create a JSON encoder
	if err := encoder.Encode(obj); err != nil {     // Encode and send JSON data
		http.Error(c.Writer, err.Error(), 500) // Send an error response if encoding fails
	}
}

// Data sends a raw data response
func (c *Context) Data(code int, data []byte) {
	c.Status(code)       // Set response status code
	c.Writer.Write(data) // Send raw data
}

// HTML sends a response in HTML format
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}
