# 仓库规则

## 项目结构与模块组织
- `cmd/admin`：服务入口与初始化。
- `internal`：应用分层，包含 `handler`（HTTP 处理）、`logic`（业务逻辑）、`dao`（数据访问）、`model`（数据库模型）、`routers`（路由/中间件）、`server`（HTTP 服务）、`types`（API 类型）。
- `configs`：YAML 配置模板（参见 `configs/admin.yml`）。
- `docs`：Swagger/OpenAPI 产物（由 `make docs` 生成）。
- `deployments` 与 `scripts`：部署与构建脚本。

## 构建、测试与开发命令
- `make run`：构建并运行服务（可选 `Config=configs/dev.yml`）。
- `make docs`：API 变更后更新 Swagger 文档。
- `make test`：运行 `go test -short`。
- `make cover`：生成覆盖率报告。
- `make ci-lint`：执行 `gofmt` 并运行 `golangci-lint`。
- `make build`：构建 Linux amd64 二进制（`cmd/admin`）。

## 编码风格与命名规范
- 使用标准 Go 格式化与导入顺序（`gofmt` / `goimports`）。
- 包名小写且简短，导出标识符使用驼峰命名。
- 遵守分层边界，`handler` 不可直接访问 `model`，须通过 `dao`。
- 行长度保持在 200 字符以内（满足 `lll` 规则）。
- 保持与现有 Sponge 模板风格一致，优先顺着现有骨架扩展，不随意切换到另一套架构或编码手法。

## 测试规范
- 测试文件与代码同目录，命名为 `*_test.go`。
- 运行 `make test` 作为常规验证步骤。

## 提交与 PR 约定
- 使用简洁的祈使句提交信息（例如：`add order cache metrics`）。
- PR 需描述清楚变更，若有 API 变更需提交 `docs/` 更新。

## 配置与部署
- 本地配置位于 `configs/*.yml`，不要提交敏感信息。
- 部署方式参见 `deployments/`（binary、Docker、K8S）。

## 项目内快捷命令（仅当前仓库）
- 这些快捷命令只在本项目会话中生效，不影响其他项目。
- `gen <表名>`：等价于 `生成代码，<表名>`，完整执行下方“生成代码使用说明”的全部步骤与失败即停规则。

## 编辑器规则（Go + Sponge）
- 遵循 Sponge 项目惯例，`handler` 保持轻薄，业务逻辑放到 `internal/logic`。
- 访问数据库模型必须通过 `internal/dao`，禁止直接操作 `internal/model`。
- API 类型统一放在 `internal/types`，API 变更后运行 `make docs`。
- 路由与中间件沿用 `internal/routers`，服务启动沿用 `internal/server`。
- 新增功能涉及配置时，更新 `configs/` 并补充配置说明。
- 控制行长度（不超过 200 字符）。
- 优先复用既有中间件、路由与服务启动逻辑。
- 代码风格遵循 Go 规范与 `gofmt`/`goimports`。
- 新文件放入正确分层目录，包名小写，导出符号驼峰命名。
- 避免直接修改生成文件；如必须修改，先说明原因。
- `handler` 仅负责参数绑定、路径参数解析、调用 `logic`、统一响应，不写业务逻辑。
- `handler` 统一使用 `response.Success`、`response.Error`、`response.Output` 返回结果。
- `handler` 向下传递 `context.Context`，通过 `middleware.WrapCtx(c)` 包装，不直接把 `*gin.Context` 传入 `logic` 或 `dao`。
- `logic` 负责业务编排、DTO 与 `model` 转换，以及将底层错误转换为 `ecode`/业务错误。
- DTO 与 `model` 转换优先复用 `copier.Copy`，无法覆盖的字段再显式补充。
- `dao` 负责数据库、缓存、singleflight、防穿透等数据访问细节，上层不要下沉实现细节。
- 路由注册遵循 `internal/routers` 中 `init() + append(...)` 的组织方式，避免集中硬编码在单一入口文件。
- 测试优先贴近现有 Sponge `gotest` / `sqlmock` 风格，保持与仓库现有测试组织一致。
- 业务语义、字段含义、状态值、接口行为、兼容策略、脏数据处理等任何不清楚或不确定的逻辑，必须先问我确认，不要自己猜，不要擅自定义规则。

