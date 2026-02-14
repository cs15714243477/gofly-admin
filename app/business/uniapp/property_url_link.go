package uniapp

import (
	"strconv"
	"strings"

	"gofly/utils/gf"
)

// GetPropertyUrlLink 获取房源推广 URL Link（微信内可点击跳转小程序）
//
// 入参：
// - property_id: 房源ID（必填）
// - env_version: 可选，develop|trial|release（不传则使用默认）
func (api *Index) GetPropertyUrlLink(c *gf.GinCtx) {
	userID := c.GetInt64("userID")
	if userID <= 0 {
		gf.Failed().SetCode(401).SetMsg("请先登录").Regin(c)
		return
	}
	businessID := wxBusinessID(c)

	propertyID := gf.Int64(c.DefaultQuery("property_id", "0"))
	if propertyID <= 0 {
		gf.Failed().SetMsg("缺少参数 property_id").Regin(c)
		return
	}

	prop, err := gf.Model("business_properties").
		Fields("id,title").
		Where("id", propertyID).
		Where("business_id", businessID).
		Where("deletetime", nil).
		Where("status", 0).
		WhereNotIn("sale_status", wxHiddenSaleStatuses()).
		Find()
	if err != nil || prop == nil || prop["id"].Int64() <= 0 {
		gf.Failed().SetMsg("房源不存在或已下架").Regin(c)
		return
	}

	envVersion := strings.TrimSpace(strings.ToLower(c.DefaultQuery("env_version", "")))
	if envVersion != "" && envVersion != "develop" && envVersion != "trial" && envVersion != "release" {
		envVersion = ""
	}

	appid, secretkey := getWxappAppidSecret(c)
	if appid == "" || secretkey == "" {
		gf.Failed().SetMsg("未配置微信小程序 appid/secretkey").Regin(c)
		return
	}
	accessToken, err := wxGetAccessToken(appid, secretkey)
	if err != nil {
		gf.Failed().SetMsg("获取微信 access_token 失败：" + err.Error()).Regin(c)
		return
	}

	path := "pages/property_detail/property_detail"
	query := "id=" + strconv.FormatInt(propertyID, 10) +
		"&public=1&from_agent_id=" + strconv.FormatInt(userID, 10) +
		"&from_style=0"

	urlLink, err := wxGenerateUrlLink(accessToken, path, query, envVersion, 60*60*24*30)
	if err != nil {
		gf.Failed().SetMsg("获取房源推广链接失败：" + err.Error()).Regin(c)
		return
	}

	gf.Success().SetMsg("获取房源推广链接").SetData(gf.Map{
		"property_id": propertyID,
		"title":       prop["title"].String(),
		"url_link":    urlLink,
	}).Regin(c)
}
