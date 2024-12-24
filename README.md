# 🚀 SeaTunnel 一键安装指南
SeaTunnel 是一个高性能、分布式的数据集成平台，支持实时和批量数据同步。本指南将帮助您快速完成 SeaTunnel 的 Zeta 集群安装部署。
Flink/Spark 模式请自行适配。
## 支持版本
| 版本 | 状态 |
|------|------|
| 2.3.8 | ✅ 已测试 |
| 2.3.7 | ✅ 已测试 |
| 2.3.6 | ✅ 已测试 |

## 兼容系统
| 操作系统 | 版本 | 状态 |
|----------|------|------|
| CentOS | 7.4+ | ✅ 已验证 |
| Rocky Linux | 9.1+ | ✅ 已验证 |
| Ubuntu | 20.04+ | 🚧 理论可行，未验证 |
| Debian | 11+ | 🚧 理论可行，未验证 |
| OpenEuler | 20.03+ | 🚧 理论可行，未验证 |
| 银河麒麟 | V10(sp1,sp2,sp3) | 🚧 理论可行，未验证 |
| 深度 | V20+ | 🚧 理论可行，未验证 |
| 统信 | V20+ | 🚧 理论可行，未验证 |

> 💡 **特别说明**：
> - 本安装指南经过严格测试和验证
> - 提供完整的部署流程和配置说明
> - 支持单节点和集群模式安装
> - 内置常用连接器和最佳实践配置

