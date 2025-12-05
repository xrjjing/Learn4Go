-- RBAC权限控制系统数据库表
-- 创建时间: 2025-12-05
-- 说明: 实现基于角色的访问控制（Role-Based Access Control）

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE COMMENT '角色名称: admin, user, guest',
    description VARCHAR(255) COMMENT '角色描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 用户-角色关联表 (多对多)
CREATE TABLE IF NOT EXISTS user_roles (
    user_id INT NOT NULL COMMENT '用户ID',
    role_id INT NOT NULL COMMENT '角色ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id),
    KEY idx_user_id (user_id),
    KEY idx_role_id (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- 预填充角色数据
INSERT INTO roles (name, description) VALUES
    ('admin', '管理员 - 拥有所有权限'),
    ('user', '普通用户 - 只能操作自己的资源'),
    ('guest', '访客 - 只读权限')
ON DUPLICATE KEY UPDATE description=VALUES(description);

-- 权限表（可选，用于未来扩展）
CREATE TABLE IF NOT EXISTS permissions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    resource VARCHAR(100) NOT NULL COMMENT '资源名称: todos, users',
    action VARCHAR(100) NOT NULL COMMENT '操作: create, read, update, delete',
    description VARCHAR(255) COMMENT '权限描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_resource_action (resource, action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';

-- 角色-权限关联表（可选，用于未来扩展）
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INT NOT NULL COMMENT '角色ID',
    permission_id INT NOT NULL COMMENT '权限ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, permission_id),
    KEY idx_role_id (role_id),
    KEY idx_permission_id (permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

-- 预填充权限数据（可选）
INSERT INTO permissions (resource, action, description) VALUES
    ('todos', 'create', '创建TODO'),
    ('todos', 'read', '读取TODO'),
    ('todos', 'update', '更新TODO'),
    ('todos', 'delete', '删除TODO')
ON DUPLICATE KEY UPDATE description=VALUES(description);

-- 预填充角色权限关联（可选）
-- Admin拥有所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
ON DUPLICATE KEY UPDATE role_id=VALUES(role_id);

-- User拥有CRUD权限（但需要所有权验证）
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'user' AND p.resource = 'todos'
ON DUPLICATE KEY UPDATE role_id=VALUES(role_id);

-- Guest只有读权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'guest' AND p.resource = 'todos' AND p.action = 'read'
ON DUPLICATE KEY UPDATE role_id=VALUES(role_id);

-- 为todos表添加user_id字段（如果不存在）
-- ALTER TABLE todos ADD COLUMN IF NOT EXISTS user_id INT NOT NULL DEFAULT 1 COMMENT '所属用户ID';
-- ALTER TABLE todos ADD INDEX IF NOT EXISTS idx_user_id (user_id);
