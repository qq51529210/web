# http-router

A http router written in GO。

## Useage

```go
root := router.NewRootRouter()
// global handler
root.Global(func (ctx *Context) {
    t := time.Now()
    ctx.Next()
    fmt.Println(time.Now().Sub(t1))
})
// handle not match, default only response 404
root.NotFound(func (ctx *Context) {
	ctx.WriteHeader(404)
})
// handle static files
root.Static("staic", "http_static_root_dir", true)
// "github.com/login" and "github.com/qq51529210"
root.GET("/login", handleLogin)
root.GET("/?", handleUser)
// import "handler/foo"
foo.Init(root.SubRouter("/api/foo"))
// server
server := web.NewServer(":80", root)
server.Serve()
```

```go
// foo package
func Init(router router.Router) {
    router.GET("", list)
    router.GET("?", get)
    router.POST("", add)
}

func list(ctx *router.Context) {

}

func get(ctx *router.Context) {

}

func add(ctx *router.Context) {

}
// 
server := web.NewServer(":80", root)
server.Serve()
```

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
