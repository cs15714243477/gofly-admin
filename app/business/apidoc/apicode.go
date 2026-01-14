package apidoc

import (
	"gofly/utils/gf"
	"strings"
)

// 用于自动注册路由
type Apicode struct{ NoNeedAuths []string }

func init() {
	fpath := Apicode{NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

// 生成api接口代码
func (api *Apicode) Installcode(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	data, err := gf.Model("common_apidoc").Where("id", param["id"]).Fields("cid,url,getdata_type,tablename,apicode_type,is_install,fields,method").Find()
	if err != nil {
		gf.Failed().SetMsg("生成api接口代码失败").Regin(c)
	} else {
		type_id, _ := gf.Model("common_apidoc_group").Where("id", data["cid"]).Value("type_id")
		rooturl, _ := gf.Model("common_apidoc_type").Where("id", type_id).Value("model_name")
		model_name := "business" //默认后台接口
		model_name_str := rooturl.String()
		if model_name_str != "" { //模块名称
			model_name = model_name_str
		}
		//创建文件
		if data["url"].String() == "" {
			gf.Failed().SetMsg("url地址为空").Regin(c)
		} else if data["tablename"].String() == "" {
			gf.Failed().SetMsg("数据库表不能为空,选择数据表表提交保存再生成代码").Regin(c)
		} else {
			CreatApicodeFile(model_name, data)
			gf.Model("common_apidoc").
				Data(map[string]interface{}{"is_install": 1}).
				Where("id", param["id"]).
				Update()
			gf.Success().SetMsg("生成api接口代码成功！").SetData(data).Regin(c)
		}
	}
}

// 卸载api接口代码-改变方法
func (api *Apicode) Uninstallcode(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	data, err := gf.Model("common_apidoc").Where("id", param["id"]).Fields("cid,url,getdata_type,tablename,apicode_type,is_install,fields,method").Find()
	if err != nil {
		gf.Failed().SetMsg("卸载失败").Regin(c)
	} else {
		UnApicodeFile(data)
		gf.Model("common_apidoc").Data(map[string]interface{}{"is_install": 2}).Where("id", param["id"]).Update()
		gf.Success().SetMsg("卸载成功！").SetData(data).Regin(c)
	}
}

// 删除文件
func (api *Apicode) RemoveFile(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	data, err := gf.Model("common_apidoc").Where("id", param["id"]).Fields("cid,url,getdata_type,tablename,apicode_type,is_install,fields,method").Find()
	if err != nil {
		gf.Failed().SetMsg("删除文件失败").SetData(err).Regin(c)
	} else {
		type_id, _ := gf.Model("common_apidoc_group").Where("id", data["cid"]).Value("type_id")
		rooturl, _ := gf.Model("common_apidoc_type").Where("id", type_id).Value("model_name")
		model_name := "business" //默认后台接口
		model_name_str := rooturl.String()
		if model_name_str != "" { //模块名称
			model_name = model_name_str
		}
		//判断删除文件
		url := data["url"].String()
		url_arr := strings.Split(url, `/`)
		filename := url_arr[len(url_arr)-1]
		model_path := strings.Split(url, filename)
		haselist, _ := gf.Model("common_apidoc").Where("url", "like", model_path[0]+"%").Where("is_install", 1).Count("*")
		if haselist == 0 {
			RemoveModel(model_name, data) //删除文件-如果没人其他文件则移除文件夹及路由
		}
		gf.Model("common_apidoc").Data(map[string]interface{}{"is_install": 0}).Where("id", param["id"]).Update()
		gf.Success().SetMsg("删除文件成功！").SetData(data).SetExdata(haselist).Regin(c)
	}
}
