package area

import (
	"gofly/utils/gf"
	"gofly/utils/tools/gmap"
)

// 仅分类数据管理
type Index struct{ NoNeedAuths []string }

func init() {
	fpath := Index{NoNeedAuths: []string{"GetMoreList"}}
	gf.Register(&fpath, fpath)
}

// 获取列表
func (api *Index) GetList(c *gf.GinCtx) {
	pageNo := gf.Int(c.DefaultQuery("page", "1"))
	pageSize := gf.Int(c.DefaultQuery("pageSize", "10"))
	//搜索添条件
	param, _ := gf.RequestParam(c)
	whereMap := gmap.New()
	isSearch := true
	if gf.DbHaseField("business_area", "business_id") {
		businessID, _ := c.Get("businessID") //当前用户ID
		whereMap.Set("business_id", businessID)
	}
	if name, ok := param["name"]; ok && name != "" {
		whereMap.Set("name like ?", "%"+gf.String(name)+"%")
		isSearch = false
	}
	if level, ok := param["level"]; ok && level != "" {
		whereMap.Set("level", level)
		isSearch = false
	}
	if zip, ok := param["zip"]; ok && zip != "" {
		whereMap.Set("zip", zip)
		isSearch = false
	}
	if isSearch {
		whereMap.Set("pid", 0)
	}
	MDB := gf.Model("business_area").Where(whereMap)
	totalCount, _ := MDB.Clone().Count()
	list, err := MDB.Fields("id,name,level,zip,lng,lat").Order("id asc").Select()
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
	} else {
		gf.Success().SetMsg("获取全部列表").SetData(gf.Map{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list}).Regin(c)
	}
}

// 获取对应的子数据
func (api *Index) GetMoreList(c *gf.GinCtx) {
	pid := c.DefaultQuery("pid", "0")
	list, _ := gf.Model("business_area").Where("pid", pid).Select()
	for _, val := range list {
		hase, _ := gf.Model("business_area").Where("pid", val["id"]).Count()
		if hase == 0 {
			val["isLeaf"] = gf.VarNew(true)
		}
	}
	gf.Success().SetMsg("获取全部列表").SetData(list).Regin(c)
}

// 保存
func (api *Index) Save(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	var f_id = gf.GetEditId(param["id"])
	if f_id == 0 {
		if gf.DbHaseField("business_area", "business_id") {
			param["business_id"], _ = c.Get("businessID") //当前用户ID
		}
		addId, err := gf.Model("business_area").Data(param).InsertAndGetId()
		if err != nil {
			gf.Failed().SetMsg("添加失败").SetData(err).Regin(c)
		} else {
			if addId != 0 {
				gf.Model("business_area").
					Data(gf.Map{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			gf.Success().SetMsg("添加成功！").SetData(addId).Regin(c)
		}
	} else {
		res, err := gf.Model("business_area").
			Data(param).
			Where("id", f_id).
			Update()
		if err != nil {
			gf.Failed().SetMsg("更新失败").SetData(err).Regin(c)
		} else {
			gf.Success().SetMsg("更新成功！").SetData(res).Regin(c)
		}
	}
}

// 更新状态
func (api *Index) UpStatus(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	res2, err := gf.Model("business_area").Where("id", param["id"]).Data(gf.Map{"status": param["status"]}).Update()
	if err != nil {
		gf.Failed().SetMsg("更新失败！").SetData(err).Regin(c)
	} else {
		msg := "更新成功！"
		if res2 == nil {
			msg = "暂无数据更新"
		}
		gf.Success().SetMsg(msg).SetData(res2).Regin(c)
	}
}

// 删除
func (api *Index) Del(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	res2, err := gf.Model("business_area").WhereIn("id", param["ids"]).Delete()
	if err != nil {
		gf.Failed().SetMsg("删除失败").SetData(err).Regin(c)
	} else {
		gf.Success().SetMsg("删除成功！").SetData(res2).Regin(c)
	}
}

// 获取内容
func (api *Index) GetContent(c *gf.GinCtx) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		gf.Failed().SetMsg("请传参数id").Regin(c)
	} else {
		data, err := gf.Model("business_area").Where("id", id).Find()
		if err != nil {
			gf.Failed().SetMsg("获取内容失败").SetData(err).Regin(c)
		} else {

			gf.Success().SetMsg("获取内容成功！").SetData(data).Regin(c)
		}
	}
}
