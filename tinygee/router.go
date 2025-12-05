package tinygee

import (
	"net/http"
	"strings"
)

// HandlerFunc 定义业务处理函数。
type HandlerFunc func(*Context)

type router struct {
	handlers map[string]HandlerFunc
	roots    map[string]*node // method -> trie root
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
		roots:    make(map[string]*node),
	}
}

// key 形如 GET-/path
func routeKey(method, pattern string) string {
	return method + "-" + pattern
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := routeKey(method, pattern)

	// 构建 trie
	root, ok := r.roots[method]
	if !ok {
		root = &node{}
		r.roots[method] = root
	}
	root.insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// getRoute 返回匹配到的节点和解析出的参数
func (r *router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for idx, part := range parts {
			if strings.HasPrefix(part, ":") {
				params[part[1:]] = searchParts[idx]
			}
			if strings.HasPrefix(part, "*") && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[idx:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := routeKey(c.Method, n.pattern)
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(ctx *Context) {
			ctx.JSON(http.StatusNotFound, map[string]any{"error": "not found"})
		})
	}
	c.Next()
}

// parsePattern 将路由 pattern 按 / 拆分，保留通配符和参数段。
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0, len(vs))
	for _, item := range vs {
		if item == "" {
			continue
		}
		parts = append(parts, item)
		if item[0] == '*' {
			break
		}
	}
	return parts
}
