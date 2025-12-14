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
            <div class="stat-label">集群总数</div>
            <div class="stat-value">{{ clusterStats.total }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="8">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">节点总数</div>
            <div class="stat-value">{{ clusterStats.nodes }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="8">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">平均 CPU 使用率</div>
            <div class="stat-value">{{ clusterStats.avgCpu }}%</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 集群列表 -->
    <el-card class="cluster-list">
      <template #header>
        <div class="card-header">
          <span>集群列表</span>
          <el-button type="primary" :icon="Plus" @click="handleAddCluster">
            注册集群
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
                {{ cluster.status === 'healthy' ? '健康' : '异常' }}
              </el-tag>
              <span class="cluster-info">
                版本: {{ cluster.version }} | 节点: {{ cluster.nodes.length }}
              </span>
            </div>
          </template>

          <!-- 节点列表 -->
          <el-table :data="cluster.nodes" style="width: 100%">
            <el-table-column prop="name" label="节点名称" width="180" />
            <el-table-column prop="role" label="角色" width="120">
              <template #default="{ row }">
                <el-tag :type="row.role === 'master' ? 'primary' : 'info'" size="small">
                  {{ row.role === 'master' ? 'Master' : 'Worker' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="ip" label="IP 地址" width="150" />
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
            <el-table-column prop="memory" label="内存" width="100">
              <template #default="{ row }">
                <el-progress
                  :percentage="row.memory"
                  :color="getProgressColor(row.memory)"
                  :show-text="false"
                />
                <span style="margin-left: 8px">{{ row.memory }}%</span>
              </template>
            </el-table-column>
            <el-table-column prop="disk" label="磁盘" width="100">
              <template #default="{ row }">
                <el-progress
                  :percentage="row.disk"
                  :color="getProgressColor(row.disk)"
                  :show-text="false"
                />
                <span style="margin-left: 8px">{{ row.disk }}%</span>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'online' ? 'success' : 'danger'" size="small">
                  {{ row.status === 'online' ? '在线' : '离线' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="lastHeartbeat" label="最后心跳" width="150" />
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="handleRestart(row)">重启</el-button>
                <el-button size="small" @click="handleDiagnose(row)">诊断</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="cluster-actions">
            <el-button :icon="Setting" @click="handleConfigCluster(cluster)">
              配置管理
            </el-button>
            <el-button :icon="Plus" @click="handleScaleOut(cluster)">扩容</el-button>
            <el-button :icon="Minus" @click="handleScaleIn(cluster)">缩容</el-button>
            <el-button type="danger" :icon="Delete" @click="handleDeleteCluster(cluster)">
              删除集群
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
  if (percentage < 60) return '#67c23a'
  if (percentage < 80) return '#e6a23c'
  return '#f56c6c'
}

const handleAddCluster = () => {
  ElMessage.info('注册集群功能开发中')
}

const handleRestart = (node: any) => {
  ElMessage.success(`已触发重启：${node.name}`)
}

const handleDiagnose = (node: any) => {
  ElMessage.info(`诊断节点：${node.name}`)
}

const handleConfigCluster = (cluster: any) => {
  ElMessage.info(`配置集群：${cluster.name}`)
}

const handleScaleOut = (cluster: any) => {
  ElMessage.info(`扩容集群：${cluster.name}`)
}

const handleScaleIn = (cluster: any) => {
  ElMessage.info(`缩容集群：${cluster.name}`)
}

const handleDeleteCluster = (cluster: any) => {
  ElMessageBox.confirm(`确定要删除集群 "${cluster.name}" 吗？`, '确认删除', {
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
