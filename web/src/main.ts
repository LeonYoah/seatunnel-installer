/**
 * 应用入口文件
 */
import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import './styles/theme.scss'
import './style.css'
import App from './App.vue'
import router from './router'
import pinia from './stores'
import i18n from './i18n'

const app = createApp(App)

// 开发环境启用 Vue DevTools
if (import.meta.env.DEV) {
  // app.config.devtools = true
  app.config.performance = true
}

// 注册插件
app.use(router)
app.use(pinia)
app.use(i18n)
app.use(ElementPlus)

// 挂载应用
app.mount('#app')
