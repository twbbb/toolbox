package errcode

import (
	"testing"

	"trpc.group/trpc-go/trpc-go/errs"
)

func codeOf(t *testing.T, bizType BizErrType, client bool, code int) int {
	t.Helper()
	be := GetBizErr(bizType)
	if client {
		if err := be.AddBizErrClient(code, "x"); err != nil {
			t.Fatalf("AddBizErrClient: %v", err)
		}
		return int(errs.Code(be.GetBizErrClient(code)))
	}
	if err := be.AddBizErrServer(code, "x"); err != nil {
		t.Fatalf("AddBizErrServer: %v", err)
	}
	return int(errs.Code(be.GetBizErrServer(code)))
}

func TestBizErrSegments(t *testing.T) {
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
		{MATHGAME, true, 1, 4100001},
	}
	for _, c := range cases {
		if got := codeOf(t, c.bizType, c.client, c.code); got != c.want {
			t.Errorf("bizType=%d client=%v code=%d: got %d, want %d", c.bizType, c.client, c.code, got, c.want)
		}
	}
}
