<!--
  PageTabs - 页面二级导航组件
  用于页面顶部的标签导航，支持右侧操作按钮区域
-->
<template>
  <div class="page-tabs">
    <div class="tabs-left">
      <div
        v-for="tab in tabs"
        :key="tab.key"
        :class="['tab-item', { active: tab.key === activeKey }]"
        @click="handleTabClick(tab)"
      >
        {{ tab.label }}
      </div>
    </div>
    <div class="tabs-right">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue'

interface Tab {
  key: string
  label: string
  route?: string
}

interface Props {
  tabs: Tab[]
  activeKey: string
}

defineProps<Props>()

const emit = defineEmits<{
  (e: 'tab-change', key: string): void
}>()

const handleTabClick = (tab: Tab) => {
  emit('tab-change', tab.key)
}
</script>

<style lang="scss" scoped>
.page-tabs {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 0 16px 0;
  border-bottom: 1px solid var(--border);
  margin-bottom: 24px;
}

.tabs-left {
  display: flex;
  gap: 32px;
  align-items: center;
}

.tab-item {
  position: relative;
  padding: 8px 0;
  font-size: 15px;
  font-weight: 500;
  color: var(--muted);
  cursor: pointer;
  transition: color 0.2s;
  user-select: none;

  &:hover {
    color: var(--text);
  }

  &.active {
    color: var(--primary);
    font-weight: 600;

    &::after {
      content: '';
      position: absolute;
      bottom: -17px;
      left: 0;
      right: 0;
      height: 2px;
      background: var(--primary);
      border-radius: 2px 2px 0 0;
    }
  }
}

.tabs-right {
  display: flex;
  gap: 12px;
  align-items: center;
}
</style>
