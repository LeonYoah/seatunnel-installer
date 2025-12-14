<!--
  插件市场页面
-->
<template>
  <div class="plugins">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>插件市场</span>
          <div class="header-actions">
            <el-button :icon="Refresh" @click="handleRefresh">刷新仓库</el-button>
            <el-button :icon="Upload" @click="handleUpload">本地上传</el-button>
          </div>
        </div>
      </template>

      <!-- 筛选栏 -->
      <div class="filter-bar">
        <el-input
          v-model="searchText"
          placeholder="搜索插件名称"
          :prefix-icon="Search"
          style="width: 300px"
          clearable
        />
        <el-select v-model="filterType" placeholder="类型筛选" style="width: 150px" clearable>
          <el-option label="全部" value="" />
          <el-option label="Source" value="Source" />
          <el-option label="Sink" value="Sink" />
          <el-option label="Transform" value="Transform" />
        </el-select>
        <el-select v-model="filterStatus" placeholder="状态筛选" style="width: 150px" clearable>
          <el-option label="全部" value="" />
          <el-option label="已安装" value="installed" />
          <el-option label="可用" value="available" />
        </el-select>
      </div>

      <!-- 插件列表 -->
      <el-table :data="filteredPlugins" style="width: 100%; margin-top: 20px">
        <el-table-column prop="name" label="插件名称" min-width="180">
          <template #default="{ row }">
            <div class="plugin-name">
              <el-icon :size="20" style="margin-right: 8px">
                <Box />
              </el-icon>
              {{ row.name }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getTypeColor(row.type)" size="small">
              {{ row.type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="version" label="版本" width="120" />
        <el-table-column prop="compatibility" label="兼容性" width="120" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.status === 'installed' ? 'success' : 'info'" size="small">
              {{ row.status === 'installed' ? '已安装' : '可用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <template v-if="row.status === 'installed'">
              <el-button size="small" @click="handleUpgrade(row)">升级</el-button>
              <el-button size="small" type="danger" @click="handleDisable(row)">禁用</el-button>
            </template>
            <template v-else>
              <el-button size="small" type="primary" @click="handleInstall(row)">
                安装
              </el-button>
            </template>
            <el-button size="small" @click="handleViewDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Upload, Search, Box } from '@element-plus/icons-vue'

const searchText = ref('')
const filterType = ref('')
const filterStatus = ref('')

// 插件数据
const plugins = ref([
  {
    id: 1,
    name: 'MySQL-CDC',
    type: 'Source',
    version: '2.3.12',
    compatibility: '2.3.x',
    description: 'MySQL 变更数据捕获连接器',
    status: 'installed'
  },
  {
    id: 2,
    name: 'Kafka',
    type: 'Sink',
    version: '2.3.12',
    compatibility: '2.3.x',
    description: 'Apache Kafka 连接器',
    status: 'installed'
  },
  {
    id: 3,
    name: 'Iceberg',
    type: 'Sink',
    version: '2.3.12',
    compatibility: '2.3.x',
    description: 'Apache Iceberg 表格式连接器',
    status: 'available'
  },
  {
    id: 4,
    name: 'Doris',
    type: 'Sink',
    version: '2.3.12',
    compatibility: '2.3.x',
    description: 'Apache Doris 连接器',
    status: 'available'
  },
  {
    id: 5,
    name: 'PostgreSQL-CDC',
    type: 'Source',
    version: '2.3.12',
    compatibility: '2.3.x',
    description: 'PostgreSQL 变更数据捕获连接器',
    status: 'installed'
  },
  {
    id: 6,
    name: 'Elasticsearch',
    type: 'Sink',
    version: '2.3.12',
    compatibility: '2.3.x',
    description: 'Elasticsearch 连接器',
    status: 'available'
  },
  {
    id: 7,
    name: 'Hive',
    type: 'Sink',
    version: '2.3.12',
    compatibility: '2.3.x',
    description: 'Apache Hive 连接器',
    status: 'installed'
  },
  {
    id: 8,
    name: 'SQL-Transform',
    type: 'Transform',
    version: '2.3.12',
    compatibility: '2.3.x',
    description: 'SQL 转换器',
    status: 'installed'
  },
  {
    id: 9,
    name: 'MongoDB',
    type: 'Source',
    version: '2.3.12',
    compatibility: '2.3.x',
    description: 'MongoDB 连接器',
    status: 'available'
  },
  {
    id: 10,
    name: 'Redis',
    type: 'Sink',
    version: '2.3.12',
    compatibility: '2.3.x',
    description: 'Redis 连接器',
    status: 'available'
  },
  {
    id: 11,
    name: 'ClickHouse',
    type: 'Sink',
    version: '2.3.12',
    compatibility: '2.3.x',
    description: 'ClickHouse 连接器',
    status: 'available'
  },
  {
    id: 12,
    name: 'S3',
    type: 'Sink',
    version: '2.3.12',
    compatibility: '2.3.x',
    description: 'Amazon S3 连接器',
    status: 'installed'
  }
])

// 过滤后的插件
const filteredPlugins = computed(() => {
  let result = plugins.value

  if (searchText.value) {
    result = result.filter(plugin =>
      plugin.name.toLowerCase().includes(searchText.value.toLowerCase())
    )
  }

  if (filterType.value) {
    result = result.filter(plugin => plugin.type === filterType.value)
  }

  if (filterStatus.value) {
    result = result.filter(plugin => plugin.status === filterStatus.value)
  }

  return result
})

const getTypeColor = (type: string) => {
  const colorMap: Record<string, any> = {
    Source: 'primary',
    Sink: 'success',
    Transform: 'warning'
  }
  return colorMap[type] || 'info'
}

const handleRefresh = () => {
  ElMessage.success('刷新成功')
}

const handleUpload = () => {
  ElMessage.info('本地上传功能开发中')
}

const handleInstall = (row: any) => {
  ElMessageBox.confirm(`确定要安装插件 "${row.name}" 吗？`, '确认安装', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  })
    .then(() => {
      ElMessage.success(`正在安装 ${row.name}...`)
    })
    .catch(() => {
      ElMessage.info('已取消安装')
    })
}

const handleUpgrade = (row: any) => {
  ElMessage.info(`升级插件：${row.name}`)
}

const handleDisable = (row: any) => {
  ElMessageBox.confirm(`确定要禁用插件 "${row.name}" 吗？`, '确认禁用', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(() => {
      ElMessage.success('禁用成功')
    })
    .catch(() => {
      ElMessage.info('已取消禁用')
    })
}

const handleViewDetail = (row: any) => {
  ElMessage.info(`查看插件详情：${row.name}`)
}
</script>

<style scoped>
.plugins {
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

.plugin-name {
  display: flex;
  align-items: center;
  font-weight: 500;
}
</style>
