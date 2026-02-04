package datacenter

import (
	"fmt"
	"gofly/utils/gf"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// 配置邮箱
type Configuration struct{}

func init() {
	fpath := Configuration{}
	gf.Register(&fpath, fpath)
}

func getConfigLineList(configLines interface{}) []string {
	switch lines := configLines.(type) {
	case []string:
		return lines
	case []interface{}:
		result := make([]string, 0, len(lines))
		for _, line := range lines {
			result = append(result, gf.String(line))
		}
		return result
	default:
		return []string{}
	}
}

func getOrderedConfigData(lineList []string, dataMap map[string]interface{}) []map[string]interface{} {
	newData := make([]map[string]interface{}, 0, len(dataMap))
	seen := make(map[string]struct{}, len(dataMap))
	inData := false

	for _, line := range lineList {
		trimLine := strings.TrimSpace(line)
		if trimLine == "" || strings.HasPrefix(trimLine, "#") {
			continue
		}
		if !inData {
			if strings.HasPrefix(trimLine, "data:") {
				inData = true
			}
			continue
		}
		if !strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "\t") {
			break
		}
		content := strings.TrimSpace(line)
		if content == "" || strings.HasPrefix(content, "#") {
			continue
		}
		parts := strings.SplitN(content, ":", 2)
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		if key == "" {
			continue
		}
		value, exists := dataMap[key]
		if !exists {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		keyname := key
		if commentIndex := strings.Index(parts[1], "#"); commentIndex >= 0 {
			comment := strings.TrimSpace(parts[1][commentIndex+1:])
			if comment != "" {
				keyname = comment
			}
		}
		newData = append(newData, gf.Map{"keyname": keyname, "keyfield": key, "keyvalue": value})
		seen[key] = struct{}{}
	}

	missingKeys := make([]string, 0, len(dataMap))
	for key := range dataMap {
		if _, ok := seen[key]; !ok {
			missingKeys = append(missingKeys, key)
		}
	}
	sort.Strings(missingKeys)
	for _, key := range missingKeys {
		newData = append(newData, gf.Map{"keyname": key, "keyfield": key, "keyvalue": dataMap[key]})
	}
	return newData
}

// 获取邮箱
func (api *Configuration) GetEmail(c *gf.GinCtx) {
	data, _ := gf.Model("common_email").Where("data_from", "common").Find()
	gf.Success().SetMsg("获取邮箱").SetData(data).Regin(c)
}

// 保存邮箱
func (api *Configuration) SaveEmail(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	GetID, _ := gf.Model("common_email").Where("data_from", "common").Value("id")
	if GetID == nil {
		param["data_from"] = "common"
		addId, err := gf.Model("common_email").Data(param).InsertAndGetId()
		if err != nil {
			gf.Failed().SetMsg("添加失败").SetData(err).Regin(c)
		} else {
			gf.Success().SetMsg("添加成功！").SetData(addId).Regin(c)
		}
	} else {
		res, err := gf.Model("common_email").Data(param).Where("id", GetID).Update()
		if err != nil {
			gf.Failed().SetMsg("更新失败").SetData(err).Regin(c)
		} else {
			gf.Success().SetMsg("更新成功！").SetData(res).Regin(c)
		}
	}
}

// 获取安装的代码仓配置
func (api *Configuration) GetCodestoreConfig(c *gf.GinCtx) {
	path, _ := os.Getwd()
	configPath := filepath.Join(path, "/resource/config")
	go_app_dir, _ := gf.GetAllFileArray(configPath)
	var list []interface{} = make([]interface{}, 0)
	for _, val := range gf.Strings(go_app_dir["files"]) {
		configstr := gf.ReaderFileByline(filepath.Join(path, "/resource/config", val))
		fileName := strings.Split(val, ".")
		install_cofig, _ := gf.GetYmlConfigData(filepath.Join(path, "/resource/config"), fileName[0])
		data_item := install_cofig.(map[string]interface{})
		if data_item["conftype"] == "configuration" {
			dataMap, ok := data_item["data"].(map[string]interface{})
			if !ok {
				continue
			}
			new_data := getOrderedConfigData(getConfigLineList(configstr), dataMap)
			data_item["name"] = fileName[0]
			data_item["data"] = new_data
			list = append(list, install_cofig)
		}
	}
	gf.Success().SetMsg("获取安装的代码仓配置列表").SetData(list).Regin(c)
}

// 修改安装的代码仓配置
func (api *Configuration) SaveCodeStoreConfig(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	path, _ := os.Getwd()
	configPath := filepath.Join(path, "/resource/config", gf.String(param["name"])+".yaml")
	var upAppconf gf.Map = make(map[string]interface{}, 0)
	for _, val := range param["data"].([]interface{}) {
		item := val.(map[string]interface{})
		var val_data interface{}
		switch item["keyvalue"].(type) {
		case string: // 处理 string 类型
			val_str := gf.String(item["keyvalue"])
			if val_str == "" {
				val_str = "\"\""
			} else {
				val_str = strconv.Quote(val_str)
			}
			val_data = val_str
		default:
			val_data = gf.String(item["keyvalue"])
		}
		upAppconf[gf.String(item["keyfield"])] = fmt.Sprintf("%s  #%v", val_data, item["keyname"])
	}
	gf.UpConfigFild(configPath, upAppconf, "  ")
	gf.Success().SetMsg("修改安装的代码仓配置").SetData(param).SetExdata(configPath).Regin(c)
}

// 更新配置使用状态
func (api *Configuration) UpConfigStatus(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	path, _ := os.Getwd()
	//如果改为true，则把其他通用标识改为false
	if gf.Bool(param["status"]) {
		configPath := filepath.Join(path, "/resource/config")
		go_app_dir, _ := gf.GetAllFileArray(configPath)
		for _, val := range gf.Strings(go_app_dir["files"]) {
			fileName := strings.Split(val, ".")
			install_cofig, _ := gf.GetYmlConfigData(filepath.Join(path, "/resource/config"), fileName[0])
			data_item := install_cofig.(map[string]interface{})
			if fileName[0] != gf.String(param["name"]) && gf.String(data_item["pluginident"]) == gf.String(param["pluginident"]) && gf.Bool(data_item["status"]) {
				configPath := filepath.Join(path, "/resource/config", val)
				gf.UpConfigFild(configPath, gf.Map{"status": false}, "")
			}
		}
	}
	configPath := filepath.Join(path, "/resource/config", gf.String(param["name"])+".yaml")
	gf.UpConfigFild(configPath, gf.Map{"status": param["status"]}, "")
	gf.Success().SetMsg("更新配置使用状态成功").SetData(param).SetExdata(configPath).Regin(c)
}
