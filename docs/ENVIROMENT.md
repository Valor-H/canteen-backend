# Windows Golang项目环境配置指南

## 项目概述

本文档详细说明在Windows系统上配置和运行Golang项目的完整流程，包括环境变量设置、依赖安装和项目启动。

项目路径：`D:/Codes/canteen-backend`

## 1. 环境要求

- Windows 10/11
- Go 1.19+ (已安装到 D:/Tools/Go)
- MySQL 5.7+ 或 8.0+
- CMD命令行环境

## 2. 环境变量配置

### 2.1 检查Go安装

```cmd
go version
```

如果安装正确，应显示Go版本信息。

### 2.2 设置系统环境变量

打开CMD命令行，执行以下命令：

```cmd
set GOROOT=D:\Tools\Go
set GOPATH=D:\GoProjects
set PATH=%GOROOT%\bin;%GOPATH%\bin;%PATH%
```

验证环境变量：

```cmd
go env GOROOT
go env GOPATH
go env
```

## 3. 项目配置

### 3.1 进入项目目录

```cmd
cd /d D:/Codes/canteen-backend
```

### 3.2 配置Go代理（解决依赖下载慢的问题）

```cmd
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn
```

### 3.3 项目依赖安装

```cmd
go mod tidy -v
```

如果`go mod tidy`卡住，可以尝试以下解决方案：

#### 解决方案1：使用国内代理

```cmd
go env -w GOPROXY=https://goproxy.io,direct
go mod tidy -v
```

#### 解决方案2：清除模块缓存

```cmd
go clean -modcache
go mod download
go mod tidy -v
```

#### 解决方案3：使用阿里云代理

```cmd
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
go mod tidy -v
```

## 4. 数据库配置

### 4.1 MySQL安装与配置

1. 安装MySQL，确保服务运行在端口3306
2. 创建数据库：

```sql
CREATE DATABASE canteen CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

3. 创建用户并授权：

```sql
CREATE USER 'canteen'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON canteen.* TO 'canteen'@'localhost';
FLUSH PRIVILEGES;
```

### 4.2 检查项目配置

查看`config/config.yaml`文件中的数据库配置，确保与实际设置匹配。

## 5. 项目启动

### 5.1 运行项目

```cmd
go run main.go
```

### 5.2 编译项目（可选）

```cmd
go build -o canteen.exe main.go
```

### 5.3 运行编译后的可执行文件

```cmd
canteen.exe
```

## 6. 项目结构

```cmd
canteen-backend

├── config

│   └── config.yaml

├── controllers

│   ├── auth.go

│   ├── canteen.go

│   ├── category.go

│   ├── dish.go

│   ├── order.go

│   └── user.go

├── database

│   └── database.go

├── main.go

├── middleware

│   └── auth.go

├── models

│   ├── category.go

│   ├── dish.go

│   ├── order.go

│   └── user.go

└── routes

    └── router.go
```


## 7. 项目依赖说明

项目主要依赖以下Go模块：

- gin-gonic/gin: Web框架
- go-sql-driver/mysql: MySQL驱动
- spf13/viper: 配置管理
- gorm.io/gorm: ORM框架
- golang-jwt/jwt: JWT认证


## 8. 参考资源

- Go官方文档：https://golang.org/doc/
- Gin框架文档：https://gin-gonic.com/docs/
- MySQL文档：https://dev.mysql.com/doc/