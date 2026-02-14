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

const (
	uniappDefaultTitle = "金牌经纪人"

	// 小程序端登录/注册相关自定义错误码（保持 HTTP 200，code 非 0 便于前端分流处理）
	wxappCodeNotRegistered = 10001 // 未注册/未提交审核资料
	wxappCodeAuditPending  = 10002 // 审核中
	wxappCodeAuditRejected = 10003 // 审核拒绝

	// business_user.audit_status 值
	userAuditStatusApproved = "approved"
	userAuditStatusPending  = "pending"
	userAuditStatusRejected = "rejected"
)

func normalizeUserAuditStatus(v string) string {
	s := strings.ToLower(strings.TrimSpace(v))
	switch s {
	case "", "ok", "pass", "passed", userAuditStatusApproved:
		return userAuditStatusApproved
	case "review", "todo", userAuditStatusPending:
		return userAuditStatusPending
	case "reject", "fail", "denied", userAuditStatusRejected:
		return userAuditStatusRejected
	default:
		return s
	}
}

func init() {
	fpath := Index{NoNeedLogin: []string{"login", "wxLogin", "registerApply", "getRegisterStatus", "getRegisterStores", "logout", "getAgentCard", "getLoginDocs"}, NoNeedAuths: []string{"*"}}
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

	// business_id：可由前端传 business_id；未传则使用请求头 Businessid（再兜底 1）
	businessID := wxBusinessID(c)
	if v, ok := param["business_id"]; ok {
		if vv := gf.Int64(v); vv > 0 {
			businessID = vv
		}
	}

	// 查找用户（手机号即账号）
	fields := "id,business_id,status,logintime"
	if gf.DbHaseField("business_user", "audit_status") {
		fields += ",audit_status"
	}
	if gf.DbHaseField("business_user", "audit_reason") {
		fields += ",audit_reason"
	}
	user, err := gf.Model("business_user").
		Fields(fields).
		Where("mobile", mobile).
		Where("deletetime", nil).
		Find()
	if err != nil || user == nil {
		gf.Failed().SetCode(wxappCodeNotRegistered).SetMsg("账号未注册，请先完善资料提交审核").SetData(gf.Map{
			"mobile":      mobile,
			"business_id": businessID,
		}).Regin(c)
		return
	}

	if user["status"].Int() == 1 {
		gf.Failed().SetMsg("账号被禁用了").Regin(c)
		return
	}

	// 审核状态校验：未通过审核不允许登录
	if gf.DbHaseField("business_user", "audit_status") {
		auditStatus := normalizeUserAuditStatus(gf.String(user["audit_status"]))
		auditReason := ""
		if gf.DbHaseField("business_user", "audit_reason") {
			auditReason = strings.TrimSpace(gf.String(user["audit_reason"]))
		}
		switch auditStatus {
		case userAuditStatusPending:
			gf.Failed().SetCode(wxappCodeAuditPending).SetMsg("资料审核中，请等待管理员审核").SetData(gf.Map{
				"mobile":       mobile,
				"audit_status": auditStatus,
			}).Regin(c)
			return
		case userAuditStatusRejected:
			msg := "审核未通过，请重新提交资料"
			if auditReason != "" {
				msg = msg + "：" + auditReason
			}
			gf.Failed().SetCode(wxappCodeAuditRejected).SetMsg(msg).SetData(gf.Map{
				"mobile":       mobile,
				"audit_status": auditStatus,
				"audit_reason": auditReason,
			}).Regin(c)
			return
		}
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

	// business_id：可由前端传 business_id；未传则使用请求头 Businessid（再兜底 1）
	businessID := wxBusinessID(c)
	if v, ok := param["business_id"]; ok {
		if vv := gf.Int64(v); vv > 0 {
			businessID = vv
		}
	}

	fields := "id,business_id,status,logintime"
	if gf.DbHaseField("business_user", "audit_status") {
		fields += ",audit_status"
	}
	if gf.DbHaseField("business_user", "audit_reason") {
		fields += ",audit_reason"
	}
	user, err := gf.Model("business_user").
		Fields(fields).
		Where("mobile", phone).
		Where("deletetime", nil).
		Find()
	if err != nil || user == nil {
		gf.Failed().SetCode(wxappCodeNotRegistered).SetMsg("账号未注册，请先完善资料提交审核").SetData(gf.Map{
			"mobile":      phone,
			"business_id": businessID,
		}).SetExdata(gf.Map{"mobile": phone, "wx_code": wxCode}).Regin(c)
		return
	}

	if user["status"].Int() == 1 {
		gf.Failed().SetMsg("账号被禁用了").Regin(c)
		return
	}

	// 审核状态校验：未通过审核不允许登录
	if gf.DbHaseField("business_user", "audit_status") {
		auditStatus := normalizeUserAuditStatus(gf.String(user["audit_status"]))
		auditReason := ""
		if gf.DbHaseField("business_user", "audit_reason") {
			auditReason = strings.TrimSpace(gf.String(user["audit_reason"]))
		}
		switch auditStatus {
		case userAuditStatusPending:
			gf.Failed().SetCode(wxappCodeAuditPending).SetMsg("资料审核中，请等待管理员审核").SetData(gf.Map{
				"mobile":       phone,
				"audit_status": auditStatus,
			}).SetExdata(gf.Map{"mobile": phone, "wx_code": wxCode}).Regin(c)
			return
		case userAuditStatusRejected:
			msg := "审核未通过，请重新提交资料"
			if auditReason != "" {
				msg = msg + "：" + auditReason
			}
			gf.Failed().SetCode(wxappCodeAuditRejected).SetMsg(msg).SetData(gf.Map{
				"mobile":       phone,
				"audit_status": auditStatus,
				"audit_reason": auditReason,
			}).SetExdata(gf.Map{"mobile": phone, "wx_code": wxCode}).Regin(c)
			return
		}
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

// 《提交注册/完善资料申请》（无需登录，审核通过后才可登录）
// POST /uniapp/registerApply
// 入参：
// - phone_code: getPhoneNumber 返回的 e.detail.code（推荐）
// - mobile + captcha: 备用（非小程序环境）
// - name: 真实姓名（必填）
// - store_id: 选择门店ID（可选）
// - store_name_text/store_address_text: 手填门店信息（store_id 为空时必填 store_name_text）
// - region_province/region_city/region_district: 所在地区（可选）
func (api *Index) RegisterApply(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)

	name := strings.TrimSpace(gf.String(param["name"]))
	if name == "" {
		gf.Failed().SetMsg("请填写真实姓名").Regin(c)
		return
	}

	// business_id：可由前端传 business_id；未传则使用请求头 Businessid（再兜底 1）
	businessID := wxBusinessID(c)
	if v, ok := param["business_id"]; ok {
		if vv := gf.Int64(v); vv > 0 {
			businessID = vv
		}
	}

	phoneCode := strings.TrimSpace(gf.String(param["phone_code"]))
	mobile := strings.TrimSpace(gf.String(param["mobile"]))
	if phoneCode != "" {
		appid, secretkey := getWxappAppidSecret(c)
		if appid == "" || secretkey == "" {
			gf.Failed().SetMsg("未配置微信小程序 appid/secretkey，无法提交审核").Regin(c)
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
		mobile = strings.TrimSpace(phone)
	}
	if mobile == "" {
		gf.Failed().SetMsg("请先授权手机号或填写手机号").Regin(c)
		return
	}
	// 非小程序环境下的验证码校验（可选）
	if phoneCode == "" {
		code, emerr := gf.GetVerifyCode(mobile)
		if emerr != nil || code != gf.Int(param["captcha"]) {
			gf.Failed().SetMsg("验证码无效").SetData(emerr).Regin(c)
			return
		}
	}

	storeID := gf.Int64(param["store_id"])
	storeNameText := strings.TrimSpace(gf.String(param["store_name_text"]))
	if storeNameText == "" {
		storeNameText = strings.TrimSpace(gf.String(param["store_name"]))
	}
	storeAddrText := strings.TrimSpace(gf.String(param["store_address_text"]))
	if storeAddrText == "" {
		storeAddrText = strings.TrimSpace(gf.String(param["store_address"]))
	}

	if storeID <= 0 && storeNameText == "" {
		gf.Failed().SetMsg("请选择门店或填写门店名称").Regin(c)
		return
	}

	// 若选择门店，校验门店有效性
	if storeID > 0 {
		store, serr := gf.Model("business_stores").
			Fields("id,status").
			Where("id", storeID).
			Where("business_id", businessID).
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
	}

	province := strings.TrimSpace(gf.String(param["region_province"]))
	city := strings.TrimSpace(gf.String(param["region_city"]))
	district := strings.TrimSpace(gf.String(param["region_district"]))

	now := time.Now()

	// 查找用户（手机号即账号）
	findFields := "id,business_id,status"
	if gf.DbHaseField("business_user", "audit_status") {
		findFields += ",audit_status"
	}
	user, err := gf.Model("business_user").
		Fields(findFields).
		Where("mobile", mobile).
		Where("deletetime", nil).
		Find()
	if err != nil {
		gf.Failed().SetMsg("查找用户失败：" + err.Error()).Regin(c)
		return
	}

	// 仅在字段存在时写入，避免老库缺字段直接报错
	buildAuditPending := func(m gf.Map) {
		if gf.DbHaseField("business_user", "audit_status") {
			m["audit_status"] = userAuditStatusPending
		}
		if gf.DbHaseField("business_user", "audit_reason") {
			m["audit_reason"] = ""
		}
		if gf.DbHaseField("business_user", "apply_time") {
			m["apply_time"] = now
		}
		if gf.DbHaseField("business_user", "audit_time") {
			m["audit_time"] = nil
		}
	}

	if user == nil {
		data := gf.Map{
			"business_id": businessID,
			"username":    mobile,
			"name":        name,
			"nickname":    mobile,
			"mobile":      mobile,
			"avatar":      "resource/uploads/static/avatar.png",
			"sex":         0,
			"role":        "user",
			"title":       uniappDefaultTitle,
			"status":      0,
		}
		// 门店与地区（可选）
		data["store_id"] = storeID
		if gf.DbHaseField("business_user", "store_name_text") {
			data["store_name_text"] = storeNameText
		}
		if gf.DbHaseField("business_user", "store_address_text") {
			data["store_address_text"] = storeAddrText
		}
		if gf.DbHaseField("business_user", "region_province") {
			data["region_province"] = province
		}
		if gf.DbHaseField("business_user", "region_city") {
			data["region_city"] = city
		}
		if gf.DbHaseField("business_user", "region_district") {
			data["region_district"] = district
		}
		buildAuditPending(data)

		addID, addErr := gf.Model("business_user").Data(data).InsertAndGetId()
		if addErr != nil || addID == 0 {
			gf.Failed().SetMsg("提交失败，请稍后再试").SetData(addErr).Regin(c)
			return
		}
		gf.Success().SetMsg("提交成功，请等待审核").SetData(gf.Map{
			"id":           addID,
			"mobile":       mobile,
			"audit_status": userAuditStatusPending,
		}).Regin(c)
		return
	}

	// 已存在用户：若未通过审核，则重置为 pending；若已通过审核，仅更新资料
	update := gf.Map{
		"name":     name,
		"store_id": storeID,
	}
	if gf.DbHaseField("business_user", "store_name_text") {
		update["store_name_text"] = storeNameText
	}
	if gf.DbHaseField("business_user", "store_address_text") {
		update["store_address_text"] = storeAddrText
	}
	if gf.DbHaseField("business_user", "region_province") {
		update["region_province"] = province
	}
	if gf.DbHaseField("business_user", "region_city") {
		update["region_city"] = city
	}
	if gf.DbHaseField("business_user", "region_district") {
		update["region_district"] = district
	}

	needPending := true
	if gf.DbHaseField("business_user", "audit_status") {
		cur := normalizeUserAuditStatus(gf.String(user["audit_status"]))
		needPending = cur != userAuditStatusApproved
	}
	if needPending {
		buildAuditPending(update)
	}

	if _, uerr := gf.Model("business_user").Where("id", user["id"]).Update(update); uerr != nil {
		gf.Failed().SetMsg("提交失败：" + uerr.Error()).Regin(c)
		return
	}

	// 返回最新审核状态
	auditStatus := userAuditStatusPending
	if gf.DbHaseField("business_user", "audit_status") {
		auditStatus = userAuditStatusPending
		if !needPending {
			auditStatus = userAuditStatusApproved
		}
	}
	gf.Success().SetMsg("提交成功，请等待审核").SetData(gf.Map{
		"id":           user["id"].Int64(),
		"mobile":       mobile,
		"audit_status": auditStatus,
	}).Regin(c)
}

// 《获取注册/审核状态》（无需登录）
// GET /uniapp/getRegisterStatus
// 入参：mobile 或 phone_code
func (api *Index) GetRegisterStatus(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)

	phoneCode := strings.TrimSpace(gf.String(param["phone_code"]))
	mobile := strings.TrimSpace(gf.String(param["mobile"]))
	if phoneCode != "" {
		appid, secretkey := getWxappAppidSecret(c)
		if appid == "" || secretkey == "" {
			gf.Failed().SetMsg("未配置微信小程序 appid/secretkey，无法查询审核状态").Regin(c)
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
		mobile = strings.TrimSpace(phone)
	}
	if mobile == "" {
		gf.Failed().SetMsg("请提交手机号").Regin(c)
		return
	}

	fields := "id,business_id,name,mobile,store_id,status"
	if gf.DbHaseField("business_user", "audit_status") {
		fields += ",audit_status"
	}
	if gf.DbHaseField("business_user", "audit_reason") {
		fields += ",audit_reason"
	}
	if gf.DbHaseField("business_user", "apply_time") {
		fields += ",apply_time"
	}
	if gf.DbHaseField("business_user", "audit_time") {
		fields += ",audit_time"
	}
	if gf.DbHaseField("business_user", "store_name_text") {
		fields += ",store_name_text"
	}
	if gf.DbHaseField("business_user", "store_address_text") {
		fields += ",store_address_text"
	}

	user, err := gf.Model("business_user").
		Fields(fields).
		Where("mobile", mobile).
		Where("deletetime", nil).
		Find()
	if err != nil {
		gf.Failed().SetMsg("获取审核状态失败：" + err.Error()).Regin(c)
		return
	}
	if user == nil {
		gf.Failed().SetCode(wxappCodeNotRegistered).SetMsg("账号未注册").SetData(gf.Map{"mobile": mobile}).Regin(c)
		return
	}

	// 门店信息：优先 store_id，其次手填字段
	storeName := "未绑定"
	storeAddr := ""
	if user["store_id"].Int64() > 0 {
		store, serr := gf.Model("business_stores").
			Fields("id,name,address").
			Where("id", user["store_id"].Int64()).
			Where("business_id", user["business_id"].Int64()).
			Where("deletetime", nil).
			Find()
		if serr == nil && store != nil {
			if store["name"].String() != "" {
				storeName = store["name"].String()
			}
			storeAddr = store["address"].String()
		}
	} else {
		if gf.DbHaseField("business_user", "store_name_text") {
			if s := strings.TrimSpace(gf.String(user["store_name_text"])); s != "" {
				storeName = s
			}
		}
		if gf.DbHaseField("business_user", "store_address_text") {
			storeAddr = strings.TrimSpace(gf.String(user["store_address_text"]))
		}
	}

	auditStatus := userAuditStatusApproved
	if gf.DbHaseField("business_user", "audit_status") {
		auditStatus = normalizeUserAuditStatus(gf.String(user["audit_status"]))
	}
	auditReason := ""
	if gf.DbHaseField("business_user", "audit_reason") {
		auditReason = strings.TrimSpace(gf.String(user["audit_reason"]))
	}

	gf.Success().SetMsg("获取审核状态").SetData(gf.Map{
		"id":           user["id"].Int64(),
		"business_id":   user["business_id"].Int64(),
		"name":         user["name"].String(),
		"mobile":       mobile,
		"store_id":     user["store_id"].Int64(),
		"store_name":   storeName,
		"store_address": storeAddr,
		"status":       user["status"].Int(),
		"audit_status": auditStatus,
		"audit_reason": auditReason,
		"apply_time":   user["apply_time"],
		"audit_time":   user["audit_time"],
		"can_login":    auditStatus == userAuditStatusApproved && user["status"].Int() == 0,
	}).Regin(c)
}

// 《注册页门店列表》（无需登录）
// GET /uniapp/getRegisterStores
func (api *Index) GetRegisterStores(c *gf.GinCtx) {
	businessID := wxBusinessID(c)
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

// 《获取用户信息》
func (api *Index) GetUserinfo(c *gf.GinCtx) {
	userID := c.GetInt64("userID")
	fields := "id,business_id,username,name,nickname,mobile,email,avatar,sex,role,store_id,title,can_manage_properties,can_manage_locks,status,createtime,updatetime"
	if gf.DbHaseField("business_user", "store_name_text") {
		fields += ",store_name_text"
	}
	if gf.DbHaseField("business_user", "store_address_text") {
		fields += ",store_address_text"
	}
	userdata, err := gf.Model("business_user").Fields(fields).Where("id", userID).Where("deletetime", nil).Find()
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
	} else {
		// 手填门店兜底
		if gf.DbHaseField("business_user", "store_name_text") {
			if s := strings.TrimSpace(gf.String(userdata["store_name_text"])); s != "" {
				storeName = s
			}
		}
		if gf.DbHaseField("business_user", "store_address_text") {
			storeAddr = strings.TrimSpace(gf.String(userdata["store_address_text"]))
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
	fields := "id,business_id,username,name,nickname,mobile,email,avatar,role,store_id,title,introduction,status,createtime,updatetime"
	if gf.DbHaseField("business_user", "store_name_text") {
		fields += ",store_name_text"
	}
	if gf.DbHaseField("business_user", "store_address_text") {
		fields += ",store_address_text"
	}
	userdata, err := gf.Model("business_user").
		Fields(fields).
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
	} else {
		// 手填门店兜底
		if gf.DbHaseField("business_user", "store_name_text") {
			if s := strings.TrimSpace(gf.String(userdata["store_name_text"])); s != "" {
				storeName = s
			}
		}
		if gf.DbHaseField("business_user", "store_address_text") {
			storeAddr = strings.TrimSpace(gf.String(userdata["store_address_text"]))
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

	// 名片推荐房源：经纪人推荐 + 系统推荐（可并存）
	businessID := userdata["business_id"].Int64()
	agentRecommendedProperties, selectedIDs, recErr := queryAgentRecommendedProperties(businessID, agentID, agentCardRecommendShowSize)
	if recErr != nil {
		agentRecommendedProperties = []gf.Map{}
		selectedIDs = []int64{}
	}
	exclude := make(map[int64]struct{}, len(selectedIDs))
	for _, id := range selectedIDs {
		exclude[id] = struct{}{}
	}
	systemRecommendedProperties, sysErr := querySystemRecommendedProperties(businessID, agentCardRecommendShowSize, exclude)
	if sysErr != nil {
		systemRecommendedProperties = []gf.Map{}
	}

	userdata["agent_recommended_properties"] = gf.VarNew(agentRecommendedProperties)
	userdata["system_recommended_properties"] = gf.VarNew(systemRecommendedProperties)
	userdata["agent_recommend_count"] = gf.VarNew(len(agentRecommendedProperties))
	userdata["system_recommend_count"] = gf.VarNew(len(systemRecommendedProperties))

	gf.Success().SetMsg("获取经纪人名片").SetData(userdata).Regin(c)
}
