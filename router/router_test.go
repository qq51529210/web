package router

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"testing"
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
	root := NewRoot()
	root.NotFound(func(ctx *Context) {
		ctx.WriteHeader(200)
	})
	handler := new(testHandler)
	request := new(http.Request)
	request.URL = new(url.URL)
	//
	request.URL.Path = "/v0/users"
	request.Method = http.MethodGet
	err := root.GET("/v0/users", func(ctx *Context) {})
	if err != nil {
		t.Fatal(err)
	}
	root.ServeHTTP(handler, request)
	//
	request.URL.Path = "/v0/users"
	request.Method = http.MethodPost
	err = root.POST("/v0/users", func(ctx *Context) {})
	if err != nil {
		t.Fatal(err)
	}
	root.ServeHTTP(handler, request)
	//
	request.URL.Path = "/v0/users/root"
	request.Method = http.MethodGet
	err = root.GET("/v0/users", func(ctx *Context) {})
	if err != nil {
		t.Fatal(err)
	}
	root.ServeHTTP(handler, request)
}

func Test_Router_Static(t *testing.T) {
	root := NewRoot()
	root.NotFound(func(ctx *Context) {
		t.FailNow()
	})
	handler := new(testHandler)
	request := new(http.Request)
	request.URL = new(url.URL)
	request.Method = http.MethodGet
	err := root.Static(http.MethodGet, "/static", ".", false)
	if err != nil {
		t.Fatal(err)
	}
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

func Benchmark_Router_Match_Static(b *testing.B) {
	// 100 route paths
	routePathCount := 100
	// 10 route deep
	routePathDeep := 10
	root := NewRoot()
	var urls []string
	for i := 0; i < routePathCount; i++ {
		staticRoute := "/static"
		paramRoute := "/param"
		halfRoute := "/half"
		staticUrl := "/static"
		paramUrl := "/param"
		halfUrl := "/half"
		for j := 0; j < routePathDeep; j++ {
			staticRoute += fmt.Sprintf("/static%d%d", i, j)
			staticUrl += fmt.Sprintf("/static%d%d", i, j)
			paramRoute += fmt.Sprintf("/:")
			paramUrl += fmt.Sprintf("/param%d%d", i, j)
			halfRoute += fmt.Sprintf("/static%d%d/:", i, j)
			halfUrl += fmt.Sprintf("/static%d%d/param%d%d", i, j, i, j)
		}
		root.GET(staticRoute)
		root.GET(paramRoute)
		root.GET(halfRoute)
		urls = append(urls, staticUrl)
		urls = append(urls, paramUrl)
		urls = append(urls, halfUrl)
	}
	//
	h := new(testHandler)
	r := new(http.Request)
	r.Method = http.MethodGet
	r.URL = new(url.URL)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < len(urls); i++ {
		r.URL.Path = urls[i]
		for j := 0; j < b.N; j++ {
			root.ServeHTTP(h, r)
		}
	}
}
