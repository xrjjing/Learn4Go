/**
 * Learn4Go 通用前端工具层。
 *
 * 文件职责：
 * - 给 admin.html、log-detective.html 这类页面提供公共 API 调用、调试日志、表格渲染和 Toast。
 * - 这里不关心某个具体业务页面，只提供“页面脚本可复用的底层能力”。
 *
 * 真实调用链：
 * - 页面按钮点击 -> 页面内联脚本调用 apiClient / renderGenericTable / toggleDebugPanel
 * - apiClient 内部负责脱敏日志 + fetch + 错误抛出
 * - 页面再决定如何把结果渲染进表格、日志区或提示区。
 *
 * 排查建议：
 * - 如果 admin.html 的调试面板没有日志，先看 logApiActivity() 是否拿到 #log-container。
 * - 如果 log-detective.html 请求失败但页面没有清晰报错，先看 apiClient() 的异常抛出链。
 */

// ============================================
// API 客户端（带脱敏日志）
// ============================================

/**
 * 脱敏处理：移除敏感字段
 * @param {string|FormData} body - 请求体
 * @returns {string} 脱敏后的字符串
 */
// sanitizeBody 只服务调试日志，不改变真实请求体。
// 作用是避免密码、token 等敏感字段直接被写进页面调试面板。
function sanitizeBody(body) {
    if (!body) return '';

    // 如果是 FormData，转换为对象
    if (body instanceof FormData) {
        const obj = {};
        for (const [key, value] of body.entries()) {
            obj[key] = value;
        }
        body = obj;
    }

    // 如果是字符串，尝试解析为对象
    if (typeof body === 'string') {
        try {
            body = JSON.parse(body);
        } catch {
            return '[Raw String Body]';
        }
    }

    // 脱敏敏感字段
    const sensitiveFields = ['password', 'token', 'secret', 'api_key', 'access_token'];
    const sanitized = { ...body };

    sensitiveFields.forEach(field => {
        if (field in sanitized) {
            sanitized[field] = '***REDACTED***';
        }
    });

    return JSON.stringify(sanitized);
}

/**
 * 记录 API 活动到调试面板
 * @param {string} type - 日志类型 (REQUEST/RESPONSE/ERROR)
 * @param {string|number} method - HTTP 方法或状态码
 * @param {string} url - 请求 URL
 * @param {any} data - 数据内容
 * @param {number} duration - 请求耗时（毫秒）
 */
// logApiActivity 把请求/响应/错误投递到页面右下角或侧边的调试面板。
// 依赖页面提前准备好 #log-container；没有容器时会直接静默跳过。
// 调试面板写入：admin.html 底部的 #log-container 就是靠这里持续追加请求/响应日志。
function logApiActivity(type, method, url, data, duration = 0) {
    const logContainer = document.getElementById('log-container');
    if (!logContainer) return; // 如果没有日志面板，跳过

    const timestamp = new Date().toLocaleTimeString('zh-CN', { hour12: false });
    const logEntry = document.createElement('div');
    logEntry.className = 'log-entry';

    let emoji = '';
    let color = '';
    let message = '';

    switch (type) {
        case 'REQUEST':
            emoji = '🟢';
            color = '#10b981';
            message = `[${method}] ${url}`;
            if (data) message += ` | Body: ${data}`;
            break;
        case 'RESPONSE':
            emoji = '🔵';
            color = '#3b82f6';
            message = `[${method}] ${duration}ms | URL: ${url}`;
            break;
        case 'ERROR':
            emoji = '🔴';
            color = '#ef4444';
            message = `[ERROR] ${url} | ${data}`;
            break;
    }

    // FIX: Use DOM API instead of innerHTML to prevent XSS
    const iconSpan = document.createElement('span');
    iconSpan.style.color = color;
    iconSpan.textContent = emoji;

    const timeSpan = document.createElement('span');
    timeSpan.style.color = 'var(--text-muted)';
    timeSpan.style.marginRight = '8px';
    timeSpan.textContent = timestamp;

    const msgSpan = document.createElement('span');
    msgSpan.textContent = message;

    logEntry.appendChild(iconSpan);
    logEntry.appendChild(timeSpan);
    logEntry.appendChild(msgSpan);

    logContainer.appendChild(logEntry);
    logContainer.scrollTop = logContainer.scrollHeight; // 自动滚动到底部
}

/**
 * 统一的 API 客户端
 * @param {string} url - 请求 URL
 * @param {object} options - fetch 选项
 * @returns {Promise<any>} 响应数据
 */
