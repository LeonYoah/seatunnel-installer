/**
 * SeaTunnel 安装向导 v2.0
 * 状态和日志分离，支持阶段切换
 */

const PHASES = {
    2: { name: '环境检查', steps: [1, 2, 3, 4, 5, 6, 7] },
    3: { name: '安装部署', steps: [8, 9, 10, 11] },
    4: { name: '分发启动', steps: [12, 13, 14, 15, 16] }
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
    loadConfig();
});

// 绑定事件
function initEvents() {
    // 导航点击
    document.querySelectorAll('.step-nav').forEach(el => {
        el.addEventListener('click', () => goToPhase(parseInt(el.dataset.phase)));
    });
    
    // 部署模式切换
    const deployModeEl = $('deployMode');
    if (deployModeEl) deployModeEl.addEventListener('change', toggleDeployMode);
    
    // 安装模式切换
    const installModeEl = $('installMode');
    if (installModeEl) installModeEl.addEventListener('change', toggleInstallMode);
    
    // 检查点存储类型切换
    const checkpointTypeEl = $('checkpointType');
    if (checkpointTypeEl) checkpointTypeEl.addEventListener('change', toggleCheckpointType);
    
    // 防火墙检查开关
    const checkFirewallEl = $('checkFirewall');
    if (checkFirewallEl) checkFirewallEl.addEventListener('change', toggleFirewallAction);
    
    // 安装包/插件下载源
    const packageRepoEl = $('packageRepo');
    if (packageRepoEl) packageRepoEl.addEventListener('change', togglePackageRepo);
    const pluginRepoEl = $('pluginRepo');
    if (pluginRepoEl) pluginRepoEl.addEventListener('change', togglePluginRepo);
    
    // 连接器安装开关
    const installConnectorsEl = $('installConnectors');
    if (installConnectorsEl) installConnectorsEl.addEventListener('change', toggleConnectorsRow);
    
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
            toggleInstallMode();
            toggleCheckpointType();
            toggleFirewallAction();
            togglePackageRepo();
            togglePluginRepo();
            toggleConnectorsRow();
        }
    } catch (e) {
        console.log('加载配置失败');
    }
}

// 切换部署模式
function toggleDeployMode() {
    const modeEl = $('deployMode');
    if (!modeEl) return;
    const mode = modeEl.value;
    const isHybrid = mode === 'hybrid';

    // 节点配置
    $('separatedConfig').classList.toggle('hidden', isHybrid);
    $('hybridConfig').classList.toggle('hidden', !isHybrid);

    // 端口配置：混合模式使用 HYBRID_PORT，分离模式使用 MASTER/WORKER_PORT
    const hybridPortRow = document.getElementById('hybridPortRow');
    const masterPortRow = document.getElementById('masterPortRow');
    const workerPortRow = document.getElementById('workerPortRow');
    if (hybridPortRow && masterPortRow && workerPortRow) {
        hybridPortRow.classList.toggle('hidden', !isHybrid);
        masterPortRow.classList.toggle('hidden', isHybrid);
        workerPortRow.classList.toggle('hidden', isHybrid);
    }

    // JVM 配置：混合模式突出 HYBRID_HEAP_SIZE，分离模式突出 Master/Worker
    const hybridHeapRow = document.getElementById('hybridHeapRow');
    const masterHeapRow = document.getElementById('masterHeapRow');
    const workerHeapRow = document.getElementById('workerHeapRow');
    if (hybridHeapRow && masterHeapRow && workerHeapRow) {
        hybridHeapRow.classList.toggle('dimmed', !isHybrid);
        masterHeapRow.classList.toggle('dimmed', isHybrid);
        workerHeapRow.classList.toggle('dimmed', isHybrid);
    }
}

// 切换安装模式（online/offline）
function toggleInstallMode() {
    const modeEl = $('installMode');
    if (!modeEl) return;
    const mode = modeEl.value;
    const row = document.getElementById('offlinePackageRow');
    if (row) row.classList.toggle('hidden', mode !== 'offline');
}

