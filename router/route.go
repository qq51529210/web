package router

import (
	"path"
	"strings"
)

const holderChar = "?"
const anyChar = "*"
const paramChars = holderChar + anyChar

type route struct {
	handleFunc  []HandleFunc
	path        string
	staticChild []*route
	paramChild  *route
	anyChild    *route
}

func (r *route) Match(ctx *Context) *route {
	_path := ctx.Request.URL.Path
	// root
	if len(_path) < len(r.path) || _path[:len(r.path)] != r.path {
		return nil
	}
	_path = _path[len(r.path):]
	if _path == "" {
		return r
	}
	ctx.Param = ctx.Param[:0]
	idx := 0
	route := r
Loop:
	for {
		for _, child := range route.staticChild {
			if len(_path) < len(child.path) || _path[:len(child.path)] != child.path {
				continue
			}
			_path = _path[len(child.path):]
			if _path == "" {
				return child
			}
			route = child
			continue Loop
		}
		if route.paramChild != nil {
			idx = strings.IndexByte(_path, '/')
			if idx < 0 {
				ctx.Param = append(ctx.Param, _path)
				return route.paramChild
			}
			ctx.Param = append(ctx.Param, _path[:idx])
			// skip '/'
			_path = _path[idx+1:]
			if _path == "" {
				return route.paramChild
			}
			route = route.paramChild
			continue Loop
		}
		if route.anyChild != nil {
			ctx.Param = append(ctx.Param, _path)
			return route.anyChild
		}
		return nil
	}
}

func (r *route) Add(routePath string, handle ...HandleFunc) {
	_routePath := path.Clean(path.Join("/", routePath))
	//
	var routePaths []string
	for _routePath != "" {
		i := strings.IndexAny(_routePath, paramChars)
		if i < 0 {
			routePaths = append(routePaths, _routePath)
			break
		}
		if i != 0 && _routePath[:i] != "" {
			routePaths = append(routePaths, _routePath[:i])
		}
		routePaths = append(routePaths, _routePath[i:i+1])
		_routePath = strings.TrimLeftFunc(_routePath[i:], func(r rune) bool { return r != '/' })
		if _routePath == "" {
			break
		}
		// skip '/'
		_routePath = _routePath[1:]
	}
	//
	current := r
	for _, p := range routePaths {
		//
		switch p {
		case holderChar:
			if current.paramChild == nil {
				current.paramChild = new(route)
				current.paramChild.path = p
			}
			current = current.paramChild
		case anyChar:
			if current.anyChild == nil {
				current.anyChild = new(route)
				current.anyChild.path = p
			}
			current = current.anyChild
		default:
			current = current.addStatic(p)
		}
	}
	current.handleFunc = handle
}

func (r *route) addStatic(routePath string) *route {
	// root
	if r.path == "" {
		r.path = routePath
		return r
	}
	// case 1
	if r.path == routePath {
		return r
	}
	// difference string
	n := len(r.path)
	if n > len(routePath) {
		n = len(routePath)
	}
	i := 0
	for ; i < n; i++ {
		if r.path[i] != routePath[i] {
			break
		}
	}
	diff1, diff2 := r.path[i:], routePath[i:]
	// case 2, r.path="/abc", routePath="/ab", diff1="c", diff2=""
	if diff2 == "" {
		child := new(route)
		child.handleFunc = r.handleFunc
		child.path = r.path[len(routePath):]
		child.staticChild = r.staticChild
		child.paramChild = r.paramChild
		//
		r.handleFunc = nil
		r.path = routePath
		r.staticChild = make([]*route, 1)
		r.staticChild[0] = child
		r.paramChild = nil
		return child
	}
	// case 3, r.path="/ab", routePath="/abc", diff1="", diff2="c".
	if diff1 == "" || diff1 == "?" {
		for _, child := range r.staticChild {
			if child.path[0] == diff2[0] {
				return child.addStatic(diff2)
			}
		}
		route := &route{path: diff2}
		r.staticChild = append(r.staticChild, route)
		return route
	}
	// case 4, r.path="/abc", routePath="/abd", diff1="c", diff2="d".
	child1 := new(route)
	child1.handleFunc = r.handleFunc
	child1.path = diff1
	child1.staticChild = r.staticChild
	child1.paramChild = r.paramChild
	//
	child2 := new(route)
	child2.handleFunc = nil
	child2.path = diff2
	child2.staticChild = make([]*route, 0)
	child2.paramChild = nil
	//
	r.handleFunc = nil
	r.path = r.path[:len(r.path)-len(diff1)]
	r.staticChild = make([]*route, 2)
	r.staticChild[0] = child1
	r.staticChild[1] = child2
	r.paramChild = nil
	return child2
}
