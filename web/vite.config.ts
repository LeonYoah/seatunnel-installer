import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import VueDevTools from 'vite-plugin-vue-devtools'
import { resolve } from 'path'
import { getVueDevToolsPluginConfig } from './src/vite/plugins/dev-tools-config'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    // 添加 VueDevTools 插件
    // 默认使用 vscode 其他 IDE 使用 package 指定命令启动
    VueDevTools(getVueDevToolsPluginConfig()),
    vue()
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
