<!--
  集群管理页面
-->
<template>
  <div class="clusters">
    <!-- 集群概览 -->
    <el-row :gutter="20" class="overview-row">
      <el-col :xs="24" :sm="8">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ t('clusters.overview.totalClusters') }}</div>
            <div class="stat-value">{{ clusterStats.total }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="8">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ t('clusters.overview.totalNodes') }}</div>
            <div class="stat-value">{{ clusterStats.nodes }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="8">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ t('clusters.overview.avgCpu') }}</div>
            <div class="stat-value">{{ clusterStats.avgCpu }}%</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 集群列表 -->
    <el-card class="cluster-list">
      <template #header>
        <div class="card-header">
          <span>{{ t('clusters.listTitle') }}</span>
          <el-button type="primary" :icon="Plus" @click="handleAddCluster">
            {{ t('clusters.register') }}
          </el-button>
        </div>
      </template>

      <el-collapse v-model="activeCluster" accordion>
        <el-collapse-item
          v-for="cluster in clusters"
          :key="cluster.id"
          :name="cluster.id"
          :title="cluster.name"
        >
          <template #title>
            <div class="cluster-title">
              <span class="cluster-name">{{ cluster.name }}</span>
              <el-tag :type="cluster.status === 'healthy' ? 'success' : 'danger'" size="small">
                {{ cluster.status === 'healthy' ? t('status.healthy') : t('status.unhealthy') }}
              </el-tag>
              <span class="cluster-info">
                {{ t('clusters.version') }}: {{ cluster.version }} | {{ t('clusters.nodes') }}: {{ cluster.nodes.length }}
              </span>
            </div>
          </template>

          <!-- 节点列表 -->
          <el-table :data="cluster.nodes" style="width: 100%">
            <el-table-column prop="name" :label="t('clusters.columns.nodeName')" width="180" />
            <el-table-column prop="role" :label="t('clusters.columns.role')" width="120">
              <template #default="{ row }">
                <el-tag :type="row.role === 'master' ? 'primary' : 'info'" size="small">
                  {{ row.role === 'master' ? 'Master' : 'Worker' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="ip" :label="t('clusters.columns.ip')" width="150" />
            <el-table-column prop="cpu" label="CPU" width="100">
              <template #default="{ row }">
                <el-progress
                  :percentage="row.cpu"
                  :color="getProgressColor(row.cpu)"
                  :show-text="false"
                />
                <span style="margin-left: 8px">{{ row.cpu }}%</span>
              </template>
            </el-table-column>
            <el-table-column prop="memory" :label="t('clusters.columns.memory')" width="100">
              <template #default="{ row }">
                <el-progress
                  :percentage="row.memory"
                  :color="getProgressColor(row.memory)"
                  :show-text="false"
                />
                <span style="margin-left: 8px">{{ row.memory }}%</span>
              </template>
            </el-table-column>
            <el-table-column prop="disk" :label="t('clusters.columns.disk')" width="100">
              <template #default="{ row }">
                <el-progress
                  :percentage="row.disk"
                  :color="getProgressColor(row.disk)"
                  :show-text="false"
                />
                <span style="margin-left: 8px">{{ row.disk }}%</span>
              </template>
            </el-table-column>
            <el-table-column prop="status" :label="t('clusters.columns.status')" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'online' ? 'success' : 'danger'" size="small">
                  {{ row.status === 'online' ? t('status.online') : t('status.offline') }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="lastHeartbeat" :label="t('clusters.columns.lastHeartbeat')" width="150" />
            <el-table-column :label="t('common.actions')" width="200" fixed="right">
              <template #default="{ row }">
                <el-button size="small" text @click="handleRestart(row)">{{ t('clusters.restart') }}</el-button>
                <el-button size="small" text @click="handleDiagnose(row)">{{ t('clusters.diagnose') }}</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="cluster-actions">
            <el-button :icon="Setting" @click="handleConfigCluster(cluster)">
              {{ t('clusters.config') }}
            </el-button>
            <el-button :icon="Plus" @click="handleScaleOut(cluster)">{{ t('clusters.scaleOut') }}</el-button>
            <el-button :icon="Minus" @click="handleScaleIn(cluster)">{{ t('clusters.scaleIn') }}</el-button>
            <el-button type="danger" text :icon="Delete" @click="handleDeleteCluster(cluster)">
              {{ t('clusters.delete') }}
            </el-button>
          </div>
        </el-collapse-item>
      </el-collapse>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Minus, Setting, Delete } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const activeCluster = ref(1)

// 集群统计
const clusterStats = ref({
  total: 3,
  nodes: 9,
  avgCpu: 42
})

// 集群数据
const clusters = ref([
  {
    id: 1,
    name: 'production-cluster',
    version: '2.3.12',
    status: 'healthy',
    nodes: [
      {
        id: 1,
        name: 'master-01',
        role: 'master',
        ip: '192.168.1.100',
        cpu: 35,
        memory: 48,
        disk: 52,
        status: 'online',
        lastHeartbeat: '5秒前'
      },
      {
        id: 2,
        name: 'worker-01',
        role: 'worker',
        ip: '192.168.1.101',
        cpu: 58,
        memory: 72,
        disk: 65,
        status: 'online',
        lastHeartbeat: '6秒前'
      },
      {
        id: 3,
        name: 'worker-02',
        role: 'worker',
        ip: '192.168.1.102',
        cpu: 42,
        memory: 55,
        disk: 48,
        status: 'online',
        lastHeartbeat: '8秒前'
      }
    ]
  },
  {
    id: 2,
    name: 'test-cluster',
    version: '2.3.11',
    status: 'healthy',
    nodes: [
      {
        id: 4,
        name: 'master-01',
        role: 'master',
        ip: '192.168.2.100',
        cpu: 28,
        memory: 35,
        disk: 40,
        status: 'online',
        lastHeartbeat: '4秒前'
      },
      {
        id: 5,
        name: 'worker-01',
        role: 'worker',
        ip: '192.168.2.101',
        cpu: 45,
        memory: 60,
        disk: 55,
        status: 'online',
        lastHeartbeat: '7秒前'
      }
    ]
  },
  {
    id: 3,
    name: 'dev-cluster',
    version: '2.3.12',
    status: 'healthy',
    nodes: [
      {
        id: 6,
        name: 'master-01',
        role: 'master',
        ip: '192.168.3.100',
        cpu: 22,
        memory: 30,
        disk: 35,
        status: 'online',
        lastHeartbeat: '5秒前'
      },
      {
        id: 7,
        name: 'worker-01',
        role: 'worker',
        ip: '192.168.3.101',
        cpu: 38,
        memory: 45,
        disk: 42,
        status: 'online',
        lastHeartbeat: '6秒前'
      },
      {
        id: 8,
        name: 'worker-02',
        role: 'worker',
        ip: '192.168.3.102',
        cpu: 50,
        memory: 68,
        disk: 58,
        status: 'online',
        lastHeartbeat: '9秒前'
      },
      {
        id: 9,
        name: 'worker-03',
        role: 'worker',
        ip: '192.168.3.103',
        cpu: 15,
        memory: 25,
        disk: 30,
        status: 'online',
        lastHeartbeat: '5秒前'
      }
    ]
  }
])

const getProgressColor = (percentage: number) => {
  // 暗黑模式下统一主色强调
  if (document.documentElement.classList.contains('dark')) {
    return 'var(--primary)'
  }
  if (percentage < 60) return '#67c23a'
  if (percentage < 80) return '#e6a23c'
  return '#f56c6c'
}

const handleAddCluster = () => {
  ElMessage.info(t('clusters.msg.registerWip'))
}

const handleRestart = (node: any) => {
  ElMessage.success(t('clusters.msg.restart', { name: node.name }))
}

const handleDiagnose = (node: any) => {
  ElMessage.info(t('clusters.msg.diagnose', { name: node.name }))
}

const handleConfigCluster = (cluster: any) => {
  ElMessage.info(t('clusters.msg.config', { name: cluster.name }))
}

const handleScaleOut = (cluster: any) => {
  ElMessage.info(t('clusters.msg.scaleOut', { name: cluster.name }))
}

const handleScaleIn = (cluster: any) => {
  ElMessage.info(t('clusters.msg.scaleIn', { name: cluster.name }))
}

const handleDeleteCluster = (cluster: any) => {
  ElMessageBox.confirm(t('clusters.msg.confirmDelete', { name: cluster.name }), t('common.confirmDelete'), {
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
.clusters {
  width: 100%;
}

.overview-row {
  margin-bottom: 20px;
}

.stat-item {
  text-align: center;
  padding: 20px 0;
}

.stat-label {
  font-size: 14px;
  color: var(--muted);
  margin-bottom: 8px;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: var(--text);
}

.cluster-list {
  margin-top: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.cluster-title {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.cluster-name {
  font-weight: 600;
  font-size: 16px;
}

.cluster-info {
  color: var(--muted);
  font-size: 14px;
  margin-left: auto;
}

.cluster-actions {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--border);
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}
</style>
