package uniapp

import (
	"strings"

	"gofly/utils/gf"
	"gofly/utils/tools/gconv"
)

// 省市区（小程序端）
type Area struct {
	NoNeedLogin []string
	NoNeedAuths []string
}

func init() {
	fpath := Area{NoNeedLogin: []string{}, NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

// 获取省份列表（pid=0）
func (api *Area) GetProvinces(c *gf.GinCtx) {
	// businessID：优先从登录态，其次 header，最后默认 1
	businessID := c.GetInt64("businessID")
	if businessID == 0 {
		businessID = gf.Int64(c.GetHeader("Businessid"))
		if businessID == 0 {
			businessID = 1
		}
	}

	MDB := gf.Model("business_area").Where("pid", 0)
	if gf.DbHaseField("business_area", "business_id") {
		MDB = MDB.Where("business_id", businessID)
	}

	list, err := MDB.Fields("id,name,level").Order("id asc").Select()
	if err != nil {
		gf.Failed().SetMsg("获取省份失败：" + err.Error()).Regin(c)
		return
	}
	gf.Success().SetMsg("获取省份").SetData(gf.Map{"items": list}).Regin(c)
}

// 获取子级数据（传 pid）
func (api *Area) GetChildren(c *gf.GinCtx) {
	// businessID：优先从登录态，其次 header，最后默认 1
	businessID := c.GetInt64("businessID")
	if businessID == 0 {
		businessID = gf.Int64(c.GetHeader("Businessid"))
		if businessID == 0 {
			businessID = 1
		}
	}

	pid := strings.TrimSpace(c.Query("pid"))
	if pid == "" {
		param, _ := gf.RequestParam(c)
		pid = strings.TrimSpace(gconv.String(param["pid"]))
	}
	pidInt := gconv.Int64(pid)
	if pidInt < 0 {
		pidInt = 0
	}

	MDB := gf.Model("business_area").Where("pid", pidInt)
	if gf.DbHaseField("business_area", "business_id") {
		MDB = MDB.Where("business_id", businessID)
	}

	list, err := MDB.Fields("id,name,level,pid,zip,lng,lat").Order("id asc").Select()
	if err != nil {
		gf.Failed().SetMsg("获取地区失败：" + err.Error()).Regin(c)
		return
	}
	gf.Success().SetMsg("获取地区").SetData(gf.Map{"items": list}).Regin(c)
}
