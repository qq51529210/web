package router

import (
	"bytes"
	"net/http"
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

func (ctx *Context) Next() {
	if len(ctx.handleFunc) > ctx.handleIdx {
		f := ctx.handleFunc[ctx.handleIdx]
		ctx.handleIdx++
		f(ctx)
	}
}
