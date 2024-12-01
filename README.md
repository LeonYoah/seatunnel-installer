# 🚀 SeaTunnel 一键安装指南

[![Apache License 2.0](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

SeaTunnel 是一个高性能、分布式的数据集成平台，支持实时和批量数据同步。本指南将帮助您快速完成 SeaTunnel 的 Zeta 集群安装部署。
Flink/Spark 模式请自行适配。

## 📑 目录

- [功能特性](#-功能特性)
- [快速开始](#-快速开始)
- [环境要求](#-环境要求)
- [部署模式](#-部署模式)
- [插件管理](#-插件管理)
- [开机自启动](#-开机自启动)
- [常见问题](#-常见问题)
- [获取帮助](#-获取帮助)
- [端口配置说明](#-端口配置说明)

## ✨ 功能特性

> 相比官方安装方式，本安装器提供了全方位的增强功能

### 1️⃣ 健壮性增强
- 🛡️ SSH/SCP操作重试机制
- 📝 增强的错误处理和日志记录
- 🔍 自动检测和验证系统依赖
- ✅ 安装包完整性校验

### 2️⃣ 用户权限管理
- 👤 基于配置的用户安装
- 🔐 自动创建和配置用户权限
- 📂 合理的文件权限设置
- 👥 多用户环境支持

### 3️⃣ 集群管理增强
- 🎮 统一的集群管理脚本
- 🔄 支持混合模式和分离模式
- 🚀 自动节点配置和分发
- 📊 集群状态检查和监控

### 4️⃣ 依赖管理优化

#### 4.1 智能下载机制
| 下载源 | 说明 | 优先级 |
|-------|------|--------|
| 阿里云 | 国内推荐 | 1 |
| 中央仓库 | 自动备选 | 2 |
| 华为云 | 可选配置 | 3 |
| 自定义 | 支持私有仓库 | - |

#### 4.2 预置连接器
> 默认集成常用连接器及其依赖，开箱即用

- JDBC系列
  * MySQL
  * PostgreSQL
  * Oracle
  * 达梦
  * 虚谷
  * 人大金仓
- 大数据生态
  * Hive

#### 4.3 依赖管理特性
- 📦 统一的lib目录管理
- 🔄 支持增量安装
- 🛠️ 灵活的版本配置
- 📥 智能重试机制

### 5️⃣ 使用体验优化
- 📊 详细的进度展示
- ❌ 清晰的错误提示
- ✅ 完整的安装检查
- 📚 丰富的使用文档

## 🚀 快速开始

### 一键安装

```bash
./install_seatunnel.sh
```

> 💡 提示：安装默认自带jdbc和hive连接器及依赖

### 常用命令

```bash
# 完整安装（含插件）
./install_seatunnel.sh

# 仅安装核心组件
./install_seatunnel.sh --no-plugins

# 单独安装/更新插件
./install_seatunnel.sh --install-plugins
```

## ⚙️ 配置说明

### 基础配置

```properties
# ==== 必选配置 ====
SEATUNNEL_VERSION=2.3.7      # 版本号
INSTALL_MODE=offline         # 安装模式(online/offline)
BASE_DIR=/data/seatunnel    # 安装目录

# ==== 可选配置 ====
DEPLOY_MODE=separated        # 部署模式(separated/hybrid)
INSTALL_USER=root           # 安装用户
INSTALL_GROUP=root          # 安装用户组
```

### 部署模式

#### 混合模式 (Hybrid)
> 适合小规模部署，配置简单

- ✅ 所有节点运行相同组件
- ✅ 维护成本低
- ❗ 资源隔离性差

#### 分离模式 (Separated)
> 适合生产环境，资源隔离好

- ✅ Master/Worker分离部署
- ✅ 更好的扩展性
- ✅ 资源利用更合理

## 🔌 插件管理

### 快速配置

```properties
# ==== 最小配置 ====
INSTALL_CONNECTORS=true
CONNECTORS=jdbc,hive

# ==== 自定义配置 ====
jdbc_libs=(
    "mysql:mysql-connector-java:8.0.27"
    "org.postgresql:postgresql:42.4.3"
)
```

### 高级配置

<details>
<summary>点击展开完整配置示例</summary>

```properties
# ==== 下载源配置 ====
MAVEN_REPO=aliyun
# CUSTOM_MAVEN_REPO=https://your-repo.com

# ==== 连接器配置 ====
CONNECTORS=jdbc,kafka,elasticsearch

# JDBC依赖
jdbc_libs=(
    "mysql:mysql-connector-java:8.0.27"
    "org.postgresql:postgresql:42.4.3"
)

# Kafka依赖
kafka_libs=(
    "org.apache.kafka:kafka-clients:3.2.3"
)
```
</details>

## 🔄 开机自启动

### 基础配置
```properties
ENABLE_AUTO_START=true
AUTO_START_DELAY=60
```

### 服务管理

| 操作 | 命令 |
|------|------|
| 启动 | `sudo systemctl start seatunnel` |
| 停止 | `sudo systemctl stop seatunnel` |
| 重启 | `sudo systemctl restart seatunnel` |
| 状态 | `sudo systemctl status seatunnel` |
| 禁用 | `sudo systemctl disable seatunnel` |

## ❓ 常见问题

<details>
<summary>1. 安装失败如何处理？</summary>

- 检查安装日志
- 确认环境要求
- 验证网络连接
- 检查用户权限
</details>

<details>
<summary>2. 插件安装失败？</summary>

- 确认Maven仓库可访问
- 检查依赖配置正确性
- 尝试切换下载源
- 查看详细错误日志
</details>

<details>
<summary>3. 服务启动失败？</summary>

- 检查端口占用
- 验证配置文件
- 确认权限正确
- 查看系统日志
</details>

## 🆘 获取帮助

- 📖 [官方文档](https://seatunnel.apache.org/docs)
- 🐛 [问题反馈](https://github.com/apache/seatunnel/issues)
- 💬 [社区支持](https://slack.seatunnel.apache.org/)

## 📦 下一步

- [配置数据源](https://seatunnel.apache.org/docs/connector-v2/source)
- [配置数据目标](https://seatunnel.apache.org/docs/connector-v2/sink)
- [开发自定义连接器](https://seatunnel.apache.org/docs/development/connector-v2)

## 🤝 贡献

欢迎提交Issue和Pull Request来帮助改进这个安装器！

## 端口配置说明

SeaTunnel安装器支持两种部署模式的端口配置：

### 混合模式端口配置
在混合模式下，所有节点使用相同的端口：
- 默认服务端口：5801
- 配置示例：
```properties
HYBRID_PORT=5801
```

### 分离模式端口配置
在分离模式下，Master和Worker节点使用不同的端口：
- Master节点默认端口：5801
- Worker节点默认端口：5802
- 配置示例：
```properties
MASTER_PORT=5801
WORKER_PORT=5802
```

### 端口配置注意事项
1. 确保配置的端口未被其他服务占用
2. 如果使用防火墙，需要开放相应端口
3. 集群内所有节点的端口配置必须一致
4. 可以在config.properties中自定义端口，如未配置将使用默认值

## 🔄 启动命令

#### 手动启动

##### 混合模式 (Hybrid)/分离模式 (Separated)
```bash
# 启动所有节点
${SEATUNNEL_HOME}/bin/seatunnel-start-cluster.sh start

# 停止所有节点
${SEATUNNEL_HOME}/bin/seatunnel-start-cluster.sh stop

# 重启所有节点
${SEATUNNEL_HOME}/bin/seatunnel-start-cluster.sh restart
```


#### Systemd服务管理

##### 混合模式
| 操作 | 命令 |
|------|------|
| 启动服务 | `sudo systemctl start seatunnel` |
| 停止服务 | `sudo systemctl stop seatunnel` |
| 重启服务 | `sudo systemctl restart seatunnel` |
| 查看状态 | `sudo systemctl status seatunnel` |
| 启用自启动 | `sudo systemctl enable seatunnel` |
| 禁用自启动 | `sudo systemctl disable seatunnel` |

##### 分离模式 - Master节点
| 操作 | 命令 |
|------|------|
| 启动服务 | `sudo systemctl start seatunnel-master` |
| 停止服务 | `sudo systemctl stop seatunnel-master` |
| 重启服务 | `sudo systemctl restart seatunnel-master` |
| 查看状态 | `sudo systemctl status seatunnel-master` |
| 启用自启动 | `sudo systemctl enable seatunnel-master` |
| 禁用自启动 | `sudo systemctl disable seatunnel-master` |

##### 分离模式 - Worker节点
| 操作 | 命令 |
|------|------|
| 启动服务 | `sudo systemctl start seatunnel-worker` |
| 停止服务 | `sudo systemctl stop seatunnel-worker` |
| 重启服务 | `sudo systemctl restart seatunnel-worker` |
| 查看状态 | `sudo systemctl status seatunnel-worker` |
| 启用自启动 | `sudo systemctl enable seatunnel-worker` |
| 禁用自启动 | `sudo systemctl disable seatunnel-worker` |

> 💡 提示：
> - 服务管理需要sudo权限
> - 服务配置文件位于 `/etc/systemd/system/` 目录
> - 修改配置后需要重新加载：`sudo systemctl daemon-reload`
> - 查看日志：`sudo journalctl -u seatunnel[-master/-worker]`
