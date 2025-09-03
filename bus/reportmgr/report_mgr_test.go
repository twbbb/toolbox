package reportmgr

import (
	"context"
	"fmt"
	"testing"
)

func TestReportMgr(t *testing.T) {
	ctx := context.Background()
	aeskey := "vDxVjrAowp5GcvXs"
	info, err := EncReportInfo(ctx, aeskey, "123", 1, "123", "123")
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
	if decInfo.Rid != "123" || decInfo.ReportType != 1 || decInfo.Sthid != "123" || decInfo.Biz != "123" {
		t.Fail()
		return
	}

}

func TestReportMgr2(t *testing.T) {
	// 测试数据
	otherBizReportInfo := "sf7+gVLOjt1DgqOh9bauEEPh+CAbrw6iHW+8NjD0dYbYWjxYNPMlEoYjsKcJWLiBtPmI+/+myweG4BahspqQQMy1WCYDlcniyKKwAVqeMpx9/EJvoVeibIpXy3P/ZW1wnM2Ih2osvlp/9OauM+1TyvHaH590MQdzk9jChs+PYNjP3ftyw2cCprpUvG80uIlhKPzN6zX3zeiuEMZG+ze3L0Bs9x7VKenk6WTeFXLL9lS25DwP6v0prd+hpxMuY2MlkcUMxXlNptU6JGviEaFAlQ=="
	fileName := "邓炜豪的体检报告"
	orgName := "香港大学深圳医院"
	fileType := "0"
	reportType := ""
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
		reportType,
	)
	if err != nil {
		t.Errorf("BuildHealthManagerURL failed: %v", err)
		return
	}
	fmt.Printf("Generated URL: %s\n", fullURL)
}
