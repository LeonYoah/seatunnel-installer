<!--
  诊断中心页面
-->
<template>
  <div class="diagnostics">
    <el-row :gutter="20">
      <!-- 一键诊断 -->
      <el-col :xs="24" :md="12">
        <el-card>
          <template #header>
            <span>{{ t('diagnostics.quick.title') }}</span>
          </template>
          <el-form :model="diagForm" label-width="100px">
            <el-form-item :label="t('diagnostics.quick.scope')">
              <el-select v-model="diagForm.scope" :placeholder="t('diagnostics.quick.select')" style="width: 100%">
                <el-option :label="t('diagnostics.scope.cluster')" value="cluster" />
                <el-option :label="t('diagnostics.scope.node')" value="node" />
                <el-option :label="t('diagnostics.scope.task')" value="task" />
              </el-select>
            </el-form-item>
            <el-form-item :label="t('diagnostics.quick.target')">
              <el-input v-model="diagForm.target" :placeholder="t('diagnostics.quick.targetPlaceholder')" />
            </el-form-item>
            <el-form-item :label="t('diagnostics.quick.items')">
              <el-checkbox-group v-model="diagForm.items">
                <el-checkbox label="logs">{{ t('diagnostics.items.logs') }}</el-checkbox>
                <el-checkbox label="config">{{ t('diagnostics.items.config') }}</el-checkbox>
                <el-checkbox label="thread">{{ t('diagnostics.items.thread') }}</el-checkbox>
                <el-checkbox label="heap">{{ t('diagnostics.items.heap') }}</el-checkbox>
                <el-checkbox label="gc">{{ t('diagnostics.items.gc') }}</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :icon="Tools" @click="handleDiagnose">
                {{ t('diagnostics.quick.generate') }}
              </el-button>
              <el-button @click="handleReset">{{ t('common.cancel') }}</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- 常见故障库 -->
      <el-col :xs="24" :md="12">
        <el-card>
          <template #header>
            <span>{{ t('diagnostics.fault.title') }}</span>
          </template>
          <div class="fault-list">
            <div v-for="fault in faults" :key="fault.id" class="fault-item">
              <div class="fault-header">
                <el-tag :type="fault.level === 'ERROR' ? 'danger' : 'warning'" size="small">
                  {{ fault.level }}
                </el-tag>
                <span class="fault-title">{{ fault.title }}</span>
              </div>
              <div class="fault-pattern">{{ t('diagnostics.fault.pattern') }}：{{ fault.pattern }}</div>
              <div class="fault-actions">
                <el-button size="small" @click="handleViewSolution(fault)">{{ t('diagnostics.fault.view') }}</el-button>
                <el-button
                  v-if="fault.fixable"
                  size="small"
                  type="primary"
                  @click="handleAutoFix(fault)"
                >
                  {{ t('diagnostics.fault.autoFix') }}
                </el-button>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 诊断历史 -->
    <el-card class="history-card">
      <template #header>
        <div class="card-header">
          <span>{{ t('diagnostics.history.title') }}</span>
          <el-button :icon="Refresh" @click="handleRefresh">{{ t('common.refresh') }}</el-button>
        </div>
      </template>
      <el-table :data="diagHistory" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="scope" :label="t('diagnostics.history.scope')" width="100">
          <template #default="{ row }">{{ getScopeText(row.scope) }}</template>
        </el-table-column>
        <el-table-column prop="target" :label="t('diagnostics.history.target')" width="180" />
        <el-table-column prop="items" :label="t('diagnostics.history.items')" min-width="200">
          <template #default="{ row }">
            <el-tag
              v-for="item in row.items"
              :key="item"
              size="small"
              style="margin-right: 4px"
            >
              {{ getItemLabel(item) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="t('diagnostics.history.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="size" :label="t('diagnostics.history.size')" width="100" />
        <el-table-column prop="createdAt" :label="t('diagnostics.history.createdAt')" width="180" />
        <el-table-column :label="t('common.actions')" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" :icon="Download" @click="handleDownload(row)">
              {{ t('common.download') }}
            </el-button>
            <el-button size="small" type="danger" @click="handleDeleteDiag(row)">
              {{ t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Tools, Refresh, Download } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

// 诊断表单
const diagForm = ref({
  scope: 'cluster',
  target: '',
  items: ['logs', 'config']
})

// 常见故障
const faults = ref([
  {
    id: 1,
    level: 'WARN',
    title: '背压导致延迟升高',
    pattern: 'backpressure / sink slow',
    fixable: false
  },
  {
    id: 2,
    level: 'ERROR',
    title: 'Connector 依赖缺失',
    pattern: 'ClassNotFoundException',
    fixable: true
  },
  {
    id: 3,
    level: 'ERROR',
    title: '内存溢出',
    pattern: 'OutOfMemoryError',
    fixable: false
  },
  {
    id: 4,
    level: 'WARN',
    title: '连接池耗尽',
    pattern: 'Connection pool exhausted',
    fixable: false
  }
])

// 诊断历史
const diagHistory = ref([
  {
    id: 1,
    scope: '集群',
    target: 'production-cluster',
    items: ['logs', 'config', 'thread'],
    status: 'completed',
    size: '125 MB',
    createdAt: '2025-12-14 10:30:15'
  },
  {
    id: 2,
    scope: '节点',
    target: 'worker-01',
    items: ['logs', 'heap', 'gc'],
    status: 'completed',
    size: '256 MB',
    createdAt: '2025-12-14 09:15:22'
  },
  {
    id: 3,
    scope: '任务',
    target: 'orders_sync',
    items: ['logs', 'config'],
    status: 'failed',
    size: '-',
    createdAt: '2025-12-14 08:45:10'
  }
])

const getItemLabel = (item: string) => {
  const labels: Record<string, string> = {
    logs: t('diagnostics.items.logs'),
    config: t('diagnostics.items.config'),
    thread: t('diagnostics.items.thread'),
    heap: t('diagnostics.items.heap'),
    gc: t('diagnostics.items.gc')
  }
  return labels[item] || item
}

const getScopeText = (scope: string) => {
  const map: Record<string, string> = {
    集群: t('diagnostics.scope.cluster'),
    节点: t('diagnostics.scope.node'),
    任务: t('diagnostics.scope.task')
  }
  return map[scope] || scope
}

const getStatusType = (status: string) => {
  const typeMap: Record<string, any> = {
    completed: 'success',
    running: 'primary',
    failed: 'danger'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    completed: t('status.completed'),
    running: t('status.processing'),
    failed: t('status.failed')
  }
  return textMap[status] || status
}

const handleDiagnose = () => {
  if (!diagForm.value.target) {
    ElMessage.warning(t('diagnostics.msg.targetRequired'))
    return
  }
  if (diagForm.value.items.length === 0) {
    ElMessage.warning(t('diagnostics.msg.itemsRequired'))
    return
  }
  ElMessage.success(t('diagnostics.msg.created'))
}

const handleReset = () => {
  diagForm.value = {
    scope: 'cluster',
    target: '',
    items: ['logs', 'config']
  }
}

const handleViewSolution = (fault: any) => {
  ElMessage.info(t('diagnostics.fault.viewSolution', { title: fault.title }))
}

const handleAutoFix = (fault: any) => {
  ElMessageBox.confirm(t('diagnostics.fault.confirmAutoFix', { title: fault.title }), t('diagnostics.fault.autoFix'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning'
  })
    .then(() => {
      ElMessage.success(t('diagnostics.fault.fixed'))
    })
    .catch(() => {
      ElMessage.info(t('common.cancelled'))
    })
}

const handleRefresh = () => {
  ElMessage.success(t('common.refreshSuccess'))
}

const handleDownload = (row: any) => {
  ElMessage.success(t('diagnostics.history.downloading', { id: row.id }))
}

const handleDeleteDiag = (row: any) => {
  ElMessageBox.confirm(t('diagnostics.history.confirmDelete', { id: row.id }), t('common.confirmDelete'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning'
  })
    .then(() => {
      ElMessage.success(t('common.deleteSuccess'))
    })
    .catch(() => {
      ElMessage.info(t('common.cancelled'))
    })
}
</script>

<style scoped>
.diagnostics {
  width: 100%;
}

.fault-list {
  max-height: 400px;
  overflow-y: auto;
}

.fault-item {
  padding: 16px;
  border: 1px solid var(--border);
  border-radius: 8px;
  margin-bottom: 12px;
}

.fault-item:last-child {
  margin-bottom: 0;
}

.fault-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.fault-title {
  font-weight: 600;
  color: var(--text);
}

.fault-pattern {
  font-size: 13px;
  color: var(--muted);
  margin-bottom: 12px;
}

.fault-actions {
  display: flex;
  gap: 8px;
}

.history-card {
  margin-top: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
