# GoFly Admin - AI Agent 项目指南

> 本文档面向 AI Coding Agent，提供项目架构、开发规范和常用操作的快速参考。
> 项目语言：中文（注释和文档主要使用中文）

---

## 一、项目概述

**GoFly Admin** 是一个基于 Go + Vue3 的全栈中后台管理系统（企业版 V3），主要用于房产/门店管理等业务场景。采用前后端分离架构，支持管理端（Admin）、业务端（Business）和微信小程序端（WxApp）三端。

### 1.1 技术栈

| 层级 | 技术 | 版本/说明 |
|------|------|-----------|
| 后端 | Go | 1.25+ |
| Web 框架 | Gin | v1.9.1 |
| 数据库 | MySQL | 5.7+ (utf8mb4) |
| 缓存 | Redis | 可选 |
| ORM | 自定义 gform | 位于 `utils/gform/` |
| 管理端前端 | Vue 3 + Vite + TS | Arco Design UI |
| 业务端前端 | Vue 3 + Vite + TS | Arco Design UI |
| 小程序 | uni-app | Vue 3 语法 |

### 1.2 项目目录结构

```
main-admin/
├── app/                      # Go 后端业务代码
│   ├── admin/               # 管理端模块（端口 8200，路径 /admin/*）
│   │   ├── business/        # 业务账号管理
│   │   ├── system/          # 系统管理（角色、权限、部门）
│   │   ├── user/            # 用户管理
│   │   └── controller.go    # 管理端路由钩子（RBAC权限验证）
│   ├── business/            # 业务端模块（端口 8200，路径 /business/*）
│   │   ├── houses/          # 房源管理（核心业务）
│   │   ├── system/          # 系统管理
│   │   ├── uniapp/          # 小程序 API
│   │   └── controller.go    # 业务端路由钩子
│   ├── common/              # 公共模块（无需登录）
│   │   ├── install/         # 系统安装
│   │   └── upload/          # 文件上传
│   └── controller.go        # 总路由注册入口
├── web/
│   ├── admin/               # 管理端前端（端口 9201）
│   ├── business/            # 业务端前端（端口 9200）
│   └── wxapp/               # 小程序源码
├── resource/                # 资源目录
│   ├── config.yaml          # 后端核心配置文件 ⚠️
│   ├── uploads/             # 上传文件存储
│   ├── webadmin/            # Admin 前端构建产物
│   └── webbusiness/         # Business 前端构建产物
├── utils/                   # Go 工具库
│   ├── gf/                  # GoFly 核心框架（路由、响应封装）
│   ├── gform/               # ORM 数据库操作
│   ├── router/              # HTTP 服务器和中间件
│   ├── drivers/             # 数据库驱动（MySQL/Redis）
│   └── tools/               # 各类工具包（缓存、日志、验证等）
├── devsource/               # 开发工具
│   ├── developer/install/   # 系统安装页面
│   └── codemarket/          # 插件市场
├── main.go                  # Go 入口文件
└── schema.sql               # 数据库初始化 SQL
```

---

## 二、关键配置文件

### 2.1 后端配置 `resource/config.yaml`

```yaml
database:
  default:
    hostname: 127.0.0.1
    hostport: 3306
    username: root
    password: your_password      # 必改
    dbname: gofly_bs             # 必改
    type: mysql
    charset: utf8mb4

app:
  port: 8200                    # 后端服务端口
  apisecret: gofly@888          # API 验证密钥（与前端 .env VITE_ENCRYPT 必须一致）
  tokensecret: gf849841325189456f489  # JWT 密钥
  allowurl: http://localhost:9200,http://localhost:9201  # 跨域白名单（必改）
  vueobjroot: E:/Work/Progect/GoFly-admin/main-admin/web  # 前端代码绝对路径（必改）
  busDirName: business          # 业务端目录名
  adminDirName: admin           # 管理端目录名
  runEnv: debug                 # debug/release/test
  loginCaptcha: true            # 登录验证码
  validityApi: true             # 开启 API 合法性验证
```

### 2.2 前端配置

**管理端** (`web/admin/public/config.js`):
```javascript
const localhost = "http://localhost:8200";
window.globalConfig = {
    Root_url_dev: `${localhost}/admin`,  // 开发环境 API
};
```

