<!--
  主机管理页面
  管理所有主机节点，包括注册、Agent 安装、状态监控
-->
<template>
  <div class="hosts">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.hosts') }}</span>
          <div class="header-actions">
            <el-button type="primary" :icon="Plus" @click="handleAddHost">{{ t('hosts.register') }}</el-button>
          </div>
        </div>
      </template>

      <!-- 筛选栏 -->
      <div class="filter-bar">
        <el-input
          v-model="searchText"
          :placeholder="t('hosts.searchPlaceholder')"
          :prefix-icon="Search"
          style="width: 300px"
          clearable
        />
        <el-select v-model="filterStatus" :placeholder="t('common.statusFilter')" style="width: 150px" clearable>
          <el-option :label="t('common.all')" value="" />
          <el-option :label="t('status.online')" value="online" />
          <el-option :label="t('status.offline')" value="offline" />
          <el-option :label="t('hosts.agentNotInstalled')" value="no-agent" />
        </el-select>
        <el-button :icon="Refresh" @click="handleRefresh">{{ t('common.refresh') }}</el-button>
      </div>

      <!-- 主机列表 -->
      <el-table :data="filteredHosts" style="width: 100%; margin-top: 20px">
        <el-table-column prop="name" :label="t('hosts.columns.name')" min-width="150" />
        <el-table-column prop="ip" :label="t('hosts.columns.ip')" width="150" />
        <el-table-column prop="port" :label="t('hosts.columns.port')" width="100" />
        <el-table-column prop="user" :label="t('hosts.columns.user')" width="120" />
        <el-table-column prop="agentStatus" :label="t('hosts.columns.agentStatus')" width="120">
          <template #default="{ row }">
            <el-tag :type="getAgentStatusType(row.agentStatus)" size="small">
              {{ getAgentStatusText(row.agentStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="t('hosts.columns.hostStatus')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'online' ? 'success' : 'info'" size="small">
              {{ row.status === 'online' ? t('status.online') : t('status.offline') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cpu" :label="t('hosts.columns.cpu')" width="120">
          <template #default="{ row }">
            <el-progress
              :percentage="row.cpu"
              :color="getProgressColor(row.cpu)"
              :show-text="true"
            />
          </template>
        </el-table-column>
        <el-table-column prop="memory" :label="t('hosts.columns.memory')" width="120">
          <template #default="{ row }">
            <el-progress
              :percentage="row.memory"
              :color="getProgressColor(row.memory)"
              :show-text="true"
            />
          </template>
        </el-table-column>
        <el-table-column prop="lastHeartbeat" :label="t('hosts.columns.lastHeartbeat')" width="180" />
        <el-table-column :label="t('common.actions')" width="300" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.agentStatus === 'not-installed'"
              size="small"
              type="primary"
              @click="handleInstallAgent(row)"
            >
              {{ t('hosts.installAgent') }}
            </el-button>
            <el-button
              v-if="row.agentStatus === 'installed'"
              size="small"
              type="warning"
              @click="handleUninstallAgent(row)"
            >
              {{ t('hosts.uninstallAgent') }}
            </el-button>
            <el-button size="small" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
            <el-button size="small" @click="handleTestConnection(row)">{{ t('hosts.testConnection') }}</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="totalHosts"
          layout="total, sizes, prev, pager, next, jumper"
        />
      </div>
    </el-card>

    <!-- 添加/编辑主机对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form :model="hostForm" :rules="hostRules" ref="hostFormRef" label-width="120px">
        <el-form-item :label="t('hosts.form.name')" prop="name">
          <el-input v-model="hostForm.name" :placeholder="t('hosts.form.exampleName')" />
        </el-form-item>
        <el-form-item :label="t('hosts.form.ip')" prop="ip">
          <el-input v-model="hostForm.ip" :placeholder="t('hosts.form.exampleIp')" />
        </el-form-item>
        <el-form-item :label="t('hosts.form.port')" prop="port">
          <el-input-number v-model="hostForm.port" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item :label="t('hosts.form.user')" prop="user">
          <el-input v-model="hostForm.user" :placeholder="t('hosts.form.exampleUser')" />
        </el-form-item>
        <el-form-item :label="t('hosts.form.authType')" prop="authType">
          <el-radio-group v-model="hostForm.authType">
            <el-radio label="password">{{ t('hosts.form.password') }}</el-radio>
            <el-radio label="key">{{ t('hosts.form.key') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="hostForm.authType === 'password'" :label="t('hosts.form.sshPassword')" prop="password">
          <el-input v-model="hostForm.password" type="password" show-password />
        </el-form-item>
        <el-form-item v-if="hostForm.authType === 'key'" :label="t('hosts.form.keyPath')" prop="keyPath">
          <el-input v-model="hostForm.keyPath" :placeholder="t('hosts.form.exampleKeyPath')" />
        </el-form-item>
        <el-form-item :label="t('hosts.form.desc')">
          <el-input
            v-model="hostForm.description"
            type="textarea"
            :rows="3"
            :placeholder="t('hosts.form.descPlaceholder')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Search, Refresh } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

// 搜索和筛选
const searchText = ref('')
const filterStatus = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref(t('hosts.register'))
const hostFormRef = ref<FormInstance>()

// 主机表单
const hostForm = ref({
  name: '',
  ip: '',
  port: 22,
  user: 'root',
  authType: 'password',
  password: '',
  keyPath: '',
  description: ''
})

// 表单验证规则
const hostRules: FormRules = {
  name: [{ required: true, message: t('hosts.valid.name'), trigger: 'blur' }],
  ip: [
    { required: true, message: t('hosts.valid.ipRequired'), trigger: 'blur' },
    {
      pattern: /^(\d{1,3}\.){3}\d{1,3}$/,
      message: t('hosts.valid.ipFormat'),
      trigger: 'blur'
    }
  ],
  port: [{ required: true, message: t('hosts.valid.port'), trigger: 'blur' }],
  user: [{ required: true, message: t('hosts.valid.user'), trigger: 'blur' }],
  password: [
    {
      required: true,
      message: t('hosts.valid.password'),
      trigger: 'blur',
      validator: (_rule, _value, callback) => {
        if (hostForm.value.authType === 'password' && !hostForm.value.password) {
          callback(new Error(t('hosts.valid.password')))
        } else {
          callback()
        }
      }
    }
  ],
  keyPath: [
    {
      required: true,
      message: t('hosts.valid.keyPath'),
      trigger: 'blur',
      validator: (_rule, _value, callback) => {
        if (hostForm.value.authType === 'key' && !hostForm.value.keyPath) {
          callback(new Error(t('hosts.valid.keyPath')))
        } else {
          callback()
        }
      }
    }
  ]
}

// 主机数据（Mock）
const hosts = ref([
  {
    id: 1,
    name: 'master-01',
    ip: '192.168.1.100',
    port: 22,
    user: 'root',
    agentStatus: 'installed',
    status: 'online',
    cpu: 45,
    memory: 62,
    lastHeartbeat: '2025-12-14 10:30:15'
  },
  {
    id: 2,
    name: 'worker-01',
    ip: '192.168.1.101',
    port: 22,
    user: 'root',
    agentStatus: 'installed',
    status: 'online',
    cpu: 78,
    memory: 85,
    lastHeartbeat: '2025-12-14 10:30:12'
  },
  {
    id: 3,
    name: 'worker-02',
    ip: '192.168.1.102',
    port: 22,
    user: 'root',
    agentStatus: 'not-installed',
    status: 'offline',
    cpu: 0,
    memory: 0,
    lastHeartbeat: '-'
  },
  {
    id: 4,
    name: 'worker-03',
    ip: '192.168.1.103',
    port: 22,
    user: 'root',
    agentStatus: 'installed',
    status: 'offline',
    cpu: 0,
    memory: 0,
    lastHeartbeat: '2025-12-14 09:15:30'
  }
])

// 过滤后的主机
const filteredHosts = computed(() => {
  let result = hosts.value

  if (searchText.value) {
    result = result.filter(
      host =>
        host.name.toLowerCase().includes(searchText.value.toLowerCase()) ||
        host.ip.includes(searchText.value)
    )
  }

  if (filterStatus.value) {
    if (filterStatus.value === 'no-agent') {
      result = result.filter(host => host.agentStatus === 'not-installed')
    } else {
      result = result.filter(host => host.status === filterStatus.value)
    }
  }

  return result
})

const totalHosts = computed(() => filteredHosts.value.length)

const getAgentStatusType = (status: string) => {
  const typeMap: Record<string, any> = {
    installed: 'success',
    'not-installed': 'warning'
  }
  return typeMap[status] || 'info'
}

const getAgentStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    installed: t('status.installed'),
    'not-installed': t('status.notInstalled')
  }
  return textMap[status] || status
}

const getProgressColor = (percentage: number) => {
  if (percentage < 60) return '#67c23a'
  if (percentage < 80) return '#e6a23c'
  return '#f56c6c'
}

const handleAddHost = () => {
  dialogTitle.value = t('hosts.register')
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = t('hosts.edit')
  hostForm.value = { ...row }
  dialogVisible.value = true
}

const handleDialogClose = () => {
  hostFormRef.value?.resetFields()
  hostForm.value = {
    name: '',
    ip: '',
    port: 22,
    user: 'root',
    authType: 'password',
    password: '',
    keyPath: '',
    description: ''
  }
}

const handleSubmit = async () => {
  if (!hostFormRef.value) return

  await hostFormRef.value.validate((valid) => {
    if (valid) {
      ElMessage.success(t('hosts.msg.saved'))
      dialogVisible.value = false
    }
  })
}

const handleInstallAgent = (row: any) => {
  ElMessageBox.confirm(
    t('hosts.msg.confirmInstallAgent', { name: row.name, ip: row.ip }),
    t('hosts.installAgent'),
    {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'info'
    }
  )
    .then(() => {
      ElMessage.success(t('hosts.msg.installSubmitted'))
      // TODO: 调用 API 安装 Agent
    })
    .catch(() => {
      ElMessage.info(t('common.cancelled'))
    })
}

const handleUninstallAgent = (row: any) => {
  ElMessageBox.confirm(
    t('hosts.msg.confirmUninstallAgent', { name: row.name, ip: row.ip }),
    t('hosts.uninstallAgent'),
    {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    }
  )
    .then(() => {
      ElMessage.success(t('hosts.msg.uninstallSuccess'))
      // TODO: 调用 API 卸载 Agent
    })
    .catch(() => {
      ElMessage.info(t('common.cancelled'))
    })
}

const handleTestConnection = (row: any) => {
  ElMessage.info(t('hosts.msg.testing', { name: row.name }))
  // TODO: 调用 API 测试连接
  setTimeout(() => {
    ElMessage.success(t('hosts.msg.testSuccess'))
  }, 1000)
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm(t('hosts.msg.confirmDelete', { name: row.name }), t('common.confirmDelete'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning'
  })
    .then(() => {
      ElMessage.success(t('common.deleteSuccess'))
      // TODO: 调用 API 删除主机
    })
    .catch(() => {
      ElMessage.info(t('common.cancelled'))
    })
}

const handleRefresh = () => {
  ElMessage.success(t('common.refreshSuccess'))
  // TODO: 调用 API 刷新主机列表
}
</script>

<style lang="scss" scoped>
.hosts {
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
