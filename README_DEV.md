# SeaTunnel 企业级平台 - 开发指南

## 架构概述

本项目采用统一的 Agent 架构，将原有的 Installer 和 Agent 合并为一个统一的组件。

### 核心组件

1. **Agent（统一代理）**：部署在每个节点上，包含两大功能模块
   - 安装管理模块：负责集群部署、卸载、升级、诊断
   - 进程管理模块：负责 SeaTunnel 进程生命周期管理、监控、日志收集

2. **Control Plane（控制面）**：提供 Web UI 和 REST API，统一管理集群

详细架构说明请参考：[docs/ARCHITECTURE_UPDATE.md](docs/ARCHITECTURE_UPDATE.md)

## 项目结构

```
.
├── cmd/                          # 主应用程序入口
│   ├── agent/                   # Agent 统一入口（包含安装和运维功能）
│   └── control-plane/           # Control Plane 入口
├── internal/                     # 私有应用代码
│   ├── agent/                   # Agent 实现
│   │   ├── installer/          # 安装管理模块
│   │   ├── process/            # 进程管理模块
│   │   ├── monitor/            # 监控模块
│   │   └── common/             # 共享代码
│   ├── controlplane/            # Control Plane 实现
│   ├── api/                     # API 处理器
│   ├── service/                 # 业务逻辑
│   ├── repository/              # 数据访问层
│   └── models/                  # 数据模型
├── pkg/                         # 公共库
│   ├── logger/                  # 日志工具
│   ├── config/                  # 配置管理
│   ├── utils/                   # 通用工具
│   └── errors/                  # 错误处理
├── web/                         # 前端代码（Vue3）
├── scripts/                     # 构建和部署脚本
├── docs/                        # 文档
│   └── ARCHITECTURE_UPDATE.md  # 架构更新说明
└── tests/                       # 测试文件
```

## 环境要求

- Go 1.21+
- Node.js 18+（用于前端开发）
- Make
- Docker（可选，用于容器化部署）

## 快速开始

### 1. 安装依赖

```bash
make deps
```

### 2. 构建所有二进制文件

```bash
make build
```

这将在 `bin/` 目录下创建两个二进制文件：
- `seatunnel-agent` - Agent 统一代理（包含安装和运维功能）
- `seatunnel-control-plane` - Control Plane 服务器

### 3. 运行组件

```bash
# 运行 Agent
make run-agent

# 运行 Control Plane
make run-control-plane
```

## 开发工作流

### 构建

```bash
# 构建所有组件
make build

# 构建特定组件
make build-agent
make build-control-plane

# 查看所有可用命令
make help
```

### 测试

```bash
# 运行所有测试
make test

# 生成测试覆盖率报告
make test-coverage

# 运行特定包的测试
go test -v ./pkg/logger/...
go test -v ./pkg/utils/...
go test -v ./pkg/errors/...
```

### 代码质量

```bash
# 格式化代码
make fmt

# 运行代码检查（需要先安装 golangci-lint）
make lint
```

### 清理

```bash
make clean
```

## Agent 命令行接口

Agent 现在提供统一的命令行接口：

### 进程管理命令

```bash
# 启动 Agent 守护进程
seatunnel-agent start [--config=/path/to/config.yaml]

# 停止 Agent 守护进程
seatunnel-agent stop

# 查看 Agent 状态
seatunnel-agent status
```

### 安装管理命令

```bash
# 安装 SeaTunnel
seatunnel-agent install [--config=/path/to/install-config.yaml]

# 卸载 SeaTunnel
seatunnel-agent uninstall

# 升级 SeaTunnel
seatunnel-agent upgrade --version=2.3.13

# 环境预检查
seatunnel-agent precheck

# 收集诊断信息
seatunnel-agent diagnose [--output=/path/to/output.tar.gz]
```

### 通用命令

```bash
# 查看版本信息
seatunnel-agent version
```

## 配置

每个组件都可以通过以下方式配置：
1. 配置文件（YAML）
2. 环境变量
3. 命令行参数

### Agent 配置示例 (`agent-config.yaml`)

```yaml
agent:
  # Agent 基本配置
  id: "agent-node1"
  name: "SeaTunnel Agent Node 1"
  
  # Control Plane 连接配置
  control_plane:
    address: "control-plane.example.com:50051"
    tls:
      enabled: true
      cert_file: "/etc/agent/certs/client.crt"
      key_file: "/etc/agent/certs/client.key"
  
  # 心跳配置
  heartbeat:
    interval: 10s
    timeout: 30s
  
  # 日志配置
  log:
    level: info
    output: /var/log/seatunnel-agent/agent.log
    max_size: 100
    max_backups: 3
    max_age: 28
  
  # 安装管理配置
  installer:
    work_dir: /tmp/seatunnel-installer
    package_cache: /var/cache/seatunnel
  
  # 进程管理配置
  process:
    seatunnel_home: /opt/seatunnel
    check_interval: 5s
    restart_on_failure: true
    max_restart_attempts: 3
```

