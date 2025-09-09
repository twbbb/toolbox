package reportmgr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.woa.com/th-go/toolbox/component/aesutil"
	"net/url"
	"strconv"
)

// ReportInfo 报告信息结构体
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
	FileType           string `json:"fileType"`           // 文件类型 1:pdf， 2:图片 ,前端使用为string
	FileUrl            string `json:"fileUrl"`            // 文件地址
	OrgName            string `json:"orgName"`            // 机构名
	ReportType         string `json:"reportType"`         // 报告类型 如'体检报告'
}

// ReportTypeToChineseName 报告类型到中文名称的映射
var ReportTypeToChineseName = map[int32]string{
	REPORT_TYPE_UNSPECIFIED: "未指定",
	REPORT_TYPE_EXAM:        "检查报告",
	REPORT_TYPE_LAB:         "检验报告",
	REPORT_TYPE_RECORD:      "病例文书报告",
}

const (
	REPORT_TYPE_UNSPECIFIED = 0 // 未指定
	REPORT_TYPE_EXAM        = 1 // 检查报告
	REPORT_TYPE_LAB         = 2 // 检验报告
	REPORT_TYPE_RECORD      = 3 // 病例文书报告
)

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

// GetReportType 获取报告类型中文名称
func GetReportType(reportType int32) string {
	return ReportTypeToChineseName[reportType]
}

// BuildHealthManagerParams 构建健康管理参数JSON字符串
func BuildHealthManagerParams(
	otherBizReportInfo string,
	fileName string,
	fileType int32,
	fileUrl string,
	orgName string,
	reportType int32,
) (string, error) {
	params := &HealthManagerParams{
		OtherBizReportInfo: otherBizReportInfo,
		FileName:           fileName,
		FileType:           strconv.Itoa(int(fileType)),
		FileUrl:            fileUrl,
		OrgName:            orgName,
		ReportType:         GetReportType(reportType),
	}

	// 将结构体转换为JSON字符串
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// BuildHealthManagerURL 构建完整的健康管理URL, fileType 文件类型  1-pdf,2-图片, reportType 报告类型对应 ReportTypeToChineseName
func BuildHealthManagerURL(
	isBizReport string,
	bizId string,
	otherBizReportInfo string,
	fileName string,
	fileType int32,
	fileUrl string,
	orgName string,
	reportType int32,
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

// BuildHealthManagerURLForOcrCenter 构建完整的健康管理URL, fileType对应OcrCenter里面的
func BuildHealthManagerURLForOcrCenter(
	ctx context.Context,
	isBizReport string,
	bizId string,
	rid string,
	sthid string,
	fileName string,
	fileType int32,
	fileUrl string,
	orgName string,
	reportType int32, aeskey string) (string, error) {

	// 从ocrcenter过来的 0是图片，1是pdf, 2是文本
	switch reportType {
	case 0:
		reportType = 1
	case 1:
		reportType = 2
	case 2:
		reportType = 3
	}
	otherBizReportInfo, err := EncReportInfo(ctx, aeskey, rid, reportType, sthid, bizId)
	fullURL, err := BuildHealthManagerURL(
		isBizReport,
		bizId,
		otherBizReportInfo,
		fileName,
		fileType,
		fileUrl,
		orgName,
		reportType,
	)
	return fullURL, err
}
