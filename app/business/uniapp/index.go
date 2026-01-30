package uniapp

import (
	"bytes"
	"encoding/json"
	"errors"
	"gofly/utils/gf"
	"gofly/utils/tools/gcfg"
	"gofly/utils/tools/gtime"
	"gofly/utils/tools/gvar"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/*
*
* uni-app / 小程序端接口
 */
type Index struct {
	NoNeedLogin []string // 忽略登录接口配置
	NoNeedAuths []string // 忽略RBAC权限认证接口配置
}

func init() {
	fpath := Index{NoNeedLogin: []string{"login", "wxLogin", "logout", "getAgentCard"}, NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

// 《登录》（手机号验证码登录）
// 入参：mobile、captcha
// 说明：小程序端用户表使用 business_user（经纪人表）
func (api *Index) Login(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	mobile := gf.String(param["mobile"])
	if mobile == "" {
		gf.Failed().SetMsg("请提交手机号").Regin(c)
		return
	}
	code, emerr := gf.GetVerifyCode(mobile)
	if emerr != nil || code != gf.Int(param["captcha"]) {
		gf.Failed().SetMsg("验证码无效").SetData(emerr).Regin(c)
		return
	}

	// business_id：可由前端传 business_id；未传则使用默认 1（兼容单商户场景）
	businessID := int64(1)
	if v, ok := param["business_id"]; ok {
		if vv := gf.Int64(v); vv > 0 {
			businessID = vv
		}
	}

	// 查找用户（手机号即账号）
	user, err := gf.Model("business_user").
		Fields("id,business_id,status,logintime").
		Where("mobile", mobile).
		Where("deletetime", nil).
		Find()
	// 未绑定则自动创建（避免再次出现“手机号未绑定账号”）
	if err != nil || user == nil {
		addID, addErr := gf.Model("business_user").Data(gf.Map{
			"business_id": businessID,
			"username":    mobile,
			"name":        "",
			"nickname":    mobile,
			"mobile":      mobile,
			"avatar":      "resource/uploads/static/avatar.png",
			"sex":         0,
			"role":        "user",
			"status":      0,
			"loginip":     gf.GetIp(c),
			"prevtime":    0,
			"logintime":   gtime.Timestamp(),
		}).InsertAndGetId()
		if addErr != nil || addID == 0 {
			gf.Failed().SetMsg("创建用户失败").SetData(addErr).Regin(c)
			return
		}
		user, _ = gf.Model("business_user").Fields("id,business_id,status,logintime").Where("id", addID).Find()
	}

	if user["status"].Int() == 1 {
		gf.Failed().SetMsg("账号被禁用了").Regin(c)
		return
	}
	token, err := gf.CreateToken(gf.Map{
		"ID":          user["id"].Int64(),
		"account_id":  0,
		"business_id": user["business_id"].Int64(),
	})
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}

	// 更新登录信息
	prev := user["logintime"].Int64()
	gf.Model("business_user").
		Where("id", user["id"]).
		Data(gf.Map{"prevtime": prev, "logintime": gtime.Timestamp(), "loginip": gf.GetIp(c)}).
		Update()

	gf.AddloginLog(c, gf.Map{"uid": user["id"], "account_id": 0, "business_id": user["business_id"], "type": "uniapp", "status": 0, "des": "小程序手机号验证码登录"})
	gf.Success().SetMsg("登录成功！").SetData(token).SetToken(gf.String(token)).Regin(c)
}

// 《微信手机号一键登录（小程序）》
// 入参：
// - wx_code: uni.login() 获取的 code
// - phone_code: getPhoneNumber 返回的 e.detail.code
func (api *Index) WxLogin(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	wxCode := gf.String(param["wx_code"])
	phoneCode := gf.String(param["phone_code"])
	if wxCode == "" || phoneCode == "" {
		gf.Failed().SetMsg("缺少参数 wx_code / phone_code").Regin(c)
		return
	}

	appid, secretkey := getWxappAppidSecret(c)
	if appid == "" || secretkey == "" {
		gf.Failed().SetMsg("未配置微信小程序 appid/secretkey，无法一键登录").Regin(c)
		return
	}

	accessToken, err := wxGetAccessToken(appid, secretkey)
	if err != nil {
		gf.Failed().SetMsg("获取微信 access_token 失败：" + err.Error()).Regin(c)
		return
	}

	phone, err := wxGetPhoneNumber(accessToken, phoneCode)
	if err != nil {
		gf.Failed().SetMsg("获取手机号失败：" + err.Error()).Regin(c)
		return
	}

	// business_id：可由前端传 business_id；未传则使用默认 1（兼容单商户场景）
	businessID := int64(1)
	if v, ok := param["business_id"]; ok {
		if vv := gf.Int64(v); vv > 0 {
			businessID = vv
		}
	}

	user, err := gf.Model("business_user").
		Fields("id,business_id,status,logintime").
		Where("mobile", phone).
		Where("deletetime", nil).
		Find()
	// 未绑定则自动创建
	if err != nil || user == nil {
		addID, addErr := gf.Model("business_user").Data(gf.Map{
			"business_id": businessID,
			"username":    phone,
			"name":        "",
			"nickname":    phone,
			"mobile":      phone,
			"avatar":      "resource/uploads/static/avatar.png",
			"sex":         0,
			"role":        "user",
			"openid":      "",
			"unionid":     "",
			"status":      0,
			"loginip":     gf.GetIp(c),
			"prevtime":    0,
			"logintime":   gtime.Timestamp(),
		}).InsertAndGetId()
		if addErr != nil || addID == 0 {
			gf.Failed().SetMsg("创建用户失败").SetData(addErr).Regin(c)
			return
		}
		user, _ = gf.Model("business_user").Fields("id,business_id,status,logintime").Where("id", addID).Find()
	}

	if user["status"].Int() == 1 {
		gf.Failed().SetMsg("账号被禁用了").Regin(c)
		return
	}

	token, err := gf.CreateToken(gf.Map{
		"ID":          user["id"].Int64(),
		"account_id":  0,
		"business_id": user["business_id"].Int64(),
	})
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}

	prev := user["logintime"].Int64()
	gf.Model("business_user").
		Where("id", user["id"]).
		Data(gf.Map{"prevtime": prev, "logintime": gtime.Timestamp(), "loginip": gf.GetIp(c)}).
		Update()
	gf.AddloginLog(c, gf.Map{"uid": user["id"], "account_id": 0, "business_id": user["business_id"], "type": "uniapp", "status": 0, "des": "小程序微信手机号一键登录"})
	gf.Success().SetMsg("登录成功！").SetData(token).SetToken(gf.String(token)).SetExdata(gf.Map{"mobile": phone, "wx_code": wxCode}).Regin(c)
}

// 《获取用户信息》
func (api *Index) GetUserinfo(c *gf.GinCtx) {
	userID := c.GetInt64("userID")
	userdata, err := gf.Model("business_user").Fields("id,business_id,username,name,nickname,mobile,email,avatar,sex,role,store_id,title,status,createtime,updatetime").Where("id", userID).Where("deletetime", nil).Find()
	if err != nil {
		gf.Failed().SetMsg("查找用户数据错误：" + err.Error()).Regin(c)
		return
	}
	// 门店信息（不让前端展示 store_id）
	storeName := "未绑定"
	storeAddr := ""
	if userdata != nil && userdata["store_id"].Int64() > 0 {
		store, serr := gf.Model("business_stores").
			Fields("id,name,address").
			Where("id", userdata["store_id"].Int64()).
			Where("business_id", userdata["business_id"].Int64()).
			Find()
		if serr == nil && store != nil {
			if store["name"].String() != "" {
				storeName = store["name"].String()
			}
			storeAddr = store["address"].String()
		}
	}
	userdata["store_name"] = gf.VarNew(storeName)
	userdata["store_address"] = gf.VarNew(storeAddr)
	if userdata["avatar"] == nil {
		userdata["avatar"] = gvar.New(gf.GetLocalUrl() + "resource/uploads/static/avatar.png")
	} else if !strings.Contains(userdata["avatar"].String(), "http") {
		userdata["avatar"] = gvar.New(gf.GetFullUrl(userdata["avatar"].String()))
	}
	//处理敏感信息
	userdata["mobile"] = gf.VarNew(gf.HideStrInfo("mobile", userdata["mobile"].String()))
	userdata["email"] = gf.VarNew(gf.HideStrInfo("email", userdata["email"].String()))
	//附件访问完整地址域名
	userdata["rooturls"] = gf.VarNew(gf.GetAllRootUrl()) //全部上传方式访问地址
	userdata["defrooturl"] = gf.VarNew(gf.GetRootUrl())  //设置的上传方式访问地址
	gf.Success().SetMsg("获取用户信息").SetData(userdata).Regin(c)
}

// 《获取我的名片资料》（不脱敏，用于名片预览/编辑）
func (api *Index) GetCardProfile(c *gf.GinCtx) {
	userID := c.GetInt64("userID")
	userdata, err := gf.Model("business_user").
		Fields("id,business_id,username,name,nickname,mobile,email,avatar,role,store_id,title,introduction,status,createtime,updatetime").
		Where("id", userID).
		Where("deletetime", nil).
		Find()
	if err != nil {
		gf.Failed().SetMsg("查找用户数据错误：" + err.Error()).Regin(c)
		return
	}
	if userdata == nil {
		gf.Failed().SetMsg("用户不存在").Regin(c)
		return
	}

	// 门店信息（不让前端展示 store_id）
	storeName := "未绑定"
	storeAddr := ""
	if userdata["store_id"].Int64() > 0 {
		store, serr := gf.Model("business_stores").
			Fields("id,name,address").
			Where("id", userdata["store_id"].Int64()).
			Where("business_id", userdata["business_id"].Int64()).
			Where("deletetime", nil).
			Find()
		if serr == nil && store != nil {
			if store["name"].String() != "" {
				storeName = store["name"].String()
			}
			storeAddr = store["address"].String()
		}
	}
	userdata["store_name"] = gf.VarNew(storeName)
	userdata["store_address"] = gf.VarNew(storeAddr)

	// 头像处理成可访问完整地址
	if userdata["avatar"] == nil || userdata["avatar"].String() == "" {
		userdata["avatar"] = gvar.New(gf.GetLocalUrl() + "resource/uploads/static/avatar.png")
	} else if !strings.Contains(userdata["avatar"].String(), "http") {
		userdata["avatar"] = gvar.New(gf.GetFullUrl(userdata["avatar"].String()))
	}

	gf.Success().SetMsg("获取名片资料").SetData(userdata).Regin(c)
}

// 《获取门店列表》（名片页选择门店）
func (api *Index) GetStores(c *gf.GinCtx) {
	businessID := c.GetInt64("businessID")
	list, err := gf.Model("business_stores").
		Fields("id,name,address").
		Where("business_id", businessID).
		Where("deletetime", nil).
		Where("status", 0).
		Order("weigh desc,id desc").
		Select()
	if err != nil {
		gf.Failed().SetMsg("获取门店失败：" + err.Error()).Regin(c)
		return
	}
	gf.Success().SetMsg("获取门店列表").SetData(list).Regin(c)
}

// 《更新我的名片资料》
// 允许字段：avatar、name、title、introduction、store_id
func (api *Index) UpdateCardProfile(c *gf.GinCtx) {
	userID := c.GetInt64("userID")
	param, _ := gf.RequestParam(c)

	update := gf.Map{}
	if v, ok := param["avatar"]; ok {
		update["avatar"] = gf.String(v)
	}
	if v, ok := param["name"]; ok {
		update["name"] = gf.String(v)
	}
	if v, ok := param["title"]; ok {
		update["title"] = gf.String(v)
	}
	if v, ok := param["introduction"]; ok {
		update["introduction"] = gf.String(v)
	}
	if v, ok := param["store_id"]; ok {
		storeID := gf.Int64(v)
		if storeID == 0 {
			update["store_id"] = 0
		} else {
			store, serr := gf.Model("business_stores").
				Fields("id,status").
				Where("id", storeID).
				Where("business_id", c.GetInt64("businessID")).
				Where("deletetime", nil).
				Find()
			if serr != nil {
				gf.Failed().SetMsg("校验门店失败：" + serr.Error()).Regin(c)
				return
			}
			if store == nil {
				gf.Failed().SetMsg("门店不存在").Regin(c)
				return
			}
			if store["status"].Int() != 0 {
				gf.Failed().SetMsg("门店不可用").Regin(c)
				return
			}
			update["store_id"] = storeID
		}
	}
	if len(update) == 0 {
		gf.Failed().SetMsg("暂无可更新字段").Regin(c)
		return
	}

	if _, err := gf.Model("business_user").Where("id", userID).Update(update); err != nil {
		gf.Failed().SetMsg("保存失败：" + err.Error()).Regin(c)
		return
	}
	gf.Success().SetMsg("保存成功").SetData(true).Regin(c)
}

// 《退出登录》
func (api *Index) Logout(c *gf.GinCtx) {
	gf.RemoveToken(c) //清除token，让当前token失效
	gf.Success().SetMsg("退出登录").SetData(true).Regin(c)
}

// 微信接口端
func getWxappAppidSecret(c *gf.GinCtx) (appid string, secretkey string) {
	// 优先读：resource/config/wxapp.yaml
	if v, _ := gcfg.Instance("wxapp").Get(c, "appid"); v != nil {
		appid = v.String()
	}
	if v, _ := gcfg.Instance("wxapp").Get(c, "secretkey"); v != nil {
		secretkey = v.String()
	}
	if appid != "" && secretkey != "" {
		return
	}
	// 兼容读：resource/config/confdemo.yaml（你当前在用）
	if appid == "" {
		if v, _ := gcfg.Instance("confdemo").Get(c, "data.appid"); v != nil {
			appid = v.String()
		}
	}
	if secretkey == "" {
		if v, _ := gcfg.Instance("confdemo").Get(c, "data.secretkey"); v != nil {
			secretkey = v.String()
		}
	}
	return
}

type wxAccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

func wxGetAccessToken(appid, secret string) (string, error) {
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appid + "&secret=" + secret
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var out wxAccessTokenResp
	if err := json.Unmarshal(body, &out); err != nil {
		return "", err
	}
	if out.ErrCode != 0 || out.AccessToken == "" {
		return "", errors.New(out.ErrMsg)
	}
	return out.AccessToken, nil
}

type wxPhoneResp struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
	} `json:"phone_info"`
}

