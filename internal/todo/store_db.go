package todo

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TodoModel GORM 模型
type TodoModel struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Title     string    `gorm:"size:256;not null"`
	Done      bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (TodoModel) TableName() string {
	return "todos"
}

// DBStore 基于 GORM 的数据库存储
type DBStore struct {
	db *gorm.DB
}

// DBConfig 数据库配置
type DBConfig struct {
	Driver   string // sqlite, mysql
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SQLite   string // SQLite 文件路径
}

// NewDBStore 创建数据库存储
func NewDBStore(cfg DBConfig) (*DBStore, error) {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
		dialector = mysql.Open(dsn)
	case "sqlite":
		path := cfg.SQLite
		if path == "" {
			path = ":memory:"
		}
		dialector = sqlite.Open(path)
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", cfg.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // Warn 级别便于排障
	})
	if err != nil {
		return nil, err
	}

	// 自动迁移
	if err := db.AutoMigrate(&TodoModel{}); err != nil {
		return nil, err
	}

	return &DBStore{db: db}, nil
}

// NewSQLiteStore 快捷创建 SQLite 存储
func NewSQLiteStore(dbPath string) (*DBStore, error) {
	return NewDBStore(DBConfig{Driver: "sqlite", SQLite: dbPath})
}

// NewMySQLStore 快捷创建 MySQL 存储
func NewMySQLStore(host string, port int, user, password, dbName string) (*DBStore, error) {
	return NewDBStore(DBConfig{
		Driver:   "mysql",
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
	})
}

// Close 关闭数据库连接池
func (s *DBStore) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// List 返回全部待办
func (s *DBStore) List() ([]Todo, error) {
	var models []TodoModel
	if err := s.db.Order("created_at DESC").Find(&models).Error; err != nil {
		log.Printf("[DBStore] List 查询失败: %v", err)
		return nil, err
	}

	todos := make([]Todo, len(models))
	for i, m := range models {
		todos[i] = modelToTodo(m)
	}
	return todos, nil
}

// ListByUser 返回指定用户的待办列表
func (s *DBStore) ListByUser(userID uint) ([]Todo, error) {
	var models []TodoModel
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&models).Error; err != nil {
		log.Printf("[DBStore] ListByUser 查询失败: %v", err)
		return nil, err
	}

	todos := make([]Todo, len(models))
	for i, m := range models {
		todos[i] = modelToTodo(m)
	}
	return todos, nil
}

// Create 新建待办
func (s *DBStore) Create(title string, userID uint) (Todo, error) {
	model := TodoModel{
		UserID: userID,
		Title:  title,
	}
	if err := s.db.Create(&model).Error; err != nil {
		log.Printf("[DBStore] Create 失败: %v", err)
		return Todo{}, err
	}
	return modelToTodo(model), nil
}

// Get 获取指定ID的待办
func (s *DBStore) Get(id int) (Todo, bool, error) {
	var model TodoModel
	if err := s.db.First(&model, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Todo{}, false, nil
		}
		log.Printf("[DBStore] Get 查询失败: %v", err)
		return Todo{}, false, err
	}
	return modelToTodo(model), true, nil
}

// Toggle 设置完成状态
func (s *DBStore) Toggle(id int, done bool) (Todo, bool, error) {
	var model TodoModel
	if err := s.db.First(&model, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Todo{}, false, nil
		}
		return Todo{}, false, err
	}

	if err := s.db.Model(&model).Update("done", done).Error; err != nil {
		log.Printf("[DBStore] Toggle 更新失败: %v", err)
		return Todo{}, false, err
	}
	model.Done = done
	return modelToTodo(model), true, nil
}

// Delete 删除待办
func (s *DBStore) Delete(id int) (bool, error) {
	result := s.db.Delete(&TodoModel{}, uint(id))
	if result.Error != nil {
		log.Printf("[DBStore] Delete 失败: %v", result.Error)
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

// modelToTodo 转换模型到 API 结构
func modelToTodo(m TodoModel) Todo {
	return Todo{
		ID:        int(m.ID),
		UserID:    m.UserID,
		Title:     m.Title,
		Done:      m.Done,
		CreatedAt: m.CreatedAt,
	}
}
