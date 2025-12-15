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
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/Dashboard.vue'),
    meta: { titleKey: 'route.dashboard' }
  },
  {
    path: '/hosts',
    name: 'Hosts',
    component: () => import('@/views/Hosts.vue'),
    meta: { titleKey: 'route.hosts' }
  },
  {
    path: '/deploy',
    name: 'Deploy',
    component: () => import('@/views/Deploy.vue'),
    meta: { titleKey: 'route.deploy' }
  },
  {
    path: '/tasks',
    name: 'Tasks',
    component: () => import('@/views/Tasks.vue'),
    meta: { titleKey: 'route.tasks' }
  },
  {
    path: '/clusters',
    name: 'Clusters',
    component: () => import('@/views/Clusters.vue'),
    meta: { titleKey: 'route.clusters' }
  },
  {
    path: '/diagnostics',
    name: 'Diagnostics',
    component: () => import('@/views/Diagnostics.vue'),
    meta: { titleKey: 'route.diagnostics' }
  },
  {
    path: '/plugins',
    name: 'Plugins',
    component: () => import('@/views/Plugins.vue'),
    meta: { titleKey: 'route.plugins' }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/Settings.vue'),
    meta: { titleKey: 'route.settings' }
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

// 路由守卫 - 权限检查（预留）
router.beforeEach((_to, _from, next) => {
  // TODO: 实现权限检查逻辑
  // const token = localStorage.getItem('token')
  // if (!token && to.path !== '/login') {
  //   next('/login')
  // } else {
  //   next()
  // }
  next()
})

export default router
