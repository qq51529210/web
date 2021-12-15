package router

import (
	"fmt"
	"strings"
)

type route struct {
	handleFunc  []HandleFunc
	path        string
	staticChild []*route
	paramChild  *route
}

func (r *route) Match(ctx *Context) *route {
	path := ctx.Request.URL.Path
	route := r
	//
	if len(path) < len(route.path) || path[:len(route.path)] != route.path {
		return nil
	}
	path = path[len(route.path):]
	if path == "" {
		return route
	}
Loop:
	for {
		if route.paramChild != nil {
			if route.paramChild.path == ":" {
				i := strings.IndexByte(path, '/')
				if i < 0 {
					ctx.Param = append(ctx.Param, path)
					return route.paramChild
				}
				ctx.Param = append(ctx.Param, path[:i])
				// skip '/'
				path = path[i+1:]
				if path == "" {
					return route.paramChild
				}
				route = route.paramChild
				continue Loop
			}
			if route.paramChild.path == "*" {
				ctx.Param = append(ctx.Param, path)
				return route.paramChild
			}
		}
		for _, child := range route.staticChild {
			if len(path) < len(child.path) || path[:len(child.path)] != child.path {
				continue
			}
			path = path[len(child.path):]
			if path == "" {
				return child
			}
			route = child
			continue Loop
		}
		return nil
	}
}

func (r *route) Add(routePath string, handle ...HandleFunc) error {
	routePaths := splitRoutePath(routePath)
	//
	current := r
	for _, p := range routePaths {
		var child *route
		if p == ":" || p == "*" {
			child = current.addParam(p)
		} else {
			child = current.addStatic(p)
		}
		if child == nil {
			return fmt.Errorf("add route %s fail", routePath)
		}
		current = child
	}
	current.handleFunc = handle
	return nil
}

func (r *route) addStatic(routePath string) *route {
	//
	if r.path == "" {
		r.path = routePath
		return r
	}
	// case 1
	if r.path == routePath {
		return r
	}
	diff1, diff2 := diffString(r.path, routePath)
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
	if diff1 == "" {
		// "/a/:" or "/a/*"
		if r.paramChild != nil {
			return nil
		}
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

func (r *route) addParam(routePath string) *route {
	// "/a/b"
	if len(r.staticChild) != 0 {
		return nil
	}
	// "/a/:" or "/a/*"
	if r.paramChild != nil {
		// "/a/*"
		if r.paramChild.path == "*" {
			if routePath != "*" {
				return nil
			}
			return r.paramChild
		}
		// "/a/:"
		if routePath != ":" {
			return nil
		}
		return r.paramChild
	}
	r.paramChild = new(route)
	r.paramChild.path = routePath
	return r.paramChild
}