// 切换检查点存储类型
function toggleCheckpointType() {
    const typeEl = $('checkpointType');
    if (!typeEl) return;
    const type = typeEl.value;
    const hdfs = document.getElementById('hdfsConfig');
    const obj = document.getElementById('objectStorageConfig');
    if (hdfs) hdfs.classList.toggle('hidden', type !== 'HDFS');
    if (obj) obj.classList.toggle('hidden', type !== 'OSS' && type !== 'S3');
}

// 防火墙检查行为联动
function toggleFirewallAction() {
    const checkEl = $('checkFirewall');
    const row = document.getElementById('firewallActionRow');
    if (!checkEl || !row) return;
    row.classList.toggle('dimmed', checkEl.value !== 'true');
}

// 安装包下载源联动
function togglePackageRepo() {
    const repoEl = $('packageRepo');
    const row = document.getElementById('customPackageRow');
    if (!repoEl || !row) return;
    row.classList.toggle('hidden', repoEl.value !== 'custom');
}

// 插件下载源联动
function togglePluginRepo() {
    const repoEl = $('pluginRepo');
    const row = document.getElementById('customPluginRow');
    if (!repoEl || !row) return;
    row.classList.toggle('hidden', repoEl.value !== 'custom');
}

// 连接器开关联动
function toggleConnectorsRow() {
    const el = $('installConnectors');
    const row = document.getElementById('connectorsRow');
    if (!el || !row) return;
    row.classList.toggle('dimmed', el.value !== 'true');
}

// 保存配置并开始
async function saveConfigAndStart() {
    const form = $('configForm');
    const formData = new FormData(form);
    configData = {};
    
    for (const [key, value] of formData.entries()) {
        if (typeof value === 'string' && value.trim()) {
            configData[key] = value.trim();
        }
    }
    
    // 基础校验
    if (!configData.BASE_DIR) {
        alert('请填写安装目录');
        return;
    }
    
    // 离线安装必须填写 PACKAGE_PATH
    if (configData.INSTALL_MODE === 'offline' && !configData.PACKAGE_PATH) {
        alert('离线安装模式下必须填写离线安装包路径 (PACKAGE_PATH)');
        return;
    }
    
    // 部署模式校验
    if (configData.DEPLOY_MODE === 'separated') {
        if (!configData.MASTER_IP || !configData.WORKER_IPS) {
            alert('分离模式下必须配置 Master IP 和 Worker IPs');
            return;
        }
    } else if (configData.DEPLOY_MODE === 'hybrid') {
        if (!configData.CLUSTER_NODES) {
            alert('混合模式下必须配置集群节点 IP 列表');
            return;
        }
    }
    
    // 检查点存储配置校验（简化）
    if (configData.CHECKPOINT_STORAGE_TYPE === 'HDFS') {
        if (!configData.HDFS_NAMENODE_HOST || !configData.HDFS_NAMENODE_PORT) {
            alert('HDFS 检查点存储需要配置 HDFS_NAMENODE_HOST 和 HDFS_NAMENODE_PORT');
            return;
        }
    }
    if (configData.CHECKPOINT_STORAGE_TYPE === 'OSS' || configData.CHECKPOINT_STORAGE_TYPE === 'S3') {
        if (!configData.STORAGE_ENDPOINT || !configData.STORAGE_ACCESS_KEY || !configData.STORAGE_SECRET_KEY || !configData.STORAGE_BUCKET) {
            alert('OSS/S3 检查点存储需要完整的 Endpoint / AccessKey / SecretKey / Bucket 配置');
            return;
        }
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
    
    // 切换到执行阶段时，重新获取状态
    if (phase >= 2 && phase <= 4) {
        refreshPhaseStatus(phase);
    }
}

// 刷新阶段状态（切换 tab 时调用）
async function refreshPhaseStatus(phase) {
    const info = PHASES[phase];
    if (!info) return;
    
    // 先初始化所有步骤为 pending（确保按钮显示）
    info.steps.forEach(s => setStepStatus(s, 'pending'));
    
    try {
        const data = await api.call('status');
        
        // 更新该阶段所有步骤的状态
        if (data.steps) {
            info.steps.forEach(s => {
                const st = data.steps[s];
                if (st) {
                    setStepStatus(s, st);
                }
            });
            
            // 检查阶段整体状态
            const lastStep = info.steps[info.steps.length - 1];
            const phaseCompleted = data.steps[lastStep] === 'completed';
            const phaseFailed = info.steps.some(s => data.steps[s] === 'failed');
            const phasePaused = info.steps.some(s => data.steps[s] === 'paused');
            
            if (phaseCompleted) {
                phaseStatus[phase] = 'completed';
                setNavStatus(phase, 'completed');
                setBadge(phase, 'completed');
                const nextBtn = $(`btnNext${phase}`);
                if (nextBtn) nextBtn.disabled = false;
            } else if (phaseFailed) {
                phaseStatus[phase] = 'failed';
                setNavStatus(phase, 'failed');
                setBadge(phase, 'failed');
            } else if (phasePaused) {
                phaseStatus[phase] = 'paused';
                setNavStatus(phase, 'paused');
                setBadge(phase, 'paused');
            }
        }
        
        // 刷新日志
        refreshLog();
    } catch (e) {
        console.error('获取状态失败:', e);
        // API 失败时，步骤已经初始化为 pending，按钮仍会显示
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
                    else if (st === 'paused') setStepStatus(s, 'paused');
                    else if (st === 'pending') setStepStatus(s, 'pending');
                });
            }
            
            const lastStep = info.steps[info.steps.length - 1];
            
            // 检查当前阶段是否完成
            const phaseCompleted = data.steps && data.steps[lastStep] === 'completed';
            const phaseFailed = data.steps && info.steps.some(s => data.steps[s] === 'failed');
            
            // 检查是否有暂停的步骤
            const phasePaused = data.steps && info.steps.some(s => data.steps[s] === 'paused');
            
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
            } else if (phasePaused) {
                // 阶段暂停，停止轮询
                stopPolling();
                phaseStatus[phase] = 'paused';
                setNavStatus(phase, 'paused');
                setBadge(phase, 'paused');
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
    const icons = { running: '●', completed: '✓', failed: '!', paused: 'Ⅱ' };
    el.textContent = icons[status] || '';
}

