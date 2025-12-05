// 认证辅助：存取 token + 自动刷新
(function () {
    const STORAGE_KEY = 'learn4go_auth';
    const apiBase = (window.AppConfig && window.AppConfig.todoApiBaseUrl) || '';

    function load() {
        try {
            const raw = localStorage.getItem(STORAGE_KEY);
            if (!raw) return { access: '', refresh: '' };
            return JSON.parse(raw);
        } catch (e) {
            return { access: '', refresh: '' };
        }
    }

    function save(tokens) {
        try {
            localStorage.setItem(STORAGE_KEY, JSON.stringify(tokens));
        } catch (e) {
            /* 本地存储失败时静默忽略 */
        }
    }

    async function refreshToken() {
        const tokens = load();
        if (!tokens.refresh) return null;
        const resp = await fetch(`${apiBase}/refresh`, {
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

    function setTokens(access, refresh) {
        save({ access, refresh });
    }

    function clear() {
        save({ access: '', refresh: '' });
    }

    window.AuthClient = {
        setTokens,
        clear,
        authFetch,
        getAccessToken: () => load().access,
        getRefreshToken: () => load().refresh
    };
})();
