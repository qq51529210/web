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
	test_Route_Add(t, r)
	test_Route_Match(t, r)
}

func test_Route_Add(t *testing.T, r *route) {
	test_Fatal(t, r.Add("/a/1"))
	test_Fatal(t, r.Add("/a/1/:"))
	test_Fatal(t, r.Add("/a/1/:/*"))
	test_Fatal(t, r.Add("/a/12"))
	test_Fail(t, r.Add("/a/:") == nil)
	test_Fail(t, r.Add("/a/1/*") == nil)
	test_Fail(t, r.Add("/a/1/:/:") == nil)
	test_Fail(t, r.Add(":") == nil)
	test_Fail(t, r.Add("*") == nil)
	test_Fatal(t, r.Add(""))
	test_Fatal(t, r.Add("/"))
	test_Fatal(t, r.Add("/b/1"))
	test_Fatal(t, r.Add("/b/1/:"))
	test_Fatal(t, r.Add("/b/1/:/*"))
	test_Fatal(t, r.Add("/b/12"))
}

func test_Route_Match(t *testing.T, r *route) {
	c := new(Context)
	c.Request = new(http.Request)
	c.Request.URL = new(url.URL)
	c.Request.URL.Path = "/a/1"
	test_Fail(t, r.Match(c) == nil)
	c.Request.URL.Path = "/a/1/1"
	c.Param = c.Param[:0]
	test_Fail(t, r.Match(c) == nil, len(c.Param) != 1 || c.Param[0] != "1")
	c.Request.URL.Path = "/a/1/2/3/4/5"
	c.Param = c.Param[:0]
	test_Fail(t, r.Match(c) == nil, len(c.Param) != 2 || c.Param[0] != "2" || c.Param[1] != "3/4/5")
	c.Request.URL.Path = "/a/12"
	test_Fail(t, r.Match(c) == nil)
	c.Request.URL.Path = "/"
	test_Fail(t, r.Match(c) == nil)
	c.Request.URL.Path = "/b/1"
	test_Fail(t, r.Match(c) == nil)
	c.Request.URL.Path = "/b/1/1"
	c.Param = c.Param[:0]
	test_Fail(t, r.Match(c) == nil, len(c.Param) != 1 || c.Param[0] != "1")
	c.Request.URL.Path = "/b/1/2/3/4/5"
	c.Param = c.Param[:0]
	test_Fail(t, r.Match(c) == nil, len(c.Param) != 2 || c.Param[0] != "2" || c.Param[1] != "3/4/5")
	c.Request.URL.Path = "/b/12"
	test_Fail(t, r.Match(c) == nil)
}
