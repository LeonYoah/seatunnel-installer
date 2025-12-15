/**
 * 用户状态管理
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authApi } from '@/api/auth'

// 用户信息接口
export interface UserInfo {
  id: string
  username: string
  email: string
  tenantId: string
  workspaceId: string
  roles: string[]
  lastLoginAt?: string
}

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref<string>('')
  const userInfo = ref<UserInfo | null>(null)
  const isLoggedIn = ref<boolean>(false)

  // Getters
  const username = ref<string>('')
  const roles = ref<string[]>([])

  // 操作
  const setToken = (newToken: string) => {
    token.value = newToken
    isLoggedIn.value = true
    localStorage.setItem('token', newToken)
  }

  const setUserInfo = (info: UserInfo) => {
    userInfo.value = info
    username.value = info.username
    roles.value = info.roles
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  const logout = async () => {
    try {
      // 调用后端登出接口
      if (token.value) {
        await authApi.logout()
      }
    } catch (error) {
      console.error('登出API调用失败:', error)
    } finally {
      // 清除本地状态
      token.value = ''
      userInfo.value = null
      username.value = ''
      roles.value = []
      isLoggedIn.value = false
      
      // 清除本地存储
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      localStorage.removeItem('refresh_token')
      sessionStorage.removeItem('refresh_token')
    }
  }

  // 检查用户是否有指定角色
  const hasRole = (role: string): boolean => {
    return roles.value.includes(role)
  }

  // 检查用户是否有任一指定角色
  const hasAnyRole = (roleList: string[]): boolean => {
    return roleList.some(role => roles.value.includes(role))
  }

  // 自动刷新令牌
  const refreshToken = async (): Promise<boolean> => {
    try {
      const refreshToken = localStorage.getItem('refresh_token') || sessionStorage.getItem('refresh_token')
      if (!refreshToken) {
        return false
      }

      const response = await authApi.refreshToken({ refresh_token: refreshToken })
      
      // 更新令牌和用户信息
      setToken(response.access_token)
      setUserInfo({
        id: response.user.id,
        username: response.user.username,
        email: response.user.email,
        tenantId: response.user.tenant_id,
        workspaceId: response.user.workspace_id,
        roles: response.user.roles,
        lastLoginAt: response.user.last_login_at
      })

      return true
    } catch (error) {
      console.error('刷新令牌失败:', error)
      // 刷新失败，清除所有认证信息
      await logout()
      return false
    }
  }

  // 初始化用户状态
  const init = async () => {
    const savedToken = localStorage.getItem('token')
    const savedUserInfo = localStorage.getItem('userInfo')
    
    if (savedToken && savedUserInfo) {
      try {
        token.value = savedToken
        const parsedUserInfo = JSON.parse(savedUserInfo)
        setUserInfo(parsedUserInfo)
        isLoggedIn.value = true

        // 尝试获取最新用户信息以验证令牌有效性
        try {
          const currentUser = await authApi.getCurrentUser()
          setUserInfo({
            id: currentUser.id,
            username: currentUser.username,
            email: currentUser.email,
            tenantId: currentUser.tenant_id,
            workspaceId: currentUser.workspace_id,
            roles: currentUser.roles,
            lastLoginAt: currentUser.last_login_at
          })
        } catch (error) {
          // 令牌可能已过期，尝试刷新
          const refreshed = await refreshToken()
          if (!refreshed) {
            // 刷新失败，清除状态
            await logout()
          }
        }
      } catch (error) {
        console.error('初始化用户状态失败:', error)
        await logout()
      }
    }
  }

  return {
    // 状态
    token,
    userInfo,
    isLoggedIn,
    username,
    roles,
    
    // 操作
    setToken,
    setUserInfo,
    logout,
    hasRole,
    hasAnyRole,
    refreshToken,
    init
  }
})
