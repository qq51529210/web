package router

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"strings"
)

const (
	bearerTokenPrefix = "Bearer "
	contentType       = "Content-Type"
)

var (
	// ContentTypeUTF8 如其名
	ContentTypeUTF8 = "charset=utf-8"
	// ContentTypeJSON 如其名
	ContentTypeJSON = mime.TypeByExtension(".json")
	// ContentTypeHTML 如其名
	ContentTypeHTML = mime.TypeByExtension(".html")
	// ContentTypeJS 如其名
	ContentTypeJS = mime.TypeByExtension(".js")
	// ContentTypeCSS 如其名
	ContentTypeCSS = mime.TypeByExtension(".css")
)

// Context 表示调用链的上下文
type Context struct {
	// 标准库的实例
	*http.Request
	// 标准库的实例
	http.ResponseWriter
	// 动态路由的路径节点，按照注册时的顺序
	Param []string
	// 用于在调用链中保存临时数据
	TempData interface{}
	// 保存调用链函数
	handleFunc []HandleFunc
	// 当前调用的函数下标
	handleIdx int
}

// handle 执行调用链中剩下的所有函数
func (ctx *Context) handle() {
	for ctx.handleIdx < len(ctx.handleFunc) {
		ctx.handleFunc[ctx.handleIdx](ctx)
		ctx.handleIdx++
	}
}

// Handle 执行调用链的下一个函数
func (ctx *Context) Handle() {
	ctx.handleIdx++
	ctx.handle()
}

// Abort 终止调用链
func (ctx *Context) Abort() {
	ctx.handleIdx = len(ctx.handleFunc)
}

// BearerToken 尝试读取 Authorization 头中的 Bearer 的值，读取失败返回空字符串串.
func (ctx *Context) BearerToken() string {
	token := ctx.Request.Header.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, bearerTokenPrefix) {
		return ""
	}
	return token[len(bearerTokenPrefix):]
}

// WriteJSON 设置 statusCode ，Content-Type: json +utf8 ，格式化 value 为 JSON 写到响应 body 中。
func (ctx *Context) WriteJSON(statusCode int, value interface{}) error {
	ctx.ResponseWriter.WriteHeader(statusCode)
	ctx.ResponseWriter.Header().Add(contentType, ContentTypeJSON)
	ctx.ResponseWriter.Header().Add(contentType, ContentTypeUTF8)
	enc := json.NewEncoder(ctx.ResponseWriter)
	return enc.Encode(value)
}

// WriteJSONBytes 设置 statusCode ，Content-Type: json +utf8 ，将 data 写到响应 body 中。
func (ctx *Context) WriteJSONBytes(statusCode int, data []byte) error {
	ctx.ResponseWriter.WriteHeader(statusCode)
	ctx.ResponseWriter.Header().Add(contentType, ContentTypeJSON)
	ctx.ResponseWriter.Header().Add(contentType, ContentTypeUTF8)
	_, err := ctx.ResponseWriter.Write(data)
	return err
}

// WriteHTML 设置 statusCode ，Content-Type: html +utf8 ，将 text 写到响应 body 中。
func (ctx *Context) WriteHTML(statusCode int, text string) error {
	ctx.ResponseWriter.WriteHeader(statusCode)
	ctx.ResponseWriter.Header().Add(contentType, ContentTypeHTML)
	ctx.ResponseWriter.Header().Add(contentType, ContentTypeUTF8)
	_, err := io.WriteString(ctx.ResponseWriter, text)
	return err
}

// WriteHTMLBytes 设置 statusCode ，Content-Type: html +utf8 ，将 data 写到响应 body 中。
func (ctx *Context) WriteHTMLBytes(statusCode int, data []byte) error {
	ctx.ResponseWriter.WriteHeader(statusCode)
	ctx.ResponseWriter.Header().Add(contentType, ContentTypeHTML)
	ctx.ResponseWriter.Header().Add(contentType, ContentTypeUTF8)
	_, err := ctx.ResponseWriter.Write(data)
	return err
}
