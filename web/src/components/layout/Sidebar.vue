<!--
  Sidebar 组件
  导航菜单
-->
<template>
  <div class="sidebar">
    <el-menu
      :default-active="activeMenu"
      class="sidebar-menu"
      :collapse="isCollapse"
      @select="handleMenuSelect"
    >
      <el-menu-item index="/dashboard">
        <el-icon><Odometer /></el-icon>
        <template #title>{{ t('menu.dashboard') }}</template>
      </el-menu-item>
      <el-menu-item index="/hosts">
        <el-icon><Monitor /></el-icon>
        <template #title>{{ t('menu.hosts') }}</template>
      </el-menu-item>
      <el-menu-item index="/deploy">
        <el-icon><Download /></el-icon>
        <template #title>{{ t('menu.deploy') }}</template>
      </el-menu-item>
      <el-menu-item index="/tasks">
        <el-icon><List /></el-icon>
        <template #title>{{ t('menu.tasks') }}</template>
      </el-menu-item>
      <el-menu-item index="/clusters">
        <el-icon><Connection /></el-icon>
        <template #title>{{ t('menu.clusters') }}</template>
      </el-menu-item>
      <el-menu-item index="/diagnostics">
        <el-icon><Tools /></el-icon>
        <template #title>{{ t('menu.diagnostics') }}</template>
      </el-menu-item>
      <el-menu-item index="/plugins">
        <el-icon><Grid /></el-icon>
        <template #title>{{ t('menu.plugins') }}</template>
      </el-menu-item>
      <el-menu-item index="/settings">
        <el-icon><Setting /></el-icon>
        <template #title>{{ t('menu.settings') }}</template>
      </el-menu-item>
    </el-menu>
    <div class="sidebar-footer"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  Odometer,
  Monitor,
  Download,
  List,
  Connection,
  Tools,
  Grid,
  Setting
} from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()

const isCollapse = ref(false)

const activeMenu = computed(() => route.path)

const handleMenuSelect = (index: string) => {
  router.push(index)
}

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}
</script>

<style scoped>
.sidebar {
  width: 216px;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--surface);
  border-right: 1px solid var(--border);
  transition: width 0.3s;
}

.sidebar-menu {
  flex: 1;
  border-right: none;
}

.sidebar-menu:not(.el-menu--collapse) {
  width: 216px;
}

.sidebar-footer {
  padding: 12px;
  border-top: 1px solid var(--border);
  text-align: center;
}
</style>
