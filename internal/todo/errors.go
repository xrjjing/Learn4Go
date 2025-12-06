package todo

import (
	"net/http"
)

// ErrorCode 定义业务错误码
type ErrorCode string

const (
	ErrCodeOK              ErrorCode = "OK"
	ErrCodeBadRequest      ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeConflict        ErrorCode = "CONFLICT"
	ErrCodeTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"
	ErrCodeInternal        ErrorCode = "INTERNAL_ERROR"

	// 业务错误码
	ErrCodeInvalidJSON     ErrorCode = "INVALID_JSON"
	ErrCodeInvalidID       ErrorCode = "INVALID_ID"
	ErrCodeTitleRequired   ErrorCode = "TITLE_REQUIRED"
	ErrCodeEmailRequired   ErrorCode = "EMAIL_REQUIRED"
	ErrCodePasswordWeak    ErrorCode = "PASSWORD_WEAK"
	ErrCodeEmailInvalid    ErrorCode = "EMAIL_INVALID"
	ErrCodeEmailExists     ErrorCode = "EMAIL_EXISTS"
	ErrCodeInvalidCreds    ErrorCode = "INVALID_CREDENTIALS"
	ErrCodeAccountLocked   ErrorCode = "ACCOUNT_LOCKED"
	ErrCodeTokenExpired    ErrorCode = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid    ErrorCode = "TOKEN_INVALID"
	ErrCodeRefreshRequired ErrorCode = "REFRESH_REQUIRED"
)

// APIError 统一 API 错误响应结构
type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Detail  string    `json:"detail,omitempty"`
}

// APIResponse 统一 API 响应结构
type APIResponse struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// errorCodeHTTPStatus 错误码到 HTTP 状态码的映射
var errorCodeHTTPStatus = map[ErrorCode]int{
	ErrCodeOK:              http.StatusOK,
	ErrCodeBadRequest:      http.StatusBadRequest,
	ErrCodeUnauthorized:    http.StatusUnauthorized,
	ErrCodeForbidden:       http.StatusForbidden,
	ErrCodeNotFound:        http.StatusNotFound,
	ErrCodeConflict:        http.StatusConflict,
	ErrCodeTooManyRequests: http.StatusTooManyRequests,
	ErrCodeInternal:        http.StatusInternalServerError,
	ErrCodeInvalidJSON:     http.StatusBadRequest,
	ErrCodeInvalidID:       http.StatusBadRequest,
	ErrCodeTitleRequired:   http.StatusBadRequest,
	ErrCodeEmailRequired:   http.StatusBadRequest,
	ErrCodePasswordWeak:    http.StatusBadRequest,
	ErrCodeEmailInvalid:    http.StatusBadRequest,
	ErrCodeEmailExists:     http.StatusConflict,
	ErrCodeInvalidCreds:    http.StatusUnauthorized,
	ErrCodeAccountLocked:   http.StatusTooManyRequests,
	ErrCodeTokenExpired:    http.StatusUnauthorized,
	ErrCodeTokenInvalid:    http.StatusUnauthorized,
	ErrCodeRefreshRequired: http.StatusUnauthorized,
}

// HTTPStatus 返回错误码对应的 HTTP 状态码
func (c ErrorCode) HTTPStatus() int {
	if status, ok := errorCodeHTTPStatus[c]; ok {
		return status
	}
	return http.StatusInternalServerError
}

// respondAPIError 使用统一错误格式响应
func respondAPIError(w http.ResponseWriter, code ErrorCode, message string) {
	respondJSON(w, APIError{
		Code:    code,
		Message: message,
	}, code.HTTPStatus())
}

// respondAPISuccess 使用统一成功格式响应
func respondAPISuccess(w http.ResponseWriter, data interface{}, status int) {
	respondJSON(w, APIResponse{
		Code: ErrCodeOK,
		Data: data,
	}, status)
}
