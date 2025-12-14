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
        <template #title>总览</template>
      </el-menu-item>
      <el-menu-item index="/hosts">
        <el-icon><Monitor /></el-icon>
        <template #title>主机管理</template>
      </el-menu-item>
      <el-menu-item index="/deploy">
        <el-icon><Download /></el-icon>
        <template #title>部署集群</template>
      </el-menu-item>
      <el-menu-item index="/tasks">
        <el-icon><List /></el-icon>
        <template #title>任务管理</template>
      </el-menu-item>
      <el-menu-item index="/clusters">
        <el-icon><Connection /></el-icon>
        <template #title>集群管理</template>
      </el-menu-item>
      <el-menu-item index="/diagnostics">
        <el-icon><Tools /></el-icon>
        <template #title>诊断中心</template>
      </el-menu-item>
      <el-menu-item index="/plugins">
        <el-icon><Grid /></el-icon>
        <template #title>插件市场</template>
      </el-menu-item>
      <el-menu-item index="/settings">
        <el-icon><Setting /></el-icon>
        <template #title>设置</template>
      </el-menu-item>
    </el-menu>
    <div class="sidebar-footer">
      <el-button
        :icon="isCollapse ? 'Expand' : 'Fold'"
        text
        @click="toggleCollapse"
      >
        {{ isCollapse ? '' : '收起' }}
      </el-button>
    </div>
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

const router = useRouter()
const route = useRoute()

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
  width: 200px;
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
  width: 200px;
}

.sidebar-footer {
  padding: 12px;
  border-top: 1px solid var(--border);
  text-align: center;
}
</style>
