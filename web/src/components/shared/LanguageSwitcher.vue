<template>
  <el-dropdown trigger="click" @command="onChange">
    <el-button text>{{ currentLabel }}</el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item command="zh-CN">{{ t('lang.zh') }}</el-dropdown-item>
        <el-dropdown-item command="en">{{ t('lang.en') }}</el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
  </template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { locale, t } = useI18n()

const onChange = (val: 'zh-CN' | 'en') => {
  locale.value = val
  localStorage.setItem('locale', val)
  document.documentElement.lang = val
}

const currentLabel = computed(() => (locale.value === 'zh-CN' ? t('lang.zh') : t('lang.en')))
</script>

<style scoped>
</style>

