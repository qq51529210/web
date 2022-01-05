# http-router

A http router written in GO。

## Useage

```go
root := router.NewRootRouter()
// Handle static files.
// Before Intercept will not be intercepted.
root.Static("staic", "http_static_root_dir", true)
// Handle not match, default only response 404.
// Before Intercept will not be intercepted.
root.NotFound(func (ctx *Context) {
	ctx.WriteHeader(404)
})
// Global handler
root.Intercept(func (ctx *Context) {
    t := time.Now()
    ctx.Handle()
    fmt.Println(time.Now().Sub(t1))
})
// Example, "github.com/login" and "github.com/qq51529210".
// Static priority is higher.
root.GET("/login", handleLogin)
root.GET("/?", handleUser)
// Holder '?' priority is higher, so this will not match.
root.GET("/*", handleUser)
// import "handler/file"
file.Init(root.SubRouter("/api/files"))
// server
server := web.NewServer(":80", root)
server.Serve()
```

```go
// file package
func Init(r router.Router) {
    // Before Intercept will not be intercepted.
    // But still intercepted by global handler.
    r.GET("", list)
    r.GET("?", get)
    r.GET("dir/*", listDir)
    r.Intercept(parseForm)
    // Call parseForm and add
    r.POST("", add)
}

func parseForm(ctx *router.Context) {}

func list(ctx *router.Context) {}

func get(ctx *router.Context) {}

func listDir(ctx *router.Context) {
    // example: ctx.Request.URL.Path = "/api/files/dir/docs/test"
    // ctx.Param[0] = "docs/test"
}

func add(ctx *router.Context) {}
// Or your server.
web.NewServer(":80", root).Serve()
```

## Call chain priority

- root intercept > [sub intercept] > route handle
- root intercept > root notfound

## Benchmark

- /static-0-count/static-0-deep
- /param-0-count/:parma-0-deep
- /half-0-count/static-0-deep/:param-0-deep

```golang
// UrlCount = 5 * 3, UrlDeep = 2
Benchmark_My-4           1300068               927 ns/op               0 B/op          0 allocs/op
Benchmark_Gin-4           927745              1102 ns/op               0 B/op          0 allocs/op
// UrlCount = 10 * 3, UrlDeep = 3
Benchmark_My-4            455833              2280 ns/op               0 B/op          0 allocs/op
Benchmark_Gin-4           431293              2465 ns/op               0 B/op          0 allocs/op
// UrlCount = 20 * 3, UrlDeep = 5
Benchmark_My-4            205758              5365 ns/op               0 B/op          0 allocs/op
Benchmark_Gin-4           199281              5566 ns/op               0 B/op          0 allocs/op
// UrlCount = 30 * 3, UrlDeep = 7
Benchmark_My-4            121921              9018 ns/op               0 B/op          0 allocs/op
Benchmark_Gin-4           119985              9731 ns/op               0 B/op          0 allocs/op
```
