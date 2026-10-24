package ping

import (
	"context"
	"testing"
)

type pingReq struct{}

type pingRsp struct{ Msg string }

type pingService interface {
	Ping(ctx context.Context, req *pingReq) (*pingRsp, error)
}

type bizService struct {
	Default[pingReq, pingRsp]
}

func TestDefaultPing(t *testing.T) {
	var svc pingService = &bizService{}
	rsp, err := svc.Ping(context.Background(), &pingReq{})
	if err != nil {
		t.Fatalf("Ping err = %v", err)
	}
	if rsp == nil {
		t.Fatal("Ping rsp = nil")
	}
}
