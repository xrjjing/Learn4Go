package todo

// 本文件集中定义请求 DTO 和基础校验规则。
//
// 当前 handler.go 里仍有部分手写 JSON / 非空校验，这里更像“统一约束入口”和后续收口点。
// 如果后续要把校验逻辑标准化，通常会先扩展这里。
import (
	"github.com/go-playground/validator/v10"
)

// validate 是包级共享验证器，避免每次请求都重复创建 validator 实例。
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

// ValidateStruct 是 handler 层可复用的统一入口。
// ValidateStruct：handler 收到请求体后最先经过的结构校验入口。
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// FormatValidationError 把 validator 的底层错误翻译成更适合前端展示的消息。
// FormatValidationError：把 validator 的底层错误转成更友好的可返回消息。
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
