# GoFly 项目开发指南

## 一、项目架构

### 技术栈
- **后端**：Go 1.25+ (Gin 框架)
- **前端 Admin**：Vue 3 + Vite + TypeScript + Arco Design
- **前端 Business**：Vue 3 + Vite + TypeScript + Arco Design
- **前端 WxApp**：uni-app (Vue 3)

### 项目结构
```
main-admin/
├── app/                    # Go 后端业务代码
│   ├── admin/             # 管理端业务模块
│   ├── business/          # 业务端业务模块
│   └── common/            # 公共模块
├── web/
│   ├── admin/             # 管理端前端 (端口: 9201)
│   ├── business/          # 业务端前端 (端口: 9200)
│   └── wxapp/             # 小程序前端 (uni-app)
├── resource/
│   ├── config.yaml        # 后端配置文件（重要！）
│   ├── webadmin/          # Admin 前端构建产物部署目录
│   └── webbusiness/       # Business 前端构建产物部署目录
├── utils/                  # Go 工具库
├── main.go                # Go 入口文件
└── schema.sql             # 数据库初始化 SQL
```

---

## 二、环境准备

### 1. 后端环境
- Go 1.25+ 
- MySQL 5.7+
- Redis（可选，用于缓存）

### 2. 前端环境
- Node.js 14.0+
- npm/yarn/pnpm（推荐 yarn）

### 3. 热加载工具（二选一）
```bash
# 方式1：fresh（Windows 推荐）
go install github.com/pilu/fresh@latest

# 方式2：air（Mac 推荐）
go install github.com/cosmtrek/air@latest
```

---

## 三、初始化步骤

### 步骤1：配置数据库

编辑 `resource/config.yaml`，修改数据库连接信息：

```yaml
database:
  default:
    hostname: 127.0.0.1
    hostport: 3306
    username: root
    password: your_password      # 修改为你的密码
    dbname: gofly_bs             # 修改为你的数据库名
    type: mysql
    charset: utf8mb4
```

### 步骤2：初始化数据库

```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE gofly_bs CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

# 导入表结构（可选，系统首次运行会自动安装）
mysql -u root -p gofly_bs < schema.sql
```

### 步骤3：配置后端

编辑 `resource/config.yaml`，关键配置：

```yaml
app:
  port: 8200                    # 后端服务端口
  apisecret: gofly@888          # API 验证密钥（重要！）
  tokensecret: gf849841325189456f489  # JWT 密钥
  allowurl: http://localhost:9200,http://localhost:9201  # 跨域白名单
  vueobjroot: E:/Work/Progect/GoFly-admin/main-admin/web  # 前端代码根目录（必改！）
  busDirName: business          # 业务端目录名
  adminDirName: admin            # 管理端目录名
  runEnv: debug                 # 开发环境：debug，生产：release
```

**⚠️ 重要**：`vueobjroot` 必须修改为你的实际前端代码路径（绝对路径）。

### 步骤4：初始化 Go 依赖

```bash
go mod tidy
```

### 步骤5：配置前端 Admin

编辑 `web/admin/public/config.js`：

```javascript
const domain = "";  // 生产环境域名（开发环境留空）
const localhost = "http://localhost:8200";  // 后端服务地址
window.globalConfig = {
    Root_url_dev: `${localhost}/admin`,  // 开发环境 API 地址
    // ... 其他配置
};
```

创建 `web/admin/.env`（如果不存在）：

```env
VITE_APP_ENV=development
VITE_ENCRYPT=gofly@888
VITE_GLOB_APP_TITLE=GoFly管理后台
```

**⚠️ 注意**：`VITE_ENCRYPT` 必须与后端 `config.yaml` 中的 `apisecret` 保持一致！

### 步骤6：配置前端 Business

编辑 `web/business/public/config.js`：

```javascript
const domain = "";
const localhost = "http://localhost:8200";
window.globalConfig = {
    Root_url_dev: `${localhost}/business`,  // 开发环境 API 地址
    // ... 其他配置
};
```

创建 `web/business/.env`（如果不存在）：

