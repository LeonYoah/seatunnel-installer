<!--
  登录页面
  实现用户登录功能，包括表单验证、JWT令牌存储、自动登录等
-->
<template>
  <div class="login-container">
    <div class="login-card">
      <!-- Logo 和标题 -->
      <div class="login-header">
        <img src="/logo.png" alt="SeaTunnel" class="logo" />
        <h1 class="title">{{ t('app.name') }}</h1>
        <p class="subtitle">{{ t('app.subtitle') }}</p>
      </div>

      <!-- 登录表单 -->
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        @submit.prevent="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            :placeholder="t('login.form.username')"
            size="large"
            prefix-icon="User"
            clearable
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            :placeholder="t('login.form.password')"
            size="large"
            prefix-icon="Lock"
            show-password
            clearable
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-form-item>
          <el-checkbox v-model="loginForm.rememberMe">
            {{ t('login.form.rememberMe') }}
          </el-checkbox>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="login-button"
            @click="handleLogin"
          >
            {{ loading ? t('login.form.loggingIn') : t('login.form.login') }}
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 语言切换 -->
      <div class="language-switcher">
        <LanguageSwitcher />
      </div>
    </div>

    <!-- 背景装饰 -->
    <div class="login-background">
      <div class="bg-shape bg-shape-1"></div>
      <div class="bg-shape bg-shape-2"></div>
      <div class="bg-shape bg-shape-3"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useUserStore } from '@/stores/user'
import LanguageSwitcher from '@/components/shared/LanguageSwitcher.vue'
import { authApi } from '@/api/auth'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()

// 表单引用
const loginFormRef = ref<FormInstance>()

// 加载状态
const loading = ref(false)

// 登录表单数据
const loginForm = reactive({
  username: '',
  password: '',
  rememberMe: false
})

// 表单验证规则
const loginRules: FormRules = {
  username: [
    { required: true, message: t('login.validation.usernameRequired'), trigger: 'blur' },
    { min: 2, max: 50, message: t('login.validation.usernameLength'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: t('login.validation.passwordRequired'), trigger: 'blur' },
    { min: 6, message: t('login.validation.passwordLength'), trigger: 'blur' }
  ]
}

// 处理登录
const handleLogin = async () => {
  if (!loginFormRef.value) return

  try {
    // 表单验证
    await loginFormRef.value.validate()
    
    loading.value = true

    // 调用登录API
    const response = await authApi.login({
      username: loginForm.username,
      password: loginForm.password
    })

    // 存储令牌和用户信息
    userStore.setToken(response.access_token)
    userStore.setUserInfo({
      id: response.user.id,
      username: response.user.username,
      email: response.user.email,
      tenantId: response.user.tenant_id,
      workspaceId: response.user.workspace_id,
      roles: response.user.roles,
      lastLoginAt: response.user.last_login_at
    })

    // 存储刷新令牌
    if (loginForm.rememberMe) {
      localStorage.setItem('refresh_token', response.refresh_token)
    } else {
      sessionStorage.setItem('refresh_token', response.refresh_token)
    }

    ElMessage.success(t('login.messages.loginSuccess'))

    // 跳转到首页
    const redirect = router.currentRoute.value.query.redirect as string
    await router.push(redirect || '/dashboard')
  } catch (error: any) {
    console.error('登录失败:', error)
    
    // 显示错误信息
    const message = error.response?.data?.message || t('login.messages.loginFailed')
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}

// 组件挂载时检查是否已登录
onMounted(() => {
  // 如果已经有令牌，直接跳转到首页
  if (userStore.token) {
    router.push('/dashboard')
  }
})
</script>

<style scoped>
.login-container {
  position: relative;
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  overflow: hidden;
}

.login-card {
  position: relative;
  z-index: 10;
  width: 400px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo {
  width: 64px;
  height: 64px;
  margin-bottom: 16px;
}

.title {
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.subtitle {
  font-size: 14px;
  color: #7f8c8d;
  margin: 0;
}

.login-form {
  margin-bottom: 24px;
}

.login-form .el-form-item {
  margin-bottom: 20px;
}

.login-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 500;
}

.language-switcher {
  display: flex;
  justify-content: center;
}

/* 背景装饰 */
.login-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 1;
}

.bg-shape {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 6s ease-in-out infinite;
}

.bg-shape-1 {
  width: 200px;
  height: 200px;
  top: 10%;
  left: 10%;
  animation-delay: 0s;
}

.bg-shape-2 {
  width: 150px;
  height: 150px;
  top: 60%;
  right: 10%;
  animation-delay: 2s;
}

.bg-shape-3 {
  width: 100px;
  height: 100px;
  bottom: 20%;
  left: 20%;
  animation-delay: 4s;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
  }
}

/* 响应式设计 */
@media (max-width: 480px) {
  .login-card {
    width: 90%;
    padding: 24px;
    margin: 20px;
  }
  
  .title {
    font-size: 24px;
  }
  
  .login-button {
    height: 44px;
    font-size: 14px;
  }
}

/* 暗色主题适配 */
@media (prefers-color-scheme: dark) {
  .login-card {
    background: rgba(30, 30, 30, 0.95);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }
  
  .title {
    color: #ffffff;
  }
  
  .subtitle {
    color: #a0a0a0;
  }
}
</style>