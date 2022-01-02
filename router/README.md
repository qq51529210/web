# http-router

A http router written in GO。

## Useage

```go
root := router.NewRootRouter()
root.Global(func (ctx *Context) {
    t := time.Now()
    ctx.Next()
    fmt.Println(time.Now().Sub(t1))
})
root.Static("staic", "http_static_root_dir", true)
root.GET("login", loginHandleFunc)
// 
users := root.SubRouter("/users")
users.GET("", getUsersHandleFunc)
users.GET("/:", getUserHandleFunc)
users.POST("", addUsersHandleFunc)
// 
server := web.NewServer(":80", root)
server.Serve()
// or 
yourServer.httpServer.Handler = root
```

## Benchmark

- /static-0-count/static-0-deep
- /param-0-count/:parma-0-deep
- /half-0-count/static-0-deep/:param-0-deep

```golang
// UrlCount = 5 * 3, UrlDeep = 2
Benchmark_My-4           1296975               933 ns/op               0 B/op          0 allocs/op
Benchmark_Gin-4           927745              1102 ns/op               0 B/op          0 allocs/op
// UrlCount = 10 * 3, UrlDeep = 3
Benchmark_My-4            469432              2242 ns/op               0 B/op          0 allocs/op
Benchmark_Gin-4           431293              2465 ns/op               0 B/op          0 allocs/op
// UrlCount = 20 * 3, UrlDeep = 5
Benchmark_My-4            223212              5332 ns/op               0 B/op          0 allocs/op
Benchmark_Gin-4           199281              5566 ns/op               0 B/op          0 allocs/op
// UrlCount = 30 * 3, UrlDeep = 7
Benchmark_My-4            122394              9305 ns/op               0 B/op          0 allocs/op
Benchmark_Gin-4           119985              9731 ns/op               0 B/op          0 allocs/op
```