```env
VITE_APP_ENV=development
VITE_ENCRYPT=gofly@888
VITE_GLOB_APP_TITLE=GoFly业务端
```

### 步骤7：安装前端依赖

```bash
# Admin 端
cd web/admin
yarn install

# Business 端
cd web/business
yarn install

# WxApp 端（可选）
cd web/wxapp
npm install
```

### 步骤8：启动后端服务

```bash
# 方式1：使用 fresh（推荐）
fresh

# 方式2：使用 air
air

# 方式3：直接运行（无热加载）
go run main.go
```

后端启动后访问：`http://localhost:8200/install` 进行系统安装（首次运行）。

### 步骤9：启动前端服务

```bash
# 启动 Admin 端（新终端）
cd web/admin
yarn serve
# 访问：http://localhost:9201

# 启动 Business 端（新终端）
cd web/business
yarn serve
# 访问：http://localhost:9200
```

---

## 四、开发要点

### 1. API 请求机制

#### 后端验证流程
1. **接口合法性验证**（`apiverify` 头）：
   - 前端：`MD5(VITE_ENCRYPT + timestamp) + "#" + timestamp` → Base64
   - 后端：验证 MD5 和时间戳（5分钟内有效）
   - 配置：`resource/config.yaml` 中的 `apisecret` 必须与前端 `VITE_ENCRYPT` 一致

2. **JWT Token 验证**：
   - 请求头：`Authorization: <token>`
   - 超时时间：`config.yaml` 中的 `tokenouttime`（默认 120 分钟）

#### 前端请求封装
- **Admin/Business**：`src/utils/http/index.ts`
  - 自动添加 `apiverify` 头
  - 自动添加 `Authorization` 头
  - 自动处理 token 刷新
  - API 地址：通过 `window.globalConfig.Root_url_dev` 获取

- **WxApp**：`utils/request/index.js`
  - 使用 `GF_API_SECRET` 环境变量生成 `apiverify`
  - 使用 `GF_DEV_BASE_URL` 配置 API 地址

### 2. 路由与权限

#### 前端路由守卫
- **登录验证**：`src/router/guard/userLoginInfo.ts`
  - 检查 token 是否存在
  - 自动获取用户信息

- **权限验证**：`src/router/guard/permission.ts`
  - 动态加载服务端菜单配置
  - 验证路由访问权限
  - 按钮权限：通过 `to.meta.btnroles` 获取

#### 后端路由注册
- 路由自动注册：在 `app/*/controller.go` 中定义
- 路由文件：`runtime/app/routers.txt`（自动生成）

### 3. 跨域配置

后端 `config.yaml`：
```yaml
app:
  allowurl: http://localhost:9200,http://localhost:9201
```

**⚠️ 注意**：如果前端端口改变，必须更新 `allowurl` 配置。

### 4. 代码生成

系统支持代码生成功能：
- 位置：`app/business/developer/generatecode/`
- 可生成：Go 后端代码、Vue 前端代码、数据库表结构

### 5. 插件安装

系统支持插件化安装：
- 插件目录：`devsource/codemarket/install/`
- 安装入口：后台管理 → 开发者 → 代码商店
- 安装流程：SQL 导入 → 菜单导入 → 代码复制

### 6. 日志与调试

- **后端日志**：`runtime/log/`
- **SQL 日志**：`runtime/log/sql/`（需在 `config.yaml` 中开启）
- **前端调试**：浏览器 DevTools
- **性能分析**：`config.yaml` 中开启 `runpprof: true`，访问 `http://localhost:8081/debug/pprof/`

### 7. 环境变量

#### Admin/Business 前端
- `VITE_APP_ENV`：环境标识（development/production）
- `VITE_ENCRYPT`：API 验证密钥（必须与后端 `apisecret` 一致）
- `VITE_GLOB_APP_TITLE`：应用标题

