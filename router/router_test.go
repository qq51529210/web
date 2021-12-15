package router

import (
	"bytes"
	"net/http"
	"net/url"
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
