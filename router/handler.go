package router

type HandleFunc func(ctx *Context)

type StaticHandler struct {
}

type CacheHandler struct {
}
