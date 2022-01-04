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

func test_Error(t *testing.T, err error) {
	if err == nil {
		t.FailNow()
	} else {
		t.Log(err)
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
	// add
	{
		test_Fatal(t, r.Add(""))
		test_Fatal(t, r.Add("/"))
		test_Fatal(t, r.Add("/a"))
		test_Fatal(t, r.Add("/a"))
		test_Error(t, r.Add("/?"))
		// "/a/" -> "/a"
		test_Fatal(t, r.Add("/a/"))
		test_Fatal(t, r.Add("/ab"))
		test_Fatal(t, r.Add("/a/b"))
		test_Fatal(t, r.Add("/a/b/c"))
		test_Fatal(t, r.Add("/b"))
		test_Fatal(t, r.Add("/b/?"))
		test_Fatal(t, r.Add("/b/?"))
		test_Error(t, r.Add("/b/b"))
		test_Fatal(t, r.Add("/b/?/?"))
		test_Fatal(t, r.Add("/bc/?"))
	}
	// match
	c.Request.URL.Path = "/"
	test_Fail(t, r.Match(c) == nil)
	c.Request.URL.Path = "/a"
	test_Fail(t, r.Match(c) == nil)
	c.Request.URL.Path = "/ab"
	test_Fail(t, r.Match(c) == nil)
	c.Request.URL.Path = "/a/b"
	test_Fail(t, r.Match(c) == nil)
	c.Request.URL.Path = "/a/b/c"
	test_Fail(t, r.Match(c) == nil)
	c.Request.URL.Path = "/b"
	test_Fail(t, r.Match(c) == nil)
	c.Request.URL.Path = "/b/1"
	c.Param = c.Param[:0]
	test_Fail(t, r.Match(c) == nil, len(c.Param) != 1 || c.Param[0] != "1")
	c.Request.URL.Path = "/b/b"
	test_Fail(t, r.Match(c) == nil)
	c.Request.URL.Path = "/b/1/2"
	c.Param = c.Param[:0]
	test_Fail(t, r.Match(c) == nil, len(c.Param) != 2 || c.Param[0] != "1" || c.Param[1] != "2")
	c.Request.URL.Path = "/bc/1"
	c.Param = c.Param[:0]
	test_Fail(t, r.Match(c) == nil, len(c.Param) != 1 || c.Param[0] != "1")
}
