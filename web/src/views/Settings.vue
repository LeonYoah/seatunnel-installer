<!--
  设置页面
-->
<template>
  <div class="settings">
    <el-row :gutter="20">
      <!-- 通知渠道 -->
      <el-col :xs="24" :md="12">
        <el-card>
          <template #header>
            <span>{{ t('settings.notify.title') }}</span>
          </template>
          <el-form :model="notifyForm" label-width="120px">
            <el-form-item :label="t('settings.notify.dingWebhook')">
              <el-input
                v-model="notifyForm.dingding"
                placeholder="https://oapi.dingtalk.com/robot/send?access_token=..."
              />
            </el-form-item>
            <el-form-item :label="t('settings.notify.wecomWebhook')">
              <el-input
                v-model="notifyForm.wecom"
                placeholder="https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=..."
              />
            </el-form-item>
            <el-form-item :label="t('settings.notify.email')">
              <el-input v-model="notifyForm.email" placeholder="ops@example.com" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveNotify">{{ t('common.save') }}</el-button>
              <el-button @click="handleTestNotify">{{ t('settings.notify.sendTest') }}</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- 系统配置 -->
      <el-col :xs="24" :md="12">
        <el-card>
          <template #header>
            <span>{{ t('settings.system.title') }}</span>
          </template>
          <el-form :model="systemForm" label-width="120px">
            <el-form-item :label="t('settings.system.heartbeatTimeout')">
              <el-input-number
                v-model="systemForm.heartbeatTimeout"
                :min="10"
                :max="300"
                :step="10"
              />
              <span style="margin-left: 8px">{{ t('settings.system.seconds') }}</span>
            </el-form-item>
            <el-form-item :label="t('settings.system.logRetention')">
              <el-input-number v-model="systemForm.logRetention" :min="1" :max="365" />
              <span style="margin-left: 8px">{{ t('settings.system.days') }}</span>
            </el-form-item>
            <el-form-item :label="t('settings.system.autoBackup')">
              <el-switch v-model="systemForm.autoBackup" />
            </el-form-item>
            <el-form-item :label="t('settings.system.backupTime')">
              <el-time-picker
                v-model="systemForm.backupTime"
                format="HH:mm"
                :placeholder="t('settings.system.pickTime')"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveSystem">{{ t('common.save') }}</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>

    <!-- 用户管理 -->
    <el-card class="users-card">
      <template #header>
        <div class="card-header">
          <span>{{ t('settings.users.title') }}</span>
          <el-button type="primary" :icon="Plus" @click="handleAddUser">{{ t('settings.users.add') }}</el-button>
        </div>
      </template>
      <el-table :data="users" style="width: 100%">
        <el-table-column prop="username" :label="t('settings.users.columns.username')" width="150" />
        <el-table-column prop="email" :label="t('settings.users.columns.email')" width="200" />
        <el-table-column prop="role" :label="t('settings.users.columns.role')" width="120">
          <template #default="{ row }">
            <el-tag :type="getRoleType(row.role)" size="small">
              {{ getRoleText(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="t('settings.users.columns.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? t('settings.users.active') : t('settings.users.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastLogin" :label="t('settings.users.columns.lastLogin')" width="180" />
        <el-table-column prop="createdAt" :label="t('settings.users.columns.createdAt')" width="180" />
        <el-table-column :label="t('common.actions')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEditUser(row)">{{ t('common.edit') }}</el-button>
            <el-button size="small" @click="handleResetPassword(row)">{{ t('settings.users.resetPwd') }}</el-button>
            <el-button size="small" type="danger" @click="handleDeleteUser(row)">
              {{ t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 审计日志 -->
    <el-card class="audit-card">
      <template #header>
        <div class="card-header">
          <span>{{ t('settings.audit.title') }}</span>
          <el-button :icon="Download" @click="handleExportAudit">{{ t('common.export') }}</el-button>
        </div>
      </template>
      <el-table :data="auditLogs" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user" :label="t('settings.audit.columns.user')" width="120" />
        <el-table-column prop="action" :label="t('settings.audit.columns.action')" width="150" />
        <el-table-column prop="resource" :label="t('settings.audit.columns.resource')" width="150" />
        <el-table-column prop="result" :label="t('settings.audit.columns.result')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.result === 'success' ? 'success' : 'danger'" size="small">
              {{ row.result === 'success' ? t('status.success') : t('status.failed') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip" :label="t('settings.audit.columns.ip')" width="150" />
        <el-table-column prop="timestamp" :label="t('settings.audit.columns.timestamp')" width="180" />
        <el-table-column prop="details" :label="t('settings.audit.columns.details')" min-width="200" show-overflow-tooltip />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Download } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

// 通知配置
const notifyForm = ref({
  dingding: '',
  wecom: '',
  email: 'ops@example.com'
})

// 系统配置
const systemForm = ref({
  heartbeatTimeout: 60,
  logRetention: 30,
  autoBackup: true,
  backupTime: new Date(2000, 1, 1, 2, 0)
})

// 用户列表
const users = ref([
  {
    id: 1,
    username: 'admin',
    email: 'admin@example.com',
    role: 'admin',
    status: 'active',
    lastLogin: '2025-12-14 10:30',
    createdAt: '2025-01-01 00:00'
  },
  {
    id: 2,
    username: 'zhangsan',
    email: 'zhangsan@example.com',
    role: 'operator',
    status: 'active',
    lastLogin: '2025-12-14 09:15',
    createdAt: '2025-03-15 10:20'
  },
  {
    id: 3,
    username: 'lisi',
    email: 'lisi@example.com',
    role: 'viewer',
    status: 'active',
    lastLogin: '2025-12-13 16:45',
    createdAt: '2025-06-20 14:30'
  }
])

// 审计日志
const auditLogs = ref([
  {
    id: 1,
    user: 'admin',
    action: '创建任务',
    resource: 'task:orders_sync',
    result: 'success',
    ip: '192.168.1.100',
    timestamp: '2025-12-14 10:30:15',
    details: '创建实时同步任务'
  },
  {
    id: 2,
    user: 'zhangsan',
    action: '重启节点',
    resource: 'node:worker-01',
    result: 'success',
    ip: '192.168.1.101',
    timestamp: '2025-12-14 09:15:22',
    details: '手动重启节点'
  },
  {
    id: 3,
    user: 'lisi',
    action: '删除任务',
    resource: 'task:test_job',
    result: 'failed',
    ip: '192.168.1.102',
    timestamp: '2025-12-14 08:45:10',
    details: '权限不足'
  }
])

const getRoleType = (role: string) => {
  const typeMap: Record<string, any> = {
    admin: 'danger',
    operator: 'warning',
    viewer: 'info'
  }
  return typeMap[role] || 'info'
}

const getRoleText = (role: string) => {
  const textMap: Record<string, string> = {
    admin: t('settings.users.roles.admin'),
    operator: t('settings.users.roles.operator'),
    viewer: t('settings.users.roles.viewer')
  }
  return textMap[role] || role
}

const handleSaveNotify = () => {
  ElMessage.success(t('settings.notify.saved'))
}

const handleTestNotify = () => {
  ElMessage.success(t('settings.notify.testSent'))
}

const handleSaveSystem = () => {
  ElMessage.success(t('settings.system.saved'))
}

const handleAddUser = () => {
  ElMessage.info(t('settings.users.addWip'))
}

const handleEditUser = (row: any) => {
  ElMessage.info(t('settings.users.edit', { name: row.username }))
}

const handleResetPassword = (row: any) => {
  ElMessageBox.confirm(t('settings.users.confirmReset', { name: row.username }), t('settings.users.resetPwd'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning'
  })
    .then(() => {
      ElMessage.success(t('settings.users.resetOk'))
    })
    .catch(() => {
      ElMessage.info(t('common.cancelled'))
    })
}

const handleDeleteUser = (row: any) => {
  ElMessageBox.confirm(t('settings.users.confirmDelete', { name: row.username }), t('common.confirmDelete'), {
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

const handleExportAudit = () => {
  ElMessage.success(t('settings.audit.exporting'))
}
</script>

<style scoped>
.settings {
  width: 100%;
}

.users-card,
.audit-card {
  margin-top: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
