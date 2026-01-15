package houses

import (
	"gofly/utils/gf"
	"gofly/utils/tools/gconv"
	"gofly/utils/tools/gmap"
	"gofly/utils/tools/gvar"
)

// 海报模板（business_poster_templates）管理
type Poster struct{ NoNeedAuths []string }

func init() {
	// 默认都需要登录与权限校验（由 business/controller.go 的 RouterHandler 处理）
	fpath := Poster{NoNeedAuths: []string{}}
	gf.Register(&fpath, fpath)
}

func posterSafeFullUrl(url string) string {
	if url == "" {
		return ""
	}
	// gf.GetFullUrl 内部会直接 url[0:4]，这里防止短字符串 panic
	if len(url) < 4 {
		return gf.GetFullUrl(url)
	}
	return gf.GetFullUrl(url)
}

func posterPickMap(src map[string]interface{}, keys []string) gf.Map {
	out := gf.Map{}
	for _, k := range keys {
		if v, ok := src[k]; ok {
			out[k] = v
		}
	}
	return out
}

// 获取海报模板列表
func (api *Poster) GetList(c *gf.GinCtx) {
	pageNo := gconv.Int(c.DefaultQuery("page", "1"))
	pageSize := gconv.Int(c.DefaultQuery("pageSize", "10"))
	param, _ := gf.RequestParam(c)

	whereMap := gmap.New()
	whereMap.Set("business_id", c.GetInt64("businessID"))
	if name, ok := param["name"]; ok && name != "" {
		whereMap.Set("name like ?", "%"+gconv.String(name)+"%")
	}
	if tp, ok := param["template_type"]; ok && tp != "" {
		whereMap.Set("template_type", tp)
	}
	if status, ok := param["status"]; ok && status != "" {
		whereMap.Set("status", status)
	}

	MDB := gf.Model("business_poster_templates").Where(whereMap)
	totalCount, _ := MDB.Clone().Count()
	list, err := MDB.Fields("id,business_id,name,preview_url,template_config,template_type,status,weigh,createtime,updatetime").
		Page(pageNo, pageSize).
		Order("weigh desc,id desc").
		Select()
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	for _, it := range list {
		it["preview_url"] = gvar.New(posterSafeFullUrl(gconv.String(it["preview_url"])))
	}
	gf.Success().SetMsg("获取海报模板列表").SetData(gf.Map{
		"page":     pageNo,
		"pageSize": pageSize,
		"total":    totalCount,
		"items":    list,
	}).Regin(c)
}

// 获取海报模板详情
func (api *Poster) GetContent(c *gf.GinCtx) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		gf.Failed().SetMsg("请传参数id").Regin(c)
		return
	}
	data, err := gf.Model("business_poster_templates").
		Where("business_id", c.GetInt64("businessID")).
		Where("id", id).
		Find()
	if err != nil {
		gf.Failed().SetMsg("获取内容失败").SetData(err).Regin(c)
		return
	}
	if data != nil {
		data["preview_url"] = gvar.New(posterSafeFullUrl(data["preview_url"].String()))
	}
	gf.Success().SetMsg("获取内容成功").SetData(data).Regin(c)
}

// 保存海报模板
func (api *Poster) Save(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	fid := gf.GetEditId(param["id"])

	saveData := posterPickMap(param, []string{
		"name",
		"preview_url",
		"template_config",
		"template_type",
		"status",
		"weigh",
	})

	if fid == 0 {
		saveData["business_id"] = c.GetInt64("businessID")
		addId, err := gf.Model("business_poster_templates").Data(saveData).InsertAndGetId()
		if err != nil {
			gf.Failed().SetMsg("添加失败").SetData(err).Regin(c)
			return
		}
		gf.Success().SetMsg("添加成功").SetData(addId).Regin(c)
		return
	}

	_, err := gf.Model("business_poster_templates").
		Where("business_id", c.GetInt64("businessID")).
		Where("id", fid).
		Update(saveData)
	if err != nil {
		gf.Failed().SetMsg("更新失败").SetData(err).Regin(c)
		return
	}
	gf.Success().SetMsg("更新成功").Regin(c)
}

// 更新状态/权重（兼容旧路由 putStatus）
func (api *Poster) UpStatus(c *gf.GinCtx)  { api.posterUpStatus(c) }
func (api *Poster) PutStatus(c *gf.GinCtx) { api.posterUpStatus(c) }

func (api *Poster) posterUpStatus(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	if param["id"] == nil || gconv.Int64(param["id"]) == 0 {
		gf.Failed().SetMsg("请传参数id").Regin(c)
		return
	}
	update := gf.Map{}
	if v, ok := param["status"]; ok && v != "" {
		update["status"] = gconv.Int(v)
	}
	if v, ok := param["weigh"]; ok && v != "" {
		update["weigh"] = gconv.Int(v)
	}
	if len(update) == 0 {
		gf.Failed().SetMsg("暂无可更新字段").Regin(c)
		return
	}
	_, err := gf.Model("business_poster_templates").
		Where("business_id", c.GetInt64("businessID")).
		Where("id", param["id"]).
		Update(update)
	if err != nil {
		gf.Failed().SetMsg("更新失败").SetData(err).Regin(c)
	} else {
		gf.Success().SetMsg("更新成功").Regin(c)
	}
}

// 删除海报模板
func (api *Poster) Del(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	if param["ids"] == nil {
		gf.Failed().SetMsg("请传参数ids").Regin(c)
		return
	}
	_, err := gf.Model("business_poster_templates").
		Where("business_id", c.GetInt64("businessID")).
		WhereIn("id", param["ids"]).
		Delete()
	if err != nil {
		gf.Failed().SetMsg("删除失败").SetData(err).Regin(c)
	} else {
		gf.Success().SetMsg("删除成功").Regin(c)
	}
}
