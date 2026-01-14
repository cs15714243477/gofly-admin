package apidoc

import (
	"gofly/utils/gf"
)

// 接口文档
type Devapigroup struct{ NoNeedAuths []string }

func init() {
	fpath := Devapigroup{NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

// 获取列表
func (api *Devapigroup) Get_list(c *gf.GinCtx) {
	list, err := gf.Model("common_apidoc_group").Order("id asc").Select()
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
	} else {
		for _, val := range list {
			typename, _ := gf.Model("common_apidoc_type").Where("id", val["type_id"]).Value("name")
			val["typename"] = typename
		}
		gf.Success().SetMsg("获取数据列表").SetData(list).Regin(c)
	}
}

// 获取父级数据
func (api *Devapigroup) Get_parent(c *gf.GinCtx) {
	list, _ := gf.Model("common_apidoc_group").Fields("id,pid,name").Order("id asc").Select()
	list = gf.GetMenuChildrenArray(list, 0, "pid")
	gf.Success().SetMsg("获取分组列表").SetData(list).Regin(c)
}

// 保存
func (api *Devapigroup) Save(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	var f_id float64 = 0
	if param["id"] != nil {
		f_id = param["id"].(float64)
	}
	if f_id == 0 {
		delete(param, "id")
		addId, err := gf.Model("common_apidoc_group").Data(param).InsertAndGetId()
		if err != nil {
			gf.Failed().SetMsg("添加失败").SetData(err).Regin(c)
		} else {
			gf.Success().SetMsg("添加成功！").SetData(addId).Regin(c)
		}
	} else {
		res, err := gf.Model("common_apidoc_group").
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
func (api *Devapigroup) UpStatus(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	res2, err := gf.Model("common_apidoc_group").Where("id", param["id"]).Data(map[string]interface{}{"status": param["status"]}).Update()
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
func (api *Devapigroup) Del(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	res2, err := gf.Model("common_apidoc_group").WhereIn("id", param["ids"]).Delete()
	if err != nil {
		gf.Failed().SetMsg("删除失败").SetData(err).Regin(c)
	} else {
		gf.Success().SetMsg("删除成功！").SetData(res2).Regin(c)
	}
}
