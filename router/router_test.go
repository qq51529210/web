package router

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	// "github.com/gin-gonic/gin"
)

type testHandler struct {
	header http.Header
	code   int
	buffer strings.Builder
	param  []string
	req    *http.Request
}

func (h *testHandler) Header() http.Header {
	return h.header
}

func (h *testHandler) Write(b []byte) (n int, err error) {
	return h.buffer.Write(b)
}

func (h *testHandler) WriteString(s string) (n int, err error) {
	return h.buffer.WriteString(s)
}

func (h *testHandler) WriteHeader(code int) {
	h.code = code
}

func (h *testHandler) Reset() {
	h.header = make(http.Header)
	h.code = 0
	h.buffer.Reset()
}

func newTestHandler() *testHandler {
	h := new(testHandler)
	h.req = new(http.Request)
	h.req.Method = http.MethodGet
	h.req.URL = new(url.URL)
	h.Reset()
	return h
}

func Test_Router(t *testing.T) {
	h := newTestHandler()
	r := NewRootRouter()
	// global
	g1, g2 := 0, 0
	r.Intercept(func(ctx *Context) {
		g1++
		ctx.Next()
	})
	// notfound
	r.NotFound(func(ctx *Context) {
		ctx.WriteHeader(404)
	})
	// static, cache file which size less than 2kb.
	r.Static("/static", ".", 1024*2)
	// handle
	r.GET("/login", func(ctx *Context) {
		io.WriteString(ctx.ResponseWriter, "get login")
	})
	r.GET("/?", func(ctx *Context) {
		io.WriteString(ctx.ResponseWriter, "get /"+ctx.Param[0])
	})
	// sub router
	s := r.SubRouter("/foo")
	s.GET("", func(ctx *Context) {
		io.WriteString(ctx.ResponseWriter, "get foo list")
	})
	s.Intercept(func(ctx *Context) {
		g2++
		ctx.Next()
	})
	s.GET("?", func(ctx *Context) {
		io.WriteString(ctx.ResponseWriter, "get foo "+ctx.Param[0])
	})
	//
	fis, err := ioutil.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(fis); i++ {
		h.req.URL.Path = fmt.Sprintf("/static/%s", fis[i].Name())
		h.Reset()
		r.ServeHTTP(h, h.req)
		d, err := ioutil.ReadFile(filepath.Join(".", fis[i].Name()))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(d, []byte(h.buffer.String())) {
			t.FailNow()
		}
	}
	//
	g1 = 0
	g2 = 0
	//
	h.req.URL.Path = "/login"
	h.Reset()
	r.ServeHTTP(h, h.req)
	if h.code == 404 || h.buffer.String() != "get login" || g1 != 1 || g2 != 0 {
		t.FailNow()
	}
	//
	h.req.URL.Path = "/qq51529210"
	h.Reset()
	r.ServeHTTP(h, h.req)
	if h.code == 404 || h.buffer.String() != "get /qq51529210" || g1 != 2 || g2 != 0 {
		t.FailNow()
	}
	//
	h.req.URL.Path = "/foo"
	h.Reset()
	r.ServeHTTP(h, h.req)
	if h.code == 404 || h.buffer.String() != "get foo list" || g1 != 3 || g2 != 0 {
		t.FailNow()
	}
	h.req.URL.Path = "/foo/qq51529210"
	h.Reset()
	r.ServeHTTP(h, h.req)
	if h.code == 404 || h.buffer.String() != "get foo qq51529210" || g1 != 4 || g2 != 1 {
		t.FailNow()
	}
}

func benchmarkRoutePaths(paramName string) ([]string, []string) {
	var routePathCount, routePathDeep = 5, 2
	// var routePathCount, routePathDeep = 10, 3
	// var routePathCount, routePathDeep = 20, 5
	// var routePathCount, routePathDeep = 30, 7
	var routes, urls []string
	for i := 0; i < routePathCount; i++ {
		staticRoute := fmt.Sprintf("/static%d", i)
		paramRoute := fmt.Sprintf("/param%d", i)
		halfRoute := fmt.Sprintf("/half%d", i)
		staticUrl := staticRoute
		paramUrl := paramRoute
		halfUrl := halfRoute
		for j := 0; j < routePathDeep; j++ {
			staticRoute += fmt.Sprintf("/static%d", j)
			staticUrl += fmt.Sprintf("/static%d", j)
			paramRoute += fmt.Sprintf("/%s%d", paramName, j)
			paramUrl += fmt.Sprintf("/param%d", j)
			halfRoute += fmt.Sprintf("/static%d/%s%d", j, paramName, j)
			halfUrl += fmt.Sprintf("/static%d/param%d", j, i)
		}
		routes = append(routes, staticRoute)
		routes = append(routes, paramRoute)
		routes = append(routes, halfUrl)
		urls = append(urls, staticUrl)
		urls = append(urls, paramUrl)
		urls = append(urls, halfUrl)
	}
	return routes, urls
}

func benchmarkServeHTTP(b *testing.B, handler http.Handler, urls []string) {
	h := new(testHandler)
	h.header = make(http.Header)
	r := new(http.Request)
	r.Method = http.MethodGet
	r.URL = new(url.URL)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < len(urls); i++ {
		r.URL.Path = urls[i]
		for j := 0; j < b.N; j++ {
			handler.ServeHTTP(h, r)
		}
	}
}

func Benchmark_My(b *testing.B) {
	routes, urls := benchmarkRoutePaths("?")
	root := NewRootRouter()
	root.NotFound(func(ctx *Context) { b.FailNow() })
	for i := 0; i < len(routes); i++ {
		root.GET(routes[i], func(ctx *Context) {})
	}
	benchmarkServeHTTP(b, root, urls)
}

// func Benchmark_Gin(b *testing.B) {
// 	routes, urls := benchmarkRoutePaths(":p")
// 	gin.SetMode(gin.ReleaseMode)
// 	root := gin.New()
// 	root.NoMethod(func(c *gin.Context) { b.FailNow() })
// 	for i := 0; i < len(routes); i++ {
// 		root.GET(routes[i], func(c *gin.Context) {})
// 	}
// 	benchmarkServeHTTP(b, root, urls)
// }
