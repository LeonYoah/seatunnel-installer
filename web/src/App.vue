<!--
  应用根组件
-->
<template>
  <el-config-provider :locale="elLocale">
    <!-- 根据路由meta决定是否显示主布局 -->
    <MainLayout v-if="!$route.meta.hideLayout" />
    <router-view v-else />
  </el-config-provider>
</template>

<script setup lang="ts">
import MainLayout from './components/layout/MainLayout.vue'
import { onMounted, computed } from 'vue'
import { useUserStore } from './stores/user'
import { useThemeStore } from './stores/theme'
import { useI18n } from 'vue-i18n'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'

const userStore = useUserStore()
const themeStore = useThemeStore()
const { locale } = useI18n()

const elLocale = computed(() => (locale.value === 'zh-CN' ? zhCn : en))

onMounted(() => {
  // 初始化用户状态
  userStore.init()
  // 初始化主题
  themeStore.init()
  // 设置 html lang
  document.documentElement.lang = String(locale.value)
})
</script>

<style>
/* 全局样式重置 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body {
  width: 100%;
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
    sans-serif;
}

#app {
  width: 100%;
  height: 100%;
}
</style>
