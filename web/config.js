/**
 * Learn4Go 前端配置总入口。
 *
 * 文件职责：
 * 1. 在页面刚加载时写入 window.AppConfig，供 portal.html、login.html、todo-login.html、admin.html、log-detective.html 等页面读取。
 * 2. 统一决定当前是 Docker/Nginx 代理模式，还是本地直连模式。
 * 3. 给后续 JS 提供页面 -> API 的基址，不直接发请求，也不负责鉴权。
 *
 * 调用关系：
 * - HTML 页面先加载本文件，再由页面内联脚本或 auth-helper.js/common.js 读取 AppConfig。
 * - 这里的 todoApiBaseUrl 直接影响 todo-login.html 登录、portal.html 服务探测、以及部分真实接口调试。
 *
 * 排查建议：
 * - 页面请求打错地址、Docker 与本地地址混用时，先看这里。
 * - 如果页面拿不到 window.AppConfig，先检查本文件是否比业务脚本更早加载。
 */
// Docker 部署时使用相对路径（Nginx 代理）
// 本地开发时使用绝对地址
(function () {
  // 这里不是严格识别容器，而是识别“是否走本地开发域名”。
  // 一旦不是 localhost / 127.0.0.1，就默认按 Nginx 反向代理入口处理。
  const isDocker =
    window.location.hostname !== "localhost" &&
    window.location.hostname !== "127.0.0.1";

  // AppConfig 是整个前端最上游的环境配置对象。
  // 页面脚本、auth-helper、portal 服务状态探测都会从这里取基址。
  window.AppConfig = isDocker
    ? {
        // Docker 部署：使用相对路径，Nginx 会代理到后端
        todoApiBaseUrl: "/api",
        gatewayUrl: "/api",
        apiBaseUrl: "/api",
        logApiBaseUrl: "/api",
        enableMock: false,
      }
    : {
        // 本地开发：直连后端服务
        todoApiBaseUrl: "http://127.0.0.1:8080",
        gatewayUrl: "http://127.0.0.1:8888",
        apiBaseUrl: "http://127.0.0.1:8080", // 统一使用 8080 端口（RBAC API 已集成）
        logApiBaseUrl: "http://127.0.0.1:8002",
        enableMock: true,
      };
})();
