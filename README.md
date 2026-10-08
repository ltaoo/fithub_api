# Fithub：Go API、h5 与 admin

本仓库统一维护服务端与两套前端：

- `cmd/`、`internal/`、`config/`、`pkg/`、`migrations/`：Go API 与静态资源服务。
- `h5/`：移动端网页，保留现有 SolidJS 项目与业务改动。
- `admin/`：Timeless shadcn 管理后台。
- `dist/h5/`、`dist/admin/`：前端构建产物。
- `bin/fithub-server`：Go 可执行文件。

## 安装、构建、启动

```sh
pnpm --dir h5 install --frozen-lockfile
pnpm --dir admin install --frozen-lockfile
pnpm run build
pnpm start
```

从仓库根目录启动，默认监听 `:8080`。现有 `.env` 与数据库继续使用；新环境复制 `.env.example` 为 `.env` 并设置数据库和 TOKEN_SECRET_KEY。

| 地址 | 内容 |
| --- | --- |
| `/` | h5，游客进入 `/workout_plan/list` |
| `/admin/` | 管理后台 |
| `/api/` | 业务 API |
| `/health` | 健康检查 |

Go 使用 `STATIC_DIR`（默认 `./dist`）中的 h5/、admin/，支持两套 SPA 子路由刷新。API 优先匹配；缺失静态文件返回 404，未构建页面返回 503。部署携带 binary、dist/、migrations/ 与本地配置即可，运行时无需 Node.js。

admin 使用本地 `../../timeless/packages` 的预构建 Timeless/shadcn，路径从 admin 目录计算，即工作区中的 timeless 仓库。

## 开发与检查

```sh
pnpm run dev:h5       # Vite 3100
pnpm run dev:admin    # Vite 3003，/admin/
pnpm run dev:server  # Go 8080
pnpm run check
pnpm run test:h5
pnpm run test:admin
pnpm run test:server
```

开发代理默认连接 `http://127.0.0.1:8080`；h5 可用 API_URL、admin 可用 ADMIN_API_URL 覆盖。生产前端默认同源访问 /api。两套前端独立保留 package.json、锁文件和依赖目录，根 package.json 仅提供组合命令。

## 游客权限

游客可浏览公开的单次/周期训练计划、动作及公开的配套视频内容。私人和已删除内容不对游客开放；登录用户可查看本人私人计划。

开始训练、应用周期计划、选择参与人、个人记录和创建/编辑仍需要登录。游客点击开始训练先进入登录，登录后返回原计划及查询参数，由用户再次确认开始；登录成功不会自动创建训练记录。后端写接口继续严格验证 token。

## 迁移与验证

源码及游客行为来自 coach_assistant_group/coach_assistant_h5 中已整合的版本，原目录保留。API 的 Git 历史、本地 .env 与现有数据库不被前端目录迁移覆盖；被更新的源文件有迁移临时备份。

非浏览器验证包括两套前端类型检查与构建、Go 全量测试、静态资源/SPA/API 路由及游客鉴权测试、h5 导航权限测试和 admin 组件测试。全量前端测试各有一项迁移前已有的 Form 对象/字符串断言失败，未修改该断言。

遵照项目约定，不自动操作浏览器。请手动验证游客浏览、登录返回计划、开始训练、admin 登录及子路由刷新。
