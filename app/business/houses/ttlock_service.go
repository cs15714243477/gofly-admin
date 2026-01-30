package houses

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gofly/utils/gf"
	"gofly/utils/tools/gcfg"
)

type ttlockCloudConfig struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
	APIBase      string
}

func getTTLockCloudConfig(c *gf.GinCtx) ttlockCloudConfig {
	cfg := ttlockCloudConfig{
		APIBase: "https://api.ttlock.com",
	}
	if v, _ := gcfg.Instance("confdemo").Get(c, "data.ttlock_client_id"); v != nil {
		cfg.ClientID = strings.TrimSpace(v.String())
	}
	if v, _ := gcfg.Instance("confdemo").Get(c, "data.ttlock_client_secret"); v != nil {
		cfg.ClientSecret = strings.TrimSpace(v.String())
	}
	if v, _ := gcfg.Instance("confdemo").Get(c, "data.ttlock_username"); v != nil {
		cfg.Username = strings.TrimSpace(v.String())
	}
	if v, _ := gcfg.Instance("confdemo").Get(c, "data.ttlock_password"); v != nil {
		cfg.Password = strings.TrimSpace(v.String())
	}
	if v, _ := gcfg.Instance("confdemo").Get(c, "data.ttlock_api_base"); v != nil {
		if s := strings.TrimSpace(v.String()); s != "" {
			cfg.APIBase = strings.TrimRight(s, "/")
		}
	}
	return cfg
}

type ttlockTokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}

func ttlockAPIURL(apiBase, path string) string {
	return strings.TrimRight(apiBase, "/") + path
}

func ttlockTokenURL(apiBase string) string {
	// TTLock cloud oauth2 endpoint
	return strings.TrimRight(apiBase, "/") + "/oauth2/token"
}

func ttlockPostForm(u string, data url.Values) ([]byte, error) {
	resp, err := http.PostForm(u, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, errors.New("http status " + strconv.Itoa(resp.StatusCode))
	}
	return bs, nil
}

var ttlockMD5Re = regexp.MustCompile(`^[a-fA-F0-9]{32}$`)

