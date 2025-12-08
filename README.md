# 食堂管理系统

## 项目简介
本系统是一个食堂管理系统的后端服务，实现了用户管理、订单处理、套餐管理等核心功能。

## 系统架构
系统采用三层架构设计：
- **表现层 (Controller)**: 处理HTTP请求，负责路由和参数校验
- **业务逻辑层 (Service)**: 实现核心业务逻辑，处理业务规则和流程
- **数据访问层 (Repository)**: 负责数据库操作，提供数据持久化功能

详细架构说明请参考 [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

## 环境要求
- Go 1.24+
- MySQL 5.7+
- Redis 6.0+

## 安装和运行

### 1. 克隆项目
```bash
git clone https://github.com/yourusername/canteen-backend.git
cd canteen-backend
```

### 2. 安装依赖
```bash
go mod download
```

### 3. 数据库初始化
创建数据库并执行初始化脚本：
```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS canteen DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"

# 执行初始化脚本
mysql -u root -p canteen < scripts/init_database.sql
```

### 4. 配置文件
复制并修改配置文件：
```bash
cp config/config.yaml.example config/config.yaml
```
根据实际环境修改数据库连接、Redis配置等信息。

### 5. 运行项目
#### 方式一：直接运行（开发模式）
```bash
# 使用旧的入口点（已重构，但仍保留）
go run main.go

# 使用新的入口点（推荐）
go run cmd/server/main.go
```

#### 方式二：编译后运行
```bash
# 编译项目
go build -o server.exe cmd/server/main.go

# 运行编译后的程序
./server.exe
```

### 6. 构建发布版本
```bash
# Linux/macOS
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o canteen-server cmd/server/main.go

# Windows
go build -ldflags "-s -w" -o canteen-server.exe cmd/server/main.go

# 交叉编译
# Linux 64位
GOOS=linux GOARCH=amd64 go build -o canteen-server-linux-amd64 cmd/server/main.go

# Windows 64位
GOOS=windows GOARCH=amd64 go build -o canteen-server-windows-amd64.exe cmd/server/main.go
```

## API文档
项目提供以下主要API接口：
- 用户认证和授权
- 套餐管理
- 订单处理
- 数据统计

## 开发指南
详细的开发指南和架构说明请参考 [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

## 许可证
本项目采用许可证授权，使用前请确保许可证有效。



