package reportmgr

import (
	"context"
	"fmt"
	"testing"
)

func TestReportMgr(t *testing.T) {
	ctx := context.Background()
	aeskey := "vDxVjrAowp5GcvXs"
	info, err := EncReportInfo(ctx, aeskey, "123", REPORT_TYPE_EXAM, "123", "123")
	if err != nil {
		t.Fail()
		return
	}

	fmt.Println(info)

	decInfo, err := DecReportInfo(ctx, aeskey, info)
	if err != nil {
		t.Fail()
		return
	}
	fmt.Println(decInfo)
	if decInfo.Rid != "123" || decInfo.ReportType != REPORT_TYPE_EXAM || decInfo.Sthid != "123" || decInfo.Biz != "123" {
		t.Fail()
		return
	}

}

func TestReportMgr2(t *testing.T) {
	// 测试数据
	otherBizReportInfo := "sf7+gVLOjt1DgqOh9bauEEPh+CAbrw6iHW+8NjD0dYbYWjxYNPMlEoYjsKcJWLiBtPmI+/+myweG4BahspqQQMy1WCYDlcniyKKwAVqeMpx9/EJvoVeibIpXy3P/ZW1wnM2Ih2osvlp/9OauM+1TyvHaH590MQdzk9jChs+PYNjP3ftyw2cCprpUvG80uIlhKPzN6zX3zeiuEMZG+ze3L0Bs9x7VKenk6WTeFXLL9lS25DwP6v0prd+hpxMuY2MlkcUMxXlNptU6JGviEaFAlQ=="
	fileName := "邓炜豪的体检报告"
	orgName := "香港大学深圳医院"
	fileType := int32(0)
	fileUrl := ""

	// 测试构建完整URL
	fullURL, err := BuildHealthManagerURL(
		"1",
		"jiankangwenwen",
		otherBizReportInfo,
		fileName,
		fileType,
		fileUrl,
		orgName,
		10,
	)
	if err != nil {
		t.Errorf("BuildHealthManagerURL failed: %v", err)
		return
	}
	fmt.Printf("Generated URL: %s\n", fullURL)
}

func TestReportMgr3(t *testing.T) {
	ctx := context.Background()
	// 0是图片，1是pdf, 2是文本
	// 	REPORT_TYPE_UNSPECIFIED = 0 // 未指定
	//	REPORT_TYPE_EXAM        = 1 // 检查报告
	//	REPORT_TYPE_LAB         = 2 // 检验报告
	//	REPORT_TYPE_RECORD      = 3 // 病例文书报告
	url, _ := BuildHealthManagerURLForOcrCenter(ctx, "1", "diagnosis",
		"14ARD08NrtyeXnpCqME7KBG81En7dzBDGojZcU52HcVs_HZXoahhHxCPgpXeLTWjJg-vLNF",
		"@athooJQLI1ozwL3eroVjJUoef6w7IqzMm8", "", 1, "", "豪贵测试医院", 0, "tZa8k5X7HO5vT9Iw")
	fmt.Println(url)
}