## 生成代码使用说明
- 方式一：你输入 `生成代码，<表名>`
- 根据 configs/admin.yml 内容找到database.mysql.dsn ,截取?之前内容，得到 `<dsn主串>`。
- 执行顺序：
  1. 命令执行后，通过本机 `mysql` CLI 获取该表结构。确认表是否存在
    - 固定执行以下检查语句：
      - `SHOW TABLES LIKE '<表名>';`
      - `SHOW CREATE TABLE <表名>;`
      - `SHOW INDEX FROM <表名>;`（用于识别唯一索引并同步到 `logic/dao` 约束）
    - 失败即停：任一语句执行失败或表不存在时，立即停止流程，不执行 `sponge web http`。
  2. 直接执行：
     `sponge web http --module-name=admin --server-name=admin --project-name=admin --repo-addr= --db-driver=mysql --db-dsn=<dsn主串>;prefix=<x> --db-table=<表名> --embed=true --suited-mono-repo=false --extended-api=false --out=$(pwd)`
    - --db-dsn 的值 <dsn主串>然后拼接;prefix=x, x是我给你的表名称的第一个_的值包括_ 比如 t_goods, 那么x就是t_。
    - --db-table <表名>。

- 方式二：你输入 `帮我生成代码 <CREATE TABLE ...>`
- 执行顺序：
  1. 先通过本机 `mysql` CLI 执行你提供的 `CREATE TABLE` 语句创建表。
  2. 再按方式一执行 `sponge web http` 生成代码。

- 根据表字段、索引、约束继续定义并实现业务逻辑。


## 生成代码后强制规则（binding + 业务确认）
- 在执行 `sponge web http` 并通过本机 `mysql` CLI 获取真实表结构后，必须立刻基于表结构同步更新 `internal/types/*_types.go` 的 `binding` 标签，不允许保留空 `binding:""`。
- 列表搜索条件强制规则：
  1. 基于真实表结构生成列表查询条件时，`created_at`、`updated_at`、`deleted_at` 三个字段固定不加入通用搜索条件。
  2. 除上述 3 个字段外，其余字段都必须在 `ListXXXRequest` 与对应 `logic.List` 中提供可用搜索条件（至少一种）。
  3. 默认映射规则：
    - `varchar`/`char`/`text`：模糊搜索（`like`）；
    - `tinyint`/`int`/`bigint`：精确搜索（`=`）；
    - `decimal`：精确搜索（`=`），如有业务需求可扩展 `Min/Max` 区间；
    - `datetime`/`date`/`timestamp`：区间搜索（`start/end`）；
    - `json`：默认不做通用搜索，除非业务明确要求。
  4. 若用户明确指定某字段的搜索策略（如“不可搜索/必须精确/必须模糊”），以用户规则优先覆盖默认映射。
- 生成后最小验收（进入业务确认前必须完成）：
  1. 扫描并确认 `internal/types/*_types.go` 不存在 `binding:""` 空标签；
  2. 运行 `make test`（至少一次）并确认无失败用例。
- `binding` 生成与 `go-playground/validator` 对齐，至少包含以下映射：
  - `NOT NULL` 且无默认值字段：创建请求使用 `required`。
  - `varchar(n)`/`char(n)`：使用 `max=n`；如业务要求固定长度可用 `len=n`。
  - `tinyint`/`int` 等枚举语义字段（如状态、来源、支付类型）：使用 `oneof=...`（枚举值来自表注释或业务规则）。
  - 金额/小数字段：若请求结构体用字符串承接，使用 `numeric`；若业务要求必须大于等于 0，增加对应数值校验（如自定义校验或在 logic 强校验）。
  - 可选更新字段（Update 接口）：使用 `omitempty,...` 组合，避免误判必填。
- 唯一索引（如 `order_no`）不能仅依赖 `binding`；必须在 `logic/dao` 做唯一性校验或冲突错误转换。
- 所有根据表结构可确定的硬约束（类型、可空、默认值、唯一索引）必须同时体现在：
  1. `types` 的 `binding` 标签；
  2. `logic` 的参数与业务校验；
  3. `dao` 的更新白名单/更新策略。

- 强制思考与确认：生成骨架后、写业务逻辑前，必须先向用户确认“业务语义规则”，至少确认以下 5 项，未确认前不得臆断实现：
  1. 状态字段及允许流转（状态机）；
  2. 金额/数量计算公式；
  3. 创建后不可修改字段；
  4. 触发动作条件（取消、退款、完成等）；
  5. 唯一性/幂等规则。
- 若用户未提供完整业务规则，只允许先完成“硬约束同步（binding + 基础参数校验）”，并明确标注待确认项；禁止自行杜撰业务规则。
