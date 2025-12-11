/**
 * SeaTunnel 安装向导 v2.0
 * 状态和日志分离，支持阶段切换
 */

const PHASES = {
    2: { name: '环境检查', steps: [1, 2, 3, 4, 5, 6, 7] },
    3: { name: '安装部署', steps: [8, 9, 10, 11, 12] },
    4: { name: '分发启动', steps: [13, 14, 15, 16] }
};

let currentPhase = 1;
let phaseStatus = {};
let pollInterval = null;
let logInterval = null;
let configData = {};

const $ = id => document.getElementById(id);

// API
const api = {
    call: async (action, params = {}) => {
        const query = new URLSearchParams({ action, ...params }).toString();
        const res = await fetch(`/cgi-bin/run.sh?${query}`);
        return res.json();
    },
    post: async (action, data) => {
        const res = await fetch(`/cgi-bin/run.sh?action=${action}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        return res.json();
    }
};

// 初始化
document.addEventListener('DOMContentLoaded', async () => {
    initEvents();
    await checkTempFiles();
    loadConfig();
});

// 检测临时文件
async function checkTempFiles() {
    try {
        const res = await api.call('check_temp');
        if (res.status === 'found') {
            const clean = confirm(`检测到上次安装遗留的临时文件 (${res.files})，是否清理？\n\n点击"确定"清理，点击"取消"保留。`);
            if (clean) {
                await api.call('clean_temp');
            }
        }
    } catch (e) {
        console.log('检测临时文件失败');
    }
}

// 绑定事件
function initEvents() {
    // 导航点击
    document.querySelectorAll('.step-nav').forEach(el => {
        el.addEventListener('click', () => goToPhase(parseInt(el.dataset.phase)));
    });
    
    // 部署模式切换
    const deployModeEl = $('deployMode');
    if (deployModeEl) deployModeEl.addEventListener('change', toggleDeployMode);
    
    // 按钮 - 安全绑定
    const bindClick = (id, fn) => { const el = $(id); if (el) el.addEventListener('click', fn); };
    bindClick('btnSaveStart', saveConfigAndStart);
    bindClick('btnRetry2', () => retryPhase(2));
    bindClick('btnRetry3', () => retryPhase(3));
    bindClick('btnRetry4', () => retryPhase(4));
    bindClick('btnNext2', () => { goToPhase(3); startPhase(3); });
    bindClick('btnNext3', () => { goToPhase(4); startPhase(4); });
    bindClick('btnNext4', () => goToPhase(5));
    bindClick('btnRestart', () => location.reload());
    bindClick('btnRefreshLog', refreshLog);
}

// 加载配置
async function loadConfig() {
    try {
        const data = await api.call('config_load');
        if (data.status === 'ok' && data.config) {
            configData = data.config;
            Object.entries(data.config).forEach(([key, value]) => {
                const el = document.querySelector(`[name="${key}"]`);
                if (el) el.value = value;
            });
            toggleDeployMode();
        }
    } catch (e) {
        console.log('加载配置失败');
    }
}

// 切换部署模式
function toggleDeployMode() {
    const mode = $('deployMode').value;
    $('separatedConfig').classList.toggle('hidden', mode !== 'separated');
    $('hybridConfig').classList.toggle('hidden', mode !== 'hybrid');
}

// 保存配置并开始
async function saveConfigAndStart() {
    const form = $('configForm');
    const formData = new FormData(form);
    configData = {};
    
    for (const [key, value] of formData.entries()) {
        if (value.trim()) configData[key] = value.trim();
    }
    
    if (!configData.BASE_DIR) {
        alert('请填写安装目录');
        return;
    }
    
    const btn = $('btnSaveStart');
    btn.disabled = true;
    btn.textContent = '保存中...';
    
    try {
        const res = await api.post('config_save', configData);
        if (res.status === 'ok') {
            setNavStatus(1, 'completed');
            goToPhase(2);
            startPhase(2);
        } else {
            alert('保存失败: ' + (res.message || ''));
        }
    } catch (e) {
        alert('保存出错');
    }
    
    btn.disabled = false;
    btn.textContent = '保存配置并开始安装';
}

// 切换阶段
function goToPhase(phase) {
    currentPhase = phase;
    
    // 更新导航
    document.querySelectorAll('.step-nav').forEach(el => {
        el.classList.toggle('active', parseInt(el.dataset.phase) === phase);
    });
    
    // 显示对应面板
    for (let i = 1; i <= 5; i++) {
        const el = $(`phase${i}`);
        if (el) el.classList.toggle('hidden', i !== phase);
    }
    
    // 完成页显示结果，停止轮询
    if (phase === 5) {
        stopPolling();
        showResult();
    }
}

// 开始执行阶段
function startPhase(phase) {
    const info = PHASES[phase];
    if (!info) return;
    
    const start = info.steps[0];
    const end = info.steps[info.steps.length - 1];
    
    // 重置状态
    phaseStatus[phase] = 'running';
    setNavStatus(phase, 'running');
    setBadge(phase, 'running');
    info.steps.forEach(s => setStepStatus(s, 'pending'));
    $(`btnNext${phase}`).disabled = true;
    
    // 立即开始轮询，不等待 run_range 响应
    startPolling(phase);
    
    // 异步发起执行请求（不阻塞）
    fetch(`/cgi-bin/run.sh?action=run_range&start=${start}&end=${end}`)
        .catch(e => console.error('启动阶段失败:', e));
}

// 重试阶段
function retryPhase(phase) {
    const info = PHASES[phase];
    if (info) info.steps.forEach(s => setStepIcon(s, '⏳'));
    startPhase(phase);
}

// 状态轮询 + 日志自动刷新
function startPolling(phase) {
    if (pollInterval) clearInterval(pollInterval);
    if (logInterval) clearInterval(logInterval);
    
    // 启动日志自动刷新
    refreshLog();
    logInterval = setInterval(refreshLog, 2000);
    
    const poll = async () => {
        try {
            const data = await api.call('status');
            const info = PHASES[phase];
            
            // 使用后端返回的步骤状态
            if (data.steps) {
                info.steps.forEach(s => {
                    const st = data.steps[s];
                    if (st === 'completed') setStepStatus(s, 'completed');
                    else if (st === 'running') setStepStatus(s, 'running');
                    else if (st === 'failed') setStepStatus(s, 'failed');
                    else if (st === 'pending') setStepStatus(s, 'pending');
                });
            }
            
            const lastStep = info.steps[info.steps.length - 1];
            
            // 检查当前阶段是否完成
            const phaseCompleted = data.steps && data.steps[lastStep] === 'completed';
            const phaseFailed = data.steps && info.steps.some(s => data.steps[s] === 'failed');
            
            if (phaseCompleted) {
                // 阶段完成，更新状态但不停止轮询（继续刷新日志）
                phaseStatus[phase] = 'completed';
                setNavStatus(phase, 'completed');
                setBadge(phase, 'completed');
                const nextBtn = $(`btnNext${phase}`);
                if (nextBtn) nextBtn.disabled = false;
            } else if (phaseFailed) {
                // 阶段失败，停止轮询
                stopPolling();
                phaseStatus[phase] = 'failed';
                setNavStatus(phase, 'failed');
                setBadge(phase, 'failed');
                refreshLog(); // 最后刷新一次日志
            }
        } catch (e) {
            console.error('轮询失败:', e);
        }
    };
    
    poll();
    pollInterval = setInterval(poll, 2000);
}

// 停止轮询和日志刷新
function stopPolling() {
    if (pollInterval) { clearInterval(pollInterval); pollInterval = null; }
    if (logInterval) { clearInterval(logInterval); logInterval = null; }
}

// 设置导航状态图标
function setNavStatus(phase, status) {
    const el = $(`navStatus${phase}`);
    if (!el) return;
    const icons = { running: '🔄', completed: '✅', failed: '❌' };
    el.textContent = icons[status] || '';
}

// 设置阶段徽章
function setBadge(phase, status) {
    const el = $(`badge${phase}`);
    if (!el) return;
    const texts = { running: '执行中', completed: '已完成', failed: '失败' };
    el.textContent = texts[status] || '';
    el.className = `status-badge ${status}`;
}

// 设置步骤图标 (旧方法，保留兼容)
function setStepIcon(step, icon) {
    const el = document.querySelector(`.step-row[data-step="${step}"] .step-icon`);
    if (el) el.textContent = icon;
}

// 设置步骤状态 (新方法，显示图标+重试按钮)
function setStepStatus(step, status) {
    const row = document.querySelector(`.step-row[data-step="${step}"]`);
    if (!row) return;
    
    const iconEl = row.querySelector('.step-icon');
    const icons = { pending: '⏳', running: '🔄', completed: '✅', failed: '❌' };
    if (iconEl) iconEl.textContent = icons[status] || '⏳';
    
    // 更新行样式
    row.className = `step-row status-${status}`;
    
    // 失败时显示重试按钮
    let retryBtn = row.querySelector('.btn-retry-step');
    if (status === 'failed') {
        if (!retryBtn) {
            retryBtn = document.createElement('button');
            retryBtn.className = 'btn-retry-step';
            retryBtn.textContent = '重试';
            retryBtn.onclick = () => retryStep(step);
            row.appendChild(retryBtn);
        }
    } else if (retryBtn) {
        retryBtn.remove();
    }
}

// 重试单个步骤
function retryStep(step) {
    setStepStatus(step, 'running');
    
    // 找到该步骤所属阶段并开始轮询
    for (const [phase, info] of Object.entries(PHASES)) {
        if (info.steps.includes(step)) {
            startPolling(parseInt(phase));
            break;
        }
    }
    
    // 异步发起请求（不阻塞）
    fetch(`/cgi-bin/run.sh?action=run_step&step=${step}`)
        .catch(e => console.error('重试步骤失败:', e));
}

// 刷新日志 - 高亮 ERROR/WARN/SUCCESS
async function refreshLog() {
    const el = $('logContent');
    el.innerHTML = '<span style="color:#64748b">加载中...</span>';
    
    try {
        const data = await api.call('log', { lines: 200 });
        if (data.log) {
            let text = data.log
                .replace(/\\n/g, '\n')
                .replace(/\x1b\[[0-9;]*m/g, '')
                .replace(/\[0;[0-9]+m/g, '')
                .replace(/\[0m/g, '')
                .replace(/\[1;[0-9]+m/g, '');
            
            // 按行处理，高亮不同级别
            const lines = text.split('\n').map(line => {
                const escaped = line.replace(/</g, '&lt;').replace(/>/g, '&gt;');
                if (line.includes('[ERROR]')) return `<span class="log-error">${escaped}</span>`;
                if (line.includes('[WARN]')) return `<span class="log-warn">${escaped}</span>`;
                if (line.includes('[SUCCESS]')) return `<span class="log-success">${escaped}</span>`;
                if (line.includes('[DEBUG]')) return `<span class="log-debug">${escaped}</span>`;
                return escaped;
            });
            el.innerHTML = lines.join('\n') || '暂无日志';
        } else {
            el.innerHTML = '暂无日志';
        }
    } catch (e) {
        el.innerHTML = '<span class="log-error">加载失败</span>';
    }
    
    el.scrollTop = el.scrollHeight;
}

// 显示结果
function showResult() {
    $('resultDetails').innerHTML = `
        <p><strong>安装路径:</strong></p>
        <code>${configData.BASE_DIR || '/home/seatunnel/seatunnel-package'}/apache-seatunnel-${configData.SEATUNNEL_VERSION || '2.3.12'}</code>
        <p><strong>启动命令:</strong></p>
        <code>systemctl start seatunnel-master</code>
        <code>systemctl start seatunnel-worker</code>
        <p><strong>查看状态:</strong></p>
        <code>systemctl status seatunnel-master</code>
        <code>systemctl status seatunnel-worker</code>
    `;
}
