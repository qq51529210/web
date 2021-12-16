package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
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

var (
	notfoundHandlerFunc = []HandleFunc{
		func(ctx *Context) {
			ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		},
	}
)

type Root interface {
	http.Handler
	Router
	NotFound(handle ...HandleFunc)
	// Add static file handlers.
	// If file is a directory, add all sub files, routePath is root route of these handlers.
	// If file name is "index.html", add extra route "/".
	// If cache is true, use CachaHandler, else use FileHandler.
	Static(method, routePath, file string, cache bool) error
}

type router struct {
	root     [_METHOD_INVALID]*route
	notfound []HandleFunc
	ctx      sync.Pool
}

func (r *router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := r.ctx.Get().(*Context)
	ctx.Request = req
	ctx.ResponseWriter = res
	ctx.Param = ctx.Param[:0]
	ctx.Data = nil
	ctx.handleIdx = -1
	//
	var root *route
	if req.Method[0] == 'G' {
		root = r.root[_METHOD_GET]
	} else if req.Method[0] == 'H' {
		root = r.root[_METHOD_HEAD]
	} else if req.Method[0] == 'D' {
		root = r.root[_METHOD_DELETE]
	} else if req.Method[0] == 'C' {
		root = r.root[_METHOD_CONNECT]
	} else if req.Method[0] == 'O' {
		root = r.root[_METHOD_OPTIONS]
	} else if req.Method[0] == 'T' {
		root = r.root[_METHOD_TRACE]
	} else if req.Method[1] == 'O' {
		root = r.root[_METHOD_POST]
	} else if req.Method[1] == 'U' {
		root = r.root[_METHOD_PUT]
	} else if req.Method[1] == 'A' {
		root = r.root[_METHOD_PATCH]
	} else {
		ctx.handleFunc = r.notfound
		ctx.Next()
		r.ctx.Put(ctx)
		return
	}
	//
	route := root.Match(ctx)
	if route == nil {
		ctx.handleFunc = r.notfound
	} else {
		ctx.handleFunc = route.handleFunc
	}
	//
	ctx.Next()
	r.ctx.Put(ctx)
}

func (r *router) NotFound(handle ...HandleFunc) {
	if len(handle) != 0 {
		r.notfound = handle
	} else {
		r.notfound = notfoundHandlerFunc
	}
}

func (r *router) Sub(routePath string) Router {
	return &sub{path: routePath, router: r}
}

func (r *router) GET(routePath string, handle ...HandleFunc) error {
	return r.root[_METHOD_GET].Add(routePath, handle...)
}

func (r *router) HEAD(routePath string, handle ...HandleFunc) error {
	return r.root[_METHOD_HEAD].Add(routePath, handle...)
}

func (r *router) POST(routePath string, handle ...HandleFunc) error {
	return r.root[_METHOD_POST].Add(routePath, handle...)
}

func (r *router) PUT(routePath string, handle ...HandleFunc) error {
	return r.root[_METHOD_PUT].Add(routePath, handle...)
}

func (r *router) PATCH(routePath string, handle ...HandleFunc) error {
	return r.root[_METHOD_PATCH].Add(routePath, handle...)
}

func (r *router) DELETE(routePath string, handle ...HandleFunc) error {
	return r.root[_METHOD_DELETE].Add(routePath, handle...)
}

func (r *router) CONNECT(routePath string, handle ...HandleFunc) error {
	return r.root[_METHOD_CONNECT].Add(routePath, handle...)
}

func (r *router) OPTIONS(routePath string, handle ...HandleFunc) error {
	return r.root[_METHOD_OPTIONS].Add(routePath, handle...)
}

func (r *router) TRACE(routePath string, handle ...HandleFunc) error {
	return r.root[_METHOD_TRACE].Add(routePath, handle...)
}

func (r *router) Static(method, routePath, file string, cache bool) error {
	var root *route
	switch strings.ToUpper(method) {
	case http.MethodGet:
		root = r.root[_METHOD_GET]
	case http.MethodHead:
		root = r.root[_METHOD_HEAD]
	case http.MethodPost:
		root = r.root[_METHOD_POST]
	case http.MethodPut:
		root = r.root[_METHOD_PUT]
	case http.MethodPatch:
		root = r.root[_METHOD_PATCH]
	case http.MethodDelete:
		root = r.root[_METHOD_DELETE]
	case http.MethodConnect:
		root = r.root[_METHOD_CONNECT]
	case http.MethodOptions:
		root = r.root[_METHOD_OPTIONS]
	default:
		return fmt.Errorf("invalid method %s", method)
	}
	return r.static(root, routePath, file, cache)
}

func (r *router) static(root *route, routePath, file string, cache bool) error {
	fi, err := os.Stat(file)
	if err != nil {
		return err
	}
	// Is a file
	if !fi.IsDir() {
		fi, err := os.Stat(file)
		if err != nil {
			return err
		}
		if !cache {
			h := &FileHandler{file: file}
			err = root.Add(routePath, h.Handle)
		} else {
			h, err := NewCacheHandlerFromFile(file)
			if err != nil {
				return err
			}
			err = root.Add(routePath, h.Handle)
		}
		if err != nil {
			return err
		}
		//
		if fi.Name() == "index.html" {
			file, _ = filepath.Split(file)
			if !cache {
				h := &FileHandler{file: file}
				return root.Add(routePath, h.Handle)
			} else {
				h, err := NewCacheHandlerFromFile(file)
				if err != nil {
					return err
				}
				return root.Add(routePath, h.Handle)
			}
		}
		return nil
	}
	// Is a dir
	fis, err := ioutil.ReadDir(file)
	if err != nil {
		return err
	}
	// Add sub
	for i := 0; i < len(fis); i++ {
		err = r.static(root, path.Join(routePath, fis[i].Name()), filepath.Join(file, fis[i].Name()), cache)
		if err != nil {
			return err
		}
	}
	return nil
}
