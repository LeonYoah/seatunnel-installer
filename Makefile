.PHONY: all build clean test install help

# 构建变量
BINARY_AGENT=bin/seatunnel-agent
BINARY_CONTROL_PLANE=bin/seatunnel-control-plane

# Go 参数
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# 版本信息
VERSION?=dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

all: clean build

build: build-agent build-control-plane

build-agent:
	@echo "构建 Agent（统一代理，包含安装和运维功能）..."
	@mkdir -p bin
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_AGENT) ./cmd/agent

build-control-plane:
	@echo "构建 Control Plane..."
	@mkdir -p bin
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_CONTROL_PLANE) ./cmd/control-plane

clean:
	@echo "清理构建产物..."
	$(GOCLEAN)
	rm -rf bin/

test:
	@echo "运行测试..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

test-coverage:
	@echo "生成测试覆盖率报告..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

deps:
	@echo "下载依赖..."
	$(GOMOD) download
	$(GOMOD) tidy

fmt:
	@echo "格式化代码..."
	$(GOFMT) ./...

lint:
	@echo "运行代码检查..."
	@which golangci-lint > /dev/null || (echo "请先安装 golangci-lint" && exit 1)
	golangci-lint run ./...

install: build
	@echo "安装二进制文件..."
	cp $(BINARY_AGENT) /usr/local/bin/
	cp $(BINARY_CONTROL_PLANE) /usr/local/bin/
	@echo "安装完成！"
	@echo "  - seatunnel-agent: 统一代理（安装+运维）"
	@echo "  - seatunnel-control-plane: 控制面"

run-agent:
	@echo "运行 Agent..."
	$(GOBUILD) -o $(BINARY_AGENT) ./cmd/agent
	$(BINARY_AGENT) start

run-control-plane:
	@echo "运行 Control Plane..."
	$(GOBUILD) -o $(BINARY_CONTROL_PLANE) ./cmd/control-plane
	$(BINARY_CONTROL_PLANE) server

# Docker 构建
docker-build-agent:
	@echo "构建 Agent Docker 镜像..."
	docker build -t seatunnel/agent:$(VERSION) -f build/docker/Dockerfile.agent .

docker-build-control-plane:
	@echo "构建 Control Plane Docker 镜像..."
	docker build -t seatunnel/control-plane:$(VERSION) -f build/docker/Dockerfile.control-plane .

docker-build: docker-build-agent docker-build-control-plane

help:
	@echo "SeaTunnel 企业级平台 - Makefile 帮助"
	@echo ""
	@echo "可用目标："
	@echo "  make build              - 构建所有组件"
	@echo "  make build-agent        - 构建 Agent（统一代理）"
	@echo "  make build-control-plane - 构建 Control Plane"
	@echo "  make clean              - 清理构建产物"
	@echo "  make test               - 运行测试"
	@echo "  make test-coverage      - 生成测试覆盖率报告"
	@echo "  make deps               - 下载依赖"
	@echo "  make fmt                - 格式化代码"
	@echo "  make lint               - 运行代码检查"
	@echo "  make install            - 安装二进制文件到 /usr/local/bin"
	@echo "  make run-agent          - 运行 Agent"
	@echo "  make run-control-plane  - 运行 Control Plane"
	@echo "  make docker-build       - 构建 Docker 镜像"
	@echo ""
	@echo "架构说明："
	@echo "  Agent 现在是统一组件，包含："
	@echo "    - 安装管理模块：集群部署、卸载、升级、诊断"
	@echo "    - 进程管理模块：SeaTunnel 进程生命周期管理"
	@echo ""