func ttlockMD5HexLower(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func ttlockGetAccessTokenByPassword(cfg ttlockCloudConfig) (ttlockTokenResp, error) {
	if cfg.ClientID == "" || cfg.ClientSecret == "" || cfg.Username == "" || cfg.Password == "" {
		return ttlockTokenResp{}, errors.New("TTLock 配置不完整（client_id/client_secret/username/password）")
	}
	password := strings.TrimSpace(cfg.Password)
	// TTLock OAuth2 password 授权要求 password 为 MD5(32位) 加密串
	if !ttlockMD5Re.MatchString(password) {
		password = ttlockMD5HexLower(password)
	}
	form := url.Values{}
	form.Set("client_id", cfg.ClientID)
	form.Set("client_secret", cfg.ClientSecret)
	form.Set("username", cfg.Username)
	form.Set("password", password)
	form.Set("grant_type", "password")

	bs, err := ttlockPostForm(ttlockTokenURL(cfg.APIBase), form)
	if err != nil {
		return ttlockTokenResp{}, err
	}
	var out ttlockTokenResp
	if e := json.Unmarshal(bs, &out); e != nil {
		return ttlockTokenResp{}, e
	}
	if out.ErrCode != 0 || out.AccessToken == "" {
		return out, errors.New(out.ErrMsg)
	}
	return out, nil
}

func ttlockRefreshAccessToken(cfg ttlockCloudConfig, refreshToken string) (ttlockTokenResp, error) {
	if cfg.ClientID == "" || cfg.ClientSecret == "" || refreshToken == "" {
		return ttlockTokenResp{}, errors.New("缺少 client_id/client_secret/refresh_token")
	}
	form := url.Values{}
	form.Set("client_id", cfg.ClientID)
	form.Set("client_secret", cfg.ClientSecret)
	form.Set("refresh_token", refreshToken)
	form.Set("grant_type", "refresh_token")

	bs, err := ttlockPostForm(ttlockTokenURL(cfg.APIBase), form)
	if err != nil {
		return ttlockTokenResp{}, err
	}
	var out ttlockTokenResp
	if e := json.Unmarshal(bs, &out); e != nil {
		return ttlockTokenResp{}, e
	}
	if out.ErrCode != 0 || out.AccessToken == "" {
		return out, errors.New(out.ErrMsg)
	}
	return out, nil
}

type ttlockAccountRecord struct {
	ID           int64
	BusinessID   int64
	AccessToken  string
	RefreshToken string
	ExpireAt     int64
}

func ttlockLoadAccount(businessID int64) (*ttlockAccountRecord, error) {
	row, err := gf.Model("business_ttlock_accounts").
		Fields("id,business_id,access_token,refresh_token,expire_at").
		Where("business_id", businessID).
		Where("deletetime", 0).
		Find()
	if err != nil || row == nil {
		return nil, err
	}
	return &ttlockAccountRecord{
		ID:           row["id"].Int64(),
		BusinessID:   row["business_id"].Int64(),
		AccessToken:  row["access_token"].String(),
		RefreshToken: row["refresh_token"].String(),
		ExpireAt:     row["expire_at"].Int64(),
	}, nil
}

func ttlockUpsertAccount(businessID int64, token ttlockTokenResp) error {
	expireAt := time.Now().Unix() + int64(token.ExpiresIn) - 60
	data := gf.Map{
		"business_id":   businessID,
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"expire_at":     expireAt,
		"updatetime":    time.Now().Unix(),
		"status":        0,
		"deletetime":    0,
	}
	exist, _ := gf.Model("business_ttlock_accounts").Where("business_id", businessID).Where("deletetime", 0).Value("id")
	if exist == nil {
		data["createtime"] = time.Now().Unix()
		_, err := gf.Model("business_ttlock_accounts").Data(data).Insert()
		return err
	}
	_, err := gf.Model("business_ttlock_accounts").Where("id", exist).Update(data)
	return err
}

func ttlockEnsureAccessToken(c *gf.GinCtx, businessID int64) (string, error) {
	cfg := getTTLockCloudConfig(c)
	acc, _ := ttlockLoadAccount(businessID)
	now := time.Now().Unix()
	if acc != nil && acc.AccessToken != "" && acc.ExpireAt > now {
		return acc.AccessToken, nil
	}
	// refresh first
	if acc != nil && acc.RefreshToken != "" {
		if tok, err := ttlockRefreshAccessToken(cfg, acc.RefreshToken); err == nil {
			_ = ttlockUpsertAccount(businessID, tok)
			return tok.AccessToken, nil
		}
	}
	// password grant (needs username/password in config)
	tok, err := ttlockGetAccessTokenByPassword(cfg)
	if err != nil {
		return "", err
	}
	_ = ttlockUpsertAccount(businessID, tok)
	return tok.AccessToken, nil
}

type ttlockListLocksResp struct {
	List     []map[string]any `json:"list"`
	Total    int              `json:"total"`
	PageNo   int              `json:"pageNo"`
	PageSize int              `json:"pageSize"`
	ErrCode  int              `json:"errcode"`
	ErrMsg   string           `json:"errmsg"`
}

func ttlockListLocks(cfg ttlockCloudConfig, accessToken string, pageNo, pageSize int) (ttlockListLocksResp, error) {
	if cfg.ClientID == "" || accessToken == "" {
		return ttlockListLocksResp{}, errors.New("缺少 client_id/access_token")
	}
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	u, _ := url.Parse(ttlockAPIURL(cfg.APIBase, "/v3/lock/list"))
	q := u.Query()
	q.Set("clientId", cfg.ClientID)
	q.Set("accessToken", accessToken)
	q.Set("pageNo", strconv.Itoa(pageNo))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("date", strconv.FormatInt(time.Now().UnixMilli(), 10))
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return ttlockListLocksResp{}, err
	}
	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return ttlockListLocksResp{}, errors.New("http status " + strconv.Itoa(resp.StatusCode))
	}
	var out ttlockListLocksResp
	if e := json.Unmarshal(bs, &out); e != nil {
		return ttlockListLocksResp{}, e
	}
	if out.ErrCode != 0 {
		return out, errors.New(out.ErrMsg)
	}
	if out.List == nil {
		out.List = []map[string]any{}
	}
	return out, nil
}

func ttlockNowMilli() int64 { return time.Now().UnixMilli() }

func ttlockParseErr(out map[string]any) (int, string) {
	if out == nil {
		return 0, ""
	}
	var errCode int
	switch v := out["errcode"].(type) {
	case float64:
		errCode = int(v)
	case int:
		errCode = v
	case int64:
		errCode = int(v)
	case string:
		if n, _ := strconv.Atoi(v); n != 0 {
			errCode = n
		}
	}
	errMsg := ""
	if v, ok := out["errmsg"]; ok {
		errMsg = fmt.Sprint(v)
	}
	if errMsg == "" {
		if v, ok := out["description"]; ok {
			errMsg = fmt.Sprint(v)
		}
	}
	return errCode, errMsg
}

