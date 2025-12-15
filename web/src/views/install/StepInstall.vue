<!--
  步骤 4: 安装部署
-->
<template>
  <el-card class="step-card">
    <template #header>
      <div class="card-header">
        <span>{{ t('install.steps.install') }}</span>
        <el-tag v-if="installing" type="primary">{{ t('install.install.installing') }}</el-tag>
        <el-tag v-else-if="completed" type="success">{{ t('install.install.completed') }}</el-tag>
      </div>
    </template>

    <div class="install-steps">
      <div v-for="step in steps" :key="step.id" class="install-step" :class="getStepClass(step)">
        <div class="step-icon">
          <el-icon v-if="step.status === 'completed'" :size="20">
            <CircleCheck />
          </el-icon>
          <el-icon v-else-if="step.status === 'running'" :size="20">
            <Loading />
          </el-icon>
          <span v-else>{{ step.id }}</span>
        </div>
        <div class="step-text">{{ step.text }}</div>
      </div>
    </div>

    <div class="log-section">
      <div class="log-header">
        <span>{{ t('install.install.logs') }}</span>
        <el-button size="small" text @click="scrollToBottom">{{ t('install.install.scrollBottom') }}</el-button>
      </div>
      <div ref="logContainer" class="log-content">
        <div v-for="(log, index) in logs" :key="index" class="log-line">{{ log }}</div>
      </div>
    </div>

    <div class="step-actions">
      <el-button @click="handlePrev" :disabled="installing">{{ t('common.prev') }}</el-button>
      <el-button type="primary" :disabled="!completed" @click="handleNext">{{ t('common.ok') }}</el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { CircleCheck, Loading } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

const emit = defineEmits(['next', 'prev'])
const { t } = useI18n()

const installing = ref(false)
const completed = ref(false)
const logContainer = ref<HTMLElement>()

const steps = ref([
  { id: 1, text: t('install.install.steps.download'), status: 'pending' },
  { id: 2, text: t('install.install.steps.unpack'), status: 'pending' },
  { id: 3, text: t('install.install.steps.configure'), status: 'pending' },
  { id: 4, text: t('install.install.steps.plugins'), status: 'pending' },
  { id: 5, text: t('install.install.steps.distribute'), status: 'pending' },
  { id: 6, text: t('install.install.steps.start'), status: 'pending' }
])

const logs = ref<string[]>([])

const getStepClass = (step: any) => {
  return `status-${step.status}`
}

const scrollToBottom = () => {
  if (logContainer.value) {
    const el = logContainer.value as HTMLElement
    el.scrollTop = el.scrollHeight
  }
}

const simulateInstall = async () => {
  installing.value = true

  for (let i = 0; i < steps.value.length; i++) {
    const step = steps.value[i]
    if (!step) continue

    step.status = 'running'
    logs.value.push(`[INFO] 开始执行: ${step.text}`)

    await new Promise(resolve => setTimeout(resolve, 1000))

    step.status = 'completed'
    logs.value.push(`[SUCCESS] 完成: ${step.text}`)

    await nextTick()
    scrollToBottom()
  }

  installing.value = false
  completed.value = true
  logs.value.push(t('install.install.done'))
}

onMounted(() => {
  simulateInstall()
})

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

.install-steps {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 20px;
}

.install-step {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  border: 1px solid var(--border);
}

.install-step.status-running {
  background: rgba(64, 158, 255, 0.1);
  border-color: var(--primary);
}

.install-step.status-completed {
  background: rgba(103, 194, 58, 0.1);
  border-color: var(--success);
}

.step-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--surface-2);
  color: var(--text);
  font-weight: 600;
}

.step-text {
  flex: 1;
  color: var(--text);
}

.log-section {
  margin-top: 20px;
}

.log-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  font-weight: 600;
  color: var(--text);
}

.log-content {
  height: 300px;
  overflow-y: auto;
  padding: 12px;
  background: var(--surface-2);
  border: 1px solid var(--border);
  border-radius: 8px;
  font-family: monospace;
  font-size: 13px;
}

.log-line {
  color: var(--text);
  line-height: 1.6;
}

.step-actions {
  margin-top: 30px;
  text-align: right;
}
</style>
