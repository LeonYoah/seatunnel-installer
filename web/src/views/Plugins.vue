<!--
  插件市场页面
-->
<template>
  <div class="plugins">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.plugins') }}</span>
          <div class="header-actions">
            <el-button :icon="Refresh" @click="handleRefresh">{{ t('plugins.refreshRepo') }}</el-button>
            <el-button :icon="Upload" @click="handleUpload">{{ t('plugins.uploadLocal') }}</el-button>
          </div>
        </div>
      </template>

      <!-- 筛选栏 -->
      <div class="filter-bar">
        <el-input
          v-model="searchText"
          :placeholder="t('plugins.searchPlaceholder')"
          :prefix-icon="Search"
          style="width: 300px"
          clearable
        />
        <el-select v-model="filterType" :placeholder="t('common.typeFilter')" style="width: 150px" clearable>
          <el-option :label="t('common.all')" value="" />
          <el-option label="Source" value="Source" />
          <el-option label="Sink" value="Sink" />
          <el-option label="Transform" value="Transform" />
        </el-select>
        <el-select v-model="filterStatus" :placeholder="t('common.statusFilter')" style="width: 150px" clearable>
          <el-option :label="t('common.all')" value="" />
          <el-option :label="t('status.installed')" value="installed" />
          <el-option :label="t('plugins.available')" value="available" />
        </el-select>
      </div>

      <!-- 插件列表 -->
      <el-table :data="filteredPlugins" style="width: 100%; margin-top: 20px">
        <el-table-column prop="name" :label="t('plugins.columns.name')" min-width="180">
          <template #default="{ row }">
            <div class="plugin-name">
              <el-icon :size="20" style="margin-right: 8px">
                <Box />
              </el-icon>
              {{ row.name }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="type" :label="t('plugins.columns.type')" width="120">
          <template #default="{ row }">
            <el-tag :type="getTypeColor(row.type)" size="small">
              {{ row.type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="version" :label="t('plugins.columns.version')" width="120" />
        <el-table-column prop="compatibility" :label="t('plugins.columns.compatibility')" width="120" />
        <el-table-column prop="description" :label="t('plugins.columns.desc')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" :label="t('plugins.columns.status')" width="120">
          <template #default="{ row }">
            <el-tag :type="row.status === 'installed' ? 'success' : 'info'" size="small">
              {{ row.status === 'installed' ? t('status.installed') : t('plugins.available') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="220" fixed="right">
          <template #default="{ row }">
            <template v-if="row.status === 'installed'">
              <el-button size="small" text @click="handleUpgrade(row)">{{ t('plugins.upgrade') }}</el-button>
              <el-button size="small" type="danger" text @click="handleDisable(row)">{{ t('plugins.disable') }}</el-button>
            </template>
            <template v-else>
              <el-button size="small" type="primary" text @click="handleInstall(row)">
                {{ t('plugins.install') }}
              </el-button>
            </template>
            <el-button size="small" text @click="handleViewDetail(row)">{{ t('common.view') }}</el-button>
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
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

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
  ElMessage.success(t('common.refreshSuccess'))
}

const handleUpload = () => {
  ElMessage.info(t('plugins.msg.uploadWip'))
}

const handleInstall = (row: any) => {
  ElMessageBox.confirm(t('plugins.msg.confirmInstall', { name: row.name }), t('plugins.install'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'info'
  })
    .then(() => {
      ElMessage.success(t('plugins.msg.installing', { name: row.name }))
    })
    .catch(() => {
      ElMessage.info(t('common.cancelled'))
    })
}

const handleUpgrade = (row: any) => {
  ElMessage.info(t('plugins.msg.upgrade', { name: row.name }))
}

const handleDisable = (row: any) => {
  ElMessageBox.confirm(t('plugins.msg.confirmDisable', { name: row.name }), t('plugins.disable'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning'
  })
    .then(() => {
      ElMessage.success(t('plugins.msg.disabled'))
    })
    .catch(() => {
      ElMessage.info(t('common.cancelled'))
    })
}

const handleViewDetail = (row: any) => {
  ElMessage.info(t('plugins.msg.view', { name: row.name }))
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
