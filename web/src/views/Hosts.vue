<!--
  主机管理页面
  管理所有主机节点，包括注册、Agent 安装、状态监控
-->
<template>
  <div class="hosts">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>主机管理</span>
          <div class="header-actions">
            <el-button type="primary" :icon="Plus" @click="handleAddHost">注册主机</el-button>
          </div>
        </div>
      </template>

      <!-- 筛选栏 -->
      <div class="filter-bar">
        <el-input
          v-model="searchText"
          placeholder="搜索主机 IP 或名称"
          :prefix-icon="Search"
          style="width: 300px"
          clearable
        />
        <el-select v-model="filterStatus" placeholder="状态筛选" style="width: 150px" clearable>
          <el-option label="全部" value="" />
          <el-option label="在线" value="online" />
          <el-option label="离线" value="offline" />
          <el-option label="Agent 未安装" value="no-agent" />
        </el-select>
        <el-button :icon="Refresh" @click="handleRefresh">刷新</el-button>
      </div>

      <!-- 主机列表 -->
      <el-table :data="filteredHosts" style="width: 100%; margin-top: 20px">
        <el-table-column prop="name" label="主机名称" min-width="150" />
        <el-table-column prop="ip" label="IP 地址" width="150" />
        <el-table-column prop="port" label="SSH 端口" width="100" />
        <el-table-column prop="user" label="SSH 用户" width="120" />
        <el-table-column prop="agentStatus" label="Agent 状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getAgentStatusType(row.agentStatus)" size="small">
              {{ getAgentStatusText(row.agentStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="主机状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'online' ? 'success' : 'info'" size="small">
              {{ row.status === 'online' ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cpu" label="CPU 使用率" width="120">
          <template #default="{ row }">
            <el-progress
              :percentage="row.cpu"
              :color="getProgressColor(row.cpu)"
              :show-text="true"
            />
          </template>
        </el-table-column>
        <el-table-column prop="memory" label="内存使用率" width="120">
          <template #default="{ row }">
            <el-progress
              :percentage="row.memory"
              :color="getProgressColor(row.memory)"
              :show-text="true"
            />
          </template>
        </el-table-column>
        <el-table-column prop="lastHeartbeat" label="最后心跳" width="180" />
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.agentStatus === 'not-installed'"
              size="small"
              type="primary"
              @click="handleInstallAgent(row)"
            >
              安装 Agent
            </el-button>
            <el-button
              v-if="row.agentStatus === 'installed'"
              size="small"
              type="warning"
              @click="handleUninstallAgent(row)"
            >
              卸载 Agent
            </el-button>
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleTestConnection(row)">测试连接</el-button>
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
        <el-form-item label="主机名称" prop="name">
          <el-input v-model="hostForm.name" placeholder="例如：master-01" />
        </el-form-item>
        <el-form-item label="IP 地址" prop="ip">
          <el-input v-model="hostForm.ip" placeholder="例如：192.168.1.100" />
        </el-form-item>
        <el-form-item label="SSH 端口" prop="port">
          <el-input-number v-model="hostForm.port" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item label="SSH 用户" prop="user">
          <el-input v-model="hostForm.user" placeholder="例如：root" />
        </el-form-item>
        <el-form-item label="认证方式" prop="authType">
          <el-radio-group v-model="hostForm.authType">
            <el-radio label="password">密码</el-radio>
            <el-radio label="key">密钥</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="hostForm.authType === 'password'" label="SSH 密码" prop="password">
          <el-input v-model="hostForm.password" type="password" show-password />
        </el-form-item>
        <el-form-item v-if="hostForm.authType === 'key'" label="私钥路径" prop="keyPath">
          <el-input v-model="hostForm.keyPath" placeholder="例如：~/.ssh/id_rsa" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="hostForm.description"
            type="textarea"
            :rows="3"
            placeholder="主机描述信息"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Search, Refresh } from '@element-plus/icons-vue'

// 搜索和筛选
const searchText = ref('')
const filterStatus = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('注册主机')
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
  name: [{ required: true, message: '请输入主机名称', trigger: 'blur' }],
  ip: [
    { required: true, message: '请输入 IP 地址', trigger: 'blur' },
    {
      pattern: /^(\d{1,3}\.){3}\d{1,3}$/,
      message: '请输入有效的 IP 地址',
      trigger: 'blur'
    }
  ],
  port: [{ required: true, message: '请输入 SSH 端口', trigger: 'blur' }],
  user: [{ required: true, message: '请输入 SSH 用户', trigger: 'blur' }],
  password: [
    {
      required: true,
      message: '请输入 SSH 密码',
      trigger: 'blur',
      validator: (_rule, _value, callback) => {
        if (hostForm.value.authType === 'password' && !hostForm.value.password) {
          callback(new Error('请输入 SSH 密码'))
        } else {
          callback()
        }
      }
    }
  ],
  keyPath: [
    {
      required: true,
      message: '请输入私钥路径',
      trigger: 'blur',
      validator: (_rule, _value, callback) => {
        if (hostForm.value.authType === 'key' && !hostForm.value.keyPath) {
          callback(new Error('请输入私钥路径'))
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
    installed: '已安装',
    'not-installed': '未安装'
  }
  return textMap[status] || status
}

const getProgressColor = (percentage: number) => {
  if (percentage < 60) return '#67c23a'
  if (percentage < 80) return '#e6a23c'
  return '#f56c6c'
}

const handleAddHost = () => {
  dialogTitle.value = '注册主机'
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑主机'
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
      ElMessage.success('主机注册成功')
      dialogVisible.value = false
    }
  })
}

const handleInstallAgent = (row: any) => {
  ElMessageBox.confirm(
    `确定要在主机 "${row.name}" (${row.ip}) 上安装 Agent 吗？`,
    '安装 Agent',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    }
  )
    .then(() => {
      ElMessage.success('Agent 安装任务已提交，请稍候...')
      // TODO: 调用 API 安装 Agent
    })
    .catch(() => {
      ElMessage.info('已取消安装')
    })
}

const handleUninstallAgent = (row: any) => {
  ElMessageBox.confirm(
    `确定要卸载主机 "${row.name}" (${row.ip}) 上的 Agent 吗？`,
    '卸载 Agent',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  )
    .then(() => {
      ElMessage.success('Agent 卸载成功')
      // TODO: 调用 API 卸载 Agent
    })
    .catch(() => {
      ElMessage.info('已取消卸载')
    })
}

const handleTestConnection = (row: any) => {
  ElMessage.info(`正在测试与主机 ${row.name} 的连接...`)
  // TODO: 调用 API 测试连接
  setTimeout(() => {
    ElMessage.success('连接测试成功')
  }, 1000)
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm(`确定要删除主机 "${row.name}" 吗？`, '确认删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(() => {
      ElMessage.success('删除成功')
      // TODO: 调用 API 删除主机
    })
    .catch(() => {
      ElMessage.info('已取消删除')
    })
}

const handleRefresh = () => {
  ElMessage.success('刷新成功')
  // TODO: 调用 API 刷新主机列表
}
</script>

<style scoped>
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
