package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
)

const (
	_METHOD_GET = iota
	_METHOD_HEAD
	_METHOD_POST
	_METHOD_PUT
	_METHOD_PATCH
	_METHOD_DELETE
	_METHOD_CONNECT
	_METHOD_OPTIONS
	_METHOD_TRACE
	_METHOD_INVALID
)

func methodString(m int) string {
	switch m {
	case _METHOD_GET:
		return http.MethodGet
	case _METHOD_HEAD:
		return http.MethodHead
	case _METHOD_POST:
		return http.MethodPost
	case _METHOD_PUT:
		return http.MethodPut
	case _METHOD_PATCH:
		return http.MethodPatch
	case _METHOD_DELETE:
		return http.MethodDelete
	case _METHOD_CONNECT:
		return http.MethodConnect
	case _METHOD_OPTIONS:
		return http.MethodOptions
	case _METHOD_TRACE:
		return http.MethodTrace
	default:
		return ""
	}
}

type Router interface {
	GET(routePath string, handle ...HandleFunc)
	HEAD(routePath string, handle ...HandleFunc)
	POST(routePath string, handle ...HandleFunc)
	PUT(routePath string, handle ...HandleFunc)
	PATCH(routePath string, handle ...HandleFunc)
	DELETE(routePath string, handle ...HandleFunc)
	CONNECT(routePath string, handle ...HandleFunc)
	OPTIONS(routePath string, handle ...HandleFunc)
	TRACE(routePath string, handle ...HandleFunc)
}

type router struct {
	root [_METHOD_INVALID]route
}

func (r *router) Add(method int, routePath string, handle ...HandleFunc) {
	if len(handle) < 1 {
		panic(fmt.Errorf("[%s] %s fail, empty handle function", methodString(method), routePath))
	}
	r.root[method].Add(routePath, handle...)
}

func (r *router) GET(routePath string, handle ...HandleFunc) {
	r.Add(_METHOD_GET, routePath, handle...)
}

func (r *router) HEAD(routePath string, handle ...HandleFunc) {
	r.Add(_METHOD_HEAD, routePath, handle...)
}

func (r *router) POST(routePath string, handle ...HandleFunc) {
	r.Add(_METHOD_POST, routePath, handle...)
}

func (r *router) PUT(routePath string, handle ...HandleFunc) {
	r.Add(_METHOD_PUT, routePath, handle...)
}

func (r *router) PATCH(routePath string, handle ...HandleFunc) {
	r.Add(_METHOD_PATCH, routePath, handle...)
}

func (r *router) DELETE(routePath string, handle ...HandleFunc) {
	r.Add(_METHOD_DELETE, routePath, handle...)
}

func (r *router) CONNECT(routePath string, handle ...HandleFunc) {
	r.Add(_METHOD_CONNECT, routePath, handle...)
}

func (r *router) OPTIONS(routePath string, handle ...HandleFunc) {
	r.Add(_METHOD_OPTIONS, routePath, handle...)
}

func (r *router) TRACE(routePath string, handle ...HandleFunc) {
	r.Add(_METHOD_TRACE, routePath, handle...)
}

func (r *router) Static(routePath, file string, cache bool) {
	fi, err := os.Stat(file)
	if err != nil {
		panic(err)
	}
	// Is a file
	if !fi.IsDir() {
		fi, err := os.Stat(file)
		if err != nil {
			panic(err)
		}
		if !cache {
			r.root[_METHOD_GET].Add(routePath, (&FileHandler{file: file}).Handle)
		} else {
			h, err := NewCacheHandlerFromFile(file)
			if err != nil {
				panic(err)
			}
			r.root[_METHOD_GET].Add(routePath, h.Handle)
		}
		//
		if fi.Name() == "index.html" {
			file, _ = filepath.Split(file)
			if !cache {
				h := &FileHandler{file: file}
				r.root[_METHOD_GET].Add(routePath, h.Handle)
			} else {
				h, err := NewCacheHandlerFromFile(file)
				if err != nil {
					panic(err)
				}
				r.root[_METHOD_GET].Add(routePath, h.Handle)
			}
		}
		return
	}
	// Is a dir
	fis, err := ioutil.ReadDir(file)
	if err != nil {
		panic(err)
	}
	// Add sub
	for i := 0; i < len(fis); i++ {
		r.Static(path.Join(routePath, fis[i].Name()), filepath.Join(file, fis[i].Name()), cache)
	}
}

type subRouter struct {
	path string
	Router
}

type RootRouter interface {
	http.Handler
	Router
	SubRouter(routePath string) Router
	Global(handle ...HandleFunc)
	NotFound(handle ...HandleFunc)
	// Add static file handlers.
	// If file is a directory, add all sub files, routePath is root route of these handlers.
	// If file name is "index.html", add extra route "/".
	// If cache is true, use CachaHandler, else use FileHandler.
	Static(routePath, file string, cache bool)
}

func NewRootRouter() RootRouter {
	r := new(rootRouter)
	r.ctx.New = func() interface{} {
		return new(Context)
	}
	return r
}

type rootRouter struct {
	router
	global   []HandleFunc
	notfound []HandleFunc
	ctx      sync.Pool
}

func (r *rootRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := r.ctx.Get().(*Context)
	ctx.Request = req
	ctx.ResponseWriter = res
	ctx.Data = nil
	ctx.globalFunc = r.global
	ctx.globalIdx = 0
	ctx.handleIdx = 0
	//
	var route *route
	if req.Method[0] == 'G' {
		route = r.root[_METHOD_GET].Match(ctx)
	} else if req.Method[0] == 'H' {
		route = r.root[_METHOD_HEAD].Match(ctx)
	} else if req.Method[0] == 'D' {
		route = r.root[_METHOD_DELETE].Match(ctx)
	} else if req.Method[0] == 'C' {
		route = r.root[_METHOD_CONNECT].Match(ctx)
	} else if req.Method[0] == 'O' {
		route = r.root[_METHOD_OPTIONS].Match(ctx)
	} else if req.Method[0] == 'T' {
		route = r.root[_METHOD_TRACE].Match(ctx)
	} else if req.Method[1] == 'O' {
		route = r.root[_METHOD_POST].Match(ctx)
	} else if req.Method[1] == 'U' {
		route = r.root[_METHOD_PUT].Match(ctx)
	} else if req.Method[1] == 'A' {
		route = r.root[_METHOD_PATCH].Match(ctx)
	}
	if route == nil {
		ctx.handleFunc = r.notfound
	} else {
		ctx.handleFunc = route.handleFunc
	}
	ctx.Next()
	r.ctx.Put(ctx)
}

func (r *rootRouter) NotFound(handle ...HandleFunc) {
	if len(handle) != 0 {
		r.notfound = handle
	} else {
		r.notfound = []HandleFunc{
			func(ctx *Context) {
				ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
			},
		}
	}
}

func (r *rootRouter) Global(handle ...HandleFunc) {
	r.global = handle
}

func (r *rootRouter) SubRouter(routePath string) Router {
	return &subRouter{path: routePath, Router: r}
}
