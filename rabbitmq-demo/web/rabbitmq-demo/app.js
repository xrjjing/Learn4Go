const logContainer = document.getElementById("log-container");
const sendForm = document.getElementById("send-form");
const sendResult = document.getElementById("send-result");
const mockBtn = document.getElementById("send-mock");
const mockResult = document.getElementById("mock-result");
const modeEl = document.getElementById("mode");
const connEl = document.getElementById("conn");
const queuesEl = document.getElementById("queues");
const errorBanner = document.getElementById("error-banner");
const errorMessage = document.getElementById("error-message");
const errorClose = document.getElementById("error-close");

let errorTimer;

// 显示全局错误提示
function showError(msg) {
  errorMessage.textContent = msg;
  errorBanner.classList.add("visible");
  errorBanner.setAttribute("aria-hidden", "false");
  clearTimeout(errorTimer);
  startErrorTimer();
}

// 启动自动隐藏计时器
function startErrorTimer() {
  errorTimer = setTimeout(hideError, 5000);
}

// 隐藏错误提示
function hideError() {
  errorBanner.classList.remove("visible");
  errorBanner.setAttribute("aria-hidden", "true");
  clearTimeout(errorTimer);
}

// 鼠标悬停暂停自动隐藏
errorBanner.addEventListener("mouseenter", () => clearTimeout(errorTimer));
errorBanner.addEventListener("mouseleave", () => startErrorTimer());
errorClose.addEventListener("click", hideError);

async function fetchLogs() {
  try {
    const res = await fetch("/api/logs");
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const data = await res.json();
    renderLogs(data);
  } catch (err) {
    console.error("拉取日志失败", err);
    showError("⚠️ 获取日志失败，请检查后端服务是否正常运行");
  }
}

async function fetchStatus() {
  try {
    const res = await fetch("/api/status");
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const data = await res.json();
    modeEl.textContent =
      data.mode === "fake" ? "内存模拟 (RABBITMQ_FAKE=1)" : "真实 MQ";
    connEl.textContent = data.rabbit_url || "-";
    renderQueues(data.queues || []);
  } catch (err) {
    modeEl.textContent = "状态获取失败";
    showError("⚠️ 获取运行状态失败，请检查后端服务");
  }
}

function renderQueues(queues) {
  queuesEl.innerHTML = "";
  queues.forEach((q) => {
    const div = document.createElement("div");
    div.className = "queue-card";
    div.innerHTML = `<div class=\"q-name\">${q.name}</div><div>消息: ${
      q.messages ?? "-"
    }</div><div>消费者: ${q.consumers ?? "-"}</div>`;
    queuesEl.appendChild(div);
  });
}

function renderLogs(logs) {
  logContainer.innerHTML = "";
  logs
    .slice()
    .reverse()
    .forEach((log) => {
      const div = document.createElement("div");
      div.className = `log log-${log.kind}`;
      const ts = new Date(log.time).toLocaleTimeString();
      div.textContent = `[${ts}] ${log.kind.toUpperCase()} ${log.id} ${
        log.type
      } - ${log.message}`;
      logContainer.appendChild(div);
    });
}

sendForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  sendResult.textContent = "发送中...";
  const type = document.getElementById("msg-type").value;
  const ttl = Number(document.getElementById("msg-ttl").value || 0);
  let payload;
  try {
    payload = JSON.parse(document.getElementById("msg-payload").value || "{}");
  } catch (err) {
    sendResult.textContent = "payload 不是合法 JSON";
    return;
  }
  const body = { type, payload };
  if (ttl > 0) body.ttl_ms = ttl;

  try {
    const res = await fetch("/api/messages", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    if (!res.ok) {
      const errText = await res.text().catch(() => "未知错误");
      throw new Error(errText);
    }
    const data = await res.json();
    sendResult.textContent = `✅ 已发送，ID=${data.id}`;
    sendResult.style.color = "var(--success)";
  } catch (err) {
    sendResult.textContent = "❌ 发送失败";
    sendResult.style.color = "var(--error)";
    showError(`❌ 消息发送失败：${err.message}`);
  }
});

mockBtn.addEventListener("click", async () => {
  mockResult.textContent = "发送中...";
  try {
    const res = await fetch("/api/messages/batch", { method: "POST" });
    if (!res.ok) {
      const errText = await res.text().catch(() => "未知错误");
      throw new Error(errText);
    }
    const data = await res.json();
    mockResult.textContent = `✅ 批量发送成功：${data.published} 条`;
    mockResult.style.color = "var(--success)";
  } catch (err) {
    mockResult.textContent = "❌ 批量发送失败";
    mockResult.style.color = "var(--error)";
    showError(`❌ 批量发送失败：${err.message}`);
  }
});

fetchLogs();
fetchStatus();
setInterval(fetchLogs, 3000);
setInterval(fetchStatus, 5000);
