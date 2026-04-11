# 仓库规则

## 规则分级（必选 / 建议）
- `必选（Gate）`：必须满足，未满足即视为未完成，不得跳过。
- `建议（Best Practice）`：在不与业务需求冲突时默认遵循；若不采用，需说明原因。
- 规则治理：
  - `必选（Gate）`：每新增一条规则，必须同步删除或合并至少一条旧规则，保持规则总量不膨胀。
  - `建议（Best Practice）`：每两周审计一次，清理重复、不可验证、低频条目。

## 项目结构
- `cmd/admin`：服务入口与初始化。
- `internal`：`handler`、`logic`、`dao`、`model`、`routers`、`server`、`types`。
- `configs`：配置模板（参见 `configs/admin.yml`）。
- `docs`：Swagger/OpenAPI 产物。
- `deployments`、`scripts`：部署与构建脚本。

## 常用命令
- `make run`：构建并运行服务（可选 `Config=configs/dev.yml`）。
- `make docs`：生成 Swagger 文档。
- `make test`：运行 `go test -short`。
- `make cover`：覆盖率报告。
- `make ci-lint`：`gofmt` + `golangci-lint`。
- `make build`：构建 Linux amd64 二进制（`cmd/admin`）。

## 开发通用规则（Go + Sponge）
- `必选（Gate）`
  - 仅修改当前任务直接相关文件；发现无关改动或脏变更先暂停并确认。
  - 严格分层：`handler -> logic -> dao`。
  - `handler` 只做绑定/解析/调用 `logic`/统一响应（`response.Success|Error|Output`），通过 `middleware.WrapCtx(c)` 下传 `context.Context`。
  - `logic` 负责业务编排与错误转换；`dao` 负责数据访问细节；禁止跨层直连 `model`。
  - 业务语义不清（字段含义、状态流转、兼容策略、脏数据处理）必须先确认，禁止臆断。
- `建议（Best Practice）`
  - 目录沿用既有结构：`internal/types`、`internal/routers`、`internal/server`。
  - DTO 与 `model` 优先 `copier.Copy`，补字段时显式处理。
  - 测试风格贴近现有 `gotest/sqlmock`。
  - 控制行长（<=200），优先复用现有中间件与启动逻辑。

## 生成代码入口
- 方式一：`gen <表名>` / `生成代码，<表名>`。
  - 从 `configs/admin.yml` 读取 `database.mysql.dsn`，取 `?` 前内容作为 `<dsn主串>`。
- 方式二：`帮我生成代码 <CREATE TABLE ...>`。
  - 先用本机 `mysql` 执行建表 SQL，再按方式一执行生成流程。

## 生成代码 Gate 流程（强制）
- 按顺序执行，任一步失败即停止并反馈“失败点 + 原因 + 下一步建议”。

1. 前置校验（mysql CLI）
  - `SHOW TABLES LIKE '<表名>';`
  - `SHOW CREATE TABLE <表名>;`
  - `SHOW INDEX FROM <表名>;`
  - 任一失败或表不存在：禁止执行 `sponge web http`。

2. 执行生成
  - `sponge web http --module-name=admin --server-name=admin --project-name=admin --repo-addr= --db-driver=mysql --db-dsn=<dsn主串>;prefix=<x> --db-table=<表名> --embed=true --suited-mono-repo=false --extended-api=false --out=$(pwd)`
  - `prefix` 取表名前缀（如 `t_goods` -> `t_`）。

3. 生成后硬约束同步（`types` + `logic` + `dao`）
  - `binding`：`internal/types/*_types.go` 禁止空 `binding:""`。
  - 映射规则：
    - `NOT NULL` 且无默认值 -> `required`
    - `varchar/char` -> `max=n`
    - 枚举语义 -> `oneof`
    - 金额字符串 -> `numeric`
    - 更新接口 -> `omitempty,...`
  - 唯一索引（`Non_unique=0` 且 `Key_name != PRIMARY`）必须落地：
    - `logic.Create`：预查重，命中返回明确业务错误。
    - `logic.Update`：预查重且排除当前 ID，命中返回明确业务错误。
    - `dao`：数据库唯一键冲突（如 MySQL `1062`）统一转换为可识别业务错误。
  - 若唯一索引语义不清（空值、大小写不敏感、联合唯一解释），必须先向用户确认。

4. 列表搜索规则（强制）
  - 黑名单字段：`password/token/secret/salt` 默认禁止搜索，且不得在列表 DTO 返回（确需开放先确认）。
  - `created_at/updated_at/deleted_at` 不做通用搜索。
  - 其余字段至少提供一种搜索能力（精确/模糊/范围其一）。
  - 搜索白名单必须在 `types`、`logic`、`docs/swagger.yaml` 三处一致。

5. 测试与文档验收（强制）
  - 受影响包测试通过（至少 `types/logic/dao`）。
  - 唯一索引相关至少包含：
    - `Create` 冲突测试
    - `Update` 冲突测试（排除自身）
  - 新增列表筛选能力至少补 1 条 `logic.List` 条件下推测试。
  - API/路由变更必须执行 `make docs` 并核对产物。
  - 提交前必须通过 `make test`（建议同时执行 `make ci-lint`）。

## 其他约定
- 命名与格式：包名小写，导出标识符驼峰，统一 `gofmt/goimports`。
- 提交信息：简洁祈使句（如 `add order cache metrics`）。
- PR：描述变更；有 API 变更时包含 `docs/` 更新。
- 配置与部署：`configs/*.yml` 不提交敏感信息，部署参考 `deployments/`（binary、Docker、K8S）。
- 有不清楚或不确定逻辑的地方，一定要问我，不要自己猜。
