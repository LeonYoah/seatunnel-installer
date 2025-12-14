<!--
  任务管理页面
-->
<template>
  <div class="tasks">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>任务管理</span>
          <div class="header-actions">
            <el-button type="primary" :icon="Plus" @click="handleCreate">新建任务</el-button>
          </div>
        </div>
      </template>

      <!-- 筛选栏 -->
      <div class="filter-bar">
        <el-input
          v-model="searchText"
          placeholder="搜索任务名称"
          :prefix-icon="Search"
          style="width: 300px"
          clearable
        />
        <el-select v-model="filterStatus" placeholder="状态筛选" style="width: 150px" clearable>
          <el-option label="全部" value="" />
          <el-option label="运行中" value="running" />
          <el-option label="成功" value="success" />
          <el-option label="失败" value="failed" />
          <el-option label="暂停" value="paused" />
        </el-select>
        <el-select v-model="filterType" placeholder="类型筛选" style="width: 150px" clearable>
          <el-option label="全部" value="" />
          <el-option label="实时" value="实时" />
          <el-option label="批处理" value="批处理" />
        </el-select>
        <el-button :icon="Refresh" @click="handleRefresh">刷新</el-button>
      </div>

      <!-- 任务列表 -->
      <el-table :data="filteredTasks" style="width: 100%; margin-top: 20px">
        <el-table-column prop="name" label="任务名称" min-width="180">
          <template #default="{ row }">
            <el-link type="primary" @click="handleView(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="100" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastRun" label="最近运行" width="180" />
        <el-table-column prop="duration" label="耗时" width="120" />
        <el-table-column prop="version" label="版本" width="100" />
        <el-table-column prop="creator" label="创建者" width="120" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleRun(row)">运行</el-button>
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleView(row)">详情</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
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
    running: '运行中',
    success: '成功',
    failed: '失败',
    paused: '暂停'
  }
  return textMap[status] || status
}

const handleCreate = () => {
  ElMessage.info('新建任务功能开发中')
}

const handleRun = (row: any) => {
  ElMessage.success(`已触发运行：${row.name}`)
}

const handleEdit = (row: any) => {
  ElMessage.info(`编辑任务：${row.name}`)
}

const handleView = (row: any) => {
  ElMessage.info(`查看任务详情：${row.name}`)
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm(`确定要删除任务 "${row.name}" 吗？`, '确认删除', {
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

const handleRefresh = () => {
  ElMessage.success('刷新成功')
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
