package generate

import (
	"fmt"
	"strings"
	"text/template"
)

const (
	overviewDesc = `## 概述

1. **服务名称与定位**  
   - 是什么类型的服务？（如：用户管理微服务、订单处理API、数据同步任务等）
   - 解决什么问题？（如：提供用户注册/登录能力、处理电商订单生命周期等）

2. **核心功能**  
   - 用1-3句话概括主要功能。

3. **服务边界**（可选）  
   - 如果是微服务，说明与其他服务的关系（如：依赖哪些服务/被哪些服务调用）。
`

	deploymentDesc = `## 部署

- [裸机部署](https://go-sponge.com/zh/deployment/binary.html)
- [Docker 部署](https://go-sponge.com/zh/deployment/docker.html)
- [K8S 部署](https://go-sponge.com/zh/deployment/kubernetes.html)

`
)

var (
	//nolint
	httpServerReadmeTmplRaw = `## 技术栈

- 编程语言: go
- Web框架: gin
- 配置管理: viper
- 日志: zap
- ORM: {{.ORMType}}
- 数据库: {{.DatabaseType}}
- 缓存: go-redis
- 监控: prometheus+grafana
- 链路追踪: opentracing+jaeger
- 其他: ...

## 目录结构

<BQ><BQ><BQ>text
.
├─ cmd                          # 应用程序入口目录
│   └─ {{.ServerName}}                     # 服务名称
│       ├─ initial              # 初始化逻辑(如配置加载、服务初始化等)
│       └─ main.go              # 主程序入口文件
├─ configs                      # 配置文件目录(yaml 格式配置模板)
├─ deployments                  # 部署相关脚本(二进制、Docker、K8S 部署)
├─ docs                         # 项目文档(API 文档、设计文档等)
├─ internal                     # 内部实现代码(对外不可见)
│   ├─ cache                    # 缓存相关实现(Redis 或本地内存缓存封装)
│   ├─ config                   # 配置解析和结构体定义
│   ├─ dao                      # 数据访问层(Database Access Object)
│   ├─ ecode                    # 错误码定义
│   ├─ handler                  # 业务逻辑处理层(类似 Controller)
│   ├─ model                    # 数据模型/实体定义
│   ├─ routers                  # 路由定义和中间件
│   ├─ server                   # 服务启动
│   └─ types                    # 请求/响应结构体定义
├─ scripts                      # 实用脚本(如代码生成、构建、运行、部署等)
├─ go.mod                       # Go 模块定义文件(声明依赖)
├─ go.sum                       # Go 模块校验文件(自动生成)
├─ Makefile                     # 项目构建自动化脚本
└─ README.md                    # 项目说明文档
<BQ><BQ><BQ>

代码采用分层架构，完整调用链路如下：

<BQ>cmd/{{.ServerName}}/main.go<BQ> → <BQ>internal/server/http.go<BQ> → <BQ>internal/routers/router.go<BQ> → <BQ>internal/handler<BQ> → <BQ>internal/dao<BQ> → <BQ>internal/model<BQ>

其中 handler 层主要负责 API 处理，若需处理更复杂业务逻辑，建议在 handler 和 dao 之间额外添加业务逻辑层（如 <BQ>service<BQ>、<BQ>logic<BQ> 或 <BQ>biz<BQ> 等，自己定义）。

## 快速开始

### 1. 生成 openapi 文档

<BQ><BQ><BQ>bash
make docs
<BQ><BQ><BQ>

注：仅当新增或修改 API 时需要执行该命令，API 未变更时无需重复执行。

### 2. 编译和运行

<BQ><BQ><BQ>bash
make run
<BQ><BQ><BQ>

### 3. 测试 API

在浏览器访问 [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)，测试 HTTP API。

## 开发指南

点击查看详细的 [**开发指南**](https://go-sponge.com/zh/guide/web/based-on-sql.html)。

`

	//nolint
	httpPbServerReadmeTmplRaw = `## 技术栈

- 编程语言: go
- Web框架: gin
- 序列化协议: protobuf
- 配置管理: viper
- 日志: zap
- 监控: prometheus+grafana
- 链路追踪: opentracing+jaeger
- 其他: ...

## 目录结构

<BQ><BQ><BQ>text
├─ api                          # API 协议定义目录(proto/OpenAPI 等)
│   └─ {{.ServerName}}                     # 服务名称
│       └─ v1                   # API 版本 v1(存放 proto 文件和生成的 pb.go 等)
├─ cmd                          # 应用程序入口目录
│   └─ {{.ServerName}}                     # 服务名称
│       ├─ initial              # 初始化逻辑(如配置加载、服务初始化等)
│       └─ main.go              # 主程序入口文件
├─ configs                      # 配置文件目录(yaml 格式配置模板)
├─ deployments                  # 部署相关脚本(二进制、Docker、K8S 部署)
├─ docs                         # 项目文档(API 文档、设计文档等)
├─ internal                     # 内部实现代码(对外不可见)
│   ├─ config                   # 配置解析和结构体定义
│   ├─ ecode                    # 错误码定义
│   ├─ handler                  # 业务逻辑处理层(类似 Controller)
│   ├─ routers                  # 路由定义和中间件
│   └─ server                   # 服务启动
├─ scripts                      # 实用脚本(如代码生成、构建、运行、部署等)
├─ third_party                  # 第三方 Protobuf 依赖/工具
├─ go.mod                       # Go 模块定义文件(声明依赖)
├─ go.sum                       # Go 模块校验文件(自动生成)
├─ Makefile                     # 项目构建自动化脚本
└─ README.md                    # 项目说明文档
<BQ><BQ><BQ>

代码采用分层架构，完整调用链路如下：

<BQ>cmd/{{.ServerName}}/main.go<BQ> → <BQ>internal/server/http.go<BQ> → <BQ>internal/routers/router.go<BQ> → <BQ>internal/handler<BQ> → <BQ>...<BQ>

其中 handler 层主要负责 API 处理，若需处理更复杂业务逻辑，建议在 handler 和 dao(或repo) 之间额外添加业务逻辑层（如 <BQ>service<BQ>、<BQ>logic<BQ> 或 <BQ>biz<BQ> 等，自己定义）。

## 快速开始

### 1. 生成 openapi 文档

<BQ><BQ><BQ>bash
make proto
<BQ><BQ><BQ>

注：仅当新增或修改 API 时需要执行该命令，API 未变更时无需重复执行。

### 2. 编译和运行

<BQ><BQ><BQ>bash
make run
<BQ><BQ><BQ>

### 3. 测试 API

在浏览器访问 [http://localhost:8080/apis/swagger/index.html](http://localhost:8080/apis/swagger/index.html)，测试 HTTP API。

## 开发指南

点击查看详细的 [**开发指南**](https://go-sponge.com/zh/guide/web/based-on-protobuf.html)。

`

	//nolint
	grpcServerReadmeTmplRaw = `## 技术栈

- 编程语言: go
- RPC 框架: gRPC
- 序列化协议: protobuf
- 配置管理: viper
- 日志: zap
- ORM: {{.ORMType}}
- 数据库: {{.DatabaseType}}
- 缓存: go-redis
- 监控: prometheus+grafana
- 链路追踪: opentracing+jaeger
- 服务注册与发现: consul/etcd/nacos
- 其他: ...

## 目录结构

<BQ><BQ><BQ>text
.
├─ api                          # API 协议定义目录(proto/OpenAPI 等)
│   └─ {{.ServerName}}                     # 服务名称
│       └─ v1                   # API 版本 v1(存放 proto 文件和生成的 pb.go 等)
├─ cmd                          # 应用程序入口目录
│   └─ {{.ServerName}}                     # 服务名称
│       ├─ initial              # 初始化逻辑(如配置加载、服务初始化等)
│       └─ main.go              # 主程序入口文件
├─ configs                      # 配置文件目录(yaml 格式配置模板)
├─ deployments                  # 部署相关脚本(二进制、Docker、K8S 部署)
├─ docs                         # 项目文档(API 文档、设计文档等)
├─ internal                     # 内部实现代码(对外不可见)
│   ├─ cache                    # 缓存相关实现(Redis 或本地内存缓存封装)
│   ├─ config                   # 配置解析和结构体定义
│   ├─ dao                      # 数据访问层(Database Access Object)
│   ├─ ecode                    # 错误码定义
│   ├─ model                    # 数据模型/实体定义
│   ├─ server                   # 服务启动
│   └─ service                  # 业务逻辑处理层
├─ scripts                      # 实用脚本(如代码生成、构建、运行、部署等)
├─ third_party                  # 第三方 Protobuf 依赖/工具
├─ go.mod                       # Go 模块定义文件(声明依赖)
├─ go.sum                       # Go 模块校验文件(自动生成)
├─ Makefile                     # 项目构建自动化脚本
└─ README.md                    # 项目说明文档
<BQ><BQ><BQ>

代码采用分层架构，完整调用链路如下：

<BQ>cmd/{{.ServerName}}/main.go<BQ> → <BQ>internal/server/grpc.go<BQ> → <BQ>internal/service<BQ> → <BQ>internal/dao<BQ> → <BQ>internal/model<BQ>

其中 service 层主要负责 API 处理，若需处理更复杂业务逻辑，建议在 service 和 dao 之间额外添加业务逻辑层（如 <BQ>logic<BQ> 或 <BQ>biz<BQ> 等，自己定义）。

## 快速开始

### 1. 生成与合并 api 相关代码

<BQ><BQ><BQ>bash
make proto
<BQ><BQ><BQ>

注：仅当新增或修改 API 时需要执行该命令，API 未变更时无需重复执行。

### 2. 编译和运行

<BQ><BQ><BQ>bash
make run
<BQ><BQ><BQ>

### 3. 测试 API

使用 VSCode 或 GoLand 等 IDE 打开文件 <BQ>internal/service/xxx_client_test.go<BQ>，测试或者压测 gRPC API。

## 开发指南

点击查看详细的 [**开发指南**](https://go-sponge.com/zh/guide/grpc/based-on-sql.html)。

`

	//nolint
	grpcPbServerReadmeTmplRaw = `## 技术栈

- 编程语言: go
- RPC 框架: gRPC
- 序列化协议: protobuf
- 配置管理: viper
- 日志: zap
- 监控: prometheus+grafana
- 链路追踪: opentracing+jaeger
- 服务注册与发现: consul/etcd/nacos
- 其他: ...

## 目录结构

<BQ><BQ><BQ>text
.
├─ api                          # API 协议定义目录(proto/OpenAPI 等)
│   └─ {{.ServerName}}                     # 服务名称
│       └─ v1                   # API 版本 v1(存放 proto 文件和生成的 pb.go 等)
├─ cmd                          # 应用程序入口目录
│   └─ {{.ServerName}}                     # 服务名称
│       ├─ initial              # 初始化逻辑(如配置加载、服务初始化等)
│       └─ main.go              # 主程序入口文件
├─ configs                      # 配置文件目录(yaml 格式配置模板)
├─ deployments                  # 部署相关脚本(二进制、Docker、K8S 部署)
├─ docs                         # 项目文档(API 文档、设计文档等)
├─ internal                     # 内部实现代码(对外不可见)
│   ├─ config                   # 配置解析和结构体定义
│   ├─ ecode                    # 错误码定义
│   ├─ server                   # 服务启动
│   └─ service                  # 业务逻辑处理层
├─ scripts                      # 实用脚本(如代码生成、构建、运行、部署等)
├─ third_party                  # 第三方 Protobuf 依赖/工具
├─ go.mod                       # Go 模块定义文件(声明依赖)
├─ go.sum                       # Go 模块校验文件(自动生成)
├─ Makefile                     # 项目构建自动化脚本
└─ README.md                    # 项目说明文档
<BQ><BQ><BQ>

代码采用分层架构，完整调用链路如下：

<BQ>cmd/{{.ServerName}}/main.go<BQ> → <BQ>internal/server/grpc.go<BQ> → <BQ>internal/service<BQ> → <BQ>...<BQ>

其中 service 层主要负责 API 处理，若需处理更复杂业务逻辑，建议在 service 和 dao(或repo) 之间额外添加业务逻辑层（如 <BQ>logic<BQ> 或 <BQ>biz<BQ> 等，自己定义）。

## 快速开始

### 1. 生成与合并 api 相关代码

<BQ><BQ><BQ>bash
make proto
<BQ><BQ><BQ>

注：仅当新增或修改 API 时需要执行该命令，API 未变更时无需重复执行。

### 2. 编译和运行

<BQ><BQ><BQ>bash
make run
<BQ><BQ><BQ>

### 3. 测试 API

使用 VSCode 或 GoLand 等 IDE 打开文件 <BQ>internal/service/xxx_client_test.go<BQ>，测试或性能压测 gRPC API。

## 开发指南

点击查看详细的 [**开发指南**](https://go-sponge.com/zh/guide/grpc/based-on-protobuf.html)。

`

	//nolint
	grpcHTTPServerReadmeTmplRaw = `## 技术栈

- 编程语言: go
- Web 框架: gin
- RPC 框架: gRPC
- 序列化协议: protobuf
- 配置管理: viper
- 日志: zap
- ORM: {{.ORMType}}
- 数据库: {{.DatabaseType}}
- 缓存: go-redis
- 监控: prometheus+grafana
- 链路追踪: opentracing+jaeger
- 服务注册与发现: consul/etcd/nacos
- 其他: ...

## 目录结构

<BQ><BQ><BQ>text
.
├─ api                          # API 协议定义目录(proto/OpenAPI 等)
│   └─ {{.ServerName}}                     # 服务名称
│       └─ v1                   # API 版本 v1(存放 proto 文件和生成的 pb.go 等)
├─ cmd                          # 应用程序入口目录
│   └─ {{.ServerName}}                     # 服务名称
│       ├─ initial              # 初始化逻辑(如配置加载、服务初始化等)
│       └─ main.go              # 主程序入口文件
├─ configs                      # 配置文件目录(yaml 格式配置模板)
├─ deployments                  # 部署相关脚本(二进制、Docker、K8S 部署)
├─ docs                         # 项目文档(API 文档、设计文档等)
├─ internal                     # 内部实现代码(对外不可见)
│   ├─ cache                    # 缓存相关实现(Redis 或本地内存缓存封装)
│   ├─ config                   # 配置解析和结构体定义
│   ├─ dao                      # 数据访问层(Database Access Object)
│   ├─ ecode                    # 错误码定义
│   ├─ handler                  # 业务逻辑处理层(类似 Controller)
│   ├─ model                    # 数据模型/实体定义
│   ├─ routers                  # 路由定义和中间件
│   ├─ server                   # 服务启动
│   └─ service                  # 业务逻辑处理层
├─ scripts                      # 实用脚本(如代码生成、构建、运行、部署等)
├─ third_party                  # 第三方 Protobuf 依赖/工具
├─ go.mod                       # Go 模块定义文件(声明依赖)
├─ go.sum                       # Go 模块校验文件(自动生成)
├─ Makefile                     # 项目构建自动化脚本
└─ README.md                    # 项目说明文档
<BQ><BQ><BQ>

代码采用分层架构，完整调用链路如下：

1. gRPC 主要调用链路：

	<BQ>cmd/{{.ServerName}}/main.go<BQ> → <BQ>internal/server/grpc.go<BQ> → <BQ>internal/service<BQ> → <BQ>internal/dao<BQ> → <BQ>internal/model<BQ>

2. http 主要调用链路：

	<BQ>cmd/{{.ServerName}}/main.go<BQ> → <BQ>internal/server/http.go<BQ> → <BQ>internal/routers/router.go<BQ> → <BQ>internal/handler<BQ> → <BQ>internal/service<BQ> → <BQ>internal/dao<BQ> → <BQ>internal/model<BQ>

> TIP: 从调用链路上可以看出，http 调用链路与 gRPC 调用链路共用业务逻辑层 service，不需要写两套代码与协议转换。

其中 service 层主要负责 API 处理，若需处理更复杂业务逻辑，建议在 service 和 dao 之间额外添加业务逻辑层（如 <BQ>logic<BQ> 或 <BQ>biz<BQ> 等，自己定义）。

## 快速开始

### 1. 生成与合并 api 相关代码

<BQ><BQ><BQ>bash
make proto
<BQ><BQ><BQ>

注：仅当新增或修改 API 时需要执行该命令，API 未变更时无需重复执行。

### 2. 编译和运行

<BQ><BQ><BQ>bash
make run
<BQ><BQ><BQ>

### 3. 测试 API

- 使用 VSCode 或 GoLand 等 IDE 打开文件 <BQ>internal/service/xxx_client_test.go<BQ>，测试或性能压测 gRPC API。
- 在浏览器访问 [http://localhost:8080/apis/swagger/index.html](http://localhost:8080/apis/swagger/index.html)，测试 HTTP API。

## 开发指南

点击查看详细的 [**开发指南**](https://go-sponge.com/zh/guide/grpc-http/based-on-sql.html)。

`

	//nolint
	grpcHTTPPbServerReadmeTmplRaw = `## 技术栈

- 编程语言: go
- Web 框架: gin
- RPC 框架: gRPC
- 序列化协议: protobuf
- 配置管理: viper
- 日志: zap
- 监控: prometheus+grafana
- 链路追踪: opentracing+jaeger
- 服务注册与发现: consul/etcd/nacos
- 其他: ...

## 目录结构

<BQ><BQ><BQ>text
.
├─ api                          # API 协议定义目录(proto/OpenAPI 等)
│   └─ {{.ServerName}}                     # 服务名称
│       └─ v1                   # API 版本 v1(存放 proto 文件和生成的 pb.go 等)
├─ cmd                          # 应用程序入口目录
│   └─ {{.ServerName}}                     # 服务名称
│       ├─ initial              # 初始化逻辑(如配置加载、服务初始化等)
│       └─ main.go              # 主程序入口文件
├─ configs                      # 配置文件目录(yaml 格式配置模板)
├─ deployments                  # 部署相关脚本(二进制、Docker、K8S 部署)
├─ docs                         # 项目文档(API 文档、设计文档等)
├─ internal                     # 内部实现代码(对外不可见)
│   ├─ config                   # 配置解析和结构体定义
│   ├─ ecode                    # 错误码定义
│   ├─ handler                  # 业务逻辑处理层(类似 Controller)
│   ├─ routers                  # 路由定义和中间件
│   ├─ server                   # 服务启动
│   └─ service                  # 业务逻辑处理层
├─ scripts                      # 实用脚本(如代码生成、构建、运行、部署等)
├─ third_party                  # 第三方 Protobuf 依赖/工具
├─ go.mod                       # Go 模块定义文件(声明依赖)
├─ go.sum                       # Go 模块校验文件(自动生成)
├─ Makefile                     # 项目构建自动化脚本
└─ README.md                    # 项目说明文档
<BQ><BQ><BQ>

代码采用分层架构，完整调用链路如下：

1. gRPC 主要调用链路：

	<BQ>cmd/{{.ServerName}}/main.go<BQ> → <BQ>internal/server/grpc.go<BQ> → <BQ>internal/service<BQ> → <BQ>...<BQ>

2. http 主要调用链路：

	<BQ>cmd/{{.ServerName}}/main.go<BQ> → <BQ>internal/server/http.go<BQ> → <BQ>internal/routers/router.go<BQ> → <BQ>internal/handler<BQ> → <BQ>internal/service<BQ> → <BQ>...<BQ>

> TIP: 从调用链路上可以看出，http 调用链路与 gRPC 调用链路共用业务逻辑层 service，不需要写两套代码与协议转换。

其中 service 层主要负责 API 处理，若需处理更复杂业务逻辑，建议在 service 和 dao(或repo) 之间额外添加业务逻辑层（如 <BQ>logic<BQ> 或 <BQ>biz<BQ> 等，自己定义）。

## 快速开始

### 1. 生成与合并 api 相关代码

<BQ><BQ><BQ>bash
make proto
<BQ><BQ><BQ>

注：仅当新增或修改 API 时需要执行该命令，API 未变更时无需重复执行。

### 2. 编译和运行

<BQ><BQ><BQ>bash
make run
<BQ><BQ><BQ>

### 3. 测试 API

- 使用 VSCode 或 GoLand 等 IDE 打开文件 <BQ>internal/service/xxx_client_test.go<BQ>，测试或性能压测 gRPC API。
- 在浏览器访问 [http://localhost:8080/apis/swagger/index.html](http://localhost:8080/apis/swagger/index.html)，测试 HTTP API。

## 开发指南

点击查看详细的 [**开发指南**](https://go-sponge.com/zh/guide/grpc-http/based-on-protobuf.html)。

`

	//nolint
	grpcGwPbServerReadmeTmplRaw = `## 技术栈

- 编程语言: go
- Web 框架: gin
- RPC 框架: gRPC
- 序列化协议: protobuf
- 配置管理: viper
- 日志: zap
- 监控: prometheus+grafana
- 链路追踪: opentracing+jaeger
- 服务注册与发现: consul/etcd/nacos
- 其他: ...

## 目录结构

<BQ><BQ><BQ>text
.
├─ api                          # API 协议定义目录(proto/OpenAPI 等)
│   └─ {{.ServerName}}                   # 服务名称
│       └─ v1                   # API 版本 v1(存放 proto 文件和生成的 pb.go 等)
├─ cmd                          # 应用程序入口目录
│   └─ {{.ServerName}}                   # 服务名称
│       ├─ initial              # 初始化逻辑(如配置加载、服务初始化等)
│       └─ main.go              # 主程序入口文件
├─ configs                      # 配置文件目录(yaml 格式配置模板)
├─ deployments                  # 部署相关脚本(二进制、Docker、K8S 部署)
├─ docs                         # 项目文档(API 文档、设计文档等)
├─ internal                     # 内部实现代码(对外不可见)
│   ├─ config                   # 配置解析和结构体定义
│   ├─ ecode                    # 错误码定义
│   ├─ routers                  # 路由定义和中间件
│   ├─ server                   # 服务启动
│   └─ service                  # 业务逻辑处理层 (调用 gRPC 服务，聚合数据等)
├─ scripts                      # 实用脚本(如代码生成、构建、运行、部署等)
├─ third_party                  # 第三方 Protobuf 依赖/工具
├─ go.mod                       # Go 模块定义文件(声明依赖)
├─ go.sum                       # Go 模块校验文件(自动生成)
├─ Makefile                     # 项目构建自动化脚本
└─ README.md                    # 项目说明文档
<BQ><BQ><BQ>

代码采用分层架构，完整调用链路如下：

<BQ>cmd/{{.ServerName}}/main.go<BQ> → <BQ>internal/server/http.go<BQ> → <BQ>internal/routers/router.go<BQ> → <BQ>internal/service<BQ>

其中 service 层主要负责 API 调用的具体实现，若有复杂的业务逻辑，建议额外添加专门处理业务逻辑层（如 <BQ>logic<BQ> 或 <BQ>biz<BQ> 等，自己定义）。

## 快速开始

### 1. 生成 openapi 文档

<BQ><BQ><BQ>bash
make proto
<BQ><BQ><BQ>

注：仅当新增或修改 API 时需要执行该命令，API 未变更时无需重复执行。

### 2. 编译和运行

<BQ><BQ><BQ>bash
make run
<BQ><BQ><BQ>

### 3. 测试 API

在浏览器访问 [http://localhost:8080/apis/swagger/index.html](http://localhost:8080/apis/swagger/index.html)，测试 HTTP API。

## 开发指南

点击查看详细的 [**开发指南**](https://go-sponge.com/zh/guide/grpc-gateway/based-on-protobuf.html)。

`
)

