package router

import (
	"bytes"
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
	Param      []string     // Dynamic route values, in the order of registration.
	Data       interface{}  // Keep user data in the handler chain.
	Buff       bytes.Buffer // A cache may be used.
	handleFunc []HandleFunc
	handleIdx  int
}

func (ctx *Context) handle() {
	for ctx.handleIdx < len(ctx.handleFunc) {
		ctx.handleFunc[ctx.handleIdx](ctx)
		ctx.handleIdx++
	}
}

func (ctx *Context) Next() {
	ctx.handleIdx++
	ctx.handle()
}

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
func (c *Context) WriteJSON(statusCode int, value interface{}) error {
	c.ResponseWriter.WriteHeader(statusCode)
	c.ResponseWriter.Header().Set("Content-Type", ContentTypeJSON)
	enc := json.NewEncoder(c.ResponseWriter)
	return enc.Encode(value)
}

// Set Content-Type and statusCode, write text to body.
func (c *Context) WriteHTML(statusCode int, text string) error {
	c.ResponseWriter.WriteHeader(statusCode)
	c.ResponseWriter.Header().Set("Content-Type", ContentTypeHTML)
	_, err := io.WriteString(c.ResponseWriter, text)
	return err
}
