package todo

import (
	"github.com/go-playground/validator/v10"
)

// 全局验证器实例
var validate = validator.New()

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// CreateTodoRequest 创建 TODO 请求
type CreateTodoRequest struct {
	Title string `json:"title" validate:"required,min=1,max=256"`
}

// UpdateTodoRequest 更新 TODO 请求
type UpdateTodoRequest struct {
	Done bool `json:"done"`
}

// RefreshRequest 刷新令牌请求
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// FormatValidationError 格式化验证错误为用户友好的消息
func FormatValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				return e.Field() + " is required"
			case "email":
				return "invalid email format"
			case "min":
				return e.Field() + " must be at least " + e.Param() + " characters"
			case "max":
				return e.Field() + " must be at most " + e.Param() + " characters"
			}
		}
	}
	return "validation failed"
}
