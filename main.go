package main

import (
	//引入数据库驱动-去这里下载：https://doc.goflys.cn/docview?id=26&fid=395
	// Redis驱动和安装说明：https://doc.goflys.cn/docview?id=26&fid=392
	_ "gofly/utils/drivers/mysql"
	// _ "gofly/utils/drivers/redis"
	"gofly/utils/router"
)

func main() {
	// 启动服务器
	router.RunServer()
}
