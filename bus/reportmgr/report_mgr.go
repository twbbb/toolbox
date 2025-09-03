package reportmgr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.woa.com/th-go/toolbox/component/aesutil"
)

type ReportInfo struct {
	Rid        string `json:"rid"`         // 报告ID  根据rid可以从cos拿到报告原文
	ReportType int32  `json:"report_type"` // 报告类型  1-pdf,2-图片
	Sthid      string `json:"sthid"`       // 用户sthid
	Biz        string `json:"biz"`         // 业务ID 健康管理侧分配
}

// HealthManagerParams 健康管理参数结构体
type HealthManagerParams struct {
	OtherBizReportInfo string `json:"otherBizReportInfo"` // 加密的其他业务报告信息
	FileName           string `json:"fileName"`           // 文件名
	FileType           string `json:"fileType"`           // 文件类型 1:pdf， 2:图片
	FileUrl            string `json:"fileUrl"`            // 文件地址
	OrgName            string `json:"orgName"`            // 机构名
	ReportType         string `json:"reportType"`         // 报告类型 如'体检报告'
}

// BuildHealthManagerParams 构建健康管理参数JSON字符串
func BuildHealthManagerParams(
	otherBizReportInfo string,
	fileName string,
	fileType string,
	fileUrl string,
	orgName string,
	reportType string,
) (string, error) {
	params := &HealthManagerParams{
		OtherBizReportInfo: otherBizReportInfo,
		FileName:           fileName,
		FileType:           fileType,
		FileUrl:            fileUrl,
		OrgName:            orgName,
		ReportType:         reportType,
	}

	// 将结构体转换为JSON字符串
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// BuildHealthManagerURL 构建完整的健康管理URL
func BuildHealthManagerURL(
	isBizReport string,
	bizId string,
	otherBizReportInfo string,
	fileName string,
	fileType string,
	fileUrl string,
	orgName string,
	reportType string,
) (string, error) {
	// 构建healthManagerParams JSON字符串
	healthManagerParams, err := BuildHealthManagerParams(
		otherBizReportInfo,
		fileName,
		fileType,
		fileUrl,
		orgName,
		reportType,
	)
	if err != nil {
		return "", err
	}

	// 对JSON字符串进行URL编码
	encodedParams := url.QueryEscape(healthManagerParams)

	// 构建完整URL
	fullURL := fmt.Sprintf("pages/health-manager/home/main?isBizReport=%s&bizId=%s&healthManagerParams=%s",
		isBizReport, bizId, encodedParams)

	return fullURL, nil
}

// ParseHealthManagerParams 解析健康管理参数
func ParseHealthManagerParams(encodedParams string) (*HealthManagerParams, error) {
	if encodedParams == "" {
		return nil, errors.New("encodedParams is empty")
	}

	// URL解码
	decodedParams, err := url.QueryUnescape(encodedParams)
	if err != nil {
		return nil, err
	}

	// JSON解析
	params := &HealthManagerParams{}
	err = json.Unmarshal([]byte(decodedParams), params)
	if err != nil {
		return nil, err
	}

	return params, nil
}

// ParseHealthManagerParamsFromURL 从完整URL中解析健康管理参数
func ParseHealthManagerParamsFromURL(fullURL string) (*HealthManagerParams, error) {
	if fullURL == "" {
		return nil, errors.New("fullURL is empty")
	}

	// 解析URL
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}

	// 获取查询参数
	queryParams := parsedURL.Query()
	healthManagerParams := queryParams.Get("healthManagerParams")

	if healthManagerParams == "" {
		return nil, errors.New("healthManagerParams not found in URL")
	}

	// 解析健康管理参数
	return ParseHealthManagerParams(healthManagerParams)
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
