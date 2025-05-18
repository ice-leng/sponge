## [English](../README.md) | 简体中文

<br>

**sponge** 是一个强大且易用的 `Go` 开发框架，其核心理念是通过解析 `SQL`、`Protobuf`、`JSON` 文件逆向生成模块化的代码，这些模块代码可灵活组合成多种类型的完整后端服务。sponge 采用模块化架构并深度集成 AI 助手，显著提升开发效率、降低技术门槛，以"低代码"方式轻松构建高性能、高可用的后端服务体系。

<br>

### 适用场景

sponge 适用于快速构建多种类型的高性能后端服务，适用场景如下：

- 开发企业内部 API 服务。
- 快速构建微服务 (Microservices)。
- 搭建后台管理系统 API。
- 构建 gRPC 服务进行服务间通信。
- 作为 Go 初学者或团队学习 Go 项目结构和最佳实践的起点。
- 需要提高开发效率、统一开发规范的团队。
- 云原生开发。

此外，开发者还可以通过自定义模板，生成满足业务需求的各类代码。

<br>

### 主要特点

1. **一键生成完整后端服务代码**  
   对于仅需 `CRUD API` 的 `Web` 或 `gRPC` 服务，无需编写任何 `Go` 代码。只需连接数据库(如 `MySQL`、`MongoDB`、`PostgreSQL`、`SQLite`)，即可一键生成完整后端服务代码，并轻松部署到 Linux 服务器、Docker 或 Kubernetes 上。

2. **高效开发通用服务**  
   开发通用的 `Web`、`gRPC`、`HTTP+gRPC` 或 `gRPC Gateway` 服务，只需专注于以下三部分：
    - 数据库表的定义；
    - 在 Protobuf 文件中定义 API 描述信息；
    - 在生成的模板中，使用内置AI助手或人工编写业务逻辑代码。  

   服务的框架代码和 CRUD API 代码均由 sponge 自动生成。

3. **支持自定义模板，灵活扩展**  
   sponge 支持通过自定义模板生成项目所需的多种代码类型，不局限于 `Go` 语言。例如：
    - 后端代码；
    - 前端代码；
    - 配置文件；
    - 测试代码；
    - 构建和部署脚本等。

4. **在页面生成代码，简单易用**  
   sponge 提供在页面生成代码，避免了复杂的命令行操作，只需在页面上简单的填写参数即可一键生成代码。

5. **与 AI 助手协同开发，形成开发闭环**  
   sponge 与 内置的 AI 助手(DeepSeek、ChatGPT、Gemini)深度融合，形成一套完整的高效开发解决方案：
    - **sponge**：负责基础设施代码生成(服务框架、CRUD API 接口、自定义 API 接口代码(缺少业务逻辑)等)。
    - **AI助手**：专注业务逻辑实现(表结构 DDL 设计、自定义 API 接口定义、业务逻辑实现代码)。

<br>

### 快速开始

