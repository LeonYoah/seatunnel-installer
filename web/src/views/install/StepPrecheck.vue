<!--
  步骤 2: 环境检查
-->
<template>
  <el-card class="step-card">
    <template #header>
      <div class="card-header">
        <span>{{ t('install.steps.precheck') }}</span>
        <el-tag v-if="allPassed" type="success">{{ t('install.precheck.allPassed') }}</el-tag>
      </div>
    </template>

    <div class="precheck-list">
      <div v-for="item in items" :key="item.id" class="precheck-item">
        <div class="precheck-icon">
          <el-icon v-if="item.status === 'success'" color="#67c23a" :size="24">
            <CircleCheck />
          </el-icon>
          <el-icon v-else-if="item.status === 'failed'" color="#f56c6c" :size="24">
            <CircleClose />
          </el-icon>
          <el-icon v-else color="#909399" :size="24">
            <Loading />
          </el-icon>
        </div>
        <div class="precheck-content">
          <div class="precheck-title">{{ item.title }}</div>
          <div v-if="item.message" class="precheck-message">{{ item.message }}</div>
        </div>
      </div>
    </div>

    <div class="step-actions">
      <el-button @click="handlePrev">{{ t('common.prev') }}</el-button>
      <el-button type="primary" :disabled="!allPassed" @click="handleNext">{{ t('common.next') }}</el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { CircleCheck, CircleClose, Loading } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

const emit = defineEmits(['next', 'prev'])
const { t } = useI18n()

const items = ref([
  { id: 1, title: t('install.precheck.items.memory'), status: 'success', message: t('install.precheck.msg.memoryOk') },
  { id: 2, title: t('install.precheck.items.cpu'), status: 'success', message: t('install.precheck.msg.cpuOk') },
  { id: 3, title: t('install.precheck.items.disk'), status: 'success', message: t('install.precheck.msg.diskOk') },
  { id: 4, title: t('install.precheck.items.ssh'), status: 'success', message: t('install.precheck.msg.sshOk') },
  { id: 5, title: t('install.precheck.items.port'), status: 'success', message: t('install.precheck.msg.portOk') },
  { id: 6, title: t('install.precheck.items.firewall'), status: 'success', message: t('install.precheck.msg.firewallOk') }
])

const allPassed = computed(() => items.value.every(item => item.status === 'success'))

const handleNext = () => {
  emit('next')
}

const handlePrev = () => {
  emit('prev')
}
</script>

<style scoped>
.step-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.precheck-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.precheck-item {
  display: flex;
  gap: 16px;
  padding: 16px;
  border: 1px solid var(--border);
  border-radius: 8px;
}

.precheck-icon {
  flex-shrink: 0;
}

.precheck-content {
  flex: 1;
}

.precheck-title {
  font-weight: 600;
  color: var(--text);
  margin-bottom: 4px;
}

.precheck-message {
  font-size: 14px;
  color: var(--muted);
}

.step-actions {
  margin-top: 30px;
  text-align: right;
}
</style>
