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
            <span>通知渠道</span>
          </template>
          <el-form :model="notifyForm" label-width="120px">
            <el-form-item label="钉钉 Webhook">
              <el-input
                v-model="notifyForm.dingding"
                placeholder="https://oapi.dingtalk.com/robot/send?access_token=..."
              />
            </el-form-item>
            <el-form-item label="企微 Webhook">
              <el-input
                v-model="notifyForm.wecom"
                placeholder="https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=..."
              />
            </el-form-item>
            <el-form-item label="邮件接收人">
              <el-input v-model="notifyForm.email" placeholder="ops@example.com" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveNotify">保存</el-button>
              <el-button @click="handleTestNotify">发送测试</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- 系统配置 -->
      <el-col :xs="24" :md="12">
        <el-card>
          <template #header>
            <span>系统配置</span>
          </template>
          <el-form :model="systemForm" label-width="120px">
            <el-form-item label="心跳超时">
              <el-input-number
                v-model="systemForm.heartbeatTimeout"
                :min="10"
                :max="300"
                :step="10"
              />
              <span style="margin-left: 8px">秒</span>
            </el-form-item>
            <el-form-item label="日志保留天数">
              <el-input-number v-model="systemForm.logRetention" :min="1" :max="365" />
              <span style="margin-left: 8px">天</span>
            </el-form-item>
            <el-form-item label="自动备份">
              <el-switch v-model="systemForm.autoBackup" />
            </el-form-item>
            <el-form-item label="备份时间">
              <el-time-picker
                v-model="systemForm.backupTime"
                format="HH:mm"
                placeholder="选择时间"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveSystem">保存</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>

    <!-- 用户管理 -->
    <el-card class="users-card">
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-button type="primary" :icon="Plus" @click="handleAddUser">添加用户</el-button>
        </div>
      </template>
      <el-table :data="users" style="width: 100%">
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="role" label="角色" width="120">
          <template #default="{ row }">
            <el-tag :type="getRoleType(row.role)" size="small">
              {{ getRoleText(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? '激活' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastLogin" label="最后登录" width="180" />
        <el-table-column prop="createdAt" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEditUser(row)">编辑</el-button>
            <el-button size="small" @click="handleResetPassword(row)">重置密码</el-button>
            <el-button size="small" type="danger" @click="handleDeleteUser(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 审计日志 -->
    <el-card class="audit-card">
      <template #header>
        <div class="card-header">
          <span>审计日志</span>
          <el-button :icon="Download" @click="handleExportAudit">导出</el-button>
        </div>
      </template>
      <el-table :data="auditLogs" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user" label="操作者" width="120" />
        <el-table-column prop="action" label="操作" width="150" />
        <el-table-column prop="resource" label="资源" width="150" />
        <el-table-column prop="result" label="结果" width="100">
          <template #default="{ row }">
            <el-tag :type="row.result === 'success' ? 'success' : 'danger'" size="small">
              {{ row.result === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP 地址" width="150" />
        <el-table-column prop="timestamp" label="时间" width="180" />
        <el-table-column prop="details" label="详情" min-width="200" show-overflow-tooltip />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Download } from '@element-plus/icons-vue'

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
    admin: '管理员',
    operator: '操作员',
    viewer: '查看者'
  }
  return textMap[role] || role
}

const handleSaveNotify = () => {
  ElMessage.success('通知配置已保存')
}

const handleTestNotify = () => {
  ElMessage.success('测试通知已发送')
}

const handleSaveSystem = () => {
  ElMessage.success('系统配置已保存')
}

const handleAddUser = () => {
  ElMessage.info('添加用户功能开发中')
}

const handleEditUser = (row: any) => {
  ElMessage.info(`编辑用户：${row.username}`)
}

const handleResetPassword = (row: any) => {
  ElMessageBox.confirm(`确定要重置用户 "${row.username}" 的密码吗？`, '确认重置', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(() => {
      ElMessage.success('密码已重置')
    })
    .catch(() => {
      ElMessage.info('已取消重置')
    })
}

const handleDeleteUser = (row: any) => {
  ElMessageBox.confirm(`确定要删除用户 "${row.username}" 吗？`, '确认删除', {
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

const handleExportAudit = () => {
  ElMessage.success('开始导出审计日志')
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
