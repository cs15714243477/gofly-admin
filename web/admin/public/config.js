const domain="";//您的域名//独立部署需要填写Go服务api接口域名(部署在Go目录下可以留空)
const localhost="http://localhost:8200";//访问本地的域名和端口，如果您改变Go服务端口，请自行修改
window.globalConfig = {
    Domain:domain,
    Main_url:`${domain}/`,//域名
    Main_url_dev:`${localhost}/`,//域名
	Root_url:`${domain}/admin`,//Api服务器域名
	Root_url_dev:`${localhost}/admin`,//Api服务器域名-开发环境
	AppTitle:"GoFly管理SaaS系统后台(企业版)",
    AppVersion:"2.0.0",
    CompanySite:"//goflys.cn/",//公司官网
    Company:"正科云滇科技有限公司",//公司名称
    Address:"中国·昆明",
    TeamSite:"//goflys.cn/",//技术团队官网
    Team:"GoFly技术团队",//技术团队，没有则填空""
    loginTitle:"登录GoFly快速后台管理系统",
    loginSubTitle:"管理账号、业务、数据统计、触达事件等",
    Copyright:"GoFly团队提供技术支持",
    RouterHome:"home",//路由默认入口
    MaxSizeImage:5,//最大上传图片大小,单位M
	MaxSizeVideo:150,//最大上传视频大小,单位M
};