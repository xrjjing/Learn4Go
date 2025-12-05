package main

import (
	"log"
	"net/http"
	"time"

	"github.com/xrjjing/Learn4Go/tinygee"
	"github.com/xrjjing/Learn4Go/tinygee/middleware"
	"github.com/xrjjing/Learn4Go/tinygee/middleware/auth"
)

func main() {
	secret := "tinygee-secret"
	jwtCfg := auth.JWTConfig{Secret: secret, TTL: time.Hour}
	token, err := auth.GenerateToken(jwtCfg, 1, "admin")
	if err != nil {
		log.Fatal(err)
	}

	app := tinygee.New()
	app.Use(middleware.Logger(), middleware.Recover())

	api := app.Group("/api")
	api.Use(auth.NewJWTMiddleware(jwtCfg))
	api.Use(auth.RBAC(auth.RBACConfig{
		RolePermissions: map[string][]string{
			"admin": {"/api"},
			"user":  {"/api/public"},
		},
	}))

	api.GET("/ping", func(c *tinygee.Context) {
		c.JSON(http.StatusOK, map[string]string{"message": "pong"})
	})

	log.Printf("use token: %s", token)
	log.Println("tinygee auth demo on :8089")
	log.Fatal(app.Run(":8089"))
}
