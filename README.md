# mini-mall-gin-gorm

一个面向初学者的 Go + Gin + Gorm 学习项目示例。

目标是：

- 项目分层接近正规商业项目。
- 流程尽量简化，只保留核心模板代码。
- 覆盖完整最小业务闭环：用户登录、商品 CRUD、购买扣库存。

## 1. 技术栈

- Go
- Gin
- Gorm
- MySQL 8.0
- JWT（登录态）

## 2. 项目结构

```text
mini-mall-gin-gorm
├── cmd
│   └── server
│       └── main.go                 # 程序入口，依赖装配、路由启动
├── configs
│   └── .env.example                # 环境变量模板
├── internal
│   ├── config
│   │   └── config.go               # 配置加载
│   ├── handler
│   │   ├── auth_handler.go         # 注册/登录接口
│   │   ├── product_handler.go      # 商品接口
│   │   ├── purchase_handler.go     # 购买接口
│   │   └── user_handler.go         # 用户信息接口
│   ├── middleware
│   │   └── jwt_auth.go             # JWT 鉴权中间件
│   ├── model
│   │   ├── product.go              # 商品模型
│   │   ├── purchase.go             # 购买模型
│   │   └── user.go                 # 用户模型
│   ├── repository
│   │   ├── product_repository.go   # 商品数据访问
│   │   ├── purchase_repository.go  # 购买数据访问
│   │   └── user_repository.go      # 用户数据访问
│   ├── router
│   │   └── router.go               # 路由注册
│   └── service
│       ├── auth_service.go         # 认证业务
│       ├── product_service.go      # 商品业务
│       ├── purchase_service.go     # 购买业务
│       └── user_service.go         # 用户业务
├── pkg
│   ├── jwtutil
│   │   └── jwt.go                  # JWT 工具
│   ├── password
│   │   └── password.go             # 密码哈希工具
│   └── response
│       └── response.go             # 统一响应格式
└── scripts
    └── sql
        └── init.sql                # MySQL 初始化脚本
```

## 3. 快速启动

### 3.1 创建数据库表

先执行：

- `scripts/sql/init.sql`

你可以使用 MySQL 客户端命令：

```bash
mysql -u root -p < scripts/sql/init.sql
```

### 3.2 配置环境变量

复制模板并按你本地环境修改：

```bash
cp configs/.env.example .env
```

然后导出环境变量（示例）：

```bash
export APP_NAME=mini-mall-gin-gorm
export APP_PORT=8080
export APP_API_PREFIX=/api
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=root
export DB_NAME=mini_mall
export JWT_SECRET=ougeishin
export JWT_EXPIRE_HOURS=24
```

### 3.3 运行服务

```bash
go mod tidy
go run ./cmd/server
```

启动后访问：

- `GET /api/ping`

API 前缀可配置：

- 设置 `APP_API_PREFIX=/api`，接口形如 `/api/products`
- 设置 `APP_API_PREFIX=/api/v1`，接口形如 `/api/v1/products`
- 设置 `APP_API_PREFIX=`（空字符串）时，接口无前缀，形如 `/products`

## 4. 核心接口

### 4.1 认证

- `POST /api/auth/register` 注册
- `POST /api/auth/login` 登录

登录后拿到 token，后续请求头带：

```text
Authorization: Bearer <token>
```

### 4.2 用户

- `GET /api/users/me` 获取当前用户信息
- `PUT /api/users/me` 更新当前用户信息

### 4.3 商品

- `GET /api/products` 商品列表（支持 page、size）
- `GET /api/products/:id` 商品详情
- `POST /api/products` 创建商品（需登录）
- `PUT /api/products/:id` 更新商品（需登录）
- `DELETE /api/products/:id` 删除商品（需登录）

### 4.4 购买

- `POST /api/purchases` 购买商品（需登录）
- `GET /api/purchases/me` 查看我的购买记录（需登录）

## 5. 购买流程说明（重点）

1. 用户发起购买请求，包含商品 ID 和数量。
2. 服务层开启数据库事务。
3. 使用 `SELECT ... FOR UPDATE` 锁定商品行，防止并发超卖。
4. 校验库存是否足够。
5. 扣减库存并写入购买记录。
6. 事务提交，购买成功。

这就是最小可用的“下单扣库存”核心链路。

## 6. 学习建议

- 第一步先把代码跑通，再看分层。
- 第二步重点看 service 层如何组织业务规则。
- 第三步尝试自己加一个“取消订单并回滚库存”的接口，练习事务。
