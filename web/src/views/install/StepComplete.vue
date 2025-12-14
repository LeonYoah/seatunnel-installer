<!--
  步骤 5: 完成
-->
<template>
  <el-card class="step-card">
    <div class="complete-content">
      <div class="complete-icon">
        <el-icon color="#67c23a" :size="80">
          <SuccessFilled />
        </el-icon>
      </div>
      <h2>SeaTunnel 集群部署成功！</h2>
      <p class="complete-desc">恭喜！您的 SeaTunnel 集群已成功部署并启动。</p>

      <div class="info-section">
        <h3>集群信息</h3>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="集群名称">production-cluster</el-descriptions-item>
          <el-descriptions-item label="SeaTunnel 版本">2.3.12</el-descriptions-item>
          <el-descriptions-item label="部署模式">分离模式 (Master/Worker)</el-descriptions-item>
          <el-descriptions-item label="节点数量">3 个节点</el-descriptions-item>
          <el-descriptions-item label="安装路径">
            /home/seatunnel/seatunnel-package/apache-seatunnel-2.3.12
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <div class="info-section">
        <h3>访问地址</h3>
        <el-alert type="info" :closable="false">
          <template #default>
            <div>Master HTTP 端口: <strong>http://192.168.1.100:8080</strong></div>
          </template>
        </el-alert>
      </div>

      <div class="info-section">
        <h3>快速开始</h3>
        <el-steps direction="vertical" :active="0">
          <el-step title="查看集群状态">
            <template #description>
              <code>systemctl status seatunnel-master</code>
            </template>
          </el-step>
          <el-step title="提交第一个任务">
            <template #description>
              前往任务管理页面创建和提交数据集成任务
            </template>
          </el-step>
          <el-step title="监控集群">
            <template #description>
              在集群管理页面查看节点状态和资源使用情况
            </template>
          </el-step>
        </el-steps>
      </div>

      <div class="action-buttons">
        <el-button type="primary" size="large" @click="goToDashboard">
          进入控制台
        </el-button>
        <el-button size="large" @click="goToClusters">查看集群</el-button>
        <el-button size="large" @click="downloadReport">下载安装报告</el-button>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { SuccessFilled } from '@element-plus/icons-vue'

const router = useRouter()

const goToDashboard = () => {
  router.push('/dashboard')
}

const goToClusters = () => {
  router.push('/clusters')
}

const downloadReport = () => {
  ElMessage.success('开始下载安装报告')
}
</script>

<style scoped>
.step-card {
  margin-bottom: 20px;
}

.complete-content {
  text-align: center;
  padding: 40px 20px;
}

.complete-icon {
  margin-bottom: 24px;
}

.complete-content h2 {
  margin: 0 0 12px 0;
  font-size: 28px;
  font-weight: 700;
  color: var(--text);
}

.complete-desc {
  margin: 0 0 40px 0;
  font-size: 16px;
  color: var(--muted);
}

.info-section {
  margin-bottom: 32px;
  text-align: left;
}

.info-section h3 {
  margin: 0 0 16px 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text);
}

.info-section code {
  padding: 4px 8px;
  background: var(--surface-2);
  border: 1px solid var(--border);
  border-radius: 4px;
  font-family: monospace;
  color: var(--text);
}

.action-buttons {
  margin-top: 40px;
  display: flex;
  gap: 12px;
  justify-content: center;
}
</style>