// wxGetPhoneNumber 通过微信小程序授权码获取用户手机号
// 该函数调用微信接口将临时授权码转换为用户手机号
//
// @param accessToken 微信接口访问令牌
// @param phoneCode 小程序获取手机号授权码
// @return string 用户手机号码，优先返回纯手机号，其次返回完整手机号
// @return error 接口调用错误或微信API返回错误
func wxGetPhoneNumber(accessToken, phoneCode string) (string, error) {
	url := "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=" + accessToken
	payload := map[string]string{"code": phoneCode}
	bs, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewReader(bs))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var out wxPhoneResp
	if err := json.Unmarshal(body, &out); err != nil {
		return "", err
	}
	if out.ErrCode != 0 {
		return "", errors.New(out.ErrMsg)
	}
	if out.PhoneInfo.PurePhoneNumber != "" {
		return out.PhoneInfo.PurePhoneNumber, nil
	}
	return out.PhoneInfo.PhoneNumber, nil
}

type wxMiniCodeErr struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type wxUrlLinkResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	URLLink string `json:"url_link"`
}

// wxGetWxaCodeUnlimit 获取微信小程序的无限码
// 该函数通过微信API获取小程序的二维码，支持自定义场景值和页面路径
//
// @param accessToken 微信接口访问令牌
// @param scene 场景值，用于传递自定义参数
// @param page 小程序页面路径
// @return []byte 二维码图片数据
// @return error 错误信息，获取失败时返回错误
func wxGetWxaCodeUnlimit(accessToken, scene, page string, checkPath bool) ([]byte, error) {
	url := "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=" + accessToken
	payload := map[string]any{
		"scene": scene,
		"page":  page,
		"width": 430,
	}
	// 测试阶段可关闭页面路径校验（避免新增页面未发布导致生成失败）
	if !checkPath {
		payload["check_path"] = false
	}
	bs, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	// 错误返回是 JSON
	if bytes.HasPrefix(bytes.TrimSpace(body), []byte("{")) {
		var out wxMiniCodeErr
		if e := json.Unmarshal(body, &out); e == nil && out.ErrCode != 0 {
			return nil, errors.New(out.ErrMsg)
		}
	}
	return body, nil
}

