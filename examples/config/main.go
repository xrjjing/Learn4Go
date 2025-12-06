// Viper 配置管理示例
// 对应章节: 05_配置与日志.md
//
// 运行方式:
//
//	go run ./examples/config
//
// 环境变量覆盖:
//
//	SERVER_PORT=9090 go run ./examples/config
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置结构体
// mapstructure 标签用于 Viper 绑定
// 类似 Spring 的 @ConfigurationProperties
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Feature  FeatureConfig  `mapstructure:"feature"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	Name   string `mapstructure:"name"`
	User   string `mapstructure:"user"`
}

type FeatureConfig struct {
	EnableCache bool `mapstructure:"enableCache"`
	MaxRetries  int  `mapstructure:"maxRetries"`
}

func main() {
	fmt.Println("=== Viper 配置管理示例 ===")

	// 设置默认值 (类似 Spring 的 @Value 默认值)
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.driver", "sqlite")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.name", "demo")
	viper.SetDefault("database.user", "root")
	viper.SetDefault("feature.enableCache", true)
	viper.SetDefault("feature.maxRetries", 3)

	// 启用环境变量覆盖 (类似 Spring 的环境变量覆盖)
	// SERVER_PORT=9090 会覆盖 server.port
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 绑定到结构体
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("配置绑定失败: %v", err)
	}

	// 打印配置
	fmt.Println("\n--- 服务器配置 ---")
	fmt.Printf("  Host: %s\n", cfg.Server.Host)
	fmt.Printf("  Port: %d\n", cfg.Server.Port)

	fmt.Println("\n--- 数据库配置 ---")
	fmt.Printf("  Driver: %s\n", cfg.Database.Driver)
	fmt.Printf("  Host: %s:%d\n", cfg.Database.Host, cfg.Database.Port)
	fmt.Printf("  Database: %s\n", cfg.Database.Name)
	fmt.Printf("  User: %s\n", cfg.Database.User)

	fmt.Println("\n--- 功能开关 ---")
	fmt.Printf("  EnableCache: %v\n", cfg.Feature.EnableCache)
	fmt.Printf("  MaxRetries: %d\n", cfg.Feature.MaxRetries)

	// 动态获取单个值 (类似 @Value)
	fmt.Println("\n--- 动态获取 ---")
	fmt.Printf("  viper.GetInt(\"server.port\"): %d\n", viper.GetInt("server.port"))
	fmt.Printf("  viper.GetBool(\"feature.enableCache\"): %v\n", viper.GetBool("feature.enableCache"))

	fmt.Println("\n--- 使用提示 ---")
	fmt.Println("  设置环境变量测试: SERVER_PORT=9090 go run ./examples/config")
	fmt.Println("  生产环境通常从 config.yaml 读取配置")
}
