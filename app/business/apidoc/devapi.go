package apidoc

import (
	"gofly/utils/gf"
	"gofly/utils/tools/gcfg"
	"gofly/utils/tools/gconv"
	"gofly/utils/tools/gctx"
	"gofly/utils/tools/gmap"
	"gofly/utils/tools/gvar"
	"os"
	"strings"

	"gofly/utils/gform"
)

// 接口文档
type Devapi struct{ NoNeedAuths []string }

func init() {
	fpath := Devapi{NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

var (
	ctx        = gctx.New()
	dbConf, _  = gcfg.Instance().Get(ctx, "database.default")
	dbConf_arr = gconv.Map(dbConf)
)

// 获取部门列表
func (api *Devapi) GetList(c *gf.GinCtx) {
	pageNo := gconv.Int(c.DefaultQuery("page", "1"))
	pageSize := gconv.Int(c.DefaultQuery("pageSize", "10"))
	//搜索添条件
	param, _ := gf.RequestParam(c)
	whereMap := gmap.New()
	if cid, ok := param["cid"]; ok && cid != "0" {
		whereMap.Set("cid", cid)
	}
	if title, ok := param["title"]; ok && title != "" {
		whereMap.Set("title like ?", "%"+gconv.String(title)+"%")
	}
	if url, ok := param["url"]; ok && url != "" {
		whereMap.Set("url like ?", "%"+gconv.String(url)+"%")
	}
	MDB := gf.Model("common_apidoc").Where(whereMap)
	totalCount, _ := MDB.Clone().Count()
	list, err := MDB.Page(pageNo, pageSize).Order("id desc").Select()
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
	} else {
		for _, val := range list {
			groupdata, _ := gf.Model("common_apidoc_group a").
				LeftJoin("common_apidoc_type t", "t.id = a.type_id").
				Where("a.id", val["cid"]).Fields("a.name,a.type_id").Find()
			val["groupname"] = groupdata["name"]
			val["type_id"] = groupdata["type_id"]
		}
		gf.Success().SetMsg("获取数据列表").SetData(gf.Map{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list}).SetExdata(whereMap).Regin(c)
	}
}

// 获取分组
func (api *Devapi) GetGroup(c *gf.GinCtx) {
	list, _ := gf.Model("common_apidoc_group").Fields("id,pid,name").Order("id asc").Select()
	list = gf.GetMenuChildrenArray(list, 0, "pid")
	gf.Success().SetMsg("获取分组列表").SetData(list).Regin(c)
}

// 保存
func (api *Devapi) Save(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	var f_id float64 = 0
	if param["id"] != nil {
		f_id = param["id"].(float64)
	}
	if f_id == 0 {
		if gf.Int(param["apicode_type"]) == 2 {
			api_id, err := gf.CreateAndUpdateApi(gf.Map{"title": param["title"], "tablename": param["tablename"], "fields": param["fields"], "istoken": param["istoken"]})
			if err != nil {
				gf.Failed().SetMsg("创建低代码接口失败").SetData(err).Regin(c)
				return
			}
			param["api_id"] = api_id
		}
		addId, err := gf.Model("common_apidoc").Data(param).InsertAndGetId()
		if err != nil {
			gf.Failed().SetMsg("添加失败").SetData(err).Regin(c)
		} else {
			gf.Success().SetMsg("添加成功！").SetData(addId).Regin(c)
		}
	} else {
		res, err := gf.Model("common_apidoc").Data(param).Where("id", f_id).Update()
		if err != nil {
			gf.Failed().SetMsg("更新失败").SetData(err).Regin(c)
		} else {
			if gf.Int(param["apicode_type"]) == 2 {
				gf.CreateAndUpdateApi(gf.Map{"api_id": param["api_id"], "title": param["title"], "tablename": param["tablename"], "fields": param["fields"], "istoken": param["istoken"]})
			}
			gf.Success().SetMsg("更新成功！").SetData(res).Regin(c)
		}
	}
}

// 更新状态
func (api *Devapi) UpStatus(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	res2, err := gf.Model("common_apidoc").Where("id", param["id"]).Data(map[string]interface{}{"status": param["status"]}).Update()
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
func (api *Devapi) Del(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	res2, err := gf.Model("common_apidoc").WhereIn("id", param["ids"]).Delete()
	if err != nil {
		gf.Failed().SetMsg("删除失败").SetData(err).Regin(c)
	} else {
		gf.Success().SetMsg("删除成功！").SetData(res2).Regin(c)
	}
}

// 获取数据库字段
func (api *Devapi) GetDBField(c *gf.GinCtx) {
	tablename := c.DefaultQuery("tablename", "")
	if tablename == "" {
		gf.Failed().SetMsg("请传数据表名称").Regin(c)
	} else {
		tablename_arr := strings.Split(tablename, ",")
		//获取数据库名
		var dielddata_list gform.Result
		for _, Val := range tablename_arr {
			dielddata, _ := gf.DB().Query(c, "select COLUMN_NAME,COLUMN_COMMENT,COLUMN_TYPE,DATA_TYPE,CHARACTER_MAXIMUM_LENGTH,IS_NULLABLE,COLUMN_DEFAULT,NUMERIC_PRECISION from information_schema.columns where TABLE_SCHEMA='"+gconv.String(dbConf_arr["dbname"])+"' AND TABLE_NAME='"+Val+"'")
			for _, data := range dielddata {
				data["tablename"] = gvar.New(Val)
				dielddata_list = append(dielddata_list, data)
			}
		}
		gf.Success().SetMsg("获取数据库字段").SetData(dielddata_list).SetExdata(tablename).Regin(c)
	}
}

// 获取所有路由列表
func (api *Devapi) GetRoutes(c *gf.GinCtx) {
	filePath := "runtime/app/routers.txt"
	list := gf.ReaderFileByline(filePath)
	gf.Success().SetMsg("获取所有路由列表").SetData(list).Regin(c)
}

// 获取模块列表
func (api *Devapi) GetModel(c *gf.GinCtx) {
	var list []string
	files, err := os.ReadDir("./app")
	if err != nil {
		gf.Failed().SetMsg("获取目录错误").SetData(err).Regin(c)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			list = append(list, file.Name())
		}
	}
	gf.Success().SetMsg("获取模块目录列表").SetData(list).Regin(c)
}

// 获取指定数据库表的数据
func (api *Devapi) GetTablelist(c *gf.GinCtx) {
	tablename := c.DefaultQuery("tablename", "")
	pageSize := c.DefaultQuery("pageSize", "10")
	if tablename == "" {
		gf.Failed().SetMsg("请传数据库名称").Regin(c)
	} else {
		seachword := c.DefaultQuery("seachword", "")
		MDB := gf.Model(tablename)
		if seachword != "" {
			MDB.Where("name", "like", "%"+seachword+"%")
		}
		list, err := MDB.Limit(gconv.Int(pageSize)).Order("id desc").Select()
		if err != nil {
			gf.Failed().SetMsg(err.Error()).Regin(c)
		} else {
			gf.Success().SetMsg("获取指定数据库表的数据").SetData(list).Regin(c)
		}
	}
}
