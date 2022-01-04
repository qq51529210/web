package router

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"testing"
	// "github.com/gin-gonic/gin"
)

type testHandler struct {
	header http.Header
	buffer bytes.Buffer
	funcs  []string
	param  []string
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

func (h *testHandler) WriteHeader(int) {
}

func (h *testHandler) Reset() {
	h.header = make(http.Header)
	h.buffer.Reset()
	h.funcs = make([]string, 0)
}

func Test_Router(t *testing.T) {
	root := NewRootRouter()
	root.NotFound(func(ctx *Context) {
		ctx.WriteHeader(200)
	})
	handler := new(testHandler)
	request := new(http.Request)
	request.URL = new(url.URL)
	//
	request.URL.Path = "/v0/users"
	request.Method = http.MethodGet
	root.GET("/v0/users", func(ctx *Context) {})
	root.ServeHTTP(handler, request)
	//
	request.URL.Path = "/v0/users"
	request.Method = http.MethodPost
	root.POST("/v0/users", func(ctx *Context) {})
	root.ServeHTTP(handler, request)
	//
	request.URL.Path = "/v0/users/root"
	request.Method = http.MethodGet
	root.GET("/v0/users", func(ctx *Context) {})
	root.ServeHTTP(handler, request)
}

func Test_Router_Static(t *testing.T) {
	root := NewRootRouter()
	root.NotFound(func(ctx *Context) {
		t.FailNow()
	})
	handler := new(testHandler)
	request := new(http.Request)
	request.URL = new(url.URL)
	request.Method = http.MethodGet
	root.Static("/static", ".", false)
	fis, err := ioutil.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(fis); i++ {
		request.URL.Path = fmt.Sprintf("/static/%s", fis[i].Name())
		handler.Reset()
		root.ServeHTTP(handler, request)
		data, err := ioutil.ReadFile(filepath.Join(".", fis[i].Name()))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(data, handler.buffer.Bytes()) {
			t.FailNow()
		}
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
