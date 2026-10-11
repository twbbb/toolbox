package reportmgr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/twbbb/toolbox/component/aesutil"
	"trpc.group/trpc-go/trpc-go/log"
)

// ReportInfo 报告信息结构体
type ReportInfo struct {
	Rid        string `json:"rid"`         // 报告ID  根据rid可以从cos拿到报告原文
	ReportType int32  `json:"report_type"` // 报告类型  1-pdf,2-图片，3-文本
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
	REPORT_TYPE_UNSPECIFIED: "",
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

// 对应健康管理的文件映射
const (
	FILE_TYPE_PDF   = 1
	FILE_TYPE_IMAGE = 2
	FILE_TYPE_TEXT  = 3
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

// BuildHealthManagerURLWithChannel 构建完整的健康管理URL，支持channel参数
func BuildHealthManagerURLWithChannel(
	isBizReport string,
	bizId string,
	channel string,
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

	// 构建完整URL，添加channel参数
	fullURL := fmt.Sprintf("pages/health-manager/home/main?isBizReport=%s&bizId=%s&channel=%s&healthManagerParams=%s",
		isBizReport, bizId, channel, encodedParams)

	return fullURL, nil
}

// ReportEncryptOption 用于生成otherBizReportInfo的加密参数
type ReportEncryptOption struct {
	Rid      string // 报告ID（必填）
	Sthid    string // 用户sthid（必填）
	FileType int32  // 文件类型 OcrCenter: 0-图片，1-pdf，2-文本（必填）
	BizId    string // 业务ID，作为biz字段加密（必填）
	AesKey   string // AES加密密钥（必填）
}

// HealthManagerURLOption 构建健康管理URL的参数
type HealthManagerURLOption struct {
	IsBizReport string // 是否业务报告（必填）
	BizId       string // 业务ID（必填）
	Channel     string // 渠道（可选）
	GroupId     string // 分组ID（可选）
	FileName    string // 文件名（可选）
	FileType    int32  // 文件类型 OcrCenter: 0-图片，1-pdf，2-文本（必填）
	FileUrl     string // 文件地址（可选）
	OrgName     string // 机构名（可选）
	ReportType  int32  // 报告类型（可选）
}

// BuildHealthManagerURLForOcrCenter 构建完整的健康管理URL, fileType对应OcrCenter里面的
// 当 encOpts 为 nil 时，不进行加密，直接构建URL（不包含 healthManagerParams）
func BuildHealthManagerURLForOcrCenter(
	ctx context.Context,
	encOpts *ReportEncryptOption,
	urlOpts *HealthManagerURLOption,
) (string, error) {

	// 构建完整URL
	fullURL := fmt.Sprintf("pages/health-manager/home/main?isBizReport=%s&bizId=%s",
		urlOpts.IsBizReport, urlOpts.BizId)
	if urlOpts.Channel != "" {
		fullURL += fmt.Sprintf("&channel=%s", urlOpts.Channel)
	}
	if urlOpts.GroupId != "" {
		fullURL += fmt.Sprintf("&groupId=%s", urlOpts.GroupId)
	}

	// encOpts 为空时不需要加密，直接返回URL
	if encOpts == nil {
		return fullURL, nil
	}

	fileType := encOpts.FileType
	// 从ocrcenter过来的 0是图片，1是pdf, 2是文本
	switch fileType {
	case 0:
		fileType = FILE_TYPE_IMAGE
	case 1:
		fileType = FILE_TYPE_PDF
	case 2:
		fileType = FILE_TYPE_TEXT
	}
	otherBizReportInfo, err := EncReportInfo(ctx, encOpts.AesKey, encOpts.Rid, fileType, encOpts.Sthid, encOpts.BizId)
	if err != nil {
		return "", err
	}

	// 构建healthManagerParams JSON字符串
	healthManagerParams, err := BuildHealthManagerParams(
		otherBizReportInfo,
		urlOpts.FileName,
		fileType,
		urlOpts.FileUrl,
		urlOpts.OrgName,
		urlOpts.ReportType,
	)
	if err != nil {
		return "", err
	}

	// 对JSON字符串进行URL编码
	encodedParams := url.QueryEscape(healthManagerParams)
	fullURL += fmt.Sprintf("&healthManagerParams=%s", encodedParams)

	return fullURL, nil
}

// BuildHealthManagerURLForOcrCenterWithChannel 构建完整的健康管理URL，支持channel参数，fileType对应OcrCenter里面的
// Deprecated: 请使用 BuildHealthManagerURLForOcrCenter，通过 HealthManagerURLOption.Channel 设置channel
func BuildHealthManagerURLForOcrCenterWithChannel(
	ctx context.Context,
	isBizReport string,
	bizId string,
	channel string,
	rid string,
	sthid string,
	fileName string,
	fileType int32,
	fileUrl string,
	orgName string,
	reportType int32,
	aeskey string) (string, error) {

	return BuildHealthManagerURLForOcrCenter(ctx, &ReportEncryptOption{
		Rid:      rid,
		Sthid:    sthid,
		FileType: fileType,
		BizId:    bizId,
		AesKey:   aeskey,
	}, &HealthManagerURLOption{
		IsBizReport: isBizReport,
		BizId:       bizId,
		Channel:     channel,
		FileName:    fileName,
		FileType:    fileType,
		FileUrl:     fileUrl,
		OrgName:     orgName,
		ReportType:  reportType,
	})
}

// BuildHealthManagerURLByPatientID 构建完整的健康管理URL
func BuildHealthManagerURLByPatientID(
	patientId string,
	bizId string,
	channel string,
) (string, error) {
	// 构建完整URL，添加channel参数
	fullURL := fmt.Sprintf("pages/health-manager/home/main?patientId=%s&bizId=%s&channel=%s",
		patientId, bizId, channel)

	return fullURL, nil
}
