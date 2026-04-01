/**
 * Learn4Go 前端基础工具库。
 *
 * 服务页面：
 * - 当前最直接服务 login.html。
 * - 也可被后续其他页面复用，用来统一 Toast、按钮 loading、表单校验和多语言文案。
 *
 * 结构说明：
 * - i18n：文案字典与语言切换。
 * - Toast：顶部提示。
 * - Loading：按钮和全屏加载态。
 * - FormValidator：表单字段校验与错误展示。
 *
 * 排查建议：
 * - 登录页按钮一直转圈，先看 Loading.setButton() 是否正确恢复。
 * - 表单校验提示不出现，先看 FormValidator.validate() 和页面里 data-error/name 是否对齐。
 */

// ==================== 国际化 (P2-F5) ====================
// i18n 是最基础的文案层；登录页提交前后的提示文本都从这里取。
const i18n = {
  locale: 'zh-CN',
  messages: {
    'zh-CN': {
      login: {
        title: '欢迎回来',
        subtitle: '登录访问 Go 学习平台',
        username: '用户名',
        password: '密码',
        submit: '登录',
        loading: '登录中...',
        error: {
          required: '请填写所有必填字段',
          email: '请输入有效的邮箱地址',
          password: '密码至少需要 8 个字符',
          network: '网络错误，请检查服务是否启动',
          invalid: '用户名或密码错误'
        }
      },
      common: {
        loading: '加载中...',
        error: '操作失败',
        success: '操作成功',
        confirm: '确认',
        cancel: '取消'
      }
    },
    'en': {
      login: {
        title: 'Welcome Back',
        subtitle: 'Sign in to access the Go Learning Platform',
        username: 'Username',
        password: 'Password',
        submit: 'Login',
        loading: 'Logging in...',
        error: {
          required: 'Please fill in all required fields',
          email: 'Please enter a valid email address',
          password: 'Password must be at least 8 characters',
          network: 'Connection error. Please ensure service is running.',
          invalid: 'Invalid username or password'
        }
      },
      common: {
        loading: 'Loading...',
        error: 'Operation failed',
        success: 'Operation successful',
        confirm: 'Confirm',
        cancel: 'Cancel'
      }
    }
  },
  t(key) {
    const keys = key.split('.');
    let value = this.messages[this.locale];
    for (const k of keys) {
      value = value?.[k];
    }
    return value || key;
  },
  setLocale(locale) {
    if (this.messages[locale]) {
      this.locale = locale;
      localStorage.setItem('locale', locale);
    }
  }
};

// 初始化语言设置。这里是文件加载时立即执行的轻量初始化入口。
try {
  const savedLocale = localStorage.getItem('locale');
  if (savedLocale && i18n.messages[savedLocale]) {
    i18n.locale = savedLocale;
  }
} catch (e) {}

// ==================== Toast 通知 (P2-F3) ====================
// Toast 负责页面右上角即时反馈，适合登录成功/失败、保存成功等短消息。
// Toast 组件：适合给登录、保存、失败等轻量反馈使用。
const Toast = {
  container: null,

  init() {
    if (this.container) return;
    this.container = document.createElement('div');
    this.container.id = 'toast-container';
    this.container.style.cssText = `
      position: fixed;
      top: 1rem;
      right: 1rem;
      z-index: 9999;
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
      max-width: 400px;
    `;
    document.body.appendChild(this.container);
  },

  show(message, type = 'info', duration = 4000) {
    this.init();

    const colors = {
      success: { bg: '#10b981', icon: '✓' },
      error: { bg: '#ef4444', icon: '✕' },
      warning: { bg: '#f59e0b', icon: '⚠' },
      info: { bg: '#3b82f6', icon: 'ℹ' }
    };

    const { bg, icon } = colors[type] || colors.info;

    const toast = document.createElement('div');
    toast.style.cssText = `
      background: ${bg};
      color: white;
      padding: 0.75rem 1rem;
      border-radius: 8px;
      box-shadow: 0 4px 12px rgba(0,0,0,0.15);
      display: flex;
      align-items: center;
      gap: 0.5rem;
      animation: slideIn 0.3s ease;
      font-size: 0.9rem;
    `;

    // 使用 DOM API 防止 XSS
    const iconSpan = document.createElement('span');
    iconSpan.style.fontSize = '1.1em';
    iconSpan.textContent = icon;

    const msgSpan = document.createElement('span');
    msgSpan.textContent = message;

    toast.appendChild(iconSpan);
    toast.appendChild(msgSpan);

    this.container.appendChild(toast);

    setTimeout(() => {
      toast.style.animation = 'slideOut 0.3s ease';
      setTimeout(() => toast.remove(), 300);
    }, duration);
  },

  success(msg) { this.show(msg, 'success'); },
  error(msg) { this.show(msg, 'error'); },
  warning(msg) { this.show(msg, 'warning'); },
  info(msg) { this.show(msg, 'info'); }
};