1. **安装 sponge**

   支持在 windows、mac、linux、docker 环境下安装 sponge，点击查看 [**安装 sponge 说明**](https://github.com/go-dev-frame/sponge/blob/main/assets/install-cn.md)。

2. **打开生成代码 UI 页面**

   安装完成后，执行命令打开 sponge UI 页面：

   ```bash
   sponge run
   ```

   在本地浏览器访问 `http://localhost:24631`，在页面上操作生成代码，如下图所示：

   <p align="center">
   <img width="1500px" src="https://raw.githubusercontent.com/go-dev-frame/sponge/main/assets/sponge-ui.png">
   </p>

   > 如果想要在跨主机的浏览器上访问，启动UI时需要指定宿主机ip或域名，示例 `sponge run -a http://your_host_ip:24631`。

<br>

### 组件

sponge 内置了丰富的组件(按需使用)：

| 组件 | 使用示例 |
| :--- | :-------- |
| Web 框架 [gin](https://github.com/gin-gonic/gin) | [gin 示例](https://github.com/go-dev-frame/sponge/blob/main/internal/routers/routers.go#L32)<br>[gin 中间件示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/gin/middleware/README.md) |
| RPC 框架 [gRPC](https://github.com/grpc/grpc-go) | [gRPC 示例](https://github.com/go-dev-frame/sponge/blob/main/internal/server/grpc.go#L312)<br>[gRPC 拦截器示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/grpc/interceptor/README.md) |
| 配置解析 [viper](https://github.com/spf13/viper) | [示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/conf/README.md#example-of-use) |
| 日志 [zap](https://github.com/uber-go/zap) | [示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/logger/README.md#example-of-use) |
| ORM 框架 [gorm](https://github.com/go-gorm/gorm), [mongo-go-driver](https://github.com/mongodb/mongo-go-driver) | [gorm 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/sgorm/README.md#examples-of-use)<br>[mongodb 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/mgo/README.md#example-of-use) |
| 缓存 [go-redis](https://github.com/go-redis/redis), [ristretto](https://github.com/dgraph-io/ristretto) | [go-redis 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/goredis/README.md#example-of-use)<br>[ristretto 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/cache/memory.go#L25) |
| 自动化api文档 [swagger](https://github.com/swaggo/swag), [protoc-gen-openapiv2](https://github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2) | - |
| 鉴权 [jwt](https://github.com/golang-jwt/jwt) | [jwt 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/jwt/README.md#example-of-use)<br>[gin 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/gin/middleware/README.md#jwt-authorization-middleware)<br>[gRPC 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/grpc/interceptor/README.md#jwt-authentication-interceptor) |
| 校验 [validator](https://github.com/go-playground/validator), [protoc-gen-validate](https://github.com/bufbuild/protoc-gen-validate) | [validator 示例](https://github.com/go-dev-frame/sponge/blob/main/internal/types/userExample_types.go#L17)<br>[protoc-gen-validate 示例](https://github.com/go-dev-frame/sponge/blob/main/api/serverNameExample/v1/userExample.proto#L156) |
| Websocket [gorilla/websocket](https://github.com/gorilla/websocket) | [示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/ws/README.md#example-of-use) |
| 定时任务 [cron](https://github.com/robfig/cron) | [示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/gocron/README.md#example-of-use) |
| 消息队列 [rabbitmq](https://github.com/rabbitmq/amqp091-go), [kafka](https://github.com/IBM/sarama) | [rabbitmq 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/rabbitmq/README.md#example-of-use)<br>[kafka 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/kafka/README.md#example-of-use) |
| 分布式事务管理器 [dtm](https://github.com/dtm-labs/dtm) | [dtm 服务发现示例](https://github.com/go-dev-frame/sponge_examples/blob/main/_11_sponge-dtm-service-registration-discovery/internal/rpcclient/dtmservice.go#L31)<br>[使用 dtm 秒杀抢购示例](https://github.com/go-dev-frame/sponge_examples/blob/main/_12_sponge-dtm-flashSale/grpc%2Bhttp/internal/service/flashSale.go#L67) |
| 分布式锁 [dlock](https://github.com/go-dev-frame/sponge/tree/main/pkg/dlock) | [示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/dlock/README.md#example-of-use) |
| 自适应限流 [ratelimit](https://github.com/go-dev-frame/sponge/tree/main/pkg/shield/ratelimit) | [gin 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/gin/middleware/README.md#rate-limiter-middleware)<br>[gRPC 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/gin/middleware/README.md#rate-limiter-interceptor) | |
| 自适应熔断 [circuitbreaker](https://github.com/go-dev-frame/sponge/tree/main/pkg/shield/circuitbreaker) | [gin 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/gin/middleware/README.md#circuit-breaker-middleware)<br>[gRPC 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/grpc/interceptor/README.md#circuit-breaker-interceptor) | |
| 链路追踪 [opentelemetry](https://github.com/open-telemetry/opentelemetry-go) | [gin 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/gin/middleware/README.md#tracing-middleware)<br>[gRPC 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/grpc/interceptor/README.md#tracing-interceptor)<br>[跨服务链路追踪示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/tracer/example-cn.md) |
| 监控 [prometheus](https://github.com/prometheus/client_golang/prometheus), [grafana](https://github.com/grafana/grafana) | [gin 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/gin/middleware/metrics/README.md#example-of-use)<br>[gRPC 示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/grpc/metrics/README.md#example-of-use)<br>[web 和 gRPC 监控示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/grpc/metrics/monitor-example-cn.md) |
| 服务注册与发现 [etcd](https://github.com/etcd-io/etcd), [consul](https://github.com/hashicorp/consul), [nacos](https://github.com/alibaba/nacos) | [服务注册示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/servicerd/registry/README.md#example-of-use)<br>[服务发现示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/servicerd/discovery/README.md#example-of-use) |
| 自适应采集 [profile](https://go.dev/blog/pprof) | [示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/prof/go-profile-cn.md) |
| 资源统计 [gopsutil](https://github.com/shirou/gopsutil) | [示例](https://github.com/go-dev-frame/sponge/tree/main/pkg/stat#example-of-use) |
| 配置中心 [nacos](https://github.com/alibaba/nacos) | [示例](https://go-sponge.com/zh/component/config-center.html) |
| 代码质量检查 [golangci-lint](https://github.com/golangci/golangci-lint) | - |
| 持续集成部署 CI/CD [kubernetes](https://github.com/kubernetes/kubernetes), [docker](https://www.docker.com/), [jenkins](https://github.com/jenkinsci/jenkins) | [示例](https://go-sponge.com/zh/deployment/kubernetes.html) |
| 生成项目业务架构图 [spograph](https://github.com/go-dev-frame/spograph) | [示例](https://github.com/go-dev-frame/spograph?tab=readme-ov-file#example-of-use) |
| 生成自定义代码 [go template](https://pkg.go.dev/text/template@go1.23.3) | [json 示例](https://go-sponge.com/zh/guide/customize/template-json.html)<br>[sql 示例](https://go-sponge.com/zh/guide/customize/template-sql.html)<br>[protobuf 示例](https://go-sponge.com/zh/guide/customize/template-protobuf.html) |
| AI助手 [DeepSeek](https://deepseek.com), [ChatGPT](https://chatgpt.com), [Gemini](https://gemini.google.com) | [示例](https://github.com/go-dev-frame/sponge/blob/main/cmd/sponge/commands/assistant/generate.go#L44) |

<br>

### 代码生成引擎

sponge 提供强大的代码生成能力，支持基于`内置模板`和`自定义模板`两种方式快速生成项目所需代码，同时集成`AI 助手`辅助生成业务逻辑代码。

1. sponge 基于内置模板生成代码框架，如下图所示：

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/go-dev-frame/sponge/main/assets/sponge-framework.png">
</p>

<br>

2. sponge 基于自定义模板生成代码框架，如下图所示：

<p align="center">
<img width="1200px" src="https://raw.githubusercontent.com/go-dev-frame/sponge/main/assets/template-framework.png">
</p>

<br>

3. sponge 基于函数及注释生成业务逻辑代码框架，如下图所示：

<p align="center">
<img width="1200px" src="https://raw.githubusercontent.com/go-dev-frame/sponge/main/assets/ai-assistant-framework.png">
</p>

<br>

### 微服务框架

sponge 是一个现代化的 Go 微服务框架，它采用典型的微服务分层架构，内置了丰富的服务治理功能，帮助开发者快速构建和维护复杂的微服务系统，框架结构如下图所示：

<p align="center">
<img width="1000px" src="https://raw.githubusercontent.com/go-dev-frame/sponge/main/assets/microservices-framework.png">
</p>

<br>

创建的 HTTP 和 gRPC 服务代码的性能测试： 50个并发，总共100万个请求。

![http-server](https://raw.githubusercontent.com/zhufuyi/microservices_framework_benchmark/main/test/assets/http-server.png)

![grpc-server](https://raw.githubusercontent.com/zhufuyi/microservices_framework_benchmark/main/test/assets/grpc-server.png)

点击查看[**测试代码**](https://github.com/zhufuyi/microservices_framework_benchmark)。

<br>

### sponge 开发指南

欢迎查阅 [Sponge 开发项目的完整技术文档](https://go-sponge.com/zh/)，该文档详尽涵盖了代码生成、开发流程、系统配置及部署方案等核心内容。

<br>

### 目录结构

sponge 创建的服务代码目录结构遵循 [project-layout](https://github.com/golang-standards/project-layout)。

sponge 支持创建 `单体应用单体仓库(monolith)`、`微服务多仓库(multi-repo)`、`微服务单体仓库(mono-repo)`三种类型的项目代码结构。

1. 创建`单体应用单体仓库(monolith)`或`微服务多仓库(multi-repo)`代码目录结构如下：

```bash
   .
   ├── api            # protobuf文件和生成的*pb.go目录
   ├── assets         # 其他与资源库一起使用的资产(图片、logo等)目录
   ├── cmd            # 程序入口目录
   ├── configs        # 配置文件的目录
   ├── deployments    # 裸机、docker、k8s部署脚本目录
   ├── docs           # 设计文档和界面文档目录
   ├── internal       # 项目内部代码目录
   │    ├── cache        # 基于业务包装的缓存目录
   │    ├── config       # Go结构的配置文件目录
   │    ├── dao          # 数据访问目录
   │    ├── database     # 数据库目录
   │    ├── ecode        # 自定义业务错误代码目录
   │    ├── handler      # http的业务功能实现目录
   │    ├── model        # 数据库模型目录
   │    ├── routers      # http路由目录
   │    ├── rpcclient    # 连接grpc服务的客户端目录
   │    ├── server       # 服务入口，包括http、grpc等
   │    ├── service      # grpc的业务功能实现目录
   │    └── types        # http的请求和响应类型目录
   ├── pkg            # 外部应用程序可以使用的库目录
   ├── scripts        # 执行脚本目录
   ├── test           # 额外的外部测试程序和测试数据
   ├── third_party    # 依赖第三方protobuf文件或其他工具的目录
   ├── Makefile       # 开发、测试、部署相关的命令集合
   ├── go.mod         # go 模块依赖关系和版本控制文件
   └── go.sum         # go 模块依赖项的密钥和校验文件
```

<br>

2. 创建`微服务单体仓库(mono-repo)`(大仓库)代码目录结构如下：

```bash
   .
   ├── api
   │    ├── server1       # 服务1的protobuf文件和生成的*pb.go目录
   │    ├── server2       # 服务2的protobuf文件和生成的*pb.go目录
   │    ├── server3       # 服务3的protobuf文件和生成的*pb.go目录
   │    └── ...
   ├── server1        # 服务1的代码目录，与微服务多仓库(multi-repo)目录结构基本一样
   ├── server2        # 服务2的代码目录，与微服务多仓库(multi-repo)目录结构基本一样
   ├── server3        # 服务3的代码目录，与微服务多仓库(multi-repo)目录结构基本一样
   ├── ...
   ├── third_party    # 依赖的第三方protobuf文件
   ├── go.mod         # go 模块依赖关系和版本控制文件
   └── go.sum         # go 模块依赖项的密钥和校验和文件
```

<br>

### 代码示例

#### sponge 创建服务代码示例

- [基于sql创建web服务(包括CRUD)](https://github.com/go-dev-frame/sponge_examples/tree/main/1_web-gin-CRUD)
- [基于sql创建grpc服务(包括CRUD)](https://github.com/go-dev-frame/sponge_examples/tree/main/2_micro-grpc-CRUD)
- [基于protobuf创建web服务](https://github.com/go-dev-frame/sponge_examples/tree/main/3_web-gin-protobuf)
- [基于protobuf创建grpc服务](https://github.com/go-dev-frame/sponge_examples/tree/main/4_micro-grpc-protobuf)
- [基于protobuf创建grpc网关服务](https://github.com/go-dev-frame/sponge_examples/tree/main/5_micro-gin-rpc-gateway)
- [基于protobuf创建grpc+http服务](https://github.com/go-dev-frame/sponge_examples/tree/main/_10_micro-grpc-http-protobuf)

#### 分布式事务示例

- [简单的分布式订单系统](https://github.com/go-dev-frame/sponge_examples/tree/main/9_order-grpc-distributed-transaction)
- [秒杀抢购活动](https://github.com/go-dev-frame/sponge_examples/tree/main/_12_sponge-dtm-flashSale)
- [电商系统](https://github.com/go-dev-frame/sponge_examples/tree/main/_14_eshop)

####  sponge+AI 助手协同开发示例

- [家电零售管理](https://github.com/go-dev-frame/sponge_examples/tree/main/_15_appliance_store)

#### sponge 开发项目示例

- [社区后端服务](https://github.com/go-dev-frame/sponge_examples/tree/main/7_community-single)
- [单体服务拆分为微服务](https://github.com/go-dev-frame/sponge_examples/tree/main/8_community-cluster)

<br>
<br>

如果对您有帮助给个star⭐，欢迎加入**go sponge微信群交流**，加微信(备注`sponge`)进群。

<img width="300px" src="https://raw.githubusercontent.com/go-dev-frame/sponge/main/assets/wechat-group.jpg">
