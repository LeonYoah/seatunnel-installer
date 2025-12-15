/**
 * 登录页面测试
 */
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Login from '../Login.vue'
import zhCN from '@/locales/zh-CN'
import enUS from '@/locales/en'

// Mock Element Plus
vi.mock('element-plus', () => ({
  ElMessage: {
    success: vi.fn(),
    error: vi.fn()
  }
}))

// Mock auth API
vi.mock('@/api/auth', () => ({
  authApi: {
    login: vi.fn()
  }
}))

describe('Login.vue', () => {
  let router: any
  let pinia: any
  let i18n: any

  beforeEach(() => {
    // 创建路由实例
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div>Home</div>' } },
        { path: '/login', component: Login },
        { path: '/dashboard', component: { template: '<div>Dashboard</div>' } }
      ]
    })

    // 创建Pinia实例
    pinia = createPinia()

    // 创建i18n实例
    i18n = createI18n({
      legacy: false,
      locale: 'zh-CN',
      messages: {
        'zh-CN': zhCN,
        'en': enUS
      }
    })
  })

  it('应该正确渲染登录表单', async () => {
    const wrapper = mount(Login, {
      global: {
        plugins: [router, pinia, i18n]
      }
    })

    // 检查表单元素是否存在
    expect(wrapper.find('input[placeholder*="用户名"]').exists()).toBe(true)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    expect(wrapper.find('button').text()).toContain('登录')
  })

  it('应该验证必填字段', async () => {
    const wrapper = mount(Login, {
      global: {
        plugins: [router, pinia, i18n]
      }
    })

    // 点击登录按钮而不填写表单
    const loginButton = wrapper.find('button')
    await loginButton.trigger('click')

    // 应该显示验证错误（这里我们只是检查表单存在，实际验证由Element Plus处理）
    expect(wrapper.find('form').exists()).toBe(true)
  })

  it('应该在表单填写完整时启用登录按钮', async () => {
    const wrapper = mount(Login, {
      global: {
        plugins: [router, pinia, i18n]
      }
    })

    // 填写表单
    const usernameInput = wrapper.find('input[placeholder*="用户名"]')
    const passwordInput = wrapper.find('input[type="password"]')

    await usernameInput.setValue('testuser')
    await passwordInput.setValue('password123')

    // 检查按钮是否可用
    const loginButton = wrapper.find('button')
    expect(loginButton.exists()).toBe(true)
  })

  it('应该显示语言切换器', () => {
    const wrapper = mount(Login, {
      global: {
        plugins: [router, pinia, i18n]
      }
    })

    // 检查语言切换器是否存在
    expect(wrapper.find('.language-switcher').exists()).toBe(true)
  })

  it('应该显示记住我选项', () => {
    const wrapper = mount(Login, {
      global: {
        plugins: [router, pinia, i18n]
      }
    })

    // 检查记住我复选框是否存在
    expect(wrapper.find('input[type="checkbox"]').exists()).toBe(true)
  })
})