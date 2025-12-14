/**
 * 用户状态管理
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref<string>('')
  const username = ref<string>('')
  const roles = ref<string[]>([])

  // 操作
  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUserInfo = (info: { username: string; roles: string[] }) => {
    username.value = info.username
    roles.value = info.roles
  }

  const logout = () => {
    token.value = ''
    username.value = ''
    roles.value = []
    localStorage.removeItem('token')
  }

  // 初始化
  const init = () => {
    const savedToken = localStorage.getItem('token')
    if (savedToken) {
      token.value = savedToken
    }
  }

  return {
    token,
    username,
    roles,
    setToken,
    setUserInfo,
    logout,
    init
  }
})
