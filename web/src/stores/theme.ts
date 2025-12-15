/**
 * 主题状态管理
 */
import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export type ThemeMode = 'light' | 'dark'

export const useThemeStore = defineStore('theme', () => {
  // 状态
  const mode = ref<ThemeMode>('light')

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

  const setPrimary = (hex: string) => {
    const setVar = (name: string, val: string) => {
      document.documentElement.style.setProperty(name, val)
    }
    // 计算加深色
    const shade = (col: string, amt: number) => {
      let usePound = false
      let c = col
      if (c[0] === '#') {
        c = c.slice(1)
        usePound = true
      }
      const num = parseInt(c, 16)
      let r = (num >> 16) + amt
      let g = ((num >> 8) & 0x00ff) + amt
      let b = (num & 0x0000ff) + amt
      r = r < 0 ? 0 : r > 255 ? 255 : r
      g = g < 0 ? 0 : g > 255 ? 255 : g
      b = b < 0 ? 0 : b > 255 ? 255 : b
      return (usePound ? '#' : '') + ((r << 16) | (g << 8) | b).toString(16).padStart(6, '0')
    }

    setVar('--primary', hex)
    setVar('--primary-2', shade(hex, -30))
    // 同步 Element Plus 主色
    setVar('--el-color-primary', hex)
  }

  const setPrimaryFromLogo = async (logoUrl = '/logo.png') => {
    try {
      const img = new Image()
      img.crossOrigin = 'anonymous'
      img.src = logoUrl
      await new Promise((resolve, reject) => {
        img.onload = resolve
        img.onerror = reject
      })
      const canvas = document.createElement('canvas')
      canvas.width = 32
      canvas.height = 32
      const ctx = canvas.getContext('2d')
      if (!ctx) return
      ctx.drawImage(img, 0, 0, 32, 32)
      const data = ctx.getImageData(0, 0, 32, 32).data
      let r = 0, g = 0, b = 0, count = 0
      for (let i = 0; i < data.length; i += 4) {
        const rr = data[i], gg = data[i + 1], bb = data[i + 2], a = data[i + 3]
        if (a < 200) continue // 忽略透明像素
        // 稍微偏向蓝色像素
        const bias = bb > rr && bb > gg ? 1.2 : 1.0
        r += rr * bias
        g += gg * bias
        b += bb * bias
        count++
      }
      if (count > 0) {
        r = Math.round(r / count)
        g = Math.round(g / count)
        b = Math.round(b / count)
        const hex = '#' + [r, g, b].map(v => v.toString(16).padStart(2, '0')).join('')
        setPrimary(hex)
      }
    } catch {
      // 忽略，维持默认主题色
    }
  }

  // 初始化
  const init = () => {
    const savedTheme = localStorage.getItem('theme') as ThemeMode
    if (savedTheme && (savedTheme === 'light' || savedTheme === 'dark')) {
      mode.value = savedTheme
    }
    applyTheme(mode.value)
    // 初始化主色调与 logo 保持一致（异步）
    // setPrimaryFromLogo()
  }

  // 监听主题变化
  watch(mode, newMode => {
    applyTheme(newMode)
  })

  return {
    mode,
    setTheme,
    toggleTheme,
    setPrimary,
    setPrimaryFromLogo,
    init
  }
})
