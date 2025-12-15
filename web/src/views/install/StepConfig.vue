<!--
  步骤 1: 配置参数
-->
<template>
  <el-card class="step-card">
    <template #header>
      <span>{{ t('install.config.title') }}</span>
    </template>
    <el-form :model="form" label-width="140px">
      <!-- 基础配置 -->
      <div class="form-section">
        <h3>{{ t('install.config.basic') }}</h3>
        <el-form-item :label="t('install.config.version')">
          <el-select v-model="form.version" style="width: 300px">
            <el-option label="2.3.12" value="2.3.12" />
            <el-option label="2.3.11" value="2.3.11" />
            <el-option label="2.3.10" value="2.3.10" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('install.config.installMode')">
          <el-radio-group v-model="form.installMode">
            <el-radio label="online">{{ t('install.config.online') }}</el-radio>
            <el-radio label="offline">{{ t('install.config.offline') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="t('install.config.baseDir')">
          <el-input v-model="form.baseDir" style="width: 500px" />
        </el-form-item>
      </div>

      <!-- 部署配置 -->
      <div class="form-section">
        <h3>{{ t('install.config.deploy') }}</h3>
        <el-form-item :label="t('install.config.deployMode')">
          <el-radio-group v-model="form.deployMode">
            <el-radio label="separated">{{ t('install.config.separated') }}</el-radio>
            <el-radio label="hybrid">{{ t('install.config.hybrid') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <template v-if="form.deployMode === 'separated'">
          <el-form-item label="Master IP">
            <el-input v-model="form.masterIp" style="width: 300px" />
          </el-form-item>
          <el-form-item label="Worker IPs">
            <el-input v-model="form.workerIps" style="width: 500px" />
          </el-form-item>
        </template>
        <template v-else>
          <el-form-item :label="t('install.config.clusterNodes')">
            <el-input v-model="form.clusterNodes" style="width: 500px" />
          </el-form-item>
        </template>
      </div>
    </el-form>

    <div class="step-actions">
      <el-button type="primary" @click="handleNext">{{ t('common.next') }}</el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const emit = defineEmits(['next'])
const { t } = useI18n()

const form = ref({
  version: '2.3.12',
  installMode: 'online',
  baseDir: '/home/seatunnel/seatunnel-package',
  deployMode: 'separated',
  masterIp: '',
  workerIps: '',
  clusterNodes: ''
})

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
