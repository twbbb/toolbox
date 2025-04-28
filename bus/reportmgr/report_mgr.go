package reportmgr

import (
	"context"
	"encoding/json"
	"errors"

	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.woa.com/th-go/toolbox/component/aesutil"
)

type ReportInfo struct {
	Rid        string `json:"rid"`         // 报告ID  根据rid可以从cos拿到报告原文
	ReportType int32  `json:"report_type"` // 报告类型  1-pdf,2-图片
	Sthid      string `json:"sthid"`       // 用户sthid
	Biz        string `json:"biz"`         // 业务ID 健康管理侧分配
}

// EncReportInfo 加密报告信息
func EncReportInfo(ctx context.Context, aesKey string,
	rid string,
	reportType int32,
	sthid string,
	biz string,
) (string, error) {

	report := &ReportInfo{
		Rid:        rid,
		ReportType: reportType,
		Sthid:      sthid,
		Biz:        biz,
	}
	reportStr, err := json.Marshal(report)
	if err != nil {
		log.ErrorContextf(ctx, "EncReportInfo,json.Marshal,err:%v", err)
		return "", err
	}
	encReport, err := aesutil.AesEncryptStr(string(reportStr), aesKey)
	if err != nil {
		log.ErrorContextf(ctx, "EncReportInfo,aesutil.AesEncryptStr,err:%v", err)
		return "", err
	}

	return encReport, nil
}

// DecReportInfo 解密报告信息
func DecReportInfo(ctx context.Context, aesKey string, reportStr string) (*ReportInfo, error) {
	if reportStr == "" {
		return nil, errors.New("reportStr is empty")
	}
	decReport, err := aesutil.AesDecryptStr(reportStr, aesKey)
	if err != nil {
		log.ErrorContextf(ctx, "DecReportInfo,aesutil.AesDecryptStr,err:%v", err)
		return nil, err
	}

	report := &ReportInfo{}
	err = json.Unmarshal([]byte(decReport), report)
	if err != nil {
		log.ErrorContextf(ctx, "DecReportInfo,json.Unmarshal,err:%v", err)
		return nil, err
	}
	return report, nil
}