### Control Plane 配置示例 (`control-plane-config.yaml`)

```yaml
server:
  port: 8080
  host: 0.0.0.0

database:
  type: sqlite
  database: seatunnel.db

logger:
  level: info
  output_paths:
    - stdout
    - /var/log/seatunnel/control-plane.log
```

## Docker 构建

```bash
# 构建 Agent 镜像
make docker-build-agent

# 构建 Control Plane 镜像
make docker-build-control-plane

# 构建所有镜像
make docker-build
```

## 开发任务进度

### 第一阶段：基础框架搭建 ✅

- [x] 1.1 实现配置管理模块
- [x] 1.2 实现日志框架
- [x] 1.3 实现工具函数库
- [x] 1.4 实现错误处理和恢复机制

### 第二阶段：Agent 统一组件开发 🚧

- [ ] 4. 实现 Agent 基础框架
  - [ ] 4.1 实现 Agent 安装管理模块
  - [ ] 4.2 实现预检查功能
  - [ ] 4.3-4.9 实现各项检查功能
- [ ] 5-7. 实现安装包处理、配置生成、插件管理
- [ ] 8. 实现 Control Plane 节点分发功能
- [ ] 9. 实现 Agent 进程管理模块
- [ ] 10-11. 实现集群启动和卸载功能

详细任务列表请参考：[.kiro/specs/seatunnel-enterprise-platform/tasks.md](.kiro/specs/seatunnel-enterprise-platform/tasks.md)

## 代码规范

### Go 代码规范

1. **注释**：所有代码注释必须使用中文
2. **命名**：使用驼峰命名法，导出的标识符首字母大写
3. **错误处理**：使用 `pkg/errors` 包进行统一的错误处理
4. **日志**：使用 `pkg/logger` 包记录日志
5. **格式化**：使用 `gofmt` 或 `make fmt` 格式化代码

### 示例代码

```go
package example

import (
    "context"
    
    "github.com/seatunnel/enterprise-platform/pkg/errors"
    "github.com/seatunnel/enterprise-platform/pkg/logger"
)

// DoSomething 执行某个操作
// 这是一个示例函数，展示如何使用错误处理和日志
func DoSomething(ctx context.Context, input string) error {
    // 参数验证
    if input == "" {
        return errors.New(errors.ErrCodeInvalidParam, "输入参数不能为空")
    }
    
    // 记录日志
    logger.Info("开始执行操作", 
        zap.String("input", input))
    
    // 执行操作
    if err := performOperation(input); err != nil {
        return errors.Wrap(err, errors.ErrCodeInternalError, "操作执行失败")
    }
    
    logger.Info("操作执行成功")
    return nil
}
```

## 测试

### 单元测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./pkg/logger/
go test ./pkg/utils/
go test ./pkg/errors/

# 查看测试覆盖率
go test -cover ./...
```

### 集成测试

```bash
# 运行集成测试（需要先启动依赖服务）
go test -tags=integration ./tests/integration/...
```

## 调试

### 使用 Delve 调试

```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试 Agent
dlv debug ./cmd/agent -- start --config=agent-config.yaml

# 调试 Control Plane
dlv debug ./cmd/control-plane -- server --config=control-plane-config.yaml
```

### 查看日志

```bash
# Agent 日志
tail -f /var/log/seatunnel-agent/agent.log

# Control Plane 日志
tail -f /var/log/seatunnel/control-plane.log

# 使用 journalctl（如果使用 systemd）
sudo journalctl -u seatunnel-agent -f
```

## 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m '添加某个特性'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 常见问题

### Q: 为什么将 Installer 和 Agent 合并？

A: 为了简化部署和运维，减少组件数量，降低资源占用，提高可维护性。详见 [docs/ARCHITECTURE_UPDATE.md](docs/ARCHITECTURE_UPDATE.md)

### Q: 如何从旧架构迁移？

A: 参考 [docs/ARCHITECTURE_UPDATE.md](docs/ARCHITECTURE_UPDATE.md) 中的迁移指南。

### Q: 如何添加新的安装步骤？

A: 在 `internal/agent/installer/steps/` 中添加新的步骤实现。

### Q: 如何添加新的运维操作？

A: 在 `internal/agent/process/actions/` 中添加新的操作实现。

## 相关文档

- [架构更新说明](docs/ARCHITECTURE_UPDATE.md)
- [需求文档](.kiro/specs/seatunnel-enterprise-platform/requirements.md)
- [设计文档](.kiro/specs/seatunnel-enterprise-platform/design.md)
- [任务列表](.kiro/specs/seatunnel-enterprise-platform/tasks.md)
- [路线图](docs/ROADMAP.md)

## 许可证

Apache License 2.0

## 本地开发
windows需要下载：https://github.com/jmeubank/tdm-gcc/releases/download/v10.3.0-tdm64-2/tdm64-gcc-10.3.0-2.exe，sqlite3需要cgo