package houses

import (
	"gofly/utils/gf"
	"gofly/utils/gform"
	"gofly/utils/tools/gconv"
	"gofly/utils/tools/gmap"
)

// 经纪人（business_user）管理
type Broker struct{ NoNeedAuths []string }

func init() {
	// 默认都需要登录与权限校验（由 business/controller.go 的 RouterHandler 处理）
	fpath := Broker{NoNeedAuths: []string{}}
	gf.Register(&fpath, fpath)
}

// 获取经纪人列表
func (api *Broker) GetList(c *gf.GinCtx) {
	pageNo := gconv.Int(c.DefaultQuery("page", "1"))
	pageSize := gconv.Int(c.DefaultQuery("pageSize", "10"))
	param, _ := gf.RequestParam(c)

	whereMap := gmap.New()
	whereMap.Set("business_id", c.GetInt64("businessID"))
	// 搜索条件
	if name, ok := param["name"]; ok && name != "" {
		whereMap.Set("name like ?", "%"+gconv.String(name)+"%")
	}
	if mobile, ok := param["mobile"]; ok && mobile != "" {
		whereMap.Set("mobile like ?", "%"+gconv.String(mobile)+"%")
	}
	if username, ok := param["username"]; ok && username != "" {
		whereMap.Set("username like ?", "%"+gconv.String(username)+"%")
	}
	if status, ok := param["status"]; ok && status != "" {
		whereMap.Set("status", status)
	}
	if storeID, ok := param["store_id"]; ok && storeID != "" && gconv.Int64(storeID) != 0 {
		whereMap.Set("store_id", storeID)
	}
	if role, ok := param["role"]; ok && role != "" {
		whereMap.Set("role", role)
	}
	if v, ok := param["can_manage_properties"]; ok && v != "" {
		whereMap.Set("can_manage_properties", v)
	}

	MDB := gf.Model("business_user").Where(whereMap)
	totalCount, _ := MDB.Clone().Count()
	list, err := MDB.Fields("id,business_id,username,name,nickname,remark,email,mobile,avatar,sex,role,can_manage_properties,store_id,title,introduction,status,createtime,updatetime").
		Page(pageNo, pageSize).
		Order("id desc").
		Select()
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}

	// 补齐门店信息（名称/地址/电话/店长）
	if len(list) > 0 {
		storeIDs := make([]int64, 0, len(list))
		storeIDSet := map[int64]struct{}{}
		for _, it := range list {
			sid := gconv.Int64(it["store_id"])
			if sid > 0 {
				if _, ok := storeIDSet[sid]; !ok {
					storeIDSet[sid] = struct{}{}
					storeIDs = append(storeIDs, sid)
				}
			}
		}
		if len(storeIDs) > 0 {
			stores, _ := gf.Model("business_stores").
				Where("business_id", c.GetInt64("businessID")).
				WhereIn("id", storeIDs).
				Fields("id,name,address,contact_phone,manager_name").
				Select()
			storeMap := map[int64]gform.Record{}
			for _, s := range stores {
				storeMap[gconv.Int64(s["id"])] = s
			}
			for _, it := range list {
				sid := gconv.Int64(it["store_id"])
				if sid == 0 {
					continue
				}
				if s, ok := storeMap[sid]; ok {
					it["store_name"] = s["name"]
					it["store_address"] = s["address"]
					it["store_contact_phone"] = s["contact_phone"]
					it["store_manager_name"] = s["manager_name"]
				}
			}
		}
	}
	gf.Success().SetMsg("获取经纪人列表").SetData(gf.Map{
		"page":     pageNo,
		"pageSize": pageSize,
		"total":    totalCount,
		"items":    list,
	}).Regin(c)
}

// 获取经纪人详情
func (api *Broker) GetContent(c *gf.GinCtx) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		gf.Failed().SetMsg("请传参数id").Regin(c)
		return
	}
	data, err := gf.Model("business_user").Where("business_id", c.GetInt64("businessID")).Where("id", id).Find()
	if err != nil {
		gf.Failed().SetMsg("获取内容失败").SetData(err).Regin(c)
		return
	}
	gf.Success().SetMsg("获取内容成功").SetData(data).Regin(c)
}

// 保存经纪人
func (api *Broker) Save(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	fid := gf.GetEditId(param["id"])

	// 仅允许写入字段
	saveData := gf.Map{}
	for _, k := range []string{
		"username", "name", "nickname", "remark",
		"password", "salt",
		"email", "mobile", "avatar",
		"sex", "role", "store_id",
		"can_manage_properties",
		"title", "introduction",
		"status",
	} {
		if v, ok := param[k]; ok {
			saveData[k] = v
		}
	}
	if _, ok := saveData["can_manage_properties"]; ok {
		v := gconv.Int(saveData["can_manage_properties"])
		if v != 0 && v != 1 {
			gf.Failed().SetMsg("can_manage_properties参数不合法").Regin(c)
			return
		}
		saveData["can_manage_properties"] = v
	}

	// 固定 business_id，避免越权
	if fid == 0 {
		saveData["business_id"] = c.GetInt64("businessID")
		addId, err := gf.Model("business_user").Data(saveData).InsertAndGetId()
		if err != nil {
			gf.Failed().SetMsg("添加失败").SetData(err).Regin(c)
			return
		}
		gf.Success().SetMsg("添加成功").SetData(addId).Regin(c)
		return
	}

	_, err := gf.Model("business_user").Where("business_id", c.GetInt64("businessID")).Where("id", fid).Update(saveData)
	if err != nil {
		gf.Failed().SetMsg("更新失败").SetData(err).Regin(c)
		return
	}
	gf.Success().SetMsg("更新成功").Regin(c)
}

// 更新状态
func (api *Broker) UpStatus(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	if param["id"] == nil || gconv.Int64(param["id"]) == 0 {
		gf.Failed().SetMsg("请传参数id").Regin(c)
		return
	}
	update := gf.Map{}
	if v, ok := param["status"]; ok && v != "" {
		update["status"] = gconv.Int(v)
	}
	if len(update) == 0 {
		gf.Failed().SetMsg("暂无可更新字段").Regin(c)
		return
	}
	_, err := gf.Model("business_user").Where("business_id", c.GetInt64("businessID")).Where("id", param["id"]).Update(update)
	if err != nil {
		gf.Failed().SetMsg("更新失败").SetData(err).Regin(c)
	} else {
		gf.Success().SetMsg("更新成功").Regin(c)
	}
}

// 删除经纪人
func (api *Broker) Del(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	if param["ids"] == nil {
		gf.Failed().SetMsg("请传参数ids").Regin(c)
		return
	}
	_, err := gf.Model("business_user").Where("business_id", c.GetInt64("businessID")).WhereIn("id", param["ids"]).Delete()
	if err != nil {
		gf.Failed().SetMsg("删除失败").SetData(err).Regin(c)
	} else {
		gf.Success().SetMsg("删除成功").Regin(c)
	}
}
