package router

import (
	"net/http"
)

var (
	notfoundHandlerFunc = []HandleFunc{
		func(ctx *Context) {
			ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		},
	}
)

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

// func NewRootRouter() RootRouter {
// 	r := new(rootRouter)
// 	r.ctx.New = func() interface{} {
// 		return new(Context)
// 	}
// 	for i := 0; i < _METHOD_INVALID; i++ {
// 		r.root[i] = new(route)
// 	}
// 	return r
// }

// type rootRouter struct {
// 	root     [_METHOD_INVALID]*route
// 	global   []HandleFunc
// 	notfound []HandleFunc
// 	ctx      sync.Pool
// }

// func (r *rootRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
// 	ctx := r.ctx.Get().(*Context)
// 	ctx.Request = req
// 	ctx.ResponseWriter = res
// 	ctx.Data = nil
// 	ctx.globalFunc = r.global
// 	ctx.globalIdx = 0
// 	ctx.handleIdx = 0
// 	//
// 	var route *route
// 	if req.Method[0] == 'G' {
// 		route = r.root[_METHOD_GET].Match(ctx)
// 	} else if req.Method[0] == 'H' {
// 		route = r.root[_METHOD_HEAD].Match(ctx)
// 	} else if req.Method[0] == 'D' {
// 		route = r.root[_METHOD_DELETE].Match(ctx)
// 	} else if req.Method[0] == 'C' {
// 		route = r.root[_METHOD_CONNECT].Match(ctx)
// 	} else if req.Method[0] == 'O' {
// 		route = r.root[_METHOD_OPTIONS].Match(ctx)
// 	} else if req.Method[0] == 'T' {
// 		route = r.root[_METHOD_TRACE].Match(ctx)
// 	} else if req.Method[1] == 'O' {
// 		route = r.root[_METHOD_POST].Match(ctx)
// 	} else if req.Method[1] == 'U' {
// 		route = r.root[_METHOD_PUT].Match(ctx)
// 	} else if req.Method[1] == 'A' {
// 		route = r.root[_METHOD_PATCH].Match(ctx)
// 	}
// 	if route == nil {
// 		ctx.handleFunc = r.notfound
// 	} else {
// 		ctx.handleFunc = route.handleFunc
// 	}
// 	ctx.Next()
// 	r.ctx.Put(ctx)
// }

// func (r *rootRouter) GET(routePath string, handle ...HandleFunc) {
// 	return r.root[_METHOD_GET].Add(routePath, handle...)
// }

// func (r *rootRouter) HEAD(routePath string, handle ...HandleFunc) {
// 	return r.root[_METHOD_HEAD].Add(routePath, handle...)
// }

// func (r *rootRouter) POST(routePath string, handle ...HandleFunc) {
// 	return r.root[_METHOD_POST].Add(routePath, handle...)
// }

// func (r *rootRouter) PUT(routePath string, handle ...HandleFunc) {
// 	return r.root[_METHOD_PUT].Add(routePath, handle...)
// }

// func (r *rootRouter) PATCH(routePath string, handle ...HandleFunc) {
// 	return r.root[_METHOD_PATCH].Add(routePath, handle...)
// }

// func (r *rootRouter) DELETE(routePath string, handle ...HandleFunc) {
// 	return r.root[_METHOD_DELETE].Add(routePath, handle...)
// }

// func (r *rootRouter) CONNECT(routePath string, handle ...HandleFunc) {
// 	return r.root[_METHOD_CONNECT].Add(routePath, handle...)
// }

// func (r *rootRouter) OPTIONS(routePath string, handle ...HandleFunc) {
// 	return r.root[_METHOD_OPTIONS].Add(routePath, handle...)
// }

// func (r *rootRouter) TRACE(routePath string, handle ...HandleFunc) {
// 	return r.root[_METHOD_TRACE].Add(routePath, handle...)
// }

// func (r *rootRouter) Static(routePath, file string, cache bool) {
// 	fi, err := os.Stat(file)
// 	if err != nil {
// 		return err
// 	}
// 	// Is a file
// 	if !fi.IsDir() {
// 		fi, err := os.Stat(file)
// 		if err != nil {
// 			return err
// 		}
// 		if !cache {
// 			h := &FileHandler{file: file}
// 			err = r.root[_METHOD_GET].Add(routePath, h.Handle)
// 		} else {
// 			h, err := NewCacheHandlerFromFile(file)
// 			if err != nil {
// 				return err
// 			}
// 			err = r.root[_METHOD_GET].Add(routePath, h.Handle)
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		//
// 		if fi.Name() == "index.html" {
// 			file, _ = filepath.Split(file)
// 			if !cache {
// 				h := &FileHandler{file: file}
// 				return r.root[_METHOD_GET].Add(routePath, h.Handle)
// 			} else {
// 				h, err := NewCacheHandlerFromFile(file)
// 				if err != nil {
// 					return err
// 				}
// 				return r.root[_METHOD_GET].Add(routePath, h.Handle)
// 			}
// 		}
// 		return nil
// 	}
// 	// Is a dir
// 	fis, err := ioutil.ReadDir(file)
// 	if err != nil {
// 		return err
// 	}
// 	// Add sub
// 	for i := 0; i < len(fis); i++ {
// 		err = r.Static(path.Join(routePath, fis[i].Name()), filepath.Join(file, fis[i].Name()), cache)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (r *rootRouter) NotFound(handle ...HandleFunc) {
// 	if len(handle) != 0 {
// 		r.notfound = handle
// 	} else {
// 		r.notfound = notfoundHandlerFunc
// 	}
// }

// func (r *rootRouter) Global(handle ...HandleFunc) {
// 	r.global = handle
// }

// func (r *rootRouter) SubRouter(routePath string) Router {
// 	return &subRouter{path: routePath, root: r}
// }
