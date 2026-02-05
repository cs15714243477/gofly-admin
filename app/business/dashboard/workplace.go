package dashboard

import (
	"gofly/utils/gf"
	"time"
)

/**
* 使用说明：
* 首页统计是根据业务需求数据来统计的，框架无法预知你的项目实际需求，我们只能内置一些方法仅供参考，
* 实际项目开发完成后，根据项目需求自己编写统计数据接口
* business_youtablebane 是你的项目实际数据表(泛指)，不是实际测存在表，切记！自己根据需求开发出对应接口
 */
type Workplace struct{ NoNeedAuths []string }

func init() {
	fpath := Workplace{NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

// 1获取快捷操作
func (api *Workplace) GetQuick(c *gf.GinCtx) {
	businessID := c.GetInt64("businessID") //当前商户ID
	list, err := gf.Model("business_home_quickop").Where("business_id", businessID).WhereOr("is_common", 1).Fields("id,uid,path_url,name,icon,type,is_common,weigh").Order("weigh asc,id asc").Select()
	if err != nil {
		gf.Failed().SetMsg("获取快捷操作失败").SetData(err).Regin(c)
	} else {
		gf.Success().SetMsg("获取快捷操作数据").SetData(list).Regin(c)
	}
}

// 3保存快捷操作
func (api *Workplace) SaveQuick(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	var f_id = gf.GetEditId(param["id"])
	if f_id == 0 {
		param["uid"] = c.GetInt64("userID")             //当前用户
		param["business_id"] = c.GetInt64("businessID") //当前商户
		addId, err := gf.Model("business_home_quickop").Data(param).InsertAndGetId()
		if err != nil {
			gf.Failed().SetMsg("添加失败").SetData(err).Regin(c)
		} else {
			if addId != 0 {
				gf.Model("business_home_quickop").Data(map[string]interface{}{"weigh": addId}).Where("id", addId).Update()
			}
			gf.Success().SetMsg("添加成功！").SetData(addId).Regin(c)
		}
	} else {
		res, err := gf.Model("business_home_quickop").Data(param).Where("id", f_id).Update()
		if err != nil {
			gf.Failed().SetMsg("更新失败").SetData(err).Regin(c)
		} else {
			gf.Success().SetMsg("更新成功！").SetData(res).Regin(c)
		}
	}
}

// 3删除快捷操作
func (api *Workplace) DelQuick(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	res2, err := gf.Model("business_home_quickop").Where("id", param["id"]).Delete()
	if err != nil {
		gf.Failed().SetMsg("删除失败").SetData(err).Regin(c)
	} else {
		gf.Success().SetMsg("删除成功！").SetData(res2).Regin(c)
	}
}

// 获取首页概况统计
// GET /business/dashboard/workplace/get_statistical
func (api *Workplace) Get_statistical(c *gf.GinCtx) {
	businessID := c.GetInt64("businessID") // 当前商户ID
	if businessID == 0 {
		gf.Success().SetMsg("获取概况统计成功").SetData(gf.Map{
			"propertyTotal":       0,
			"propertyOnSale":      0,
			"lockBindTotal":       0,
			"unlockPendingTotal":  0,
			"todayPropertyAdd":    0,
			"todayUnlockRequests": 0,
		}).Regin(c)
		return
	}

	// 1) 房源统计
	propertyTotal, err := gf.Model("business_properties").
		Where("business_id", businessID).
		Count()
	if err != nil {
		gf.Failed().SetMsg("获取房源统计失败").SetData(err).Regin(c)
		return
	}
	propertyOnSale, err := gf.Model("business_properties").
		Where("business_id", businessID).
		Where("sale_status", "on_sale").
		Count()
	if err != nil {
		gf.Failed().SetMsg("获取在售房源统计失败").SetData(err).Regin(c)
		return
	}

	// 2) 智能锁绑定统计
	lockBindTotal, err := gf.Model("business_property_locks").
		Where("business_id", businessID).
		Where("bind_status", 1).
		Count()
	if err != nil {
		gf.Failed().SetMsg("获取智能锁统计失败").SetData(err).Regin(c)
		return
	}

	// 3) 待审核开锁申请（按房源归属商户过滤）
	unlockPendingTotal, err := gf.Model("business_unlock_requests", "ur").
		InnerJoin("business_properties", "p", "p.id=ur.property_id").
		Where("p.business_id", businessID).
		Where("ur.request_status", "pending").
		Count()
	if err != nil {
		gf.Failed().SetMsg("获取开锁申请统计失败").SetData(err).Regin(c)
		return
	}

	// 4) 今日新增房源 / 今日开锁申请
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayPropertyAdd, err := gf.Model("business_properties").
		Where("business_id", businessID).
		Where("createtime", ">=", todayStart).
		Count()
	if err != nil {
		gf.Failed().SetMsg("获取今日新增房源失败").SetData(err).Regin(c)
		return
	}
	todayUnlockRequests, err := gf.Model("business_unlock_requests", "ur").
		InnerJoin("business_properties", "p", "p.id=ur.property_id").
		Where("p.business_id", businessID).
		Where("ur.createtime", ">=", todayStart).
		Count()
	if err != nil {
		gf.Failed().SetMsg("获取今日开锁申请失败").SetData(err).Regin(c)
		return
	}

	gf.Success().SetMsg("获取概况统计成功").SetData(gf.Map{
		"propertyTotal":       propertyTotal,
		"propertyOnSale":      propertyOnSale,
		"lockBindTotal":       lockBindTotal,
		"unlockPendingTotal":  unlockPendingTotal,
		"todayPropertyAdd":    todayPropertyAdd,
		"todayUnlockRequests": todayUnlockRequests,
	}).Regin(c)
}

// 获取近 7 天趋势（默认：新增房源）
// GET /business/dashboard/workplace/get_visitlist
func (api *Workplace) Get_visitlist(c *gf.GinCtx) {
	businessID := c.GetInt64("businessID") // 当前商户ID

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start := todayStart.AddDate(0, 0, -6)

	// 先构造 7 天日期轴，补齐缺失日期
	dates := make([]string, 0, 7)
	dateCountMap := make(map[string]int, 7)
	for i := 0; i < 7; i++ {
		key := start.AddDate(0, 0, i).Format("2006-01-02")
		dates = append(dates, key)
		dateCountMap[key] = 0
	}

	// businessID 为 0 时直接返回 0 数据
	if businessID != 0 {
		list, err := gf.Model("business_properties").
			Fields("DATE_FORMAT(createtime,'%Y-%m-%d') as x, COUNT(*) as y").
			Where("business_id", businessID).
			Where("createtime", ">=", start).
			Group("DATE_FORMAT(createtime,'%Y-%m-%d')").
			Order("x asc").
			Select()
		if err != nil {
			gf.Failed().SetMsg("获取趋势数据失败").SetData(err).Regin(c)
			return
		}
		for _, item := range list {
			x := item["x"].String()
			y := item["y"].Int()
			dateCountMap[x] = y
		}
	}

	result := make([]gf.Map, 0, 7)
	for _, x := range dates {
		result = append(result, gf.Map{
			"x": x,
			"y": dateCountMap[x],
		})
	}
	gf.Success().SetMsg("获取趋势数据成功").SetData(result).Regin(c)
}

// 获取热门房源（按浏览量排序）
// GET /business/dashboard/workplace/get_popular
func (api *Workplace) Get_popular(c *gf.GinCtx) {
	businessID := c.GetInt64("businessID") // 当前商户ID
	if businessID == 0 {
		gf.Success().SetMsg("获取热门房源成功").SetData([]interface{}{}).Regin(c)
		return
	}

	list, err := gf.Model("business_properties").
		Fields("id,title,view_count as viewCount,follow_count as followCount,showing_count as showingCount,sale_status as saleStatus,price,price_unit as priceUnit,area").
		Where("business_id", businessID).
		Order("view_count desc,id desc").
		Limit(10).
		Select()
	if err != nil {
		gf.Failed().SetMsg("获取热门房源失败").SetData(err).Regin(c)
		return
	}
	gf.Success().SetMsg("获取热门房源成功").SetData(list).Regin(c)
}