type readmeTemp struct {
	ServerType   string
	ServerName   string
	ORMType      string
	DatabaseType string
	ModuleName   string
	RepoType     string
}

func newReadmeTemp(moduleName, serverName, serverType, dbDriver string, suitedMonoRepo bool) *readmeTemp {
	var repoType string
	if suitedMonoRepo {
		repoType = "mono-repo"
	} else {
		if serverType == codeNameHTTP {
			repoType = "monolith"
		} else {
			repoType = "multi-repo"
		}
	}
	var ormType string
	if dbDriver != "" {
		if dbDriver == DBDriverMongodb {
			ormType = DBDriverMongodb
		} else {
			ormType = "gorm"
		}
	}

	return &readmeTemp{
		ServerType:   serverType,
		ServerName:   serverName,
		ORMType:      ormType,
		DatabaseType: dbDriver,
		ModuleName:   moduleName,
		RepoType:     repoType,
	}
}

func (r *readmeTemp) genReadmeContent() (string, error) {
	var (
		err            error
		readmeTemplate *template.Template
		readmeContent  = fmt.Sprintf("# %s (%s, %s)\n\n%s\n", r.ServerName, r.ServerType, r.RepoType, overviewDesc)
	)

	switch r.ServerType {
	case codeNameHTTP:
		readmeTemplate, err = template.New(r.ServerType).Parse(httpServerReadmeTmplRaw)
		if err != nil {
			return readmeContent, err
		}
	case codeNameHTTPPb:
		readmeTemplate, err = template.New(r.ServerType).Parse(httpPbServerReadmeTmplRaw)
		if err != nil {
			return readmeContent, err
		}
	case codeNameGRPC:
		readmeTemplate, err = template.New(r.ServerType).Parse(grpcServerReadmeTmplRaw)
		if err != nil {
			return readmeContent, err
		}
	case codeNameGRPCPb:
		readmeTemplate, err = template.New(r.ServerType).Parse(grpcPbServerReadmeTmplRaw)
		if err != nil {
			return readmeContent, err
		}
	case codeNameGRPCHTTP:
		readmeTemplate, err = template.New(r.ServerType).Parse(grpcHTTPServerReadmeTmplRaw)
		if err != nil {
			return readmeContent, err
		}
	case codeNameGRPCHTTPPb:
		readmeTemplate, err = template.New(r.ServerType).Parse(grpcHTTPPbServerReadmeTmplRaw)
		if err != nil {
			return readmeContent, err
		}
	case codeNameGRPCGW:
		readmeTemplate, err = template.New(r.ServerType).Parse(grpcGwPbServerReadmeTmplRaw)
		if err != nil {
			return readmeContent, err
		}
	}

	builder := strings.Builder{}
	err = readmeTemplate.Execute(&builder, *r)
	if err != nil {
		return readmeContent, fmt.Errorf("readmeTemplate.Execute error: %v", err)
	}
	partContent := strings.ReplaceAll(builder.String(), "<BQ>", "`")

	readmeContent += partContent + deploymentDesc
	return readmeContent, err
}
