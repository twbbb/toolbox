package ocr

// PatientInfo 患者信息
type PatientInfo struct {
	Name                     string `json:"Name"`                     // 姓名
	Sex                      string `json:"Sex"`                      // 性别
	Age                      string `json:"Age"`                      // 年龄
	Phone                    string `json:"Phone"`                    // 电话
	Address                  string `json:"Address"`                  // 地址
	IdCard                   string `json:"IdCard"`                   // 身份证号
	HealthCardNo             string `json:"HealthCardNo"`             // 健康卡号
	SocialSecurityCardNo     string `json:"SocialSecurityCardNo"`     // 社保卡号
	BedNo                    string `json:"BedNo"`                    // 床号
	Birthday                 string `json:"Birthday"`                 // 生日
	Ethnicity                string `json:"Ethnicity"`                // 民族
	Married                  string `json:"Married"`                  // 婚姻状况
	Profession               string `json:"Profession"`               // 职业
	EducationBackground      string `json:"EducationBackground"`      // 教育背景
	Nationality              string `json:"Nationality"`              // 国籍
	BirthPlace               string `json:"BirthPlace"`               // 出生地
	MedicalInsuranceType     string `json:"MedicalInsuranceType"`     // 医保类型
	AgeNorm                  string `json:"AgeNorm"`                  // 年龄规范
	Nation                   string `json:"Nation"`                   // 国家
	MarriedCode              string `json:"MarriedCode"`              // 婚姻状况代码
	ProfessionCode           string `json:"ProfessionCode"`           // 职业代码
	MedicalInsuranceTypeCode string `json:"MedicalInsuranceTypeCode"` // 医保类型代码
}

// ReportInfo 报告信息
type ReportInfo struct {
	Hospital         string `json:"Hospital"`         // 医院
	DepartmentName   string `json:"DepartmentName"`   // 科室名称
	BillingTime      string `json:"BillingTime"`      // 计费时间
	ReportTime       string `json:"ReportTime"`       // 报告时间
	InspectTime      string `json:"InspectTime"`      // 检查时间
	CheckNum         string `json:"CheckNum"`         // 检查号
	ImageNum         string `json:"ImageNum"`         // 影像号
	RadiationNum     string `json:"RadiationNum"`     // 放射号
	TestNum          string `json:"TestNum"`          // 检验号
	OutpatientNum    string `json:"OutpatientNum"`    // 门诊号
	PathologyNum     string `json:"PathologyNum"`     // 病理号
	InHospitalNum    string `json:"InHospitalNum"`    // 住院号
	SampleNum        string `json:"SampleNum"`        // 样本号
	SampleType       string `json:"SampleType"`       // 样本类型
	MedicalRecordNum string `json:"MedicalRecordNum"` // 病历号
	ReportName       string `json:"ReportName"`       // 报告名称
	UltraNum         string `json:"UltraNum"`         // 超声号
	Diagnose         string `json:"Diagnose"`         // 诊断
	PrintTime        string `json:"PrintTime"`        // 打印时间
	DiagnoseTime     string `json:"DiagnoseTime"`     // 诊断时间
	OtherTime        string `json:"OtherTime"`        // 其他时间
	BedNo            string `json:"BedNo"`            // 床号
	HealthCheckupNum string `json:"HealthCheckupNum"` // 体检号
	CheckItem        string `json:"CheckItem"`        // 检查项目
	CheckMethod      string `json:"CheckMethod"`      // 检查方法
}

// CheckText 检查文本信息（当TextReportType为1时该字段必填）
type CheckText struct {
	Description string `json:"Description"` // 描述
	Summary     string `json:"Summary"`     // 总结
}

// IndicatorText 指标文本信息
type IndicatorText struct {
	Code       string `json:"Code"`       // 代码
	ItemString string `json:"ItemString"` // 项目名称
	Result     string `json:"Result"`     // 结果
	Unit       string `json:"Unit"`       // 单位
	Range      string `json:"Range"`      // 范围
	Arrow      string `json:"Arrow"`      // 箭头方向
}

// MedDocText 医疗文档文本
type MedDocText struct {
	MedDoc string `json:"MedDoc"` // 医疗文档
	Origin string `json:"Origin"` // 来源
}

// OCRData OCR识别数据主结构体
type OCRData struct {
	PatientInfo       PatientInfo     `json:"PatientInfo"`         // 患者信息
	ReportInfo        ReportInfo      `json:"ReportInfo"`          // 报告信息
	TextReportType    int             `json:"TextReportType"`      // 文本报告类型：1-检查类 2-检验类 3-病例文书类
	CheckText         *CheckText      `json:"CheckText,omitempty"` // 检查文本（当TextReportType为1时必填）
	IndicatorTextList []IndicatorText `json:"IndicatorTextList"`   // 指标文本列表
	MedDocText        MedDocText      `json:"MedDocText"`          // 医疗文档文本
}

// TextPngRedisData 文字转图片Redis存储数据结构
type TextPngRedisData struct {
	BizId string `json:"bizId"` // 业务ID
	Rid   string `json:"rid"`   // 资源ID
	Md5   string `json:"md5"`   // MD5值
	Sthid string `json:"sthid"` // 存储ID
}