[![Apache License 2.0](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

## 项目初衷

这个一键安装工具的设计初衷是:

1. 🎯 **降低使用门槛**
   - 面向小白用户,提供最简单的部署方式
   - 自动处理各种依赖和配置,避免繁琐的手动设置
   - 提供清晰的中文提示和引导

2. 🚀 **快速体验新版本**
   - 让开发者能快速部署和体验最新版SeaTunnel
   - 便于评估是否需要升级现有环境
   - 支持多种部署模式,方便测试验证

3. 💡 **简化集群部署**
   - 自动化处理集群配置和节点分发
   - 内置最佳实践配置
   - 提供完整的部署检查和验证

4. 🛠 **开箱即用**
   - 预置常用连接器和依赖
   - 自动配置开机自启
   - 提供完整的运维命令

> 💡 提示：本工具特别适合以下场景:
> - 快速搭建测试/开发环境
> - 评估新版本特性
> - 临时部署验证概念
> - 学习和熟悉SeaTunnel



## 目录

- [快速部署](#快速部署)
- [✨ 功能特性](#-功能特性)
- [📦 快速开始](#-快速开始)
- [⚙️ 配置说明](#️-配置说明)
- [🔄 启动命令](#-启动命令)
- [🔌 端口配置](#-端口配置)
- [🔧 部署模式](#-部署模式)
- [📂 插件管理](#-插件管理)
- [🚀 开机自启动](#-开机自启动)
- [💫 安装模式](#-安装模式)
- [🔄 部署模式](#-部署模式)
- [🛡️ 安全配置](#-安全配置)
- [🔍 系统检查](#-系统检查)
- [❓ 常见问题](#-常见问题)
- [💡 获取帮助](#-获取帮助)
- [📦 下一步](#-下一步)
- [🤝 贡献](#-贡献)

## 快速部署

### 0. 用户权限配置

#### 方式一：使用root用户安装
```bash
# 直接使用root用户执行安装脚本即可
sudo su -
./install_seatunnel.sh
```

#### 方式二：使用普通用户 + sudo权限安装（推荐）

1. 创建安装用户和用户组
```bash
# 创建用户组
sudo groupadd seatunnel

# 创建用户并加入用户组
sudo useradd -m -g seatunnel seatunnel

# 设置用户密码
sudo passwd seatunnel
```

2. 配置sudo权限
```bash
# 创建sudo权限配置文件
sudo tee /etc/sudoers.d/seatunnel << EOF
Defaults:seatunnel !authenticate
seatunnel ALL=(ALL:ALL) NOPASSWD: ALL
EOF

# 设置正确的权限
sudo chmod 440 /etc/sudoers.d/seatunnel
```

3. 切换到安装用户
```bash
# 切换用户
su - seatunnel

# 验证sudo权限
sudo whoami  # 应该输出 root
sudo ls /root  # 应该能访问root目录
sudo systemctl status  # 应该能执行系统管理命令 （如果报错，请检查是否安装了systemctl）
```

4. 修改config.properties中的用户配置
```properties
# 设置安装用户和用户组
INSTALL_USER=seatunnel
INSTALL_GROUP=seatunnel
```

> ⚠️ 注意：
> - 建议使用普通用户 + sudo权限的方式安装
> - 安装用户必须具有sudo权限
> - 如果使用root用户安装，脚本会给出警告提示
> - 安装目录的所有者会被设置为INSTALL_USER:INSTALL_GROUP
> - !!!安装脚本会自动禁用SELinux，以避免权限问题导致的各种错误

### 1. 单节点安装(默认root用户)

#### 方式一：GitHub下载（国外推荐）
```bash
# 第一步：下载并解压
curl -s https://api.github.com/repos/LeonYoah/seatunnel-installer/releases/latest | grep "tag_name" | cut -d '"' -f 4 | xargs -I {} sh -c 'mkdir -p ~/seatunnel-installer && cd ~/seatunnel-installer && wget https://github.com/LeonYoah/seatunnel-installer/archive/refs/tags/{}.tar.gz -O- | tar -xz'

# 第二步：进入目录并执行安装
cd ~/seatunnel-installer/seatunnel-installer-* && chmod +x install_seatunnel.sh

# 完整安装（含插件）
./install_seatunnel.sh

# 仅安装核心组件（不含插件）
./install_seatunnel.sh --no-plugins

# 在已有安装的seatunnel中更新插件
./install_seatunnel.sh --install-plugins
```

#### 方式二：Gitee下载（国内推荐）
```bash
# 第一步：下载并解压
curl -s https://gitee.com/api/v5/repos/lyb173/seatunnel-installer/releases/latest | grep -o '"tag_name":"[^"]*' | cut -d'"' -f4 | xargs -I {} sh -c 'mkdir -p ~/seatunnel-installer && cd ~/seatunnel-installer && wget https://gitee.com/lyb173/seatunnel-installer/repository/archive/{}.tar.gz -O- | tar -xz'

# 第二步：进入目录并执行安装
cd ~/seatunnel-installer/seatunnel-installer-* && chmod +x install_seatunnel.sh

# 完整安装（含插件）
./install_seatunnel.sh

# 仅安装核心组件（不含插件）
./install_seatunnel.sh --no-plugins

# 在已有安装的seatunnel中更新插件
./install_seatunnel.sh --install-plugins
```

> 💡 提示：
> - 默认安装目录为 `/home/seatunnel/seatunnel-package`
> - 如需修改安装目录，请编辑 config.properties 中的 BASE_DIR 配置项
> - `--no-plugins`: 仅安装核心组件，不安装任何插件
> - `--install-plugins`: 单独安装或更新插件，可用于已安装环境
> - GitHub最新版本：[![Latest Release](https://img.shields.io/github/v/release/LeonYoah/seatunnel-installer)](https://github.com/LeonYoah/seatunnel-installer/releases/latest)
> - Gitee仓库：[![Gitee](https://img.shields.io/badge/Gitee-Repository-red)](https://gitee.com/lyb173/seatunnel-installer/releases)

### 2. 多节点安装(默认root用户)

#### 2.1 配置SSH免密登录
```bash
ssh-keygen -t rsa
ssh-copy-id user@node1
ssh-copy-id user@node2
# ... 对所有节点执行
```

#### 2.2 下载并解压安装包

##### 方式一：GitHub下载（国外推荐）
```bash
# 第一步：下载并解压
curl -s https://api.github.com/repos/LeonYoah/seatunnel-installer/releases/latest | grep "tag_name" | cut -d '"' -f 4 | xargs -I {} sh -c 'mkdir -p ~/seatunnel-installer && cd ~/seatunnel-installer && wget https://github.com/LeonYoah/seatunnel-installer/archive/refs/tags/{}.tar.gz -O- | tar -xz'

# 第二步：进入目录修改config.properties
cd ~/seatunnel-installer/seatunnel-installer-* && chmod +x install_seatunnel.sh && vim config.properties
```

##### 方式二：Gitee下载（国内推荐）
```bash
# 第一步：下载并解压
curl -s https://gitee.com/api/v5/repos/lyb173/seatunnel-installer/releases/latest | grep -o '"tag_name":"[^"]*' | cut -d'"' -f4 | xargs -I {} sh -c 'mkdir -p ~/seatunnel-installer && cd ~/seatunnel-installer && wget https://gitee.com/lyb173/seatunnel-installer/repository/archive/{}.tar.gz -O- | tar -xz'

# 第二步：进入目录修改config.properties
cd ~/seatunnel-installer/seatunnel-installer-* && chmod +x install_seatunnel.sh && vim config.properties
```

#### 2.3 配置节点IP
修改 config.properties 中的以下部分：
```properties
# ==== 分离模式 ====
MASTER_IP=192.168.1.100,192.168.1.101
WORKER_IPS=192.168.1.102,192.168.1.103,192.168.1.104

# ==== 或者使用混合模式 ====
CLUSTER_NODES=192.168.1.100,192.168.1.101,192.168.1.102
```

#### 2.4 执行安装
```bash
# 完整安装（含插件）
./install_seatunnel.sh

# 仅安装核心组件（不含插件）
./install_seatunnel.sh --no-plugins

# 在已有安装的seatunnel中更新插件
./install_seatunnel.sh --install-plugins
```

### 3. 卸载 SeaTunnel
如需卸载 SeaTunnel，请执行以下命令：
```bash
# 下载卸载脚本 ，注意！！改脚本需要配合config.properties使用
cd ~/seatunnel-installer/seatunnel-installer-* && chmod +x uninstall_seatunnel.sh

# 执行卸载
./uninstall_seatunnel.sh
```

> ⚠️ 注意：
> - 卸载操作将停止所有 SeaTunnel 服务
> - 删除安装目录及所有相关文件
> - 移除系统服务配置
> - 清理环境变量设置

> 💡 提示：
> - 默认已包含常用连接器(jdbc,hive)
> - 其他配置项使用默认值，可按需调整

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

### 系统要求

#### Java环境
- 支持 Java 8 或 Java 11
- 在线安装模式下,如未安装Java会提示自动安装:

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
SEATUNNEL_VERSION=2.3.7
INSTALL_MODE=offline
BASE_DIR=/data/seatunnel

# ==== 可选配置 ====
DEPLOY_MODE=separated
INSTALL_USER=root
INSTALL_GROUP=root
```

## 🔄 启动命令

### 手动启动

#### 混合模式 (Hybrid)/分离模式 (Separated)
```bash
# 启动集群
${BASE_DIR}/bin/seatunnel-cluster.sh start

# 停止集群
${BASE_DIR}/bin/seatunnel-cluster.sh stop

# 重启集群
${BASE_DIR}/bin/seatunnel-cluster.sh restart

# 查看日志
tail -n 100 $SEATUNNEL_HOME/logs/seatunnel-engine[-master/-worker/-server].log
```



### 使用Systemd服务

#### 混合模式
| 操作 | 命令 |
|------|------|
| 启动服务 | `sudo systemctl start seatunnel` |
| 停止服务 | `sudo systemctl stop seatunnel` |
| 重启服务 | `sudo systemctl restart seatunnel` |
| 查看状态 | `sudo systemctl status seatunnel` |
| 启用自启动 | `sudo systemctl enable seatunnel` |
| 禁用自启动 | `sudo systemctl disable seatunnel` |

#### 分离模式 - Master节点
| 操作 | 命令 |
|------|------|
| 启动服务 | `sudo systemctl start seatunnel-master` |
| 停止服务 | `sudo systemctl stop seatunnel-master` |
| 重启服务 | `sudo systemctl restart seatunnel-master` |
| 查看状态 | `sudo systemctl status seatunnel-master` |
| 启用自启动 | `sudo systemctl enable seatunnel-master` |
| 禁用自启动 | `sudo systemctl disable seatunnel-master` |

#### 分离模式 - Worker节点
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
> - 查看启动日志：`sudo journalctl -u seatunnel[-master/-worker] -n 100 --no-pager`
> - 查看运行日志：`tail -n 100 $SEATUNNEL_HOME/logs/seatunnel-engine[-master/-worker/-server].log`

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
jdbc_libs="mysql:mysql-connector-java:8.0.27","org.postgresql:postgresql:42.4.3"

# hive依赖
hive_libs="org.apache.hive:hive-exec:3.1.3","org.apache.hive:hive-service:3.1.3"
```
</details>

## 🚀 systemd管理开机自启动

### 基础配置
```properties
ENABLE_AUTO_START=true
```

## 💫 安装模式

### 在线安装
```properties
INSTALL_MODE=online
PACKAGE_REPO=aliyun
# 可选：指定下载源
DOWNLOAD_URL=https://archive.apache.org/dist/seatunnel/${VERSION}/apache-seatunnel-${VERSION}-bin.tar.gz
```

### 离线安装
```properties
INSTALL_MODE=offline
PACKAGE_PATH=apache-seatunnel-${VERSION}.tar.gz
```

### 镜像源配置
支持多种镜像源加速下载：
- Apache官方源
- 阿里云镜像
- 华为云镜像

## 🔄 部署模式


### 分离模式 (默认)
Master和Worker分开部署：
```properties
DEPLOY_MODE=separated
# Master节点
MASTER_IP=192.168.1.100,192.168.1.101
# Worker节点 
WORKER_IPS=192.168.1.102,192.168.1.103
```

### 混合模式
所有节点对等部署：
```properties
DEPLOY_MODE=hybrid
# 所有节点IP
CLUSTER_NODES=192.168.1.100,192.168.1.101,192.168.1.102
```

## 🛡️ 安全配置

### 用户权限
```properties
# 安装用户(需要sudo权限)
INSTALL_USER=root
INSTALL_GROUP=root
```

### SSH配置
```properties
# SSH端口
SSH_PORT=22
# 超时设置(秒)
SSH_TIMEOUT=10
```

### 自动重试机制
- 最大重试次数：3次
- 失败自动回滚
- 详细错误日志

## 🔍 系统检查

安装前自动检查：
- Java环境检查
- 内存要求检查
- 端口占用检查
- 依赖组件检查
- 下载源可用性检查


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


