package uniapp

import (
	"strings"

	"gofly/utils/gf"
	"gofly/utils/tools/gconv"
)

func wxDocItem(key, title, content, url string) gf.Map {
	return gf.Map{
		"key":     key,
		"title":   strings.TrimSpace(title),
		"content": strings.TrimSpace(content),
		"url":     strings.TrimSpace(url),
	}
}

// GetLoginDocs 获取登录页文档配置（用户协议/隐私政策/帮助中心）
//
// 配置文件：
// - resource/config/wxapp_login_docs.yaml
// - conftype=configuration 且 status=true 时生效
func (api *Index) GetLoginDocs(c *gf.GinCtx) {
	userAgreementTitle := "用户协议"
	userAgreementContent := ""
	userAgreementURL := ""
	privacyPolicyTitle := "隐私政策"
	privacyPolicyContent := ""
	privacyPolicyURL := ""
	helpCenterTitle := "帮助中心"
	helpCenterContent := ""
	helpCenterURL := ""

	if conf, err := gf.GetConfByFile("wxapp_login_docs"); err == nil && conf != nil {
		if m, ok := conf.(map[string]interface{}); ok {
			if gconv.String(m["conftype"]) == "configuration" && gf.Bool(m["status"]) {
				if data, ok := m["data"].(map[string]interface{}); ok {
					if v := strings.TrimSpace(gconv.String(data["user_agreement_title"])); v != "" {
						userAgreementTitle = v
					}
					userAgreementContent = strings.TrimSpace(gconv.String(data["user_agreement_content"]))
					userAgreementURL = strings.TrimSpace(gconv.String(data["user_agreement_url"]))

					if v := strings.TrimSpace(gconv.String(data["privacy_policy_title"])); v != "" {
						privacyPolicyTitle = v
					}
					privacyPolicyContent = strings.TrimSpace(gconv.String(data["privacy_policy_content"]))
					privacyPolicyURL = strings.TrimSpace(gconv.String(data["privacy_policy_url"]))

					if v := strings.TrimSpace(gconv.String(data["help_center_title"])); v != "" {
						helpCenterTitle = v
					}
					helpCenterContent = strings.TrimSpace(gconv.String(data["help_center_content"]))
					helpCenterURL = strings.TrimSpace(gconv.String(data["help_center_url"]))
				}
			}
		}
	}

	gf.Success().SetMsg("获取登录文档配置").SetData(gf.Map{
		"docs": []gf.Map{
			wxDocItem("user_agreement", userAgreementTitle, userAgreementContent, userAgreementURL),
			wxDocItem("privacy_policy", privacyPolicyTitle, privacyPolicyContent, privacyPolicyURL),
			wxDocItem("help_center", helpCenterTitle, helpCenterContent, helpCenterURL),
		},
	}).Regin(c)
}
