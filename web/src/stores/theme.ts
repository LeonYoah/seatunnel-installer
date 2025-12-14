/**
 * 主题状态管理
 */
import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export type ThemeMode = 'light' | 'dark'

export const useThemeStore = defineStore('theme', () => {
  // 状态
  const mode = ref<ThemeMode>('dark')

  // 操作
  const setTheme = (newMode: ThemeMode) => {
    mode.value = newMode
    applyTheme(newMode)
    localStorage.setItem('theme', newMode)
  }

  const toggleTheme = () => {
    const newMode = mode.value === 'light' ? 'dark' : 'light'
    setTheme(newMode)
  }

  const applyTheme = (theme: ThemeMode) => {
    if (theme === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  // 初始化
  const init = () => {
    const savedTheme = localStorage.getItem('theme') as ThemeMode
    if (savedTheme && (savedTheme === 'light' || savedTheme === 'dark')) {
      mode.value = savedTheme
    }
    applyTheme(mode.value)
  }

  // 监听主题变化
  watch(mode, newMode => {
    applyTheme(newMode)
  })

  return {
    mode,
    setTheme,
    toggleTheme,
    init
  }
})
