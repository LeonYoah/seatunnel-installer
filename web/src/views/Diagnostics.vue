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
            <span>一键诊断</span>
          </template>
          <el-form :model="diagForm" label-width="100px">
            <el-form-item label="诊断范围">
              <el-select v-model="diagForm.scope" placeholder="请选择" style="width: 100%">
                <el-option label="整个集群" value="cluster" />
                <el-option label="指定节点" value="node" />
                <el-option label="指定任务" value="task" />
              </el-select>
            </el-form-item>
            <el-form-item label="目标对象">
              <el-input v-model="diagForm.target" placeholder="如 worker-01 / orders_sync" />
            </el-form-item>
            <el-form-item label="收集内容">
              <el-checkbox-group v-model="diagForm.items">
                <el-checkbox label="logs">日志文件</el-checkbox>
                <el-checkbox label="config">配置快照</el-checkbox>
                <el-checkbox label="thread">线程 Dump</el-checkbox>
                <el-checkbox label="heap">堆内存 Dump</el-checkbox>
                <el-checkbox label="gc">GC 日志</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :icon="Tools" @click="handleDiagnose">
                生成诊断包
              </el-button>
              <el-button @click="handleReset">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- 常见故障库 -->
      <el-col :xs="24" :md="12">
        <el-card>
          <template #header>
            <span>常见故障库</span>
          </template>
          <div class="fault-list">
            <div v-for="fault in faults" :key="fault.id" class="fault-item">
              <div class="fault-header">
                <el-tag :type="fault.level === 'ERROR' ? 'danger' : 'warning'" size="small">
                  {{ fault.level }}
                </el-tag>
                <span class="fault-title">{{ fault.title }}</span>
              </div>
              <div class="fault-pattern">特征：{{ fault.pattern }}</div>
              <div class="fault-actions">
                <el-button size="small" @click="handleViewSolution(fault)">查看方案</el-button>
                <el-button
                  v-if="fault.fixable"
                  size="small"
                  type="primary"
                  @click="handleAutoFix(fault)"
                >
                  一键修复
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
          <span>诊断历史</span>
          <el-button :icon="Refresh" @click="handleRefresh">刷新</el-button>
        </div>
      </template>
      <el-table :data="diagHistory" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="scope" label="范围" width="100" />
        <el-table-column prop="target" label="目标" width="180" />
        <el-table-column prop="items" label="收集内容" min-width="200">
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
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="大小" width="100" />
        <el-table-column prop="createdAt" label="创建时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" :icon="Download" @click="handleDownload(row)">
              下载
            </el-button>
            <el-button size="small" type="danger" @click="handleDeleteDiag(row)">
              删除
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
    logs: '日志',
    config: '配置',
    thread: '线程',
    heap: '堆内存',
    gc: 'GC'
  }
  return labels[item] || item
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
    completed: '已完成',
    running: '进行中',
    failed: '失败'
  }
  return textMap[status] || status
}

const handleDiagnose = () => {
  if (!diagForm.value.target) {
    ElMessage.warning('请输入目标对象')
    return
  }
  if (diagForm.value.items.length === 0) {
    ElMessage.warning('请选择至少一项收集内容')
    return
  }
  ElMessage.success('诊断任务已创建，正在收集数据...')
}

const handleReset = () => {
  diagForm.value = {
    scope: 'cluster',
    target: '',
    items: ['logs', 'config']
  }
}

const handleViewSolution = (fault: any) => {
  ElMessage.info(`查看故障解决方案：${fault.title}`)
}

const handleAutoFix = (fault: any) => {
  ElMessageBox.confirm(`确定要自动修复 "${fault.title}" 吗？`, '确认修复', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(() => {
      ElMessage.success('修复成功')
    })
    .catch(() => {
      ElMessage.info('已取消修复')
    })
}

const handleRefresh = () => {
  ElMessage.success('刷新成功')
}

const handleDownload = (row: any) => {
  ElMessage.success(`开始下载诊断包 #${row.id}`)
}

const handleDeleteDiag = (row: any) => {
  ElMessageBox.confirm(`确定要删除诊断包 #${row.id} 吗？`, '确认删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(() => {
      ElMessage.success('删除成功')
    })
    .catch(() => {
      ElMessage.info('已取消删除')
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