func wxGenerateUrlLink(accessToken, path, query, envVersion string, expireIntervalSec int) (string, error) {
	url := "https://api.weixin.qq.com/wxa/generate_urllink?access_token=" + accessToken
	payload := map[string]any{
		"path": path,
	}
	if strings.TrimSpace(query) != "" {
		payload["query"] = query
	}
	if strings.TrimSpace(envVersion) != "" {
		payload["env_version"] = envVersion
	}
	if expireIntervalSec > 0 {
		payload["expire_type"] = 1
		payload["expire_interval"] = expireIntervalSec
	}

	bs, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewReader(bs))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var out wxUrlLinkResp
	if err := json.Unmarshal(body, &out); err != nil {
		return "", err
	}
	if out.ErrCode != 0 || out.URLLink == "" {
		return "", errors.New(out.ErrMsg)
	}
	return out.URLLink, nil
}

// 《获取经纪人小程序码》（用于名片二维码）
// 入参：style（可选，名片样式下标，放入 scene 方便客户端复原样式）
func (api *Index) GetAgentWxaCode(c *gf.GinCtx) {
	userID := c.GetInt64("userID")
	style := c.DefaultQuery("style", "0")
	styleIdx, _ := strconv.Atoi(style)
	if styleIdx < 0 {
		styleIdx = 0
	}
	checkPath := true
	checkPathStr := strings.TrimSpace(strings.ToLower(c.DefaultQuery("check_path", "true")))
	if checkPathStr == "false" || checkPathStr == "0" {
		checkPath = false
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

	scene := "a" + strconv.FormatInt(userID, 10) + "_s" + strconv.Itoa(styleIdx)
	page := "pages/agent_public_card/agent_public_card"

	// 文件缓存：同一个 agent+style 复用，减少微信调用次数
	relDir := filepath.Join("resource", "uploads", "wxacode")
	if mkErr := os.MkdirAll(relDir, 0o755); mkErr != nil {
		gf.Failed().SetMsg("创建目录失败：" + mkErr.Error()).Regin(c)
		return
	}
	filename := "agent_" + strconv.FormatInt(userID, 10) + "_s" + strconv.Itoa(styleIdx) + ".png"
	relPath := filepath.Join(relDir, filename)
	if _, statErr := os.Stat(relPath); statErr == nil {
		gf.Success().SetMsg("获取小程序码").SetData(gf.Map{"url": gf.GetFullUrl(filepath.ToSlash(filepath.Join("resource/uploads/wxacode", filename)))}).Regin(c)
		return
	}

	img, err := wxGetWxaCodeUnlimit(accessToken, scene, page, checkPath)
	if err != nil {
		gf.Failed().SetMsg("获取小程序码失败：" + err.Error()).Regin(c)
		return
	}
	if wErr := os.WriteFile(relPath, img, 0o644); wErr != nil {
		gf.Failed().SetMsg("保存小程序码失败：" + wErr.Error()).Regin(c)
		return
	}

	gf.Success().SetMsg("获取小程序码").SetData(gf.Map{"url": gf.GetFullUrl(filepath.ToSlash(filepath.Join("resource/uploads/wxacode", filename)))}).Regin(c)
}

// 《获取经纪人 URL Link》（用于“复制链接”，微信内打开可跳转到小程序落地页）
// 入参：style（可选），env_version（可选：develop|trial|release）
func (api *Index) GetAgentUrlLink(c *gf.GinCtx) {
	userID := c.GetInt64("userID")
	styleIdx := gf.Int(c.DefaultQuery("style", "0"))
	if styleIdx < 0 {
		styleIdx = 0
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

	path := "pages/agent_public_card/agent_public_card"
	query := "agent_id=" + strconv.FormatInt(userID, 10) + "&style=" + strconv.Itoa(styleIdx)
	// 默认 30 天，测试阶段足够；如需永久可调整为更长时间
	urlLink, err := wxGenerateUrlLink(accessToken, path, query, envVersion, 60*60*24*30)
	if err != nil {
		gf.Failed().SetMsg("获取 URLLink 失败：" + err.Error()).Regin(c)
		return
	}
	gf.Success().SetMsg("获取 URLLink").SetData(gf.Map{"url_link": urlLink}).Regin(c)
}

// 《获取经纪人名片公开信息》（给客户扫码落地页使用，无需登录）
// 入参：agent_id
func (api *Index) GetAgentCard(c *gf.GinCtx) {
	agentID := gf.Int64(c.DefaultQuery("agent_id", "0"))
	if agentID <= 0 {
		gf.Failed().SetMsg("缺少参数 agent_id").Regin(c)
		return
	}

	userdata, err := gf.Model("business_user").
		Fields("id,business_id,name,nickname,mobile,avatar,title,introduction,store_id,status").
		Where("id", agentID).
		Where("deletetime", nil).
		Find()
	if err != nil {
		gf.Failed().SetMsg("查找经纪人失败：" + err.Error()).Regin(c)
		return
	}
	if userdata == nil {
		gf.Failed().SetMsg("经纪人不存在").Regin(c)
		return
	}
	if userdata["status"].Int() != 0 {
		gf.Failed().SetMsg("经纪人不可用").Regin(c)
		return
	}

	storeName := "未绑定"
	storeAddr := ""
	if userdata["store_id"].Int64() > 0 {
		store, serr := gf.Model("business_stores").
			Fields("id,name,address,status").
			Where("id", userdata["store_id"].Int64()).
			Where("business_id", userdata["business_id"].Int64()).
			Where("deletetime", nil).
			Find()
		if serr == nil && store != nil && store["status"].Int() == 0 {
			if store["name"].String() != "" {
				storeName = store["name"].String()
			}
			storeAddr = store["address"].String()
		}
	}
	userdata["store_name"] = gf.VarNew(storeName)
	userdata["store_address"] = gf.VarNew(storeAddr)

	// 头像处理成可访问完整地址
	if userdata["avatar"] == nil || userdata["avatar"].String() == "" {
		userdata["avatar"] = gvar.New(gf.GetLocalUrl() + "resource/uploads/static/avatar.png")
	} else if !strings.Contains(userdata["avatar"].String(), "http") {
		userdata["avatar"] = gvar.New(gf.GetFullUrl(userdata["avatar"].String()))
	}

	if userdata["name"].String() == "" {
		userdata["name"] = userdata["nickname"]
	}

	gf.Success().SetMsg("获取经纪人名片").SetData(userdata).Regin(c)
}