// 设置阶段徽章
function setBadge(phase, status) {
    const el = $(`badge${phase}`);
    if (!el) return;
    const texts = { running: '执行中', completed: '已完成', failed: '失败', paused: '已暂停' };
    el.textContent = texts[status] || '';
    el.className = `status-badge ${status}`;
}

// 设置步骤图标 (旧方法，保留兼容)
function setStepIcon(step, icon) {
    const el = document.querySelector(`.step-row[data-step="${step}"] .step-icon`);
    if (el) el.textContent = icon;
}

// 设置步骤状态 (显示图标+操作按钮)
function setStepStatus(step, status) {
    const row = document.querySelector(`.step-row[data-step="${step}"]`);
    if (!row) return;
    
    const iconEl = row.querySelector('.step-icon');
    const icons = { pending: '○', running: '●', completed: '✓', failed: '✗', paused: 'Ⅱ' };
    if (iconEl) iconEl.textContent = icons[status] || '○';
    
    // 更新行样式
    row.className = `step-row status-${status}`;
    
    // 清除旧的操作按钮容器
    let actionsEl = row.querySelector('.step-actions');
    if (!actionsEl) {
        actionsEl = document.createElement('div');
        actionsEl.className = 'step-actions';
        row.appendChild(actionsEl);
    }
    actionsEl.innerHTML = '';
    
    // 根据状态显示不同按钮
    switch (status) {
        case 'pending':
            // 待执行：显示"运行"按钮
            actionsEl.innerHTML = `
                <button class="btn-step btn-run" onclick="runSingleStep(${step})">运行</button>
                <button class="btn-step btn-continue" onclick="continueFromStep(${step})">从此继续</button>
            `;
            break;
        case 'running':
            // 执行中：显示"暂停"按钮
            actionsEl.innerHTML = `
                <button class="btn-step btn-pause" onclick="pauseExecution()">暂停</button>
            `;
            break;
        case 'completed':
            // 已完成：显示"重新运行"按钮（可选）
            actionsEl.innerHTML = `
                <button class="btn-step btn-rerun" onclick="runSingleStep(${step})">重新运行</button>
            `;
            break;
        case 'failed':
            // 失败：显示"重试"和"从此继续"按钮
            actionsEl.innerHTML = `
                <button class="btn-step btn-retry" onclick="runSingleStep(${step})">重试</button>
                <button class="btn-step btn-continue" onclick="continueFromStep(${step})">从此继续</button>
            `;
            break;
        case 'paused':
            // 已暂停：显示"继续"按钮
            actionsEl.innerHTML = `
                <button class="btn-step btn-resume" onclick="continueFromStep(${step})">继续执行</button>
            `;
            break;
    }
}

