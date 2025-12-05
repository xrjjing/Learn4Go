package todo

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// Role 角色定义
type Role string

const (
	RoleAdmin Role = "admin" // 管理员 - 完全权限
	RoleUser  Role = "user"  // 普通用户 - 只能操作自己的资源
	RoleGuest Role = "guest" // 访客 - 只读权限
)

// Action 操作类型
type Action string

const (
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

// Resource 资源类型
type Resource string

const (
	ResourceTodos Resource = "todos"
)

// Permission 权限定义
type Permission struct {
	Resource Resource
	Action   Action
}

// RBACManager RBAC管理器
type RBACManager struct {
	permissions map[Role][]Permission
}

// NewRBACManager 创建RBAC管理器
func NewRBACManager() *RBACManager {
	return &RBACManager{
		permissions: map[Role][]Permission{
			RoleAdmin: {
				{ResourceTodos, ActionCreate},
				{ResourceTodos, ActionRead},
				{ResourceTodos, ActionUpdate},
				{ResourceTodos, ActionDelete},
			},
			RoleUser: {
				{ResourceTodos, ActionCreate},
				{ResourceTodos, ActionRead},
				{ResourceTodos, ActionUpdate},
				{ResourceTodos, ActionDelete},
			},
			RoleGuest: {
				{ResourceTodos, ActionRead},
			},
		},
	}
}

// CheckPermission 检查角色是否有权限执行操作
func (m *RBACManager) CheckPermission(role Role, resource Resource, action Action) bool {
	permissions, ok := m.permissions[role]
	if !ok {
		return false
	}

	for _, perm := range permissions {
		if perm.Resource == resource && perm.Action == action {
			return true
		}
	}
	return false
}

// GetResourceFromRequest 从请求中解析资源
func GetResourceFromRequest(r *http.Request) Resource {
	path := strings.Trim(r.URL.Path, "/")
	if strings.Contains(path, "todos") {
		return ResourceTodos
	}
	return ""
}

// GetActionFromRequest 从HTTP方法解析操作
func GetActionFromRequest(r *http.Request) Action {
	switch r.Method {
	case http.MethodPost:
		return ActionCreate
	case http.MethodGet:
		return ActionRead
	case http.MethodPut, http.MethodPatch:
		return ActionUpdate
	case http.MethodDelete:
		return ActionDelete
	default:
		return ""
	}
}

// GetTodoIDFromRequest 从请求中提取TODO ID
func GetTodoIDFromRequest(r *http.Request) (int, error) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	// 查找todos后面的ID
	for i, part := range parts {
		if part == "todos" && i+1 < len(parts) {
			return strconv.Atoi(parts[i+1])
		}
	}

	return 0, errors.New("no todo ID in request")
}

// authzMiddleware RBAC授权中间件
func (s *Server) authzMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 获取用户信息
		userID, ok := GetUserID(r.Context())
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		// 获取用户角色
		user, err := s.userStore.FindByID(r.Context(), userID)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "user not found")
			return
		}

		// 解析资源和操作
		resource := GetResourceFromRequest(r)
		action := GetActionFromRequest(r)

		// 检查基本权限
		if !s.rbacManager.CheckPermission(user.Role, resource, action) {
			respondError(w, http.StatusForbidden, "insufficient permissions")
			return
		}

		// 对于普通用户，需要验证资源所有权
		if user.Role == RoleUser && (action == ActionRead || action == ActionUpdate || action == ActionDelete) {
			todoID, err := GetTodoIDFromRequest(r)
			if err == nil {
				// 验证TODO所有权
				todo, exists, err := s.store.Get(todoID)
				if err != nil {
					respondError(w, http.StatusInternalServerError, "internal error")
					return
				}
				if !exists {
					respondError(w, http.StatusNotFound, "not found")
					return
				}
				if todo.UserID != userID {
					respondError(w, http.StatusForbidden, "you don't own this resource")
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
