/**
 * HTTP请求工具
 * 配置axios拦截器，处理JWT令牌、自动刷新、错误处理等
 */
import axios, { type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

// 创建axios实例
const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 是否正在刷新令牌
let isRefreshing = false
// 等待刷新的请求队列
let failedQueue: Array<{
  resolve: (value?: any) => void
  reject: (reason?: any) => void
}> = []

// 处理队列中的请求
const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) {
      reject(error)
    } else {
      resolve(token)
    }
  })
  
  failedQueue = []
}

// 请求拦截器
http.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 添加JWT令牌到请求头
    const token = localStorage.getItem('token')
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }
    
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
http.interceptors.response.use(
  (response: AxiosResponse) => {
    return response
  },
  async (error) => {
    const originalRequest = error.config
    
    // 如果是401错误且不是登录或刷新令牌请求
    if (error.response?.status === 401 && !originalRequest._retry) {
      // 如果是登录或刷新令牌请求失败，直接返回错误
      if (originalRequest.url?.includes('/auth/login') || originalRequest.url?.includes('/auth/refresh')) {
        return Promise.reject(error)
      }
      
      // 如果正在刷新令牌，将请求加入队列
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        }).then(() => {
          return http(originalRequest)
        }).catch(err => {
          return Promise.reject(err)
        })
      }
      
      originalRequest._retry = true
      isRefreshing = true
      
      try {
        // 尝试刷新令牌
        const { useUserStore } = await import('@/stores/user')
        const userStore = useUserStore()
        
        const refreshed = await userStore.refreshToken()
        
        if (refreshed) {
          // 刷新成功，处理队列中的请求
          processQueue(null, userStore.token)
          // 重试原始请求
          return http(originalRequest)
        } else {
          // 刷新失败，跳转到登录页
          processQueue(error, null)
          await router.push('/login')
          return Promise.reject(error)
        }
      } catch (refreshError) {
        // 刷新失败，清除认证信息并跳转到登录页
        processQueue(refreshError, null)
        const { useUserStore } = await import('@/stores/user')
        const userStore = useUserStore()
        await userStore.logout()
        await router.push('/login')
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }
    
    // 处理其他错误
    if (error.response) {
      const { status, data } = error.response
      
      switch (status) {
        case 403:
          ElMessage.error('权限不足')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          if (data?.message) {
            ElMessage.error(data.message)
          } else {
            ElMessage.error(`请求失败 (${status})`)
          }
      }
    } else if (error.request) {
      ElMessage.error('网络连接失败，请检查网络设置')
    } else {
      ElMessage.error('请求配置错误')
    }
    
    return Promise.reject(error)
  }
)

export default http