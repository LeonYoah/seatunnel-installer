# 🚀 SeaTunnel 一键安装指南

[![Apache License 2.0](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

SeaTunnel 是一个高性能、分布式的数据集成平台，支持实时和批量数据同步。本指南将帮助您快速完成 SeaTunnel 的 Zeta 集群安装部署。
Flink/Spark 模式请自行适配。


## 目录

- [快速部署](#快速部署)
  * [1. 准备安装目录](#1-准备安装目录)
  * [2. 配置SSH免密登录](#2-配置ssh免密登录)
  * [3. 配置节点IP](#3-配置节点ip)
  * [4. 执行安装](#4-执行安装)
- [✨ 功能特性](#-功能特性)
- [📦 快速开始](#-快速开始)
- [⚙️ 配置说明](#️-配置说明)
- [🔄 启动命令](#-启动命令)
- [🔌 端口配置](#-端口配置)
- [🔧 部署模式](#-部署模式)
- [📂 插件管理](#-插件管理)
- [🚀 开机自启动](#-开机自启动)
- [❓ 常见问题](#-常见问题)
- [💡 获取帮助](#-获取帮助)
- [🤝 贡献](#-贡献)

## 快速部署

### 1. 准备安装目录
```bash
# 创建安装目录
mkdir -p ~/seatunnel-installer && cd ~/seatunnel-installer

# 下载安装脚本和配置文件
wget https://github.com/LeonYoah/seatunnel-installer/raw/main/install_seatunnel.sh
wget https://github.com/LeonYoah/seatunnel-installer/raw/main/config.properties

# 添加执行权限
chmod +x install_seatunnel.sh
```

> 💡 提示：
> - 默认安装目录为 `/data/seatunnel`
> - 如需修改安装目录，请编辑 config.properties 中的 BASE_DIR 配置项
> ```properties
> # 修改为你想要的安装目录
> BASE_DIR=/your/custom/path/seatunnel
> ```

### 2. 配置SSH免密登录
```bash
# 在所有节点间配置SSH免密登录
ssh-keygen -t rsa  # 如果已经有密钥对，可以跳过
ssh-copy-id user@node1
ssh-copy-id user@node2
# ... 对所有节点执行
```

### 3. 配置节点IP（默认是localhost）
只需修改config.properties中的以下部分：
```properties
# ==== 分离模式 ====
# Master节点IP
MASTER_IP=192.168.1.100,192.168.1.101
# Worker节点IP
WORKER_IPS=192.168.1.102,192.168.1.103,192.168.1.104

# ==== 或者使用混合模式 ====
# 所有节点IP
CLUSTER_NODES=192.168.1.100,192.168.1.101,192.168.1.102
```

### 4. 执行安装
```bash
./install_seatunnel.sh
```

> 💡 提示：
> - 默认已包含常用连接器(jdbc,hive)
> - 其他配置项使用默认值，可按需调整
> - 详细配置说明请继续往下阅读

### ⚠️ 重要提醒：分布式部署必读
如果您正在部署分布式集群（多节点部署），请选择合适的配置分布式存储作为checkpoint存储，否则将影响以下功能：
- 流式处理连接器（如：Kafka）无法正常运行
- CDC连接器(如：ORACLE-CDC)的断点续传功能无法使用

推荐配置以下任一存储：
```properties
# ==== 在config.properties中配置 ====

# 方式1：配置HDFS（推荐）
CHECKPOINT_STORAGE_TYPE=HDFS
CHECKPOINT_NAMESPACE=/seatunnel/checkpoint
HDFS_NAMENODE_HOST=hdfs-namenode-host
HDFS_NAMENODE_PORT=8020

# 方式2：配置OSS
CHECKPOINT_STORAGE_TYPE=OSS
CHECKPOINT_NAMESPACE=/seatunnel/checkpoint
STORAGE_ENDPOINT=http://oss-cn-hangzhou.aliyuncs.com
STORAGE_ACCESS_KEY=your_access_key
STORAGE_SECRET_KEY=your_secret_key
STORAGE_BUCKET=your_bucket

# 方式3：配置S3
CHECKPOINT_STORAGE_TYPE=S3
CHECKPOINT_NAMESPACE=/seatunnel/checkpoint
STORAGE_ENDPOINT=http://s3.amazonaws.com
STORAGE_ACCESS_KEY=your_access_key
STORAGE_SECRET_KEY=your_secret_key
STORAGE_BUCKET=your_bucket
```

> ⚠️ 注意：默认的LOCAL_FILE存储模式只适用于单节点测试环境，不建议在生产环境使用。

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

## 📦 快速开始

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

## 🔄 启动命令

### 手动启动

#### 混合模式 (Hybrid)/分离模式 (Separated)
```bash
# 启动集群
${SEATUNNEL_HOME}/bin/seatunnel-start-cluster.sh start


# 停止集群
${SEATUNNEL_HOME}/bin/seatunnel-start-cluster.sh stop

# 启动/停止/重启集群
${SEATUNNEL_HOME}/bin/seatunnel-start-cluster.sh restart
```

### 使用Systemd服务

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

## 🔌 端口配置

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

## 🔧 部署模式

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

## 📂 插件管理

### 高级配置

<details>
<summary>点击展开完整配置示例</summary>

```properties
# ==== 连接器配置 ====
CONNECTORS=jdbc,hive

# JDBC依赖
jdbc_libs=(
    "mysql:mysql-connector-java:8.0.27"
    "org.postgresql:postgresql:42.4.3"
)

# hive依赖
hive_libs=(
    "org.apache.hive:hive-exec:3.1.3"
    "org.apache.hive:hive-service:3.1.3"
)
```
</details>

## 🚀 开机自启动

### 基础配置
```properties
ENABLE_AUTO_START=true
```



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

## 💡 获取帮助

- 📖 [官方文档](https://seatunnel.apache.org/docs)
- 🐛 [问题反馈](https://github.com/apache/seatunnel/issues)
- 💬 [社区支持](https://slack.seatunnel.apache.org/)

## 📦 下一步

- [配置数据源](https://seatunnel.apache.org/docs/connector-v2/source)
- [配置数据目标](https://seatunnel.apache.org/docs/connector-v2/sink)
- [开发自定义连接器](https://seatunnel.apache.org/docs/development/connector-v2)

## 🤝 贡献

欢迎提交Issue和Pull Request来帮助改进这个安装器！
