<!--
  Header 组件
  包含 Logo、用户信息、退出按钮
-->
<template>
  <div class="header">
    <div class="header-left">
      <div class="logo">
        <img src="/logo.png" alt="SeaTunnel" class="logo-img" />
        <div class="logo-text">
          <h1 class="logo-title">SeaTunnel</h1>
          <p class="logo-subtitle">企业级管理平台</p>
        </div>
      </div>
    </div>
    <div class="header-right">
      <el-button
        :icon="themeIcon"
        circle
        @click="toggleTheme"
        :title="themeMode === 'dark' ? '切换到明亮模式' : '切换到暗色模式'"
      />
      <el-dropdown trigger="click">
        <div class="user-info">
          <el-avatar :size="32" icon="User" />
          <span class="username">{{ username || '管理员' }}</span>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="handleProfile">个人信息</el-dropdown-item>
            <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useThemeStore } from '@/stores/theme'
import { ElMessage } from 'element-plus'
import { Sunny, Moon } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const themeStore = useThemeStore()

const username = computed(() => userStore.username)
const themeMode = computed(() => themeStore.mode)
const themeIcon = computed(() => (themeMode.value === 'dark' ? Sunny : Moon))

const toggleTheme = () => {
  themeStore.toggleTheme()
}

const handleProfile = () => {
  ElMessage.info('个人信息功能开发中')
}

const handleLogout = () => {
  userStore.logout()
  ElMessage.success('已退出登录')
  router.push('/login')
}
</script>

<style scoped>
.header {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  background: var(--surface);
  border-bottom: 1px solid var(--border);
  box-shadow: var(--shadow);
}

.header-left {
  display: flex;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-img {
  width: 40px;
  height: 40px;
  object-fit: contain;
}

.logo-text {
  display: flex;
  flex-direction: column;
}

.logo-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text);
  margin: 0;
  line-height: 1.2;
}

.logo-subtitle {
  font-size: 12px;
  color: var(--muted);
  margin: 0;
  line-height: 1.2;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: var(--surface-2);
}

.username {
  font-size: 14px;
  color: var(--text);
}
</style>
