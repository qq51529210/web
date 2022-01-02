# http-router

A http router written in GO。

## Useage

```go
root := router.NewRoot()
root.Global(func (ctx *Context) {
    t := time.Now()
    ctx.Next()
    fmt.Println(time.Now().Sub(t1))
})
root.Static("staic", "http_static_root_dir", true)
root.GET("login", loginHandleFunc)
// 
users := root.Sub("/users")
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

### UrlCount = 10 * 3, UrlDeep = 3

```golang
goos: darwin
goarch: amd64
pkg: github.com/qq51529210/web/router
Benchmark_My-4            477474              2220 ns/op               0 B/op          0 allocs/op
Benchmark_Gin-4           391414              2562 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/qq51529210/web/router        2.797s
```

### UrlCount = 20 * 3, UrlDeep = 5

```golang
goos: darwin
goarch: amd64
pkg: github.com/qq51529210/web/router
Benchmark_My-4            215452              5094 ns/op               0 B/op          0 allocs/op
Benchmark_Gin-4           214322              5867 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/qq51529210/web/router        2.834s
```

### UrlCount = 30 * 3, UrlDeep = 7

```golang
goos: darwin
goarch: amd64
pkg: github.com/qq51529210/web/router
Benchmark_My-4            134044              8996 ns/op               0 B/op          0 allocs/op
Benchmark_Gin-4           125269              9564 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/qq51529210/web/router        3.021s
```

### UrlCount = 50 * 3, UrlDeep = 10

```golang
goos: darwin
goarch: amd64
pkg: github.com/qq51529210/web/router
Benchmark_My-4             69570             16809 ns/op               0 B/op          0 allocs/op
Benchmark_Gin-4            64596             17831 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/qq51529210/web/router        3.080s
```
