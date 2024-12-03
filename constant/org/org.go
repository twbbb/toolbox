// Package org 机构
package org

// mapDataReq 请求参数
const (
	REQ_AREACODE     = "areaCode"     // 区域编码
	REQ_SHOWALLENTRY = "showAllEntry" // true显示所有入口
	REQ_ORGTYPE      = "orgType"      // 医疗机构类型
	REQ_LEVEL        = "level"        // 医院等级
	REQ_ORGCATEGORY  = "orgCategory"  // 机构经营类别
	REQ_ACTIVITY     = "activity"     // 参与的活动
	REQ_LAT          = "lat"          // 纬度
	REQ_LNG          = "lng"          // 经度
	REQ_SORTTYPE     = "sortType"     // 排序类型
	REQ_EXAMPKG      = "examPkg"      // 体检套餐过滤
	REQ_TAGTYPE      = "tagType"      // 标签类型
	REQ_CFGBIZ       = "cfgBiz"       // 待获取的配置类型，逗号分隔
	// REQ_ENTRYID 入口名称，注意优先使用的map外的EntryFilter。register 预约挂号；examination 体检；healthcare 消费医疗；vaccination 疫苗
	REQ_ENTRYID = "entryId"
)

// orgInfo 机构信息
const (
	ORGINFO_CITYCODE        = "cityCode"        // 城市代码
	ORGINFO_CITYNAME        = "cityName"        // 城市名称
	ORGINFO_AREACODE        = "areaCode"        // 区域代码
	ORGINFO_ORGID           = "orgId"           // 医疗机构ID
	ORGINFO_ORGTYPE         = "orgType"         // 医疗机构类型
	ORGINFO_ORGNAME         = "orgName"         // 医疗机构名称
	ORGINFO_ORGADDR         = "orgAddr"         // 医疗机构地址
	ORGINFO_ORGLEVEL        = "orgLevel"        // 医院等级
	ORGINFO_STATUS          = "status"          // 状态
	ORGINFO_LOGOURL         = "logoUrl"         // logo地址
	ORGINFO_ORGPHONE        = "orgPhone"        // 机构电话
	ORGINFO_ORGCLASS        = "orgClass"        // 机构性质
	ORGINFO_ORGCATEGORY     = "orgCategory"     // 机构经营类别
	ORGINFO_ORGCATEGORYNAME = "orgCategoryName" // 机构经营类别名
	ORGINFO_BUSIINFOS       = "busiInfos"       // 业务信息
	ORGINFO_LAT             = "lat"             // 纬度
	ORGINFO_LNG             = "lng"             // 经度
	ORGINFO_DESC            = "desc"            // 机构简介
	ORGINFO_CFG             = "cfg"             // 配置信息
	ORGINFO_STRORGLEVEL     = "strOrgLevel"     // 医院级别描述
	ORGINFO_PICURL          = "picUrl"          // 背景图,医院主页使用
	ORGINFO_OFFICETIME      = "officeTime"      // 营业时间,医院主页使用
	ORGINFO_ORGALIASNAME    = "orgAliasName"    // 机构别名,医院主页使用
)

// entryArr 入口信息
const (
	ENTRYARR_LOGOURL       = "logoUrl"       // logo图片链接
	ENTRYARR_STATUS        = "status"        // 状态 0 - 正常状态, >0 - 失效状态
	ENTRYARR_SVRDESC       = "svrDesc"       // 服务描述
	ENTRYARR_ENTRYID       = "entryId"       // 入口项ID
	ENTRYARR_CARDICON      = "cardIcon"      // 卡片图标
	ENTRYARR_ONLINETIME    = "onlineTime"    // 入口上线时间
	ENTRYARR_WHITECHANNELS = "whiteChannels" // 白牌渠道
	ENTRYARR_PAYINFO       = "payInfo"       // 支付信息
	ENTRYARR_ALIASSUBTITLE = "aliasSubTitle" // 别名副标题
	ENTRYARR_BGIMG         = "bgImg"         // 背景图片
	ENTRYARR_STYLETYPE     = "styleType"     // 样式类型
	ENTRYARR_TITLE         = "title"         // 标题
	ENTRYARR_ALIASTITLE    = "aliasTitle"    // 别名标题
	ENTRYARR_JUMPLINK      = "jumpLink"      // 跳转链接信息
	ENTRYARR_PARTNERIDS    = "partnerIds"    // 合作方ID列表
)

