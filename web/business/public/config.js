const domain="";//您的域名//独立部署需要填写Go服务api接口域名(部署在Go目录下可以留空)
const localhost="http://localhost:8200";//访问本地的域名和端口，如果您改变Go服务端口，请自行修改
window.globalConfig = {
	Main_url:`${domain}`,//域名
	Main_url_dev:`${localhost}`,//域名
	Root_url:`${domain}/business`,//Api服务器域名
    Root_url_dev:`${localhost}/business`,//Api服务器域名-开发环境
    TencentMapKey:"YBUBZ-EXEKL-ELNPT-EBIW6-7RADV-ODBBF",//腾讯地图 JS API Key（用于房源经纬度地图选点）
	AppTitle_zhCN:"GoFly开发框架(企业版)",
    AppTitle_zhTW:"GoFly開發框架",
	AppTitle_enUS:"GoFly Dev frame",
    AppVersion:"2.5.0",
    CompanySite:"//goflys.cn/",//公司官网
    ICP:"滇ICP备19001128号-3",//备案号
    Company:"正科云滇科技有限公司",//公司名称
    Address:"中国·昆明",
    
    TeamSite:"//goflys.cn/",//技术团队官网
    Team:"GoFly技术团队",//技术团队，没有则填空""
    loginTitle:"登录GoFly后台管理系统",
    loginSubTitle:"管理账号、业务、数据统计、触达事件等",
    Copyright:"GoFly团队提供技术支持",
    RouterHome:"home",//路由默认入口
    MaxSizeImage:5,//最大上传图片大小,单位M
	MaxSizeVideo:150,//最大上传视频大小,单位M
};
