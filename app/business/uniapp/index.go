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
	"strings"
	"time"
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
	fpath := Index{NoNeedLogin: []string{"login", "wxLogin", "logout"}, NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

// 《登录》（手机号验证码登录）
// 入参：mobile、captcha
// 说明：当前沿用 business_account 作为小程序端账号表（与 GetUserinfo 保持一致）
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

	data, err := gf.Model("business_account").
		Fields("id,account_id,business_id,status,lock_time").
		Where("mobile", mobile).
		Find()
	if data == nil || err != nil {
		gf.Failed().SetMsg("手机账号不存在！").Regin(c)
		return
	}
	if data["status"].Int() == 1 {
		gf.Failed().SetMsg("账号被禁用了").Regin(c)
		return
	}
	if time.Now().Before(data["lock_time"].Time()) {
		gf.Failed().SetMsg("账户已被锁定，请稍后再试").Regin(c)
		return
	}
	token, err := gf.CreateToken(gf.Map{
		"ID":          data["id"].Int64(),
		"account_id":  data["account_id"].Int64(),
		"business_id": data["business_id"].Int64(),
	})
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	gf.Model("business_account").
		Where("id", data["id"]).
		Data(gf.Map{"loginstatus": 1, "last_login_time": gtime.Timestamp(), "last_login_ip": gf.GetIp(c)}).
		Update()
	gf.AddloginLog(c, gf.Map{"uid": data["id"], "account_id": data["account_id"], "business_id": data["business_id"], "type": "business", "status": 0, "des": "小程序手机号验证码登录"})
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

	data, err := gf.Model("business_account").
		Fields("id,account_id,business_id,status,lock_time").
		Where("mobile", phone).
		Find()
	if err != nil || data == nil {
		gf.Failed().SetMsg("手机号未绑定账号").Regin(c)
		return
	}
	if data["status"].Int() == 1 {
		gf.Failed().SetMsg("账号被禁用了").Regin(c)
		return
	}
	if time.Now().Before(data["lock_time"].Time()) {
		gf.Failed().SetMsg("账户已被锁定，请稍后再试").Regin(c)
		return
	}

	token, err := gf.CreateToken(gf.Map{
		"ID":          data["id"].Int64(),
		"account_id":  data["account_id"].Int64(),
		"business_id": data["business_id"].Int64(),
	})
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	gf.Model("business_account").
		Where("id", data["id"]).
		Data(gf.Map{"loginstatus": 1, "last_login_time": gtime.Timestamp(), "last_login_ip": gf.GetIp(c)}).
		Update()
	gf.AddloginLog(c, gf.Map{"uid": data["id"], "account_id": data["account_id"], "business_id": data["business_id"], "type": "business", "status": 0, "des": "小程序微信手机号一键登录"})
	gf.Success().SetMsg("登录成功！").SetData(token).SetToken(gf.String(token)).SetExdata(gf.Map{"mobile": phone, "wx_code": wxCode}).Regin(c)
}

// 《获取用户信息》
func (api *Index) GetUserinfo(c *gf.GinCtx) {
	userID := c.GetInt64("userID")
	userdata, err := gf.Model("business_account").Fields("id,business_id,name,nickname,mobile,email,avatar,status,createtime,pwd_reset_time").Where("id", userID).Find()
	if err != nil {
		gf.Failed().SetMsg("查找用户数据错误：" + err.Error()).Regin(c)
		return
	}
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

// 《退出登录》
func (api *Index) Logout(c *gf.GinCtx) {
	user, err := gf.ParseTokenNoValid(c) //当前用户
	if err == nil {
		gf.Model("business_account").Where("id", user.ID).Data(gf.Map{"loginstatus": 0}).Update()
	}
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