// tagArr 标签信息
const (
	TAGARR_LABLE = "lable" // 标签名
	TAGARR_TYPE  = "type"  // 类型
	TAGARR_ORDER = "order" // 排序
	TAGARR_STYLE = "style" // 样式
)

// userInfo 用户信息
const (
	USERINFO_DISTANCE    = "distance"    // 距离
	USERINFO_ISCONCERNED = "isConcerned" // 是否关注
	USERINFO_ISUSED      = "isUsed"      // 是否使用
	USERINFO_LSTUSETIME  = "lstUseTime"  // 最后使用时间
	USERINFO_ISNEARST    = "isNearst"    // 是否最近
)

// registerData 挂号信息
const (
	// --- 挂号基础信息 ---
	REGISTERDATA_BASIC               = "basic"         // 挂号基础信息
	REGISTERDATA_BASIC_REGISTERMODE  = "registerMode"  // 挂号模式
	REGISTERDATA_BASIC_PLATFORMID    = "platformId"    // 城市服务平台ID
	REGISTERDATA_BASIC_DATAVALID     = "dataValid"     // 挂号数据是否可用
	REGISTERDATA_BASIC_IMPORTANTDEP  = "importantDep"  // 挂号优势科室
	REGISTERDATA_BASIC_ORGVIDEOTITLE = "orgVideoTitle" // 机构视频号标题
	REGISTERDATA_BASIC_ORGVIDEOICON  = "orgVideoIcon"  // 机构视频号图标
	REGISTERDATA_BASIC_ORGVIDEOCOUNT = "orgVideoCount" // 机构视频号总数

	// --- 用户信息 -----
	REGISTERDATA_USER                = "user"           // 挂号用户信息
	REGISTERDATA_USER_REGISTERID     = "registerId"     // 用户相关的挂号ID
	REGISTERDATA_USER_CITYSVCHASAUTH = "citySvcHasAuth" // 城市服务授权信息
)

// examData
const (
	// ---体检基础信息---
	EXAMDATA_BASIC              = "basic"        // 体检基础信息
	EXAMDATA_BASIC_PLATFORMID   = "platformId"   // 体检合作方ID
	EXAMDATA_BASIC_PARTNERORGID = "partnerOrgId" // 体检外部医院ID

	// --- 体检套餐信息 -----
	EXAMDATA_PKG                = "pkg"                // 体检套餐信息
	EXAMDATA_PKG_EXAMPKGCNT     = "pkg.examPkgCnt"     // 套餐数量
	EXAMDATA_PKG_OUTSHOWPKGLIST = "pkg.outShowPkgList" // 体检套餐
)

// activityArr
const (
	ACTIVITYARR_ACTIVITYID   = "activityID"   // 活动ID
	ACTIVITYARR_ACTIVITYNAME = "activityName" // 活动名称
	ACTIVITYARR_PRIORITY     = "priority"     // 活动优先级
	ACTIVITYARR_TAG          = "tag"          // 活动标签
	ACTIVITYARR_WORDING      = "wording"      // 规则文案
	ACTIVITYARR_CATEGORY     = "category"     // 分类
	ACTIVITYARR_WORKON       = "workon"       // 入口服务
	ACTIVITYARR_ICON         = "icon"         // 图标
	ACTIVITYARR_TAGSTYLE     = "tagStyle"     // 标签样式
	ACTIVITYARR_COUPONIDS    = "couponIds"    // 涉及到的优惠券
	ACTIVITYARR_SCORE        = "score"        // 积分金额
	ACTIVITYARR_PARTNERSARR  = "partnersArr"  // 合作方ID列表
	ACTIVITYARR_PKGIDSARR    = "pkgIdsArr"    // 套餐ID
)

// cfgArr
const (
	CFGARR_BIZID = "bizID" // 业务id
	CFGARR_CFG   = "cfg"   // 配置数据
)
