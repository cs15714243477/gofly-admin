package uniapp

import (
	"strings"
	"sync"

	"gofly/utils/gf"
	"gofly/utils/tools/gcfg"
	"gofly/utils/tools/gconv"
	"gofly/utils/tools/gctx"
)

var (
	wxPreSaleViewKeySecretOnce sync.Once
	wxPreSaleViewKeySecret     string
)

func wxGetPreSaleViewKeySecret() string {
	wxPreSaleViewKeySecretOnce.Do(func() {
		ctx := gctx.New()
		appConf, err := gcfg.Instance().Get(ctx, "app")
		if err == nil {
			m := gconv.Map(appConf)
			wxPreSaleViewKeySecret = strings.TrimSpace(gconv.String(m["tokensecret"]))
			if wxPreSaleViewKeySecret == "" {
				wxPreSaleViewKeySecret = strings.TrimSpace(gconv.String(m["apisecret"]))
			}
		}
		if wxPreSaleViewKeySecret == "" {
			wxPreSaleViewKeySecret = "gofly_presale"
		}
	})
	return wxPreSaleViewKeySecret
}

func wxBuildPreSaleViewKey(businessID, userID, propertyID int64) string {
	secret := wxGetPreSaleViewKeySecret()
	raw := secret + "#presale#" + gconv.String(businessID) + "#" + gconv.String(userID) + "#" + gconv.String(propertyID)
	return gf.Md5Str(raw)
}

func wxCheckPreSaleViewKey(businessID, userID, propertyID int64, viewKey string) bool {
	key := strings.TrimSpace(viewKey)
	if key == "" {
		return false
	}
	expect := wxBuildPreSaleViewKey(businessID, userID, propertyID)
	return strings.EqualFold(expect, key)
}
