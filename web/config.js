// Learn4Go 前端配置
// Docker 部署时使用相对路径（Nginx 代理）
// 本地开发时使用绝对地址
(function() {
    const isDocker = window.location.hostname !== 'localhost' && window.location.hostname !== '127.0.0.1';

    window.AppConfig = isDocker ? {
        // Docker 部署：使用相对路径，Nginx 会代理到后端
        todoApiBaseUrl: '/api/todos',
        gatewayUrl: '/api',
        apiBaseUrl: '/api',
        logApiBaseUrl: '/api',
        enableMock: false
    } : {
        // 本地开发：直连后端服务
        todoApiBaseUrl: 'http://127.0.0.1:8080',
        gatewayUrl: 'http://127.0.0.1:8888',
        apiBaseUrl: 'http://127.0.0.1:8000',
        logApiBaseUrl: 'http://127.0.0.1:8002',
        enableMock: true
    };
})();
