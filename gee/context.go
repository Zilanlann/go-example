package gee

import (
	"encoding/json" // Import the package for encoding and decoding JSON
	"fmt"           // Import the package for formatted I/O functions
	"net/http"      // Import the package for HTTP server and client
)

// H is a convenient alias for constructing JSON data
type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Params     map[string]string
	Path       string
	Method     string
	handlers   []HandlerFunc
	StatusCode int
	index      int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:   req.URL.Path,
		Method: req.Method,
		Req:    req,
		Writer: w,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

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
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html") // Set response type as HTML
	c.Status(code)                           // Set response status code
	c.Writer.Write([]byte(html))             // Send HTML content
}
