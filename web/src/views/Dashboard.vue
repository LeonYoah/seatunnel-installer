<!--
  Dashboard 总览页面
  展示集群健康状态、任务统计、告警信息等
-->
<template>
  <div class="dashboard">
    <!-- KPI 卡片 -->
    <el-row :gutter="20" class="kpi-row">
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="kpi-card">
          <div class="kpi-content">
            <div class="kpi-icon success">
              <el-icon :size="32"><Connection /></el-icon>
            </div>
            <div class="kpi-info">
              <div class="kpi-value">{{ clusterStats.healthy }}/{{ clusterStats.total }}</div>
              <div class="kpi-label">{{ t('dashboard.kpis.health') }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="kpi-card">
          <div class="kpi-content">
            <div class="kpi-icon primary">
              <el-icon :size="32"><List /></el-icon>
            </div>
            <div class="kpi-info">
              <div class="kpi-value">{{ taskStats.successRate }}%</div>
              <div class="kpi-label">{{ t('dashboard.kpis.taskSuccess') }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="kpi-card">
          <div class="kpi-content">
            <div class="kpi-icon warning">
              <el-icon :size="32"><Timer /></el-icon>
            </div>
            <div class="kpi-info">
              <div class="kpi-value">{{ taskStats.p95Latency }} ms</div>
              <div class="kpi-label">{{ t('dashboard.kpis.latencyP95') }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="kpi-card">
          <div class="kpi-content">
            <div class="kpi-icon" :class="alertStats.count > 0 ? 'error' : 'success'">
              <el-icon :size="32"><Bell /></el-icon>
            </div>
            <div class="kpi-info">
              <div class="kpi-value">{{ alertStats.count }}</div>
              <div class="kpi-label">{{ t('dashboard.kpis.activeAlerts') }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 趋势图表和告警 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :xs="24" :md="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ t('dashboard.trend.title') }}</span>
              <el-tag size="small" type="info">{{ t('dashboard.trend.last24h') }}</el-tag>
            </div>
          </template>
          <div ref="chartContainer" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="8">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ t('dashboard.alerts.title') }}</span>
              <el-link type="primary" :underline="false" @click="goToAlerts">{{ t('common.viewAll') }}</el-link>
            </div>
          </template>
          <div class="alert-list">
            <div v-for="alert in recentAlerts" :key="alert.id" class="alert-item">
              <el-tag :type="alert.level === 'ERROR' ? 'danger' : 'warning'" size="small">
                {{ alert.level }}
              </el-tag>
              <div class="alert-content">
                <div class="alert-target">{{ alert.target }}</div>
                <div class="alert-message">{{ alert.message }}</div>
                <div class="alert-time">{{ alert.time }}</div>
              </div>
            </div>
            <el-empty v-if="recentAlerts.length === 0" :description="t('dashboard.alerts.empty')" :image-size="80" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近任务 -->
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ t('dashboard.recentTasks.title') }}</span>
              <el-link type="primary" :underline="false" @click="goToTasks">{{ t('common.viewAll') }}</el-link>
            </div>
          </template>
          <el-table :data="recentTasks" style="width: 100%">
            <el-table-column prop="name" :label="t('tasks.columns.name')" min-width="180" />
            <el-table-column prop="type" :label="t('tasks.columns.type')" width="100" />
            <el-table-column prop="status" :label="t('tasks.columns.status')" width="120">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="lastRun" :label="t('tasks.columns.lastRun')" width="180" />
            <el-table-column prop="duration" :label="t('tasks.columns.duration')" width="120" />
            <el-table-column prop="version" :label="t('tasks.columns.version')" width="100" />
            <el-table-column :label="t('common.actions')" width="200" fixed="right">
              <template #default="{ row }">
                <el-button size="small" text @click="handleTaskAction(row, 'run')">{{ t('common.retry') }}</el-button>
                <el-button size="small" text @click="handleTaskAction(row, 'view')">{{ t('common.view') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Connection,
  List,
  Timer,
  Bell
} from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import * as echarts from 'echarts'

const router = useRouter()
const { t } = useI18n()

// 图表容器引用
const chartContainer = ref<HTMLElement>()
let chartInstance: echarts.ECharts | null = null

// 集群统计数据
const clusterStats = ref({
  total: 3,
  healthy: 3
})

// 任务统计数据
const taskStats = ref({
  successRate: 99.2,
  p95Latency: 820,
  throughput: 180000
})

// 告警统计数据
const alertStats = ref({
  count: 2
})

// 最近告警
const recentAlerts = ref([
  {
    id: 1,
    level: 'WARN',
    target: 'task:orders_sync',
    message: '延迟抖动，P95 > 1s',
    time: '10分钟前'
  },
  {
    id: 2,
    level: 'ERROR',
    target: 'node:worker-03',
    message: '心跳超时（已恢复）',
    time: '2小时前'
  }
])

// 最近任务
const recentTasks = ref([
  {
    id: 1,
    name: 'orders_sync',
    type: '实时',
    status: 'running',
    lastRun: '持续运行',
    duration: '-',
    version: 'v14'
  },
  {
    id: 2,
    name: 'users_dim',
    type: '批处理',
    status: 'success',
    lastRun: '2025-12-13 09:02',
    duration: '3m12s',
    version: 'v8'
  },
  {
    id: 3,
    name: 'refunds_backfill',
    type: '批处理',
    status: 'failed',
    lastRun: '2025-12-13 08:40',
    duration: '1m05s',
    version: 'v3'
  },
  {
    id: 4,
    name: 'cdc_mysql_kafka',
    type: '实时',
    status: 'paused',
    lastRun: '暂停',
    duration: '-',
    version: 'v21'
  }
])

const getStatusType = (status: string) => {
  const typeMap: Record<string, any> = {
    running: 'primary',
    success: 'success',
    failed: 'danger',
    paused: 'warning'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    running: t('status.running'),
    success: t('status.success'),
    failed: t('status.failed'),
    paused: t('status.paused')
  }
  return textMap[status] || status
}

const goToAlerts = () => {
  ElMessage.info(t('dashboard.msg.alertsWip'))
}

const goToTasks = () => {
  router.push('/tasks')
}

const handleTaskAction = (row: any, action: string) => {
  if (action === 'run') {
    ElMessage.success(t('dashboard.msg.retryTriggered', { name: row.name }))
  } else if (action === 'view') {
    ElMessage.info(t('tasks.msg.view', { name: row.name }))
  }
}

// 初始化图表
const initChart = () => {
  if (!chartContainer.value) return

  chartInstance = echarts.init(chartContainer.value)
  
  // 生成过去24小时的模拟数据
  const hours = []
  const successData = []
  const failedData = []
  const now = new Date()
  
  for (let i = 23; i >= 0; i--) {
    const time = new Date(now.getTime() - i * 60 * 60 * 1000)
    hours.push(time.getHours() + ':00')
    // 模拟成功任务数据（80-120之间波动）
    successData.push(Math.floor(Math.random() * 40) + 80)
    // 模拟失败任务数据（0-10之间波动）
    failedData.push(Math.floor(Math.random() * 10))
  }

  const option = {
    title: {
      text: t('dashboard.trend.title'),
      textStyle: {
        fontSize: 16,
        fontWeight: 'normal'
      }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      }
    },
    legend: {
      data: [t('dashboard.chart.success'), t('dashboard.chart.failed')],
      bottom: 10
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: hours,
      axisLabel: {
        fontSize: 12
      }
    },
    yAxis: {
      type: 'value',
      name: t('dashboard.chart.taskCount'),
      axisLabel: {
        fontSize: 12
      }
    },
    series: [
      {
        name: t('dashboard.chart.success'),
        type: 'line',
        smooth: true,
        data: successData,
        itemStyle: {
          color: '#67c23a'
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(103, 194, 58, 0.3)' },
            { offset: 1, color: 'rgba(103, 194, 58, 0.1)' }
          ])
        }
      },
      {
        name: t('dashboard.chart.failed'),
        type: 'line',
        smooth: true,
        data: failedData,
        itemStyle: {
          color: '#f56c6c'
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(245, 108, 108, 0.3)' },
            { offset: 1, color: 'rgba(245, 108, 108, 0.1)' }
          ])
        }
      }
    ]
  }

  chartInstance.setOption(option)
}