func ttlockPostFormJSON(u string, data url.Values) (map[string]any, error) {
	resp, err := http.PostForm(u, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	if len(bs) == 0 {
		if resp.StatusCode >= 400 {
			return nil, errors.New("http status " + strconv.Itoa(resp.StatusCode))
		}
		return map[string]any{}, nil
	}
	var out map[string]any
	if e := json.Unmarshal(bs, &out); e != nil {
		if resp.StatusCode >= 400 {
			return nil, errors.New("http status " + strconv.Itoa(resp.StatusCode))
		}
		return nil, e
	}
	if resp.StatusCode >= 400 {
		if code, msg := ttlockParseErr(out); code != 0 {
			return out, errors.New(msg)
		}
		return out, errors.New("http status " + strconv.Itoa(resp.StatusCode))
	}
	if code, msg := ttlockParseErr(out); code != 0 {
		if msg == "" {
			msg = "TTLock 请求失败"
		}
		return out, errors.New(msg)
	}
	return out, nil
}

type ttlockKeyListResp struct {
	List     []map[string]any `json:"list"`
	Total    int              `json:"total"`
	PageNo   int              `json:"pageNo"`
	PageSize int              `json:"pageSize"`
	ErrCode  int              `json:"errcode"`
	ErrMsg   string           `json:"errmsg"`
}

func ttlockListKeys(cfg ttlockCloudConfig, accessToken string, pageNo, pageSize, sdkVersion int) (ttlockKeyListResp, error) {
	if cfg.ClientID == "" || accessToken == "" {
		return ttlockKeyListResp{}, errors.New("缺少 client_id/access_token")
	}
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}
	u, _ := url.Parse(ttlockAPIURL(cfg.APIBase, "/v3/key/list"))
	q := u.Query()
	q.Set("clientId", cfg.ClientID)
	q.Set("accessToken", accessToken)
	q.Set("pageNo", strconv.Itoa(pageNo))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("date", strconv.FormatInt(ttlockNowMilli(), 10))
	if sdkVersion > 0 {
		q.Set("sdkVersion", strconv.Itoa(sdkVersion))
	}
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return ttlockKeyListResp{}, err
	}
	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return ttlockKeyListResp{}, errors.New("http status " + strconv.Itoa(resp.StatusCode))
	}
	var out ttlockKeyListResp
	if e := json.Unmarshal(bs, &out); e != nil {
		return ttlockKeyListResp{}, e
	}
	if out.ErrCode != 0 {
		return out, errors.New(out.ErrMsg)
	}
	if out.List == nil {
		out.List = []map[string]any{}
	}
	return out, nil
}

func ttlockFindKeyForLock(cfg ttlockCloudConfig, accessToken string, lockID int64, sdkVersion int) (map[string]any, error) {
	pageNo := 1
	pageSize := 200
	for pageNo <= 20 {
		out, err := ttlockListKeys(cfg, accessToken, pageNo, pageSize, sdkVersion)
		if err != nil {
			return nil, err
		}
		for _, it := range out.List {
			if it == nil {
				continue
			}
			if gf.Int64(it["lockId"]) != lockID {
				continue
			}
			// keyStatus: 110401 正常
			if s := gf.Int(it["keyStatus"]); s != 0 && s != 110401 {
				continue
			}
			return it, nil
		}
		if out.Total > 0 && pageNo*out.PageSize >= out.Total {
			break
		}
		if len(out.List) == 0 {
			break
		}
		pageNo++
	}
	return nil, errors.New("未找到该锁对应的钥匙数据（key/list）")
}

func ttlockUnlock(cfg ttlockCloudConfig, accessToken string, lockID int64) (map[string]any, error) {
	if cfg.ClientID == "" || accessToken == "" || lockID == 0 {
		return nil, errors.New("缺少 client_id/access_token/lock_id")
	}
	form := url.Values{}
	form.Set("clientId", cfg.ClientID)
	form.Set("accessToken", accessToken)
	form.Set("lockId", strconv.FormatInt(lockID, 10))
	form.Set("date", strconv.FormatInt(ttlockNowMilli(), 10))
	return ttlockPostFormJSON(ttlockAPIURL(cfg.APIBase, "/v3/lock/unlock"), form)
}

