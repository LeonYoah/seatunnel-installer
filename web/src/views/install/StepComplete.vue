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
      <h2>{{ t('install.complete.title') }}</h2>
      <p class="complete-desc">{{ t('install.complete.subtitle') }}</p>

      <div class="info-section">
        <h3>{{ t('install.complete.clusterInfo') }}</h3>
        <el-descriptions :column="1" border>
          <el-descriptions-item :label="t('install.complete.clusterName')">production-cluster</el-descriptions-item>
          <el-descriptions-item :label="t('install.complete.version')">2.3.12</el-descriptions-item>
          <el-descriptions-item :label="t('install.complete.deployMode')">{{ t('install.config.separated') }} (Master/Worker)</el-descriptions-item>
          <el-descriptions-item :label="t('install.complete.nodeCount')">3</el-descriptions-item>
          <el-descriptions-item :label="t('install.complete.installPath')">
            /home/seatunnel/seatunnel-package/apache-seatunnel-2.3.12
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <div class="info-section">
        <h3>{{ t('install.complete.access') }}</h3>
        <el-alert type="info" :closable="false">
          <template #default>
            <div>{{ t('install.complete.masterHttp') }} <strong>http://192.168.1.100:8080</strong></div>
          </template>
        </el-alert>
      </div>

      <div class="info-section">
        <h3>{{ t('install.complete.quickStart') }}</h3>
        <el-steps direction="vertical" :active="0">
          <el-step :title="t('install.complete.stepCheck')">
            <template #description>
              <code>systemctl status seatunnel-master</code>
            </template>
          </el-step>
          <el-step :title="t('install.complete.stepSubmit')">
            <template #description>
              {{ t('install.complete.stepSubmitDesc') }}
            </template>
          </el-step>
          <el-step :title="t('install.complete.stepMonitor')">
            <template #description>
              {{ t('install.complete.stepMonitorDesc') }}
            </template>
          </el-step>
        </el-steps>
      </div>

      <div class="action-buttons">
        <el-button type="primary" size="large" @click="goToDashboard">
          {{ t('install.complete.goConsole') }}
        </el-button>
        <el-button size="large" @click="goToClusters">{{ t('install.complete.viewClusters') }}</el-button>
        <el-button size="large" @click="downloadReport">{{ t('install.complete.downloadReport') }}</el-button>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { SuccessFilled } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

const router = useRouter()
const { t } = useI18n()

const goToDashboard = () => {
  router.push('/dashboard')
}

const goToClusters = () => {
  router.push('/clusters')
}

const downloadReport = () => {
  ElMessage.success(t('install.complete.reportStart'))
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
