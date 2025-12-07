# 食堂管理系统架构文档

## 概述

本项目是一个基于Go语言开发的食堂管理系统，采用三层架构模式，实现了清晰的代码分离和组织结构。

## 目录结构

```
canteen-backend/
├── cmd/                    # 应用入口
│   └── server/
│       └── main.go        # 主程序入口
├── internal/               # 内部应用代码
│   ├── app/              # 应用初始化和生命周期管理
│   │   └── app.go
│   ├── controller/        # 表现层 - 处理HTTP请求和响应
│   │   ├── card/         # 卡片相关控制器
│   │   ├── health/       # 健康检查控制器
│   │   ├── tempDirect/   # 临时直接访问控制器
│   │   └── uploadFile/   # 文件上传控制器
│   ├── service/          # 业务逻辑层 - 处理业务规则和流程
│   │   ├── card/         # 卡片相关业务逻辑
│   │   ├── meal/         # 餐食相关业务逻辑
│   │   ├── order/        # 订单相关业务逻辑
│   │   └── user/         # 用户相关业务逻辑
│   ├── repository/       # 数据访问层 - 处理数据库操作
│   │   ├── card/         # 卡片数据访问
│   │   ├── meal/         # 餐食数据访问
│   │   ├── order/        # 订单数据访问
│   │   └── user/         # 用户数据访问
│   ├── model/            # 数据模型和DTO
│   │   ├── user.go       # 用户模型
│   │   ├── order.go      # 订单模型
│   │   ├── meal.go       # 餐食模型
│   │   └── card.go       # 卡片模型
│   ├── infrastructure/   # 基础设施层
│   │   ├── config/       # 配置管理
│   │   ├── database/     # 数据库连接和初始化
│   │   ├── cache/        # 缓存管理
│   │   └── logging/     # 日志处理
│   ├── router/           # 路由配置
│   │   └── router.go
│   └── middleware/       # 中间件
├── pkg/                  # 可复用的库代码
│   └── utils/            # 工具函数
│       ├── cron_tasks.go  # 定时任务
│       ├── embed_key.go   # 嵌入密钥
│       └── public.pem    # 公钥文件
├── api/                  # API定义和文档
├── config/               # 配置文件
├── logs/                 # 日志文件
├── docs/                 # 项目文档
│   └── ARCHITECTURE.md
├── model/                # 原有模型（已迁移至internal/model）
├── controller/           # 原有控制器（已迁移至internal/controller）
├── router/               # 原有路由（已迁移至internal/router）
├── utils/                # 原有工具（已迁移至pkg/utils）
├── main.go               # 原有主程序入口
└── server.exe            # 编译后的可执行文件
```

## 架构层次

### 1. 表现层 (Controller Layer)

位置：`internal/controller/`

职责：
- 处理HTTP请求和响应
- 参数验证和格式化
- 调用业务逻辑层处理业务请求
- 返回适当的HTTP状态码和响应

主要组件：
- `card/card_controller.go`: 处理卡片相关请求
- `health/health_controller.go`: 处理健康检查请求
- `tempDirect/tempDirect_controller.go`: 处理临时直接访问请求
- `uploadFile/upload_controller.go`: 处理文件上传请求

### 2. 业务逻辑层 (Service Layer)

位置：`internal/service/`

职责：
- 实现业务规则和流程
- 协调多个数据访问操作
- 事务管理
- 业务逻辑验证

主要组件：
- `card/card_service.go`: 卡片相关业务逻辑
- `meal/meal_service.go`: 餐食相关业务逻辑
- `order/order_service.go`: 订单相关业务逻辑
- `user/user_service.go`: 用户相关业务逻辑

### 3. 数据访问层 (Repository Layer)

位置：`internal/repository/`

职责：
- 数据库操作抽象
- SQL查询执行
- 数据模型映射
- 缓存集成

主要组件：
- `card/card_repository.go`: 卡片数据访问
- `meal/meal_repository.go`: 餐食数据访问
- `order/order_repository.go`: 订单数据访问
- `user/user_repository.go`: 用户数据访问

## 基础设施层

### 配置管理

位置：`internal/infrastructure/config/`

职责：
- 应用配置加载和管理
- 配置验证
- 环境特定配置

### 数据库管理

位置：`internal/infrastructure/database/`

职责：
- 数据库连接初始化
- 连接池管理
- 数据库生命周期管理

### 缓存管理

位置：`internal/infrastructure/cache/`

职责：
- Redis连接管理
- 缓存操作抽象
- 缓存策略实现

### 日志处理

位置：`internal/infrastructure/logging/`

职责：
- 日志格式化和输出
- 日志轮转和归档
- 不同级别日志分离

## 依赖注入

依赖关系遵循以下原则：
- Controller依赖于Service接口
- Service依赖于Repository接口
- Repository依赖于数据库抽象

依赖注入在控制器初始化时进行，通过构造函数或设置方法注入依赖。

## 应用生命周期

应用启动流程：
1. 验证许可证
2. 初始化配置
3. 创建应用实例
4. 初始化数据库和缓存
5. 注册控制器
6. 启动后台任务
7. 启动HTTP服务器
8. 处理关闭信号
9. 优雅关闭

## 代码组织原则

1. **单一职责原则**：每个类、模块只负责一项功能
2. **依赖倒置原则**：高层模块不依赖低层模块，两者都依赖抽象
3. **开闭原则**：对扩展开放，对修改关闭
4. **接口隔离原则**：使用多个专门的接口，而不是单一的通用接口
5. **领域驱动设计**：代码组织反映业务领域

## 重构收益

1. **代码组织更清晰**：按功能模块和职责层次组织代码
2. **依赖关系明确**：通过接口定义依赖，降低耦合
3. **测试更容易**：可以轻松模拟依赖进行单元测试
4. **维护成本降低**：修改某一层不会影响其他层
5. **扩展性更好**：新功能可以通过添加新的服务或实现来扩展
6. **团队协作更高效**：不同层次可以由不同团队成员并行开发

## 未来改进方向

1. **引入依赖注入框架**：如Wire或Dig，进一步简化依赖管理
2. **完善单元测试**：为各层添加全面的单元测试
3. **添加中间件机制**：认证、授权、限流等
4. **完善错误处理**：统一错误处理机制和错误码
5. **添加监控和指标**：系统运行状态监控和性能指标
6. **API文档生成**：使用Swagger等工具自动生成API文档