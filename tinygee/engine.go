package tinygee

import "net/http"

// Engine 实现最小的 HTTP 路由引擎。
type Engine struct {
	router *router
	groups []*RouterGroup
	// 全局中间件
	middlewares []HandlerFunc
}

// New 创建引擎。
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.groups = []*RouterGroup{{engine: engine}}
	return engine
}

// addRoute 注册路由。
func (e *Engine) addRoute(method, pattern string, handler HandlerFunc) {
	e.router.addRoute(method, pattern, handler)
}

// GET 注册 GET 路由。
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute(http.MethodGet, pattern, handler)
}

// POST 注册 POST 路由。
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute(http.MethodPost, pattern, handler)
}

// Use 注册全局中间件。
func (e *Engine) Use(m ...HandlerFunc) {
	e.middlewares = append(e.middlewares, m...)
}

// ServeHTTP 实现 http.Handler。
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	// 选择命中的分组中间件 + 全局中间件
	for _, group := range e.groups {
		if len(group.prefix) == 0 || hasPrefix(req.URL.Path, group.prefix) {
			c.handlers = append(c.handlers, group.middlewares...)
		}
	}
	c.handlers = append(c.handlers, e.middlewares...)
	e.router.handle(c)
}

// Run 启动 HTTP 服务。
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

// hasPrefix 判断路由是否以分组前缀开头（保证 / 分隔）
func hasPrefix(path, prefix string) bool {
	if len(prefix) == 0 {
		return true
	}
	if len(path) < len(prefix) {
		return false
	}
	if path[:len(prefix)] == prefix {
		if len(path) == len(prefix) {
			return true
		}
		return path[len(prefix)] == '/'
	}
	return false
}

// MatchPrefix 暴露给中间件使用，判断路径是否带有特定前缀（以 / 边界为准）。
func MatchPrefix(path, prefix string) bool {
	return hasPrefix(path, prefix)
}
