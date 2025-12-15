/**
 * 路由配置
 * 定义应用的所有路由规则
 */
import { createRouter, createWebHashHistory } from 'vue-router'
import i18n from '@/i18n'

// 路由定义
const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { 
      titleKey: 'login.form.login',
      requiresAuth: false,
      hideLayout: true
    }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/Dashboard.vue'),
    meta: { titleKey: 'route.dashboard', requiresAuth: true }
  },
  {
    path: '/hosts',
    name: 'Hosts',
    component: () => import('@/views/Hosts.vue'),
    meta: { titleKey: 'route.hosts', requiresAuth: true }
  },
  {
    path: '/deploy',
    name: 'Deploy',
    component: () => import('@/views/Deploy.vue'),
    meta: { titleKey: 'route.deploy', requiresAuth: true }
  },
  {
    path: '/tasks',
    name: 'Tasks',
    component: () => import('@/views/Tasks.vue'),
    meta: { titleKey: 'route.tasks', requiresAuth: true }
  },
  {
    path: '/clusters',
    name: 'Clusters',
    component: () => import('@/views/Clusters.vue'),
    meta: { titleKey: 'route.clusters', requiresAuth: true }
  },
  {
    path: '/diagnostics',
    name: 'Diagnostics',
    component: () => import('@/views/Diagnostics.vue'),
    meta: { titleKey: 'route.diagnostics', requiresAuth: true }
  },
  {
    path: '/plugins',
    name: 'Plugins',
    component: () => import('@/views/Plugins.vue'),
    meta: { titleKey: 'route.plugins', requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/Settings.vue'),
    meta: { titleKey: 'route.settings', requiresAuth: true }
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// 路由守卫 - 设置页面标题
router.beforeEach((to, _from, next) => {
  // 设置页面标题（根据当前语言）
  const { t, locale } = i18n.global
  const titleKey = to.meta.titleKey as string | undefined
  if (titleKey) {
    const page = t(titleKey)
    const suffix = t('app.titleSuffix')
    document.title = `${page} - ${suffix}`
    document.documentElement.lang = String(locale.value)
  }
  next()
})

// 路由守卫 - 权限检查
router.beforeEach(async (to, _from, next) => {
  // 动态导入用户store以避免循环依赖
  const { useUserStore } = await import('@/stores/user')
  const userStore = useUserStore()
  
  const requiresAuth = to.meta.requiresAuth !== false // 默认需要认证
  const isLoginPage = to.path === '/login'
  
  // 如果是登录页面且已登录，重定向到首页
  if (isLoginPage && userStore.isLoggedIn) {
    next('/dashboard')
    return
  }
  
  // 如果需要认证但未登录，重定向到登录页
  if (requiresAuth && !userStore.isLoggedIn) {
    next({
      path: '/login',
      query: { redirect: to.fullPath }
    })
    return
  }
  
  next()
})

export default router
