<!--
  主机管理页面
  管理所有主机节点，包括注册、Agent 安装、状态监控
-->
<template>
  <div class="hosts">
    <!-- 页面导航标签 -->
    <PageTabs :tabs="pageTabs" :active-key="activeTab" @tab-change="handleTabChange">
      <template #actions>
        <el-button :icon="Refresh" @click="handleRefresh">{{ t('common.refresh') }}</el-button>
        <el-button type="primary" :icon="Plus" @click="handleAddHost">{{ t('hosts.register') }}</el-button>
      </template>
    </PageTabs>

    <el-card>
      <!-- 工具栏 -->
      <div class="toolbar">
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
      </div>

      <!-- 主机列表 -->
      <el-table 
        :data="filteredHosts" 
        :row-class-name="getRowClassName"
        class="hosts-table"
        style="width: 100%; margin-top: 16px"
      >
        <!-- 展开列 -->
        <el-table-column type="expand" width="50">
          <template #default="{ row }">
            <div class="expand-content">
              <!-- 基本信息 -->
              <div class="info-section">
                <h4 class="section-title">{{ t('hosts.expandSections.basicInfo') }}</h4>
                <el-descriptions :column="3" border size="small">
                  <el-descriptions-item :label="t('hosts.details.cpuModel')">{{ row.cpuModel }}</el-descriptions-item>
                  <el-descriptions-item :label="t('hosts.details.totalMemory')">{{ row.totalMemory }}</el-descriptions-item>
                  <el-descriptions-item :label="t('hosts.details.registerTime')">{{ row.registerTime }}</el-descriptions-item>
                  <el-descriptions-item :label="t('hosts.details.os')">{{ row.os }}</el-descriptions-item>
                  <el-descriptions-item :label="t('hosts.details.kernel')">{{ row.kernel }}</el-descriptions-item>
                  <el-descriptions-item :label="t('hosts.details.hostname')">{{ row.hostname }}</el-descriptions-item>
                </el-descriptions>
              </div>

              <!-- 网卡信息 -->
              <div class="info-section">
                <h4 class="section-title">{{ t('hosts.expandSections.networkInfo') }}</h4>
                <el-descriptions :column="2" border size="small">
                  <el-descriptions-item 
                    v-for="(nic, index) in row.networkInterfaces" 
                    :key="index"
                    :label="nic.name"
                  >
                    {{ nic.ip }} / {{ nic.speed }}
                  </el-descriptions-item>
                </el-descriptions>
              </div>

              <!-- 磁盘信息 -->
              <div class="info-section">
                <h4 class="section-title">{{ t('hosts.expandSections.diskInfo') }}</h4>
                <el-table :data="row.disks" size="small" border class="disk-table">
                  <el-table-column prop="name" :label="t('hosts.disk.name')" width="120" />
                  <el-table-column prop="used" :label="t('hosts.disk.used')" width="120" />
                  <el-table-column prop="total" :label="t('hosts.disk.total')" width="120" />
                  <el-table-column :label="t('hosts.disk.usage')" width="200">
                    <template #default="{ row: disk }">
                      <el-progress
                        :percentage="disk.usage"
                        :color="getProgressColor(disk.usage)"
                        :show-text="true"
                      />
                    </template>
                  </el-table-column>
                  <el-table-column prop="ioUtil" :label="t('hosts.disk.ioUtil')" width="100" />
                  <el-table-column prop="mountPoint" :label="t('hosts.disk.mountPoint')" min-width="150" />
                </el-table>
              </div>
            </div>
          </template>
        </el-table-column>

        <!-- 主要列 -->
        <el-table-column prop="ip" :label="t('hosts.columns.ipAddress')" min-width="140" />
        <el-table-column prop="agentStatus" :label="t('hosts.columns.agentStatus')" min-width="110">
          <template #default="{ row }">
            <el-tag :type="getAgentStatusType(row.agentStatus)" size="small">
              {{ getAgentStatusText(row.agentStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="t('hosts.columns.hostStatus')" min-width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'online' ? 'success' : 'info'" size="small">
              {{ row.status === 'online' ? t('status.online') : t('status.offline') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('hosts.columns.cpuUsage')" min-width="150">
          <template #default="{ row }">
            <el-progress
              :percentage="row.cpu"
              :color="getProgressColor(row.cpu)"
              :show-text="true"
              :stroke-width="8"
            />
          </template>
        </el-table-column>
        <el-table-column :label="t('hosts.columns.memoryUsage')" min-width="150">
          <template #default="{ row }">
            <el-progress
              :percentage="row.memory"
              :color="getProgressColor(row.memory)"
              :show-text="true"
              :stroke-width="8"
            />
          </template>
        </el-table-column>
        <el-table-column prop="networkIO" :label="t('hosts.columns.networkIO')" min-width="130">
          <template #default="{ row }">
            <div class="network-io">
              <div class="io-item">↑ {{ row.networkIO.upload }}</div>
              <div class="io-item">↓ {{ row.networkIO.download }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="avgLoad" :label="t('hosts.columns.avgLoad')" min-width="150">
          <template #default="{ row }">
            <div class="avg-load">
              {{ row.avgLoad.join(' / ') }}
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="300" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button
                v-if="row.agentStatus === 'not-installed'"
                size="small"
                type="primary" text
                @click="handleInstallAgent(row)"
              >
                {{ t('hosts.installAgent') }}
              </el-button>
              <el-button
                v-if="row.agentStatus === 'installed'"
                size="small"
                type="warning" text
                @click="handleUninstallAgent(row)"
              >
                {{ t('hosts.uninstallAgent') }}
              </el-button>
              <el-button size="small" text @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
              <el-button size="small" type="danger" text @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
            </div>
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
import PageTabs from '@/components/layout/PageTabs.vue'

const { t } = useI18n()

// 页面标签
const activeTab = ref('hosts')
const pageTabs = [
  { key: 'hosts', label: t('hosts.tabs.list') },
  { key: 'monitor', label: t('hosts.tabs.monitor') },
  { key: 'alerts', label: t('hosts.tabs.alerts') }
]

const handleTabChange = (key: string) => {
  activeTab.value = key
  // TODO: 实现标签切换逻辑
}

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
    lastHeartbeat: '2025-12-14 10:30:15',
    cpuModel: 'Intel(R) Xeon(R) CPU E5-2680 v4 @ 2.40GHz',
    totalMemory: '64 GB',
    registerTime: '2025-12-01 09:30:00',
    os: 'Ubuntu 20.04.3 LTS',
    kernel: '5.4.0-91-generic',
    hostname: 'master-node-01',
    networkInterfaces: [
      { name: 'eth0', ip: '192.168.1.100', speed: '1000Mbps' },
      { name: 'eth1', ip: '10.0.0.100', speed: '10Gbps' }
    ],
    networkIO: { upload: '12.5 MB/s', download: '8.3 MB/s' },
    avgLoad: [1.24, 1.58, 1.82],
    disks: [
      { name: '/dev/sda1', used: '180 GB', total: '500 GB', usage: 36, ioUtil: '12%', mountPoint: '/' },
      { name: '/dev/sdb1', used: '850 GB', total: '2 TB', usage: 43, ioUtil: '8%', mountPoint: '/data' }
    ]
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
    lastHeartbeat: '2025-12-14 10:30:12',
    cpuModel: 'Intel(R) Xeon(R) CPU E5-2680 v4 @ 2.40GHz',
    totalMemory: '128 GB',
    registerTime: '2025-12-01 10:15:00',
    os: 'CentOS 7.9',
    kernel: '3.10.0-1160.el7.x86_64',
    hostname: 'worker-node-01',
    networkInterfaces: [
      { name: 'eth0', ip: '192.168.1.101', speed: '1000Mbps' }
    ],
    networkIO: { upload: '25.8 MB/s', download: '15.2 MB/s' },
    avgLoad: [3.24, 3.15, 2.98],
    disks: [
      { name: '/dev/sda1', used: '280 GB', total: '500 GB', usage: 56, ioUtil: '28%', mountPoint: '/' },
      { name: '/dev/sdb1', used: '1.5 TB', total: '2 TB', usage: 75, ioUtil: '45%', mountPoint: '/data' },
      { name: '/dev/sdc1', used: '3.2 TB', total: '4 TB', usage: 80, ioUtil: '52%', mountPoint: '/backup' }
    ]
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
    lastHeartbeat: '-',
    cpuModel: 'Intel(R) Xeon(R) CPU E5-2680 v4 @ 2.40GHz',
    totalMemory: '64 GB',
    registerTime: '2025-12-02 14:20:00',
    os: 'Ubuntu 22.04 LTS',
    kernel: '5.15.0-56-generic',
    hostname: 'worker-node-02',
    networkInterfaces: [
      { name: 'eth0', ip: '192.168.1.102', speed: '1000Mbps' }
    ],
    networkIO: { upload: '-', download: '-' },
    avgLoad: [0, 0, 0],
    disks: [
      { name: '/dev/sda1', used: '-', total: '500 GB', usage: 0, ioUtil: '-', mountPoint: '/' }
    ]
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
    lastHeartbeat: '2025-12-14 09:15:30',
    cpuModel: 'Intel(R) Xeon(R) CPU E5-2680 v4 @ 2.40GHz',
    totalMemory: '96 GB',
    registerTime: '2025-12-03 11:10:00',
    os: 'Ubuntu 20.04.3 LTS',
    kernel: '5.4.0-91-generic',
    hostname: 'worker-node-03',
    networkInterfaces: [
      { name: 'eth0', ip: '192.168.1.103', speed: '1000Mbps' }
    ],
    networkIO: { upload: '-', download: '-' },
    avgLoad: [0, 0, 0],
    disks: [
      { name: '/dev/sda1', used: '120 GB', total: '500 GB', usage: 24, ioUtil: '-', mountPoint: '/' },
      { name: '/dev/sdb1', used: '600 GB', total: '2 TB', usage: 30, ioUtil: '-', mountPoint: '/data' }
    ]
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
  // 暗黑模式下统一使用主色作为强调色
  if (document.documentElement.classList.contains('dark')) {
    return 'var(--primary)'
  }
  if (percentage < 60) return '#67c23a'
  if (percentage < 80) return '#e6a23c'
  return '#f56c6c'
}

const getRowClassName = ({ row }: { row: any }) => {
  // 高负载行特别高亮
  if (row.cpu >= 80 || row.memory >= 80) {
    return 'high-load-row'
  }
  return ''
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

// const handleTestConnection = (row: any) => {
//   ElMessage.info(t('hosts.msg.testing', { name: row.name }))
//   // TODO: 调用 API 测试连接
//   setTimeout(() => {
//     ElMessage.success(t('hosts.msg.testSuccess'))
//   }, 1000)
// }

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

.toolbar {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  padding: 16px;
  background: var(--surface-2);
  border-radius: 6px;
  margin-bottom: 8px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

// 展开内容样式
.expand-content {
  padding: 24px 48px;
  background: var(--surface-2);
}

.info-section {
  margin-bottom: 24px;

  &:last-child {
    margin-bottom: 0;
  }
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border);
}

.disk-table {
  :deep(.el-table__header) {
    th {
      background: rgba(0, 0, 0, 0.02);
      font-weight: 600;
    }
  }
}

// 网络IO样式
.network-io {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 12px;

  .io-item {
    font-family: 'Consolas', 'Monaco', monospace;
  }
}

// 平均负载样式
.avg-load {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 13px;
}

// 操作按钮样式
.action-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: nowrap;
  align-items: center;

  .el-button {
    flex-shrink: 0;
  }
}

// 表格样式优化
.hosts-table {
  :deep(.el-table__header) {
    th {
      background: var(--surface-2);
      color: var(--text);
      font-weight: 600;
      border-color: var(--border);
    }
  }

  :deep(.el-table__body) {
    td {
      border-color: var(--border);
    }
  }

  // 高负载行高亮
  :deep(.high-load-row) {
    background: var(--text-soft-bg) !important;
    border-left: 4px solid var(--primary) !important;

    &:hover {
      background: var(--surface-2) !important;
    }

    // 第一个单元格需要补偿左边框占用的空间
    td:first-child {
      padding-left: 12px;
    }
  }
}
</style>