// 响应式调整图表大小
const handleResize = () => {
  if (chartInstance) {
    chartInstance.resize()
  }
}

onMounted(() => {
  // 延迟初始化图表，确保DOM已渲染
  setTimeout(() => {
    initChart()
  }, 100)
  
  // 监听窗口大小变化
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  // 清理图表实例和事件监听器
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.dashboard {
  width: 100%;
}

.kpi-row {
  margin-bottom: 20px;
}

.kpi-card {
  height: 100%;
}

.kpi-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.kpi-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.kpi-icon.success {
  background: rgba(103, 194, 58, 0.1);
  color: var(--success);
}

.kpi-icon.primary {
  background: rgba(64, 158, 255, 0.1);
  color: var(--primary);
}

.kpi-icon.warning {
  background: rgba(230, 162, 60, 0.1);
  color: var(--warning);
}

.kpi-icon.error {
  background: rgba(245, 108, 108, 0.1);
  color: var(--error);
}

.kpi-info {
  flex: 1;
}

.kpi-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--text);
  line-height: 1.2;
  margin-bottom: 4px;
}

.kpi-label {
  font-size: 14px;
  color: var(--muted);
}

.chart-row {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.chart-container {
  height: 300px;
  width: 100%;
}

.alert-list {
  max-height: 300px;
  overflow-y: auto;
}

.alert-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  border-bottom: 1px solid var(--border);
}

.alert-item:last-child {
  border-bottom: none;
}

.alert-content {
  flex: 1;
}

.alert-target {
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
  margin-bottom: 4px;
}

.alert-message {
  font-size: 13px;
  color: var(--text);
  margin-bottom: 4px;
}

.alert-time {
  font-size: 12px;
  color: var(--muted);
}
</style>
