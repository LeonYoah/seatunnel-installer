<!--
  步骤 3: 插件选择
-->
<template>
  <el-card class="step-card">
    <template #header>
      <div class="card-header">
        <span>插件选择</span>
        <el-tag>已选 {{ selected.length }} 个插件</el-tag>
      </div>
    </template>

    <el-checkbox-group v-model="selected">
      <div v-for="category in categories" :key="category.name" class="plugin-category">
        <h3>{{ category.label }}</h3>
        <div class="plugin-list">
          <div v-for="plugin in category.plugins" :key="plugin.id" class="plugin-item">
            <el-checkbox :label="plugin.id">
              {{ plugin.name }}
            </el-checkbox>
            <span class="plugin-desc">{{ plugin.desc }}</span>
          </div>
        </div>
      </div>
    </el-checkbox-group>

    <div class="step-actions">
      <el-button @click="handlePrev">上一步</el-button>
      <el-button type="primary" @click="handleNext">下一步</el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits(['next', 'prev'])

const selected = ref(['mysql-cdc', 'kafka'])

const categories = ref([
  {
    name: 'source',
    label: 'Source 连接器',
    plugins: [
      { id: 'mysql-cdc', name: 'MySQL-CDC', desc: 'MySQL 变更数据捕获' },
      { id: 'postgres-cdc', name: 'PostgreSQL-CDC', desc: 'PostgreSQL 变更数据捕获' }
    ]
  },
  {
    name: 'sink',
    label: 'Sink 连接器',
    plugins: [
      { id: 'kafka', name: 'Kafka', desc: 'Apache Kafka' },
      { id: 'iceberg', name: 'Iceberg', desc: 'Apache Iceberg' }
    ]
  }
])

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

.plugin-category {
  margin-bottom: 30px;
}

.plugin-category h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text);
}

.plugin-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.plugin-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.plugin-desc {
  font-size: 14px;
  color: var(--muted);
}

.step-actions {
  margin-top: 30px;
  text-align: right;
}
</style>
