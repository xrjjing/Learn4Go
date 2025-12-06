// GORM 数据库访问示例
// 对应章节: 06_数据库访问.md
//
// 运行方式:
//
//	go run ./examples/database
//
// 本示例使用内存 SQLite，无需额外配置数据库
package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User 用户模型 (类似 JPA @Entity)
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:64;uniqueIndex"`
	Email     string    `gorm:"size:128"`
	Age       int       `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// Todo 待办模型
type Todo struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	Title     string    `gorm:"size:256;not null"`
	Done      bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User User `gorm:"foreignKey:UserID"` // 关联
}

func main() {
	fmt.Println("=== GORM 数据库访问示例 ===")

	// 1. 连接数据库 (内存 SQLite)
	fmt.Println("\n--- 连接数据库 ---")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	fmt.Println("  SQLite 内存数据库连接成功")

	// 2. 自动迁移 (类似 Hibernate ddl-auto=update)
	fmt.Println("\n--- 自动迁移 ---")
	if err := db.AutoMigrate(&User{}, &Todo{}); err != nil {
		log.Fatalf("迁移失败: %v", err)
	}
	fmt.Println("  表结构创建完成: users, todos")

	// 3. 创建 (Create)
	fmt.Println("\n--- 创建记录 ---")
	users := []User{
		{Name: "张三", Email: "zhang@example.com", Age: 25},
		{Name: "李四", Email: "li@example.com", Age: 30},
		{Name: "王五", Email: "wang@example.com", Age: 28},
	}
	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			log.Printf("创建用户失败: %v", err)
			continue
		}
		fmt.Printf("  创建用户: ID=%d, Name=%s\n", users[i].ID, users[i].Name)
	}

	// 创建关联的待办
	todos := []Todo{
		{UserID: 1, Title: "学习 Go 基础语法"},
		{UserID: 1, Title: "学习 Gin 框架"},
		{UserID: 2, Title: "复习 GORM"},
	}
	db.Create(&todos)
	fmt.Printf("  创建待办: %d 条\n", len(todos))

	// 4. 查询 (Read)
	fmt.Println("\n--- 查询记录 ---")

	// 单条查询 (类似 JPA findById)
	var user User
	db.First(&user, 1)
	fmt.Printf("  First(id=1): %s <%s>\n", user.Name, user.Email)

	// 条件查询 (类似 JPA findByName)
	var found User
	db.Where("name = ?", "李四").First(&found)
	fmt.Printf("  Where(name=李四): ID=%d, Age=%d\n", found.ID, found.Age)

	// 查询所有 (类似 JPA findAll)
	var allUsers []User
	db.Find(&allUsers)
	fmt.Printf("  Find all: %d 条记录\n", len(allUsers))

	// 5. 更新 (Update)
	fmt.Println("\n--- 更新记录 ---")
	db.Model(&user).Update("Age", 26)
	fmt.Printf("  更新 %s 年龄为 %d\n", user.Name, 26)

	// 部分更新
	db.Model(&User{}).Where("id = ?", 2).Updates(map[string]interface{}{
		"Age":   31,
		"Email": "lisi_new@example.com",
	})
	fmt.Println("  批量更新 ID=2 的多个字段")

	// 6. 事务 (类似 @Transactional)
	fmt.Println("\n--- 事务示例 ---")
	err = db.Transaction(func(tx *gorm.DB) error {
		// 事务内的操作
		if err := tx.Create(&User{Name: "事务用户", Email: "tx@example.com"}).Error; err != nil {
			return err // 返回错误将回滚
		}
		if err := tx.Create(&Todo{UserID: 4, Title: "事务待办"}).Error; err != nil {
			return err
		}
		return nil // 返回 nil 自动提交
	})
	if err != nil {
		fmt.Printf("  事务失败: %v\n", err)
	} else {
		fmt.Println("  事务提交成功")
	}

	// 7. 关联查询 (Preload)
	fmt.Println("\n--- 关联查询 ---")
	var userWithTodos User
	db.Preload("Todos").First(&userWithTodos, 1)
	fmt.Printf("  用户 %s 的待办:\n", userWithTodos.Name)

	var userTodos []Todo
	db.Where("user_id = ?", 1).Find(&userTodos)
	for _, t := range userTodos {
		status := "[ ]"
		if t.Done {
			status = "[x]"
		}
		fmt.Printf("    %s %s\n", status, t.Title)
	}

	// 8. 删除 (Delete)
	fmt.Println("\n--- 删除记录 ---")
	db.Delete(&User{}, 3)
	fmt.Println("  删除 ID=3 的用户")

	// 最终统计
	var count int64
	db.Model(&User{}).Count(&count)
	fmt.Printf("\n--- 最终统计 ---\n")
	fmt.Printf("  剩余用户: %d\n", count)
	db.Model(&Todo{}).Count(&count)
	fmt.Printf("  剩余待办: %d\n", count)

	fmt.Println("\n--- Java 对照 ---")
	fmt.Println("  db.Create()      → EntityManager.persist()")
	fmt.Println("  db.First()       → EntityManager.find()")
	fmt.Println("  db.Save()        → EntityManager.merge()")
	fmt.Println("  db.Delete()      → EntityManager.remove()")
	fmt.Println("  db.Transaction() → @Transactional")
}
