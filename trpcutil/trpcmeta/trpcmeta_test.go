package trpcmeta

import (
	"context"
	"testing"

	"trpc.group/trpc-go/trpc-go/codec"
)

func ctxWithMeta(md codec.MetaData) context.Context {
	ctx, msg := codec.WithNewMessage(context.Background())
	msg.WithServerMetaData(md)
	return ctx
}

func TestGatewayIdentity(t *testing.T) {
	ctx := ctxWithMeta(codec.MetaData{
		"uid":       []byte("10086"),
		"app":       []byte("mathgame"),
		"cid":       []byte("3"),
		"requestId": []byte("req-1"),
		"clientIp":  []byte("1.2.3.4"),
	})

	if got := GetUid(ctx); got != 10086 {
		t.Errorf("GetUid = %d, want 10086", got)
	}
	if got := GetApp(ctx); got != "mathgame" {
		t.Errorf("GetApp = %q, want mathgame", got)
	}
	if got := GetCid(ctx); got != 3 {
		t.Errorf("GetCid = %d, want 3", got)
	}
	if got := GetRequestId(ctx); got != "req-1" {
		t.Errorf("GetRequestId = %q, want req-1", got)
	}
	if got := GetClientIp(ctx); got != "1.2.3.4" {
		t.Errorf("GetClientIp = %q, want 1.2.3.4", got)
	}
}

func TestGatewayIdentityMissing(t *testing.T) {
	for name, ctx := range map[string]context.Context{
		"no message":    context.Background(),
		"empty meta":    ctxWithMeta(codec.MetaData{}),
		"bad uid & cid": ctxWithMeta(codec.MetaData{"uid": []byte("abc"), "cid": []byte("")}),
	} {
		t.Run(name, func(t *testing.T) {
			if got := GetUid(ctx); got != 0 {
				t.Errorf("GetUid = %d, want 0", got)
			}
			if got := GetCid(ctx); got != 0 {
				t.Errorf("GetCid = %d, want 0", got)
			}
			if got := GetApp(ctx); got != "" {
				t.Errorf("GetApp = %q, want empty", got)
			}
			if got := GetRequestId(ctx); got != "" {
				t.Errorf("GetRequestId = %q, want empty", got)
			}
			if got := GetClientIp(ctx); got != "" {
				t.Errorf("GetClientIp = %q, want empty", got)
			}
		})
	}
}
