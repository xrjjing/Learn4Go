package tinygee

import "net/http"

// RouterGroup 支持路由分组与分组中间件。
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	engine      *Engine
}

// Group 创建子分组。
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		engine: g.engine,
	}
	g.engine.groups = append(g.engine.groups, newGroup)
	return newGroup
}

// Use 为分组注册中间件。
func (g *RouterGroup) Use(m ...HandlerFunc) {
	g.middlewares = append(g.middlewares, m...)
}

// addRoute 带分组前缀的路由注册。
func (g *RouterGroup) addRoute(method, comp string, handler HandlerFunc) {
	pattern := g.prefix + comp
	g.engine.router.addRoute(method, pattern, handler)
}

func (g *RouterGroup) GET(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodGet, pattern, handler)
}

func (g *RouterGroup) POST(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodPost, pattern, handler)
}

// Engine 的 Group 代理方法
func (e *Engine) Group(prefix string) *RouterGroup {
	return e.groups[0].Group(prefix)
}
