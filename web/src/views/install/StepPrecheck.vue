<!--
  步骤 2: 环境检查
-->
<template>
  <el-card class="step-card">
    <template #header>
      <div class="card-header">
        <span>环境检查</span>
        <el-tag v-if="allPassed" type="success">全部通过</el-tag>
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
      <el-button @click="handlePrev">上一步</el-button>
      <el-button type="primary" :disabled="!allPassed" @click="handleNext">下一步</el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { CircleCheck, CircleClose, Loading } from '@element-plus/icons-vue'

const emit = defineEmits(['next', 'prev'])

const items = ref([
  { id: 1, title: '内存检查', status: 'success', message: '可用内存: 16GB' },
  { id: 2, title: 'CPU 检查', status: 'success', message: 'CPU 核心数: 8' },
  { id: 3, title: '磁盘空间检查', status: 'success', message: '可用空间: 100GB' },
  { id: 4, title: 'SSH 连通性检查', status: 'success', message: '所有节点连接正常' },
  { id: 5, title: '端口占用检查', status: 'success', message: '端口可用' },
  { id: 6, title: '防火墙状态检查', status: 'success', message: '防火墙已关闭' }
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