**管理端环境变量** (`web/admin/.env`):
```env
VITE_APP_ENV=development
VITE_ENCRYPT=gofly@888           # 必须与后端 apisecret 一致
VITE_GLOB_APP_TITLE=GoFly管理后台
```

---

## 三、开发环境搭建

### 3.1 依赖安装

```bash
# 1. 后端依赖
go mod tidy

# 2. 热重载工具（二选一）
# Windows 推荐:
go install github.com/pilu/fresh@latest
# Mac 推荐:
go install github.com/cosmtrek/air@latest

# 3. 管理端前端
cd web/admin
yarn install

# 4. 业务端前端
cd web/business
yarn install
```

### 3.2 启动服务

```bash
# 终端 1: 启动后端
cd project-root
fresh   # 或 air
# 首次访问 http://localhost:8200/install 进行系统安装

# 终端 2: 启动管理端前端
cd web/admin
yarn serve
# 访问 http://localhost:9201

# 终端 3: 启动业务端前端
cd web/business
yarn serve
# 访问 http://localhost:9200
```

### 3.3 热重载配置

- **fresh**: 配置见 `runner.conf`
- **air**: 配置见 `.air.conf`

**⚠️ 注意**: 如果热重载太慢，将 `web/` 目录移到项目外，并修改 `config.yaml` 中的 `vueobjroot` 为新路径。

---

## 四、代码开发规范

### 4.1 后端开发

#### 控制器结构

```go
package houses

import (
    "gofly/utils/gf"
    "gofly/utils/gform"
)

// Houses 控制器结构体
type Houses struct {
    NoNeedLogin []string  // 无需登录的方法列表，["*"]表示全部
    NoNeedAuths []string  // 无需权限验证的方法列表
}

// init 注册控制器
func init() {
    fpath := Houses{NoNeedAuths: []string{"GetList"}}  // GetList 无需权限
    gf.Register(&fpath, fpath)
}

// GetList 获取列表（自动映射为 GET /business/houses/getList）
func (api *Houses) GetList(c *gf.GinCtx) {
    // 1. 获取参数
    pageNo := gf.GetQueryInt(c, "page", 1)
    pageSize := gf.GetQueryInt(c, "pageSize", 12)
    
    // 2. 查询数据
    db := gform.DB().Table("business_properties")
    // ... 构建查询条件
    
    // 3. 响应
    gf.Success().SetData(list).SetCount(total).Regin(c)
}

// Add 添加（自动映射为 POST /business/houses/add）
func (api *Houses) Add(c *gf.GinCtx) {
    param, _ := gf.RequestParam(c)  // 获取请求参数
    
    // 插入数据
    result, err := gform.DB().Table("business_properties").Data(param).Insert()
    if err != nil {
        gf.Failed().SetMsg("添加失败").Regin(c)
        return
    }
    gf.Success().SetMsg("添加成功").SetData(result).Regin(c)
}
```

#### 路由自动生成规则

路由格式: `/{module}/{controller}/{action}`

| HTTP 方法 | 方法名前缀 | 示例方法名 | 生成路由 |
|-----------|-----------|-----------|---------|
| GET | `Get` | `GetList` | `GET /business/houses/getList` |
| POST | 无/Post | `Add` | `POST /business/houses/add` |
| PUT | `Put` | `PutUpdate` | `PUT /business/houses/putUpdate` |
| DELETE | `Del` | `DelById` | `DELETE /business/houses/delById` |
| GET+POST | `GetPost` | `GetPostData` | 同时注册 GET 和 POST |

路径参数: 方法名以 `ParId` 结尾时，自动添加 `/:id` 路径参数。

#### 数据库操作 (gform ORM)

```go
// 查询单条
record, err := gform.DB().Table("business_user").Where("id", 1).One()

// 查询列表
list, err := gform.DB().Table("business_properties").
    Where("status", 0).
    WhereLike("title", "%keyword%").
    Order("weigh desc, id desc").
    Limit(10).
    All()

// 分页查询
total, err := db.Count()
list, err := db.Page(pageNo, pageSize).All()

// 插入
result, err := gform.DB().Table("business_stores").Data(gf.Map{
    "name": "门店名称",
    "address": "地址",
}).Insert()

// 更新
result, err := gform.DB().Table("business_properties").
    Where("id", 1).
    Data(gf.Map{"status": 1}).
    Update()

// 删除（软删除）
result, err := gform.DB().Table("business_properties").Where("id", 1).Delete()

// 事务
err := gform.DB().Transaction(func(tx *gform.DB) error {
    tx.Table("table1").Data(...).Insert()
    tx.Table("table2").Where(...).Update()
    return nil  // 返回 nil 提交，返回 error 回滚
})
```

