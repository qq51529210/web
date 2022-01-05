package router

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"strings"
)

const (
	bearerTokenPrefix = "Bearer "
)

var (
	ContentTypeUTF8 = "charset=utf-8"
	ContentTypeJSON = mime.TypeByExtension(".json")
	ContentTypeHTML = mime.TypeByExtension(".html")
	ContentTypeJS   = mime.TypeByExtension(".js")
	ContentTypeCSS  = mime.TypeByExtension(".css")
)

type Context struct {
	*http.Request
	http.ResponseWriter
	Param      []string    // Dynamic route values, in the order of registration.
	TempData   interface{} // Keep user data duraing the call chain.
	handleFunc []HandleFunc
	handleIdx  int
}

func (ctx *Context) handle() {
	for ctx.handleIdx < len(ctx.handleFunc) {
		ctx.handleFunc[ctx.handleIdx](ctx)
		ctx.handleIdx++
	}
}

// Call next handler and run in current function.
// Example: chains: f1->f2->f3->f4 run in ServeHTTP function.
// Call in f1, f2->f3->f4 run in f1
func (ctx *Context) Handle() {
	ctx.handleIdx++
	ctx.handle()
}

// Abort handler chains.
// Example: chains: f1->f2->f3.
// Call Abort in f1, f2->f3 will not be called.
func (ctx *Context) Abort() {
	ctx.handleIdx = len(ctx.handleFunc)
}

// Return header["Authorization"] Bearer token.
func (c *Context) BearerToken() string {
	token := c.Request.Header.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, bearerTokenPrefix) {
		return ""
	}
	return token[len(bearerTokenPrefix):]
}

// Set Content-Type and statusCode, convert data to JSON and write to body.
func (c *Context) JSON(statusCode int, value interface{}) error {
	c.ResponseWriter.WriteHeader(statusCode)
	c.ResponseWriter.Header().Add("Content-Type", ContentTypeJSON)
	enc := json.NewEncoder(c.ResponseWriter)
	return enc.Encode(value)
}

// Set Content-Type and statusCode, data is JSON format.
func (c *Context) JSONBytes(statusCode int, data []byte) error {
	c.ResponseWriter.WriteHeader(statusCode)
	c.ResponseWriter.Header().Add("Content-Type", ContentTypeJSON)
	_, err := c.ResponseWriter.Write(data)
	return err
}

// Set Content-Type and statusCode, write text to body.
func (c *Context) HTML(statusCode int, text string) error {
	c.ResponseWriter.WriteHeader(statusCode)
	c.ResponseWriter.Header().Add("Content-Type", ContentTypeHTML)
	_, err := io.WriteString(c.ResponseWriter, text)
	return err
}
