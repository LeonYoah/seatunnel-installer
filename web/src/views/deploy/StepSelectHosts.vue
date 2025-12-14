<!--
  步骤 1: 选择主机
  从已注册的主机中选择要部署 SeaTunnel 的节点
-->
<template>
  <el-card class="step-card">
    <template #header>
      <span>选择主机</span>
    </template>

    <el-alert
      title="提示"
      type="info"
      :closable="false"
      style="margin-bottom: 20px"
    >
      请从已注册的主机中选择要部署 SeaTunnel 集群的节点。如果没有可用主机，请先前往
      <el-link type="primary" @click="goToHosts">主机管理</el-link>
      注册主机并安装 Agent。
    </el-alert>

    <el-form :model="form" label-width="140px">
      <!-- 部署模式 -->
      <div class="form-section">
        <h3>部署模式</h3>
        <el-form-item label="部署模式">
          <el-radio-group v-model="form.deployMode">
            <el-radio label="separated">分离模式（Master/Worker）</el-radio>
            <el-radio label="hybrid">混合模式（所有节点角色相同）</el-radio>
          </el-radio-group>
        </el-form-item>
      </div>

      <!-- 分离模式 - 选择 Master -->
      <div v-if="form.deployMode === 'separated'" class="form-section">
        <h3>Master 节点</h3>
        <el-form-item label="选择 Master">
          <el-select
            v-model="form.masterHost"
            placeholder="请选择 Master 节点"
            style="width: 400px"
            filterable
          >
            <el-option
              v-for="host in availableHosts"
              :key="host.id"
              :label="`${host.name} (${host.ip})`"
              :value="host.id"
              :disabled="!host.agentInstalled"
            >
              <div style="display: flex; justify-content: space-between; align-items: center">
                <span>{{ host.name }} ({{ host.ip }})</span>
                <el-tag v-if="!host.agentInstalled" type="warning" size="small">
                  Agent 未安装
                </el-tag>
                <el-tag v-else-if="host.status === 'offline'" type="info" size="small">
                  离线
                </el-tag>
                <el-tag v-else type="success" size="small">在线</el-tag>
              </div>
            </el-option>
          </el-select>
        </el-form-item>
      </div>

      <!-- 分离模式 - 选择 Workers -->
      <div v-if="form.deployMode === 'separated'" class="form-section">
        <h3>Worker 节点</h3>
        <el-form-item label="选择 Workers">
          <el-select
            v-model="form.workerHosts"
            placeholder="请选择 Worker 节点（可多选）"
            style="width: 400px"
            multiple
            filterable
          >
            <el-option
              v-for="host in availableHosts"
              :key="host.id"
              :label="`${host.name} (${host.ip})`"
              :value="host.id"
              :disabled="!host.agentInstalled || host.id === form.masterHost"
            >
              <div style="display: flex; justify-content: space-between; align-items: center">
                <span>{{ host.name }} ({{ host.ip }})</span>
                <el-tag v-if="!host.agentInstalled" type="warning" size="small">
                  Agent 未安装
                </el-tag>
                <el-tag v-else-if="host.status === 'offline'" type="info" size="small">
                  离线
                </el-tag>
                <el-tag v-else type="success" size="small">在线</el-tag>
              </div>
            </el-option>
          </el-select>
        </el-form-item>
      </div>

      <!-- 混合模式 - 选择所有节点 -->
      <div v-if="form.deployMode === 'hybrid'" class="form-section">
        <h3>集群节点</h3>
        <el-form-item label="选择节点">
          <el-select
            v-model="form.clusterHosts"
            placeholder="请选择集群节点（可多选）"
            style="width: 400px"
            multiple
            filterable
          >
            <el-option
              v-for="host in availableHosts"
              :key="host.id"
              :label="`${host.name} (${host.ip})`"
              :value="host.id"
              :disabled="!host.agentInstalled"
            >
              <div style="display: flex; justify-content: space-between; align-items: center">
                <span>{{ host.name}} ({{ host.ip }})</span>
                <el-tag v-if="!host.agentInstalled" type="warning" size="small">
                  Agent 未安装
                </el-tag>
                <el-tag v-else-if="host.status === 'offline'" type="info" size="small">
                  离线
                </el-tag>
                <el-tag v-else type="success" size="small">在线</el-tag>
              </div>
            </el-option>
          </el-select>
        </el-form-item>
      </div>

      <!-- 已选主机预览 -->
      <div v-if="selectedHosts.length > 0" class="form-section">
        <h3>已选主机（{{ selectedHosts.length }} 台）</h3>
        <el-table :data="selectedHosts" style="width: 100%">
          <el-table-column prop="name" label="主机名称" width="150" />
          <el-table-column prop="ip" label="IP 地址" width="150" />
          <el-table-column prop="role" label="角色" width="120">
            <template #default="{ row }">
              <el-tag :type="row.role === 'Master' ? 'danger' : 'primary'" size="small">
                {{ row.role }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="cpu" label="CPU" width="100" />
          <el-table-column prop="memory" label="内存" width="100" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'online' ? 'success' : 'info'" size="small">
                {{ row.status === 'online' ? '在线' : '离线' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-form>

    <div class="step-actions">
      <el-button type="primary" :disabled="!canProceed" @click="handleNext">下一步</el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const emit = defineEmits(['next'])

const form = ref({
  deployMode: 'separated',
  masterHost: null as number | null,
  workerHosts: [] as number[],
  clusterHosts: [] as number[]
})

// Mock 数据 - 已注册的主机
const availableHosts = ref([
  {
    id: 1,
    name: 'master-01',
    ip: '192.168.1.100',
    agentInstalled: true,
    status: 'online',
    cpu: '8 核',
    memory: '16 GB'
  },
  {
    id: 2,
    name: 'worker-01',
    ip: '192.168.1.101',
    agentInstalled: true,
    status: 'online',
    cpu: '16 核',
    memory: '32 GB'
  },
  {
    id: 3,
    name: 'worker-02',
    ip: '192.168.1.102',
    agentInstalled: false,
    status: 'offline',
    cpu: '16 核',
    memory: '32 GB'
  },
  {
    id: 4,
    name: 'worker-03',
    ip: '192.168.1.103',
    agentInstalled: true,
    status: 'online',
    cpu: '16 核',
    memory: '32 GB'
  }
])

// 已选主机列表
const selectedHosts = computed(() => {
  const hosts = []

  if (form.value.deployMode === 'separated') {
    // Master
    if (form.value.masterHost) {
      const master = availableHosts.value.find(h => h.id === form.value.masterHost)
      if (master) {
        hosts.push({ ...master, role: 'Master' })
      }
    }
    // Workers
    form.value.workerHosts.forEach(id => {
      const worker = availableHosts.value.find(h => h.id === id)
      if (worker) {
        hosts.push({ ...worker, role: 'Worker' })
      }
    })
  } else {
    // 混合模式
    form.value.clusterHosts.forEach(id => {
      const host = availableHosts.value.find(h => h.id === id)
      if (host) {
        hosts.push({ ...host, role: 'Hybrid' })
      }
    })
  }

  return hosts
})

// 是否可以继续
const canProceed = computed(() => {
  if (form.value.deployMode === 'separated') {
    return form.value.masterHost !== null && form.value.workerHosts.length > 0
  } else {
    return form.value.clusterHosts.length > 0
  }
})

const goToHosts = () => {
  router.push('/hosts')
}

const handleNext = () => {
  emit('next')
}
</script>

<style scoped>
.step-card {
  margin-bottom: 20px;
}

.form-section {
  margin-bottom: 30px;
}

.form-section h3 {
  margin: 0 0 20px 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text);
}

.step-actions {
  margin-top: 30px;
  text-align: right;
}
</style>
