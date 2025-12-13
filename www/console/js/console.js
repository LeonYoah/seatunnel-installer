/* SeaTunnel Console Prototype (no framework) */

const qs = (selector, root = document) => root.querySelector(selector);
const qsa = (selector, root = document) => Array.from(root.querySelectorAll(selector));

const store = {
  get(key, fallback) {
    try {
      const raw = localStorage.getItem(key);
      return raw ? JSON.parse(raw) : fallback;
    } catch {
      return fallback;
    }
  },
  set(key, value) {
    try {
      localStorage.setItem(key, JSON.stringify(value));
    } catch {
      // ignore
    }
  }
};

const state = {
  cluster: store.get("st.console.cluster", "default"),
  dataMode: "prototype",
  lastRefreshAt: null
};

function escapeHtml(text) {
  return String(text)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

function parseHash() {
  const raw = location.hash || "#/dashboard";
  const [pathPart, queryPart] = raw.replace(/^#/, "").split("?");
  const path = pathPart.startsWith("/") ? pathPart : `/${pathPart}`;
  const query = new URLSearchParams(queryPart || "");
  return { path, query };
}

function setTitle(title) {
  const el = qs("#pageTitle");
  if (el) el.textContent = title;
  document.title = `${title} - SeaTunnel 运维控制台（原型）`;
}

function setActiveNav(path) {
  qsa(".nav-item").forEach(a => {
    a.classList.toggle("active", a.dataset.route === path);
  });
}

function toast(message) {
  const el = qs("#toast");
  if (!el) return;
  el.textContent = message;
  el.classList.add("show");
  clearTimeout(toast._t);
  toast._t = setTimeout(() => el.classList.remove("show"), 1800);
}

function nowTime() {
  const d = new Date();
  return d.toLocaleString();
}

function mockData() {
  return {
    cluster: {
      name: state.cluster,
      version: "2.3.12",
      nodes: 6,
      up: 6,
      cpu: 0.42,
      mem: 0.61
    },
    kpis: {
      successRate: 99.2,
      p95LatencyMs: 820,
      throughputRows: 180000,
      alerts: 2
    },
    alerts: [
      { level: "WARN", target: "task:orders_sync", at: "10m ago", msg: "延迟抖动，P95 > 1s" },
      { level: "ERROR", target: "node:worker-03", at: "2h ago", msg: "心跳超时（已恢复）" }
    ],
    tasks: [
      { id: "orders_sync", name: "orders_sync", type: "实时", status: "running", lastRun: "持续运行", duration: "-", version: "v14" },
      { id: "users_dim", name: "users_dim", type: "批", status: "success", lastRun: "2025-12-13 09:02", duration: "3m12s", version: "v8" },
      { id: "refunds_backfill", name: "refunds_backfill", type: "批", status: "failed", lastRun: "2025-12-13 08:40", duration: "1m05s", version: "v3" },
      { id: "cdc_mysql_kafka", name: "cdc_mysql_kafka", type: "实时", status: "paused", lastRun: "暂停", duration: "-", version: "v21" }
    ],
    nodes: [
      { name: "master-01", role: "master", cpu: 0.21, mem: 0.34, disk: 0.48, status: "up", lastSeen: "5s ago" },
      { name: "worker-01", role: "worker", cpu: 0.48, mem: 0.61, disk: 0.52, status: "up", lastSeen: "5s ago" },
      { name: "worker-02", role: "worker", cpu: 0.55, mem: 0.72, disk: 0.60, status: "up", lastSeen: "6s ago" },
      { name: "worker-03", role: "worker", cpu: 0.10, mem: 0.22, disk: 0.40, status: "up", lastSeen: "9s ago" }
    ]
  };
}

function badgeForTaskStatus(status) {
  const map = {
    running: { cls: "running", text: "运行中" },
    success: { cls: "success", text: "成功" },
    failed: { cls: "failed", text: "失败" },
    paused: { cls: "paused", text: "暂停" }
  };
  const v = map[status] || { cls: "", text: status };
  return `<span class="badge ${v.cls}">${escapeHtml(v.text)}</span>`;
}

function chipForAlertLevel(level) {
  if (level === "ERROR") return `<span class="chip chip-bad">ERROR</span>`;
  if (level === "WARN") return `<span class="chip chip-warn">WARN</span>`;
  return `<span class="chip chip-info">${escapeHtml(level)}</span>`;
}

function renderTasksTable(tasks, opts = {}) {
  const compact = Boolean(opts.compact);
  return `
    <table>
      <thead>
        <tr>
          <th style="width:26%">任务</th>
          <th style="width:10%">类型</th>
          <th style="width:12%">状态</th>
          <th style="width:16%">最近运行</th>
          <th style="width:12%">耗时</th>
          <th style="width:10%">版本</th>
          <th style="width:${compact ? "14" : "24"}%">操作</th>
        </tr>
      </thead>
      <tbody>
        ${tasks.map(t => `
          <tr>
            <td><a class="link" href="#/task?id=${encodeURIComponent(t.id)}">${escapeHtml(t.name)}</a></td>
            <td>${escapeHtml(t.type)}</td>
            <td>${badgeForTaskStatus(t.status)}</td>
            <td class="muted">${escapeHtml(t.lastRun)}</td>
            <td class="muted">${escapeHtml(t.duration)}</td>
            <td class="muted">${escapeHtml(t.version)}</td>
            <td>
              <div class="table-actions">
                <button class="btn btn-sm btn-outline" type="button" data-action="task-run" data-id="${escapeHtml(t.id)}">重试</button>
                <button class="btn btn-sm btn-outline" type="button" data-action="task-rollback" data-id="${escapeHtml(t.id)}">回滚</button>
                ${compact ? "" : `<button class="btn btn-sm btn-danger" type="button" data-action="task-stop" data-id="${escapeHtml(t.id)}">停止</button>`}
              </div>
            </td>
          </tr>
        `).join("")}
      </tbody>
    </table>
  `;
}

function renderDashboard(data) {
  setTitle("总览");
  const c = data.cluster;
  const k = data.kpis;
  return `
    <div class="grid">
      <div class="col-3 card">
        <div class="card-title">集群健康 <span class="chip chip-ok">${escapeHtml(c.up)}/${escapeHtml(c.nodes)} UP</span></div>
        <div class="kpi">
          <div class="kpi-value">${escapeHtml(c.name)}</div>
          <div class="kpi-label">SeaTunnel ${escapeHtml(c.version)}</div>
        </div>
      </div>
      <div class="col-3 card">
        <div class="card-title">任务成功率 <span class="chip chip-ok">SLO</span></div>
        <div class="kpi">
          <div class="kpi-value">${escapeHtml(k.successRate)}%</div>
          <div class="kpi-label">过去 24 小时</div>
        </div>
      </div>
      <div class="col-3 card">
        <div class="card-title">延迟 P95 <span class="chip chip-warn">监控</span></div>
        <div class="kpi">
          <div class="kpi-value">${escapeHtml(k.p95LatencyMs)} ms</div>
          <div class="kpi-label">过去 1 小时</div>
        </div>
      </div>
      <div class="col-3 card">
        <div class="card-title">告警 <span class="chip ${k.alerts ? "chip-warn" : "chip-ok"}">${escapeHtml(k.alerts)} 条</span></div>
        <div class="kpi">
          <div class="kpi-value">${escapeHtml(k.throughputRows.toLocaleString())}</div>
          <div class="kpi-label">吞吐（rows/min）</div>
        </div>
      </div>

      <div class="col-8 card">
        <div class="card-title">趋势（占位） <span class="card-subtitle">后续接入 Prometheus/Grafana 数据</span></div>
        <div class="empty">这里将展示：成功率/延迟/吞吐/资源利用趋势图。</div>
      </div>
      <div class="col-4 card">
        <div class="card-title">告警（最近） <span class="card-subtitle">支持静默/跳转</span></div>
        <div class="row" style="flex-direction:column; align-items:stretch; gap:10px">
          ${data.alerts
            .map(a => `
              <div class="panel">
                <div class="row">
                  ${chipForAlertLevel(a.level)}
                  <div class="muted">${escapeHtml(a.target)}</div>
                  <div class="spacer"></div>
                  <div class="muted">${escapeHtml(a.at)}</div>
                </div>
                <div style="margin-top:6px">${escapeHtml(a.msg)}</div>
              </div>`)
            .join("")}
        </div>
      </div>

      <div class="col-12 card">
        <div class="card-title">最近任务 <span class="card-subtitle">点击进入详情</span></div>
        ${renderTasksTable(data.tasks, { compact: true })}
      </div>
    </div>
  `;
}

function renderTasks(data, query) {
  setTitle("任务");
  const filter = (query.get("status") || "").toLowerCase();
  const q = (query.get("q") || "").trim().toLowerCase();
  let tasks = data.tasks.slice();
  if (filter) tasks = tasks.filter(t => t.status === filter);
  if (q) tasks = tasks.filter(t => t.name.toLowerCase().includes(q) || t.id.toLowerCase().includes(q));

  const chips = [
    { label: "全部", href: "#/tasks", active: !filter },
    { label: "运行中", href: "#/tasks?status=running", active: filter === "running" },
    { label: "成功", href: "#/tasks?status=success", active: filter === "success" },
    { label: "失败", href: "#/tasks?status=failed", active: filter === "failed" },
    { label: "暂停", href: "#/tasks?status=paused", active: filter === "paused" }
  ];

  return `
    <div class="grid">
      <div class="col-12 card">
        <div class="card-title">
          任务列表
          <div class="row">
            ${chips.map(c => `<a class="chip ${c.active ? "chip-info" : ""}" href="${c.href}">${escapeHtml(c.label)}</a>`).join("")}
            <div class="spacer"></div>
            <a class="btn btn-sm" href="#/task-new">新建任务</a>
          </div>
        </div>
        <div class="row" style="margin-bottom:10px">
          <input class="input input-sm" id="taskSearch" placeholder="按任务名搜索…" value="${escapeHtml(query.get("q") || "")}" />
          <button class="btn btn-sm btn-outline" type="button" id="btnTaskSearch">搜索</button>
          <div class="spacer"></div>
          <button class="btn btn-sm btn-outline" type="button" data-action="task-import">导入</button>
          <button class="btn btn-sm btn-outline" type="button" data-action="task-export">导出</button>
          <button class="btn btn-sm btn-outline" type="button" data-action="task-batch-retry">批量重试</button>
        </div>
        ${renderTasksTable(tasks)}
      </div>
    </div>
  `;
}

function renderTaskDetail(data, query) {
  const id = query.get("id") || "";
  const task = data.tasks.find(t => t.id === id);
  if (!task) {
    setTitle("任务详情");
    return `<div class="empty">未找到任务：${escapeHtml(id)}。返回 <a class="link" href="#/tasks">任务列表</a>。</div>`;
  }

  setTitle(`任务：${task.name}`);
  const mockLog = [
    "[INFO] Build pipeline ...",
    "[INFO] Source=mysql -> Transform=sql -> Sink=iceberg",
    "[WARN] P95 latency > 1s (backpressure detected)",
    "[INFO] checkpoint completed: id=3821",
    "[INFO] throughput=180k rows/min"
  ].join("\\n");

  return `
    <div class="grid">
      <div class="col-12 card">
        <div class="card-title">
          <div class="row">
            <a class="link" href="#/tasks">← 返回列表</a>
            <div style="font-weight:800">${escapeHtml(task.name)}</div>
            ${badgeForTaskStatus(task.status)}
            <span class="chip chip-info">${escapeHtml(task.type)}</span>
            <span class="chip">版本 ${escapeHtml(task.version)}</span>
          </div>
          <div class="row">
            <button class="btn btn-sm btn-outline" type="button" data-action="task-retry" data-id="${escapeHtml(task.id)}">重试</button>
            <button class="btn btn-sm btn-outline" type="button" data-action="task-rollback" data-id="${escapeHtml(task.id)}">回滚</button>
            <button class="btn btn-sm btn-danger" type="button" data-action="task-stop" data-id="${escapeHtml(task.id)}">停止</button>
          </div>
        </div>

        <div class="grid">
          <div class="col-3 card">
            <div class="card-title">延迟 P95</div>
            <div class="kpi"><div class="kpi-value">820 ms</div><div class="kpi-label">过去 15 分钟</div></div>
          </div>
          <div class="col-3 card">
            <div class="card-title">吞吐</div>
            <div class="kpi"><div class="kpi-value">180k</div><div class="kpi-label">rows/min</div></div>
          </div>
          <div class="col-3 card">
            <div class="card-title">失败率</div>
            <div class="kpi"><div class="kpi-value">0.8%</div><div class="kpi-label">过去 24 小时</div></div>
          </div>
          <div class="col-3 card">
            <div class="card-title">队列/背压</div>
            <div class="kpi"><div class="kpi-value">轻微</div><div class="kpi-label">建议检查 Sink</div></div>
          </div>
        </div>
      </div>

      <div class="col-6 card">
        <div class="card-title">运行时间线 <span class="card-subtitle">提交 → 调度 → 运行 → 结束/失败</span></div>
        <div class="panel">
          <div class="row"><span class="chip chip-info">提交</span><span class="muted">2025-12-13 08:40</span></div>
          <div class="row" style="margin-top:8px"><span class="chip chip-info">调度</span><span class="muted">+3s</span></div>
          <div class="row" style="margin-top:8px"><span class="chip chip-info">运行</span><span class="muted">持续运行 / 耗时 ${escapeHtml(task.duration)}</span></div>
          <div class="row" style="margin-top:8px"><span class="chip chip-warn">事件</span><span class="muted">检测到背压，已自动降速（占位）</span></div>
        </div>
      </div>

      <div class="col-6 card">
        <div class="card-title">实时日志 <span class="card-subtitle">后续接入 /cgi-bin/run.sh?action=task_log</span></div>
        <div class="panel"><pre>${escapeHtml(mockLog)}</pre></div>
      </div>

      <div class="col-12 card">
        <div class="card-title">配置与版本 <span class="card-subtitle">YAML/表单编辑、Diff/回滚（占位）</span></div>
        <div class="grid">
          <div class="col-6">
            <div class="panel"><pre>env {
  parallelism = 2
}

source {
  MySQL-CDC { ... }
}

sink {
  Iceberg { ... }
}</pre></div>
          </div>
          <div class="col-6">
            <div class="empty">这里将展示版本 Diff、发布策略（滚动/金丝雀）与预检结果。</div>
          </div>
        </div>
      </div>
    </div>
  `;
}

function renderTaskNew() {
  setTitle("新建任务");
  return `
    <div class="grid">
      <div class="col-12 card">
        <div class="card-title">
          新建任务 <span class="card-subtitle">模板 + YAML/表单双模式（原型）</span>
          <div class="row">
            <a class="btn btn-sm btn-outline" href="#/tasks">返回</a>
          </div>
        </div>
        <div class="row">
          <label class="field" style="min-width:260px">
            <span class="field-label">模板</span>
            <select class="input" id="tpl">
              <option value="stream">实时同步（CDC→Kafka）</option>
              <option value="batch">离线同步（JDBC→Iceberg）</option>
              <option value="dag">DAG（多源多汇）</option>
            </select>
          </label>
          <label class="field" style="flex:1">
            <span class="field-label">任务名</span>
            <input class="input" id="taskName" placeholder="orders_sync" />
          </label>
        </div>
        <div class="row" style="margin-top:10px">
          <label class="field" style="flex:1">
            <span class="field-label">YAML</span>
            <textarea class="input" id="taskYaml" style="min-height:240px; resize:vertical">env {
  parallelism = 2
}

source {
  MySQL-CDC { ... }
}

sink {
  Kafka { ... }
}</textarea>
          </label>
        </div>
        <div class="row" style="margin-top:12px">
          <button class="btn" type="button" data-action="task-save">保存版本</button>
          <button class="btn btn-outline" type="button" data-action="task-validate">预检</button>
          <button class="btn btn-outline" type="button" data-action="task-submit">提交运行</button>
          <div class="spacer"></div>
          <span class="muted">建议：接入 GitOps（保存即 PR），支持 Diff/回滚。</span>
        </div>
      </div>
    </div>
  `;
}

function renderClusters(data) {
  setTitle("集群/节点");
  const c = data.cluster;
  return `
    <div class="grid">
      <div class="col-4 card">
        <div class="card-title">集群概览 <span class="card-subtitle">${escapeHtml(c.name)}</span></div>
        <div class="row">
          <div class="kpi"><div class="kpi-value">${escapeHtml(c.nodes)}</div><div class="kpi-label">节点数</div></div>
          <div class="kpi"><div class="kpi-value">${escapeHtml(Math.round(c.cpu*100))}%</div><div class="kpi-label">CPU 水位</div></div>
          <div class="kpi"><div class="kpi-value">${escapeHtml(Math.round(c.mem*100))}%</div><div class="kpi-label">内存水位</div></div>
        </div>
      </div>
      <div class="col-8 card">
        <div class="card-title">
          节点列表
          <div class="row">
            <button class="btn btn-sm btn-outline" type="button" data-action="node-restart">重启</button>
            <button class="btn btn-sm btn-outline" type="button" data-action="node-isolate">隔离</button>
            <button class="btn btn-sm btn-outline" type="button" data-action="cluster-scale">扩/缩容</button>
          </div>
        </div>
        <table>
          <thead>
            <tr>
              <th style="width:22%">节点</th>
              <th style="width:12%">角色</th>
              <th style="width:12%">CPU</th>
              <th style="width:12%">内存</th>
              <th style="width:12%">磁盘</th>
              <th style="width:12%">状态</th>
              <th style="width:18%">心跳</th>
            </tr>
          </thead>
          <tbody>
            ${data.nodes.map(n => `
              <tr>
                <td>${escapeHtml(n.name)}</td>
                <td class="muted">${escapeHtml(n.role)}</td>
                <td>${escapeHtml(Math.round(n.cpu*100))}%</td>
                <td>${escapeHtml(Math.round(n.mem*100))}%</td>
                <td>${escapeHtml(Math.round(n.disk*100))}%</td>
                <td>${n.status === "up" ? `<span class="badge success">UP</span>` : `<span class="badge failed">DOWN</span>`}</td>
                <td class="muted">${escapeHtml(n.lastSeen)}</td>
              </tr>
            `).join("")}
          </tbody>
        </table>
      </div>

      <div class="col-12 card">
        <div class="card-title">发布与扩缩容（占位） <span class="card-subtitle">滚动/金丝雀/蓝绿 + 自动策略</span></div>
        <div class="empty">这里将展示发布进度、回滚、一键扩缩容策略（阈值/预测）与执行记录。</div>
      </div>
    </div>
  `;
}

function renderDiagnostics() {
  setTitle("诊断");
  return `
    <div class="grid">
      <div class="col-6 card">
        <div class="card-title">一键诊断 <span class="card-subtitle">生成诊断包供下载/工单</span></div>
        <div class="row">
          <label class="field" style="min-width:220px">
            <span class="field-label">范围</span>
            <select class="input" id="diagScope">
              <option value="cluster">集群</option>
              <option value="node">指定节点</option>
              <option value="task">指定任务</option>
            </select>
          </label>
          <label class="field" style="min-width:220px">
            <span class="field-label">对象</span>
            <input class="input" id="diagTarget" placeholder="如 worker-01 / orders_sync" />
          </label>
        </div>
        <div class="row" style="margin-top:10px">
          <label class="chip"><input type="checkbox" checked /> 日志</label>
          <label class="chip"><input type="checkbox" checked /> 配置快照</label>
          <label class="chip"><input type="checkbox" /> 线程 dump</label>
          <label class="chip"><input type="checkbox" /> GC/Heap</label>
        </div>
        <div class="row" style="margin-top:12px">
          <button class="btn" type="button" data-action="diag-run">生成诊断包</button>
          <button class="btn btn-outline" type="button" data-action="diag-history">查看历史</button>
        </div>
        <div class="empty" style="margin-top:12px">后端落地建议：新增 diagnose.sh/Agent，上报诊断进度；UI 展示下载链接。</div>
      </div>

      <div class="col-6 card">
        <div class="card-title">常见故障库 <span class="card-subtitle">日志指纹 → 建议/一键修复</span></div>
        <div class="panel" style="margin-bottom:10px">
          <div class="row">
            <span class="chip chip-warn">WARN</span>
            <div><div>背压导致延迟升高</div><div class="muted">特征：backpressure / sink slow</div></div>
            <div class="spacer"></div>
            <button class="btn btn-sm btn-outline" type="button" data-action="kb-open">查看</button>
          </div>
        </div>
        <div class="panel">
          <div class="row">
            <span class="chip chip-bad">ERROR</span>
            <div><div>Connector 依赖缺失</div><div class="muted">特征：ClassNotFoundException</div></div>
            <div class="spacer"></div>
            <button class="btn btn-sm btn-outline" type="button" data-action="kb-fix">一键修复</button>
          </div>
        </div>
        <div class="empty" style="margin-top:12px">建议把故障库做成可扩展的 JSON/YAML（社区可贡献）。</div>
      </div>
    </div>
  `;
}

function renderPlugins() {
  setTitle("插件市场");
  const rows = [
    { name: "MySQL-CDC", type: "Source", ver: "2.3.x", status: "installed" },
    { name: "Kafka", type: "Sink", ver: "2.3.x", status: "installed" },
    { name: "Iceberg", type: "Sink", ver: "2.3.x", status: "available" },
    { name: "Doris", type: "Sink", ver: "2.3.x", status: "available" }
  ];
  return `
    <div class="grid">
      <div class="col-12 card">
        <div class="card-title">
          Connector/插件
          <div class="row">
            <button class="btn btn-sm btn-outline" type="button" data-action="plugin-refresh">刷新仓库</button>
            <button class="btn btn-sm btn-outline" type="button" data-action="plugin-upload">本地上传</button>
          </div>
        </div>
        <table>
          <thead>
            <tr>
              <th style="width:30%">名称</th>
              <th style="width:12%">类型</th>
              <th style="width:12%">兼容性</th>
              <th style="width:16%">状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            ${rows.map(r => `
              <tr>
                <td>${escapeHtml(r.name)}</td>
                <td class="muted">${escapeHtml(r.type)}</td>
                <td class="muted">${escapeHtml(r.ver)}</td>
                <td>${r.status === "installed" ? `<span class="badge success">已安装</span>` : `<span class="badge running">可用</span>`}</td>
                <td>
                  <div class="table-actions">
                    ${r.status === "installed"
                      ? `<button class="btn btn-sm btn-outline" type="button" data-action="plugin-upgrade">升级</button>
                         <button class="btn btn-sm btn-danger" type="button" data-action="plugin-disable">禁用</button>`
                      : `<button class="btn btn-sm" type="button" data-action="plugin-install">安装</button>`}
                    <button class="btn btn-sm btn-outline" type="button" data-action="plugin-detail">详情</button>
                  </div>
                </td>
              </tr>
            `).join("")}
          </tbody>
        </table>
        <div class="empty" style="margin-top:12px">后续可接：离线插件包/在线仓库、依赖检查、版本锁定与回滚。</div>
      </div>
    </div>
  `;
}

function renderSettings() {
  setTitle("设置");
  return `
    <div class="grid">
      <div class="col-6 card">
        <div class="card-title">通知渠道 <span class="card-subtitle">钉钉/企微/邮件</span></div>
        <div class="row">
          <label class="field" style="flex:1">
            <span class="field-label">Webhook（钉钉/企微）</span>
            <input class="input" placeholder="https://..." />
          </label>
        </div>
        <div class="row" style="margin-top:10px">
          <label class="field" style="flex:1">
            <span class="field-label">邮件接收人</span>
            <input class="input" placeholder="ops@example.com" />
          </label>
        </div>
        <div class="row" style="margin-top:12px">
          <button class="btn" type="button" data-action="settings-save">保存</button>
          <button class="btn btn-outline" type="button" data-action="settings-test">发送测试</button>
        </div>
      </div>

      <div class="col-6 card">
        <div class="card-title">RBAC / 多租户 <span class="card-subtitle">占位</span></div>
        <div class="empty">后续落地：用户/角色/命名空间/队列隔离 + 操作审计。</div>
      </div>

      <div class="col-12 card">
        <div class="card-title">审计日志（占位）</div>
        <div class="empty">这里将展示：操作时间、操作者、对象、结果；支持导出。</div>
      </div>
    </div>
  `;
}

const routes = {
  "/dashboard": ({ data }) => renderDashboard(data),
  "/tasks": ({ data, query }) => renderTasks(data, query),
  "/task": ({ data, query }) => renderTaskDetail(data, query),
  "/task-new": () => renderTaskNew(),
  "/clusters": ({ data }) => renderClusters(data),
  "/diagnostics": () => renderDiagnostics(),
  "/plugins": () => renderPlugins(),
  "/settings": () => renderSettings()
};

function render() {
  const { path, query } = parseHash();
  const data = mockData();
  const fn = routes[path] || routes["/dashboard"];
  const routePath = routes[path] ? path : "/dashboard";
  setActiveNav(routePath);

  const content = qs("#appContent");
  content.innerHTML = fn({ data, query });

  const meta = qs("#sidebarMeta");
  if (meta) {
    meta.textContent = `cluster=${state.cluster} · ${state.lastRefreshAt ? `updated ${state.lastRefreshAt}` : "not refreshed"}`;
  }

  const searchBtn = qs("#btnTaskSearch");
  if (searchBtn) {
    searchBtn.addEventListener("click", () => {
      const q = qs("#taskSearch").value.trim();
      const params = new URLSearchParams(location.hash.split("?")[1] || "");
      if (q) params.set("q", q);
      else params.delete("q");
      location.hash = `#/tasks?${params.toString()}`;
    });
  }
}

function bindGlobalEvents() {
  const sel = qs("#clusterSelect");
  sel.value = state.cluster;
  sel.addEventListener("change", () => {
    state.cluster = sel.value;
    store.set("st.console.cluster", state.cluster);
    toast(`切换集群：${state.cluster}`);
    render();
  });

  qs("#btnRefresh").addEventListener("click", () => {
    state.lastRefreshAt = nowTime();
    toast("已刷新（原型数据）");
    render();
  });

  qs("#globalSearch").addEventListener("keydown", (e) => {
    if (e.key !== "Enter") return;
    const q = qs("#globalSearch").value.trim();
    location.hash = `#/tasks?q=${encodeURIComponent(q)}`;
  });

  document.body.addEventListener("click", (e) => {
    const btn = e.target.closest("[data-action]");
    if (!btn) return;
    const action = btn.dataset.action;
    const id = btn.dataset.id || "";
    if (action === "task-run" || action === "task-retry") return toast(`已触发重试：${id}（占位）`);
    if (action === "task-rollback") return toast(`已触发回滚：${id}（占位）`);
    if (action === "task-stop") return toast(`已触发停止：${id}（占位）`);
    if (action === "task-import") return toast("导入：占位（后续支持 YAML/JSON）");
    if (action === "task-export") return toast("导出：占位");
    if (action === "task-batch-retry") return toast("批量重试：占位");
    if (action === "diag-run") return toast("生成诊断包：占位（后续接入 diagnose.sh/Agent）");
    if (action === "settings-save") return toast("设置已保存（占位）");
    if (action === "settings-test") return toast("已发送测试通知（占位）");
    if (action === "task-save") return toast("已保存版本（占位）");
    if (action === "task-validate") return toast("预检通过（占位）");
    if (action === "task-submit") return toast("已提交运行（占位）");
    return toast(`${action}：占位`);
  });

  window.addEventListener("hashchange", render);
}

document.addEventListener("DOMContentLoaded", () => {
  const chip = qs("#dataModeChip");
  if (chip) chip.textContent = state.dataMode === "prototype" ? "原型数据" : "实时数据";
  state.lastRefreshAt = nowTime();
  bindGlobalEvents();
  render();
});