func ttlockLockDetail(cfg ttlockCloudConfig, accessToken string, lockID int64) (map[string]any, error) {
	if cfg.ClientID == "" || accessToken == "" || lockID == 0 {
		return nil, errors.New("缺少 client_id/access_token/lock_id")
	}
	form := url.Values{}
	form.Set("clientId", cfg.ClientID)
	form.Set("accessToken", accessToken)
	form.Set("lockId", strconv.FormatInt(lockID, 10))
	form.Set("date", strconv.FormatInt(ttlockNowMilli(), 10))
	return ttlockPostFormJSON(ttlockAPIURL(cfg.APIBase, "/v3/lock/detail"), form)
}

// ttlockAddKeyboardPwdGateway 通过网关下发密码（addType=2，不需要 SDK）
func ttlockAddKeyboardPwdGateway(cfg ttlockCloudConfig, accessToken string, lockID int64, keyboardPwd, keyboardPwdName string, startDateMs, endDateMs int64) (map[string]any, error) {
	if cfg.ClientID == "" || accessToken == "" || lockID == 0 {
		return nil, errors.New("缺少 client_id/access_token/lock_id")
	}
	if strings.TrimSpace(keyboardPwd) == "" {
		return nil, errors.New("缺少 keyboard_pwd")
	}
	if startDateMs <= 0 {
		startDateMs = ttlockNowMilli()
	}
	if endDateMs <= 0 {
		endDateMs = startDateMs + 24*60*60*1000
	}
	form := url.Values{}
	form.Set("clientId", cfg.ClientID)
	form.Set("accessToken", accessToken)
	form.Set("lockId", strconv.FormatInt(lockID, 10))
	form.Set("keyboardPwd", strings.TrimSpace(keyboardPwd))
	if strings.TrimSpace(keyboardPwdName) != "" {
		form.Set("keyboardPwdName", strings.TrimSpace(keyboardPwdName))
	}
	form.Set("startDate", strconv.FormatInt(startDateMs, 10))
	form.Set("endDate", strconv.FormatInt(endDateMs, 10))
	form.Set("addType", "2")
	form.Set("date", strconv.FormatInt(ttlockNowMilli(), 10))
	return ttlockPostFormJSON(ttlockAPIURL(cfg.APIBase, "/v3/keyboardPwd/add"), form)
}

func ttlockSendEkey(cfg ttlockCloudConfig, accessToken string, lockID int64, receiverUsername, keyName, remarks string, startDateMs, endDateMs int64) (map[string]any, error) {
	if cfg.ClientID == "" || accessToken == "" || lockID == 0 {
		return nil, errors.New("缺少 client_id/access_token/lock_id")
	}
	if strings.TrimSpace(receiverUsername) == "" {
		return nil, errors.New("缺少 receiver_username")
	}
	if startDateMs <= 0 {
		startDateMs = ttlockNowMilli()
	}
	if endDateMs <= 0 {
		endDateMs = startDateMs + 24*60*60*1000
	}
	form := url.Values{}
	form.Set("clientId", cfg.ClientID)
	form.Set("accessToken", accessToken)
	form.Set("lockId", strconv.FormatInt(lockID, 10))
	form.Set("receiverUsername", strings.TrimSpace(receiverUsername))
	if strings.TrimSpace(keyName) != "" {
		form.Set("keyName", strings.TrimSpace(keyName))
	}
	if strings.TrimSpace(remarks) != "" {
		form.Set("remarks", strings.TrimSpace(remarks))
	}
	form.Set("startDate", strconv.FormatInt(startDateMs, 10))
	form.Set("endDate", strconv.FormatInt(endDateMs, 10))
	form.Set("date", strconv.FormatInt(ttlockNowMilli(), 10))
	return ttlockPostFormJSON(ttlockAPIURL(cfg.APIBase, "/v3/key/send"), form)
}

func ttlockInsertEvent(businessID, userID, propertyID, lockID int64, eventType string, result bool, errCode int, errMsg string, meta any) {
	now := time.Now().Unix()
	metaStr := ""
	if meta != nil {
		if bs, e := json.Marshal(meta); e == nil {
			metaStr = string(bs)
		}
	}
	data := gf.Map{
		"business_id":    businessID,
		"user_id":        userID,
		"property_id":    propertyID,
		"ttlock_lock_id": lockID,
		"event_type":     eventType,
		"result": func() int {
			if result {
				return 1
			}
			return 0
		}(),
		"err_code":   errCode,
		"err_msg":    errMsg,
		"meta_json":  metaStr,
		"createtime": now,
	}
	_, _ = gf.Model("business_lock_events").Data(data).Insert()
}
