# 管理后台 Timeless 迁移

迁移目标是本仓库 `@fithub/admin`，不是 `timeless/apps/web-solid`。

## 运行

本地框架来源默认是 `../../timeless/packages`（即 `/Users/litao/Documents/workspace/timeless/packages`）。使用该项目构建好的 core、DOM renderer 和 shadcn ESM 产物。业务 API 默认为 `http://127.0.0.1:8080`，开发时由 Vite 代理 `/api`；可用 `ADMIN_API_URL` 覆盖。

```sh
pnpm dev
pnpm check
pnpm test
pnpm build
```

## 实现边界

- 保留当前路由、业务模型和接口，包括迁移开始前已有的未提交业务修改。
- TSX 由 `scripts/timeless-jsx.ts` 编译为 Timeless VNode，运行时使用 Timeless ref、View、For、Show 和 DOM renderer；不调用 Solid 渲染器。
- `src/timeless/index.ts` 保留旧页面使用的响应式函数名称，减少业务层变更。SolidJS 仅作为开发期 JSX 类型来源，Lucide 仅提供已复制的 ISC SVG 数据；生产构建会检查两者没有运行时模块进入产物。
- 共享 Button、Input、Textarea、Checkbox、Dialog、DropdownMenu、ContextMenu、Popover、DatePicker、Tabs、Progress、Label、Skeleton，以及 Card、Badge、FieldGroup/FieldSet/FieldLegend/FieldDescription 使用本地 shadcn。适配层连接原业务 store 与 Timeless vm，保留提交、取消、加载、禁用、选项更新和日期值同步。
- Select 使用 Timeless SelectPrimitive 与 shadcn 语义样式组合：本地库 0.34.2 的 Select 下拉内容仍为 `Som` 占位，直接调用会缺少选项。项目适配支持完整选项、键盘操作、动态选项和数值类型。SimpleSelect 复用该实现。
- Toast 保留原有定时关闭行为并采用语义配色：库 0.34.2 的 Toast 可视提示内容尚未实现。文件输入保留 FileList，媒体上传、拖拽、滚动与路由等业务组件继续通过 Timeless 渲染层运行。
- 登录/注册页、侧栏导航、首页、肌肉与器械卡片采用共享组件。肌肉、器械、动作、训练计划和计划集合的动态表单已复用共享输入、选择与复选控件。
- 直接加载本地 `@timeless/shadcn/globals.css`，采用库的亮/暗主题、`--primary`、`--background`、`--foreground`、`--border`、`--input`、`--accent` 和 `--radius`。旧 `--weui-*` 名称仅作为业务视图到语义 token 的别名，HTML 中旧色板已移除。Tailwind 3 保留用于现有页面布局，库的 `tt-*` 类使用预编译 CSS。
- 没有发现本项目 `.design` 合同；视觉规范以指定组件库的 `packages/shadcn/THEME_DESIGN.md` 为依据。未另建设计系统。

## 验证说明

渲染层测试覆盖文本/属性更新、列表删除、条件显示、ref 与挂载顺序、清理、布尔属性和条件列表根节点。真实 API 回归应使用独立数据库，不能将构建通过等同于业务功能完全一致。

非浏览器检查：`pnpm check`、`pnpm build` 通过；生产模块图没有 SolidJS / lucide-solid 运行时。新增组件 DOM 测试覆盖输入/多行文本双向更新、密码类型、按钮加载、复选框禁用、数值选择与动态选项、弹窗提交/取消/重新打开、动态菜单、导航样式与属性以及日期值同步。

完整测试结果为 23 项通过、1 项失败；新增的 10 项组件测试全部通过。失败项在迁移前已存在：`src/domains/ui/formv2/__tests__/index.test.ts:74` 将对象类型的 `form.value` 断言为字符串。未修改该业务断言，也未将其计为新组件测试通过。

遵照 AGENTS.md，没有启动或自动操作浏览器。请手动验证登录/注册/退出、侧栏选中状态、亮暗主题、列表筛选、创建/编辑表单、弹窗关闭、日期与多选以及媒体上传。真实 API 与浏览器视觉行为仍需手动确认，构建与 DOM 测试不代表已完成真实业务回归。
