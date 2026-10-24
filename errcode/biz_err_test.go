package errcode

import (
	"testing"

	"trpc.group/trpc-go/trpc-go/errs"
)

func codeOf(t *testing.T, bizType BizErrType, client bool, code int) int {
	t.Helper()
	be := GetBizErr(bizType)
	if client {
		_ = be.AddBizErrClient(code, "x")
		return int(errs.Code(be.GetBizErrClient(code)))
	}
	_ = be.AddBizErrServer(code, "x")
	return int(errs.Code(be.GetBizErrServer(code)))
}

func TestBizErrSegments(t *testing.T) {
	if !(GATEWAY < AUTH && AUTH < WXTOKEN && WXTOKEN < BizStart) {
		t.Fatalf("segments overlap: gateway=%d auth=%d wxtoken=%d biz=%d", GATEWAY, AUTH, WXTOKEN, BizStart)
	}

	cases := []struct {
		bizType BizErrType
		client  bool
		code    int
		want    int
	}{
		{GATEWAY, true, 1, 4001001},
		{GATEWAY, false, 999, 5001999},
		{AUTH, true, 1, 4002001},
		{WXTOKEN, false, 1, 5003001},
		{BizStart, true, 1, 4100001},
		{MATHGAME, true, 1, 4100001},
	}
	for _, c := range cases {
		if got := codeOf(t, c.bizType, c.client, c.code); got != c.want {
			t.Errorf("bizType=%d client=%v code=%d: got %d, want %d", c.bizType, c.client, c.code, got, c.want)
		}
	}
}
