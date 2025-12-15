<!--
  步骤 3: 插件选择
-->
<template>
  <el-card class="step-card">
    <template #header>
      <div class="card-header">
        <span>{{ t('install.steps.plugins') }}</span>
        <el-tag>{{ t('install.plugins.selected', { count: selected.length }) }}</el-tag>
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
      <el-button @click="handlePrev">{{ t('common.prev') }}</el-button>
      <el-button type="primary" @click="handleNext">{{ t('common.next') }}</el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const emit = defineEmits(['next', 'prev'])
const { t } = useI18n()

const selected = ref(['mysql-cdc', 'kafka'])

const categories = ref([
  {
    name: 'source',
    label: t('install.plugins.source'),
    plugins: [
      { id: 'mysql-cdc', name: 'MySQL-CDC', desc: 'MySQL CDC' },
      { id: 'postgres-cdc', name: 'PostgreSQL-CDC', desc: 'PostgreSQL CDC' }
    ]
  },
  {
    name: 'sink',
    label: t('install.plugins.sink'),
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
