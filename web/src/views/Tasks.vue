<!--
  任务管理页面
-->
<template>
  <div class="tasks">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.tasks') }}</span>
          <div class="header-actions">
            <el-button type="primary" :icon="Plus" @click="handleCreate">{{ t('common.create') }}</el-button>
          </div>
        </div>
      </template>

      <!-- 筛选栏 -->
      <div class="filter-bar">
        <el-input
          v-model="searchText"
          :placeholder="t('tasks.searchPlaceholder')"
          :prefix-icon="Search"
          style="width: 300px"
          clearable
        />
        <el-select v-model="filterStatus" :placeholder="t('common.statusFilter')" style="width: 150px" clearable>
          <el-option :label="t('common.all')" value="" />
          <el-option :label="t('status.running')" value="running" />
          <el-option :label="t('status.success')" value="success" />
          <el-option :label="t('status.failed')" value="failed" />
          <el-option :label="t('status.paused')" value="paused" />
        </el-select>
        <el-select v-model="filterType" :placeholder="t('common.typeFilter')" style="width: 150px" clearable>
          <el-option :label="t('common.all')" value="" />
          <el-option :label="t('tasks.type.streaming')" value="实时" />
          <el-option :label="t('tasks.type.batch')" value="批处理" />
        </el-select>
        <el-button :icon="Refresh" @click="handleRefresh">{{ t('common.refresh') }}</el-button>
      </div>

      <!-- 任务列表 -->
      <el-table :data="filteredTasks" style="width: 100%; margin-top: 20px">
        <el-table-column prop="name" :label="t('tasks.columns.name')" min-width="180">
          <template #default="{ row }">
            <el-link type="primary" @click="handleView(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="type" :label="t('tasks.columns.type')" width="100">
          <template #default="{ row }">{{ getTypeText(row.type) }}</template>
        </el-table-column>
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
        <el-table-column prop="creator" :label="t('tasks.columns.creator')" width="120" />
        <el-table-column :label="t('common.actions')" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" text @click="handleRun(row)">{{ t('common.run') }}</el-button>
            <el-button size="small" text @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
            <el-button size="small" text @click="handleView(row)">{{ t('common.view') }}</el-button>
            <el-button size="small" type="danger" text @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="totalTasks"
          layout="total, sizes, prev, pager, next, jumper"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

// 搜索和筛选
const searchText = ref('')
const filterStatus = ref('')
const filterType = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

// 任务数据
const tasks = ref([
  {
    id: 1,
    name: 'orders_sync',
    type: '实时',
    status: 'running',
    lastRun: '持续运行',
    duration: '-',
    version: 'v14',
    creator: '张三'
  },
  {
    id: 2,
    name: 'users_dim',
    type: '批处理',
    status: 'success',
    lastRun: '2025-12-14 09:02',
    duration: '3m12s',
    version: 'v8',
    creator: '李四'
  },
  {
    id: 3,
    name: 'refunds_backfill',
    type: '批处理',
    status: 'failed',
    lastRun: '2025-12-14 08:40',
    duration: '1m05s',
    version: 'v3',
    creator: '王五'
  },
  {
    id: 4,
    name: 'cdc_mysql_kafka',
    type: '实时',
    status: 'paused',
    lastRun: '暂停',
    duration: '-',
    version: 'v21',
    creator: '张三'
  },
  {
    id: 5,
    name: 'products_etl',
    type: '批处理',
    status: 'success',
    lastRun: '2025-12-14 10:15',
    duration: '5m30s',
    version: 'v12',
    creator: '李四'
  },
  {
    id: 6,
    name: 'logs_aggregation',
    type: '实时',
    status: 'running',
    lastRun: '持续运行',
    duration: '-',
    version: 'v6',
    creator: '王五'
  },
  {
    id: 7,
    name: 'inventory_sync',
    type: '批处理',
    status: 'success',
    lastRun: '2025-12-14 09:45',
    duration: '2m18s',
    version: 'v9',
    creator: '张三'
  },
  {
    id: 8,
    name: 'customer_profile',
    type: '批处理',
    status: 'failed',
    lastRun: '2025-12-14 08:20',
    duration: '45s',
    version: 'v4',
    creator: '李四'
  }
])

// 过滤后的任务
const filteredTasks = computed(() => {
  let result = tasks.value

  if (searchText.value) {
    result = result.filter(task =>
      task.name.toLowerCase().includes(searchText.value.toLowerCase())
    )
  }

  if (filterStatus.value) {
    result = result.filter(task => task.status === filterStatus.value)
  }

  if (filterType.value) {
    result = result.filter(task => task.type === filterType.value)
  }

  return result
})

const totalTasks = computed(() => filteredTasks.value.length)

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

const getTypeText = (type: string) => {
  const textMap: Record<string, string> = {
    '实时': t('tasks.type.streaming'),
    '批处理': t('tasks.type.batch')
  }
  return textMap[type] || type
}

const handleCreate = () => {
  ElMessage.info(t('tasks.msg.createWip'))
}

const handleRun = (row: any) => {
  ElMessage.success(t('tasks.msg.runTriggered', { name: row.name }))
}

const handleEdit = (row: any) => {
  ElMessage.info(t('tasks.msg.edit', { name: row.name }))
}

const handleView = (row: any) => {
  ElMessage.info(t('tasks.msg.view', { name: row.name }))
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm(t('tasks.msg.confirmDelete', { name: row.name }), t('common.confirmDelete'), {
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

const handleRefresh = () => {
  ElMessage.success(t('common.refreshSuccess'))
}
</script>

<style scoped>
.tasks {
  width: 100%;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.filter-bar {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