#### WxApp 前端
- `GF_DEV_BASE_URL`：开发环境 API 地址
- `GF_BASE_URL`：生产环境 API 地址
- `GF_API_SECRET`：API 验证密钥
- `GF_DEV_PORT`：开发服务器端口
- `GF_API_PATH`：API 路径前缀
- `GF_API_MODEL`：API 模块（如 `/uniapp`）
- `GF_BUSINESUSSID`：业务 ID

### 8. 构建与部署

#### 前端构建
```bash
# Admin 端
cd web/admin
yarn build
# 产物输出到：resource/webadmin/

# Business 端
cd web/business
yarn build
# 产物输出到：resource/webbusiness/

# WxApp 端
cd web/wxapp
npm run build:mp-weixin  # 微信小程序
npm run build:h5         # H5
# 产物输出到：unpackage/dist/build/
```

**⚠️ WxApp 编译超时问题**：
- 确保 `.env` 文件存在且配置正确
- 参考 `web/wxapp/README.md` 中的"解决编译超时问题"章节

#### 后端构建
```bash
# Windows 构建 Linux 版本（推荐）
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build main.go

# Windows 构建 Windows 版本
SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64
go build main.go
```

---

## 五、常见问题

### 1. 前端请求 403 错误
- **原因**：跨域配置不正确
- **解决**：检查 `resource/config.yaml` 中的 `allowurl` 是否包含前端地址

### 2. 前端请求 "请求不合法"
- **原因**：`apiverify` 验证失败
- **解决**：确保前端 `VITE_ENCRYPT` 与后端 `apisecret` 一致

### 3. 热加载太慢
- **原因**：`web/` 目录被监听，文件变化触发重建
- **解决**：将前端代码移到项目外，修改 `config.yaml` 中的 `vueobjroot`

### 4. 数据库连接失败
- **检查**：`config.yaml` 中的数据库配置
- **检查**：MySQL 服务是否启动
- **检查**：数据库用户权限

### 5. 前端页面空白
- **检查**：`public/config.js` 是否正确加载
- **检查**：浏览器控制台是否有错误
- **检查**：后端服务是否正常运行

### 6. 路由 404
- **检查**：后端路由是否注册（查看 `runtime/app/routers.txt`）
- **检查**：前端路由配置是否正确
- **检查**：权限配置是否正确

---

## 六、开发规范

### 1. 代码风格
- **Go**：遵循 Go 官方代码规范
- **Vue**：遵循 Vue 3 官方风格指南
- **TypeScript**：使用严格模式

### 2. 提交规范
- 参考：`cursor-rules/other/git.md`
- 使用 Gitflow 工作流（参考：`cursor-rules/other/gitflow.md`）

### 3. 命名规范
- **Go**：驼峰命名，公开函数首字母大写
- **Vue**：组件使用 PascalCase，文件使用 kebab-case
- **数据库**：表名使用下划线命名，字段名使用下划线命名

### 4. 目录结构
- **后端**：按业务模块划分（`app/admin/`, `app/business/`）
- **前端**：按功能模块划分（`src/views/`, `src/components/`）

---

## 七、快速启动检查清单

- [ ] 数据库已创建并配置
- [ ] `resource/config.yaml` 已正确配置（数据库、端口、跨域、vueobjroot）
- [ ] 后端依赖已安装（`go mod tidy`）
- [ ] 后端服务已启动（`fresh` 或 `air`）
- [ ] 系统已安装（访问 `http://localhost:8200/install`）
- [ ] 前端 `public/config.js` 已配置
- [ ] 前端 `.env` 已创建（`VITE_ENCRYPT` 与后端一致）
- [ ] 前端依赖已安装（`yarn install`）
- [ ] 前端服务已启动（`yarn serve`）
- [ ] 浏览器可正常访问前端页面
- [ ] 登录功能正常

---

## 八、参考文档

- GoFly 官方文档：https://doc.goflys.cn
- Gin 框架文档：https://gin-gonic.com/docs/
- Vue 3 文档：https://vuejs.org/
- Arco Design 文档：https://arco.design/vue/docs/start
- uni-app 文档：https://uniapp.dcloud.net.cn/

---

**最后更新**：2026-01-16
