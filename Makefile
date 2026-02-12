# Catwalk 项目 Makefile
# 功能完善的构建脚本，支持项目的主要开发流程

# 定义变量
GO := go
GOPATH := $(shell $(GO) env GOPATH)
GOBIN := $(GOPATH)/bin

# 项目名称
PROJECT := catwalk-cn

# 源代码目录
SRC_DIR := .

# 命令目录
CMD_DIR := cmd

# 内部包目录
INTERNAL_DIR := internal

# 公共包目录
PKG_DIR := pkg

# 构建输出目录
BUILD_DIR := build

# 定义目标
.PHONY: all build compile test clean install init tidy format lint run help

# 默认目标
all: build test

# 构建项目
# 运行 `make build` 来构建项目
build:
	@echo "=== 构建项目 ==="
	@$(GO) build -v ./...

# 编译可执行文件
# 运行 `make compile` 来编译可执行文件
compile:
	@echo "=== 编译可执行文件 ==="
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(PROJECT) -v .

# 运行测试
# 运行 `make test` 来执行所有测试
# 运行 `make test TEST=./path/to/test` 来执行特定测试
TEST ?= ./...
test:
	@echo "=== 运行测试 ==="
	@$(GO) test -v $(TEST)

# 清理构建产物
# 运行 `make clean` 来清理构建产物和临时文件
clean:
	@echo "=== 清理构建产物 ==="
	@$(GO) clean ./...
	@rm -rf $(BUILD_DIR)
	@find . -name "*.test" -type f -delete

# 安装项目
# 运行 `make install` 来安装项目到 GOPATH
install:
	@echo "=== 安装项目 ==="
	@$(GO) install -v .

# 初始化模块
# 运行 `make init` 来初始化 Go 模块
init:
	@echo "=== 初始化模块 ==="
	@$(GO) mod init github.com/purpose168/catwalk-cn
	@$(MAKE) tidy

# 整理依赖
# 运行 `make tidy` 来整理和更新依赖

tidy:
	@echo "=== 整理依赖 ==="
	@$(GO) mod tidy

# 格式化代码
# 运行 `make format` 来格式化代码
format:
	@echo "=== 格式化代码 ==="
	@$(GO) fmt ./...

# 运行 lint 检查
# 运行 `make lint` 来执行代码风格检查
lint:
	@echo "=== 运行 lint 检查 ==="
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "安装 golangci-lint..."; \
		$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@golangci-lint run

# 运行项目
# 运行 `make run` 来运行项目
run:
	@echo "=== 运行项目 ==="
	@$(GO) run .

# 显示帮助信息
# 运行 `make help` 来显示此帮助信息
help:
	@echo "=== Catwalk 项目 Makefile 帮助 ==="
	@echo "可用目标:"
	@echo "  make all          - 构建并测试项目（默认目标）"
	@echo "  make build        - 构建项目"
	@echo "  make compile      - 编译可执行文件"
	@echo "  make test         - 运行测试（可指定 TEST=./path/to/test）"
	@echo "  make clean        - 清理构建产物"
	@echo "  make install      - 安装项目到 GOPATH"
	@echo "  make init         - 初始化 Go 模块"
	@echo "  make tidy         - 整理和更新依赖"
	@echo "  make format       - 格式化代码"
	@echo "  make lint         - 运行 lint 检查"
	@echo "  make run          - 运行项目"
	@echo "  make help         - 显示此帮助信息"
	@echo ""
	@echo "使用示例:"
	@echo "  make build        # 构建项目"
	@echo "  make test         # 运行所有测试"
	@echo "  make test TEST=./pkg/path/to/test # 运行特定测试"
	@echo "  make clean        # 清理构建产物"
	@echo "  make run          # 运行项目"
