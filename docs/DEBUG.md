# Golang 断点调试完整指南

## 概述

Golang提供了强大的调试支持，通过Delve调试器可以实现断点调试、变量检查、调用堆栈查看等功能。本文将介绍在主流开发环境中配置Go调试的方法。

## 1. 调试工具基础 - Delve

### 安装Delve

```bash
# 安装最新版本
go install github.com/go-delve/delve/cmd/dlv@latest

# 设置环境变量
set PATH=%PATH%;%USERPROFILE%\go\bin

# 验证安装
dlv version
```

### Delve基本命令

```bash
# 调试当前目录下的main包
dlv debug

# 调试特定包
dlv debug ./package

# 调试测试文件
dlv test

# 附加到正在运行的进程
dlv attach <pid>

# 调试二进制文件
dlv exec ./binary
```

## 2. VS Code集成方案

### 必需插件

1. **Go扩展** (官方推荐)
   - 插件名称: `Go`
   - 发布者: `golang.go`
   - 功能: Go语言支持、调试、测试等

2. **Delve调试器** (随Go扩展自动安装)
   - Go扩展会自动提示安装`dlv`调试器

### VS Code配置步骤

1. **安装Go扩展**
   - 打开VS Code扩展面板
   - 搜索"Go"并安装官方扩展

2. **安装工具**
   - 安装扩展后，按`Ctrl+Shift+P`
   - 输入`Go: Install/Update Tools`
   - 选择安装所有工具（包括`dlv`）

3. **创建launch.json配置**

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}"
        },
        {
            "name": "Launch File",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${file}"
        },
        {
            "name": "Attach to Process",
            "type": "go",
            "request": "attach",
            "mode": "local",
            "processId": 0
        },
        {
            "name": "Launch Test",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${workspaceFolder}"
        }
    ]
}
```

4. **VS Code调试功能**
   - F5: 开始调试
   - F9: 切换断点
   - F10: 单步执行
   - F11: 进入函数
   - Shift+F11: 跳出函数

## 3. GoLand集成方案 (IntelliJ)

### 配置步骤

1. **安装Go插件**
   - 打开File > Settings > Plugins
   - 搜索"Go"并安装JetBrains官方插件

2. **配置Go SDK**
   - File > Settings > Go > GOROOT
   - 选择Go安装目录

3. **配置GOPATH**
   - File > Settings > Go > GOPATH
   - 设置工作区路径

4. **创建运行/调试配置**
   - Run > Edit Configurations
   - 点击"+"选择"Go Build"
   - 配置以下参数:
     - Name: 调试配置名称
     - Run kind: 选择"Directory"或"File"
     - Directory/File: 选择要调试的目录或文件
     - Output directory: 编译输出目录

5. **调试功能**
   - 点击代码行左侧设置断点
   - 使用Debug按钮启动调试
   - 支持变量查看、表达式计算等

## 4. Vim/Neovim集成方案

### 必需插件

1. **vim-go** (Vim)
   ```vim
   Plug 'fatih/vim-go'
   ```

2. **nvim-dap** (Neovim)
   ```lua
   use { "mfussenegger/nvim-dap" }
   use { "leoluz/nvim-dap-go", requires = {"mfussenegger/nvim-dap"} }
   ```

### 配置示例 (Neovim)

```lua
require('dap-go').setup()
```

## 5. 命令行调试详解

### Delve交互命令

```bash
# 启动调试会话
dlv debug ./main.go

# Delve交互命令
(dlv) break main.go:20          # 在第20行设置断点
(dlv) breakpoints               # 列出所有断点
(dlv) clear main.go:20         # 清除断点
(dlv) continue                 # 继续执行
(dlv) next                     # 单步执行
(dlv) step                     # 进入函数
(dlv) stepout                  # 跳出函数
(dlv) print variableName       # 打印变量值
(dlv) locals                   # 显示所有局部变量
(dlv) args                     # 显示函数参数
(dlv) stack                    # 显示调用栈
(dlv) goroutines               # 显示所有goroutine
(dlv) goroutine <id>           # 切换到指定goroutine
(dlv) exit                     # 退出调试器
```

### 高级调试技巧

```bash
# 条件断点
(dlv) break main.go:20 if x > 10

# 临时断点
(dlv) trace main.go:20

# 监视点
(dlv) watch variableName

# 调用函数
(dlv) call someFunction(args)
```