// apiClient 是 common.js 最核心的统一请求入口。
// 它不处理 token 自动刷新；若页面需要 refresh 机制，应走 auth-helper.js 的 AuthClient。
// 统一请求入口：负责 fetch、日志记录、错误抛出和 JSON 解析。
async function apiClient(url, options = {}) {
    const method = options.method || 'GET';
    const startTime = Date.now();

    // 脱敏处理：不记录完整敏感信息
    const safeBody = sanitizeBody(options.body);
    logApiActivity('REQUEST', method, url, safeBody);

    try {
        const response = await fetch(url, options);
        const data = await response.json();
        const duration = Date.now() - startTime;

        logApiActivity('RESPONSE', response.status, url, data, duration);

        if (!response.ok) {
            throw new Error(data.detail || data.message || `HTTP ${response.status}`);
        }

        return data;
    } catch (error) {
        logApiActivity('ERROR', '---', url, error.message);
        throw error;
    }
}

// ============================================
// 通用表格渲染器
// ============================================

/**
 * 安全的通用表格渲染器
 * @param {string} elementId - tbody 元素的 ID
 * @param {Array} data - 数据数组
 * @param {Array} columns - 列定义数组
 * @example
 * renderGenericTable('users-table', users, [
 *   { key: 'id' },
 *   { key: 'username' },
 *   { key: 'roles', render: (roles) => roles.map(r => r.name).join(', ') }
 * ]);
 */
// renderGenericTable 负责把“数组数据”安全灌进指定 tbody。
// admin.html 的用户/角色/权限列表，后续都可复用这一层。
// 通用表格渲染：把后端返回数组映射成 tbody 内容，admin.html 的多个列表都复用它。
function renderGenericTable(elementId, data, columns) {
    const tbody = document.getElementById(elementId);

    // 边界检查
    if (!tbody) {
        console.error(`Table element #${elementId} not found`);
        return;
    }

    // 安全清空表格
    tbody.replaceChildren();

    // 处理空数据
    if (!data || data.length === 0) {
        const tr = document.createElement('tr');
        const td = document.createElement('td');
        td.colSpan = columns.length;
        td.textContent = 'No data available';
        td.style.textAlign = 'center';
        td.style.color = 'var(--text-muted)';
        tr.appendChild(td);
        tbody.appendChild(tr);
        return;
    }

    // 渲染数据行
    data.forEach(item => {
        const tr = document.createElement('tr');

        columns.forEach(col => {
            const td = document.createElement('td');

            if (col.render) {
                // 使用自定义渲染函数
                const content = col.render(item[col.key], item);

                if (content instanceof HTMLElement) {
                    // 安全插入 DOM 节点
                    td.appendChild(content);
                } else if (content != null) {
                    // 安全设置文本（转换为字符串）
                    td.textContent = String(content);
                } else {
                    td.textContent = '-';
                }
            } else {
                // 默认文本渲染（使用 nullish 合并处理 0/false）
                const value = item[col.key];
                td.textContent = value ?? '-';
            }

            tr.appendChild(td);
        });

        tbody.appendChild(tr);
    });
}

// ============================================
// 调试面板管理
// ============================================

/**
 * 切换调试面板显示状态
 * @param {boolean} show - 是否显示
 */
// 调试面板开关，常见于 admin.html 这种需要边查接口边看日志的页面。
// 调试面板显隐：与 admin.html 右上角的“调试模式”开关联动。
function toggleDebugPanel(show) {
    const debugPanel = document.getElementById('debug-panel');
    if (debugPanel) {
        debugPanel.style.display = show ? 'block' : 'none';
    }
}

/**
 * 清空调试日志
 */
// 清空调试面板现有输出，不会影响真实业务状态。
// 调试日志清理：只清空前端面板内容，不影响真实网络请求。
function clearDebugLogs() {
    const logContainer = document.getElementById('log-container');
    if (logContainer) {
        logContainer.replaceChildren();
    }
}

// ============================================
// 工具函数
// ============================================

/**
 * 显示 Toast 提示
 * @param {string} message - 提示消息
 * @param {string} type - 类型 (success/error/info)
 */
// 轻量 Toast；当页面没有引入更完整的 Toast 体系时，可直接用它给出结果反馈。
// 轻量提示：给没有引入完整组件体系的页面复用一个最小消息提示。
function showToast(message, type = 'info') {
    // 简单实现，可以后续增强为更美观的 Toast 组件
    const colors = {
        success: '#10b981',
        error: '#ef4444',
        info: '#3b82f6'
    };

    const toast = document.createElement('div');
    toast.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        background: ${colors[type] || colors.info};
        color: white;
        padding: 12px 20px;
        border-radius: 8px;
        box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        z-index: 10000;
        animation: slideIn 0.3s ease;
    `;
    toast.textContent = message;

    document.body.appendChild(toast);

    setTimeout(() => {
        toast.style.animation = 'slideOut 0.3s ease';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}

// ============================================
// 导出（如果使用模块化）
// ============================================

// 如果使用 ES6 模块，可以取消注释
// export { apiClient, renderGenericTable, toggleDebugPanel, clearDebugLogs, showToast };