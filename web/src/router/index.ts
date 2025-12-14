/**
 * 路由配置
 * 定义应用的所有路由规则
 */
import { createRouter, createWebHashHistory } from 'vue-router'

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
    meta: { title: '总览' }
  },
  {
    path: '/hosts',
    name: 'Hosts',
    component: () => import('@/views/Hosts.vue'),
    meta: { title: '主机管理' }
  },
  {
    path: '/deploy',
    name: 'Deploy',
    component: () => import('@/views/Deploy.vue'),
    meta: { title: '部署集群' }
  },
  {
    path: '/tasks',
    name: 'Tasks',
    component: () => import('@/views/Tasks.vue'),
    meta: { title: '任务管理' }
  },
  {
    path: '/clusters',
    name: 'Clusters',
    component: () => import('@/views/Clusters.vue'),
    meta: { title: '集群管理' }
  },
  {
    path: '/diagnostics',
    name: 'Diagnostics',
    component: () => import('@/views/Diagnostics.vue'),
    meta: { title: '诊断中心' }
  },
  {
    path: '/plugins',
    name: 'Plugins',
    component: () => import('@/views/Plugins.vue'),
    meta: { title: '插件市场' }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/Settings.vue'),
    meta: { title: '设置' }
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// 路由守卫 - 设置页面标题
router.beforeEach((to, _from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - SeaTunnel 企业级管理平台`
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
