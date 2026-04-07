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

## 测试规范
- 测试文件与代码同目录，命名为 `*_test.go`。
- 运行 `make test` 作为常规验证步骤。

## 提交与 PR 约定
- 使用简洁的祈使句提交信息（例如：`add order cache metrics`）。
- PR 需描述清楚变更，若有 API 变更需提交 `docs/` 更新。

## 配置与部署
- 本地配置位于 `configs/*.yml`，不要提交敏感信息。
- 部署方式参见 `deployments/`（binary、Docker、K8S）。

## 编辑器规则（Go + Sponge）
- 遵循 Sponge 项目惯例，`handler` 保持轻薄，业务逻辑放到 `internal/logic`。
- 访问数据库模型必须通过 `internal/dao`，禁止直接操作 `internal/model`。
- API 类型统一放在 `internal/types`，API 变更后运行 `make docs`。
- 路由与中间件沿用 `internal/routers`，服务启动沿用 `internal/server`。
- 新增功能涉及配置时，更新 `configs/` 并补充配置说明。
- 控制行长度（不超过 200 字符）。
- 优先复用既有中间件、路由与服务启动逻辑。
- 代码风格遵循_ Go 规范与 `gofmt`/`goimports`。
- 新文件放入正确分层目录，包名小写，导出符号驼峰命名。
- 避免直接修改生成文件；如必须修改，先说明原因。
- 有不清楚或不确定逻辑的地方，一定要问我，不要自己猜。