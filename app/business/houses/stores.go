package houses

import (
	"gofly/utils/gf"
	"gofly/utils/tools/gconv"
	"gofly/utils/tools/gmap"
	"strings"
)

// 门店管理
type Stores struct{}

func init() {
	fpath := Stores{}
	gf.Register(&fpath, fpath)
}

func storesNormalizeCommaText(v interface{}) string {
	s := strings.TrimSpace(gconv.String(v))
	if s == "" {
		return ""
	}
	parts := gf.SplitAndStr(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, ",")
}

func storesPickMap(param map[string]interface{}, keys ...string) gf.Map {
	out := gf.Map{}
	for _, k := range keys {
		if v, ok := param[k]; ok {
			out[k] = v
		}
	}
	return out
}

// 获取列表
func (api *Stores) GetList(c *gf.GinCtx) {
	pageNo := gf.Int(c.DefaultQuery("page", "1"))
	pageSize := gf.Int(c.DefaultQuery("pageSize", "10"))
	//搜索添条件
	param, _ := gf.RequestParam(c)
	whereMap := gmap.New()
	whereMap.Set("business_id", c.GetInt64("businessID"))
	if name, ok := param["name"]; ok && name != "" {
		whereMap.Set("name like ?", "%"+gf.String(name)+"%")
	}
	if address, ok := param["address"]; ok && address != "" {
		whereMap.Set("address like ?", "%"+gf.String(address)+"%")
	}
	if contact_phone, ok := param["contact_phone"]; ok && contact_phone != "" {
		whereMap.Set("contact_phone like ?", "%"+gf.String(contact_phone)+"%")
	}
	if manager_name, ok := param["manager_name"]; ok && manager_name != "" {
		whereMap.Set("manager_name", manager_name)
	}
	if status, ok := param["status"]; ok && status != "" {
		whereMap.Set("status", status)
	}
	MDB := gf.Model("business_stores").Where(whereMap)
	totalCount, _ := MDB.Clone().Count()
	list, err := MDB.Fields("id,name,address,contact_phone,manager_name,images,status,weigh,createtime").Page(pageNo, pageSize).Order("id desc").Select()
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

// 保存
func (api *Stores) Save(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	var f_id = gf.GetEditId(param["id"])

	// 仅允许写入字段，禁止前端传 business_id 等敏感字段
	saveData := storesPickMap(param,
		"name", "address", "contact_phone", "manager_name",
		"images", "status", "weigh",
	)
	if _, ok := saveData["images"]; ok {
		saveData["images"] = storesNormalizeCommaText(saveData["images"])
	}

	if f_id == 0 {
		saveData["business_id"] = c.GetInt64("businessID") //当前用户商户ID
		addId, err := gf.Model("business_stores").Data(saveData).InsertAndGetId()
		if err != nil {
			gf.Failed().SetMsg("添加失败").SetData(err).Regin(c)
		} else {
			if addId != 0 {
				gf.Model("business_stores").Where("id", addId).Update(gf.Map{"weigh": addId})
			}
			gf.Success().SetMsg("添加成功！").SetData(addId).Regin(c)
		}
	} else {
		res, err := gf.Model("business_stores").
			Where("business_id", c.GetInt64("businessID")).
			Where("id", f_id).
			Update(saveData)
		if err != nil {
			gf.Failed().SetMsg("更新失败").SetData(err).Regin(c)
		} else {
			gf.Success().SetMsg("更新成功！").SetData(res).Regin(c)
		}
	}
}

// 更新状态
func (api *Stores) UpStatus(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	if param["id"] == nil || gconv.Int64(param["id"]) == 0 {
		gf.Failed().SetMsg("请传参数id").Regin(c)
		return
	}
	// 仅允许更新指定字段，避免任意字段更新
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
	res, err := gf.Model("business_stores").
		Where("business_id", c.GetInt64("businessID")).
		Where("id", param["id"]).
		Update(update)
	if err != nil {
		gf.Failed().SetMsg("更新失败！").SetData(err).Regin(c)
	} else {
		msg := "更新成功！"
		if res == nil {
			msg = "暂无数据更新"
		}
		gf.Success().SetMsg(msg).SetData(res).Regin(c)
	}
}

// 删除
func (api *Stores) Del(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	res, err := gf.Model("business_stores").Where("business_id", c.GetInt64("businessID")).WhereIn("id", param["ids"]).Delete()
	if err != nil {
		gf.Failed().SetMsg("删除失败").SetData(err).Regin(c)
	} else {
		gf.Success().SetMsg("删除成功！").SetData(res).Regin(c)
	}
}

// 获取内容
func (api *Stores) GetContent(c *gf.GinCtx) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		gf.Failed().SetMsg("请传参数id").Regin(c)
	} else {
		data, err := gf.Model("business_stores").Where("business_id", c.GetInt64("businessID")).Where("id", id).Find()
		if err != nil {
			gf.Failed().SetMsg("获取内容失败").SetData(err).Regin(c)
		} else {

			gf.Success().SetMsg("获取内容成功！").SetData(data).Regin(c)
		}
	}
}
