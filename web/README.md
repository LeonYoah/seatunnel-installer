# SeaTunnel 企业级管理平台 - 前端项目

## 项目简介

基于 Vue 3 + TypeScript + Vite + Element Plus 构建的现代化企业级管理平台前端。

## 技术栈

- **框架**: Vue 3.5+ (Composition API)
- **语言**: TypeScript 5.9+
- **构建工具**: Vite 7.2+
- **UI 组件库**: Element Plus 2.5+
- **路由**: Vue Router 4.2+
- **国际化**: vue-i18n（简体中文/English，默认简体中文）
- **状态管理**: Pinia 2.1+
- **HTTP 客户端**: Axios 1.6+
- **代码规范**: ESLint + Prettier

## 项目结构

```
web/
├── public/              # 静态资源
│   └── logo.png        # Logo 图片
├── src/
│   ├── assets/         # 资源文件
│   ├── components/     # 组件
│   │   └── layout/     # 布局组件
│   │       ├── Header.vue      # 头部组件
│   │       ├── Sidebar.vue     # 侧边栏组件
│   │       ├── Footer.vue      # 底部组件
│   │       └── MainLayout.vue  # 主布局组件
│   ├── router/         # 路由配置
│   │   └── index.ts    # 路由定义
│   ├── stores/         # 状态管理
│   │   ├── index.ts    # Pinia 入口
│   │   └── user.ts     # 用户状态
│   ├── views/          # 页面视图
│   │   ├── Dashboard.vue    # 总览页面
│   │   ├── Install.vue      # 安装向导
│   │   ├── Tasks.vue        # 任务管理
│   │   ├── Clusters.vue     # 集群管理
│   │   ├── Diagnostics.vue  # 诊断中心
│   │   ├── Plugins.vue      # 插件市场
│   │   └── Settings.vue     # 设置
│   ├── App.vue         # 根组件
│   ├── main.ts         # 应用入口
│   └── style.css       # 全局样式
├── .eslintrc.cjs       # ESLint 配置
├── .prettierrc.json    # Prettier 配置
├── vite.config.ts      # Vite 配置
├── tsconfig.json       # TypeScript 配置
└── package.json        # 项目依赖

```

## 开发指南

### 安装依赖

```bash
cd web
npm install
```

### 启动开发服务器

```bash
npm run dev
```

访问 http://localhost:5173

### 国际化（i18n）

- 默认语言：简体中文（`zh-CN`）。
- 切换语言：右上角语言下拉框，或 `localStorage.locale`。
- 语言文件：`src/locales/zh-CN.ts`, `src/locales/en.ts`。
- i18n 初始化：`src/i18n/index.ts`。

### 主题与品牌色

- 默认主题：明亮（白/蓝）。可在右上角切换暗色主题。
- 品牌蓝：运行时根据 `public/logo.png` 采样，自动对齐到 `--primary` 并同步 Element Plus 主色。
- 覆盖变量：见 `src/styles/theme.scss`（`--bg/--surface/--text/--primary` 等）。

### 构建生产版本

```bash
npm run build
```

### 预览生产构建

```bash
npm run preview
```

## 功能模块

### 已实现

#### 基础框架
- ✅ Vue 3 + TypeScript 项目初始化
- ✅ Vite 构建工具配置
- ✅ Element Plus UI 库集成
- ✅ ESLint + Prettier 代码规范
- ✅ Vue Router 路由框架
- ✅ Pinia 状态管理

#### 布局和主题
- ✅ 基础布局组件（Header、Sidebar、Footer、MainLayout）
- ✅ 明亮/暗色主题切换
- ✅ 主题持久化
- ✅ 完整的暗色主题适配（参考 www/console）
- ✅ 路由守卫（页面标题、权限检查预留）
- ✅ 白/蓝主题（默认白色），品牌蓝与 logo 自动对齐

#### 页面功能
- ✅ **Dashboard 总览页面** - KPI 卡片、趋势图表、告警列表、任务列表
- ✅ **Tasks 任务管理** - 任务列表、搜索筛选、分页、操作按钮
- ✅ **Clusters 集群管理** - 集群列表、节点管理、资源监控、扩缩容
- ✅ **Diagnostics 诊断中心** - 一键诊断、故障库、诊断历史
- ✅ **Plugins 插件市场** - 插件列表、安装管理、搜索筛选
- ✅ **Settings 设置** - 通知配置、系统配置、用户管理、审计日志
- ✅ **Install 安装向导** - 页面占位

#### 数据展示
- ✅ Mock 数据填充
- ✅ 表格、卡片、图表占位
- ✅ 交互反馈（消息提示、确认框）

### 待实现

- ⏳ 安装向导完整流程
- ⏳ 任务创建和编辑（YAML/表单双模式）
- ⏳ 任务版本管理和回滚
- ⏳ 实时日志查看（WebSocket）
- ⏳ ECharts 图表集成
- ⏳ 用户认证与授权
- ⏳ API 接口集成
- ⏳ WebSocket 实时通信
- ⏳ 国际化支持（中英文）

## 路由配置

| 路径 | 组件 | 说明 |
|------|------|------|
| / | Dashboard | 重定向到总览 |
| /dashboard | Dashboard | 总览页面 |
| /install | Install | 安装向导 |
| /tasks | Tasks | 任务管理 |
| /clusters | Clusters | 集群管理 |
| /diagnostics | Diagnostics | 诊断中心 |
| /plugins | Plugins | 插件市场 |
| /settings | Settings | 设置 |

## 代码规范

### 组件命名

- 组件文件使用 PascalCase：`MyComponent.vue`
- 组件名称使用 PascalCase：`<MyComponent />`

### 样式规范

- 使用 scoped 样式避免污染全局
- 使用 Element Plus 的设计规范
- 自定义样式使用 BEM 命名规范

### TypeScript 规范

- 启用严格模式
- 为函数参数和返回值添加类型注解
- 避免使用 `any`，使用 `unknown` 或具体类型

## API 代理配置

开发环境下，API 请求会被代理到后端服务：

```
/api/* -> http://localhost:8080/api/*
```

## 注意事项

1. 本项目使用 Vue 3 Composition API，请熟悉相关语法
2. 所有代码注释和文档必须使用中文
3. 提交代码前请运行 ESLint 检查
4. 遵循 Git 提交规范

## 参考资料

- [Vue 3 文档](https://cn.vuejs.org/)
- [Element Plus 文档](https://element-plus.org/zh-CN/)
- [Vite 文档](https://cn.vitejs.dev/)
- [Pinia 文档](https://pinia.vuejs.org/zh/)
- [Vue Router 文档](https://router.vuejs.org/zh/)