#### 响应封装

```go
// 成功响应
gf.Success().SetData(data).Regin(c)
gf.Success().SetMsg("操作成功").SetData(data).SetCount(total).Regin(c)

// 失败响应
gf.Failed().SetMsg("操作失败").Regin(c)
gf.Failed().SetCode(401).SetMsg("未登录").SetExdata(errMsg).Regin(c)
```

### 4.2 前端开发

#### API 请求封装

位置: `web/admin/src/utils/http/index.ts` (管理端)

```typescript
import axios from 'axios';

// 自动处理 apiverify 和 Authorization 头
// apiverify = Base64(MD5(VITE_ENCRYPT + timestamp) + "#" + timestamp)
// 时间戳 5 分钟内有效
```

#### 权限控制

- 路由权限: 动态加载服务端菜单配置
- 按钮权限: 通过 `to.meta.btnroles` 获取

---

## 五、构建与部署

### 5.1 后端构建

```bash
# Windows 本地运行
go build main.go

# Windows 构建 Linux 可执行文件（推荐用于服务器部署）
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build main.go

# Mac 构建 Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
```

### 5.2 前端构建

```bash
# 管理端构建（输出到 resource/webadmin/）
cd web/admin
yarn build

# 业务端构建（输出到 resource/webbusiness/）
cd web/business
yarn build
```

### 5.3 部署检查清单

- [ ] 修改 `config.yaml` 中的 `runEnv: release`
- [ ] 修改 `apisecret` 和 `tokensecret` 为强密码
- [ ] 修改 `allowurl` 为生产域名
- [ ] 修改数据库配置为生产环境
- [ ] 确保前端 `.env` 中 `VITE_ENCRYPT` 与后端 `apisecret` 一致
- [ ] 构建前端并确认资源在 `resource/webadmin/` 和 `resource/webbusiness/`

---

## 六、安全机制

### 6.1 API 合法性验证

请求头: `apiverify: Base64(MD5(secret + timestamp) + "#" + timestamp)`

- secret: 后端 `config.yaml` 中的 `apisecret`
- timestamp: Unix 时间戳（秒）
- 有效期: 180 秒（由 `apiouttime` 配置）

### 6.2 JWT Token 验证

请求头: `Authorization: <token>`

- 登录时获取 token
- 超时时间: 120 分钟（由 `tokenouttime` 配置）
- 密钥: `tokensecret`

### 6.3 RBAC 权限

- 用户 → 角色 → 权限（菜单/按钮）
- 超级管理员角色 ID: 在代码中硬编码检查
- 权限验证位置: `app/admin/controller.go` 和 `app/business/controller.go`

---

## 七、常见问题

| 问题 | 原因 | 解决 |
|------|------|------|
| 前端请求 403 | 跨域配置不正确 | 检查 `allowurl` 是否包含前端地址 |
| "请求不合法" | apiverify 验证失败 | 确认 `VITE_ENCRYPT` 与 `apisecret` 一致 |
| 热加载太慢 | web/ 目录被监听 | 将前端移到项目外，修改 `vueobjroot` |
| 数据库连接失败 | 配置错误或权限问题 | 检查 `config.yaml` 数据库配置 |
| 路由 404 | 路由未注册 | 查看 `runtime/app/routers.txt` 确认路由已生成 |
| 登录后跳转失败 | token 验证失败 | 检查 `tokensecret` 配置和 Redis 连接 |

---

## 八、参考资源

- GoFly 官方文档: https://doc.goflys.cn
- Gin 框架文档: https://gin-gonic.com/docs/
- Vue 3 文档: https://vuejs.org/
- Arco Design: https://arco.design/vue/docs/start
- uni-app 文档: https://uniapp.dcloud.net.cn/

---

**文档版本**: 2026-02-02  
**最后更新**: 基于 GoFly V3.0.13
