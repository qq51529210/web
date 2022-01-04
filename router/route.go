package router

import (
	"fmt"
	"path"
	"strings"
)

const paramChar = '?'

type route struct {
	handleFunc  []HandleFunc
	path        string
	staticChild []*route
	paramChild  *route
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
			continue
		}
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
		return nil
	}
}

func (r *route) Add(routePath string, handle ...HandleFunc) error {
	routePath = path.Clean(path.Join("/", routePath))
	_routePath := routePath
	//
	var routePaths []string
	for _routePath != "" {
		i := strings.IndexByte(_routePath, paramChar)
		if i < 0 {
			routePaths = append(routePaths, _routePath)
			break
		}
		if _routePath[:i] != "" {
			routePaths = append(routePaths, _routePath[:i])
		}
		routePaths = append(routePaths, string(paramChar))
		_routePath = strings.TrimLeftFunc(_routePath[i:], func(r rune) bool { return r != '/' })
		if _routePath == "" {
			break
		}
		// skip '/'
		_routePath = _routePath[1:]
	}
	//
	var rootPath strings.Builder
	var child *route
	var err error
	current := r
	if current.path != "" {
		rootPath.WriteString(current.path)
	}
	for _, p := range routePaths {
		//
		if p[0] == paramChar {
			child, err = current.addParam(p, &rootPath)
		} else {
			child, err = current.addStatic(p, &rootPath)
		}
		if err != nil {
			return fmt.Errorf(`add "%s" failed, %v`, routePath, err)
		}
		if current != child {
			if current.path[0] == paramChar {
				rootPath.WriteByte('/')
			}
			rootPath.WriteString(current.path)
			current = child
		}
	}
	current.handleFunc = handle
	return nil
}

func (r *route) addStatic(routePath string, rootPath *strings.Builder) (*route, error) {
	// root
	if r.path == "" {
		r.path = routePath
		return r, nil
	}
	// case 1
	if r.path == routePath {
		return r, nil
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
		return child, nil
	}
	// case 3, r.path="/ab", routePath="/abc", diff1="", diff2="c".
	if diff1 == "" {
		if r.paramChild != nil {
			return nil, fmt.Errorf(`"%v" has sub parameter "?"`, rootPath)
		}
		for _, child := range r.staticChild {
			if child.path[0] == diff2[0] {
				rootPath.WriteString(child.path)
				return child.addStatic(diff2, rootPath)
			}
		}
		route := &route{path: diff2}
		r.staticChild = append(r.staticChild, route)
		return route, nil
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
	return child2, nil
}

func (r *route) addParam(routePath string, rootPath *strings.Builder) (*route, error) {
	if len(r.staticChild) != 0 {
		return nil, fmt.Errorf(`"%v" has sub static "%s"`, rootPath, r.staticChild[0].path)
	}
	if r.paramChild == nil {
		r.paramChild = new(route)
		r.paramChild.path = routePath
	}
	return r.paramChild, nil
}
