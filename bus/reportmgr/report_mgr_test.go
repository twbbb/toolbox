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
