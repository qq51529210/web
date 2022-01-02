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
	globalFunc []HandleFunc
	globalIdx  int
	handleFunc []HandleFunc
	handleIdx  int
}

func (ctx *Context) Next() {
	if len(ctx.globalFunc) > ctx.globalIdx {
		f := ctx.globalFunc[ctx.globalIdx]
		ctx.globalIdx++
		f(ctx)
	} else {
		if len(ctx.handleFunc) > ctx.handleIdx {
			f := ctx.handleFunc[ctx.handleIdx]
			ctx.handleIdx++
			f(ctx)
		}
	}
}
