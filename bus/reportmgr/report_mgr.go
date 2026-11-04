package reportmgr

// ReportInfo 报告信息结构体
type ReportInfo struct {
	Rid        string `json:"rid"`         // 报告ID  根据rid可以从cos拿到报告原文
	ReportType int32  `json:"report_type"` // 报告类型  1-pdf,2-图片，3-文本
	Sthid      string `json:"sthid"`       // 用户sthid
	Biz        string `json:"biz"`         // 业务ID 健康管理侧分配
}