// 运行单个步骤
function runSingleStep(step) {
    setStepStatus(step, 'running');
    
    // 找到该步骤所属阶段并开始轮询
    for (const [phase, info] of Object.entries(PHASES)) {
        if (info.steps.includes(step)) {
            phaseStatus[phase] = 'running';
            setNavStatus(phase, 'running');
            setBadge(phase, 'running');
            startPolling(parseInt(phase));
            break;
        }
    }
    
    // 异步发起请求（不阻塞）
    fetch(`/cgi-bin/run.sh?action=run_step&step=${step}`)
        .catch(e => console.error('执行步骤失败:', e));
}

// 从指定步骤继续执行到阶段结束
function continueFromStep(step) {
    // 找到该步骤所属阶段
    for (const [phase, info] of Object.entries(PHASES)) {
        if (info.steps.includes(step)) {
            const endStep = info.steps[info.steps.length - 1];
            
            // 更新状态
            phaseStatus[phase] = 'running';
            setNavStatus(phase, 'running');
            setBadge(phase, 'running');
            
            // 将从当前步骤到结束的所有步骤设为 pending
            for (let s = step; s <= endStep; s++) {
                if (info.steps.includes(s)) {
                    setStepStatus(s, 'pending');
                }
            }
            setStepStatus(step, 'running');
            
            // 开始轮询
            startPolling(parseInt(phase));
            
            // 发起从当前步骤到阶段结束的执行请求
            fetch(`/cgi-bin/run.sh?action=run_range&start=${step}&end=${endStep}`)
                .catch(e => console.error('继续执行失败:', e));
            break;
        }
    }
}

// 暂停执行
function pauseExecution() {
    fetch(`/cgi-bin/run.sh?action=pause`)
        .then(res => res.json())
        .then(data => {
            if (data.status === 'ok') {
                // 更新当前运行中的步骤为 paused
                document.querySelectorAll('.step-row.status-running').forEach(row => {
                    const step = parseInt(row.dataset.step);
                    setStepStatus(step, 'paused');
                });
            }
        })
        .catch(e => console.error('暂停失败:', e));
}

// 日志缓存，用于增量更新
let lastLogHash = '';
let lastLogLineCount = 0;
let userScrolled = false;  // 用户是否手动滚动过

// 刷新日志 - 高亮 ERROR/WARN/SUCCESS，增量更新
async function refreshLog() {
    const el = $('logContent');
    
    // 首次加载显示加载中
    if (!lastLogHash) {
        el.innerHTML = '<span style="color:#64748b">加载中...</span>';
    }
    
    try {
        const data = await api.call('log', { lines: 200 });
        if (data.log) {
            let text = data.log
                .replace(/\\n/g, '\n')
                .replace(/\x1b\[[0-9;]*m/g, '')
                .replace(/\[0;[0-9]+m/g, '')
                .replace(/\[0m/g, '')
                .replace(/\[1;[0-9]+m/g, '');
            
            // 计算日志哈希（简单使用长度+最后100字符）
            const logHash = text.length + '_' + text.slice(-100);
            
            // 如果日志没有变化，不更新页面
            if (logHash === lastLogHash) {
                return;
            }
            lastLogHash = logHash;
            
            // 记录当前滚动位置
            const wasAtBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 50;
            
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
            lastLogLineCount = lines.length;
            
            // 只有当用户在底部时才自动滚动到底部
            if (wasAtBottom) {
                el.scrollTop = el.scrollHeight;
            }
        } else {
            if (lastLogHash !== 'empty') {
                el.innerHTML = '暂无日志';
                lastLogHash = 'empty';
            }
        }
    } catch (e) {
        if (lastLogHash !== 'error') {
            el.innerHTML = '<span class="log-error">加载失败</span>';
            lastLogHash = 'error';
        }
    }
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
