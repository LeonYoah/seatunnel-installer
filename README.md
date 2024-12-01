# SeaTunnel 一键安装指南

[![Apache License 2.0](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

SeaTunnel 是一个高性能、分布式的数据集成平台，支持实时和批量数据同步。本指南将帮助您快速完成 SeaTunnel 的安装部署。

## 🚀 快速开始

只需一个命令即可完成安装：

```bash
./install_seatunnel.sh
```

安装脚本会自动完成所有配置和部署步骤，包括环境检查、依赖安装、集群配置等。

## ✅ 环境要求

- Java 8 或更高版本
- 足够的磁盘空间（建议 > 10GB）
- 支持的操作系统：CentOS 7+/Ubuntu 18.04+
- 安装用户需要 sudo 权限
- 各节点间需要免密 SSH 访问

## 📋 安装前配置

1. 编辑 config.properties 文件，设置基本参数：

```properties
# 必选配置
SEATUNNEL_VERSION=2.3.7      # SeaTunnel 版本
INSTALL_MODE=offline         # 安装模式：online/offline
BASE_DIR=/data/seatunnel    # 安装目录

# 可选配置
DEPLOY_MODE=separated       # 部署模式：separated/hybrid
```

2. 确保安装包位于正确位置（离线安装模式）：
   - 将 `apache-seatunnel-${VERSION}.tar.gz` 放在脚本同目录下

## 🔧 部署模式

### 混合模式 (Hybrid)
- 所有节点运行相同的组件
- 适合小规模部署
- 配置简单，维护方便

### 分离模式 (Separated)
- Master 和 Worker 分开部署
- 适合大规模生产环境
- 更好的资源隔离和扩展性

## 🎯 一键部署步骤

1. 下载安装包：
```bash
wget https://github.com/LeonYoah/seatunnel-installer/archive/main.zip
unzip main.zip
```

2. 进入安装目录：
```bash
cd seatunnel-installer
```

3. 执行安装：
```bash
./install_seatunnel.sh
```

安装过程会自动：
- ✅ 检查系统环境
- ✅ 部署 SeaTunnel 组件
- ✅ 启动服务
- ✅ 验证安装结果

## 🔍 验证安装

安装完成后，可以通过以下命令验证：

```bash
# 检查服务状态
sudo systemctl status seatunnel

# 查看集群状态
${BASE_DIR}/bin/seatunnel status
```

## 📝 常见问题

1. 安装失败如何处理？
   - 检查 logs 目录下的安装日志
   - 确保满足所有环境要求
   - 重新执行安装脚本

2. 服务无法启动？
   - 检查端口占用情况
   - 验证配置文件正确性
   - 查看系统日志

3. 性能调优建议？
   - 适当调整 JVM 参数
   - 优化系统参数
   - 参考性能优化指南

## 🆘 获取帮助

- 查看详细文档：[官方文档](https://seatunnel.apache.org/)
- 提交 Issue：[GitHub Issues](https://github.com/apache/seatunnel/issues)
- 社区支持：[Slack Channel](https://slack.seatunnel.apache.org/)

## 📦 下一步

- [配置数据源](https://seatunnel.apache.org/docs/connector-v2/source)
- [配置数据目标](https://seatunnel.apache.org/docs/connector-v2/sink)
- [开发自定义连接器](https://seatunnel.apache.org/docs/development/connector-v2)
