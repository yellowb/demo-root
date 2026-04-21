# AGENTS.md

## 项目定位
- 这是 Agent Harness 分享用的 Todo List baseline codebase。
- 当前仓库的目标是提供 **live demo 前的基线项目**。
- 在没有明确要求前，**不要提前实现** `priority`、过滤、搜索、标签、截止日期、登录、多人协作。
- 做 live demo feature 前，先阅读 `docs/live-demo-context.md`，理解这次演示服务的是 Agent Harness 工作流，而不只是功能开发。

## 技术栈与结构
- 前端：`frontend/`，React + Vite + TypeScript
- 后端：`backend/`，Go + Gin + `database/sql`
- 数据库：`backend/data/todos.db`，SQLite 本地文件
- 前端开发端口：`5173`
- 后端开发端口：`8080`
- 前端通过 `/api` 访问后端

## 开发命令
- 首次进入仓库先执行：`make setup`
- 启动开发环境：`make dev`
  - 默认会在前端就绪后自动打开 `http://localhost:5173/`
  - 如需关闭，使用：`AUTO_OPEN_BROWSER=0 make dev`
- 停止开发环境并释放端口：`make stop`
- 运行后端 lint：`make lint`
- 运行 OpenSpec 规格验证：`make validate-specs`
- 运行验证：`make test`
- 重置数据库并恢复种子数据：`make reset-db`

## 开发约束
- 保持项目轻量：不要引入 Docker、MySQL、Redis、Postgres 或外部云服务依赖。
- 保持目录边界清楚：
  - `frontend/src/api/` 放最小 API 客户端
  - `frontend/src/features/todos/` 放 Todo 主业务 UI
  - `backend/internal/httpapi/` 放路由和 handler
  - `backend/internal/todos/` 放模型、仓储、服务逻辑
  - `backend/internal/store/` 放 SQLite、schema、seed
- UI 方向保持“简洁但像真实产品”，不要做花哨动画或复杂多页流转。
- 编辑 Todo 优先保持 inline 或轻量交互，不要改成复杂路由页面。

## 数据与 API 约束
- 基线 Todo 字段只有：
  - `id`
  - `title`
  - `notes`
  - `completed`
  - `created_at`
  - `updated_at`
- `title` 必填，`notes` 可空，默认按 `created_at desc` 展示。
- 基线 API 只包含：
  - `GET /api/health`
  - `GET /api/todos`
  - `POST /api/todos`
  - `PATCH /api/todos/:id`
  - `DELETE /api/todos/:id`

## 修改时的要求
- 任何会影响前后端联动的改动，都要同步更新：
  - 后端 handler / repository / schema
  - 前端 API client / UI
  - 种子数据
  - README（如果启动或验证方式变化）
- 完成前必须运行 `make test`。
- `make test` 包含 OpenSpec 规格验证、后端 `golangci-lint`、后端测试、前端类型检查和前端构建。
- 本 repo 配置了 Codex `Stop` hook：当 Codex App 在本仓库根目录运行且检测到未提交改动时，hook 会自动触发 `make test`。
- 最终回复必须说明验证结果；如果验证由 Codex hook 触发，也要明确说明。
- 如果改动影响演示初始状态，确保 `make reset-db` 后仍然有适合 demo 的种子数据。
