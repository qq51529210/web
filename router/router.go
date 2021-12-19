package router

import (
	"path"
	"strings"
)

type Router interface {
	Sub(routePath string) Router
	GET(routePath string, handle ...HandleFunc) error
	HEAD(routePath string, handle ...HandleFunc) error
	POST(routePath string, handle ...HandleFunc) error
	PUT(routePath string, handle ...HandleFunc) error
	PATCH(routePath string, handle ...HandleFunc) error
	DELETE(routePath string, handle ...HandleFunc) error
	CONNECT(routePath string, handle ...HandleFunc) error
	OPTIONS(routePath string, handle ...HandleFunc) error
	TRACE(routePath string, handle ...HandleFunc) error
}

func NewRoot() Root {
	r := new(rootRouter)
	r.ctx.New = func() interface{} {
		return new(Context)
	}
	for i := 0; i < _METHOD_INVALID; i++ {
		r.root[i] = new(route)
	}
	return r
}

func splitRoutePath(routePath string) []string {
	if routePath == "" || routePath == "/" {
		return []string{"/"}
	}
	routePath = path.Clean(routePath)
	//
	var routePaths []string
	var static []string
	for _, p := range strings.Split(routePath, "/") {
		if p == "" {
			continue
		}
		if p[0] == ':' || p[0] == '*' {
			if len(static) > 0 {
				staticPath := strings.Join(static, "/")
				static = static[:0]
				if len(routePaths) == 0 {
					routePaths = append(routePaths, "/"+staticPath+"/")
				} else {
					if routePaths[len(routePaths)-1] == ":" ||
						routePaths[len(routePaths)-1] == "*" {
						routePaths = append(routePaths, staticPath+"/")
					}
				}
			} else {
				if len(routePaths) == 0 {
					routePaths = append(routePaths, "/")
				}
			}
			routePaths = append(routePaths, string(p[0]))
			if p[0] == '*' {
				break
			}
			continue
		}
		static = append(static, p)
	}
	if len(static) > 0 {
		staticPath := strings.Join(static, "/")
		if len(routePaths) == 0 {
			routePaths = append(routePaths, "/"+staticPath)
		} else {
			if routePaths[len(routePaths)-1] == ":" ||
				routePaths[len(routePaths)-1] == "*" {
				routePaths = append(routePaths, staticPath)
			}
		}
	}
	return routePaths
}

func diffString(s1, s2 string) (string, string) {
	n := len(s1)
	if n > len(s2) {
		n = len(s2)
	}
	i := 0
	for ; i < n; i++ {
		if s1[i] != s2[i] {
			break
		}
	}
	return s1[i:], s2[i:]
}
