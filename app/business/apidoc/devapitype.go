package apidoc

import (
	"gofly/utils/gf"
)

// 接口分类
type Devapitype struct{ NoNeedAuths []string }

func init() {
	fpath := Devapitype{NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

// 获取列表
func (api *Devapitype) Get_list(c *gf.GinCtx) {
	list, err := gf.Model("common_apidoc_type").Order("id asc").Select()
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
	} else {
		gf.Success().SetMsg("获取数据列表").SetData(list).Regin(c)
	}
}

// 获取单条数据
func (api *Devapitype) Get_typeinfo(c *gf.GinCtx) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		gf.Failed().SetMsg("参数id不能为空").Regin(c)
	} else {
		data, _ := gf.Model("common_apidoc_type").Where("id", id).Find()
		gf.Success().SetMsg("获取单条接口类型").SetData(data).Regin(c)
	}
}

// 保存
func (api *Devapitype) Save(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	//当前用户
	var f_id float64 = 0
	if param["id"] != nil {
		f_id = param["id"].(float64)
	}
	if f_id == 0 {
		delete(param, "id")
		addId, err := gf.Model("common_apidoc_type").Data(param).InsertAndGetId()
		if err != nil {
			gf.Failed().SetMsg("添加失败").SetData(err).Regin(c)
		} else {
			gf.Success().SetMsg("添加成功！").SetData(addId).Regin(c)
		}
	} else {
		res, err := gf.Model("common_apidoc_type").
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

// 删除
func (api *Devapitype) Del(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	res, err := gf.Model("common_apidoc_type").WhereIn("id", param["ids"]).Delete()
	if err != nil {
		gf.Failed().SetMsg("删除失败").SetData(err).Regin(c)
	} else {
		gf.Success().SetMsg("删除成功！").SetData(res).Regin(c)
	}
}
