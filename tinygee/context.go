package tinygee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Context 封装一次 HTTP 请求的上下文。
// 后续会扩展 Params、中间件等字段。
type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request

	Path   string
	Method string

	StatusCode int
	Params     map[string]string

	handlers []HandlerFunc
	index    int
}

// NewContext 创建上下文对象。
func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// Next 执行下一个中间件/处理器
func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}

// Param 获取路由参数
func (c *Context) Param(key string) string {
	if c.Params == nil {
		return ""
	}
	return c.Params[key]
}

// Status 设置状态码。
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置响应头。
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// String 返回字符串。
func (c *Context) String(code int, format string, values ...any) {
	c.SetHeader("Content-Type", "text/plain; charset=utf-8")
	c.Status(code)
	_, _ = fmt.Fprintf(c.Writer, format, values...)
}

// JSON 返回 JSON。
func (c *Context) JSON(code int, obj any) {
	c.SetHeader("Content-Type", "application/json; charset=utf-8")
	c.Status(code)
	_ = json.NewEncoder(c.Writer).Encode(obj)
}

// Data 返回字节数据。
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	_, _ = c.Writer.Write(data)
}

// HTML 简易返回 HTML。
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html; charset=utf-8")
	c.Status(code)
	_, _ = c.Writer.Write([]byte(html))
}
