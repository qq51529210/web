package router

import (
	"path"
)

type sub struct {
	path   string
	router *router
}

func (r *sub) Sub(routePath string) Router {
	return &sub{path: routePath, router: r.router}
}

func (r *sub) GET(routePath string, handle ...HandleFunc) error {
	return r.router.GET(path.Join(r.path, routePath), handle...)
}

func (r *sub) HEAD(routePath string, handle ...HandleFunc) error {
	return r.router.HEAD(path.Join(r.path, routePath), handle...)
}

func (r *sub) POST(routePath string, handle ...HandleFunc) error {
	return r.router.POST(path.Join(r.path, routePath), handle...)
}

func (r *sub) PUT(routePath string, handle ...HandleFunc) error {
	return r.router.PUT(path.Join(r.path, routePath), handle...)
}

func (r *sub) PATCH(routePath string, handle ...HandleFunc) error {
	return r.router.PATCH(path.Join(r.path, routePath), handle...)
}

func (r *sub) DELETE(routePath string, handle ...HandleFunc) error {
	return r.router.DELETE(path.Join(r.path, routePath), handle...)
}

func (r *sub) CONNECT(routePath string, handle ...HandleFunc) error {
	return r.router.CONNECT(path.Join(r.path, routePath), handle...)
}

func (r *sub) OPTIONS(routePath string, handle ...HandleFunc) error {
	return r.router.OPTIONS(path.Join(r.path, routePath), handle...)
}

func (r *sub) TRACE(routePath string, handle ...HandleFunc) error {
	return r.router.TRACE(path.Join(r.path, routePath), handle...)
}
