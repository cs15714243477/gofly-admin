package houses

import (
	"gofly/utils/gf"
	"gofly/utils/gform"
	"gofly/utils/tools/gconv"
	"gofly/utils/tools/gmap"
	"strings"
	"time"
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
	if v, ok := param["can_manage_locks"]; ok && v != "" {
		whereMap.Set("can_manage_locks", v)
	}
	if v, ok := param["audit_status"]; ok && v != "" && gf.DbHaseField("business_user", "audit_status") {
		whereMap.Set("audit_status", v)
	}

	MDB := gf.Model("business_user").Where(whereMap)
	totalCount, _ := MDB.Clone().Count()
	fields := "id,business_id,username,name,nickname,remark,email,mobile,avatar,sex,role,can_manage_properties,can_manage_locks,store_id,title,introduction,status,createtime,updatetime"
	// 兼容：手填门店与审核字段可能是增量字段
	for _, col := range []string{
		"store_name_text", "store_address_text",
		"region_province", "region_city", "region_district",
		"audit_status", "audit_reason", "apply_time", "audit_time",
	} {
		if gf.DbHaseField("business_user", col) {
			fields += "," + col
		}
	}
	list, err := MDB.Fields(fields).
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
		// 手填门店兜底：当 store_id 为空，或门店表缺失/已删除时可展示
		for _, it := range list {
			sid := gconv.Int64(it["store_id"])
			if sid != 0 {
				// 已有门店数据则不覆盖
				if v, ok := it["store_name"]; ok && strings.TrimSpace(gconv.String(v)) != "" {
					continue
				}
			}
			if gf.DbHaseField("business_user", "store_name_text") {
				if s := strings.TrimSpace(gconv.String(it["store_name_text"])); s != "" {
					it["store_name"] = gf.VarNew(s)
				}
			}
			if gf.DbHaseField("business_user", "store_address_text") {
				if s := strings.TrimSpace(gconv.String(it["store_address_text"])); s != "" {
					it["store_address"] = gf.VarNew(s)
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
		"store_name_text", "store_address_text",
		"region_province", "region_city", "region_district",
		"audit_status", "audit_reason",
		"can_manage_properties",
		"can_manage_locks",
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
	if _, ok := saveData["can_manage_locks"]; ok {
		v := gconv.Int(saveData["can_manage_locks"])
		if v != 0 && v != 1 {
			gf.Failed().SetMsg("can_manage_locks参数不合法").Regin(c)
			return
		}
		saveData["can_manage_locks"] = v
	}

	// 审核字段校验（新增字段：老库可不存在）
	if v, ok := saveData["audit_status"]; ok {
		if !gf.DbHaseField("business_user", "audit_status") {
			delete(saveData, "audit_status")
		} else {
			as := strings.ToLower(strings.TrimSpace(gconv.String(v)))
			switch as {
			case "pending", "approved", "rejected":
				saveData["audit_status"] = as
			default:
				gf.Failed().SetMsg("audit_status参数不合法").Regin(c)
				return
			}

			// 审核时间自动写入（仅当字段存在）
			if gf.DbHaseField("business_user", "audit_time") {
				saveData["audit_time"] = time.Now()
			}

			// 拒绝时要求填写原因
			if as == "rejected" {
				reason := strings.TrimSpace(gconv.String(saveData["audit_reason"]))
				if reason == "" {
					gf.Failed().SetMsg("拒绝时请填写审核原因").Regin(c)
					return
				}
				if gf.DbHaseField("business_user", "audit_reason") {
					saveData["audit_reason"] = reason
				} else {
					delete(saveData, "audit_reason")
				}
			} else {
				// 非拒绝：若原因字段存在且未传，则清空（避免残留）
				if gf.DbHaseField("business_user", "audit_reason") {
					if _, has := saveData["audit_reason"]; !has {
						saveData["audit_reason"] = ""
					}
				} else {
					delete(saveData, "audit_reason")
				}
			}
		}
	} else {
		// 未传 audit_status：若 audit_reason 单独传入则丢弃（避免无意义更新）
		if _, ok := saveData["audit_reason"]; ok && !gf.DbHaseField("business_user", "audit_reason") {
			delete(saveData, "audit_reason")
		}
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