// 添加动画样式
const style = document.createElement('style');
style.textContent = `
  @keyframes slideIn {
    from { transform: translateX(100%); opacity: 0; }
    to { transform: translateX(0); opacity: 1; }
  }
  @keyframes slideOut {
    from { transform: translateX(0); opacity: 1; }
    to { transform: translateX(100%); opacity: 0; }
  }
  .btn-loading {
    position: relative;
    color: transparent !important;
  }
  .btn-loading::after {
    content: '';
    position: absolute;
    width: 1rem;
    height: 1rem;
    top: 50%;
    left: 50%;
    margin: -0.5rem 0 0 -0.5rem;
    border: 2px solid rgba(255,255,255,0.3);
    border-top-color: white;
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
  }
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
  .form-error {
    color: #ef4444;
    font-size: 0.8rem;
    margin-top: 0.25rem;
  }
  .input-error {
    border-color: #ef4444 !important;
  }
`;
document.head.appendChild(style);

// ==================== 加载状态 (P2-F2) ====================
// Loading 同时覆盖“按钮 loading”与“全屏遮罩”两类场景。
// Loading 组件：负责按钮忙碌态和覆盖层忙碌态。
const Loading = {
  setButton(btn, loading) {
    if (loading) {
      btn.classList.add('btn-loading');
      btn.disabled = true;
      btn.dataset.originalText = btn.textContent;
    } else {
      btn.classList.remove('btn-loading');
      btn.disabled = false;
      if (btn.dataset.originalText) {
        btn.textContent = btn.dataset.originalText;
      }
    }
  },

  overlay: null,

  show(message = i18n.t('common.loading')) {
    if (this.overlay) return;
    this.overlay = document.createElement('div');
    this.overlay.style.cssText = `
      position: fixed;
      inset: 0;
      background: rgba(0,0,0,0.5);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 9998;
    `;

    // 使用 DOM API 防止 XSS
    const container = document.createElement('div');
    container.style.cssText = 'background:var(--bg-panel,white);padding:2rem;border-radius:8px;text-align:center';

    const spinner = document.createElement('div');
    spinner.style.cssText = 'width:40px;height:40px;border:3px solid var(--border,#e5e7eb);border-top-color:var(--primary,#4f46e5);border-radius:50%;margin:0 auto 1rem;animation:spin 0.8s linear infinite';

    const text = document.createElement('p');
    text.style.color = 'var(--text-main,#334155)';
    text.textContent = message;

    container.appendChild(spinner);
    container.appendChild(text);
    this.overlay.appendChild(container);

    document.body.appendChild(this.overlay);
  },

  hide() {
    this.overlay?.remove();
    this.overlay = null;
  }
};

// ==================== 表单验证 (P2-F1) ====================
// FormValidator 是 login.html 提交前的第一道前端校验。
// 它只做页面体验层的拦截，后端校验仍以 API 返回为准。
// 表单校验器：login.html 会通过规则表驱动最小校验逻辑。
const FormValidator = {
  rules: {
    required: (value) => value?.trim() ? null : 'required',
    email: (value) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value) ? null : 'email',
    minLength: (min) => (value) => value?.length >= min ? null : `minLength:${min}`,
    maxLength: (max) => (value) => value?.length <= max ? null : `maxLength:${max}`,
    password: (value) => value?.length >= 8 ? null : 'password'
  },

  validate(form, config) {
    const errors = {};
    let isValid = true;

    for (const [field, rules] of Object.entries(config)) {
      const input = form.querySelector(`[name="${field}"]`);
      if (!input) continue;

      const value = input.value;
      const errorEl = form.querySelector(`[data-error="${field}"]`);

      input.classList.remove('input-error');
      if (errorEl) errorEl.textContent = '';

      for (const rule of rules) {
        const error = typeof rule === 'function' ? rule(value) : this.rules[rule]?.(value);
        if (error) {
          errors[field] = error;
          isValid = false;
          input.classList.add('input-error');
          if (errorEl) {
            errorEl.textContent = i18n.t(`login.error.${error.split(':')[0]}`) || error;
          }
          break;
        }
      }
    }

    return { isValid, errors };
  },

  clearErrors(form) {
    form.querySelectorAll('.input-error').forEach(el => el.classList.remove('input-error'));
    form.querySelectorAll('[data-error]').forEach(el => el.textContent = '');
  }
};

// 导出到全局。页面一般通过 window.Learn4Go.xxx 读取，而不是模块化 import。
// 全局暴露：静态 HTML 页面不走模块打包，因此通过 window.Learn4Go 给内联脚本复用。
window.Learn4Go = { i18n, Toast, Loading, FormValidator };
