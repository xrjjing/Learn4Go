/**
 * 前端认证辅助层。
 *
 * 服务页面：
 * - 直接服务 web/todo-login.html。
 * - 也适合作为后续真实后端页面的统一鉴权薄封装。
 *
 * 真实调用链：
 * 页面点击登录 -> login() 调 /v1/login -> AuthClient.setTokens() 落本地存储
 * -> 后续页面调用 AuthClient.authFetch() -> 自动带 Bearer Token
 * -> 若收到 401，则这里内部调用 /v1/refresh 刷新 access token 并重试一次。
 *
 * 角色定位：
 * - 它不是完整状态管理器，只是一个“本地 token 仓库 + 自动 refresh 适配层”。
 * - 若登录后请求还是 401，优先看 refreshToken() 和 authFetch()。
 */
// 认证辅助：存取 token + 自动刷新
(function () {
    // 所有 access/refresh token 都集中保存在这一个 key 下，便于定位和清理。
    const STORAGE_KEY = 'learn4go_auth';
    const apiBase = (window.AppConfig && window.AppConfig.todoApiBaseUrl) || '';

    // 读取本地 token；任何 JSON 解析失败都降级为空对象，避免页面直接崩。
    function load() {
        try {
            const raw = localStorage.getItem(STORAGE_KEY);
            if (!raw) return { access: '', refresh: '' };
            return JSON.parse(raw);
        } catch (e) {
            return { access: '', refresh: '' };
        }
    }

    // 写回 token；这里只负责存储，不负责校验 token 是否有效。
    function save(tokens) {
        try {
            localStorage.setItem(STORAGE_KEY, JSON.stringify(tokens));
        } catch (e) {
            /* 本地存储失败时静默忽略 */
        }
    }

    // refreshToken 是自动续期的核心入口。
    // 只有在 authFetch 收到 401 且本地仍有 refresh token 时才会触发。
    // refresh 链路：当 access token 失效时，使用 refresh token 向 /v1/refresh 申请新的 access。
    async function refreshToken() {
        const tokens = load();
        if (!tokens.refresh) return null;
        const resp = await fetch(`${apiBase}/v1/refresh`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refresh_token: tokens.refresh })
        });
        if (!resp.ok) {
            clear();
            return null;
        }
        const data = await resp.json();
        const next = {
            access: data.token,
            refresh: data.refresh_token
        };
        save(next);
        return next.access;
    }

    // authFetch 是页面层最应该优先调用的带鉴权 fetch。
    // 它负责三件事：补 Authorization、遇到 401 尝试刷新、刷新成功后重放一次原请求。
    // 带认证请求入口：先附带 Authorization，再在收到 401 时自动尝试 refresh 并重放请求。
    async function authFetch(input, init = {}) {
        const tokens = load();
        const headers = new Headers(init.headers || {});
        if (tokens.access) {
            headers.set('Authorization', `Bearer ${tokens.access}`);
        }
        let resp = await fetch(input, { ...init, headers });
        if (resp.status === 401 && tokens.refresh) {
            const refreshed = await refreshToken();
            if (refreshed) {
                headers.set('Authorization', `Bearer ${refreshed}`);
                resp = await fetch(input, { ...init, headers });
            }
        }
        return resp;
    }

    // 登录成功后的写入口。
    function setTokens(access, refresh) {
        save({ access, refresh });
    }

    // 退出登录或 refresh 失败后的清理入口。
    function clear() {
        save({ access: '', refresh: '' });
    }

    // 暴露给页面的全局对象。todo-login.html 就是通过这里串起登录态和后续 API 调用。
    window.AuthClient = {
        setTokens,
        clear,
        authFetch,
        getAccessToken: () => load().access,
        getRefreshToken: () => load().refresh
    };
})();
