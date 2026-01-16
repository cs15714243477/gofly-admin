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
