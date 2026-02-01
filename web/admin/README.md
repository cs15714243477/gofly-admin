##  GoFly单服务版 -管理前端代码
 你开发的业务代码在这里编写，这里的超级管理账号admin端“业务端管理”里的账号添加和编辑。
### 快速运行项目
建议使用 yarn 
```
# 初始化项目
```
npm install 或者 yarn install

# 运行项目
npm run serve 或者  yarn serve

# 打包项目
npm run build 或者 yarn build 

```
如果 install 失败请检查镜像
可以切换淘宝镜像：
```
yarn config set registry http://registry.npm.taobao.org/
```
打包注意事项：
1.如果您是独立域名部署，打包时候需要修改两个地方：
（1）、config/vite.config.prod.ts中的mergeConfig下的：base:""。
（2）、src/router/index.ts中的createRouter下的： history: createWebHashHistory('') 。
2.如果您是部署在后端代码二级目录下，同样打包时候需要修改两个地方（例如业务端webadmin，按照您的需求替换“webadmin”）：
（1）、config/vite.config.prod.ts中的mergeConfig下的：base: process.env.NODE_ENV === 'production' ? '/webadmin/' : ''。
（2）、src/router/index.ts中的createRouter下的：history: createWebHashHistory(process.env.NODE_ENV === 'production' ? '/webadmin/' : '')。

### 开发说明
您在开发时候请您参考文档说明 https://doc.goflys.cn/docview?id=25
 

 