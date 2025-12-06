package todo

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig    `mapstructure:"server"`
	Database DatabaseConfig  `mapstructure:"database"`
	JWT      JWTConfigParams `mapstructure:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Addr         string        `mapstructure:"addr"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	Path            string        `mapstructure:"path"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// JWTConfigParams JWT 配置参数
type JWTConfigParams struct {
	Secret     string        `mapstructure:"secret"`
	TTL        time.Duration `mapstructure:"ttl"`
	RefreshTTL time.Duration `mapstructure:"refresh_ttl"`
}

// LoadConfig 从环境变量和配置文件加载配置
func LoadConfig() (*Config, error) {
	v := viper.New()

	// 设置默认值
	v.SetDefault("server.addr", ":8080")
	v.SetDefault("server.read_timeout", 15*time.Second)
	v.SetDefault("server.write_timeout", 15*time.Second)
	v.SetDefault("server.idle_timeout", 60*time.Second)

	v.SetDefault("database.driver", "memory")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.user", "root")
	v.SetDefault("database.password", "root")
	v.SetDefault("database.name", "learn4go")
	v.SetDefault("database.path", "todos.db")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", 5*time.Minute)

	v.SetDefault("jwt.ttl", 24*time.Hour)
	v.SetDefault("jwt.refresh_ttl", 7*24*time.Hour)

	// 环境变量映射
	v.SetEnvPrefix("TODO")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 手动绑定关键环境变量（兼容旧变量名）
	_ = v.BindEnv("database.driver", "TODO_STORAGE")
	_ = v.BindEnv("server.addr", "TODO_ADDR")
	_ = v.BindEnv("database.host", "TODO_DB_HOST")
	_ = v.BindEnv("database.port", "TODO_DB_PORT")
	_ = v.BindEnv("database.user", "TODO_DB_USER")
	_ = v.BindEnv("database.password", "TODO_DB_PASS")
	_ = v.BindEnv("database.name", "TODO_DB_NAME")
	_ = v.BindEnv("database.path", "TODO_DB_PATH")
	_ = v.BindEnv("jwt.secret", "JWT_SECRET")

	// 尝试加载配置文件（可选）
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	_ = v.ReadInConfig() // 忽略文件不存在的错误

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
