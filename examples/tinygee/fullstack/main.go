package main

import (
	"expvar"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/xrjjing/Learn4Go/tinygee"
	"github.com/xrjjing/Learn4Go/tinygee/middleware"
	"github.com/xrjjing/Learn4Go/tinygee/middleware/auth"
	"github.com/xrjjing/Learn4Go/tinygee/middleware/ratelimit"
	"github.com/xrjjing/Learn4Go/tinygee/render"
)

// Fullstack 示例：路由 + 中间件 + 模板 + JWT/RBAC + 限流 + /metrics
func main() {
	secret := "tinygee-fullstack"
	jwtCfg := auth.JWTConfig{Secret: secret, TTL: time.Hour}
	token, err := auth.GenerateToken(jwtCfg, 1, "admin")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("示例 token: %s", token)

	r := tinygee.New()
	r.Use(middleware.Logger(), middleware.Recover())

	// 限流（每秒 5 次）
	r.Use(ratelimit.New(time.Second, 5).Middleware())

	// 模板与静态
	funcMap := template.FuncMap{"upper": func(s string) string { return template.HTMLEscapeString(s) }}
	renderer, err := render.New("examples/tinygee/day6/templates/*.html", funcMap)
	if err != nil {
		log.Fatal(err)
	}
	r.GET("/", func(c *tinygee.Context) {
		renderer.HTML(c, http.StatusOK, "hello.html", map[string]any{"Name": "TinyGee Fullstack"})
	})
	r.GET("/static/*filepath", render.Static("/static/", http.Dir("examples/tinygee/day6/static")))

	// 安全分组
	api := r.Group("/api")
	api.Use(auth.NewJWTMiddleware(jwtCfg))
	api.Use(auth.RBAC(auth.RBACConfig{
		RolePermissions: map[string][]string{
			"admin": {"/api"},
			"user":  {"/api/public"},
		},
	}))
	api.GET("/secure", func(c *tinygee.Context) {
		c.JSON(http.StatusOK, map[string]string{"message": "secure ok"})
	})

	// metrics
	r.GET("/metrics", func(c *tinygee.Context) {
		expvar.Handler().ServeHTTP(c.Writer, c.Req)
	})

	log.Println("tinygee fullstack on :8092")
	log.Fatal(r.Run(":8092"))
}
