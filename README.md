# 🚀 SeaTunnel 一键安装指南
SeaTunnel 是一个高性能、分布式的数据集成平台，支持实时和批量数据同步。本指南将帮助您快速完成 SeaTunnel 的 Zeta 集群安装部署。
Flink/Spark 模式请自行适配。
## 支持版本
| 版本 | 状态 |
|------|------|
| 2.3.12 | ✅ 已测试 |
| 2.3.11 | ✅ 已测试 |
| 2.3.10 | ✅ 已测试 |
| 2.3.9 | ✅ 已测试 |
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


[![Apache License 2.0](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

## 项目初衷

这个一键安装工具的目标很简单：

- 降低 SeaTunnel 安装和配置门槛，快速拉起一个可用集群
- 支持单节点 / 多节点、混合 / 分离等多种部署模式，方便体验和测试
- 内置常用连接器、systemd 服务和运维脚本，开箱即用

## 快速开始

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
curl -s https://gitee.com/api/v5/repos/lyb173/seatunnel-installer/releases/latest | grep -o '"tag_name":"[^\"]*' | cut -d'"' -f4 | xargs -I {} sh -c 'mkdir -p ~/seatunnel-installer && cd ~/seatunnel-installer && wget https://gitee.com/lyb173/seatunnel-installer/repository/archive/{}.tar.gz -O- | tar -xz'

# 第二步：进入目录并执行安装
cd ~/seatunnel-installer/seatunnel-installer-* && chmod +x install_seatunnel.sh

# 完整安装（含插件）
./install_seatunnel.sh

# 仅安装核心组件（不含插件）
./install_seatunnel.sh --no-plugins

# 在已有安装的seatunnel中更新插件
./install_seatunnel.sh --install-plugins
```

> 如果你已经手动下载并进入安装目录，也可以直接执行：
> 
> ```bash
> chmod +x install_seatunnel.sh
> ./install_seatunnel.sh          # 使用 config.properties 中的配置
> ```

### 2. Web 安装向导（推荐）

本仓库内置了一个 Web 安装向导，适合希望通过页面一步步完成配置和安装的场景。

```bash
chmod +x start_web.sh
./start_web.sh start              # 默认端口启动



## 其他参数
./start_web.sh -p 9000 start      # 端口9000启动
./start_web.sh -c start           # 清理后启动
./start_web.sh -c -p 9000 start   # 清理后端口9000启动
./start_web.sh stop               # 停止
./start_web.sh clean              # 仅清理临时文件
```

启动成功后终端会输出类似信息：

```text
============================================
SeaTunnel Web 安装向导已启动!
============================================

访问: http://<当前机器IP>:8888
CLI:  ./install_seatunnel.sh --help
```

在浏览器中访问上述地址：

- 在「安装配置」页填写基础信息（安装目录、部署模式、节点 IP、安装模式 online/offline 等）
- 支持在线/离线安装、混合/分离部署、HDFS/OSS/S3 检查点存储等配置
- 点击「保存配置并开始安装」，右侧日志区域实时展示安装过程
- 每个步骤都有状态和操作按钮，可以单步执行、重试、从指定步骤继续等

示例界面截图：

![SeaTunnel Web 安装向导](image/install01.png)

> 提示：Web 安装向导本质上还是调用同一个 `install_seatunnel.sh`，只是通过页面帮你编辑 `config.properties` 并按步骤执行。

---

## 常用配置示例（config.properties）

`config.properties` 是所有安装方式的唯一配置入口，以下是一个典型示例（仅保留常用项）：

```properties
SEATUNNEL_VERSION=2.3.12

# 安装模式
INSTALL_MODE=online          # online / offline
PACKAGE_PATH=/path/to/apache-seatunnel-${SEATUNNEL_VERSION}-bin.tar.gz

# 安装目录
BASE_DIR=/home/seatunnel/seatunnel-package

# 部署模式
DEPLOY_MODE=separated        # separated / hybrid

# 分离模式节点
MASTER_IP=192.168.102.101
WORKER_IPS=192.168.102.102

# 混合模式节点
CLUSTER_NODES=192.168.102.101,192.168.102.102

# 端口配置
HYBRID_PORT=5801
MASTER_PORT=5801
WORKER_PORT=5802
MASTER_HTTP_PORT=8080

# JVM 内存配置（GB）
HYBRID_HEAP_SIZE=3
MASTER_HEAP_SIZE=1
WORKER_HEAP_SIZE=3

# 检查点存储
CHECKPOINT_STORAGE_TYPE=LOCAL_FILE   # LOCAL_FILE / HDFS / OSS / S3
CHECKPOINT_NAMESPACE=/tmp/seatunnel/checkpoint/

# systemd 自启动
ENABLE_AUTO_START=true
```

> 建议：先用 Web 向导在浏览器里把配置填好并保存，然后再根据需要手动查看/微调 `config.properties`。

---

## 多节点 / 部署模式要点

- **分离模式（DEPLOY_MODE=separated）**
  - Master 负责控制与协调，Worker 负责执行任务
  - 必须配置 `MASTER_IP` 和 `WORKER_IPS`
  - 端口主要使用 `MASTER_PORT` / `WORKER_PORT` / `MASTER_HTTP_PORT`

- **混合模式（DEPLOY_MODE=hybrid）**
  - 所有节点角色相同，统一写在 `CLUSTER_NODES`
  - 端口主要使用 `HYBRID_PORT`（集群）和 `MASTER_HTTP_PORT`（Web/REST）

> 生产环境部署建议使用分离模式，并结合 HDFS / OSS / S3 作为检查点存储。

---

## 卸载 SeaTunnel

卸载脚本会根据 `config.properties` 中的配置，安全地停止服务并删除安装目录、Java 软链接、systemd 配置等。

```bash
chmod +x uninstall_seatunnel.sh
./uninstall_seatunnel.sh
```

> 卸载前请确认：
> - 不再需要当前集群和相关数据
> - 所有重要配置和日志已自行备份

---

## systemd 服务管理（安装完成后）

如果在配置中开启了 `ENABLE_AUTO_START=true`，安装脚本会自动生成 systemd 服务：

- 混合模式：`seatunnel`
- 分离模式：`seatunnel-master`、`seatunnel-worker`

常用命令示例：

```bash
# Master 节点
sudo systemctl start seatunnel-master
sudo systemctl status seatunnel-master

# Worker 节点
sudo systemctl start seatunnel-worker
sudo systemctl status seatunnel-worker

# 查看日志
sudo journalctl -u seatunnel-master -n 100 --no-pager
sudo journalctl -u seatunnel-worker -n 100 --no-pager
```

---

## 更多信息

- 更复杂的连接器/依赖配置，请直接参考 `config.properties` 中的注释
- SeaTunnel 官方文档：https://seatunnel.apache.org/docs

本 README 仅保留最常用的安装和运行方式，便于快速上手，其余细节以实际脚本和配置文件为准。
