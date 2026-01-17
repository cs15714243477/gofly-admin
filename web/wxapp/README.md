# WxApp (uni-app) 项目配置说明

## 一、环境配置

### 1. 创建环境变量文件

项目根目录下已自动创建 `.env` 文件，包含以下配置：

```env
# 开发环境 API 地址
GF_DEV_BASE_URL=http://localhost:8200
# 开发服务器端口
GF_DEV_PORT=9000
# API 路径前缀
GF_API_PATH=/uniapp
# API 模块
GF_API_MODEL=/uniapp
# API 验证密钥（必须与后端 resource/config.yaml 中的 apisecret 一致）
GF_API_SECRET=gofly@888
# 业务 ID
GF_BUSINESUSSID=
# 版本号
GF_VERSION=1.0.0
# 静态资源地址
GF_STATIC_URL=
```

**⚠️ 重要**：
- `GF_API_SECRET` 必须与后端 `resource/config.yaml` 中的 `apisecret` 保持一致
- `GF_DEV_BASE_URL` 应指向后端服务地址（默认 `http://localhost:8200`）

### 2. 安装依赖

```bash
cd web/wxapp
npm install
# 或
yarn install
```

## 二、解决编译超时问题

### 问题现象
```
timeout(env: Windows,mp,1.06.2504060; lib: 3.13.0)
```

### 解决方案

#### 1. 检查环境变量
确保 `.env` 文件存在且配置正确，特别是 `GF_DEV_PORT` 必须有值。

#### 2. 清理构建缓存
```bash
# 删除构建产物
rm -rf unpackage/dist
# Windows PowerShell
Remove-Item -Recurse -Force unpackage/dist
```

#### 3. 优化编译配置
已优化 `vite.config.js`：
- 添加默认端口（9000）
- 优化依赖预构建
- 优化构建性能

#### 4. 微信开发者工具设置
1. 打开微信开发者工具
2. 设置 → 编辑器设置 → 编译设置
3. 增加编译超时时间（如果有此选项）
4. 关闭"上传代码时样式自动补全"（可加快编译速度）

#### 5. 网络问题
如果编译时出现网络超时：
- 检查网络连接
- 尝试使用代理
- 关闭防火墙或杀毒软件（临时）

#### 6. 重新安装依赖
```bash
# 删除 node_modules 和锁文件
rm -rf node_modules package-lock.json
# 重新安装
npm install
```

## 三、开发调试

### 1. H5 开发
```bash
npm run dev:h5
# 或使用 HBuilderX 运行到浏览器
```

### 2. 微信小程序开发
```bash
# 使用 HBuilderX：
# 1. 打开项目
# 2. 运行 → 运行到小程序模拟器 → 微信开发者工具

# 或使用命令行（需要配置微信开发者工具路径）
npm run dev:mp-weixin
```

### 3. 配置微信开发者工具
1. 打开微信开发者工具
2. 设置 → 安全 → 开启"服务端口"
3. 在 HBuilderX 中配置微信开发者工具路径

## 四、API 请求配置

### 请求封装位置
- 文件：`utils/request/index.js`
- 自动添加 `apiverify` 头（使用 `GF_API_SECRET`）
- 自动添加 `Authorization` 头（JWT Token）
- 自动处理 token 刷新

### 使用示例
```javascript
import request from '@/utils/request'

// GET 请求
request({
  url: '/user/info',
  method: 'GET'
})

// POST 请求
request({
  url: '/user/login',
  method: 'POST',
  data: {
    username: 'admin',
    password: '123456'
  }
})
```

## 五、常见问题

### 1. 编译超时
- **原因**：编译时间过长或网络问题
- **解决**：参考"解决编译超时问题"章节

### 2. API 请求失败
- **检查**：`GF_API_SECRET` 是否与后端一致
- **检查**：`GF_DEV_BASE_URL` 是否正确
- **检查**：后端服务是否正常运行

### 3. 环境变量未生效
- **检查**：`.env` 文件是否在项目根目录
- **检查**：环境变量前缀是否为 `GF_`
- **重启**：重启开发服务器

### 4. 微信开发者工具无法连接
- **检查**：微信开发者工具是否开启服务端口
- **检查**：防火墙是否阻止连接
- **尝试**：手动在微信开发者工具中打开项目

## 六、项目结构

```
wxapp/
├── api/              # API 接口定义
├── components/       # 组件
├── pages/            # 页面
├── static/           # 静态资源
├── store/           # Pinia 状态管理
├── utils/           # 工具函数
│   ├── config/      # 配置（环境变量）
│   ├── request/     # 请求封装
│   └── router/      # 路由工具
├── App.vue          # 应用入口
├── main.js          # 主入口文件
├── manifest.json    # 应用配置
├── pages.json       # 页面配置
└── vite.config.js   # Vite 配置
```

## 七、构建发布

### 1. 构建 H5
```bash
npm run build:h5
# 产物输出到：unpackage/dist/build/h5
```

### 2. 构建微信小程序
```bash
npm run build:mp-weixin
# 产物输出到：unpackage/dist/build/mp-weixin
# 使用微信开发者工具上传代码
```

### 3. 构建 App
```bash
npm run build:app-plus
# 使用 HBuilderX 打包 App
```

## 八、参考文档

- uni-app 官方文档：https://uniapp.dcloud.net.cn/
- Vite 配置文档：https://vitejs.dev/config/
- 微信小程序文档：https://developers.weixin.qq.com/miniprogram/dev/framework/

---

**最后更新**：2026-01-16
