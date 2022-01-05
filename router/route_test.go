package router

import (
	"net/http"
	"net/url"
	"testing"
)

func test_Fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func test_Fail(t *testing.T, oks ...bool) {
	for _, ok := range oks {
		if ok {
			t.FailNow()
		}
	}
}

func Test_Route_Add_Match(t *testing.T) {
	r := new(route)
	c := new(Context)
	c.Request = new(http.Request)
	c.Request.URL = new(url.URL)
	r.Add("")
	r.Add("/")
	c.Request.URL.Path = "/"
	test_Fail(t, r.Match(c) == nil)
	r.Add("/a")
	r.Add("/?")
	r.Add("/?abcd")
	c.Request.URL.Path = "/a"
	test_Fail(t, r.Match(c) == nil)
	c.Request.URL.Path = "/1"
	test_Fail(t, r.Match(c) == nil, len(c.Param) != 1 || c.Param[0] != "1")
	r.Add("/ab")
	c.Request.URL.Path = "/ab"
	test_Fail(t, r.Match(c) == nil)
	r.Add("/a/b")
	c.Request.URL.Path = "/a/b"
	test_Fail(t, r.Match(c) == nil)
	r.Add("/b")
	c.Request.URL.Path = "/b"
	test_Fail(t, r.Match(c) == nil)
	c.Request.URL.Path = "/b/"
	test_Fail(t, r.Match(c) != nil) // noted
	r.Add("/b/?")
	c.Request.URL.Path = "/b/1"
	test_Fail(t, r.Match(c) == nil, len(c.Param) != 1 || c.Param[0] != "1")
	c.Request.URL.Path = "/b/"
	test_Fail(t, r.Match(c) == nil) // noted
	r.Add("/b/b")
	c.Request.URL.Path = "/b/b"
	test_Fail(t, r.Match(c) == nil)
	r.Add("/b/b/*")
	c.Request.URL.Path = "/b/b/1/2"
	test_Fail(t, r.Match(c) == nil, len(c.Param) != 1 || c.Param[0] != "1/2")
	r.Add("/b/?/?/b")
	c.Request.URL.Path = "/b/1/2/b"
	test_Fail(t, r.Match(c) == nil, len(c.Param) != 2 || c.Param[0] != "1" || c.Param[1] != "2")
}
