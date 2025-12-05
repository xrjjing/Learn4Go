package auth

import (
	"net/http"

	"github.com/xrjjing/Learn4Go/tinygee"
)

// RBACConfig 定义角色到可访问路由的映射（简单版：基于前缀）。
type RBACConfig struct {
	RolePermissions map[string][]string // role -> allowed prefixes
}

// RBAC 返回授权中间件。
func RBAC(cfg RBACConfig) tinygee.HandlerFunc {
	return func(c *tinygee.Context) {
		role := c.Params["role"]
		if role == "" {
			c.JSON(http.StatusForbidden, map[string]string{"error": "role required"})
			return
		}
		prefixes := cfg.RolePermissions[role]
		for _, p := range prefixes {
			if tinygee.MatchPrefix(c.Path, p) {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, map[string]string{"error": "insufficient permissions"})
	}
}
