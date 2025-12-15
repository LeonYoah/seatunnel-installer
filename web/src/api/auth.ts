/**
 * 认证相关API
 */
import http from '@/utils/http'

// 登录请求接口
export interface LoginRequest {
  username: string
  password: string
}

// 登录响应接口
export interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_at: string
  user: {
    id: string
    username: string
    email: string
    tenant_id: string
    workspace_id: string
    roles: string[]
    last_login_at?: string
  }
}

// 刷新令牌请求接口
export interface RefreshTokenRequest {
  refresh_token: string
}

// 用户信息接口
export interface UserInfo {
  id: string
  username: string
  email: string
  tenant_id: string
  workspace_id: string
  roles: string[]
  last_login_at?: string
}

// API响应格式
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 认证API类
class AuthAPI {
  private baseURL = '/api/v1/auth'

  /**
   * 用户登录
   */
  async login(request: LoginRequest): Promise<LoginResponse> {
    const response = await http.post<ApiResponse<LoginResponse>>(
      `${this.baseURL}/login`,
      request
    )
    return response.data.data
  }

  /**
   * 刷新访问令牌
   */
  async refreshToken(request: RefreshTokenRequest): Promise<LoginResponse> {
    const response = await http.post<ApiResponse<LoginResponse>>(
      `${this.baseURL}/refresh`,
      request
    )
    return response.data.data
  }

  /**
   * 获取当前用户信息
   */
  async getCurrentUser(): Promise<UserInfo> {
    const response = await http.get<ApiResponse<UserInfo>>(`${this.baseURL}/me`)
    return response.data.data
  }

  /**
   * 用户登出
   */
  async logout(): Promise<void> {
    await http.post(`${this.baseURL}/logout`)
  }
}

// 导出API实例
export const authApi = new AuthAPI()