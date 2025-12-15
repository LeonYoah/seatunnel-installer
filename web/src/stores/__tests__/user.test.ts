/**
 * 用户状态管理测试
 */
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUserStore } from '../user'

// Mock auth API
vi.mock('@/api/auth', () => ({
  authApi: {
    logout: vi.fn().mockResolvedValue(undefined),
    refreshToken: vi.fn(),
    getCurrentUser: vi.fn()
  }
}))

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
}
Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
})

// Mock sessionStorage
const sessionStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
}
Object.defineProperty(window, 'sessionStorage', {
  value: sessionStorageMock
})

describe('useUserStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('应该正确设置令牌', () => {
    const userStore = useUserStore()
    const token = 'test-token'

    userStore.setToken(token)

    expect(userStore.token).toBe(token)
    expect(userStore.isLoggedIn).toBe(true)
    expect(localStorageMock.setItem).toHaveBeenCalledWith('token', token)
  })

  it('应该正确设置用户信息', () => {
    const userStore = useUserStore()
    const userInfo = {
      id: '1',
      username: 'testuser',
      email: 'test@example.com',
      tenantId: 'tenant1',
      workspaceId: 'workspace1',
      roles: ['admin'],
      lastLoginAt: '2023-01-01T00:00:00Z'
    }

    userStore.setUserInfo(userInfo)

    expect(userStore.userInfo).toEqual(userInfo)
    expect(userStore.username).toBe(userInfo.username)
    expect(userStore.roles).toEqual(userInfo.roles)
    expect(localStorageMock.setItem).toHaveBeenCalledWith('userInfo', JSON.stringify(userInfo))
  })

  it('应该正确检查用户角色', () => {
    const userStore = useUserStore()
    userStore.roles = ['admin', 'user']

    expect(userStore.hasRole('admin')).toBe(true)
    expect(userStore.hasRole('guest')).toBe(false)
  })

  it('应该正确检查用户是否有任一角色', () => {
    const userStore = useUserStore()
    userStore.roles = ['admin', 'user']

    expect(userStore.hasAnyRole(['admin', 'guest'])).toBe(true)
    expect(userStore.hasAnyRole(['guest', 'visitor'])).toBe(false)
  })

  it('应该正确登出', async () => {
    const userStore = useUserStore()
    
    // 设置初始状态
    userStore.setToken('test-token')
    userStore.setUserInfo({
      id: '1',
      username: 'testuser',
      email: 'test@example.com',
      tenantId: 'tenant1',
      workspaceId: 'workspace1',
      roles: ['admin']
    })

    await userStore.logout()

    // 检查状态是否被清除
    expect(userStore.token).toBe('')
    expect(userStore.userInfo).toBeNull()
    expect(userStore.username).toBe('')
    expect(userStore.roles).toEqual([])
    expect(userStore.isLoggedIn).toBe(false)

    // 检查本地存储是否被清除
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('token')
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('userInfo')
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('refresh_token')
    expect(sessionStorageMock.removeItem).toHaveBeenCalledWith('refresh_token')
  })
})